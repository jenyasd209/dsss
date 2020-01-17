package models

import (
	"crypto/sha256"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleData_Hash(t *testing.T) {
	content := []byte("content")
	expectedHash := sha256.Sum256(content)
	data := SimpleData{content}

	assert.Equal(t, Hash32(expectedHash), data.Hash())
}

func TestSimpleData_MarshalBinary(t *testing.T) {
	content := []byte("content")
	data := SimpleData{content}

	expectedJSON, err := json.Marshal(data)
	assert.Nil(t, err)

	JSON, err := data.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, JSON)
}

func TestSimpleData_UnmarshalBinary(t *testing.T) {
	expectedData := SimpleData{[]byte("content")}

	bytes, err := expectedData.MarshalBinary()
	assert.Nil(t, err)

	data := SimpleData{}
	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedData, data)
}
