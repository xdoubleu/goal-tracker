package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"goal-tracker/api/internal/helper"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/services"
	"goal-tracker/api/pkg/steam"
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
	return strconv.Itoa(int(models.SteamSource.ID))
}

func (j SteamJob) RunEvery() *time.Duration {
	//nolint:mnd //no magic number
	period := 24 * time.Hour
	return &period
}

func (j SteamJob) Run(logger slog.Logger) error {
	ctx := context.Background()

	logger.Debug("fetching owned games")
	ownedGames, err := j.steamService.GetOwnedGames(ctx)
	if err != nil {
		return err
	}
	logger.Debug(
		fmt.Sprintf("fetched %d games", len(ownedGames)),
	)

	totalAchievementsPerGame := map[int]int{}
	achievementsPerGame := map[int][]steam.Achievement{}
	for i, game := range ownedGames {
		logger.Debug(
			fmt.Sprintf("fetching achievements for game %d", (i + 1)),
		)

		var achievementsForGame []steam.Achievement
		achievementsForGame, err = j.steamService.GetAchievementsForGame(ctx, game)
		if err != nil {
			return err
		}

		achievementsPerGame[game] = achievementsForGame
		totalAchievementsPerGame[game] = len(achievementsPerGame[game])
	}

	grapher := helper.NewAchievementsGrapher(totalAchievementsPerGame)

	totalAchievedAchievements := 0
	for gameID, achievements := range achievementsPerGame {
		for _, achievement := range achievements {
			if achievement.Achieved == 0 {
				continue
			}

			totalAchievedAchievements++

			time := time.Unix(achievement.UnlockTime, 0)
			grapher.AddPoint(time, gameID)
		}
	}

	logger.Debug(
		fmt.Sprintf("achieved %d achievements in total", totalAchievedAchievements),
	)

	progressLabels, progressValues := grapher.ToSlices()

	logger.Debug("saving progress")
	return j.goalService.SaveProgress(
		ctx,
		models.SteamCompletionRate.ID,
		progressLabels,
		progressValues,
	)
}
