package main

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/iorhachovyevhen/dsss/models"
)

var route = "http://127.0.0.1:8080/files"

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
	req.SetRequestURI(route)
	req.Header.SetMethod("PUT")
	req.SetBody(data)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)
	require.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
}

func TestGet(t *testing.T) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(route + "/")
	req.Header.SetMethod("GET")

	req.URI().QueryArgs().Add("id", file.ID().String())
	req.URI().QueryArgs().Add("type", string(file.Type()))

	resp := fasthttp.AcquireResponse()

	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	require.Nil(t, err, err)

	obtainedData := models.NewSimpleData(models.MetaData{}, nil)
	err = obtainedData.UnmarshalBinary(resp.Body())
	require.Nil(t, err, err)
	require.Equal(t, file, obtainedData)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func TestDelete(t *testing.T) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(route + "/")
	req.Header.SetMethod("DELETE")

	req.URI().QueryArgs().Add("id", file.ID().String())
	req.URI().QueryArgs().Add("type", string(file.Type()))

	resp := fasthttp.AcquireResponse()

	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	require.Nil(t, err, err)

	id := string(resp.Body())
	require.Equal(t, file.ID().String(), id)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())

	req.SetRequestURI(route + "/")
	req.Header.SetMethod("GET")

	req.URI().QueryArgs().Add("id", file.ID().String())
	req.URI().QueryArgs().Add("type", string(file.Type()))

	err = client.Do(req, resp)
	println(resp.String())
	require.Nil(t, err, err)
	require.Equal(t, fasthttp.StatusNotFound, resp.StatusCode())
}
