package server

import (
	"errors"
	"github.com/iorhachovyevhen/dsss/models"
	db "github.com/iorhachovyevhen/dsss/storage"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
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
		return nil
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
	fileType := ctx.QueryArgs().Peek("type")
	if len(fileType) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return nil
	}
	data := models.NewEmptyData(db.ByteToDataType(fileType))

	file := ctx.PostBody()
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
