package sqlite3

import (
	"context"
	"database/sql"
	"sync"
)

var mapMtx sync.Mutex
var mutexConStrings = make(map[string]*sync.Mutex)

type Sqlite struct {
	*sql.DB
	mutex *sync.Mutex
}

type SqliteTx struct {
	*sql.Tx
	mutex *sync.Mutex
}

func Open(driver string, connectionString string) (*Sqlite, error) {
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}

	if err := applyPragmas(db, "journal_mode=WAL", `synchronous="1"`); err != nil {
		db.Close()
		return nil, err
	}
	sqlite := Sqlite{}
	sqlite.DB = db
	mapMtx.Lock()
	mutex, ok := mutexConStrings[connectionString]
	if !ok {
		mutex = &sync.Mutex{}
		mutexConStrings[connectionString] = mutex
	}
	mapMtx.Unlock()
	sqlite.mutex = mutex
	return &sqlite, nil
}
func applyPragmas(db *sql.DB, pragmas ...string) error {
	for _, p := range pragmas {
		if _, err := db.Exec("PRAGMA " + p); err != nil {
			return err
		}
	}
	return nil
}
func (s Sqlite) Exec(query string, args ...interface{}) (sql.Result, error) {
	s.mutex.Lock()
	res, err := s.DB.Exec(query, args...)
	s.mutex.Unlock()
	return res, err
}

func (s Sqlite) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	s.mutex.Lock()
	res, err := s.DB.ExecContext(ctx, query, args...)
	s.mutex.Unlock()
	return res, err
}

func (s Sqlite) Begin() (*SqliteTx, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	sqliteTx := SqliteTx{}
	sqliteTx.Tx = tx
	sqliteTx.mutex = s.mutex
	sqliteTx.mutex.Lock()
	return &sqliteTx, nil
}

func (sTx SqliteTx) Rollback() error {
	err := sTx.Tx.Rollback()
	sTx.mutex.Unlock()
	return err
}

func (sTx SqliteTx) Commit() error {
	err := sTx.Tx.Commit()
	sTx.mutex.Unlock()
	return err
}
