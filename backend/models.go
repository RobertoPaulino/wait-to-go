package main

import (
	"database/sql"
	"time"
)

type App struct {
	db      *sql.DB
	queue   *[]Entry
	history *[]Entry
}

type Entry struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Status      string    `json:"status"`
	JoinTime    time.Time `json:"joinTime"`
}

const (
	StatusWaiting  = "waiting"
	StatusNotified = "notified"
	StatusServed   = "served"
)
