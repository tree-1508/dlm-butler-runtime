//go:build !windows

package dldruntime

// OpenPlatformStore is test/development-only on non-Windows platforms.
// DLD-051 ships only the Windows target; credentials are never persisted
// by this fallback.
func OpenPlatformStore(stateDir string) (Store, error) {
	return NewMemoryStore(), nil
}
