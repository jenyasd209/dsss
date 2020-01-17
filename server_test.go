package main

import (
	"fmt"
	"testing"

	"dsss/models"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestAdd(t *testing.T) {
	file := models.SimpleData{Content:[]byte("content")}
	client := &fasthttp.Client{}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.URI().Update("127.0.0.1:8080")
	req.Header.SetMethodBytes([]byte("PUT"))

	data, err := file.MarshalBinary()
	assert.Nil(t, err, err)
	req.AppendBody(data)
	fmt.Println(string(req.Body()))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client.Do(req, resp)
}