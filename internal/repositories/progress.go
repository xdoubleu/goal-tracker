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

func (repo ProgressRepository) GetByTypeIDAndDates(
	ctx context.Context,
	typeID int64,
	userID string,
	dateStart time.Time,
	dateEnd time.Time,
) ([]models.Progress, error) {
	query := `
		SELECT value, date 
		FROM progress 
		WHERE type_id = $1 AND user_id = $2 AND date >= $3 AND date <= $4
		ORDER BY date ASC
	`

	rows, err := repo.db.Query(ctx, query, typeID, userID, dateStart, dateEnd)
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

func (repo ProgressRepository) Upsert(
	ctx context.Context,
	typeID int64,
	userID string,
	dateStr string,
	value string,
) (*models.Progress, error) {
	query := `
		INSERT INTO progress (type_id, user_id, date, value)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (type_id, user_id, date)
		DO UPDATE SET value = $4
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
		userID,
		date,
		value,
	).Scan(&progress.Date)

	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return &progress, nil
}
