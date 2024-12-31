package jobs

import (
	"context"
	"strconv"
	"time"

	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/services"
)

type SteamJob struct {
	steamService services.SteamService
	goalService  services.GoalService
}

func NewSteamJob(
	steamService services.SteamService,
	goalService services.GoalService,
) SteamJob {
	return SteamJob{
		steamService: steamService,
		goalService:  goalService,
	}
}

func (j SteamJob) ID() string {
	return strconv.Itoa(int(models.SteamCompletionRate.ID))
}

func (j SteamJob) RunEvery() *time.Duration {
	period := time.Hour
	return &period
}

func (j SteamJob) Run() error {
	ctx := context.Background()

	progressLabels, progressValues, err := j.steamService.GetSteamCompletionRateProgress(
		ctx,
	)
	if err != nil {
		return err
	}

	return j.goalService.SaveProgress(
		ctx,
		models.SteamCompletionRate.ID,
		progressLabels,
		progressValues,
	)
}
