package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

var (
	addr    = "http://localhost:8080"
	dsssAPI = API(addr)
)

func TestFileRoute_Add(t *testing.T) {
	filename := "text.mp3"
	content := randomContent()

	key, err := dsssAPI.Files().Add(filename, content)
	require.Nil(t, err, err)
	assert.NotNil(t, key)
}

func TestFileRoute_Get(t *testing.T) {
	filename := "text.mp3"
	content := randomContent()

	key, err := dsssAPI.Files().Add(filename, content)
	require.Nil(t, err, err)
	assert.NotNil(t, key)

	data, err := dsssAPI.Files().Get(key)
	require.Nil(t, err, err)
	require.NotNil(t, data)
	assert.Equal(t, content, []byte(data.Body()))
	assert.Equal(t, filename, data.Title())
}

func TestFileRoute_Delete(t *testing.T) {
	filename := "text.mp3"
	content := randomContent()

	key, err := dsssAPI.Files().Add(filename, content)
	require.Nil(t, err, err)
	assert.NotNil(t, key)

	data, err := dsssAPI.Files().Get(key)
	require.Nil(t, err, err)
	assert.NotNil(t, data)
	assert.Equal(t, content, []byte(data.Body()))
	assert.Equal(t, filename, data.Title())

	_, err = dsssAPI.Files().Delete(key)
	require.Nil(t, err, err)

	d, err := dsssAPI.Files().Get(key)
	assert.NotNil(t, err, err)
	require.Nil(t, d)
}

func randomContent() []byte {
	rand.Seed(time.Now().UnixNano())

	bytes := make([]byte, 4)
	rand.Read(bytes)

	return bytes
}
