package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"wait-to-go/auth"
)

// CORS middleware
func enableCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

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

	//validate we have a name and phone number
	if (entry.FirstName == "") || (len(entry.FirstName) > 30) || (entry.LastName == "") || (len(entry.LastName) > 30) || (entry.PhoneNumber == "") {
		http.Error(w, "Invalid name fields or missing phone number", http.StatusBadRequest)
		return
	}

	id, err := addEntry(entry, a.queue, a.db)
	if err != nil {
		http.Error(w, "Failed to add entry", http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(id, entry.PhoneNumber)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"id":     id,
		"token":  token,
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

	// Get claims from context (set by auth middleware)
	claims, ok := auth.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify that the token matches the requested entry
	if claims.ID != entryID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	// Calculate position in queue
	position := 0
	if entry.Status == StatusWaiting {
		for _, e := range *a.queue {
			if e.Status == StatusWaiting && e.JoinTime.Before(entry.JoinTime) {
				position++
			}
		}
		position++ // Add 1 because we want 1-based position
	}

	response := struct {
		Entry    Entry `json:"entry"`
		Position int   `json:"position"`
	}{
		Entry:    entry,
		Position: position,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
