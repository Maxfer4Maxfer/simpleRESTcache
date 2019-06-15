package service

import (
	"fmt"
	"strings"

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

	// prepare a masked DSN
	t1 := strings.Split(s.cfg.DSN, ":")
	t2 := strings.Split(t1[1], "@")
	t2[0] = "**********"
	t1[1] = strings.Join(t2, "@")
	maskDSN := strings.Join(t1, ":")

	r := []string{}
	r = append(r, fmt.Sprintf("%v<->%v", "APIAddr", s.cfg.APIAddr))
	r = append(r, fmt.Sprintf("%v<->%v", "ExpiredPeriod", s.cfg.ExpiredPeriod))
	r = append(r, fmt.Sprintf("%v<->%v", "SLA", s.cfg.SLA))
	r = append(r, fmt.Sprintf("%v<->%v", "HTTPAddr", s.cfg.HTTPAddr))
	r = append(r, fmt.Sprintf("%v<->%v", "CtlAddr", s.cfg.CtlAddr))
	r = append(r, fmt.Sprintf("%v<->%v", "DSN", maskDSN))
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
