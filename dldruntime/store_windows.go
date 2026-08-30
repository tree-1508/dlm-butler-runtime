//go:build windows

package dldruntime

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type dpapiStore struct {
	dir string
}

func OpenPlatformStore(stateDir string) (Store, error) {
	if stateDir == "" {
		base := os.Getenv("LOCALAPPDATA")
		if base == "" {
			return nil, errors.New("LOCALAPPDATA is not set")
		}
		stateDir = filepath.Join(base, "DoujinLibraryManager", "ProviderState", "itchio")
	}
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return nil, fmt.Errorf("create provider state directory: %w", err)
	}
	return &dpapiStore{dir: stateDir}, nil
}

func (s *dpapiStore) profilePath(id int64) string {
	return filepath.Join(s.dir, fmt.Sprintf("profile-%d.dpapi", id))
}

func protectDPAPI(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("refusing to protect empty payload")
	}
	in := windows.DataBlob{Size: uint32(len(src)), Data: &src[0]}
	var out windows.DataBlob
	if err := windows.CryptProtectData(&in, nil, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	res := make([]byte, int(out.Size))
	copy(res, unsafe.Slice(out.Data, int(out.Size)))
	return res, nil
}

func unprotectDPAPI(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("refusing to unprotect empty payload")
	}
	in := windows.DataBlob{Size: uint32(len(src)), Data: &src[0]}
	var out windows.DataBlob
	if err := windows.CryptUnprotectData(&in, nil, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out); err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	res := make([]byte, int(out.Size))
	copy(res, unsafe.Slice(out.Data, int(out.Size)))
	return res, nil
}

func (s *dpapiStore) SaveProfile(p StoredProfile) error {
	plain, err := json.Marshal(p)
	if err != nil {
		return err
	}
	cipher, err := protectDPAPI(plain)
	for i := range plain {
		plain[i] = 0
	}
	if err != nil {
		return fmt.Errorf("protect provider state: %w", err)
	}
	path := s.profilePath(p.ID)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, cipher, 0o600); err != nil {
		return fmt.Errorf("write provider state: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("replace provider state: %w", err)
	}
	return nil
}

func (s *dpapiStore) LoadProfile(id int64) (StoredProfile, error) {
	cipher, err := os.ReadFile(s.profilePath(id))
	if errors.Is(err, os.ErrNotExist) {
		return StoredProfile{}, ErrProfileNotFound
	}
	if err != nil {
		return StoredProfile{}, err
	}
	plain, err := unprotectDPAPI(cipher)
	if err != nil {
		return StoredProfile{}, fmt.Errorf("unprotect provider state: %w", err)
	}
	defer func() {
		for i := range plain {
			plain[i] = 0
		}
	}()
	var p StoredProfile
	if err := json.Unmarshal(plain, &p); err != nil {
		return StoredProfile{}, err
	}
	if p.ID != id || p.APIKey == "" || p.User == nil {
		return StoredProfile{}, errors.New("provider state integrity check failed")
	}
	return p, nil
}

func (s *dpapiStore) ListProfiles() ([]StoredProfile, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	var res []StoredProfile
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "profile-") || !strings.HasSuffix(entry.Name(), ".dpapi") {
			continue
		}
		idText := strings.TrimSuffix(strings.TrimPrefix(entry.Name(), "profile-"), ".dpapi")
		id, err := strconv.ParseInt(idText, 10, 64)
		if err != nil {
			continue
		}
		p, err := s.LoadProfile(id)
		if err == nil {
			res = append(res, p)
		}
	}
	return res, nil
}
