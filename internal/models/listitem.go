package models

type ListItem struct {
	ID        int64  `json:"id"`
	GoalID    string `json:"goalId"`
	Value     string `json:"value"`
	Completed bool   `json:"completed"`
}
