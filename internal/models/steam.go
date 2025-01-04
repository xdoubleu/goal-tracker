package models

import "time"

type Game struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IsDelisted bool   `json:"isDelisted"`
}

type Achievement struct {
	Name       string     `json:"name"`
	GameID     int        `json:"gameId"`
	Achieved   bool       `json:"achieved"`
	UnlockTime *time.Time `json:"unlockTime"`
}
