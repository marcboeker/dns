package server

import (
	"crypto/tls"
	"errors"
	"log"

	"github.com/marcboeker/dns/plugin"
	"github.com/miekg/dns"
)

const (
	baseAddr = ":53"
	tlsAddr  = ":853"

	netUDP    = "udp"
	netTCP    = "tcp"
	netTCPTLS = "tcp4-tls"
)

type handler struct {
	plugins *plugin.Stack
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg, err := h.plugins.Handler(r)
	if err != nil {
		log.Println(err)
		m := dns.Msg{}
		m.SetRcode(r, dns.RcodeServerFailure)
		return
	}
	w.WriteMsg(msg)
}

func newHandler(plugins *plugin.Stack) dns.Handler {
	return &handler{plugins: plugins}
}

func UDPServer(plugins *plugin.Stack) error {
	h := newHandler(plugins)
	srv := dns.Server{Addr: baseAddr, Net: netUDP, Handler: h}
	return srv.ListenAndServe()
}

func TCPServer(plugins *plugin.Stack) error {
	h := newHandler(plugins)
	srv := dns.Server{Addr: baseAddr, Net: netTCP, Handler: h}
	return srv.ListenAndServe()
}

func DOTServer(hostname string, tlsConfig *tls.Config, plugins *plugin.Stack) error {
	h := newHandler(plugins)
	srv := dns.Server{Addr: tlsAddr, Net: netTCPTLS, Handler: h, TLSConfig: tlsConfig}

	srv.TLSConfig.VerifyConnection = func(cs tls.ConnectionState) error {
		if cs.ServerName != hostname {
			return errNotFound
		}
		return nil
	}

	return srv.ListenAndServe()
}

var errNotFound = errors.New("not found")
