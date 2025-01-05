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
	Period      Period             `json:"period"`
	DueTime     *time.Time         `json:"time"`
	Order       int                `json:"order"`
	Config      *map[string]string `json:"config"`
}

type Period = int

const (
	Year    Period = iota
	Quarter Period = iota
	Month   Period = iota
)

func TodoistDueStringToPeriod(dueString *string) *Period {
	if dueString == nil {
		return nil
	}

	period := Year

	return &period
}
