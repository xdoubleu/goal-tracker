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
	authService      *services.AuthService
	goodreadsService *services.GoodreadsService
	goalService      *services.GoalService
}

func NewGoodreadsJob(
	authService *services.AuthService,
	goodreadsService *services.GoodreadsService,
	goalService *services.GoalService,
) GoodreadsJob {
	return GoodreadsJob{
		authService:      authService,
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

	users, err := j.authService.GetAllUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		err = j.updateProgress(ctx, logger, user.ID)
		if err != nil {
			return err
		}

		logger.Debug("checking goals which track specific tags")
		err = j.specificTags(ctx, user.ID)
		if err != nil {
			return err
		}

		logger.Debug("checking goals which track specific books")
		err = j.specificBooks(ctx, user.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (j GoodreadsJob) updateProgress(
	ctx context.Context,
	logger slog.Logger,
	userID string,
) error {
	logger.Debug("fetching books")
	books, err := j.goodreadsService.ImportAllBooks(ctx, userID)
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

	for _, book := range books {
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
		userID,
		progressLabels,
		progressValues,
	)
}

func (j GoodreadsJob) specificTags(ctx context.Context, userID string) error {
	goals, err := j.goalService.GetGoalsByTypeID(
		ctx,
		models.BooksFromSpecificTag.ID,
		userID,
	)
	if err != nil {
		return err
	}

	for _, goal := range goals {
		var books []goodreads.Book
		books, err = j.goodreadsService.GetBooksByTag(
			ctx,
			(*goal.Config)["tag"],
			userID,
		)
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
				userID,
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

func (j GoodreadsJob) specificBooks(ctx context.Context, userID string) error {
	goals, err := j.goalService.GetGoalsByTypeID(ctx, models.SpecificBooks.ID, userID)
	if err != nil {
		return err
	}

	for _, goal := range goals {
		var listItems []models.ListItem
		listItems, err = j.goalService.GetListItemsByGoalID(ctx, goal.ID, userID)
		if err != nil {
			return err
		}

		bookIDs := []int64{}
		for _, listItem := range listItems {
			bookIDs = append(bookIDs, listItem.ID)
		}

		var books []goodreads.Book
		books, err = j.goodreadsService.GetBooksByIDs(ctx, bookIDs, userID)
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
				userID,
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
