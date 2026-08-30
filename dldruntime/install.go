package dldruntime

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func supportedInstallerType(filename string) (string, bool) {
	name := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(name, ".zip"):
		return "archive/zip", true
	case strings.HasSuffix(name, ".tar.gz"), strings.HasSuffix(name, ".tgz"):
		return "archive/tar+gzip", true
	case strings.HasSuffix(name, ".tar"):
		return "archive/tar", true
	case strings.HasSuffix(name, ".7z"), strings.HasSuffix(name, ".rar"), strings.HasSuffix(name, ".xz"), strings.HasSuffix(name, ".tar.xz"):
		return "unsupported", false
	default:
		return "naked", true
	}
}

func requireContained(root, child string) error {
	if root == "" {
		return errors.New("root is not configured")
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	absChild, err := filepath.Abs(child)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(absRoot, absChild)
	if err != nil {
		return err
	}
	if rel == "." {
		return nil
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return errors.New("path escapes configured root")
	}
	return nil
}

func (r *Runtime) performInstall(ctx context.Context, p InstallQueueParams) (string, error) {
	if err := os.MkdirAll(p.StagingFolder, 0o700); err != nil {
		return "", err
	}
	if err := os.MkdirAll(p.InstallFolder, 0o700); err != nil {
		return "", err
	}
	profile, err := r.store.LoadProfile(p.ProfileID)
	if err != nil {
		return "", err
	}
	downloadURL := r.apiBaseURL + "/uploads/" + fmt.Sprintf("%d", p.Upload.ID) + "/download?api_key=" + url.QueryEscape(profile.APIKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", err
	}
	res, err := r.httpClient.Do(req)
	if err != nil {
		return "", redactError(err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("download failed with HTTP status %d", res.StatusCode)
	}
	name := filepath.Base(p.Upload.Filename)
	if name == "." || name == string(filepath.Separator) || name == "" {
		name = "upload.bin"
	}
	staged := filepath.Join(p.StagingFolder, name+".part")
	f, err := os.OpenFile(staged, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return "", err
	}
	_, copyErr := io.Copy(f, res.Body)
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(staged)
		return "", copyErr
	}
	if closeErr != nil {
		_ = os.Remove(staged)
		return "", closeErr
	}
	defer os.Remove(staged)

	typeName, ok := supportedInstallerType(p.Upload.Filename)
	if !ok {
		return "", errors.New("unsupported upload format")
	}
	switch typeName {
	case "archive/zip":
		if err := extractZip(staged, p.InstallFolder); err != nil {
			return "", err
		}
	case "archive/tar":
		f, err := os.Open(staged)
		if err != nil {
			return "", err
		}
		err = extractTar(f, p.InstallFolder)
		_ = f.Close()
		if err != nil {
			return "", err
		}
	case "archive/tar+gzip":
		f, err := os.Open(staged)
		if err != nil {
			return "", err
		}
		gz, err := gzip.NewReader(f)
		if err != nil {
			_ = f.Close()
			return "", err
		}
		err = extractTar(gz, p.InstallFolder)
		_ = gz.Close()
		_ = f.Close()
		if err != nil {
			return "", err
		}
	default:
		target := filepath.Join(p.InstallFolder, name)
		if err := copyFile(staged, target); err != nil {
			return "", err
		}
	}
	return typeName, nil
}

func safeArchiveTarget(root, name string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(name))
	if clean == "." || clean == "" || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", errors.New("archive path escapes destination")
	}
	target := filepath.Join(root, clean)
	if err := requireContained(root, target); err != nil {
		return "", err
	}
	return target, nil
}

func extractZip(path, dest string) error {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, zf := range zr.File {
		if zf.Mode()&os.ModeSymlink != 0 {
			return errors.New("symlinks are not allowed in DLD archives")
		}
		if strings.Contains(zf.Name, "..") {
			return errors.New("archive path contains prohibited traversal element")
		}
		target, err := safeArchiveTarget(dest, zf.Name)
		if err != nil {
			return err
		}
		if zf.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o700); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return err
		}
		rc, err := zf.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
		if err != nil {
			_ = rc.Close()
			return err
		}
		_, copyErr := io.Copy(out, rc)
		closeOutErr := out.Close()
		closeInErr := rc.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeOutErr != nil {
			return closeOutErr
		}
		if closeInErr != nil {
			return closeInErr
		}
	}
	return nil
}

func extractTar(reader io.Reader, dest string) error {
	tr := tar.NewReader(reader)
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if strings.Contains(h.Name, "..") {
			return errors.New("archive path contains prohibited traversal element")
		}
		target, err := safeArchiveTarget(dest, h.Name)
		if err != nil {
			return err
		}
		switch h.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o700); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
			if err != nil {
				return err
			}
			_, copyErr := io.Copy(out, tr)
			closeErr := out.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		default:
			return errors.New("archive contains unsupported link or special entry")
		}
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
