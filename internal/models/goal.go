package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sgreben/piecewiselinear"
)

type Goal struct {
	ID          string            `json:"id"`
	ParentID    *string           `json:"parentId"`
	Name        string            `json:"name"`
	IsLinked    bool              `json:"isLinked"`
	TargetValue *int64            `json:"targetValue"`
	SourceID    *int64            `json:"sourceId"`
	TypeID      *int64            `json:"typeId"`
	StateID     string            `json:"stateId"`
	Period      *Period           `json:"period"`
	DueTime     *time.Time        `json:"time"`
	Order       int               `json:"order"`
	Config      map[string]string `json:"config"`
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

func (goal Goal) PeriodStart() time.Time {
	switch *goal.Period {
	case Year:
		return goal.DueTime.AddDate(-1, 0, 1)
	case Quarter:
		return goal.DueTime.AddDate(0, -3, 1)
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

func (goal Goal) AdaptiveTargetValues(startProgress int) []string {
	secondsInADay := 86400

	f := piecewiselinear.Function{
		X: []float64{
			float64(goal.PeriodStart().Unix() / int64(secondsInADay)),
			float64(goal.PeriodEnd().Unix() / int64(secondsInADay)),
		},
		Y: []float64{float64(startProgress), float64(*goal.TargetValue)},
	}

	result := []string{}
	for i := goal.PeriodStart(); i.Before(goal.PeriodEnd()); i = i.AddDate(0, 0, 1) {
		result = append(
			result,
			fmt.Sprintf("%.2f", f.At(float64(i.Unix()/int64(secondsInADay)))),
		)
	}

	return result
}
