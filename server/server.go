package server

import (
	db "github.com/iorhachovyevhen/dsss/storage"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"log"
)

var storage db.DataKeeper

func StartServer(config string) error {
	log.Println("Setting config...")

	err := readConfig(config)
	if err != nil {
		return err
	}

	storage = initStorage()

	serverCfg := viper.GetStringMapString("server")
	addr := serverCfg["host"] + serverCfg["port"]

	server := &fasthttp.Server{
		Name:    viper.GetString("name"),
		Handler: router().HandleRequest,
	}

	log.Printf("Server listening on %s...", addr)
	return server.ListenAndServe(addr)
}

func readConfig(path string) error {
	log.Println("Reading config...")

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Config reading finished with error: %s\n", err)
		return err
	}

	return nil
}

func initStorage() db.DataKeeper {
	log.Println("Initiate storage...")

	return db.NewStorageWithOptions(
		db.NewOptions().WithDir(viper.GetString("storage")).
			WithValueDir(viper.GetString("storage")).
			WithValueLogFileSize(2 << 20),
	)
}
