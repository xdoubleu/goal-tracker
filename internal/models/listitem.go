package models

import "time"

type ListItem struct {
	ID            int64     `json:"id"`
	ImageURL      string    `json:"imageURL"`
	Value         string    `json:"value"`
	CompletedDate time.Time `json:"completedDate"`
}
