package server

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"simpleRestCache/pkg/config"
	"simpleRestCache/pkg/service"
)

// Server represents HTTP Server
type Server struct {
	server http.Server
}

// New creates an instance of the Server
func New(cfg *config.Config, service *service.Service) *Server {
	s := &Server{}

	m := NewHandler(cfg, service)

	s.server = http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      m,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Info("Presentation subsystem has been initialized")

	return s
}

// Run starts HTTP Server
func (s *Server) Run() {
	log.Info("HTTP server started ", s.server.Addr)
	s.server.ListenAndServe()
}
