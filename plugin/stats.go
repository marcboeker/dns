package plugin

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/miekg/dns"
)

type Stats struct {
	db   *sql.DB
	opts StatsOpts
}
type StatsOpts struct {
	DB           *sql.DB
	TrackStats   bool
	TrackQueries bool
}

func NewStats(opts StatsOpts) (*Stats, error) {
	s := Stats{opts: opts, db: opts.DB}

	if err := s.initDB(); err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *Stats) HandlerFunc(next HandlerFunc) HandlerFunc {
	return func(r *dns.Msg) (*dns.Msg, error) {
		for _, q := range r.Question {
			if s.opts.TrackStats {
				if _, err := s.db.Exec(`
				INSERT INTO stats (hostname, count)
				VALUES (?, 1)
				ON CONFLICT(hostname)
				DO UPDATE SET count = count + 1`,
					q.Name,
				); err != nil {
					log.Fatal(err)
				}
			}

			if s.opts.TrackQueries {
				if _, err := s.db.Exec(`
				INSERT INTO queries (hostname, timestamp)
				VALUES (?, ?)`,
					q.Name, time.Now().UTC(),
				); err != nil {
					log.Fatal(err)
				}
			}
		}

		return next(r)
	}
}

func (s *Stats) initDB() error {
	if s.opts.TrackStats {
		if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS stats (
			hostname TEXT PRIMARY KEY,
			count INTEGER
		)`); err != nil {
			return err
		}
	}

	if s.opts.TrackQueries {
		if _, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS queries (
			hostname TEXT,
			timestamp TEXT
		);
		CREATE INDEX IF NOT EXISTS queries_domain ON queries (hostname);`); err != nil {
			return err
		}
	}

	return nil
}

func (s *Stats) Close() error {
	return s.db.Close()
}
