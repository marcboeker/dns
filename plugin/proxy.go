package plugin

import (
	"log"
	"net/url"

	"github.com/miekg/dns"
)

type Proxy struct {
	host   string
	client *dns.Client
}

func NewProxy(upstream string) (*Proxy, error) {
	url, err := url.Parse(upstream)
	if err != nil {
		log.Fatalf("could not parse upstream URL: %s", err)
	}
	return &Proxy{host: url.Host, client: &dns.Client{Net: url.Scheme}}, nil
}

func (p *Proxy) HandlerFunc(next HandlerFunc) HandlerFunc {
	return func(r *dns.Msg) (*dns.Msg, error) {
		req := dns.Msg{
			Question: r.Question,
			MsgHdr:   dns.MsgHdr{RecursionDesired: r.RecursionDesired},
		}

		in, _, err := p.client.Exchange(&req, p.host)
		if err != nil {
			log.Printf("could not resolve query using upstream: %s\n", err)
			res := dns.Msg{}
			res.SetRcode(r, dns.RcodeServerFailure)
			return &res, err
		}

		msg := dns.Msg{}
		msg.SetRcode(r, in.Rcode)

		msg.RecursionAvailable = in.RecursionAvailable
		msg.Authoritative = in.Authoritative
		msg.Truncated = in.Truncated
		msg.Zero = in.Zero
		msg.CheckingDisabled = in.CheckingDisabled
		msg.Ns = in.Ns
		msg.Extra = in.Extra
		msg.Answer = in.Answer

		return next(&msg)
	}
}

func (p *Proxy) Close() error {
	return nil
}
