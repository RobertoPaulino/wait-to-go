package main

import (
	"sort"
	"time"
)

func getTests() ([]Entry, []Entry) {
	queue := []Entry{
		{
			FirstName:   "Alice",
			LastName:    "Johnson",
			Email:       "alice@example.com",
			PhoneNumber: "555-0100",
			Status:      StatusWaiting,
			JoinTime:    time.Date(2025, 4, 20, 9, 0, 0, 0, time.UTC),
		},
		{
			FirstName:   "Bob",
			LastName:    "Smith",
			Email:       "bob@example.com",
			PhoneNumber: "555-0101",
			Status:      StatusWaiting,
			JoinTime:    time.Date(2025, 4, 20, 9, 1, 0, 0, time.UTC),
		},
		{
			FirstName:   "Charlie",
			LastName:    "Lee",
			Email:       "charlie@example.com",
			PhoneNumber: "555-0102",
			Status:      StatusNotified,
			JoinTime:    time.Date(2025, 4, 20, 8, 58, 0, 0, time.UTC),
		},
		{FirstName: "Dana",
			LastName:    "Khan",
			Email:       "dana@example.com",
			PhoneNumber: "555-0103",
			Status:      StatusServed,
			JoinTime:    time.Date(2025, 4, 20, 8, 50, 0, 0, time.UTC),
		},
		{
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
