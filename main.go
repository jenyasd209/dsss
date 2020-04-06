package main

import (
	"github.com/iorhachovyevhen/dsss/server"
	"log"
)

func main() {
	storageServer := server.NewStorageServer()
	if err := storageServer.Start(); err != nil {
		log.Fatalf("Error starting %s", err)
	}
}
