package mocks

import (
	"context"
	"goal-tracker/api/pkg/steam"
)

type MockSteamClient struct {
}

func NewMockSteamClient() steam.Client {
	return MockSteamClient{}
}

func (client MockSteamClient) GetOwnedGames(ctx context.Context, userID string) (*steam.OwnedGamesResponse, error) {
	return nil, nil
}

func (client MockSteamClient) GetPlayerAchievements(ctx context.Context, steamID string, appID int) (*steam.AchievementsResponse, error) {
	return nil, nil
}
