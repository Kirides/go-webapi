package models

// User type for UsersController
type User struct {
	ID   string
	Name string
	Hash []byte
}
