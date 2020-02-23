package storage

import (
	"github.com/dgraph-io/badger"
	"github.com/iorhachovyevhen/dsss/models"
	"github.com/pkg/errors"
)

var (
	ErrIDNotFound    = errors.New("ID not found")
	ErrAlreadyUsedID = errors.New("ID already is used")
	ErrInvalidID     = errors.New("Invalid ID")
)

type DataKeeper interface {
	Add(data models.Data) (models.ID, error)
	Read(key []byte, data models.Data) error
	Delete(key []byte) error

	Close() error
}

type Storage struct {
	db *badger.DB
}

func NewOptions() badger.Options {
	return badger.DefaultOptions("")
}

func NewStorageWithOptions(opts badger.Options) *Storage {
	return &Storage{
		db: openDB(opts),
	}
}

func NewDefaultStorage(savePath string) *Storage {
	return &Storage{
		db: openDB(badger.DefaultOptions(savePath)),
	}
}

func openDB(opt badger.Options) *badger.DB {
	db, err := badger.Open(opt)
	if err != nil {
		panic(err)
	}

	return db
}

func (s *Storage) Add(data models.Data) (models.ID, error) {
	if s.Exist(data.ID()) {
		return nil, ErrAlreadyUsedID
	}

	err := s.db.Update(func(txn *badger.Txn) error {
		val, err := data.MarshalBinary()
		if err != nil {
			return err
		}

		return txn.Set(data.ID(), val)
	})

	if err != nil {
		return nil, errors.New("writing data finished with error: " + err.Error())
	}

	return data.ID(), nil
}

func (s *Storage) Read(key []byte, data models.Data) error {
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		if item.IsDeletedOrExpired() {
			return errors.New("Item is deleted or expired")
		}

		err = item.Value(func(val []byte) error {
			return data.UnmarshalBinary(val)
		})

		return err
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return ErrIDNotFound
		}
		return errors.New("reading finished with error: " + err.Error())
	}

	return nil
}

func (s *Storage) Delete(key []byte) error {
	if !s.Exist(key) {
		return ErrIDNotFound
	}

	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
	if err != nil {
		return errors.New("deleting finished with error: " + err.Error())
	}

	return nil
}

func (s *Storage) Exist(key []byte) bool {
	txn := s.db.NewTransaction(false)
	defer txn.Discard()

	item, err := txn.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return false
		}
		return true
	}

	if item.IsDeletedOrExpired() {
		return false
	}

	return true
}

func (s *Storage) Close() error {
	return s.db.Close()
}
