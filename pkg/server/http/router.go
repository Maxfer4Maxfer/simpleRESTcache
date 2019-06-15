package server

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"simpleRestCache/pkg/config"
	"simpleRestCache/pkg/service"

	log "github.com/sirupsen/logrus"
)

// NewHandler return a router for handling http request
func NewHandler(cfg *config.Config, service *service.Service) http.Handler {
	// create router
	m := http.NewServeMux()

	// setup a router path based on APIAddr geted from config
	handelPath := "/" + strings.Join(strings.Split(cfg.APIAddr, "/")[3:], "/")

	m.HandleFunc(handelPath, func(w http.ResponseWriter, req *http.Request) {
		handlePlaces(w, req, service)
	})

	return m
}

func handlePlaces(w http.ResponseWriter, req *http.Request, srv *service.Service) {
	switch req.Method {
	case "GET":
		// get a query part of a URL
		i := strings.Index(req.URL.String(), "?")
		rq := ""
		if i >= 0 {
			rq = req.URL.String()[i:]
		}

		id := uuid.New().String()

		log.WithFields(log.Fields{
			"id": id,
			"rq": rq,
		}).Info("New request is received")

		rs, sc, err := srv.HandelRequest(service.Request{ID: id, Q: rq})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			w.Write([]byte("\n"))
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(sc)
		w.Write(rs)
	default:
		log.Info("Not GET request received")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 bad request"))
	}

}
