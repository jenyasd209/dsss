package main

import (
	"github.com/iorhachovyevhen/dsss/server"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if err := server.StartServer(defaultConfig()); err != nil {
		log.Fatalf("Error starting %s", err)
	}
}

func defaultConfig() string {
	defaultCfg := filepath.Join(os.Getenv("HOME"), ".dsss/default_config.yml")

	log.Println("Setting default config...")

	storagePath := filepath.Join(os.Getenv("HOME"), ".dsss/storage")
	viper.SetDefault("storage", storagePath)
	viper.SetDefault("server", map[string]string{
		"host": "localhost",
		"port": ":8080",
		"name": "DSSS",
	})

	err := viper.WriteConfigAs(defaultCfg)
	if err != nil {
		panic(err)
	}

	return defaultCfg
}
