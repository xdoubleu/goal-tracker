package services

import (
	"context"
	"log/slog"

	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/steam"
)

type SteamService struct {
	logger slog.Logger
	client steam.Client
	games  repositories.GamesRepository
	userID string
}

func (service SteamService) GetOwnedGames(ctx context.Context) ([]int, error) {
	ownedGamesResponse, err := service.client.GetOwnedGames(ctx, service.userID)
	if err != nil {
		return nil, err
	}

	ownedGames := []int{}

	for _, game := range ownedGamesResponse.Response.Games {
		ownedGames = append(ownedGames, game.AppID)
	}

	return ownedGames, nil
}

func (service SteamService) GetAchievementsForGame(ctx context.Context, appID int) ([]steam.Achievement, error) {
	achievementsForGame, err := service.client.GetPlayerAchievements(
		ctx,
		service.userID,
		appID,
	)
	if err != nil {
		return nil, err
	}

	return achievementsForGame.PlayerStats.Achievements, nil
}
