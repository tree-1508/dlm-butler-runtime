package dldruntime

import (
	"errors"
	"sync"
)

var ErrProfileNotFound = errors.New("profile not found")

type StoredProfile struct {
	ID     int64  `json:"id"`
	APIKey string `json:"apiKey"`
	User   *User  `json:"user"`
}

type Store interface {
	SaveProfile(StoredProfile) error
	LoadProfile(id int64) (StoredProfile, error)
	ListProfiles() ([]StoredProfile, error)
}

type MemoryStore struct {
	mu       sync.Mutex
	profiles map[int64]StoredProfile
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{profiles: make(map[int64]StoredProfile)} }

func (s *MemoryStore) SaveProfile(p StoredProfile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.profiles[p.ID] = p
	return nil
}

func (s *MemoryStore) LoadProfile(id int64) (StoredProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.profiles[id]
	if !ok {
		return StoredProfile{}, ErrProfileNotFound
	}
	return p, nil
}

func (s *MemoryStore) ListProfiles() ([]StoredProfile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]StoredProfile, 0, len(s.profiles))
	for _, p := range s.profiles {
		res = append(res, p)
	}
	return res, nil
}
