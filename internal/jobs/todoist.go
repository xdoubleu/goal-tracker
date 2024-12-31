package jobs

import (
	"context"
	"time"

	"goal-tracker/api/internal/services"
)

type TodoistJob struct {
	goalService services.GoalService
}

func NewTodoistJob(goalService services.GoalService) TodoistJob {
	return TodoistJob{
		goalService: goalService,
	}
}

func (j TodoistJob) ID() string {
	return "todoist"
}

func (j TodoistJob) RunEvery() *time.Duration {
	//nolint:mnd //no magic number
	period := 24 * time.Hour
	return &period
}

func (j TodoistJob) Run() error {
	ctx := context.Background()
	return j.goalService.ImportFromTodoist(ctx)
}
