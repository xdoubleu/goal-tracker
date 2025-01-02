package models

import (
	"time"
)

type Goal struct {
	ID          string     `json:"id"`
	ParentID    *string    `json:"parentId"`
	Name        string     `json:"name"`
	IsLinked    bool       `json:"isLinked"`
	TargetValue *int64     `json:"targetValue"`
	TypeID      *int64     `json:"typeId"`
	StateID     string     `json:"stateId"`
	DueTime     *time.Time `json:"time"`
	Order       int        `json:"order"`
} //	@name	Goal
