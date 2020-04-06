package api

import (
	"github.com/iorhachovyevhen/dsss/models"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func API(address string) *api {
	return &api{
		address: address,
		client: fasthttp.Client{
			Name: "DSSS API",
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
		route: a.address + "/file",
	}
}

func (f *fileRoute) Add(fileName string, content []byte) (models.ID, error) {
	dt := DataTypeFromFilename(fileName)
	data := models.NewData(dt, fileName, content)

	body, err := data.MarshalBinary()
	if err != nil {
		return nil, err
	}

	resp, err := f.doRequest(f.route, "POST", body, nil)
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusCreated {
		return nil, errors.Errorf("%v: %s", resp.StatusCode(), resp.Body())
	}

	return resp.Body(), nil
}

func (f *fileRoute) Get(key models.ID) (models.Data, error) {
	resp, err := f.doRequest(f.route, "GET", nil, map[string][]byte{"key": key})
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, errors.Errorf("%v: %s", resp.StatusCode(), resp.Body())
	}

	dt, err := models.DataTypeFromID(key)
	if err != nil {
		return nil, err
	}

	obj := models.NewEmptyData(dt)
	err = obj.UnmarshalBinary(resp.Body())
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (f *fileRoute) Delete(key models.ID) ([]byte, error) {
	resp, err := f.doRequest(f.route, "DELETE", nil, map[string][]byte{"key": key})
	if err != nil {
		return nil, err
	}
	defer fasthttp.ReleaseResponse(resp)

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, errors.Errorf("%v: %s", resp.StatusCode(), resp.Body())
	}

	return resp.Body(), nil
}

func (a *api) doRequest(addr string, method string, body []byte, args map[string][]byte) (*fasthttp.Response, error) {
	req := &fasthttp.Request{}
	defer fasthttp.ReleaseRequest(req)

	resp := &fasthttp.Response{}

	req.SetRequestURI(addr)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/json")
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
