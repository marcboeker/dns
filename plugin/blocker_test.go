package plugin

import (
	"database/sql"
	"net"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestBlocker(t *testing.T) {
	t.Parallel()

	db, err := newBlocklistDB()
	require.NoError(t, err)
	defer db.Close()

	blocker, err := NewBlocker(db)
	require.NoError(t, err)

	hf := func(r *dns.Msg) (*dns.Msg, error) {
		if r.Question[0].Name != "example.com." {
			r.Answer = append(r.Answer, &dns.A{A: net.IPv4(127, 0, 0, 1)})
		}
		return r, nil
	}
	handler := blocker.HandlerFunc(hf)

	tests := []struct {
		name      string
		isBlocked bool
	}{
		{"example.com.", true},
		{"example.org.", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in := dns.Msg{Question: []dns.Question{{Name: test.name, Qtype: dns.TypeA}}}
			out, err := handler(&in)
			require.NoError(t, err)
			require.Equal(t, test.isBlocked, len(out.Answer) == 0)
		})
	}
}

func newBlocklistDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS blocked_hosts (
			hostname string PRIMARY KEY
		)
	`); err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO blocked_hosts VALUES(?)", "example.com")
	if err != nil {
		return nil, err
	}

	return db, nil
}
