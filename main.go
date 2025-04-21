package main

import (
	"fmt"
	"sort"
	"time"
)

type Entry struct {
	ID          int
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

func getTests() []Entry {
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
	return queue
}

func notifyNext(queue *[]Entry) {
	sort.Sort(ByJoinTime(*queue))

	for index, value := range *queue {
		if value.Status == StatusWaiting {
			fmt.Printf("Now Notifying: %v \n", value.FirstName)
			(*queue)[index].Status = StatusNotified
			break
		}
	}

	// Debug output
	for index, value := range *queue {
		if value.Status == StatusWaiting {
			fmt.Printf("Position: %v Name:%v Status:%v \n", index, value.FirstName, value.Status)
		}
	}

	fmt.Print("\n\n\n----------------------------------\n\n\n")
}

func addEntry(entry Entry, queue *[]Entry) {
	*queue = append(*queue, entry)
}

func main() {
	testEntries := getTests()

	newEntry := Entry{
		ID:          10,
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "JDoe@ABC.net",
		PhoneNumber: "111-111-1111",
		Status:      StatusWaiting,
		JoinTime:    time.Now(),
	}

	notifyNext(&testEntries)
	time.Sleep(3 * time.Second)

	addEntry(newEntry, &testEntries)
	time.Sleep(3 * time.Second)

	notifyNext(&testEntries)
	time.Sleep(3 * time.Second)

	notifyNext(&testEntries)
	time.Sleep(3 * time.Second)

	notifyNext(&testEntries)
	time.Sleep(3 * time.Second)

	notifyNext(&testEntries)
	time.Sleep(3 * time.Second)

}
