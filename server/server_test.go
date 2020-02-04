package server

import (
	"fmt"
	db "github.com/iorhachovyevhen/dsss/storage"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/iorhachovyevhen/dsss/models"
)

var route = "http://localhost:8080/files"

var file = models.NewSimpleData(
	models.MetaData{
		Title:    "test",
		DataType: models.Simple,
	},
	[]byte("content"),
)

func TestAdd(t *testing.T) {
	data, err := file.MarshalBinary()
	require.Nil(t, err, err)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")

	req.SetRequestURI(route)
	req.URI().QueryArgs().Add("type", string(file.Type()))

	req.SetBody(data)

	resp := fasthttp.AcquireResponse()

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)
	require.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
	require.NotNil(t, resp.Body())
}

func TestGet(t *testing.T) {
	var testFile = models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)

	data, err := testFile.MarshalBinary()
	require.Nil(t, err, err)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")

	req.SetRequestURI(route)
	req.URI().QueryArgs().Add("type", string(testFile.Type()))

	req.SetBody(data)

	resp := fasthttp.AcquireResponse()

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)

	key := string(resp.Body())

	req = fasthttp.AcquireRequest()
	req.SetRequestURI(route + "/")
	req.Header.SetMethod("GET")

	req.URI().QueryArgs().Add("key", key)

	resp = fasthttp.AcquireResponse()

	err = client.Do(req, resp)
	require.Nil(t, err, err)

	obtainedData := models.NewEmptyData(db.DataTypeFromKey([]byte(key)))
	err = obtainedData.UnmarshalBinary(resp.Body())
	fmt.Println(resp.Body())
	require.Nil(t, err, err)
	require.Equal(t, file, obtainedData)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func TestDelete(t *testing.T) {
	var testFile = models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)

	data, err := testFile.MarshalBinary()
	require.Nil(t, err, err)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")

	req.SetRequestURI(route)
	req.URI().QueryArgs().Add("type", string(testFile.Type()))

	req.SetBody(data)

	resp := fasthttp.AcquireResponse()

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)
	require.NotNil(t, resp.Body())

	key := string(resp.Body())

	req = fasthttp.AcquireRequest()
	req.SetRequestURI(route + "/")
	req.Header.SetMethod("DELETE")

	req.URI().QueryArgs().Add("key", key)

	resp = fasthttp.AcquireResponse()

	err = client.Do(req, resp)
	require.Nil(t, err, err)

	deletedKey := string(resp.Body())
	require.Equal(t, key, deletedKey)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())

	req.SetRequestURI(route + "/")
	req.Header.SetMethod("GET")

	req.URI().QueryArgs().Add("key", key)

	err = client.Do(req, resp)
	require.Nil(t, err, err)
	require.Equal(t, fasthttp.StatusNotFound, resp.StatusCode())
}
