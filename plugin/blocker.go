package plugin

import (
	"database/sql"
	"log"

	"github.com/miekg/dns"
)

type Blocker struct {
	db    *sql.DB
	hosts map[string]bool
}

func NewBlocker(db *sql.DB) (*Blocker, error) {
	hosts := make(map[string]bool)

	rows, err := db.Query("SELECT hostname FROM blocked_hosts")
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var host string
		if err := rows.Scan(&host); err != nil {
			log.Fatalln(err)
		}
		hosts[host+"."] = true
	}

	log.Printf("loaded %d hosts in blocklist\n", len(hosts))

	return &Blocker{db: db, hosts: hosts}, nil
}

func (b *Blocker) HandlerFunc(next HandlerFunc) HandlerFunc {
	return func(r *dns.Msg) (*dns.Msg, error) {
		for _, q := range r.Question {
			if _, ok := b.hosts[q.Name]; ok {
				log.Printf("host %s is blocked\n", q.Name)
				m := dns.Msg{}
				m.SetReply(r)
				return &m, nil
			}
		}

		return next(r)
	}
}

func (b *Blocker) Close() error {
	return b.db.Close()
}
