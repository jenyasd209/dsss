package storage

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/iorhachovyevhen/dsss/models"
)

func TestNewStorageWithOptions(t *testing.T) {
	opt := NewOptions().
		WithDir("/tmp/badger").
		WithValueDir("/tmp/badger").
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
		[]byte("content"),
	)

	key := composeKey(data.ID(), data.Type())

	obtainedKey, err := storage.Add(data)
	assert.Nil(t, err, err)
	assert.Equal(t, key, obtainedKey)
}

func TestStorage_Read(t *testing.T) {
	storage := NewDefaultStorage("/tmp/badger")
	defer storage.Close()

	expectedData := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		randomContent(),
	)

	key := composeKey(expectedData.ID(), expectedData.Type())

	obtainedKey, err := storage.Add(expectedData)
	require.Nil(t, err, err)
	assert.Equal(t, key, obtainedKey)

	data, err := storage.Read(obtainedKey)
	assert.Nil(t, err, err)
	assert.Equal(t, expectedData, data)
}

func TestStorage_Delete(t *testing.T) {
	storage := NewDefaultStorage("/tmp/badger")
	defer storage.Close()

	expectedData := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)

	obtainedKey, err := storage.Add(expectedData)
	assert.Nil(t, err, err)

	err = storage.Delete(obtainedKey)
	assert.Nil(t, err, err)

	data, err := storage.Read(obtainedKey)
	assert.NotNil(t, err, err)
	assert.Nil(t, data)
}

func TestStorage_AddConcurrent(t *testing.T) {
	storage := NewDefaultStorage("/tmp/badger")
	defer storage.Close()

	datas := randomData(10)
	key := composeKey(datas[0].ID(), datas[0].Type())

	var err error
	var title string
	var wg sync.WaitGroup
	for _, data := range datas {
		wg.Add(1)
		go func() {
			_, err = storage.Add(data)
			title = data.Title()

			wg.Done()
		}()
	}

	wg.Wait()

	obtainedData, err2 := storage.Read(key)
	require.Nil(t, err2, err2)
	assert.Equal(t, title, obtainedData.Title())

	err = storage.Delete(key)
	assert.Nil(t, err, err)
}

func TestStorage_DeleteConcurrent(t *testing.T) {
	storage := NewDefaultStorage("/tmp/badger")
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

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			_ = storage.Delete(obtainedKey)

			wg.Done()
		}()
	}

	wg.Wait()

	data, err := storage.Read(obtainedKey)
	assert.NotNil(t, err, err)
	assert.Nil(t, data)
}

func TestDataTypeFromKey(t *testing.T) {
	data := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)
	id := composeKey(data.ID(), data.Type())

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			dt, err := DataTypeFromKey(id)
			assert.Nil(t, err)
			assert.Equal(t, data.DataType, dt)
			wg.Done()
		}()
	}
}

func randomData(count int) (datas []models.Data) {
	content := randomContent()

	for i := 0; i < count; i++ {
		datas = append(datas, models.NewSimpleData(
			models.MetaData{
				Title:    "test.txt",
				DataType: models.Audio,
			},
			content,
		))
	}
	return
}

func randomContent() []byte {
	b := make([]byte, 4)
	rand.Read(b)

	return b
}
