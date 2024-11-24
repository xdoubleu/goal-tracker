package services

import (
	"context"
	"fmt"
	"goal-tracker/api/pkg/steam"
	"math"
)

type SteamService struct {
	client steam.Client
	userID string
}

func (service SteamService) GetSteamCompletionRate(ctx context.Context) (int, error) {
	fmt.Println("fetching owned games")
	ownedGamesResponse, err := service.client.GetOwnedGames(ctx, service.userID)
	if err != nil {
		return 0, err
	}
	fmt.Printf("fetched %d games\n", ownedGamesResponse.Response.GameCount)

	achievements := map[int][]steam.Achievement{}
	for i, game := range ownedGamesResponse.Response.Games {
		fmt.Printf("fetching achievements for game %d\n", (i + 1))
		achievementsForGame, err := service.client.GetPlayerAchievements(ctx, service.userID, game.AppID)
		if err != nil {
			return 0, err
		}
		achievements[game.AppID] = achievementsForGame.PlayerStats.Achievements
	}

	return calculateSteamCompletionRate(achievements), nil
}

func calculateSteamCompletionRate(achievements map[int][]steam.Achievement) int {
	totalCompletionRate := 0.0
	startedGames := 0

	for _, achievementsForGame := range achievements {
		totalAchievementsGame := len(achievementsForGame)
		achievedAchievementsGame := 0

		for _, achievement := range achievementsForGame {
			if achievement.Achieved == 1 {
				achievedAchievementsGame++
			}
		}

		if achievedAchievementsGame > 0 {
			startedGames++
			totalCompletionRate += float64(achievedAchievementsGame) / float64(totalAchievementsGame)
		}
	}

	return int(math.Floor(totalCompletionRate * 100.0 / float64(startedGames)))
}
