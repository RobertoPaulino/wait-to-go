package main

import "time"

type Entry struct {
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
