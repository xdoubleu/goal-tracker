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

type GoodreadsListBooksJob struct {
	goodreadsService services.GoodreadsService
	goalService      services.GoalService
}

func NewGoodreadsListBooksJob(
	goodreadsService services.GoodreadsService,
	goalService services.GoalService,
) GoodreadsListBooksJob {
	return GoodreadsListBooksJob{
		goodreadsService: goodreadsService,
		goalService:      goalService,
	}
}

func (j GoodreadsListBooksJob) ID() string {
	return strconv.Itoa(int(models.SpecificBooks.ID))
}

func (j GoodreadsListBooksJob) RunEvery() *time.Duration {
	//nolint:mnd //no magic number
	period := 24 * time.Hour
	return &period
}

func (j GoodreadsListBooksJob) Run(logger slog.Logger) error {
	ctx := context.Background()

	logger.Debug("checking goals which track specific books")
	goals, err := j.goalService.GetGoalsByTypeID(ctx, models.SpecificBooks.ID)
	if err != nil {
		return err
	}

	for _, goal := range goals {
		var listItems []models.ListItem
		listItems, err = j.goalService.GetListItemsByGoalID(ctx, goal.ID)
		if err != nil {
			return err
		}

		bookIDs := []int64{}
		for _, listItem := range listItems {
			bookIDs = append(bookIDs, listItem.ID)
		}

		var books []goodreads.Book
		books, err = j.goodreadsService.GetBooksByIDs(ctx, bookIDs)
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
