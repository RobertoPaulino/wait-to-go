package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func (a *App) handleJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	entry.Status = StatusWaiting
	entry.JoinTime = time.Now()

	//validate we have a name
	if (entry.FirstName == "") || (len(entry.FirstName) > 30) || (entry.LastName == "") || (len(entry.LastName) > 30) {
		http.Error(w, "Invalid name fields", http.StatusBadRequest)
		return
	}

	id, err := addEntry(entry, a.queue, a.db)
	if err != nil {
		http.Error(w, "Failed to add entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     id,
	})
}

func (a *App) handleQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entries, err := getWaitingEntry(a.db)
	if err != nil {
		http.Error(w, "Failed to get queue", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (a *App) handleNext(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := notifyNext(a.queue, a.history, a.db)
	if err != nil {
		http.Error(w, "Failed to notify next", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (a *App) handleServe(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = markServed(&entry, a.db)
	if err != nil {
		http.Error(w, "Failed to mark as served", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (a *App) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	id := r.URL.Path[len("/status/"):]
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	entryID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	entry, err := getEntryByID(a.db, entryID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Entry not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entry)
}

func (a *App) handleClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := clearQueueInMemory(a.queue, a.db)
	if err != nil {
		http.Error(w, "Failed to clear queue", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
