package models

import "time"

type Goal struct {
	ID          string     `json:"id"`
	UserID      string     `json:"userId"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
	Value       *int64     `json:"value"`
	SourceID    *int64     `json:"sourceId"`
	TypeID      *int64     `json:"typeId"`
	Score       int64      `json:"score"`
	StateID     int64      `json:"stateId"`
} //	@name	Goal
