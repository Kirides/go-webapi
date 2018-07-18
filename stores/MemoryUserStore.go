package stores

import (
	"log"
	"sync"

	"github.com/Kirides/simpleApi/models"
)

// InMemoryUserStore ...
type InMemoryUserStore struct {
	users []models.User
	m     *sync.Mutex
}

// NewMemoryUserStore Creates a new In-Memory UserStore
func NewMemoryUserStore() *InMemoryUserStore {
	store := &InMemoryUserStore{
		m: new(sync.Mutex),
	}
	if err := store.Insert(models.User{ID: "1", Name: "abc", Hash: []byte("$2a$10$WX3dM2ElqQFOTgtnOzjP9.snX3d0HbfQ1t.1uOWeSUeucz5RB8rEa")}); err != nil {
		log.Printf("Error inserting Demo data. Error: %v", err)
	}
	return store
}

// GetPage ...
func (s *InMemoryUserStore) GetPage(offset int64, limit int64) ([]models.User, error) {
	total := offset + limit
	if total > int64(len(s.users)) {
		return s.users[offset:], nil
	}
	return s.users[offset : offset+limit], nil
}

// Get ...
func (s *InMemoryUserStore) Get(id string) (models.User, error) {
	s.m.Lock()
	for _, v := range s.users {
		if v.ID == id {
			s.m.Unlock()
			return v, nil
		}
	}
	s.m.Unlock()
	return models.User{}, nil
}

// GetByName ...
func (s *InMemoryUserStore) GetByName(name string) (models.User, error) {
	s.m.Lock()
	for _, v := range s.users {
		if v.Name == name {
			s.m.Unlock()
			return v, nil
		}
	}
	s.m.Unlock()
	return models.User{}, nil
}

// Update ...
func (s *InMemoryUserStore) Update(u models.User) error {
	return nil
}

// InsertAll ...
func (s *InMemoryUserStore) InsertAll(users []models.User) error {
	s.m.Lock()
	s.users = append(s.users, users...)
	s.m.Unlock()
	return nil
}

// Insert ...
func (s *InMemoryUserStore) Insert(user models.User) error {
	s.m.Lock()
	s.users = append(s.users, user)
	s.m.Unlock()
	return nil
}
