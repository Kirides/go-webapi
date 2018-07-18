package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Kirides/simpleApi/stores"
	"github.com/gorilla/mux"
)

// UsersController ...
type UsersController struct {
	store            stores.UserStore
	MaxUsersReturned int64
}

// NewUsersController ...
func NewUsersController(store stores.UserStore) *UsersController {
	return &UsersController{
		store:            store,
		MaxUsersReturned: 100,
	}
}

// HandleUsersAPI registers the /users endpoint onto the provided router
func (uc *UsersController) HandleUsersAPI(r *mux.Router) {
	r.Path("/users").Methods(http.MethodGet).Handler(uc.handleUsers())
	r.Path("/users/{id:[0-9]+}").Methods(http.MethodGet).Handler(uc.handleUserByID())
	log.Println("registered users-endpoint")
}

func (uc *UsersController) handleUserByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user, err := uc.store.Get(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		b, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Could not format result", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		w.Write(b)
	})
}

func (uc *UsersController) handleUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset, err := getOffset(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		limit, err := getLimit(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if limit > uc.MaxUsersReturned {
			limit = uc.MaxUsersReturned
		}
		users, err := uc.store.GetPage(offset, limit)
		if err != nil {
			http.Error(w, "Could not retrieve result", http.StatusBadRequest)
			return
		}
		if len(users) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		b, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "Could not format result", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", strconv.Itoa(len(b)))
		w.Write(b)
	})
}

func getOffset(r *http.Request) (int64, error) {
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}
	if offset < 0 {
		return 0, fmt.Errorf("offset must be greater than or equal to '0'")
	}
	return offset, nil
}

func getLimit(r *http.Request) (int64, error) {
	limitInt, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limitInt = 20
	}
	if limitInt < 0 {
		return 0, fmt.Errorf("limit must be greater than or equal to '0'")
	}
	return limitInt, nil
}
