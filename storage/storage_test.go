package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iorhachovyevhen/dsss/models"
)

func TestWriteData(t *testing.T) {
	data := models.SimpleData{Content: []byte("test")}

	value, err := data.MarshalBinary()
	assert.Nil(t, err, err)

	err = NewData(data.Hash(), value)
	assert.Nil(t, err, err)
}

func TestReadSimple(t *testing.T) {
	expectedData := models.SimpleData{
		Content: []byte("test"),
	}

	bytes, err := expectedData.MarshalBinary()
	assert.Nil(t, err, err)

	err = NewData(expectedData.Hash(), bytes)
	assert.Nil(t, err, err)

	value, err := ReadData(expectedData.Hash())
	assert.Nil(t, err, err)

	data := models.SimpleData{}
	err = data.UnmarshalBinary(value)
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

	value, err := expectedData.MarshalBinary()
	assert.Nil(t, err, err)

	err = NewData(expectedData.Hash(), value)
	assert.Nil(t, err, err)

	value, err = updatedData.MarshalBinary()
	assert.Nil(t, err, err)

	err = UpdateData(expectedData.Hash(), updatedData.Hash(), value)
	assert.Nil(t, err, err)

	value, err = ReadData(expectedData.Hash())
	assert.NotNil(t, err, err)

	value, err = ReadData(updatedData.Hash())
	assert.Nil(t, err, err)

	data := models.SimpleData{}

	err = data.UnmarshalBinary(value)
	assert.Nil(t, err, err)

	assert.Equal(t, updatedData, data)
}

func TestDeleteData(t *testing.T) {
	data := models.SimpleData{Content: []byte("test")}

	value, err := data.MarshalBinary()
	assert.Nil(t, err, err)

	err = NewData(data.Hash(), value)
	assert.Nil(t, err, err)

	err = DeleteData(data.Hash())
	assert.Nil(t, err, err)

	value, err = ReadData(data.Hash())
	assert.NotNil(t, err, err)
	assert.Nil(t, value)
}
