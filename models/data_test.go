package models

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

var content = Content([]byte("test content"))

func TestDataTypeFromID(t *testing.T) {
	expectedSimpleData, err := NewData(NewMetaData("test", Simple), content)
	assert.Nil(t, err, err)

	dt, err := DataTypeFromID(expectedSimpleData.Meta().GetID())
	assert.Nil(t, err, err)
	assert.Equal(t, Simple, dt)

	expectedSimpleData, err = NewData(NewMetaData("test", JSON), content)
	assert.Nil(t, err, err)

	dt, err = DataTypeFromID(expectedSimpleData.Meta().GetID())
	assert.Nil(t, err, err)
	assert.Equal(t, JSON, dt)

	expectedSimpleData, err = NewData(NewMetaData("test", Audio), content)
	assert.Nil(t, err, err)

	dt, err = DataTypeFromID(expectedSimpleData.Meta().GetID())
	assert.Nil(t, err, err)
	assert.Equal(t, Audio, dt)

	expectedSimpleData, err = NewData(NewMetaData("test", Video), content)
	assert.Nil(t, err, err)

	dt, err = DataTypeFromID(expectedSimpleData.Meta().GetID())
	assert.Nil(t, err, err)
	assert.Equal(t, Video, dt)
}

func TestNewData(t *testing.T) {
	expectedSimpleData, err := NewData(NewMetaData("simple", Simple), content)
	assert.Nil(t, err, err)
	assert.Equal(t, "simple", expectedSimpleData.Meta().GetTitle())
	assert.Equal(t, content, *expectedSimpleData.Body())
}

func TestSimpleData_MarshalBinary(t *testing.T) {
	expectedSimpleData, err := NewData(NewMetaData("simple", Simple), content)
	assert.Nil(t, err, err)

	expectedJSON, err := json.Marshal(expectedSimpleData)
	assert.Nil(t, err)

	JSON, err := expectedSimpleData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, JSON)
}

func TestSimpleData_UnmarshalBinary(t *testing.T) {
	expectedSimpleData, err := NewData(NewMetaData("simple", Simple), content)
	assert.Nil(t, err, err)

	data := &simpleData{}

	bytes, err := expectedSimpleData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedSimpleData, data)
}

func TestJsonData_MarshalBinary(t *testing.T) {
	expectedJsonData, err := NewData(NewMetaData("json", JSON), content)
	assert.Nil(t, err, err)

	expectedJSON, err := json.Marshal(expectedJsonData)
	assert.Nil(t, err)

	JSON, err := expectedJsonData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, JSON)
}

func TestJsonData_UnmarshalBinary(t *testing.T) {
	expectedJsonData, err := NewData(NewMetaData("json", JSON), content)
	assert.Nil(t, err, err)
	data := &jsonData{}

	bytes, err := expectedJsonData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedJsonData, data)
}

func TestAudioData_MarshalBinary(t *testing.T) {
	expectedAudioData, err := NewData(NewMetaData("audio", Audio), content)
	assert.Nil(t, err, err)

	expected, err := json.Marshal(expectedAudioData)
	assert.Nil(t, err)

	data, err := expectedAudioData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}

func TestAudioData_UnmarshalBinary(t *testing.T) {
	expectedAudioData, err := NewData(NewMetaData("audio", Audio), content)
	assert.Nil(t, err, err)
	data := &audioData{}

	bytes, err := expectedAudioData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedAudioData, data)
}

func TestVideoData_MarshalBinary(t *testing.T) {
	expectedVideoData, err := NewData(NewMetaData("video", Video), content)
	assert.Nil(t, err, err)

	expected, err := json.Marshal(expectedVideoData)
	assert.Nil(t, err)

	data, err := expectedVideoData.MarshalBinary()
	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}

func TestVideoData_UnmarshalBinary(t *testing.T) {
	expectedVideoData, err := NewData(NewMetaData("video", Video), content)
	assert.Nil(t, err, err)

	data := &videoData{}

	bytes, err := expectedVideoData.MarshalBinary()
	assert.Nil(t, err)

	err = data.UnmarshalBinary(bytes)
	assert.Nil(t, err)
	assert.Equal(t, expectedVideoData, data)
}

func TestConvertToDataType(t *testing.T) {
	s := "3"
	f := 3.0
	ip := 3
	in := -3

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
	expectedVideoData, err := NewData(NewMetaData("video", Video), content)
	assert.Nil(t, err, err)

	dt, err := DataTypeFromID(expectedVideoData.Meta().GetID())
	require.Nil(t, err, err)

	assert.Equal(t, expectedVideoData.Meta().GetType(), dt)
}
