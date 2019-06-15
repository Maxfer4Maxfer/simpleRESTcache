package service

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// TopN returns N most visited request from cache
func (s *Service) TopN(n int) ([]Cache, error) {
	log.WithFields(log.Fields{
		"n": n,
	}).Info("Top N records are requested")
	r, err := s.storage.TopN(n)
	if err != nil {
		return []Cache{}, err
	}
	return r, nil
}

// LastN returns N most unvisited request from cache
func (s *Service) LastN(n int) ([]Cache, error) {
	log.WithFields(log.Fields{
		"n": n,
	}).Info("Last N records are requested")
	r, err := s.storage.LastN(n)
	if err != nil {
		return []Cache{}, err
	}
	return r, nil
}

// All returns all cache records
func (s *Service) All() ([]Cache, error) {
	log.Info("All records are requested")
	r, err := s.storage.All()
	if err != nil {
		return []Cache{}, err
	}
	return r, nil
}

// Settings returns all cache settings
func (s *Service) Settings() []string {
	log.Info("Settings are requested")
	r := []string{}
	r = append(r, fmt.Sprintf("%v<->%v", "APIAddr", s.cfg.APIAddr))
	r = append(r, fmt.Sprintf("%v<->%v", "ExpiredPeriod", s.cfg.ExpiredPeriod))
	r = append(r, fmt.Sprintf("%v<->%v", "SLA", s.cfg.SLA))
	r = append(r, fmt.Sprintf("%v<->%v", "HTTPAddr", s.cfg.HTTPAddr))
	r = append(r, fmt.Sprintf("%v<->%v", "CtlAddr", s.cfg.CtlAddr))
	r = append(r, fmt.Sprintf("%v<->%v", "DSN", s.cfg.DSN))
	r = append(r, fmt.Sprintf("%v<->%v", "Debug", s.cfg.Debug))
	return r
}

// Clean deletes all cache records
func (s *Service) Clean() error {
	log.Info("Deleting all records in cache")
	err := s.storage.Clean()
	if err != nil {
		return err
	}
	return nil
}
