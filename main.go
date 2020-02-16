package main

import (
	"github.com/iorhachovyevhen/dsss/server"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var (
	folder  = ".dsss"
	cfg     = "default_config.yml"
	storage = "storage"
)

func main() {
	if err := server.StartServer(defaultConfig()); err != nil {
		log.Fatalf("Error starting %s", err)
	}
}

func defaultConfig() string {
	folder = filepath.Join(os.Getenv("HOME"), folder)

	if _, err := os.Stat(folder); err != nil {
		err = os.MkdirAll(folder, os.ModePerm)
		if err != nil {
			log.Fatalf("can't create dir by path: %v: %v", folder, err)
		}
	}

	defaultCfg := filepath.Join(folder, cfg)

	_, err := os.OpenFile(defaultCfg, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	log.Println("Setting default config...")

	storagePath := filepath.Join(folder, storage)
	if _, err := os.Stat(storage); err != nil {
		err = os.MkdirAll(storage, os.ModePerm)
		if err != nil {
			log.Fatalf("can't create dir by path: %v: %v", storage, err)
		}
	}

	viper.SetDefault("storage", storagePath)
	viper.SetDefault("server", map[string]string{
		"host": "192.168.43.178",
		"port": ":8080",
		"name": "DSSS",
	})

	err = viper.WriteConfigAs(defaultCfg)
	if err != nil {
		panic(err)
	}

	return defaultCfg
}
