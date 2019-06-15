package grpc

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"simpleRestCache/pkg/config"
	"simpleRestCache/pkg/service"

	"simpleRestCache/pb"
)

// Server represents gRPC Server
type Server struct {
	cfg    *config.Config
	server *grpc.Server
}

// New creates an instance of the Server
func New(cfg *config.Config, srv *service.Service) *Server {
	s := &Server{
		cfg: cfg,
	}

	s.server = grpc.NewServer()

	pb.RegisterSrcctlServer(s.server, NewHandler(srv))

	log.Info("gRPC server has been initialized")

	return s
}

// Run starts gRPC Server
func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.cfg.CtlAddr)
	if err != nil {
		log.WithFields(log.Fields{
			"port": s.cfg.CtlAddr,
			"err":  err,
		}).Error("Cannot bind a port", err)
	}

	log.Info("Starting gRPC server ", s.cfg.CtlAddr)
	return s.server.Serve(lis)
}

// Close stops gRPC Server
func (s *Server) Close() {
	log.Info("Stopping gRPC server ", s.cfg.CtlAddr)
	s.server.Stop()
}
