package main

import (
	"dsss/models"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"log"
)

const addr  = "127.0.0.1:8080"

func main() {
	router := routing.New()

	files := router.Group("/files")

	files.Get("", listFiles).Post(addFile)
	files.Get("/<hash>", getFile).Delete(deleteFile)

	server := &fasthttp.Server{
		Name: "DSSS server",
		Handler: router.HandleRequest,
	}

	log.Println("Server starts...")
	if err := server.ListenAndServe(addr); err != nil{
		log.Fatalf("Error listening %s", err)
	}
}

func listFiles(ctx *routing.Context) error {
	ctx.Response.Header.Add("Content-Type", "text/html")
	ctx.Response.AppendBody([]byte("yep"))

	return nil
}

func getFile(ctx *routing.Context) error {
	return nil
}

func addFile(ctx *routing.Context) error {
	file := new(models.SimpleData)

	err := file.UnmarshalBinary(ctx.Request.Body())
	if err != nil{
		log.Println(err)
		return err
	}

	return nil
}

func deleteFile(context *routing.Context) error {
	return nil
}