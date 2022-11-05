package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	file := os.Args[1]
	if _, err := os.Stat(file); err != nil {
		log.Fatalln(err)
	}

	if len(os.Args[2]) == 0 {
		log.Fatalln("no database file specified")
	}
	db, err := sql.Open("sqlite3", os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS blocked_hosts (
			hostname string PRIMARY KEY
		)
	`); err != nil {
		log.Fatalln(err)
	}

	fh, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}

	stmt, err := db.Prepare("INSERT INTO blocked_hosts VALUES(?) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Fatalln(err)
	}

	count := 0
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "0.0.0.0") {
			line = strings.ReplaceAll(line, "0.0.0.0", "")
			line = strings.TrimSpace(line)

			if _, err := stmt.Exec(line); err != nil {
				log.Fatalln(err)
			}
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	log.Printf("inserted %d hosts into database", count)
}
