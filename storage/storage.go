package storage

import (
	"errors"
	"log"

	"github.com/dgraph-io/badger"

	"github.com/iorhachovyevhen/dsss/models"
)

func openDB() *badger.DB {
	opt := badger.DefaultOptions("/tmp/badger")

	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func ReadData(hash models.Hash32) ([]byte, error) {
	var data []byte

	db := openDB()
	defer db.Close()

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(hash[:])
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			data = append([]byte{}, val...)

			return nil
		})

		return err
	})

	if err != nil {
		return nil, errors.New("reading finished with error: " + err.Error())
	}

	return data, nil
}

func NewData(hash models.Hash32, data []byte) error {
	db := openDB()
	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set(hash[:], data)
	})

	if err != nil {
		return errors.New("writing data finished with error: " + err.Error())
	}

	return err
}

func UpdateData(oldHash, newHash models.Hash32, data []byte) error {
	err := DeleteData(oldHash)
	if err != nil {
		return err
	}

	err = NewData(newHash, data)

	return err
}

func DeleteData(hash models.Hash32) error {
	db := openDB()
	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		return txn.Delete(hash[:])
	})
	if err != nil {
		return errors.New("deleting finished with error: " + err.Error())
	}

	return err
}
