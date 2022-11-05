package config

import (
	"database/sql"
	"os"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"

	"github.com/marcboeker/dns/plugin"
)

// Generate dev certs via `mkcert -install && mkcert localhost`
// https://github.com/FiloSottile/mkcert

//go:embed dev/ca.pem
var devCACert []byte

//go:embed dev/cert.pem
var devCert []byte

//go:embed dev/key.pem
var devKey []byte

// dev prepares a config for local development.
func dev() (*Config, error) {
	db, err := sql.Open("sqlite3", "dns.db")
	if err != nil {
		return nil, err
	}

	// Possible URL schemes are tcp://, udp:// or tcp-tls:// for DoT.
	proxy, err := plugin.NewProxy("tcp-tls://one.one.one.one:853")
	if err != nil {
		return nil, err
	}

	stats, err := plugin.NewStats(plugin.StatsOpts{
		DB:           db,
		TrackStats:   true, // Count queries per hostname
		TrackQueries: true, // Log all queries with a timestamp
	})
	if err != nil {
		return nil, err
	}

	// Update blocklist via `make update-blocklist`
	blocker, err := plugin.NewBlocker(db)
	if err != nil {
		return nil, err
	}

	logger, err := plugin.NewLogger(os.Stdout)
	if err != nil {
		return nil, err
	}

	tlsConfig, err := newTLSConfig(devCACert, devCert, devKey)
	if err != nil {
		return nil, err
	}

	return &Config{
		Plugins: []plugin.Plugin{
			blocker, stats, proxy, logger,
		},
		DNSOverHTTP: &DNSOverHTTP{
			Hostname: "localhost", // Listen on localhost 443
			Path:     "/dns-query",
		},
		DNSOverTLS: &DNSOverTLS{
			Hostname: "localhost", // Listen on localhost 853
		},
		ListenUDP: true, // Listen on UDP port 53
		ListenTCP: true, // Listen on TCP port 53
		// Specify certs for DNS over HTTPS and DNS over TLS.
		TLSConfig: tlsConfig,
	}, nil
}
