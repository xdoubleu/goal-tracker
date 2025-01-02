package repositories

import (
	"context"
	"time"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"goal-tracker/api/internal/models"
)

type ProgressRepository struct {
	db postgres.DB
}

func (repo ProgressRepository) Fetch(
	ctx context.Context,
	typeID int64,
	dateStart time.Time,
	dateEnd time.Time,
) ([]models.Progress, error) {
	query := `
		SELECT value, date 
		FROM progress 
		WHERE type_id = $1 AND date >= $2 AND date <= $3
		ORDER BY date ASC
	`

	rows, err := repo.db.Query(ctx, query, typeID, dateStart, dateEnd)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	progresses := []models.Progress{}

	for rows.Next() {
		//nolint:exhaustruct //other fields are assigned later
		progress := models.Progress{
			TypeID: typeID,
		}

		err = rows.Scan(
			&progress.Value,
			&progress.Date,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		progresses = append(progresses, progress)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return progresses, nil
}

func (repo ProgressRepository) Save(
	ctx context.Context,
	typeID int64,
	dateStr string,
	value string,
) (*models.Progress, error) {
	query := `
		INSERT INTO progress (type_id, value, date)
		VALUES ($1, $2, $3)
		ON CONFLICT (type_id, date)
		DO UPDATE SET type_id = $1, value = $2, date = $3
		RETURNING date
	`

	date, _ := time.Parse(models.ProgressDateFormat, dateStr)

	//nolint:exhaustruct //other fields are optional
	progress := models.Progress{
		TypeID: typeID,
		Value:  value,
	}

	err := repo.db.QueryRow(
		ctx,
		query,
		typeID,
		value,
		date,
	).Scan(&progress.Date)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &progress, nil
}
