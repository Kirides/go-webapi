package stores

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Kirides/simpleApi/models"
)

// SQLUserStore Store that enables Saving and Reading Users
type SQLUserStore struct {
	db *sql.DB
}

// NewSQLiteUserStore Creates a new UserStore that uses Sqlite3
func NewSQLiteUserStore(db *sql.DB) (*SQLUserStore, error) {
	store := &SQLUserStore{db: db}

	if err := store.initialize(); err != nil {
		return nil, err
	}

	return store, nil
}
func (s *SQLUserStore) initialize() error {
	if err := s.createTableSQLite(); err != nil {
		return err
	}

	// if _, err := s.db.Exec(`INSERT INTO Users (Username, Hash) VALUES ("abc", "$2a$10$WX3dM2ElqQFOTgtnOzjP9.snX3d0HbfQ1t.1uOWeSUeucz5RB8rEa")`); err != nil {
	// 	return err
	// }

	return nil
}

func (s SQLUserStore) createTableSQLite() error {
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS Users (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Username TEXT NOT NULL,
		Hash TEXT NOT NULL
		)`); err != nil {
		return err
	}
	return nil
}

// GetPage Retrieves a paginated arary of Users
func (s SQLUserStore) GetPage(offset int64, limit int64) ([]models.User, error) {
	rows, err := s.db.Query("SELECT Id, Username, Hash FROM Users LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve Users: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing SQL rows. Error: %v", err)
		}
	}()
	var rowData []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Hash); err != nil {
			return nil, err
		}
		rowData = append(rowData, u)
	}

	return rowData, nil
}

// Get returns a single User by its Id
func (s SQLUserStore) Get(id string) (models.User, error) {
	row := s.db.QueryRow("SELECT Id, Username, Hash FROM Users WHERE Id = ?", id)
	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Hash)
	return u, err
}

// GetByName ...
func (s SQLUserStore) GetByName(name string) (models.User, error) {
	row := s.db.QueryRow("SELECT Id, Username, Hash FROM Users WHERE Username = ?", name)
	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Hash)
	return u, err
}

// Insert adds a user to the store
func (s SQLUserStore) Insert(u models.User) error {
	_, err := s.db.Exec("INSERT INTO Users (Username, Hash) VALUES (?,?)", u.Name, string(u.Hash))
	return err
}

// InsertAll adds all specified users to the store
func (s SQLUserStore) InsertAll(users []models.User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	for _, u := range users {
		if _, err := tx.Exec("INSERT INTO Users Id VALUES (NULL)", u.ID); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Could not commit Insert. Error: %v", err)
		return tx.Rollback()
	}
	return nil
}

// Update updates the specified User
func (s SQLUserStore) Update(u models.User) error {
	_, err := s.db.Exec("UPDATE Users SET xy=Z WHERE Id=?", u.ID)
	return err
}
