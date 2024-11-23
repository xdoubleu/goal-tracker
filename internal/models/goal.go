package models

import (
	"goal-tracker/api/pkg/todoist"
)

type Goal struct {
	ID          string  `json:"id"`
	ParentID    *string `json:"parentId"`
	UserID      string  `json:"userId"`
	Name        string  `json:"name"`
	IsLinked    bool    `json:"isLinked"`
	TargetValue *int64  `json:"targetValue"`
	TypeID      *int64  `json:"typeId"`
	State       string  `json:"state"`
} //	@name	Goal

func NewGoalFromTask(task todoist.Task, userID string, state string) Goal {
	return Goal{
		ID:          task.ID,
		ParentID:    task.ParentID,
		UserID:      userID,
		Name:        task.Content,
		IsLinked:    false,
		TargetValue: nil,
		TypeID:      nil,
		State:       state,
	}
}
