package api

import (
	"github.com/iorhachovyevhen/dsss/models"
	"github.com/iorhachovyevhen/dsss/storage"
	"github.com/valyala/fasthttp"
)

func API(address string) *api {
	return &api{
		address: address,
		client: fasthttp.Client{
			Name: "API",
		},
	}
}

type api struct {
	address string
	client  fasthttp.Client
}

type fileRoute struct {
	*api
	route string
}

func (a *api) Files() *fileRoute {
	return &fileRoute{
		api:   a,
		route: a.address + "/files",
	}
}

func (f *fileRoute) Add(fileName string, body []byte) (string, error) {
	resp, err := f.doRequest(f.route, "POST", body, nil)
	if err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}

func (f *fileRoute) Get(key []byte) (models.Data, error) {
	resp, err := f.doRequest(f.route, "GET", nil, map[string][]byte{"key": key})
	if err != nil {
		return nil, err
	}

	dt := storage.DataTypeFromKey(key)

	obj := models.NewEmptyData(dt)
	err = obj.UnmarshalBinary(resp.Body())
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (f *fileRoute) Delete(key []byte) (fileName string, err error) {
	resp, err := f.doRequest(f.route, "DELETE", nil, map[string][]byte{"key": key})
	if err != nil {
		return "", err
	}

	return string(resp.Body()), nil
}

func (a *api) doRequest(addr string, method string, body []byte, args map[string][]byte) (*fasthttp.Response, error) {
	req := &fasthttp.Request{}
	resp := &fasthttp.Response{}

	req.SetRequestURI(addr)
	req.Header.SetMethod(method)
	req.SetBody(body)

	for k, v := range args {
		req.URI().QueryArgs().AddBytesV(k, v)
	}

	err := a.client.Do(req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
