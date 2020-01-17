package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dsss/models"
)

func TestWriteData(t *testing.T) {
	data := models.SimpleData{Content: []byte("test")}

	value, err := data.MarshalBinary()
	assert.Nil(t, err)

	err = WriteData(data.Hash(), value)
	assert.Nil(t, err)
}

func TestReadData(t *testing.T) {
	expectedData := models.SimpleData{
		Content: []byte("test"),
	}

	bytes, err := expectedData.MarshalBinary()
	assert.Nil(t, err)

	err = WriteData(expectedData.Hash(), bytes)
	assert.Nil(t, err)

	reade, err := ReadData(expectedData.Hash())
	assert.Nil(t, err)

	data := models.SimpleData{}
	err = data.UnmarshalBinary(reade)
	assert.Nil(t, err)

	assert.Equal(t, expectedData, data)
}

func TestUpdateData(t *testing.T) {
	content := []byte("test")
	updatedContent := []byte("test_up")
	expectedData := models.SimpleData{Content: content}
	data := models.SimpleData{}

	value, err := expectedData.MarshalBinary()
	assert.Nil(t, err)

	err = WriteData(expectedData.Hash(), value)
	assert.Nil(t, err)

	expectedData.Content = updatedContent
	value, err = expectedData.MarshalBinary()
	assert.Nil(t, err)

	err = WriteData(expectedData.Hash(), value)
	assert.Nil(t, err)

	value, err = ReadData(expectedData.Hash())
	assert.Nil(t, err)

	err = data.UnmarshalBinary(value)
	assert.Nil(t, err)

	assert.Equal(t, expectedData, data)
}

func TestDeleteData(t *testing.T) {
	data := models.SimpleData{Content: []byte("test")}

	value, err := data.MarshalBinary()
	assert.Nil(t, err)

	err = WriteData(data.Hash(), value)
	assert.Nil(t, err)

	err = DeleteData(data.Hash())
	assert.Nil(t, err)

	value, err = ReadData(data.Hash())
	assert.NotNil(t, err)
	assert.Nil(t, value)
}