package main

import "time"

// Password represents a stored password entry with metadata
type Password struct {
	Name         string    `json:"name"`
	Value        string    `json:"value"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`
}

// NewPassword creates a new Password instance with current timestamps
func NewPassword(name, value, category string) Password {
	return Password{Value: value, Name: name, Category: category, CreatedAt: time.Now(), LastModified: time.Now()}
}
