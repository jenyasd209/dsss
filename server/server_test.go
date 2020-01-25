package main

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/iorhachovyevhen/dsss/models"
)

var route = "http://127.0.0.1:8080/files"

func TestAdd(t *testing.T) {
	file := models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		[]byte("content"),
	)

	data, err := file.MarshalBinary()
	require.Nil(t, err, err)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(route)
	req.Header.Add("User-Agent", "Test-Agent")
	req.Header.SetMethod("POST")
	req.SetBody(data)

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)

	req.SetRequestURI(route + "/" + file.ID().String())
	req.Header.SetMethod("GET")
	req.SetBody(data)

	println(req.Header.String())

	err = client.Do(req, resp)
	require.Nil(t, err, err)
}
