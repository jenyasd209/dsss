package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var content = []byte("test content")

var expectedSimpleData = NewSimpleData(
	MetaData{
		Title: "text_test",
	},
	content,
)

var expectedJsonData = NewJSONData(
	MetaData{
		Title: "text_json",
	},
	content,
)

var expectedAudioData = NewAudioData(
	MetaData{
		Title: "test_audio",
	},
	content,
)

var expectedVideoData = NewVideoData(
	MetaData{
		Title: "test_video",
	},
	[]byte("frame1"),
)

func TestSimpleData_Hash(t *testing.T) {
	expectedHash := sha256.Sum256(expectedSimpleData.Content)
	assert.Equal(t, hex.EncodeToString(expectedHash[:]), expectedSimpleData.CachedHash().String())
}

func TestSimpleData_MarshalBinary(t *testing.T) {
	expectedJSON, err := json.Marshal(expectedSimpleData)
	assert.Nil(t, err)

	JSON, err := expectedSimpleData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, JSON)
}

func TestSimpleData_UnmarshalBinary(t *testing.T) {
	data := &simpleData{}

	bytes, err := expectedSimpleData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedSimpleData, data)
}

func TestJsonData_MarshalBinary(t *testing.T) {
	expectedJSON, err := json.Marshal(expectedJsonData)
	assert.Nil(t, err)

	JSON, err := expectedJsonData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, JSON)
}

func TestJsonData_UnmarshalBinary(t *testing.T) {
	data := &jsonData{}

	bytes, err := expectedJsonData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedJsonData, data)
}

func TestAudioData_Hash(t *testing.T) {
	expectedHash := sha256.Sum256(expectedAudioData.Content)
	assert.Equal(t, hex.EncodeToString(expectedHash[:]), expectedAudioData.CachedHash().String())
}

func TestAudioData_MarshalBinary(t *testing.T) {
	expected, err := json.Marshal(expectedAudioData)
	assert.Nil(t, err)

	data, err := expectedAudioData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}

func TestAudioData_UnmarshalBinary(t *testing.T) {
	data := &audioData{}

	bytes, err := expectedAudioData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedAudioData, data)
}

func TestVideoData_Hash(t *testing.T) {
	expectedHash := sha256.Sum256(expectedVideoData.Frames)
	assert.Equal(t, hex.EncodeToString(expectedHash[:]), expectedVideoData.CachedHash().String())
}

func TestVideoData_MarshalBinary(t *testing.T) {
	expected, err := json.Marshal(expectedVideoData)
	assert.Nil(t, err)

	data, err := expectedVideoData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}

func TestVideoData_UnmarshalBinary(t *testing.T) {
	data := &videoData{}

	bytes, err := expectedVideoData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedVideoData, data)
}
