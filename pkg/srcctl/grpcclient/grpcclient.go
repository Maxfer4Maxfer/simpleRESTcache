package grpcclien

import (
	"context"
	"fmt"
	"os"
	"simpleRestCache/pb"
	"strconv"
	"strings"

	timestamp "github.com/golang/protobuf/ptypes"
	"github.com/olekukonko/tablewriter"
	"google.golang.org/grpc"
)

// Handler redirects request to a service and output result to a console
type Handler struct {
	addr string
}

// New returns new handler
func New(addr string) *Handler {
	return &Handler{
		addr: addr,
	}
}

// All returns N most visited request from cache
func (h *Handler) All() {
	grcpConn, err := grpc.Dial(
		h.addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}
	defer grcpConn.Close()

	service := pb.NewSrcctlClient(grcpConn)

	ctx := context.Background()
	res, err := service.All(ctx, &pb.AllRequest{})
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"Request", "Status", "Refresh Date", "Request Date", "Count"})

	for _, c := range res.Cache {
		refDate, err := timestamp.Timestamp(c.RefreshDate)
		if err != nil {
			fmt.Println("Cannot parse a responce")
			fmt.Println("Error = ", err)
			return
		}
		reqDate, err := timestamp.Timestamp(c.RequestDate)
		if err != nil {
			fmt.Println("Cannot parse a responce")
			fmt.Println("Error = ", err)
			return
		}
		status := strconv.FormatInt(int64(c.ResStatus), 10)
		count := strconv.FormatInt(int64(c.AskCount), 10)

		table.Append([]string{c.Request, status, refDate.Format("2006-01-02 15:04:05"), reqDate.Format("2006-01-02 15:04:05"), count})
	}

	table.Render()
}

// TopN returns N most visited request from cache
func (h *Handler) TopN(n int) {
	grcpConn, err := grpc.Dial(
		h.addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}
	defer grcpConn.Close()

	service := pb.NewSrcctlClient(grcpConn)

	ctx := context.Background()
	res, err := service.TopN(ctx, &pb.TopNRequest{N: int32(n)})
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"Request", "Status", "Refresh Date", "Request Date", "Count"})

	for _, c := range res.Cache {
		refDate, err := timestamp.Timestamp(c.RefreshDate)
		if err != nil {
			fmt.Println("Cannot parse a responce")
			fmt.Println("Error = ", err)
			return
		}
		reqDate, err := timestamp.Timestamp(c.RequestDate)
		if err != nil {
			fmt.Println("Cannot parse a responce")
			fmt.Println("Error = ", err)
			return
		}
		status := strconv.FormatInt(int64(c.ResStatus), 10)
		count := strconv.FormatInt(int64(c.AskCount), 10)

		table.Append([]string{c.Request, status, refDate.Format("2006-01-02 15:04:05"), reqDate.Format("2006-01-02 15:04:05"), count})
	}

	table.Render()
}

// LastN returns N most visited request from cache
func (h *Handler) LastN(n int) {
	grcpConn, err := grpc.Dial(
		h.addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}
	defer grcpConn.Close()

	service := pb.NewSrcctlClient(grcpConn)

	ctx := context.Background()
	res, err := service.LastN(ctx, &pb.LastNRequest{N: int32(n)})
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"Request", "Status", "Refresh Date", "Request Date", "Count"})

	for _, c := range res.Cache {
		refDate, err := timestamp.Timestamp(c.RefreshDate)
		if err != nil {
			fmt.Println("Cannot parse a responce")
			fmt.Println("Error = ", err)
			return
		}
		reqDate, err := timestamp.Timestamp(c.RequestDate)
		if err != nil {
			fmt.Println("Cannot parse a responce")
			fmt.Println("Error = ", err)
			return
		}
		status := strconv.FormatInt(int64(c.ResStatus), 10)
		count := strconv.FormatInt(int64(c.AskCount), 10)

		table.Append([]string{c.Request, status, refDate.Format("2006-01-02 15:04:05"), reqDate.Format("2006-01-02 15:04:05"), count})
	}

	table.Render()
}

// Settings returns all cache settings
func (h *Handler) Settings() {
	grcpConn, err := grpc.Dial(
		h.addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}
	defer grcpConn.Close()

	service := pb.NewSrcctlClient(grcpConn)

	ctx := context.Background()
	res, err := service.Settings(ctx, &pb.SettingsRequest{})
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"Name", "Value"})
	for _, s := range res.Settings {
		tmp := strings.Split(s, "<->")
		table.Append([]string{tmp[0], tmp[1]})
	}
	table.Render()
}

// Clean deletes all cached records
func (h *Handler) Clean() {
	grcpConn, err := grpc.Dial(
		h.addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}
	defer grcpConn.Close()

	service := pb.NewSrcctlClient(grcpConn)

	ctx := context.Background()
	_, err = service.Clean(ctx, &pb.CleanRequest{})
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}

	fmt.Println("Cache was cleaned")
}

// Refresh renews all cache records
func (h *Handler) Refresh() {
	grcpConn, err := grpc.Dial(
		h.addr,
		grpc.WithInsecure(),
	)
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}
	defer grcpConn.Close()

	service := pb.NewSrcctlClient(grcpConn)

	ctx := context.Background()
	_, err = service.Refresh(ctx, &pb.RefreshRequest{})
	if err != nil {
		fmt.Println("Cannot connect to the service")
		fmt.Println("Error = ", err)
		return
	}

	fmt.Println("Cache was refreshed")
}
