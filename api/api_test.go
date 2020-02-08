package api

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	addr    = "http://localhost:8080"
	dsssAPI = API(addr)

	filename = "text.txt"
	content  = []byte("body")
)

func TestFileRoute_Add(t *testing.T) {
	key, err := dsssAPI.Files().Add(filename, content)
	require.Nil(t, err, err)
	assert.NotNil(t, key)
}

func TestFileRoute_Get(t *testing.T) {
	key, err := dsssAPI.Files().Add(filename, content)
	require.Nil(t, err, err)
	assert.NotNil(t, key)

	data, err := dsssAPI.Files().Get(key)
	require.Nil(t, err, err)
	require.NotNil(t, data)
	assert.Equal(t, content, []byte(data.Body()))
}

func TestFileRoute_Delete(t *testing.T) {
	key, err := dsssAPI.Files().Add(filename, content)
	require.Nil(t, err, err)
	assert.NotNil(t, key)

	data, err := dsssAPI.Files().Get(key)
	require.Nil(t, err, err)
	require.NotNil(t, data)
	assert.Equal(t, content, []byte(data.Body()))

	_, err = dsssAPI.Files().Delete(key)
	require.Nil(t, err, err)

	data, err = dsssAPI.Files().Get(key)
	require.NotNil(t, err, err)
	require.Nil(t, data)
}
