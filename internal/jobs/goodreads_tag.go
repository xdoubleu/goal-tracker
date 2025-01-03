package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/services"
	"goal-tracker/api/pkg/goodreads"
)

type GoodreadsTagJob struct {
	goodreadsService services.GoodreadsService
	goalService      services.GoalService
}

func NewGoodreadsTagJob(
	goodreadsService services.GoodreadsService,
	goalService services.GoalService,
) GoodreadsTagJob {
	return GoodreadsTagJob{
		goodreadsService: goodreadsService,
		goalService:      goalService,
	}
}

func (j GoodreadsTagJob) ID() string {
	return strconv.Itoa(int(models.BooksFromSpecificTag.ID))
}

func (j GoodreadsTagJob) RunEvery() *time.Duration {
	//nolint:mnd //no magic number
	period := 24 * time.Hour
	return &period
}

func (j GoodreadsTagJob) Run(logger slog.Logger) error {
	ctx := context.Background()

	logger.Debug("checking goals which track specific tags")
	goals, err := j.goalService.GetGoalsByTypeID(ctx, models.BooksFromSpecificTag.ID)
	if err != nil {
		return err
	}

	for _, goal := range goals {
		var books []goodreads.Book
		books, err = j.goodreadsService.GetBooksByTag(ctx, (*goal.Config)["tag"])
		if err != nil {
			return err
		}

		for _, book := range books {
			//nolint:godox //I know
			//TODO: a book can also be reread,
			// have to check if a book was read in the past period
			_, err = j.goalService.SaveListItem(
				ctx,
				book.ID,
				goal.ID,
				fmt.Sprintf("%s - %s", book.Title, book.Author),
				len(book.DatesRead) > 0,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
