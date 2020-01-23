package storage

import (
	"errors"
	"github.com/dgraph-io/badger"

	"github.com/iorhachovyevhen/dsss/models"
)

type Prefix string

const (
	PrefixSimple Prefix = "simple"
	PrefixJSON   Prefix = "json"
	PrefixAudio  Prefix = "audio"
	PrefixVideo  Prefix = "video"
)

var DataPrefix = map[models.DataType]Prefix{
	models.Simple: PrefixSimple,
	models.JSON:   PrefixJSON,
	models.Audio:  PrefixAudio,
	models.Video:  PrefixVideo,
}

type DataKeeper interface {
	Add(data models.Data) error
	Read(key string, data models.Data) error
	Delete(key string) error

	Close() error
}

type Storage struct {
	db      *badger.DB
	options badger.Options
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

func (s *Storage) Add(data models.Data) error {
	if !data.IsCorrect() {
		return errors.New("no correct")
	}

	key := composeKey(data.CachedHash(), data.Type())

	err := s.db.Update(func(txn *badger.Txn) error {
		val, err := data.MarshalBinary()
		if err != nil {
			return err
		}

		return txn.Set(key, val)
	})

	if err != nil {
		return errors.New("writing data finished with error: " + err.Error())
	}

	return nil
}

func (s *Storage) Read(hash models.Hash32, data models.Data) error {
	key := composeKey(hash, data.Type())

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
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

func (s *Storage) Delete(hash models.Hash32, dataType models.DataType) error {
	key := composeKey(hash, dataType)

	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
	if err != nil {
		return errors.New("deleting finished with error: " + err.Error())
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func composeKey(hash32 models.Hash32, dataType models.DataType) (key []byte) {
	prefix := []byte(DataPrefix[dataType])

	key = append(key, prefix...)
	key = append(key, hash32[:]...)

	return key
}
