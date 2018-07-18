package stores

import (
	"bytes"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"github.com/Kirides/simpleApi/models"
	bolt "github.com/coreos/bbolt"
)

// BoltDBUserStore ...
type BoltDBUserStore struct {
	db *bolt.DB
}

var (
	keyID   = getUInt64Bytes(1) //[]byte("id")
	keyHash = getUInt64Bytes(2) //[]byte("hash")
	keyName = getUInt64Bytes(3) //[]byte("name")
)

// NewBoltDBUserStore Creates a new BoltDB-Based UserStore
func NewBoltDBUserStore(db *bolt.DB) (*BoltDBUserStore, error) {
	store := &BoltDBUserStore{db: db}
	return store, store.initialize()
}

func (s *BoltDBUserStore) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket(boltkeyUsersBucket)
		if err != nil {
			return fmt.Errorf("Could not create bucket 'user'. Error: %v", err)
		}
		for i := 0; i < 1; i++ {
			id, _ := bucket.NextSequence()
			curUserBucket, _ := bucket.CreateBucketIfNotExists(getUInt64Bytes(id))

			if curUserBucket == nil {
				return fmt.Errorf("Could not create bucket for user '%d'", id)
			}
			hash, err := bcrypt.GenerateFromPassword([]byte("1234567890"), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			if err := curUserBucket.Put(keyID, getUInt64Bytes(id)); err != nil {
				return err
			}
			if err := curUserBucket.Put(keyName, []byte("abc"+strconv.Itoa(i))); err != nil {
				return err
			}
			if err := curUserBucket.Put(keyHash, hash); err != nil {
				return err
			}
		}
		return nil
	})
}

// GetPage ...
func (s *BoltDBUserStore) GetPage(offset, limit uint64) ([]models.User, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("Limit cannot be less-or-equal to 0")
	}
	var users []models.User
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltkeyUsersBucket)
		cur := bucket.Cursor()

		firstK, firstV := cur.First()
		if offset > 0 {
			for i := uint64(0); i < offset; i++ {
				firstK, firstV = cur.Next()
				if firstK == nil {
					return nil
				}
			}
		}
		if firstK == nil {
			return nil
		}
		u := make([]models.User, limit)
		rowsFetched := uint64(0)
		for k, _ := firstK, firstV; k != nil && rowsFetched < limit; k, _ = cur.Next() {
			user, _ := userFromBucket(bucket.Bucket(k))

			u[rowsFetched] = user
			rowsFetched++
		}
		users = u[:rowsFetched]
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Could not get page '%d'->'%d'. Error: %v", offset, limit, err)
	}
	return users, nil
}
func userFromBucket(bucket *bolt.Bucket) (models.User, error) {
	user := models.User{
		ID:   getStringFromUInt64Bytes(bucket.Get(keyID)),
		Name: string(bucket.Get(keyName)),
		Hash: bucket.Get(keyHash),
	}
	return user, nil
}

// Get ...
func (s *BoltDBUserStore) Get(id string) (models.User, error) {
	var user models.User
	idAsInt, err := strconv.ParseUint(id, 10, sizeOfUInt64)
	if err != nil {
		return user, err
	}

	if err := s.db.View(func(tx *bolt.Tx) error {
		usrBucket := tx.Bucket(boltkeyUsersBucket)
		idAsBytes := getUInt64Bytes(idAsInt)
		reqUsrBucket := usrBucket.Bucket(idAsBytes)
		if reqUsrBucket == nil {
			return fmt.Errorf("Could not locate user")
		}
		foundUser, err := userFromBucket(reqUsrBucket)
		if err != nil {
			return fmt.Errorf("Could not locate user")
		}
		user = foundUser
		return nil
	}); err != nil {
		return user, fmt.Errorf("Could not find user '%s'. Error: %v", id, err)
	}
	return user, nil
}

// GetByName ...
func (s *BoltDBUserStore) GetByName(name string) (models.User, error) {
	var user models.User

	if err := s.db.View(func(tx *bolt.Tx) error {
		usrBucket := tx.Bucket(boltkeyUsersBucket)
		nameBytes := []byte(name)
		cur := usrBucket.Cursor()
		for k, _ := cur.First(); k != nil; k, _ = cur.Next() {
			if bytes.Equal(usrBucket.Bucket(k).Get(keyName), nameBytes) {
				foundUser, err := userFromBucket(usrBucket.Bucket(k))
				if err != nil {
					return fmt.Errorf("Internal Server Error: User")
				}
				user = foundUser
				return nil
			}
		}
		return fmt.Errorf("Could not locate user")
	}); err != nil {
		return user, fmt.Errorf("Could not find user '%s'. Error: %v", name, err)
	}
	return user, nil
}

// Update ...
func (s *BoltDBUserStore) Update(u models.User) error {
	return nil
}

// InsertAll ...
func (s *BoltDBUserStore) InsertAll(users []models.User) error {
	return nil
}

// Insert ...
func (s *BoltDBUserStore) Insert(user models.User) error {
	return nil
}
