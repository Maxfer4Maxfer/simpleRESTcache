package inmem

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"simpleRestCache/pkg/config"
	service "simpleRestCache/pkg/service"
)

// Storage stores objects in memory
type Storage struct {
	cache map[string]service.Cache
	sync.RWMutex
}

// New returns a storage object
func New(cfg *config.Config) *Storage {
	s := &Storage{
		cache: make(map[string]service.Cache),
	}
	log.Info("Storage subsystem has been initialized")
	return s
}

// Close is an empty function for closing inmem storage
func (s *Storage) Close() {
	log.Info("Storage subsystem has been closed")
}

// Cache returns cache for a requested string
func (s *Storage) Cache(r string) service.Cache {
	s.RLock()
	defer s.RUnlock()

	_, ok := s.cache[r]
	if !ok {
		return service.Cache{Err: service.ErrCacheNotFound}
	}
	return s.cache[r]
}

// SaveCache saves record to cache
func (s *Storage) SaveCache(c service.Cache) {
	s.Lock()
	defer s.Unlock()

	if v, ok := s.cache[c.Request]; ok {
		// if a record is exist update save RequestDate and AskCount
		s.cache[c.Request] = service.Cache{
			Request:     c.Request,
			Responce:    c.Responce,
			ResStatus:   c.ResStatus,
			RefreshDate: time.Now(),
			RequestDate: v.RequestDate,
			AskCount:    v.AskCount,
		}
	} else {
		// if new then add a new record
		s.cache[c.Request] = service.Cache{
			Request:     c.Request,
			Responce:    c.Responce,
			ResStatus:   c.ResStatus,
			RefreshDate: time.Now(),
		}
	}
}

// UpdateStat updates statistic of a partitional cache record
func (s *Storage) UpdateStat(req service.Request) {
	s.Lock()
	defer s.Unlock()

	c := s.cache[req.Q]
	ac := c.AskCount + 1
	s.cache[req.Q] = service.Cache{
		Request:     c.Request,
		Responce:    c.Responce,
		ResStatus:   c.ResStatus,
		RefreshDate: c.RefreshDate,
		RequestDate: time.Now(),
		AskCount:    ac,
	}
	log.WithFields(log.Fields{
		"id":    req.ID,
		"query": req.Q,
		"count": ac,
	}).Info("View count and request date fields has been increased ")
}
