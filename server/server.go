package main

import (
	db "github.com/iorhachovyevhen/dsss/storage"
	"log"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const addr = "127.0.0.1:8080"

var storage db.DataKeeper

func main() {
	storage = db.NewStorageWithOptions(
		db.NewOptions().WithValueDir("/home/$USER/badger/value").
			WithDir("/home/$USER/badger").
			WithValueLogFileSize(2 << 20),
	)

	router := routing.New()

	files := router.Group("/files")

	files.Post("", addFile)
	files.Get("/<hash>", getFile).Delete(deleteFile)

	server := &fasthttp.Server{
		Name:    "DSSS server",
		Handler: router.HandleRequest,
	}

	log.Println("Server starts...")
	if err := server.ListenAndServe(addr); err != nil {
		log.Fatalf("Error listening %s", err)
	}
}

func getFile(ctx *routing.Context) error {
	return nil
}

func addFile(ctx *routing.Context) error {
	return nil
}

func deleteFile(context *routing.Context) error {
	return nil
}
