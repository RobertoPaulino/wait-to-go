package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (a *App) apiStart() {
	mux := http.NewServeMux()

	//Entrypoints
	mux.HandleFunc("/join", a.handleJoin)
	mux.HandleFunc("/queue", a.handleQueue)
	mux.HandleFunc("/next", a.handleNext)
	mux.HandleFunc("/serve", a.handleServe)
	mux.HandleFunc("/status/", a.handleStatus)

	fmt.Println("Serverlistening to 8080")
	http.ListenAndServe(":8080", mux)

}

func (a *App) handleJoin(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	entry.Status = StatusWaiting
	entry.JoinTime = time.Now()

	//validate we have a name
	if (entry.FirstName == "") || (len(entry.FirstName) > 30) || (entry.LastName == "") || (len(entry.LastName) > 30) {
		http.Error(w, "Invalid name fields", http.StatusBadRequest)
		return
	}

	addEntry(entry, a.queue, a.db)

	w.WriteHeader(http.StatusCreated)
}

func (a *App) handleQueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	entries := getWaitingEntry(a.db)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(entries)

}

func (a *App) handleNext(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	notifyNext(a.queue, a.history)
}

func (a *App) handleServe(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	markServed(&entry, a.db)
}

// TODO: we need to review previous code so when we query the database
// we also get the ID and attach it to the entry, it will be easier to manipulate the db and even auth down the line
func (a *App) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var entry Entry
	err := json.NewDecoder(r.Body).Decode(&entry)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}
