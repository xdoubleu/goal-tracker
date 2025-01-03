package models

import (
	"time"
)

type Goal struct {
	ID          string             `json:"id"`
	ParentID    *string            `json:"parentId"`
	Name        string             `json:"name"`
	IsLinked    bool               `json:"isLinked"`
	TargetValue *int64             `json:"targetValue"`
	SourceID    *int64             `json:"sourceId"`
	TypeID      *int64             `json:"typeId"`
	StateID     string             `json:"stateId"`
	DueTime     *time.Time         `json:"time"`
	Order       int                `json:"order"`
	Config      *map[string]string `json:"config"`
} //	@name	Goal
