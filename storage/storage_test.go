package storage

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iorhachovyevhen/dsss/models"
)

func TestStorage_Add(t *testing.T) {
	storage := NewDefaultStorage()
	defer storage.Close()

	data := models.SimpleData{
		MetaData: models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		Content: []byte("content"),
	}
	err := data.Submit()
	require.Nil(t, err, err)

	err = storage.Add(&data)
	assert.Nil(t, err, err)
}

func TestStorage_Read(t *testing.T) {
	storage := NewDefaultStorage()
	defer storage.Close()

	expectedData := models.SimpleData{
		MetaData: models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		Content: []byte("content"),
	}
	err := expectedData.Submit()
	require.Nil(t, err, err)

	err = storage.Add(&expectedData)
	assert.Nil(t, err, err)

	data := models.SimpleData{}

	err = storage.Read(expectedData.CachedHash(), &data)
	assert.Nil(t, err, err)
	assert.Equal(t, expectedData.Content, data.Content)
	assert.Equal(t, expectedData.MetaData.Title, data.MetaData.Title)
	assert.Equal(t, expectedData.MetaData.DataType, data.MetaData.DataType)
}

func TestDeleteData(t *testing.T) {
	storage := NewDefaultStorage()
	defer storage.Close()

	expectedData := models.SimpleData{
		MetaData: models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		Content: []byte("content"),
	}
	err := expectedData.Submit()
	require.Nil(t, err, err)

	err = storage.Add(&expectedData)
	assert.Nil(t, err, err)

	err = storage.Delete(expectedData.CachedHash())
	assert.Nil(t, err, err)

	data := models.SimpleData{}

	err = storage.Read(expectedData.CachedHash(), &data)
	assert.NotNil(t, err, err)
	assert.Equal(t, models.SimpleData{}, data)
}
