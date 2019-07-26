package bolt

import (
	"github.com/boltdb/bolt"
)

type Storage struct {
	db *bolt.DB
}

func NewStorage(DBName string) *Storage {
	db, err := bolt.Open(DBName+".db", 0600, nil)
	if err != nil {
		panic(err)
	}

	return &Storage{db: db}
}

func (b *Storage) Close() {
	err := b.db.Close()
	if err != nil {
		panic(err)
	}
}

func (b *Storage) GetAll(collection string) (map[string]string, error) {
	values := make(map[string]string)
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(collection))
		err := bucket.ForEach(func(k []byte, v []byte) error {
			values[string(k)] = string(v)

			return nil
		})

		return err
	})

	return values, err
}

func (b *Storage) GetAllRaw(collection string) (map[string][]byte, error) {
	values := make(map[string][]byte)
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(collection))
		err := bucket.ForEach(func(k []byte, v []byte) error {
			values[string(k)] = v

			return nil
		})

		return err
	})

	return values, err
}

func (b *Storage) Put(collection string, key string, value []byte) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(collection))

		err := bucket.Put([]byte(key), value)

		return err
	})

	return err
}

func (b *Storage) CreateCollectionIfNotExist(name string) {
	err := b.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))

		return err
	})

	if err != nil {
		panic(err)
	}
}
