package plugin

import (
	"net"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestProxy(t *testing.T) {
	t.Parallel()

	db, err := newBlocklistDB()
	require.NoError(t, err)
	defer db.Close()

	proxy, err := NewProxy("tcp-tls://one.one.one.one:853")
	require.NoError(t, err)

	hf := func(r *dns.Msg) (*dns.Msg, error) {
		return r, nil
	}
	handler := proxy.HandlerFunc(hf)

	tests := []struct {
		name string
		IP   net.IP
	}{
		{"example.com.", net.IPv4(93, 184, 216, 34)},
		{"example.org.", net.IPv4(93, 184, 216, 34)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := dns.Msg{
				MsgHdr: dns.MsgHdr{
					Id:               1,
					RecursionDesired: true,
				},
				Question: []dns.Question{{Name: test.name, Qtype: dns.TypeA, Qclass: dns.ClassINET}},
			}
			out, err := handler(&in)
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(out.Answer), 1)
			require.Equal(t, test.IP.String(), out.Answer[0].(*dns.A).A.String())
		})
	}
}
