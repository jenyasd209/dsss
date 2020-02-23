package storage

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
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
	defer storage.Close()
	require.NotNil(t, storage)
}

func TestStorage_Add(t *testing.T) {
	storage := NewDefaultStorage("/tmp/badger")
	defer storage.Close()

	data := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		randomContent(),
	)

	obtainedKey, err := storage.Add(data)
	assert.Nil(t, err, err)
	assert.Equal(t, data.ID(), obtainedKey)
}

func TestStorage_Read(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "badger")
	require.Nil(t, err, err)

	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test")
	storage := NewDefaultStorage(tmpFile)
	defer storage.Close()

	expectedData := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		randomContent(),
	)

	obtainedKey, err := storage.Add(expectedData)
	require.Nil(t, err, err)
	assert.Equal(t, expectedData.ID(), obtainedKey)

	data := models.NewEmptyData(expectedData.DataType)
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
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		randomContent(),
	)

	obtainedKey, err := storage.Add(expectedData)
	require.Nil(t, err, err)
	require.NotNil(t, obtainedKey)
	assert.Equal(t, expectedData.ID(), obtainedKey)

	err = storage.Delete(obtainedKey)
	assert.Nil(t, err, err)

	data := models.NewEmptyData(expectedData.DataType)
	err = storage.Read(obtainedKey, data)
	assert.NotNil(t, err, err)
	assert.Equal(t, models.NewEmptyData(expectedData.DataType), data)

	err = storage.Delete(obtainedKey)
	assert.Equal(t, ErrIDNotFound, err)
}

func TestStorage_AddConcurrent(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "badger")
	require.Nil(t, err, err)

	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test")
	storage := NewDefaultStorage(tmpFile)
	defer storage.Close()

	datas := sameData(10, models.Audio)

	errCount := 0
	wg := sync.WaitGroup{}
	for _, data := range datas {
		wg.Add(1)
		go func() {
			_, err := storage.Add(data)
			if err != nil {
				println(err.Error())
				errCount++
			}

			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, len(datas)-1, errCount)

	obtainedData := models.NewEmptyData(datas[0].Type())
	err = storage.Read(datas[0].ID(), obtainedData)
	require.Nil(t, err, err)
	require.Nil(t, err, err)
}

func TestStorage_DeleteConcurrent(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "badger")
	require.Nil(t, err, err)

	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test")
	storage := NewDefaultStorage(tmpFile)
	defer storage.Close()

	expectedData := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		randomContent(),
	)

	obtainedKey, err := storage.Add(expectedData)
	assert.Nil(t, err, err)

	errCount := 0
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			err := storage.Delete(obtainedKey)
			if err != nil {
				println(err.Error())
				errCount++
			}

			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, 9, errCount)

	data := models.NewEmptyData(expectedData.DataType)
	err = storage.Read(obtainedKey, data)
	assert.NotNil(t, err, err)
	assert.Equal(t, models.NewEmptyData(expectedData.DataType), data)
}

func TestDataTypeFromKey(t *testing.T) {
	data := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)

	for i := 0; i < 10; i++ {
		dt, err := models.DataTypeFromID(data.ID())
		assert.Nil(t, err)
		assert.Equal(t, data.DataType, dt)
	}
}

func sameData(count int, dataType models.DataType) (datas []models.Data) {
	content := randomContent()

	for i := 0; i < count; i++ {
		datas = append(datas, randomData(content, dataType))
	}
	return
}

func randomData(content models.Content, dataType models.DataType) (data models.Data) {
	if content == nil {
		content = randomContent()
	}

	return models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: dataType,
		},
		content,
	)
}

func randomContent() []byte {
	b := make([]byte, 10)
	rand.Read(b)

	return b
}
