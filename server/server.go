package server

import (
	db "github.com/iorhachovyevhen/dsss/storage"
	"github.com/valyala/fasthttp"
)

//var storage db.DataKeeper

type Server interface {
	Start() error
	Shutdown() error
}

func NewStorageServer() *StorageServer {
	return NewStorageServerWithConfig(DefaultConfig())
}

func NewStorageServerWithConfig(config *Config) *StorageServer {
	return newStorageServer(config)
}

func newStorageServer(config *Config) *StorageServer {
	storage := db.NewDefaultStorage(config.StoragePath)
	globalRouter := InitRouter(storage)

	return &StorageServer{
		storage: storage,
		server: &fasthttp.Server{
			Name:    config.ServerName,
			Handler: globalRouter.HandleRequest,
		},
		host: config.ServerHost,
		port: config.ServerPort,
	}
}

type StorageServer struct {
	storage db.DataKeeper
	server  *fasthttp.Server

	host string
	port string
}

func (s *StorageServer) Start() error {
	return s.server.ListenAndServe(s.host + s.port)
}

func (s *StorageServer) Shutdown() error {
	return s.server.Shutdown()
}
