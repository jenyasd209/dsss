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

func ReadData(key string, data models.Data) error {
	db := openDB()
	defer db.Close()

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			return data.UnmarshalBinary(val)
		})

		return err
	})

	if err != nil {
		return errors.New("reading finished with error: " + err.Error())
	}

	return nil
}

func NewData(data models.Data) error {
	db := openDB()
	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		val, err := data.MarshalBinary()
		if err != nil {
			return err
		}

		key := data.Hash().String()

		return txn.Set([]byte(key), val)
	})

	if err != nil {
		return errors.New("writing data finished with error: " + err.Error())
	}

	return err
}

func UpdateData(oldKey string, data models.Data) error {
	err := DeleteData(oldKey)
	if err != nil {
		return err
	}

	return NewData(data)
}

func DeleteData(key string) error {
	db := openDB()
	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	if err != nil {
		return errors.New("deleting finished with error: " + err.Error())
	}

	return err
}
