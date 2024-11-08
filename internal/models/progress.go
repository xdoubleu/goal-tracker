package models

import "time"

type Progress struct {
	ID        int64     `json:"id"`
	GoalID    string    `json:"goalId"`
	Value     int64     `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
} //	@name	Progress
