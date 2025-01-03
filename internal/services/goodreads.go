package services

import (
	"context"
	"fmt"
	"log/slog"

	"goal-tracker/api/internal/repositories"
	"goal-tracker/api/pkg/goodreads"
)

type GoodreadsService struct {
	logger     slog.Logger
	goodreads  repositories.GoodreadsRepository
	client     goodreads.Client
	profileURL string
}

func (service GoodreadsService) ImportAllBooks(
	ctx context.Context,
) ([]goodreads.Book, error) {
	userID, err := service.client.GetUserID(service.profileURL)
	if err != nil {
		return nil, err
	}

	books, err := service.client.GetBooks(*userID)
	if err != nil {
		return nil, err
	}

	service.logger.Debug(fmt.Sprintf("saving %d books", len(books)))
	for i, book := range books {
		service.logger.Debug(fmt.Sprintf("saving book %d", i))
		err = service.goodreads.UpsertBook(ctx, book)
		if err != nil {
			return nil, err
		}
	}

	return books, nil
}

func (service GoodreadsService) GetAllBooks(
	ctx context.Context,
) ([]goodreads.Book, error) {
	return service.goodreads.GetAllBooks(ctx)
}

func (service GoodreadsService) GetAllTags(ctx context.Context) ([]string, error) {
	return service.goodreads.GetAllTags(ctx)
}

func (service GoodreadsService) GetBooksByTag(
	ctx context.Context,
	tag string,
) ([]goodreads.Book, error) {
	return service.goodreads.GetBooksByTag(ctx, tag)
}

func (service GoodreadsService) GetBooksByIDs(
	ctx context.Context,
	ids []int64,
) ([]goodreads.Book, error) {
	return service.goodreads.GetBooksByIDs(ctx, ids)
}
