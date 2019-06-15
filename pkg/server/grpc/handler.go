package grpc

import (
	"context"

	timestamp "github.com/golang/protobuf/ptypes"
	// log "github.com/sirupsen/logrus"

	"simpleRestCache/pb"
	"simpleRestCache/pkg/service"
)

// Handler handles gRPC request and pass them to a service
type Handler struct {
	service *service.Service
}

// NewHandler creates a new gRPC handler
func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		service: srv,
	}
}

// TopN returns a responce from TopN function of the service
func (h *Handler) TopN(ctx context.Context, req *pb.TopNRequest) (*pb.TopNReply, error) {
	n := int(req.GetN())
	cache, err := h.service.TopN(n)
	if err != nil {
		return &pb.TopNReply{}, err
	}

	// convert datatypes from different packages
	// storage.Cache -> pb.Cache
	pbc := []*pb.Cache{}
	for _, c := range cache {
		refDate, err := timestamp.TimestampProto(c.RefreshDate)
		if err != nil {
			return &pb.TopNReply{}, err
		}

		reqDate, err := timestamp.TimestampProto(c.RequestDate)
		if err != nil {
			return &pb.TopNReply{}, err
		}

		pbc = append(pbc, &pb.Cache{
			Request:     c.Request,
			Responce:    c.Responce,
			ResStatus:   int32(c.ResStatus),
			RefreshDate: refDate,
			RequestDate: reqDate,
			AskCount:    int32(c.AskCount),
		})
	}

	return &pb.TopNReply{Cache: pbc}, nil
}

// LastN returns a responce from LastN function of the service
func (h *Handler) LastN(ctx context.Context, req *pb.LastNRequest) (*pb.LastNReply, error) {
	n := int(req.GetN())
	cache, err := h.service.LastN(n)
	if err != nil {
		return &pb.LastNReply{}, err
	}

	// convert datatypes from different packages
	// storage.Cache -> pb.Cache
	pbc := []*pb.Cache{}
	for _, c := range cache {
		refDate, err := timestamp.TimestampProto(c.RefreshDate)
		if err != nil {
			return &pb.LastNReply{}, err
		}

		reqDate, err := timestamp.TimestampProto(c.RequestDate)
		if err != nil {
			return &pb.LastNReply{}, err
		}

		pbc = append(pbc, &pb.Cache{
			Request:     c.Request,
			Responce:    c.Responce,
			ResStatus:   int32(c.ResStatus),
			RefreshDate: refDate,
			RequestDate: reqDate,
			AskCount:    int32(c.AskCount),
		})
	}

	return &pb.LastNReply{Cache: pbc}, nil
}

// All returns a responce from All function of the service
func (h *Handler) All(ctx context.Context, req *pb.AllRequest) (*pb.AllReply, error) {
	cache, err := h.service.All()
	if err != nil {
		return &pb.AllReply{}, err
	}

	// convert datatypes from different packages
	// storage.Cache -> pb.Cache
	pbc := []*pb.Cache{}
	for _, c := range cache {
		refDate, err := timestamp.TimestampProto(c.RefreshDate)
		if err != nil {
			return &pb.AllReply{}, err
		}

		reqDate, err := timestamp.TimestampProto(c.RequestDate)
		if err != nil {
			return &pb.AllReply{}, err
		}

		pbc = append(pbc, &pb.Cache{
			Request:     c.Request,
			Responce:    c.Responce,
			ResStatus:   int32(c.ResStatus),
			RefreshDate: refDate,
			RequestDate: reqDate,
			AskCount:    int32(c.AskCount),
		})
	}

	return &pb.AllReply{Cache: pbc}, nil
}

// Settings returns a responce from Settings function of the service
func (h *Handler) Settings(ctx context.Context, req *pb.SettingsRequest) (*pb.SettingsReply, error) {
	ss := h.service.Settings()

	return &pb.SettingsReply{Settings: ss}, nil
}

// Clean returns a responce from Clean function of the service
func (h *Handler) Clean(ctx context.Context, req *pb.CleanRequest) (*pb.CleanReply, error) {
	err := h.service.Clean()
	if err != nil {
		return &pb.CleanReply{}, err
	}
	return &pb.CleanReply{}, nil
}

// Refresh renews all cache records
func (h *Handler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshReply, error) {
	err := h.service.Refresh()
	if err != nil {
		return &pb.RefreshReply{}, err
	}
	return &pb.RefreshReply{}, nil
}
