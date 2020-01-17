package storage

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger"

	"dsss/models"
)

func openDB() *badger.DB{
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func ReadData (hash models.Hash32) ([]byte, error){
	var data []byte

	db := openDB()
	defer db.Close()

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(hash[:])
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			fmt.Printf("The answer is: %s\n", val)

			data = append([]byte{}, val...)

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

func WriteData (hash models.Hash32, data []byte) error{
	db := openDB()
	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Set(hash[:], data)
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func DeleteData (hash models.Hash32) error{
	db := openDB()
	defer db.Close()

	err := db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(hash[:])

		return err
	})

	if err != nil {
		return err
	}

	return nil
}