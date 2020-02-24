package server

import (
	"errors"
	"github.com/iorhachovyevhen/dsss/storage"
	"log"

	"github.com/iorhachovyevhen/dsss/models"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var ErrorBadID = errors.New("wrong id")
var ErrorBadJSON = errors.New("bad json")

func InitRouter(storage storage.DataKeeper) *Router {
	log.Println("Set up router...")

	router := &Router{
		routing.New(),
		storage,
	}

	InitFileRouter(router)

	return router
}

type Router struct {
	*routing.Router
	storage storage.DataKeeper
}

func InitFileRouter(router *Router) {
	log.Println("Init file router...")

	fr := FileRouter{router}

	fr.Group("/file").
		Post("", fr.Post).
		Get(fr.Get).
		Delete(fr.Delete)
}

type FileRouter struct {
	*Router
}

//Post - file endpoint that add a new file
func (r *FileRouter) Post(ctx *routing.Context) error {
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

	key, err := r.storage.Add(data)
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

//Get - file endpoint that add a new file
func (r *FileRouter) Get(ctx *routing.Context) error {
	key := ctx.QueryArgs().Peek("key")
	if key == nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(ErrorBadID.Error()),
		)
		return nil
	}

	dt, err := models.DataTypeFromID(key)
	if err != nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(err.Error()),
		)
		return nil
	}

	data := models.NewEmptyData(dt)

	err = r.storage.Read(key, data)
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

//Delete - file endpoint that add a new file
func (r *FileRouter) Delete(ctx *routing.Context) error {
	key := ctx.QueryArgs().Peek("key")
	if key == nil {
		makeResponse(&ctx.Response,
			fasthttp.StatusBadRequest,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(ErrorBadID.Error()),
		)
		return nil
	}

	err := r.storage.Delete(key)
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
