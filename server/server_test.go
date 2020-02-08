package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/iorhachovyevhen/dsss/models"
)

var route = "http://localhost:8080/file"

func TestAdd(t *testing.T) {
	var testFile = models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Audio,
		},
		randomContent(),
	)

	data, err := testFile.MarshalBinary()
	require.Nil(t, err, err)

	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.SetRequestURI(route)
	req.SetBody(data)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)
	assert.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
	assert.Equal(t, string(resp.Body()), testFile.ID().String())
}

func TestGet(t *testing.T) {
	var testFile = models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		randomContent(),
	)

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

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)
	assert.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
	assert.Equal(t, string(resp.Body()), testFile.ID().String())

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

	obtainedData := models.NewEmptyData(dt)
	err = obtainedData.UnmarshalBinary(resp.Body())
	require.Nil(t, err, err)
	require.Equal(t, testFile, obtainedData)
	require.Equal(t, fasthttp.StatusOK, resp.StatusCode())
}

func TestDelete(t *testing.T) {
	var testFile = models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Simple,
		},
		randomContent(),
	)

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

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)
	assert.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
	assert.Equal(t, string(resp.Body()), testFile.ID().String())

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

	obtainedData := models.NewEmptyData(dt)
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
}

func TestAddSame(t *testing.T) {
	content := randomContent()
	var testFile = models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Audio,
		},
		content,
	)

	var sameFile = models.NewSimpleData(
		models.MetaData{
			Title:    "test",
			DataType: models.Audio,
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

	client := &fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)
	assert.Equal(t, fasthttp.StatusCreated, resp.StatusCode())
	assert.Equal(t, string(resp.Body()), testFile.ID().String())

	data, err = sameFile.MarshalBinary()
	require.Nil(t, err, err)
	req.SetBody(data)

	resp = fasthttp.AcquireResponse()

	err = client.Do(req, resp)
	require.Nil(t, err, err)
	assert.Equal(t, fasthttp.StatusBadRequest, resp.StatusCode())
	assert.NotNil(t, resp.Body())
}

func randomContent() []byte {
	b := make([]byte, 4)
	rand.Read(b)

	return b
}
