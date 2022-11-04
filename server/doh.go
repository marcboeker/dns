package server

import (
	"crypto/tls"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/marcboeker/dns/plugin"
	"github.com/miekg/dns"
)

func DOHServer(plugin *plugin.Stack, hostname, path string, tlsConfig *tls.Config) error {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		// Make sure, that the server only replies to request on the configured hostname.
		if r.Host != hostname {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		var (
			buf []byte
			err error
		)
		switch r.Method {
		case http.MethodGet:
			buf, err = base64.RawURLEncoding.DecodeString(r.URL.Query().Get("dns"))
			if len(buf) == 0 || err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
		case http.MethodPost:
			if r.Header.Get("Content-Type") != "application/dns-message" {
				http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
				return
			}
			buf, err = io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		req := new(dns.Msg)
		if err := req.Unpack(buf); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		msg, err := plugin.Handler(req)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		packed, err := msg.Pack()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/dns-message")

		_, _ = w.Write(packed)
	})

	s := http.Server{Addr: ":https", TLSConfig: tlsConfig}
	return s.ListenAndServeTLS("", "")
}
