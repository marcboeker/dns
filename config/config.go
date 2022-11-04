package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/marcboeker/dns/plugin"
)

type Config struct {
	Plugins     []plugin.Plugin
	DNSOverHTTP *DNSOverHTTP
	DNSOverTLS  *DNSOverTLS
	ListenUDP   bool
	ListenTCP   bool
	TLSConfig   *tls.Config
}

type DNSOverHTTP struct {
	Hostname string
	Path     string
}

type DNSOverTLS struct {
	Hostname string
}

// Use `dev` as the current config environment.
var Get = dev

// newTLSConfig loads certs and returns a TLS config.
func newTLSConfig(caCert, cert, key []byte) (*tls.Config, error) {
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("could not read CA cert")
	}
	c, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("could not read cert key pair: %s", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{c},
		ClientCAs:    caCertPool,
	}, nil
}
