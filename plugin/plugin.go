package plugin

import (
	"github.com/miekg/dns"
)

type HandlerFunc func(*dns.Msg) (*dns.Msg, error)

type Plugin interface {
	HandlerFunc(next HandlerFunc) HandlerFunc
	Close() error
}

type Stack struct {
	list []Plugin
}

// Load registers a list of plugins for later use.
func Load(plugins []Plugin) *Stack {
	s := Stack{}
	s.list = append(s.list, plugins...)
	return &s
}

// Close closes all plugins and returns an error if a plugin has failed.
func (s *Stack) Close() error {
	for _, p := range s.list {
		if err := p.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Handler iterates over all loaded plugins.
// The DNS request is passed from one plugin to the next and returned in the end.
func (s *Stack) Handler(r *dns.Msg) (*dns.Msg, error) {
	h := func(r *dns.Msg) (*dns.Msg, error) {
		return r, nil
	}
	for i := len(s.list) - 1; i >= 0; i-- {
		h = s.list[i].HandlerFunc(h)
	}
	return h(r)
}
