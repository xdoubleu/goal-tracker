package models

import (
	"strings"
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
	Period      *Period            `json:"period"`
	DueTime     *time.Time         `json:"time"`
	Order       int                `json:"order"`
	Config      *map[string]string `json:"config"`
}

type GoalWithSubGoals struct {
	Goal
	SubGoals []GoalWithSubGoals `json:"subGoals"`
}

type Period = int

const (
	Year    Period = iota
	Quarter Period = iota
	Month   Period = iota
)

func (goal Goal) PeriodStart(includePreviousPeriod bool) time.Time {
	multiplier := 1
	if includePreviousPeriod {
		multiplier = 2
	}

	switch *goal.Period {
	case Year:
		return goal.DueTime.AddDate(-1*multiplier, 0, 1)
	case Quarter:
		return goal.DueTime.AddDate(0, -3*multiplier, 1)
	default:
		panic("not implemented")
	}
}

func (goal Goal) PeriodEnd() time.Time {
	return *goal.DueTime
}

func TodoistDueStringToPeriod(dueString string) *Period {
	if dueString == "" {
		return nil
	}

	dueStringClean := strings.Split(dueString, "every ")[1]

	var period Period
	switch dueStringClean {
	case "year":
		period = Year
	case "3 months":
		period = Quarter
	default:
		return nil
	}

	return &period
}
