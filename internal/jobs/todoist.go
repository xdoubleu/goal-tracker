package jobs

import (
	"context"
	"log/slog"
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

func (j TodoistJob) Run(logger slog.Logger) error {
	ctx := context.Background()

	logger.Debug("importing states")
	err := j.goalService.ImportStatesFromTodoist(ctx)
	if err != nil {
		return err
	}

	logger.Debug("importing goals")
	return j.goalService.ImportGoalsFromTodoist(ctx)
}
