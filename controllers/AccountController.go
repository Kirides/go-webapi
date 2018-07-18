package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/Kirides/simpleApi/models"
	"github.com/Kirides/simpleApi/stores"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const authCookie = "auth"

type userLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Remember bool   `json:"remember_me"`
}
type userRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// AccountController ...
type AccountController struct {
	userStore  stores.UserStore
	rxUsername *regexp.Regexp
	rxEmail    *regexp.Regexp
}

// NewAccountController ...
func NewAccountController(us stores.UserStore) *AccountController {
	return &AccountController{
		userStore:  us,
		rxUsername: regexp.MustCompile("^[A-Za-z0-9]+(?:[_-][A-Za-z0-9]+)*$"),
		rxEmail:    regexp.MustCompile(`^(?:(?:[^<>()[\]\\.,;:\s@"]+(?:\.[^<>()[\]\\.,;:\s@"]+)*)|(?:".+"))@(?:(?:\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(?:(?:[a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`),
	}
}

// HandeAccountAPI ...
func (ac *AccountController) HandeAccountAPI(r *mux.Router) {
	r.Path("/register").Methods(http.MethodPost).HandlerFunc(ac.handleRegister)
}

func (ac *AccountController) handleRegister(w http.ResponseWriter, r *http.Request) {
	registerRequest := userRegister{}

	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if !ac.rxUsername.MatchString(registerRequest.Username) {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}
	if !ac.rxEmail.MatchString(registerRequest.Email) {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}
	if _, err := ac.userStore.GetByName(registerRequest.Username); err == nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ac.userStore.Insert(models.User{
		Name: registerRequest.Username,
		Hash: passHash,
	})
}
