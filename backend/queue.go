package main

import (
	"database/sql"
	"fmt"
	"slices"
	"sort"
)

type ByJoinTime []Entry

func (q ByJoinTime) Len() int           { return len(q) }
func (q ByJoinTime) Less(i, j int) bool { return q[i].JoinTime.Before(q[j].JoinTime) }
func (q ByJoinTime) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func notifyNext(queue *[]Entry, history *[]Entry, db *sql.DB) error {
	if len(*queue) == 0 {
		return fmt.Errorf("queue is empty")
	}

	sort.Sort(ByJoinTime(*queue))

	if (*queue)[0].Status == StatusWaiting {
		(*queue)[0].Status = StatusNotified
		// Update status in database
		if err := updateStatusByEntry(db, (*queue)[0]); err != nil {
			return fmt.Errorf("failed to update status in database: %w", err)
		}
		*history = append(*history, (*queue)[0])
		*queue = slices.Delete(*queue, 0, 1)
		return nil
	}

	return fmt.Errorf("next entry is not in waiting status")
}

func markServed(entry *Entry, db *sql.DB) error {
	if entry == nil {
		return fmt.Errorf("entry is nil")
	}

	entry.Status = StatusServed
	return updateStatusByEntry(db, *entry)
}

func addEntry(entry Entry, queue *[]Entry, db *sql.DB) (int, error) {
	if entry.Status != StatusWaiting {
		return 0, fmt.Errorf("entry must be in waiting status")
	}

	id, err := insertEntry(db, entry)
	if err != nil {
		return 0, fmt.Errorf("failed to insert entry: %w", err)
	}

	entry.ID = id
	*queue = append(*queue, entry)
	sort.Sort(ByJoinTime(*queue))

	return id, nil
}

func clearQueueInMemory(queue *[]Entry, db *sql.DB) error {
	if err := clearQueue(db); err != nil {
		return fmt.Errorf("failed to clear queue in database: %w", err)
	}
	*queue = []Entry{}
	return nil
}
