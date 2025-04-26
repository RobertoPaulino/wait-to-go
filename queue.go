package main

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"sort"
)

type ByJoinTime []Entry

func (q ByJoinTime) Len() int           { return len(q) }
func (q ByJoinTime) Less(i, j int) bool { return q[i].JoinTime.Before(q[j].JoinTime) }
func (q ByJoinTime) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func notifyNext(queue *[]Entry, history *[]Entry) {

	sort.Sort(ByJoinTime(*queue))

	if (*queue)[0].Status == StatusWaiting {
		(*queue)[0].Status = StatusNotified
		//TODO: Notification system.
		//notify(*queue[0].phonenumber) or something like that
		*history = append(*history, (*queue)[0])
		*queue = slices.Delete(*queue, 0, 1)
	}

}

func markServed(entry *Entry, db *sql.DB) {
	entry.Status = StatusServed
	err := updateStatusByEntry(db, *entry)
	if err != nil {
		log.Fatal(err)
	}
}

func addEntry(entry Entry, queue *[]Entry, db *sql.DB) {

	if entry.Status == StatusWaiting {
		*queue = append(*queue, entry)
		sort.Sort(ByJoinTime(*queue))
		insertEntry(db, entry)
	} else {
		fmt.Println("Tried adding entry with wrong status")
	}

}
