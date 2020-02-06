package storage

import (
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"

	"github.com/iorhachovyevhen/dsss/models"
)

type Prefix string

const (
	PrefixSimple Prefix = "simple"
	PrefixJSON   Prefix = "json"
	PrefixAudio  Prefix = "audio"
	PrefixVideo  Prefix = "video"
)

var DataPrefixMap = map[models.DataType]Prefix{
	models.Simple: PrefixSimple,
	models.JSON:   PrefixJSON,
	models.Audio:  PrefixAudio,
	models.Video:  PrefixVideo,
}

var DataTypeMap = map[Prefix]models.DataType{
	PrefixSimple: models.Simple,
	PrefixJSON:   models.JSON,
	PrefixAudio:  models.Audio,
	PrefixVideo:  models.Video,
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
	dt, err := DataTypeFromKey(key)
	if err != nil {
		return nil, err
	}
	data := models.NewEmptyData(dt)

	err = s.db.View(func(txn *badger.Txn) error {
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
	prefix := []byte(DataPrefixMap[dataType])

	key = append(key, prefix...)
	key = append(key, hash32[:]...)

	return
}

func DataTypeFromKey(key []byte) (models.DataType, error) {
	prefix := Prefix(key[:len(key)-32])
	dt, ok := DataTypeMap[prefix]
	if !ok {
		return dt, errors.Errorf("can't get data type from key")
	}
	return dt, nil
}
