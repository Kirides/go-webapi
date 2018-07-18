package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/Kirides/simpleApi/models"
	"github.com/Kirides/simpleApi/sqlite3"

	"github.com/Kirides/simpleApi/controllers"
	"github.com/Kirides/simpleApi/stores"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// var inMemoryDb = "file::memory:?mode=memory&cache=shared"
var (
	inMemoryDb      = "file:demo.db?cache=shared&mode=rwc&_busy_timeout=5000"
	tokenController *controllers.TokenController
	usersController *controllers.UsersController
	tokenStore      stores.TokenStore
	tokenSecret     = []byte("Secret")
)
var srv = &http.Server{
	Addr:              "127.0.0.1:5001",
	IdleTimeout:       15 * time.Second,
	ReadTimeout:       30 * time.Second,
	ReadHeaderTimeout: 15 * time.Second,
	WriteTimeout:      15 * time.Second,
}

func main() {
	r := mux.NewRouter()

	db, err := sqlite3.Open("sqlite3", "file:api.db?cache=shared&mode=rwc&_busy_timeout=20000")
	// boltdb, err := bolt.Open("boltDb.db", 0600, &bolt.Options{Timeout: time.Second * 30})
	if err != nil {
		log.Fatalf("Could not initialize Database. Error: %v", err)
	}
	// defer boltdb.Close()
	defer db.Close()
	// boltUserStore, _ := stores.NewBoltDBUserStore(db)
	// sqlDb, _ := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/GoAPI")
	userStore, err := stores.NewSQLiteUserStore(db.DB)
	if err != nil {
		panic(err)
	}
	// userStore := stores.NewMemoryUserStore()

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(authentication(jwtAuthentication))
	apiRouter.Use(accessControlAllowOrigin)

	usersController = controllers.NewUsersController(userStore)
	usersController.HandleUsersAPI(apiRouter)

	tokenController = controllers.NewTokenController(tokenSecret, userStore)
	tokenController.SetJwtSigningKey([]byte("MyNewTopSecretSecret"))
	tokenController.HandleTokenAPI(r.PathPrefix("/api").Subrouter())

	// boltTokenStore, _ := stores.NewBoltDBTokenStore(boltdb)
	// tokenStore = boltTokenStore
	// sqlTokenStore, err := stores.NewSQLTokenStore(db.DB)
	// tokenStore = sqlTokenStore
	accountController := controllers.NewAccountController(userStore)
	accountController.HandeAccountAPI(r.PathPrefix("/account").Subrouter())

	r.PathPrefix("/").Methods(http.MethodGet).Handler(http.StripPrefix("/", http.FileServer(http.Dir("./wwwroot"))))
	srv.Handler = r

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Println("Starting server on", srv.Addr)
		wg.Done()
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalln(err)
			}
		}
	}()
	wg.Wait()
	<-interrupt
	handleShutdown()
}

func handleShutdown() {
	log.Println("Started shutdown sequence (this might take a while)")
	srv.Shutdown(context.Background())
	log.Println("Shutdown completed")
}

func accessControlAllowOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

type authenticationFunc func(*http.Request) (context.Context, error)

func authentication(auths ...authenticationFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, a := range auths {
				if ctx, err := a(r); err == nil {
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
		})
	}
}

func jwtAuthentication(r *http.Request) (context.Context, error) {
	const authScheme = "Bearer "
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return r.Context(), fmt.Errorf("No Authorization header found")
	}
	if !strings.HasPrefix(authHeader, authScheme) {
		return r.Context(), fmt.Errorf("Invalid Authorization Scheme. Required: " + authScheme)
	}
	authToken := authHeader[len(authScheme):]
	token, err := jwt.Parse(authToken, tokenController.JwtTokenKeyFunc)
	if err != nil || !token.Valid {
		return r.Context(), fmt.Errorf("Invalid Authorization Token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return r.Context(), fmt.Errorf("Invalid Authorization Token")
	}
	c := context.WithValue(r.Context(), models.KeyTokenUsername, claims["username"])
	// -----------------------------
	// --- Token Revocation Demo ---
	// -----------------------------
	// revoked, tokenID, err := isTokenRevoked(token, tokenStore)
	// if err != nil {
	// 	return fmt.Errorf("Invalid Authorization Token")
	// }
	// if revoked {
	// 	return fmt.Errorf("Token revoked")
	// }
	// if err := tokenStore.Set(tokenID, int64(token.Claims.(jwt.MapClaims)["exp"].(float64))); err != nil {
	// 	log.Println(err)
	// }
	return c, nil
}

func isTokenRevoked(token *jwt.Token, tokenStore stores.TokenStore) (bool, string, error) {
	claims := token.Claims.(jwt.MapClaims)
	tokenID, ok := claims["jti"].(string)
	if !ok {
		return false, "", fmt.Errorf("Invalid Authorization Token")
	}

	if rejToken, err := tokenStore.Get(tokenID); err == nil {
		if rejToken.Date > int64(time.Now().Unix()) {
			return true, tokenID, nil
		}
	}
	return false, tokenID, nil
}
