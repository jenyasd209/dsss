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

func NewStorageWithOptions() *Storage {
	return &Storage{}
}

func NewDefaultStorage() *Storage {
	return &Storage{
		db: openDB(badger.DefaultOptions("/tmp/badger")),
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
	if data.CachedHash().IsEmpty() {
		return errors.New("content not submitted")
	}

	err := s.db.Update(func(txn *badger.Txn) error {
		val, err := data.MarshalBinary()
		if err != nil {
			return err
		}

		key := data.CachedHash()

		return txn.Set(key[:], val)
	})

	if err != nil {
		return errors.New("writing data finished with error: " + err.Error())
	}

	return err
}

func (s *Storage) Read(hash models.Hash32, data models.Data) error {
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(hash[:])
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

func (s *Storage) Delete(hash models.Hash32) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(hash[:])
	})
	if err != nil {
		return errors.New("deleting finished with error: " + err.Error())
	}

	return err
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func ComposeKey(hash32 models.Hash32, dataType models.DataType) (key []byte) {
	prefix := []byte(DataPrefix[dataType])

	key = append(key, prefix...)
	key = append(key, hash32[:]...)

	return key
}
