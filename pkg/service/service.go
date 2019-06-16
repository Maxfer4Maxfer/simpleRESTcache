package service

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"

	"simpleRestCache/pkg/config"
	parser "simpleRestCache/pkg/parser/aviasalesru/placesjsonv2"
)

// Storage declares methods that the real storage object should implement
type Storage interface {
	Cache(r string) Cache
	SaveCache(c Cache)
	Clean() error
	UpdateStat(req Request)
	TopN(n int) ([]Cache, error)
	LastN(n int) ([]Cache, error)
	All() ([]Cache, error)
}

// Request represents a request
type Request struct {
	ID string // inner system ID
	Q  string // query body
}

// Cache represents cache
type Cache struct {
	Request     string
	Responce    string
	ResStatus   int
	RefreshDate time.Time
	RequestDate time.Time
	AskCount    int
	Err         error
}

// APIResp is a responce from the endpoint
type APIResp struct {
	Resp   []byte
	Status int
	Err    error
}

// Service is a central component of the system. It contains all business logic.
type Service struct {
	storage Storage
	cfg     *config.Config
}

var (
	// ErrCacheNotFound arise when there is no record in cache
	ErrCacheNotFound = errors.New("Cache is not found")

	// ErrCannotParseMessage arise when a parser cannot parse a message
	ErrCannotParseMessage = errors.New("Cannot parse a message")

	// ErrEndpointAPIUnavailable arise if something happened with a endpoint
	ErrEndpointAPIUnavailable = errors.New("Endpoint API is unavailable")

	// ErrStorageUnavailable arise when something wrong with a storage subsystem
	ErrStorageUnavailable = errors.New("Storage subsystem is unavailable")
)

// New retunrs new Service
func New(cfg *config.Config, store Storage) *Service {
	s := &Service{
		storage: store,
		cfg:     cfg,
	}
	return s
}

// HandelRequest redirects request to endpoint and also stores a responce in cache
// HandelRequest returns Parsed/Unparsed, Status Code, error
func (s *Service) HandelRequest(req Request) ([]byte, int, error) {

	ctxAPI, cancelAPI := context.WithCancel(context.Background())

	// two channal for interact with two parallel request
	chRespStorage := make(chan Cache)
	chRespAPI := make(chan APIResp)

	// send a request to Storage
	go func() {
		chRespStorage <- s.storage.Cache(req.Q)
		close(chRespStorage)
	}()

	// send a request to Endpoint
	go func() {
		select {
		case <-time.After(s.cfg.SLA / 10):
			chRespAPI <- s.requestToAPI(req)
		case <-ctxAPI.Done():
		}
		close(chRespAPI)
	}()

	sla := time.NewTimer(s.cfg.SLA)
	respStorage := Cache{}
	for {
		select {
		case respStorage = <-chRespStorage: // got a responce from Storage
			if respStorage != (Cache{}) && respStorage.Err == nil {
				log.WithFields(log.Fields{
					"id": req.ID,
				}).Info("Find a responce in cache...")
				// if cache is not expired immediately return it
				if time.Since(respStorage.RefreshDate) <= s.cfg.ExpiredPeriod {
					log.WithFields(log.Fields{
						"id": req.ID,
					}).Info("...and cache is not expired")

					// cancel API request
					cancelAPI()
					// start gorutine to wait a responce from API for savint it to Storage
					go func(req Request, chRespAPI chan APIResp) {
						select {
						case respAPI, ok := <-chRespAPI:
							if ok {
								log.WithFields(log.Fields{
									"id": req.ID,
								}).Info("Did not have time to stop the request to the Endpoin. Refresh cache.")
								s.storage.SaveCache(Cache{
									Request:   req.Q,
									Responce:  string(respAPI.Resp),
									ResStatus: respAPI.Status,
								})
							}
						case <-time.After(s.cfg.SLA * 2):
						}
					}(req, chRespAPI)

					// update statistic
					go s.storage.UpdateStat(req)

					log.WithFields(log.Fields{
						"id": req.ID,
					}).Info("...returinint a cache record for a responce")
					return []byte(respStorage.Responce), respStorage.ResStatus, nil
				}
			}
		case respAPI := <-chRespAPI: // got a responce from Endpoint
			log.WithFields(log.Fields{
				"id": req.ID,
			}).Info("Recived a responce from Endpoint")

			// save a responce to cache and update statistic
			go func() {
				s.storage.SaveCache(Cache{
					Request:   req.Q,
					Responce:  string(respAPI.Resp),
					ResStatus: respAPI.Status,
				})
				s.storage.UpdateStat(req)
			}()

			return respAPI.Resp, respAPI.Status, respAPI.Err
		case <-sla.C: // Reached SLA
			log.WithFields(log.Fields{
				"id":  req.ID,
				"sla": s.cfg.SLA,
			}).Warn("Reached SLA...")
			// check responce from Storage. If it not empty return it
			// if empty then wait a responce from the endpoint
			if respStorage != (Cache{}) && respStorage.Err == nil {
				cancelAPI()
				log.WithFields(log.Fields{
					"id":  req.ID,
					"sla": s.cfg.SLA,
				}).Warn("...returning expired cache")
				return []byte(respStorage.Responce), respStorage.ResStatus, nil
			}
			log.Warn("...and don't have cache")
		}
	}
}

func (s *Service) requestToAPI(req Request) APIResp {
	log.WithFields(log.Fields{
		"id":      req.ID,
		"rq":      req.Q,
		"endpoin": s.cfg.APIAddr,
	}).Info("Start processing request to Endpoint")

	// Send a request to the endpoint
	resp, err := http.Get(s.cfg.APIAddr + req.Q)
	if err != nil {
		log.WithFields(log.Fields{
			"id":  req.ID,
			"err": err,
			"url": s.cfg.APIAddr + req.Q,
		}).Error("Error while calling endpoint")
		return APIResp{
			Resp:   []byte{},
			Status: http.StatusInternalServerError,
			Err:    ErrEndpointAPIUnavailable,
		}
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"id":  req.ID,
			"err": err,
		}).Error("Error read a response's body")
		return APIResp{
			Resp:   []byte{},
			Status: http.StatusInternalServerError,
			Err:    ErrEndpointAPIUnavailable,
		}
	}

	r := []byte{}
	if resp.StatusCode == http.StatusOK {
		// parse a request
		r, err = parser.Parse(respBody)
		if err != nil {
			log.WithFields(log.Fields{
				"id":  req.ID,
				"err": err,
			}).Error("Error while parsing a responce")
		}
	} else {
		log.WithFields(log.Fields{
			"id":          req.ID,
			"status_code": resp.StatusCode,
		}).Warn("Did not parse responce because status code of a responce is not OK")
		r = respBody
	}

	log.WithFields(log.Fields{
		"id": req.ID,
	}).Info("Successfully parse a request")

	return APIResp{
		Resp:   r,
		Status: resp.StatusCode,
		Err:    nil,
	}
}

// Refresh renews all cache records
func (s *Service) Refresh() error {
	cache, err := s.All()
	if err != nil {
		log.Error("Error while requested all cache records")
	}

	for _, c := range cache {
		r := s.requestToAPI(Request{
			ID: uuid.New().String(),
			Q:  c.Request,
		})
		s.storage.SaveCache(Cache{
			Request:   c.Request,
			Responce:  string(r.Resp),
			ResStatus: r.Status,
		})
	}
	return nil
}
