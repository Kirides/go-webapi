package stores

import "github.com/Kirides/simpleApi/models"

// UserStore contains the logic to persist users
type UserStore interface {
	GetPage(offset int64, limit int64) ([]models.User, error)
	Get(id string) (models.User, error)
	GetByName(name string) (models.User, error)
	Update(u models.User) error
	InsertAll(users []models.User) error
	Insert(users models.User) error
}

// TokenStore ...
type TokenStore interface {
	Get(id string) (models.TokenStruct, error)
	Set(id string, date int64) error
	Remove(id string) error
}
