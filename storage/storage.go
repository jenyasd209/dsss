package storage

import (
	"errors"
	"fmt"
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
	Add(data models.Data) ([]byte, error)
	Read(key []byte) (models.Data, error)
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

func (s *Storage) Add(data models.Data) ([]byte, error) {
	key := composeKey(data.ID(), data.Type())

	err := s.db.Update(func(txn *badger.Txn) error {
		val, err := data.MarshalBinary()
		if err != nil {
			return err
		}

		return txn.Set(key, val)
	})

	if err != nil {
		return nil, errors.New("writing data finished with error: " + err.Error())
	}

	return key, nil
}

func (s *Storage) Read(key []byte) (models.Data, error) {
	dt := dataTypeFromKey(key)
	data := models.NewEmptyData(dt)

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
		return nil, errors.New("reading finished with error: " + err.Error())
	}

	return data, nil
}

func (s *Storage) Delete(key []byte) error {
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
	fmt.Println(string(prefix))

	key = append(key, prefix...)
	key = append(key, hash32[:]...)

	return
}

func dataTypeFromKey(id []byte) models.DataType {
	return byteToDataType(id[:len(id)-32])
}

func byteToDataType(b []byte) models.DataType {
	return models.DataType(b[len(b)-1] << (8 * len(b)))
}
