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
	assert.Equal(t, expectedData, data)
}

//
//func TestUpdateData(t *testing.T) {
//	expectedData := models.SimpleData{
//		Content: []byte("test"),
//	}
//	updatedData := models.SimpleData{
//		Content: []byte("update"),
//	}
//
//	err := NewData(&expectedData)
//	assert.Nil(t, err, err)
//
//	err = UpdateData(expectedData.Hash().String(), &updatedData)
//	assert.Nil(t, err, err)
//
//	data := models.SimpleData{}
//	err = ReadData(expectedData.Hash().String(), &data)
//	assert.NotNil(t, err, err)
//
//	err = ReadData(updatedData.Hash().String(), &data)
//	assert.Nil(t, err, err)
//
//	assert.Equal(t, updatedData, data)
//}
//
//func TestDeleteData(t *testing.T) {
//	data := models.SimpleData{Content: []byte("test")}
//
//	err := NewData(&data)
//	assert.Nil(t, err, err)
//
//	err = DeleteData(data.Hash().String())
//	assert.Nil(t, err, err)
//
//	expectedData := models.SimpleData{}
//	err = ReadData(data.Hash().String(), &expectedData)
//	assert.NotNil(t, err, err)
//	assert.Equal(t, models.SimpleData{}, expectedData)
//}
