package jobs

import (
	"context"
	"goal-tracker/api/internal/services"
	"time"
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
	period := 24 * time.Hour
	return &period
}

func (j TodoistJob) Run() error {
	ctx := context.Background()
	return j.goalService.ImportFromTodoist(ctx)
}
