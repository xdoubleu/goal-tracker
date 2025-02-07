package jobs

import (
	"context"
	"log/slog"
	"time"

	"goal-tracker/api/internal/services"
)

type TodoistJob struct {
	authService *services.AuthService
	goalService *services.GoalService
}

func NewTodoistJob(
	authService *services.AuthService,
	goalService *services.GoalService,
) TodoistJob {
	return TodoistJob{
		authService: authService,
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

func (j TodoistJob) Run(logger *slog.Logger) error {
	ctx := context.Background()

	users, err := j.authService.GetAllUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		logger.Debug("importing states")
		err = j.goalService.ImportStatesFromTodoist(ctx, user.ID)
		if err != nil {
			return err
		}

		logger.Debug("importing goals")
		err = j.goalService.ImportGoalsFromTodoist(ctx, user.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
