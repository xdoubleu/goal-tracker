package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"goal-tracker/api/internal/helper"
	"goal-tracker/api/internal/models"
	"goal-tracker/api/internal/services"
	"goal-tracker/api/pkg/goodreads"
)

type GoodreadsJob struct {
	goodreadsService services.GoodreadsService
	goalService      services.GoalService
}

func NewGoodreadsJob(
	goodreadsService services.GoodreadsService,
	goalService services.GoalService,
) GoodreadsJob {
	return GoodreadsJob{
		goodreadsService: goodreadsService,
		goalService:      goalService,
	}
}

func (j GoodreadsJob) ID() string {
	return strconv.Itoa(int(models.GoodreadsSource.ID))
}

func (j GoodreadsJob) RunEvery() *time.Duration {
	//nolint:mnd //no magic number
	period := 24 * time.Hour
	return &period
}

func (j GoodreadsJob) Run(logger slog.Logger) error {
	ctx := context.Background()

	err := j.updateProgress(ctx, logger)
	if err != nil {
		return err
	}

	logger.Debug("checking goals which track specific tags")
	err = j.specificTags(ctx)
	if err != nil {
		return err
	}

	logger.Debug("checking goals which track specific books")
	return j.specificBooks(ctx)
}

func (j GoodreadsJob) updateProgress(ctx context.Context, logger slog.Logger) error {
	logger.Debug("fetching books")
	books, err := j.goodreadsService.ImportAllBooks(ctx)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("fetched %d books", len(books)))

	graphers := map[int]*helper.Grapher[int]{}

	graphers[time.Now().Year()] = helper.NewGrapher[int](helper.Cumulative)
	graphers[time.Now().Year()].AddPoint(
		time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC),
		0,
	)
	graphers[time.Now().Year()].AddPoint(
		time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.UTC),
		0,
	)

	for i, book := range books {
		logger.Debug(fmt.Sprintf("processing book %d", i))

		if len(book.DatesRead) == 0 {
			continue
		}

		for _, dateRead := range book.DatesRead {
			grapher, ok := graphers[dateRead.Year()]
			if !ok {
				graphers[dateRead.Year()] = helper.NewGrapher[int](helper.Cumulative)
				graphers[dateRead.Year()].AddPoint(
					time.Date(dateRead.Year(), 1, 1, 0, 0, 0, 0, time.UTC),
					0,
				)
				graphers[dateRead.Year()].AddPoint(
					time.Date(dateRead.Year(), 12, 31, 0, 0, 0, 0, time.UTC),
					0,
				)
				grapher = graphers[dateRead.Year()]
			}

			grapher.AddPoint(dateRead, 1)
		}
	}

	progressLabels, progressValues := []string{}, []string{}
	for _, grapher := range graphers {
		pL, pV := grapher.ToStringSlices()
		progressLabels = append(progressLabels, pL...)
		progressValues = append(progressValues, pV...)
	}

	logger.Debug("saving progress")
	return j.goalService.SaveProgress(
		ctx,
		models.FinishedBooksThisYear.ID,
		progressLabels,
		progressValues,
	)
}

func (j GoodreadsJob) specificTags(ctx context.Context) error {
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

func (j GoodreadsJob) specificBooks(ctx context.Context) error {
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
