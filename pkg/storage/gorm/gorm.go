package inmem

import (
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	"simpleRestCache/pkg/config"
	service "simpleRestCache/pkg/service"
)

// Cache represents cache in a database
type Cache struct {
	Request     string `gorm:"primary_key"`
	Responce    string `gorm:"type:text"`
	ResStatus   int
	RefreshDate time.Time
	RequestDate time.Time
	AskCount    int
	Err         error `gorm:"-"`
}

// Storage stores objects in memory
type Storage struct {
	db     *gorm.DB
	dsn    string
	cancel context.CancelFunc
}

// New returns a storage object
func New(cfg *config.Config) *Storage {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Storage{
		db:     nil,
		dsn:    cfg.DSN,
		cancel: cancel,
	}

	go s.connectToDB(ctx)

	log.Info("Storage subsystem has been initialized")
	return s
}

func (s *Storage) connectToDB(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			if s.db == nil {
				db, err := gorm.Open("mysql", s.dsn)
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Error("Error while connected to a database")
				} else {
					log.Info("Connection to a database is established")
					s.db = db
					s.db.AutoMigrate(&Cache{})
				}

			}
			if s.db != nil {
				err := s.db.DB().Ping()
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Error("Lost connection to a database")
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			s.db.Close()
		}
	}
}

// Close is an empty function for closing inmem storage
func (s *Storage) Close() {
	s.cancel()
	log.Info("Storage subsystem has been closed")
}

// Cache returns cache for a requested string
func (s *Storage) Cache(r string) service.Cache {
	c := Cache{}
	if s.db != nil {
		var count int
		s.db.Where("Request = ?", r).First(&c).Count(&count)
		if count == 0 {
			return service.Cache{Err: service.ErrCacheNotFound}
		}
		// convert datatypes from different packages
		// gorm.Cache -> service.Cache
		res := service.Cache{
			Request:     c.Request,
			Responce:    c.Responce,
			ResStatus:   c.ResStatus,
			RefreshDate: c.RefreshDate,
			RequestDate: c.RequestDate,
			AskCount:    c.AskCount,
		}
		return res
	}
	return service.Cache{Err: service.ErrStorageUnavailable}
}

// SaveCache saves record to cache
func (s *Storage) SaveCache(c service.Cache) {
	// convert datatypes from different packages
	// service.Cache -> gorm.Cache
	lc := Cache{
		Request:     c.Request,
		Responce:    c.Responce,
		ResStatus:   c.ResStatus,
		RefreshDate: time.Now(),
	}
	tmp := Cache{}
	if s.db != nil {
		var count int
		s.db.Where("Request = ?", c.Request).First(&tmp).Count(&count)
		if count != 0 {
			lc.RequestDate = tmp.RequestDate
			lc.AskCount = tmp.AskCount
		}
		s.db.Save(&lc)
	}
}

// UpdateStat updates statistic of a partitional cache record
func (s *Storage) UpdateStat(req service.Request) {
	if s.db != nil {
		c := Cache{}
		var count int
		s.db.Where("Request = ?", req.Q).First(&c).Count(&count)
		if count != 1 {
			log.WithFields(log.Fields{
				"req": req.Q,
			}).Error("Two cache records for the same request")
			return
		}
		ac := c.AskCount + 1
		c.AskCount = ac
		c.RequestDate = time.Now()

		s.db.Save(&c)

		log.WithFields(log.Fields{
			"id":    req.ID,
			"query": req.Q,
			"count": ac,
		}).Info("Statistics has been updated ")
	}
}

// All returns all requests from cache
func (s *Storage) All() ([]service.Cache, error) {
	cache := []Cache{}
	result := []service.Cache{}
	if s.db != nil {
		s.db.Find(&cache)
		for _, c := range cache {
			// convert datatypes from different packages
			// gorm.Cache -> service.Cache
			result = append(result, service.Cache{
				Request:     c.Request,
				Responce:    c.Responce,
				ResStatus:   c.ResStatus,
				RefreshDate: c.RefreshDate,
				RequestDate: c.RequestDate,
				AskCount:    c.AskCount,
			})
		}
		return result, nil
	}
	return result, service.ErrStorageUnavailable
}

// Clean deletes all cache records
func (s *Storage) Clean() error {
	if s.db != nil {
		s.db.Where("responce LIKE ?", "%").Delete(Cache{})
		return nil
	}
	return service.ErrStorageUnavailable
}

// TopN returns N most visited requests from cache
func (s *Storage) TopN(n int) ([]service.Cache, error) {
	cache := []Cache{}
	result := []service.Cache{}
	if s.db != nil {
		s.db.Order("ask_count desc, request_date desc").Limit(n).Find(&cache)
		for _, c := range cache {
			// convert datatypes from different packages
			// gorm.Cache -> service.Cache
			result = append(result, service.Cache{
				Request:     c.Request,
				Responce:    c.Responce,
				ResStatus:   c.ResStatus,
				RefreshDate: c.RefreshDate,
				RequestDate: c.RequestDate,
				AskCount:    c.AskCount,
			})
		}
		return result, nil
	}
	return result, service.ErrStorageUnavailable
}

// LastN returns N most unvisited requests from cache
func (s *Storage) LastN(n int) ([]service.Cache, error) {
	cache := []Cache{}
	result := []service.Cache{}
	if s.db != nil {
		s.db.Order("ask_count, request_date").Limit(n).Find(&cache)
		for _, c := range cache {
			// convert datatypes from different packages
			// gorm.Cache -> service.Cache
			result = append(result, service.Cache{
				Request:     c.Request,
				Responce:    c.Responce,
				ResStatus:   c.ResStatus,
				RefreshDate: c.RefreshDate,
				RequestDate: c.RequestDate,
				AskCount:    c.AskCount,
			})
		}
		return result, nil
	}
	return result, service.ErrStorageUnavailable
}
