package models

type State struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
} //	@name	State

var States = []State{
	BacklogState,
	PlannedState,
	InProgressState,
	DoneState,
}

var BacklogState = State{
	ID:   0,
	Name: "Backlog",
}

var PlannedState = State{
	ID:   1,
	Name: "Planned",
}

var InProgressState = State{
	ID:   2,
	Name: "In Progress",
}

var DoneState = State{
	ID:   3,
	Name: "Done",
}
