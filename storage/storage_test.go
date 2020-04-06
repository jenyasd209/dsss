package storage

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iorhachovyevhen/dsss/models"
)

func TestNewStorageWithOptions(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "badger")
	require.Nil(t, err, err)

	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test")
	opt := NewOptions().
		WithDir(tmpFile).
		WithValueDir(tmpFile).
		WithValueLogFileSize(2 << 20)

	storage := NewStorageWithOptions(opt)
	require.NotNil(t, storage)

	err = storage.Close()
	require.Nil(t, err)
}

func TestStorage_Add(t *testing.T) {
	storage := NewDefaultStorage("/tmp/badger")
	defer storage.Close()

	data := models.NewSimpleData(
		models.NewMetaData("test", models.Simple),
		randomContent(),
	)

	obtainedKey, err := storage.Add(data)
	assert.Nil(t, err, err)
	assert.Equal(t, data.Meta().GetID(), obtainedKey)
}

func TestStorage_Read(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "badger")
	require.Nil(t, err, err)

	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test")
	storage := NewDefaultStorage(tmpFile)
	defer storage.Close()

	expectedData := models.NewSimpleData(
		models.NewMetaData("test", models.Simple),
		randomContent(),
	)

	obtainedKey, err := storage.Add(expectedData)
	require.Nil(t, err, err)
	assert.Equal(t, expectedData.Meta().GetID(), obtainedKey)

	data, err := models.NewData(models.NewMetaData("test", models.Simple), []byte(""))
	assert.Nil(t, err, err)

	err = storage.Read(obtainedKey, data)
	assert.Nil(t, err, err)
	assert.Equal(t, expectedData, data)
}

func TestStorage_Delete(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "badgerTmp")
	require.Nil(t, err, err)

	defer os.RemoveAll(tmpDir)

	storage := NewDefaultStorage(tmpDir)
	defer storage.Close()

	expectedData := models.NewSimpleData(
		models.NewMetaData("test", models.Simple),
		randomContent(),
	)

	obtainedKey, err := storage.Add(expectedData)
	require.Nil(t, err, err)
	require.NotNil(t, obtainedKey)
	assert.Equal(t, expectedData.Meta().GetID(), obtainedKey)

	err = storage.Delete(obtainedKey)
	assert.Nil(t, err, err)

	data, err := models.NewData(models.NewMetaData("", models.Simple), []byte(""))
	assert.Nil(t, err, err)

	cpData, err := models.NewData(models.NewMetaData("", models.Simple), []byte(""))
	assert.Nil(t, err, err)

	err = storage.Read(obtainedKey, data)
	assert.NotNil(t, err, err)
	assert.Equal(t, cpData.Body(), data.Body())
	assert.Equal(t, cpData.Meta().GetType(), data.Meta().GetType())
	assert.Equal(t, cpData.Meta().GetTitle(), data.Meta().GetTitle())

	err = storage.Delete(obtainedKey)
	assert.Equal(t, ErrIDNotFound, err)
}

func TestDataTypeFromKey(t *testing.T) {
	data := models.NewSimpleData(
		models.NewMetaData("test", models.Simple),
		[]byte("content"),
	)

	for i := 0; i < 10; i++ {
		dt, err := models.DataTypeFromID(data.Meta().GetID())
		assert.Nil(t, err)
		assert.Equal(t, data.DataType, dt)
	}
}

func randomContent() []byte {
	b := make([]byte, 10)
	rand.Read(b)

	return b
}
