package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/iorhachovyevhen/dsss/models"
)

var route = "http://localhost:8080/file"

func TestNewStorageServer(t *testing.T) {
	s := NewStorageServer()
	go func() {
		if err := s.Start(); err != nil {
			panic(err)
		}
	}()
	defer s.Shutdown()

	client := &fasthttp.Client{}
	time.Sleep(time.Second)

	t.Run("FilePost", func(t *testing.T) {
		testFile, err := models.NewData(
			models.NewMetaData("simple", models.Simple),
			randomContent(),
		)
		require.Nil(t, err, err)

		data, err := testFile.MarshalBinary()
		require.Nil(t, err, err)

		req := fasthttp.AcquireRequest()
		req.Header.SetMethod("POST")
		req.SetRequestURI(route)
		req.SetBody(data)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		err = client.Do(req, resp)
		require.Nil(t, err, err)

		assert.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
		assert.Equal(t, string(resp.Body()), testFile.Meta().GetID().String())
	})

	t.Run("FileGet", func(t *testing.T) {
		testFile, err := models.NewData(
			models.NewMetaData("simple", models.Simple),
			randomContent(),
		)
		require.Nil(t, err, err)

		data, err := testFile.MarshalBinary()
		require.Nil(t, err, err)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.Header.SetMethod("POST")
		req.Header.Add("Content-Type", "application/json")
		req.SetRequestURI(route)
		req.SetBody(data)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		err = client.Do(req, resp)
		require.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
		assert.Equal(t, string(resp.Body()), testFile.Meta().GetID().String())

		key := string(resp.Body())

		req = fasthttp.AcquireRequest()
		req.SetRequestURI(route)
		req.Header.SetMethod("GET")

		req.URI().QueryArgs().Add("key", key)

		resp = fasthttp.AcquireResponse()

		err = client.Do(req, resp)
		require.Nil(t, err, err)

		dt, err := models.DataTypeFromID([]byte(key))
		require.Nil(t, err, err)

		obtainedData, err := models.NewData(
			models.NewMetaData("simple", dt),
			[]byte(""),
		)
		require.Nil(t, err, err)

		err = obtainedData.UnmarshalBinary(resp.Body())
		require.Nil(t, err, err)
		require.Equal(t, testFile, obtainedData)
		require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	})

	t.Run("FileDelete", func(t *testing.T) {
		testFile, err := models.NewData(
			models.NewMetaData("simple", models.Simple),
			randomContent(),
		)
		require.Nil(t, err, err)

		data, err := testFile.MarshalBinary()
		require.Nil(t, err, err)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.Header.SetMethod("POST")
		req.Header.Add("Content-Type", "application/json")
		req.SetRequestURI(route)
		req.SetBody(data)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		err = client.Do(req, resp)
		require.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
		assert.Equal(t, string(resp.Body()), testFile.Meta().GetID().String())

		key := resp.Body()

		req = fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.SetRequestURI(route)
		req.Header.SetMethod("GET")

		req.URI().QueryArgs().Add("key", string(key))

		resp = fasthttp.AcquireResponse()

		err = client.Do(req, resp)
		require.Nil(t, err, err)

		dt, err := models.DataTypeFromID(key)
		require.Nil(t, err, err)

		obtainedData, err := models.NewData(
			models.NewMetaData("simple", dt),
			[]byte(""),
		)
		require.Nil(t, err, err)

		err = obtainedData.UnmarshalBinary(resp.Body())
		require.Nil(t, err, err)
		require.Equal(t, testFile, obtainedData)
		require.Equal(t, fasthttp.StatusOK, resp.StatusCode())

		req = fasthttp.AcquireRequest()
		req.SetRequestURI(route)
		req.Header.SetMethod("DELETE")

		req.URI().QueryArgs().Add("key", string(key))

		resp = fasthttp.AcquireResponse()

		err = client.Do(req, resp)
		require.Nil(t, err, err)
		require.Equal(t, resp.Body(), key)
		require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
	})

	t.Run("AddSame", func(t *testing.T) {
		content := randomContent()
		var testFile = models.NewSimpleData(
			models.NewMetaData("simple", models.Simple),
			content,
		)

		var sameFile = models.NewSimpleData(
			&models.MetaData{
				Title:    "test",
				ID:       testFile.ID,
				DataType: models.Simple,
			},
			content,
		)

		data, err := testFile.MarshalBinary()
		require.Nil(t, err, err)

		req := fasthttp.AcquireRequest()
		defer fasthttp.ReleaseRequest(req)

		req.Header.SetMethod("POST")
		req.Header.SetContentType("application/json")
		req.SetRequestURI(route)
		req.SetBody(data)

		resp := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(resp)

		err = client.Do(req, resp)
		require.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
		assert.Equal(t, string(resp.Body()), testFile.Meta().GetID().String())

		data, err = sameFile.MarshalBinary()
		require.Nil(t, err, err)
		req.SetBody(data)

		resp = fasthttp.AcquireResponse()

		err = client.Do(req, resp)
		require.Nil(t, err, err)
		assert.Equal(t, fasthttp.StatusBadRequest, resp.StatusCode())
		assert.NotNil(t, resp.Body())
	})
}

func TestConfig(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		c := Config{
			StoragePath: ".badger",
			ServerName:  "DSSS",
			ServerHost:  "localhost",
			ServerPort:  ":8080",
		}

		_, err := c.Save(".")
		require.Nil(t, err, err)
	})

	t.Run("Load", func(t *testing.T) {
		c := Config{
			StoragePath: ".badger",
			ServerName:  "DSSS",
			ServerHost:  "localhost",
			ServerPort:  ":8080",
		}

		p, err := c.Save(".")
		require.Nil(t, err, err)

		c2 := Config{}

		err = c2.Load(p)
		require.Nil(t, err, err)
		assert.Equal(t, c, c2)
	})
}

func randomContent() []byte {
	b := make([]byte, 4)
	rand.Read(b)

	return b
}
