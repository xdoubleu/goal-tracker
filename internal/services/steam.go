package services

import (
	"context"
	"log/slog"

	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/steam"
)

type SteamService struct {
	logger *slog.Logger
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

	err = service.steam.UpsertGames(
		ctx,
		gamesMap,
		userID,
	)
	if err != nil {
		return nil, err
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

	err = service.steam.UpsertAchievements(
		ctx,
		achievementsForGame.PlayerStats.Achievements,
		userID,
		game.ID,
	)
	if err != nil {
		return nil, err
	}

	if len(achievementsForGame.PlayerStats.Achievements) == 0 {
		var achievementSchemasForGame *steam.GetSchemaForGameResponse
		achievementSchemasForGame, err = service.client.GetSchemaForGame(ctx, game.ID)
		if err != nil {
			return nil, err
		}

		err = service.steam.UpsertAchievementSchemas(
			ctx,
			achievementSchemasForGame.Game.AvailableGameStats.Achievements,
			userID,
			game.ID,
		)
		if err != nil {
			return nil, err
		}
	}

	return service.steam.GetAchievementsForGame(ctx, game.ID, userID)
}
