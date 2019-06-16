package config

import (
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// Config contains all configuration of App
type Config struct {
	APIAddr       string
	DSN           string
	HTTPAddr      string
	CtlAddr       string
	ExpiredPeriod time.Duration
	SLA           time.Duration
	Debug         bool
}

// GetConfig returns a fulfilled Config
func GetConfig() *Config {
	fs := flag.NewFlagSet("simpleRESTcache", flag.ExitOnError)
	var (
		apiAddr       = fs.String("api-URL", "https://places.aviasales.ru/v2/places.json", "URL of an endpoint API")
		dsn           = fs.String("dsn", "root:root@tcp(mysql:3306)/tasks?charset=utf8&parseTime=True&loc=Local", "Database Source Name")
		httpAddr      = fs.String("http-addr", ":8080", "HTTP listen address")
		ctlAddr       = fs.String("control-addr", ":8081", "Control listen address")
		expiredPeriod = fs.Duration("expiredPeriod", 24*time.Hour, "Expired cache duration. Valid time units are \"m\", \"h\"")
		sla           = fs.Duration("sla", 3*time.Second, "SLA time is a period for which a response to a client must be provided. Valid time units are \"ms\", \"s\", \"m\", \"h\"")
		debug         = fs.Bool("debug", false, "Set debug mode")
	)

	fs.Parse(os.Args[1:])

	// Check inputs
	if *expiredPeriod < 1*time.Minute {
		log.Error("Expired cache duration should be more then 1 minute")
		os.Exit(1)
	}

	if *sla < 1*time.Millisecond {
		log.Error("SLA time should be more then 1ms")
		os.Exit(1)
	}

	cfg := Config{
		APIAddr:       *apiAddr,
		DSN:           *dsn,
		HTTPAddr:      *httpAddr,
		CtlAddr:       *ctlAddr,
		ExpiredPeriod: *expiredPeriod,
		SLA:           *sla,
		Debug:         *debug,
	}

	return &cfg
}
