package repositories

import (
	"github.com/XDoubleU/essentia/pkg/database/postgres"
)

type Repositories struct {
	Goals    GoalRepository
	States   StateRepository
	Progress ProgressRepository
	Games    GamesRepository
}

func New(db postgres.DB) Repositories {
	goals := GoalRepository{db: db}
	states := StateRepository{db: db}
	progress := ProgressRepository{db: db}
	games := GamesRepository{db: db}

	return Repositories{
		Goals:    goals,
		States:   states,
		Progress: progress,
		Games:    games,
	}
}
