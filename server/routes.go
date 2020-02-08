package server

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"

	"github.com/iorhachovyevhen/dsss/models"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var ErrorBadID = errors.New("wrong id")
var ErrorBadJSON = errors.New("bad json")
var ErrorBadDataType = errors.New("Bad 'data_type' value")

func router() *routing.Router {
	log.Println("Create router...")

	router := routing.New()
	router.Group("/file").
		Post("", addFile).
		Get(getFile).
		Delete(deleteFile)

	return router
}

func addFile(ctx *routing.Context) error {
	value, err := jsonValue(ctx.Request.Body(), "data_type")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	dt, err := convertToDataType(value)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	data := models.NewEmptyData(dt)

	err = data.UnmarshalBinary(ctx.Request.Body())
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return ErrorBadJSON
	}

	key, err := storage.Add(data)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	makeResponse(&ctx.Response, fasthttp.StatusCreated, map[string]string{"Content-Type": "text/plain"}, key)

	return nil
}

func getFile(ctx *routing.Context) error {
	key := ctx.QueryArgs().Peek("key")
	if key == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return ErrorBadID
	}

	data, err := storage.Read(key)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return err
	}

	body, err := data.MarshalBinary()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	makeResponse(&ctx.Response, fasthttp.StatusOK, map[string]string{"Content-Type": "application/json"}, body)

	return nil
}

func deleteFile(ctx *routing.Context) error {
	key := ctx.QueryArgs().Peek("key")
	if key == nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return ErrorBadID
	}

	err := storage.Delete(key)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return err
	}

	makeResponse(&ctx.Response, fasthttp.StatusOK, map[string]string{"Content-Type": "text/plain"}, key)

	return nil
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
