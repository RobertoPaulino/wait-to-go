package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {

	entryQueue := []Entry{}
	historySlice := []Entry{}

	connStr := "postgres://postgres:sicreto@localhost:5432/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	//initializes db
	createEntryTable(db)

	//Loads all waiting entries to the slice
	getWaitingEntry(db)

	app := App{
		db:      db,
		queue:   &entryQueue,
		history: &historySlice,
	}

	app.apiStart()

}
