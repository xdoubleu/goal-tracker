package models

import (
	"goal-tracker/api/pkg/todoist"
)

type Goal struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	IsLinked    bool   `json:"isLinked"`
	TargetValue *int64 `json:"targetValue"`
	SourceID    *int64 `json:"sourceId"`
	TypeID      *int64 `json:"typeId"`
	State       string `json:"state"`
} //	@name	Goal

func NewGoalFromTask(task todoist.Task, userId string, state string) Goal {
	return Goal{
		ID:          task.Id,
		UserID:      userId,
		Name:        task.Content,
		IsLinked:    false,
		TargetValue: nil,
		SourceID:    nil,
		TypeID:      nil,
		State:       state,
	}
}
