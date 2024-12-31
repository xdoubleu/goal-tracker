package services

import (
	"context"
	"fmt"
	"goal-tracker/api/internal/helper"
	"goal-tracker/api/pkg/steam"
	"time"
)

type SteamService struct {
	client steam.Client
	userID string
}

func (service SteamService) GetSteamCompletionRateProgress(ctx context.Context) ([]string, []int64, error) {
	achievementsPerGame, totalAchievementsPerGame, err := service.fetchAchievementsFromApi(ctx)
	if err != nil {
		return nil, nil, err
	}

	grapher := helper.NewGrapher(totalAchievementsPerGame)

	for gameID, achievements := range achievementsPerGame {
		for _, achievement := range achievements {
			if achievement.Achieved == 0 {
				continue
			}

			time := time.Unix(achievement.UnlockTime, 0)
			grapher.AddPoint(time, gameID)
		}
	}

	dates, percentages := grapher.ToSlices()
	return dates, percentages, nil
}

func (service SteamService) fetchAchievementsFromApi(ctx context.Context) (map[int][]steam.Achievement, map[int]int, error) {
	fmt.Println("fetching owned games")
	ownedGamesResponse, err := service.client.GetOwnedGames(ctx, service.userID)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("fetched %d games\n", ownedGamesResponse.Response.GameCount)

	totalAchievementsPerGame := map[int]int{}
	achievements := map[int][]steam.Achievement{}
	for i, game := range ownedGamesResponse.Response.Games {
		fmt.Printf("fetching achievements for game %d (%s)\n", (i + 1), game.Name)
		achievementsForGame, err := service.client.GetPlayerAchievements(ctx, service.userID, game.AppID)
		if err != nil {
			return nil, nil, err
		}
		achievements[game.AppID] = achievementsForGame.PlayerStats.Achievements
		totalAchievementsPerGame[game.AppID] = len(achievements[game.AppID])
	}

	return achievements, totalAchievementsPerGame, nil
}
