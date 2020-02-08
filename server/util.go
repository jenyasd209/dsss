package server

import (
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/iorhachovyevhen/dsss/models"
	db "github.com/iorhachovyevhen/dsss/storage"
	"github.com/valyala/fasthttp"
)

func DataTypeFromKey(key []byte) (models.DataType, error) {
	return db.DataTypeFromKey(key)
}

func DataTypeFromJSON(j []byte) (models.DataType, error) {
	value, err := jsonValue(j, "data_type")
	if err != nil {
		return 0, err
	}

	return convertToDataType(value)
}

func makeResponse(resp *fasthttp.Response, statusCode int, contentTypes map[string]string, body []byte) {
	resp.SetStatusCode(statusCode)

	for k, v := range contentTypes {
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

func convertToDataType(value interface{}) (models.DataType, error) {
	dt, ok := value.(float64)
	if !ok {
		return models.DataType(dt), ErrorBadDataType
	}
	return models.DataType(dt), nil
}
