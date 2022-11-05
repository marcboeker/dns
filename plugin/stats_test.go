package plugin

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/require"
)

func TestStats(t *testing.T) {
	t.Parallel()

	db, err := newStatsDB()
	require.NoError(t, err)
	defer db.Close()

	opts := StatsOpts{
		DB:           db,
		TrackStats:   true,
		TrackQueries: true,
	}
	stats, err := NewStats(opts)
	require.NoError(t, err)

	hf := func(r *dns.Msg) (*dns.Msg, error) {
		return r, nil
	}
	handler := stats.HandlerFunc(hf)

	in := dns.Msg{Question: []dns.Question{{Name: "example.com.", Qtype: dns.TypeA}}}
	_, err = handler(&in)
	require.NoError(t, err)

	row := db.QueryRow("SELECT count FROM stats WHERE hostname = 'example.com.'")
	require.NoError(t, err)
	var count int
	require.NoError(t, row.Scan(&count))
	require.Equal(t, count, 1)

	row = db.QueryRow("SELECT count(*) FROM stats")
	require.NoError(t, err)
	require.NoError(t, row.Scan(&count))
	require.Equal(t, count, 1)
}

func newStatsDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS stats (
			hostname TEXT PRIMARY KEY,
			count INTEGER
		)`); err != nil {
		return nil, err
	}

	return db, nil
}
