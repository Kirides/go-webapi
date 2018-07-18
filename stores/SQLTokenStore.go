package stores

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Kirides/simpleApi/models"
)

// SQLTokenStore Store that enables Saving and Reading Users
type SQLTokenStore struct {
	db *sql.DB
}

// NewSQLTokenStore Creates a new UserStore that uses Sqlite3
func NewSQLTokenStore(db *sql.DB) (*SQLTokenStore, error) {
	store := &SQLTokenStore{db: db}

	return store, store.initialize()
}
func (s *SQLTokenStore) initialize() error {
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS Tokens (
Id INTEGER PRIMARY KEY AUTOINCREMENT,
TokenId TEXT NOT NULL,
Date INTEGER NOT NULL
)`); err != nil {
		log.Printf("Could not initialize Table. Error: %v", err)
	}
	return nil
}

// Get ...
func (s SQLTokenStore) Get(id string) (models.TokenStruct, error) {
	var tokenStruct models.TokenStruct
	row := s.db.QueryRow("SELECT TokenId, Date FROM Tokens WHERE TokenId = ? LIMIT 1", id)
	if err := row.Scan(&tokenStruct.Token, &tokenStruct.Date); err != nil {
		return tokenStruct, fmt.Errorf("Could not find token '%s'. Error: %v", id, err)
	}

	return tokenStruct, nil
}

// Remove ...
func (s SQLTokenStore) Remove(id string) error {
	return nil
}

// Set ...
func (s SQLTokenStore) Set(id string, date int64) error {
	r := s.db.QueryRow("SELECT Id FROM Tokens WHERE TokenId = ? LIMIT 1", id)
	exist := false
	if err := r.Scan(1); err == nil {
		exist = true
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if exist {
		r, err := tx.Exec("UPDATE Tokens SET Date = ? WHERE TokenId = ?", date, id)
		if err != nil {
			return fmt.Errorf("Error executing SQL. Error: %v", err)
		}
		if _, err := r.RowsAffected(); err != nil {
			return fmt.Errorf("Error could not update token. Error: %v", err)
		}
	} else {
		r, err := tx.Exec("INSERT INTO Tokens (TokenId, Date) VALUES (?, ?)", id, date)
		if err != nil {
			return fmt.Errorf("Error executing SQL. Error: %v", err)
		}

		if _, err := r.LastInsertId(); err != nil {
			return fmt.Errorf("Error could not insert token into database. Error: %v", err)
		}
	}
	if err := tx.Commit(); err != nil {
		log.Printf("Could not commit changes. Error: %v", err)
		if err := tx.Rollback(); err != nil {
			log.Printf("WARNING: Could not rollback changes!. Error: %v", err)
		}
	}
	return nil
}
