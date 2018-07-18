package services

import (
	"fmt"

	"github.com/Kirides/simpleApi/models"
	"github.com/Kirides/simpleApi/stores"
	"golang.org/x/crypto/bcrypt"
)

// SignInManager ...
type SignInManager struct {
	us stores.UserStore
}

// NewSignInManager ...
func NewSignInManager(us stores.UserStore) (*SignInManager, error) {
	if us == nil {
		return nil, fmt.Errorf("No valid userstore was provided")
	}
	return &SignInManager{
		us: us,
	}, nil
}

// LogIn ...
func (sim *SignInManager) LogIn(name string, password []byte) (models.User, error) {
	user, err := sim.us.GetByName(name)
	if err != nil {
		return models.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword(user.Hash, password); err != nil {
		return models.User{}, err
	}
	return user, nil
}

// LogOut ...
func (sim *SignInManager) LogOut(u models.User) (bool, error) {

	return true, nil
}
