package stores

import (
	"bytes"
	"fmt"

	"github.com/Kirides/simpleApi/models"

	bolt "github.com/coreos/bbolt"
)

// BoltDBTokenStore ...
type BoltDBTokenStore struct {
	db *bolt.DB
}

// NewBoltDBTokenStore Creates a new BoltDB-Based UserStore
func NewBoltDBTokenStore(db *bolt.DB) (*BoltDBTokenStore, error) {
	store := &BoltDBTokenStore{db: db}
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(boltkeyTokenBucket)
		return err
	}); err != nil {
		return nil, err
	}
	return store, nil
}

// Get ...
func (s BoltDBTokenStore) Get(id string) (models.TokenStruct, error) {
	var tokenStruct models.TokenStruct
	if err := s.db.View(func(tx *bolt.Tx) error {
		cur := tx.Bucket(boltkeyTokenBucket).Cursor()
		idBytes := []byte(id)
		k, v := cur.Seek(idBytes)
		if k == nil || !bytes.Equal(k, idBytes) {
			return fmt.Errorf("Token '%s' not found", id)
		}
		unixDate := int64(boltByteOrder.Uint64(v))
		tokenStruct.Token = id
		tokenStruct.Date = unixDate
		return nil
	}); err != nil {
		return tokenStruct, fmt.Errorf("Could not find token '%s'. Error: %v", id, err)
	}
	return tokenStruct, nil
}

// Remove ...
func (s BoltDBTokenStore) Remove(id string) error {
	return nil
}

// Set ...
func (s BoltDBTokenStore) Set(id string, date uint64) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(boltkeyTokenBucket).Put([]byte(id), getUInt64Bytes(date)); err != nil {
			return fmt.Errorf("Could not add Token to bucket. Error: %v", err)
		}
		return nil
	})
	return err
}
