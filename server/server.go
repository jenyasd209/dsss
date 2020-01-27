package main

import (
	"errors"
	db "github.com/iorhachovyevhen/dsss/storage"
	"log"
	"os/user"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const addr = "127.0.0.1:8080"

var ErrorWrongID = errors.New("id is wrong")

var storage db.DataKeeper

type StorageServer interface {
	Start() error
	Shutdown() error
}

type BadgerServer struct {
	storage db.DataKeeper
	server  *fasthttp.Server
}

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	storage = db.NewStorageWithOptions(
		db.NewOptions().WithDir(usr.HomeDir + "/.badger").
			WithValueDir(usr.HomeDir + "/.badger").
			WithValueLogFileSize(2 << 20),
	)

	router := router()

	server := &fasthttp.Server{
		Name:    "DSSS server",
		Handler: router.HandleRequest,
	}

	log.Println("Server starts...")
	if err := server.ListenAndServe(addr); err != nil {
		log.Fatalf("Error listening %s", err)
	}
}

func router() *routing.Router {
	router := routing.New()

	files := router.Group("/files")
	files.Put("", addFile)
	files.Get("/<hash>", getFile).
		Delete(deleteFile)

	return router
}

func getFile(ctx *routing.Context) error {
	dataType := ctx.QueryArgs().Peek("type")
	data := newData(dataType)

	idArg := ctx.QueryArgs().Peek("id")
	if string(idArg) == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return ErrorWrongID
	}

	id, err := byteToHash32(idArg)
	if err != nil {
		return err
	}

	err = storage.Read(id, data)
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
	fileType := ctx.PostArgs().Peek("type")
	data := newData(fileType)

	file := ctx.PostBody()
	err := data.UnmarshalBinary(file)
	if err != nil {
		return err
	}

	err = storage.Add(data)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return err
	}

	ctx.SetContentType("text/plain")
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBodyString(data.ID().String())

	return nil
}

func deleteFile(ctx *routing.Context) error {
	dataType := byteToDataType(ctx.QueryArgs().Peek("type"))

	idArg := ctx.QueryArgs().Peek("id")
	if string(idArg) == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return ErrorWrongID
	}

	id, err := byteToHash32(idArg)
	if err != nil {
		return err
	}

	err = storage.Delete(id, dataType)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return err
	}

	ctx.SetContentType("text/plain")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString(id.String())

	return nil
}
