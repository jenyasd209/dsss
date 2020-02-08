package main

import (
	"github.com/iorhachovyevhen/dsss/cmd/dsss/cli"
	"log"
	"os"
)

func main() {
	err := cli.RootApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
