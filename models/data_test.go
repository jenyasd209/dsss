package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var content = []byte("test content")

var simpleData = SimpleData{
	MetaData: MetaData{
		Title: "text_test",
	},
	Content: content,
}

var jsonData = JsonData{
	MetaData{
		Title: "text_json",
	},
	content,
}

var audioData = AudioData{
	MetaData: MetaData{
		Title: "test_audio",
	},
	Content: content,
}

var videoData = VideoData{
	MetaData: MetaData{
		Title: "test_video",
	},
	Frames: []Content{
		[]byte("frame1"),
		[]byte("frame3"),
		[]byte("frame4"),
	},
}

func TestSimpleData_ApplyContent(t *testing.T) {
	err := simpleData.Submit()
	assert.Nil(t, err, err)

	err = simpleData.Submit()
	assert.NotNil(t, err, err)
}

func TestSimpleData_Hash(t *testing.T) {
	expectedHash := sha256.Sum256(simpleData.Content)
	assert.Equal(t, hex.EncodeToString(expectedHash[:]), simpleData.CachedHash().String())
}

func TestSimpleData_MarshalBinary(t *testing.T) {
	expectedJSON, err := json.Marshal(simpleData)
	assert.Nil(t, err)

	JSON, err := simpleData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, JSON)
}

func TestSimpleData_UnmarshalBinary(t *testing.T) {
	data := SimpleData{}

	bytes, err := simpleData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, simpleData, data)
}

func TestJsonData_Hash(t *testing.T) {
	jsonData.Submit()

	expectedHash := sha256.Sum256(jsonData.Content)
	assert.Equal(t, hex.EncodeToString(expectedHash[:]), jsonData.CachedHash().String())
}

func TestJsonData_MarshalBinary(t *testing.T) {
	expectedJSON, err := json.Marshal(jsonData)
	assert.Nil(t, err)

	JSON, err := jsonData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, JSON)
}

func TestJsonData_UnmarshalBinary(t *testing.T) {
	data := JsonData{}

	bytes, err := jsonData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, jsonData, data)
}

func TestAudioData_Hash(t *testing.T) {
	audioData.Submit()

	expectedHash := sha256.Sum256(audioData.Content)
	assert.Equal(t, hex.EncodeToString(expectedHash[:]), audioData.CachedHash().String())
}

func TestAudioData_MarshalBinary(t *testing.T) {
	expected, err := json.Marshal(audioData)
	assert.Nil(t, err)

	data, err := audioData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}

func TestAudioData_UnmarshalBinary(t *testing.T) {
	data := AudioData{}

	bytes, err := audioData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, audioData, data)
}

func TestVideoData_Hash(t *testing.T) {
	videoData.Submit()

	expectedHash := sha256.Sum256(videoData.Frames[0])
	assert.Equal(t, hex.EncodeToString(expectedHash[:]), videoData.CachedHash().String())
}

func TestVideoData_MarshalBinary(t *testing.T) {
	expected, err := json.Marshal(videoData)
	assert.Nil(t, err)

	data, err := videoData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}

func TestVideoData_UnmarshalBinary(t *testing.T) {
	data := VideoData{}

	bytes, err := videoData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, videoData, data)
}
