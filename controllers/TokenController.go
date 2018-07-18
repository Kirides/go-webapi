package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Kirides/simpleApi/helpers"
	"github.com/Kirides/simpleApi/models"
	"github.com/Kirides/simpleApi/stores"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// ApplicationClaims ...
type ApplicationClaims struct {
	*jwt.StandardClaims
	Scope    string `json:"scope,omitempty"`
	Username string `json:"username,omitempty"`
}

// TokenController ...
type TokenController struct {
	jwtTokenSecret       []byte
	DefaultTokenLifetime time.Duration
	UserStore            stores.UserStore
}

// ErrInvalidCredentials ...
var ErrInvalidCredentials = errors.New("Invalid Credentials")

// NewTokenController creates a default TokenController with the SecretKey = "Secret" and defaultTokenLifetime = time.Hour
func NewTokenController(secret []byte, userStore stores.UserStore) *TokenController {
	tc := &TokenController{
		jwtTokenSecret:       secret,
		DefaultTokenLifetime: time.Minute * 10,
		UserStore:            userStore,
	}
	return tc
}

// HandleTokenAPI registers the /users endpoint onto the provided router
func (tc *TokenController) HandleTokenAPI(r *mux.Router) {
	r.Path("/token").Methods(http.MethodPost).HandlerFunc(tc.jwtTokenHandler)
	log.Println("registered token-endpoint")
}

// JwtTokenKeyFunc Function that provides the Signing-Key to validate the Token
func (tc TokenController) JwtTokenKeyFunc(tkn *jwt.Token) (interface{}, error) {
	return tc.jwtTokenSecret, nil
}

// SetJwtSigningKey Changes the key used for signing the JWT Tokens
func (tc *TokenController) SetJwtSigningKey(key []byte) {
	tc.jwtTokenSecret = key
}

func (tc *TokenController) jwtTokenHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	usr, err := tc.validateTokenRequest(r.Form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	tokenID, err := helpers.UUIDv4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tokenTime := time.Now()
	claims := &ApplicationClaims{
		StandardClaims: &jwt.StandardClaims{
			IssuedAt:  tokenTime.Unix(),
			ExpiresAt: tokenTime.Add(tc.DefaultTokenLifetime).Unix(),
			Issuer:    "jwt-host",
			Subject:   usr.ID,
			Id:        tokenID,
		},
		Username: usr.Name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tc.jwtTokenSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tokenResponse, err := json.Marshal(map[string]interface{}{
		"token_type":   "Bearer",
		"access_token": tokenString,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Write(tokenResponse)
}

func (tc *TokenController) validateTokenRequest(v url.Values) (models.User, error) {
	switch v.Get("grant_type") {
	case "password":
		return tc.validateResourceTokenRequest(v)
	}
	return models.User{}, fmt.Errorf("Invalid validation type '%s'", v.Get("grant_type"))
}

func (tc *TokenController) validateResourceTokenRequest(v url.Values) (models.User, error) {
	usr, err := tc.UserStore.GetByName(v.Get("username"))
	if err != nil {
		return models.User{}, ErrInvalidCredentials
	}
	pass := v.Get("password")

	if bcrypt.CompareHashAndPassword(usr.Hash, []byte(pass)) == nil {
		return usr, nil
	}
	return models.User{}, ErrInvalidCredentials
}
