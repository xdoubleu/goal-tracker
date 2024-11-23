package repositories

import (
	"github.com/XDoubleU/essentia/pkg/database/postgres"
)

type Repositories struct {
	Goals    GoalRepository
	Progress ProgressRepository
}

func New(db postgres.DB) Repositories {
	goals := GoalRepository{db: db}
	progress := ProgressRepository{db: db}

	return Repositories{
		Goals:    goals,
		Progress: progress,
	}
}
