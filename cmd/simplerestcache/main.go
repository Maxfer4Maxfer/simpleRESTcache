package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"simpleRestCache/pkg/config"
	control "simpleRestCache/pkg/server/grpc"
	server "simpleRestCache/pkg/server/http"

	group "github.com/oklog/run"

	service "simpleRestCache/pkg/service"
	storage "simpleRestCache/pkg/storage/gorm"
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
	control := control.New(cfg, service)

	var g group.Group
	{
		// start HTTP Server - presentation layer
		g.Add(func() error {
			return server.Run()
		}, func(error) {
			server.Close()
		})
	}
	{
		// start gRPC Server - control layer
		g.Add(func() error {
			return control.Run()
		}, func(error) {
			control.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	g.Run()
}
