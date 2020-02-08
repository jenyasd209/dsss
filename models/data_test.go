package models

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

var content = []byte("test content")

var (
	expectedSimpleData = NewData(Simple, "text", content)
	expectedJsonData   = NewData(JSON, "json", content)
	expectedAudioData  = NewData(Audio, "audio", content)
	expectedVideoData  = NewData(Video, "video", content)
)

func TestSimpleData_Hash(t *testing.T) {
	expectedHash := sha256.Sum256(expectedSimpleData.Body())
	h := composeID(expectedHash, expectedSimpleData.Type())
	assert.Equal(t, string(h), expectedSimpleData.ID().String())
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
	expectedHash := sha256.Sum256(expectedAudioData.Body())
	h := composeID(expectedHash, expectedAudioData.Type())
	assert.Equal(t, string(h), expectedAudioData.ID().String())
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
	expectedHash := sha256.Sum256(expectedVideoData.Body())
	h := composeID(expectedHash, expectedVideoData.Type())
	assert.Equal(t, h.String(), expectedVideoData.ID().String())
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

func TestConvertToDataType(t *testing.T) {
	s := "2"
	f := 2.0
	ip := 2
	in := -2

	dt, err := ConvertToDataType(s)
	require.Nil(t, err, err)
	assert.Equal(t, Audio, dt)

	dt, err = ConvertToDataType(f)
	require.Nil(t, err, err)
	assert.Equal(t, Audio, dt)

	dt, err = ConvertToDataType(ip)
	require.Nil(t, err, err)
	assert.Equal(t, Audio, dt)

	dt, err = ConvertToDataType(in)
	require.NotNil(t, err, err)

	dt, err = ConvertToDataType("sdfs")
	require.Equal(t, ErrorBadDataType, err)
}

func TestDataTypeFromKey(t *testing.T) {
	expectedHash := sha256.Sum256(expectedVideoData.Body())
	id := composeID(expectedHash, expectedVideoData.Type())
	assert.Equal(t, id.String(), expectedVideoData.ID().String())

	dt, err := DataTypeFromID(expectedVideoData.ID())
	require.Nil(t, err, err)

	assert.Equal(t, expectedVideoData.Type(), dt)
}
