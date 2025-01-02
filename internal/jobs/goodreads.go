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
	return strconv.Itoa(int(models.FinishedBooksThisYear.ID))
}

func (j GoodreadsJob) RunEvery() *time.Duration {
	period := 24 * time.Hour
	return &period
}

func (j GoodreadsJob) Run(logger slog.Logger) error {
	ctx := context.Background()

	logger.Debug("fetching books")
	books, err := j.goodreadsService.GetAllBooks()
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("fetched %d books", len(books)))

	graphers := map[int]*helper.Grapher[int]{}

	graphers[time.Now().Year()] = helper.NewGrapher[int](helper.Cumulative)
	graphers[time.Now().Year()].AddPoint(time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.UTC), 0)
	graphers[time.Now().Year()].AddPoint(time.Date(time.Now().Year(), 12, 31, 0, 0, 0, 0, time.UTC), 0)

	for i, book := range books {
		logger.Debug(fmt.Sprintf("processing book %d", i))

		if len(book.DatesRead) == 0 {
			continue
		}

		for _, dateRead := range book.DatesRead {
			grapher, ok := graphers[dateRead.Year()]
			if !ok {
				graphers[dateRead.Year()] = helper.NewGrapher[int](helper.Cumulative)
				graphers[dateRead.Year()].AddPoint(time.Date(dateRead.Year(), 1, 1, 0, 0, 0, 0, time.UTC), 0)
				graphers[dateRead.Year()].AddPoint(time.Date(dateRead.Year(), 12, 31, 0, 0, 0, 0, time.UTC), 0)
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
