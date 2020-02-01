package storage

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iorhachovyevhen/dsss/models"
)

func TestNewStorageWithOptions(t *testing.T) {
	opt := NewOptions().
		WithDir("/tmp/badger").
		WithValueDir("/tmp/badger").
		WithValueLogFileSize(2 << 20)

	storage := NewStorageWithOptions(opt)
	defer storage.Close()
	require.NotNil(t, storage)
}

func TestStorage_Add(t *testing.T) {
	storage := NewDefaultStorage("/tmp/badger")
	defer storage.Close()

	data := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)

	key := composeKey(data.ID(), data.Type())

	obtainedKey, err := storage.Add(data)
	assert.Nil(t, err, err)
	assert.Equal(t, key, obtainedKey)
}

func TestStorage_Read(t *testing.T) {
	storage := NewDefaultStorage("/tmp/badger")
	defer storage.Close()

	expectedData := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)

	obtainedKey, err := storage.Add(expectedData)
	assert.Nil(t, err, err)

	data, err := storage.Read(obtainedKey)
	assert.Nil(t, err, err)
	assert.Equal(t, expectedData, data)
}

func TestDeleteData(t *testing.T) {
	storage := NewDefaultStorage("/tmp/badger")
	defer storage.Close()

	expectedData := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)

	obtainedKey, err := storage.Add(expectedData)
	assert.Nil(t, err, err)

	err = storage.Delete(obtainedKey)
	assert.Nil(t, err, err)

	data, err := storage.Read(obtainedKey)
	assert.NotNil(t, err, err)
	assert.Nil(t, data)
}

func TestDataTypeFromKey(t *testing.T) {
	data := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)
	id := composeKey(data.ID(), data.Type())
	dt := DataTypeFromKey(id)
	assert.Equal(t, data.DataType, dt)
}
