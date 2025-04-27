package main

import (
	"database/sql"
	"fmt"
	"time"
)

func createEntryTable(db *sql.DB) error {
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
		return fmt.Errorf("failed to create table: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func insertEntry(db *sql.DB, entry Entry) (int, error) {
	query := `INSERT INTO entry (firstName, lastName, email, phoneNumber, status, joinTime) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var pk int
	err := db.QueryRow(query, entry.FirstName, entry.LastName, entry.Email, entry.PhoneNumber, entry.Status, entry.JoinTime).Scan(&pk)
	if err != nil {
		return 0, fmt.Errorf("failed to insert entry: %w", err)
	}

	return pk, nil
}

func updateStatusByEntry(db *sql.DB, entry Entry) error {
	query := `UPDATE entry SET status = $1 WHERE id = $2`
	_, err := db.Exec(query, entry.Status, entry.ID)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

func getWaitingEntry(db *sql.DB) ([]Entry, error) {
	var entries []Entry

	query := `SELECT id, firstName, lastName, email, phoneNumber, status, joinTime FROM entry WHERE status = 'waiting'`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query waiting entries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var firstName, lastName, email, phoneNumber, status string
		var joinTime time.Time

		err := rows.Scan(&id, &firstName, &lastName, &email, &phoneNumber, &status, &joinTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		entries = append(entries, Entry{id, firstName, lastName, email, phoneNumber, status, joinTime})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return entries, nil
}

func backupHistory(db *sql.DB, historyList *[]Entry) error {
	for _, entry := range *historyList {
		_, err := insertEntry(db, entry)
		if err != nil {
			return fmt.Errorf("failed to backup history entry: %w", err)
		}
	}

	*historyList = []Entry{}
	return nil
}

func getEntryByID(db *sql.DB, id int) (Entry, error) {
	query := `SELECT id, firstName, lastName, email, phoneNumber, status, joinTime FROM entry WHERE id = $1`
	var entry Entry
	err := db.QueryRow(query, id).Scan(
		&entry.ID,
		&entry.FirstName,
		&entry.LastName,
		&entry.Email,
		&entry.PhoneNumber,
		&entry.Status,
		&entry.JoinTime,
	)
	if err != nil {
		return Entry{}, fmt.Errorf("failed to get entry: %w", err)
	}
	return entry, nil
}

func clearQueue(db *sql.DB) error {
	query := `UPDATE entry SET status = $1 WHERE status = $2`
	_, err := db.Exec(query, StatusServed, StatusWaiting)
	if err != nil {
		return fmt.Errorf("failed to clear queue: %w", err)
	}
	return nil
}
