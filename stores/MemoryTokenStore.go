package stores

import (
	"fmt"

	"github.com/Kirides/simpleApi/models"
)

// MemoryTokenStore ...
type MemoryTokenStore struct {
	db map[string]int64
}

// NewMemoryTokenStore Creates a new BoltDB-Based UserStore
func NewMemoryTokenStore(db map[string]int64) (*MemoryTokenStore, error) {
	store := &MemoryTokenStore{db: db}
	return store, nil
}

// Get ...
func (s MemoryTokenStore) Get(id string) (models.TokenStruct, error) {
	tokenStruct := models.TokenStruct{}
	if v, ok := s.db[id]; ok {
		tokenStruct.Token = id
		tokenStruct.Date = v
	} else {
		return tokenStruct, fmt.Errorf("Could not locate token")
	}
	return tokenStruct, nil
}

// Remove ...
func (s *MemoryTokenStore) Remove(id string) error {
	return nil
}

// Set ...
func (s *MemoryTokenStore) Set(id string, date int64) error {
	s.db[id] = date
	return nil
}
