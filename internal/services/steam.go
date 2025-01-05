package services

import (
	"context"
	"log/slog"
	"time"

	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/steam"
)

type SteamService struct {
	logger slog.Logger
	client steam.Client
	userID string
	steam  *repositories.SteamRepository
}

func (service *SteamService) ImportOwnedGames(
	ctx context.Context,
	userID string,
) ([]models.Game, error) {
	ownedGamesResponse, err := service.client.GetOwnedGames(ctx, service.userID)
	if err != nil {
		return nil, err
	}

	gamesMap := map[int]steam.Game{}
	for _, game := range ownedGamesResponse.Response.Games {
		gamesMap[game.AppID] = game
	}

	games, err := service.steam.GetAllGames(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, game := range games {
		_, ok := gamesMap[game.ID]

		if ok {
			continue
		}

		err = service.steam.MarkGameAsDelisted(ctx, &game, userID)
		if err != nil {
			return nil, err
		}
	}

	for _, game := range gamesMap {
		err = service.steam.UpsertGame(
			ctx,
			game.AppID,
			userID,
			game.Name,
		)
		if err != nil {
			return nil, err
		}
	}

	return service.steam.GetAllGames(ctx, userID)
}

func (service *SteamService) ImportAchievementsForGame(
	ctx context.Context,
	game models.Game,
	userID string,
) ([]models.Achievement, error) {
	achievementsForGame, err := service.client.GetPlayerAchievements(
		ctx,
		service.userID,
		game.ID,
	)
	if err != nil {
		return nil, err
	}

	for _, achievement := range achievementsForGame.PlayerStats.Achievements {
		var unlockTime *time.Time
		if achievement.Achieved == 1 {
			value := time.Unix(achievement.UnlockTime, 0)
			unlockTime = &value
		}

		err = service.steam.UpsertAchievement(
			ctx,
			achievement.APIName,
			userID,
			game.ID,
			achievement.Achieved == 1,
			unlockTime,
		)
		if err != nil {
			return nil, err
		}
	}

	if len(achievementsForGame.PlayerStats.Achievements) == 0 {
		var achievementSchemasForGame *steam.GetSchemaForGameResponse
		achievementSchemasForGame, err = service.client.GetSchemaForGame(ctx, game.ID)
		if err != nil {
			return nil, err
		}

		//nolint:lll //it is what it is
		for _, achievement := range achievementSchemasForGame.Game.AvailableGameStats.Achievements {
			err = service.steam.UpsertAchievementSchema(
				ctx,
				achievement.Name,
				userID,
				game.ID,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return service.steam.GetAchievementsForGame(ctx, game.ID, userID)
}
