package main

import (
	db "github.com/iorhachovyevhen/dsss/storage"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const addr = "127.0.0.1:8080"

var cfgFile = filepath.Join(os.Getenv("HOME"), ".dsss/config.yml")
var storage db.DataKeeper

func main() {
	if err := startServer(); err != nil {
		log.Fatalf("Error starting %s", err)
	}
}

func startServer() error {
	log.Println("Start setting config...")
	err := readConfig()
	if err != nil {
		return err
	}
	log.Println("Config are set.")

	log.Println("Initiate storage...")
	storage = initStorage()
	log.Println("Storage initiated.")

	router := router()

	server := &fasthttp.Server{
		Name:    viper.GetString("name"),
		Handler: router.HandleRequest,
	}

	log.Println("Server starts...")
	return server.ListenAndServe(addr)
}

func readConfig() error {
	log.Println("Reading config...")

	viper.SetConfigFile(cfgFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Config reading finished with error: %s\n", err)
		return defaultConfig(cfgFile)
	}

	log.Println("Config are red.")
	return nil
}

func defaultConfig(path string) error {
	log.Println("Setting default config...")

	storagePath := filepath.Join(os.Getenv("HOME"), ".dsss/storage")
	viper.SetDefault("storage", storagePath)
	viper.SetDefault("server", map[string]string{
		"host": "localhost",
		"port": "8080",
		"name": "DSSS",
	})

	err := viper.WriteConfigAs(path)
	if err != nil {
		log.Printf("Can't set default config: %s\n", err)
		return defaultConfig(cfgFile)
	}

	return nil
}

func initStorage() db.DataKeeper {
	return db.NewStorageWithOptions(
		db.NewOptions().WithDir(viper.GetString("storage")).
			WithValueDir(viper.GetString("storage")).
			WithValueLogFileSize(2 << 20),
	)
}

func router() *routing.Router {
	router := routing.New()

	files := router.Group("/files")
	files.Put("", addFile)
	files.Get("/<hash>", getFile).
		Delete(deleteFile)

	return router
}
