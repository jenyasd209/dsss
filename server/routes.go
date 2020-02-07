package server

import (
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"log"

	"github.com/iorhachovyevhen/dsss/models"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var ErrorWrongID = errors.New("id is wrong")

func router() *routing.Router {
	log.Println("Create router...")

	router := routing.New()

	files := router.Group("/files")
	files.Post("", addFile)
	files.Get("/<hash>", getFile).
		Delete(deleteFile)

	return router
}

func getFile(ctx *routing.Context) error {
	key := ctx.QueryArgs().Peek("key")
	if string(key) == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return ErrorWrongID
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

	ctx.SetContentType("application/octet-stream")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)

	return nil
}

func addFile(ctx *routing.Context) error {
	file := ctx.Request.Body()
	dataType := gjson.Get(string(file), "data_type")
	data := models.NewEmptyData(models.DataType(dataType.Uint()))

	err := data.UnmarshalBinary(file)
	if err != nil {
		return err
	}

	key, err := storage.Add(data)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	ctx.SetContentType("text/plain")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBodyString(string(key))

	return nil
}

func deleteFile(ctx *routing.Context) error {
	key := ctx.QueryArgs().Peek("key")
	if string(key) == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return ErrorWrongID
	}

	err := storage.Delete(key)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return err
	}

	ctx.SetContentType("text/plain")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(string(key))

	return nil
}
