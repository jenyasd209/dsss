package server

import (
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/iorhachovyevhen/dsss/models"
	"github.com/valyala/fasthttp"
)

func dataTypeFromJSON(j []byte) (models.DataType, error) {
	value, err := jsonValue(j, "data_type")
	if err != nil {
		return 0, err
	}

	return models.ConvertToDataType(value)
}

func makeResponse(resp *fasthttp.Response, statusCode int, headerArgs map[string]string, body []byte) {
	resp.SetStatusCode(statusCode)

	for k, v := range headerArgs {
		resp.Header.Add(k, v)
	}

	resp.SetBody(body)
}

func jsonValue(body []byte, key string) (interface{}, error) {
	var j map[string]interface{}

	err := json.Unmarshal(body, &j)
	if err != nil {
		return nil, err
	}

	value, ok := j[key]
	if !ok {
		return nil, errors.Errorf("bad key: %v", key)
	}

	return value, nil
}
