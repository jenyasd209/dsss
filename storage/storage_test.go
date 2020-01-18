package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iorhachovyevhen/dsss/models"
)

func TestWriteData(t *testing.T) {
	data := models.SimpleData{Content: []byte("test")}

	err := NewData(&data)
	assert.Nil(t, err, err)
}

func TestReadData(t *testing.T) {
	expectedData := models.SimpleData{
		Content: []byte("test"),
	}

	err := NewData(&expectedData)
	assert.Nil(t, err, err)

	data := models.SimpleData{}
	err = ReadData(expectedData.Hash().String(), &data)
	assert.Nil(t, err, err)

	assert.Equal(t, expectedData, data)
}

func TestUpdateData(t *testing.T) {
	expectedData := models.SimpleData{
		Content: []byte("test"),
	}
	updatedData := models.SimpleData{
		Content: []byte("update"),
	}

	err := NewData(&expectedData)
	assert.Nil(t, err, err)

	err = UpdateData(expectedData.Hash().String(), &updatedData)
	assert.Nil(t, err, err)

	data := models.SimpleData{}
	err = ReadData(expectedData.Hash().String(), &data)
	assert.NotNil(t, err, err)

	err = ReadData(updatedData.Hash().String(), &data)
	assert.Nil(t, err, err)

	assert.Equal(t, updatedData, data)
}

func TestDeleteData(t *testing.T) {
	data := models.SimpleData{Content: []byte("test")}

	err := NewData(&data)
	assert.Nil(t, err, err)

	err = DeleteData(data.Hash().String())
	assert.Nil(t, err, err)

	expectedData := models.SimpleData{}
	err = ReadData(data.Hash().String(), &expectedData)
	assert.NotNil(t, err, err)
	assert.Equal(t, models.SimpleData{}, expectedData)
}
