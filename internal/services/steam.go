package services

import (
	"context"
	"log/slog"
	"sync"

	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/steam"

	"github.com/XDoubleU/essentia/pkg/threading"
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

func (service *SteamService) ImportAchievementsForGames(
	ctx context.Context,
	games []models.Game,
	userID string,
) (map[int][]models.Achievement, error) {
	var err error

	amountWorkers := (len(games) / 10) + 1
	workerPool := threading.NewWorkerPool(service.logger, amountWorkers, len(games))

	mu := sync.Mutex{}
	achievementsPerGame := map[int][]steam.Achievement{}
	for _, game := range games {
		workerPool.EnqueueWork(func(ctx context.Context, logger *slog.Logger) {
			achievementsForGame, errIn := service.client.GetPlayerAchievements(
				ctx,
				service.userID,
				game.ID,
			)
			if errIn != nil {
				err = errIn
			}

			mu.Lock()
			achievementsPerGame[game.ID] = achievementsForGame.PlayerStats.Achievements
			mu.Unlock()
		})
	}

	workerPool.WaitUntilDone()
	workerPool.Stop()
	if err != nil {
		return nil, err
	}

	gameIDs := []int{}
	for gameID, achievements := range achievementsPerGame {
		gameIDs = append(gameIDs, gameID)

		err = service.steam.UpsertAchievements(
			ctx,
			achievements,
			userID,
			gameID,
		)
		if err != nil {
			return nil, err
		}

		if len(achievements) != 0 {
			continue
		}

		var achievementSchemasForGame *steam.GetSchemaForGameResponse
		achievementSchemasForGame, err = service.client.GetSchemaForGame(ctx, gameID)
		if err != nil {
			return nil, err
		}

		err = service.steam.UpsertAchievementSchemas(
			ctx,
			achievementSchemasForGame.Game.AvailableGameStats.Achievements,
			userID,
			gameID,
		)
		if err != nil {
			return nil, err
		}
	}

	return service.steam.GetAchievementsForGames(ctx, gameIDs, userID)
}
