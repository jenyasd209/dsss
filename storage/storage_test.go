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

	err := storage.Add(data)
	assert.Nil(t, err, err)
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

	err := storage.Add(expectedData)
	assert.Nil(t, err, err)

	data := models.NewSimpleData(models.MetaData{}, nil)

	err = storage.Read(expectedData.ID(), data)
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

	err := storage.Add(expectedData)
	assert.Nil(t, err, err)

	err = storage.Delete(expectedData.ID(), expectedData.Type())
	assert.Nil(t, err, err)

	data := models.NewSimpleData(models.MetaData{}, nil)

	err = storage.Read(expectedData.ID(), data)
	assert.NotNil(t, err, err)
	assert.Equal(t, models.NewSimpleData(models.MetaData{}, nil), data)
}
