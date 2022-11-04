package plugin

import (
	"io"
	"log"

	"github.com/miekg/dns"
)

type Logger struct{}

func NewLogger(destination io.Writer) (*Logger, error) {
	return &Logger{}, nil
}

func (l *Logger) HandlerFunc(next HandlerFunc) HandlerFunc {
	return func(r *dns.Msg) (*dns.Msg, error) {
		for _, a := range r.Answer {
			log.Println(a.String())
		}

		return next(r)
	}
}

func (l *Logger) Close() error {
	return nil
}
