package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/marcboeker/dns/config"
	"github.com/marcboeker/dns/plugin"
	"github.com/marcboeker/dns/server"
)

func main() {
	// Config can be modified in config/config.go
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("could not build config: %s\n", err)
	}

	plugins := plugin.Load(cfg.Plugins)

	if cfg.ListenUDP {
		go server.UDPServer(plugins)
	}

	if cfg.ListenTCP {
		go server.TCPServer(plugins)
	}

	if cfg.DNSOverTLS != nil && cfg.TLSConfig != nil {
		go server.DOTServer(cfg.DNSOverTLS.Hostname, cfg.TLSConfig, plugins)
	}

	if cfg.DNSOverHTTP != nil && cfg.TLSConfig != nil {
		go server.DOHServer(plugins, cfg.DNSOverHTTP.Hostname, cfg.DNSOverHTTP.Path, cfg.TLSConfig)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Wait for CTRL+C and gracefully shutdown server.
	var sig = make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		s := <-sig
		log.Printf("caught kill signal: %+v", s)
		if err := plugins.Close(); err != nil {
			log.Fatalf("could not shutdown plugins: %s", err)
		}
		wg.Done()
		os.Exit(0)
	}()

	log.Println("server has been started, shutdown with CTRL+C")

	wg.Wait()
}
