package repositories

import (
	"github.com/XDoubleU/essentia/pkg/database/postgres"
)

type Repositories struct {
	Goals GoalRepository
}

func New(db postgres.DB) Repositories {
	goals := GoalRepository{db: db}

	return Repositories{
		Goals: goals,
	}
}
