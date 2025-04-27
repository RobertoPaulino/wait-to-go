package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func createEntryTable(db *sql.DB) {

	query := `CREATE TABLE IF NOT EXISTS entry (
		id SERIAL PRIMARY KEY,
		firstName VARCHAR(30) NOT NULL,
		lastName VARCHAR(30) NOT NULL,
		email VARCHAR(50),
		phoneNumber VARCHAR(10),
		status VARCHAR(20) NOT NULL,
		joinTime timestamp DEFAULT NOW()
	)`

	_, err := db.Exec(query)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

}

func insertEntry(db *sql.DB, entry Entry) int {
	query := `INSERT INTO entry (firstName, lastName, email, phoneNumber, status, joinTime) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var pk int
	err := db.QueryRow(query, entry.FirstName, entry.LastName, entry.Email, entry.PhoneNumber, entry.Status, entry.JoinTime).Scan(&pk)

	if err != nil {
		log.Fatal(err)
	}

	return pk
}

func updateStatusByEntry(db *sql.DB, entry Entry) error {
	query := `UPDATE entry SET status = $1 WHERE firstName = $2 AND lastName = $3 AND phoneNumber = $4`
	_, err := db.Exec(query, entry.Status, entry.FirstName, entry.LastName, entry.PhoneNumber)
	return err
}

func getWaitingEntry(db *sql.DB) []Entry {

	var entries []Entry

	query := `SELECT id, firstName, lastName, email, phoneNumber, status, joinTime FROM entry WHERE status = 'waiting'`
	rows, err := db.Query(query)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Found no rows waiting")
		} else {
			log.Fatal(err)
		}
	}

	defer rows.Close()

	for rows.Next() {

		var id int
		var firstName, lastName, email, phoneNumber, status string
		var joinTime time.Time

		err := rows.Scan(&id, &firstName, &lastName, &email, &phoneNumber, &status, &joinTime)

		if err != nil {
			log.Fatal(err)
		}

		entries = append(entries, Entry{id, firstName, lastName, email, phoneNumber, status, joinTime})
	}
	return entries
}

func backupHistory(db *sql.DB, historyList *[]Entry) {
	for _, entry := range *historyList {
		insertEntry(db, entry)
	}

	*historyList = []Entry{}
}
