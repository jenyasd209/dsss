package server

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var ErrEmptyStoragePath = errors.New("storage path is empty")
var ErrEmptyServerHost = errors.New("server host is empty")
var ErrEmptyServerPort = errors.New("server port is empty")

func DefaultConfig() *Config {
	return &Config{
		StoragePath: "./badger",
		ServerName:  "DSSS",
		ServerHost:  "localhost",
		ServerPort:  ":8080",
	}
}

type Config struct {
	StoragePath string
	ServerName  string
	ServerHost  string
	ServerPort  string
}

func (c *Config) Load(path string) error {
	log.Println("Reading config...")

	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Config reading finished with error: %s\n", err)
		return err
	}

	serverCfg := viper.GetStringMapString("server")
	c.ServerHost = serverCfg["host"]
	c.ServerPort = serverCfg["port"]
	c.ServerName = serverCfg["name"]
	c.StoragePath = viper.GetString("storage")

	return c.check()
}

func (c *Config) Save(folder string) (path string, err error) {
	log.Println("Saving config...")

	err = c.check()
	if err != nil {
		return
	}

	if _, err = os.Stat(folder); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(folder, os.ModePerm)
			if err != nil {
				log.Fatalf("can't create dir by path: %v: %v", folder, err)
			}
		}

		return
	}

	viper.SetDefault("storage", c.StoragePath)
	viper.SetDefault("server", map[string]string{
		"host": c.ServerHost,
		"port": c.ServerPort,
		"name": c.ServerName,
	})

	path = filepath.Join(folder, "config.yml")

	err = viper.WriteConfigAs(path)
	if err != nil {
		panic(err)
	}

	return
}

func (c *Config) check() error {
	if c.StoragePath == "" {
		return ErrEmptyStoragePath
	}

	if c.ServerHost == "" {
		return ErrEmptyServerHost
	}

	if c.ServerPort == "" {
		return ErrEmptyServerPort
	}

	if c.ServerName == "" {
		c.ServerName = "DSSS"
	}

	return nil
}
