package main

import (
	log "github.com/sirupsen/logrus"

	"simpleRestCache/pkg/config"
	server "simpleRestCache/pkg/server/http"

	service "simpleRestCache/pkg/service"
	storage "simpleRestCache/pkg/storage/inmem"
)

func main() {

	cfg := config.GetConfig()

	// setup logger
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{})
	if cfg.Debug {
		log.SetReportCaller(true)
		log.Warn("Debug mode is activated")
	}

	storage := storage.New(cfg)
	defer storage.Close()

	service := service.New(cfg, storage)

	server := server.New(cfg, service)

	server.Run()
}
