package server

import (
	"github.com/pkg/errors"
	"log"

	"github.com/iorhachovyevhen/dsss/models"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var ErrorBadID = errors.New("wrong id")
var ErrorBadJSON = errors.New("bad json")

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
	dt, err := dataTypeFromJSON(ctx.Request.Body())
	if err != nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()),
		)
		return nil
	}

	data := models.NewEmptyData(dt)

	err = data.UnmarshalBinary(ctx.Request.Body())
	if err != nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(ErrorBadJSON.Error()),
		)
		return nil
	}

	key, err := storage.Add(data)
	if err != nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()),
		)
		return nil
	}

	makeResponse(&ctx.Response, fasthttp.StatusCreated, map[string]string{"Content-Type": "text/plain"}, key)

	return nil
}

func getFile(ctx *routing.Context) error {
	key := ctx.QueryArgs().Peek("key")
	if key == nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(ErrorBadID.Error()),
		)
		return nil
	}

	data, err := storage.Read(key)
	if err != nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusNotFound,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()),
		)
		return nil
	}

	body, err := data.MarshalBinary()
	if err != nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()),
		)
		return nil
	}

	makeResponse(&ctx.Response, fasthttp.StatusOK, map[string]string{"Content-Type": "application/json"}, body)

	return nil
}

func deleteFile(ctx *routing.Context) error {
	key := ctx.QueryArgs().Peek("key")
	if key == nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(ErrorBadID.Error()),
		)
		return nil
	}

	err := storage.Delete(key)
	if err != nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()),
		)
		return nil
	}

	makeResponse(&ctx.Response, fasthttp.StatusOK, map[string]string{"Content-Type": "text/plain"}, key)

	return nil
}
