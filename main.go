package main

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"sort"
	"time"

	_ "github.com/lib/pq"
)

type Entry struct {
	ID          int
	Position    int
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
	Status      string
	JoinTime    time.Time
}

const (
	StatusWaiting  = "waiting"
	StatusNotified = "notified"
	StatusServed   = "served"
)

type ByJoinTime []Entry

func (q ByJoinTime) Len() int           { return len(q) }
func (q ByJoinTime) Less(i, j int) bool { return q[i].JoinTime.Before(q[j].JoinTime) }
func (q ByJoinTime) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func getTests() ([]Entry, []Entry) {
	queue := []Entry{
		{
			ID:          1,
			FirstName:   "Alice",
			LastName:    "Johnson",
			Email:       "alice@example.com",
			PhoneNumber: "555-0100",
			Status:      StatusWaiting,
			JoinTime:    time.Date(2025, 4, 20, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:          2,
			FirstName:   "Bob",
			LastName:    "Smith",
			Email:       "bob@example.com",
			PhoneNumber: "555-0101",
			Status:      StatusWaiting,
			JoinTime:    time.Date(2025, 4, 20, 9, 1, 0, 0, time.UTC),
		},
		{
			ID:          3,
			FirstName:   "Charlie",
			LastName:    "Lee",
			Email:       "charlie@example.com",
			PhoneNumber: "555-0102",
			Status:      StatusNotified,
			JoinTime:    time.Date(2025, 4, 20, 8, 58, 0, 0, time.UTC),
		},
		{
			ID:          4,
			FirstName:   "Dana",
			LastName:    "Khan",
			Email:       "dana@example.com",
			PhoneNumber: "555-0103",
			Status:      StatusServed,
			JoinTime:    time.Date(2025, 4, 20, 8, 50, 0, 0, time.UTC),
		},
		{
			ID:          5,
			FirstName:   "Eli",
			LastName:    "Garcia",
			Email:       "eli@example.com",
			PhoneNumber: "555-0104",
			Status:      StatusWaiting,
			JoinTime:    time.Date(2025, 4, 20, 9, 2, 0, 0, time.UTC),
		}}

	sort.Sort(ByJoinTime(queue))
	return queue, []Entry{}
}

func notifyNext(queue *[]Entry, history *[]Entry) {

	if (*queue)[0].Status == StatusWaiting {
		fmt.Printf("Now Notifying: %v \n", (*queue)[0].FirstName)
		(*queue)[0].Status = StatusNotified
		*history = append(*history, (*queue)[0])
		*queue = slices.Delete(*queue, 0, 1)

		// Debug output
		for index, value := range *queue {
			if value.Status == StatusWaiting {
				fmt.Printf("Position: %v Name:%v Status:%v \n", index, value.FirstName, value.Status)
			}
		}

		fmt.Print("\n\n\n----------------------------------\n\n\n")
	}

}

func addEntry(entry Entry, queue *[]Entry) {

	if entry.Status == StatusWaiting {
		*queue = append(*queue, entry)
		sort.Sort(ByJoinTime(*queue))
	} else {
		fmt.Println("Tried adding entry with wrong status")
	}
}

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

func main() {

	entryList, _ := getTests()
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

	pk := insertEntry(db, entryList[0])

	print(pk)
}
