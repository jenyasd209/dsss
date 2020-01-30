package main

import (
	"errors"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var ErrorWrongID = errors.New("id is wrong")

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
