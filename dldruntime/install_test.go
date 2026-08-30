package dldruntime

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestRequireContained(t *testing.T) {
	root := t.TempDir()
	if err := requireContained(root, filepath.Join(root, "game")); err != nil {
		t.Fatal(err)
	}
	if err := requireContained(root, filepath.Join(root, "..", "escape")); err == nil {
		t.Fatal("expected escape to be rejected")
	}
}

func TestZipTraversalRejected(t *testing.T) {
	root := t.TempDir()
	zipPath := filepath.Join(t.TempDir(), "bad.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	w, err := zw.Create("../escape.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.Write([]byte("bad"))
	_ = zw.Close()
	_ = f.Close()
	if err := extractZip(zipPath, root); err == nil {
		t.Fatal("expected traversal archive to be rejected")
	}
}

func TestUnsupportedArchiveFormatsFailClosed(t *testing.T) {
	for _, name := range []string{"game.7z", "game.rar", "game.xz", "game.tar.xz"} {
		if _, ok := supportedInstallerType(name); ok {
			t.Fatalf("%s must fail closed", name)
		}
	}
}
