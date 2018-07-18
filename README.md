# go-webapi
Basic implementation of a go webapi featuring separated packages and middleware

This repository functions as a Template for very basic WebAPIs that need authentication and a web-frontend.

It is built loosely coupled, allowing replacing some of its functionality by different implementations
Controllers, are handlers for specific endpoints, which reside in the `controllers`-package

```golang
// UserStore allows to persist and retrieve users
type UserStore interface {
	GetPage(offset int64, limit int64) ([]models.User, error)
	Get(id string) (models.User, error)
	GetByName(name string) (models.User, error)
	Update(u models.User) error
	InsertAll(users []models.User) error
	Insert(users models.User) error
}

// TokenStore allows to retrieving, setting and removing of validation tokens
type TokenStore interface {
	Get(id string) (models.TokenStruct, error)
	Set(id string, date int64) error
	Remove(id string) error
}
```

it has `UserStore` and `TokenStore`implementations for both `BoltDb` (native go) and `SQLite` (needs gcc, not portable)

It has a very basic, but nice looking Frontend, powered by VueJs and Bootstrap.
It has built in client-side and server-side validation for user registration
currently missing is a "password forgotten"-feature
