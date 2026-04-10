package main

import "time"

type Password struct {
	Name         string    `json:"name"`
	Value        string    `json:"value"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

func NewPassword(value, name, category string) Password {
	return Password{Value: value, Name: name, Category: category, CreatedAt: time.Now(), LastModified: time.Now()}
}
