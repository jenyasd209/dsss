package api

import (
	"github.com/iorhachovyevhen/dsss/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"testing"
)

var addr = "http://localhost:8080"

var dsssAPI = API(addr)

func TestDoRequest(t *testing.T) {
	data := models.NewEmptyData(1)

	body, err := data.MarshalBinary()
	require.Nil(t, err, err)

	req := &fasthttp.Request{}
	defer fasthttp.ReleaseRequest(req)

	resp := &fasthttp.Response{}
	req.Header.SetContentType("application/json")

	req.SetRequestURI(addr + "/files")
	req.Header.SetMethod("POST")
	req.SetBody(body)

	client := fasthttp.Client{}
	err = client.Do(req, resp)
	require.Nil(t, err, err)
	assert.NotNil(t, resp.Body())
}
