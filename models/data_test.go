package models

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var content = []byte("test content")

//var json = jsonData{content}
//
//var audioData = AudioData{
//	metaData: metaData{
//		Title: "test_audio",
//	},
//	Content: content,
//}
//
//var videoData = VideoData{
//	metaData: metaData{
//		Title: "test_video",
//	},
//	Frames: []Content{
//		[]byte("frame1"),
//		[]byte("frame3"),
//		[]byte("frame4"),
//	},
//}

func TestSimpleDataBuilder(t *testing.T) {
	simple := NewSimpleData().
		SetMetadata(
			NewMetaData().SetTitle("title").Build(),
		).SetContent(content).Build()

	simple2 := simpleData{
		metaData{
			title:"title",
			dataType:Simple,
			cashedHash: sha256.Sum256(content),
		},
		content,
	}

	assert.Equal(t, simple, simple2)
}

func TestSimpleData_Hash(t *testing.T) {
	simple := NewSimpleData().
		SetMetadata(
			NewMetaData().SetTitle("title").SetType(Simple).Build(),
		).SetContent(content).Build()

	content := simple.GetContent()
	expectedHash := sha256.Sum256(*content)

	assert.Equal(t, hex.EncodeToString(expectedHash[:]), simple.GetMetadata().GetCashedHash())
}

func TestSimpleData_MarshalBinary(t *testing.T) {
	s := simpleData{
		metaData{
			title:"title",
			dataType:Simple,
			cashedHash: sha256.Sum256(content),
		},
		content,
	}

	expected, err := json.Marshal(&struct {
		Title string
		CashedHash Hash32
		DataType DataType
		Content Content
	}{
		Title: s.metaData.title,
		CashedHash: s.metaData.cashedHash,
		DataType: s.metaData.dataType,
		Content: s.content,
	})
	assert.Nil(t, err)

	println(string(expected))

	b, err := s.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expected, b)
}

func TestSimpleData_UnmarshalBinary(t *testing.T) {
	simple := NewSimpleData().
		SetMetadata(
			NewMetaData().SetTitle("title").SetType(Simple).Build(),
		).SetContent(content).Build()

	data := simpleData{}

	bytes, err := simple.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, simple, data)
}

//func TestJsonData_Hash(t *testing.T) {
//	bytes, err := jsonData.MarshalBinary()
//	assert.Nil(t, err, err)
//
//	expectedHash := sha256.Sum256(bytes)
//	assert.Equal(t, hex.EncodeToString(expectedHash[:]), jsonData.Hash().String())
//}
//
//func TestJsonData_MarshalBinary(t *testing.T) {
//	expectedJSON, err := json.Marshal(jsonData)
//	assert.Nil(t, err)
//
//	JSON, err := jsonData.MarshalBinary()
//	assert.Nil(t, err)
//	assert.Equal(t, expectedJSON, JSON)
//}
//
//func TestJsonData_UnmarshalBinary(t *testing.T) {
//	data := jsonData{}
//
//	bytes, err := jsonData.MarshalBinary()
//	assert.Nil(t, err)
//
//	err = data.UnmarshalBinary(bytes)
//	assert.Nil(t, err)
//	assert.Equal(t, jsonData, data)
//}
//
//func TestAudioData_Hash(t *testing.T) {
//	bytes, err := audioData.MarshalBinary()
//	assert.Nil(t, err, err)
//
//	expectedHash := sha256.Sum256(bytes)
//	assert.Equal(t, hex.EncodeToString(expectedHash[:]), audioData.Hash().String())
//}
//
//func TestAudioData_MarshalBinary(t *testing.T) {
//	expected, err := json.Marshal(audioData)
//	assert.Nil(t, err)
//
//	data, err := audioData.MarshalBinary()
//	assert.Nil(t, err)
//	assert.Equal(t, expected, data)
//}
//
//func TestAudioData_UnmarshalBinary(t *testing.T) {
//	data := AudioData{}
//
//	bytes, err := audioData.MarshalBinary()
//	assert.Nil(t, err)
//
//	err = data.UnmarshalBinary(bytes)
//	assert.Nil(t, err)
//	assert.Equal(t, audioData, data)
//}
//
//func TestVideoData_Hash(t *testing.T) {
//	bytes, err := videoData.MarshalBinary()
//	assert.Nil(t, err, err)
//
//	expectedHash := sha256.Sum256(bytes)
//	assert.Equal(t, hex.EncodeToString(expectedHash[:]), videoData.Hash().String())
//}
//
//func TestVideoData_MarshalBinary(t *testing.T) {
//	expected, err := json.Marshal(videoData)
//	assert.Nil(t, err)
//
//	data, err := videoData.MarshalBinary()
//	assert.Nil(t, err)
//	assert.Equal(t, expected, data)
//}
//
//func TestVideoData_UnmarshalBinary(t *testing.T) {
//	data := VideoData{}
//
//	bytes, err := videoData.MarshalBinary()
//	assert.Nil(t, err)
//
//	err = data.UnmarshalBinary(bytes)
//	assert.Nil(t, err)
//	assert.Equal(t, videoData, data)
//}
