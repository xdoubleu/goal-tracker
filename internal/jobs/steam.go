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
)

type SteamJob struct {
	authService  *services.AuthService
	steamService *services.SteamService
	goalService  *services.GoalService
}

func NewSteamJob(
	authService *services.AuthService,
	steamService *services.SteamService,
	goalService *services.GoalService,
) SteamJob {
	return SteamJob{
		authService:  authService,
		steamService: steamService,
		goalService:  goalService,
	}
}

func (j SteamJob) ID() string {
	return strconv.Itoa(int(models.SteamSource.ID))
}

func (j SteamJob) RunEvery() time.Duration {
	//nolint:mnd //no magic number
	return 24 * time.Hour
}

func (j SteamJob) Run(ctx context.Context, logger *slog.Logger) error {
	users, err := j.authService.GetAllUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		logger.Debug("fetching owned games")
		var ownedGames []models.Game
		ownedGames, err = j.steamService.ImportOwnedGames(ctx, user.ID)
		if err != nil {
			return err
		}
		logger.Debug(
			fmt.Sprintf("fetched %d games", len(ownedGames)),
		)

		achievementsPerGame := map[int][]models.Achievement{}
		totalAchievementsPerGame := map[int]int{}

		var achievementsForGame map[int][]models.Achievement
		achievementsForGame, err = j.steamService.GetAchievementsForGames(
			ctx,
			ownedGames,
			user.ID,
		)
		if err != nil {
			return err
		}

		for _, game := range ownedGames {
			achievementsPerGame[game.ID] = achievementsForGame[game.ID]
			totalAchievementsPerGame[game.ID] = len(achievementsForGame[game.ID])
		}

		grapher := helper.NewAchievementsGrapher(totalAchievementsPerGame)

		totalAchievedAchievements := 0
		for gameID, achievements := range achievementsPerGame {
			for _, achievement := range achievements {
				if !achievement.Achieved {
					continue
				}

				totalAchievedAchievements++

				grapher.AddPoint(*achievement.UnlockTime, gameID)
			}
		}

		logger.Debug(
			fmt.Sprintf("achieved %d achievements in total", totalAchievedAchievements),
		)

		progressLabels, progressValues := grapher.ToSlices()

		logger.Debug("saving progress")
		err = j.goalService.SaveProgress(
			ctx,
			models.SteamCompletionRate.ID,
			user.ID,
			progressLabels,
			progressValues,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
