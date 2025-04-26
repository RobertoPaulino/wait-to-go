package main

import (
	"database/sql"
	"log"
	"time"

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

	createEntryTable(db)

	//test loads test data to DB
	testEntry, _ := getTests()
	for _, entry := range testEntry {
		insertEntry(db, entry)
	}

	//Loads all waiting entries to the slice
	getWaitingEntry(db, &entryQueue)

	//test simulate usage of queue and db
	addEntry(Entry{
		FirstName:   "Fiona",
		LastName:    "McAllister",
		Email:       "fiona@example.com",
		PhoneNumber: "555-0105",
		Status:      StatusWaiting,
		JoinTime:    time.Date(2025, 4, 20, 9, 3, 0, 0, time.UTC),
	}, &entryQueue, db)

	addEntry(Entry{
		FirstName:   "George",
		LastName:    "Nguyen",
		Email:       "george@example.com",
		PhoneNumber: "555-0106",
		Status:      StatusWaiting,
		JoinTime:    time.Date(2025, 4, 20, 9, 4, 0, 0, time.UTC),
	}, &entryQueue, db)

	addEntry(Entry{
		FirstName:   "Hannah",
		LastName:    "Patel",
		Email:       "hannah@example.com",
		PhoneNumber: "555-0107",
		Status:      StatusWaiting,
		JoinTime:    time.Date(2025, 4, 20, 9, 5, 0, 0, time.UTC),
	}, &entryQueue, db)

	for range entryQueue {
		notifyNext(&entryQueue, &historySlice)
	}

	for i, _ := range historySlice {
		markServed(&historySlice[i], db)
	}

}
