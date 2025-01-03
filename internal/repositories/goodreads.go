package repositories

import (
	"context"

	"github.com/XDoubleU/essentia/pkg/database/postgres"

	"goal-tracker/api/pkg/goodreads"
)

type GoodreadsRepository struct {
	db postgres.DB
}

func (repo GoodreadsRepository) GetAllBooks(
	ctx context.Context,
) ([]goodreads.Book, error) {
	query := `
		SELECT id, shelf, tags, title, author, dates_read
		FROM goodreads_books 
	`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	books := []goodreads.Book{}
	for rows.Next() {
		var book goodreads.Book

		err = rows.Scan(
			&book.ID,
			&book.Shelf,
			&book.Tags,
			&book.Title,
			&book.Author,
			&book.DatesRead,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return books, nil
}

func (repo GoodreadsRepository) GetAllTags(ctx context.Context) ([]string, error) {
	query := `
		SELECT ARRAY_AGG(DISTINCT tag) AS tags 
		FROM 
			goodreads_books,
			UNNEST(tags) as tag
		WHERE tags <> '{}';
	`

	tags := []string{}
	err := repo.db.QueryRow(ctx, query).Scan(&tags)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return tags, nil
}

func (repo GoodreadsRepository) GetBooksByTag(
	ctx context.Context,
	tag string,
) ([]goodreads.Book, error) {
	query := `
		SELECT id, shelf, tags, title, author, dates_read
		FROM goodreads_books 
		WHERE $1 = ANY(tags)
	`

	rows, err := repo.db.Query(ctx, query, tag)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	books := []goodreads.Book{}
	for rows.Next() {
		var book goodreads.Book

		err = rows.Scan(
			&book.ID,
			&book.Shelf,
			&book.Tags,
			&book.Title,
			&book.Author,
			&book.DatesRead,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return books, nil
}

func (repo GoodreadsRepository) GetBooksByIDs(
	ctx context.Context,
	ids []int64,
) ([]goodreads.Book, error) {
	query := `
		SELECT id, shelf, tags, title, author, dates_read
		FROM goodreads_books 
		WHERE id = ANY($1)
	`

	rows, err := repo.db.Query(ctx, query, ids)
	if err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	books := []goodreads.Book{}
	for rows.Next() {
		var book goodreads.Book

		err = rows.Scan(
			&book.ID,
			&book.Shelf,
			&book.Tags,
			&book.Title,
			&book.Author,
			&book.DatesRead,
		)

		if err != nil {
			return nil, postgres.PgxErrorToHTTPError(err)
		}

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, postgres.PgxErrorToHTTPError(err)
	}

	return books, nil
}

func (repo GoodreadsRepository) UpsertBook(
	ctx context.Context,
	book goodreads.Book,
) error {
	query := `
		INSERT INTO goodreads_books (id, shelf, tags, title, author, dates_read)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id)
		DO UPDATE SET shelf = $2, tags = $3, title = $4, author = $5, dates_read = $6
		RETURNING id
	`

	err := repo.db.QueryRow(
		ctx,
		query,
		book.ID,
		book.Shelf,
		book.Tags,
		book.Title,
		book.Author,
		book.DatesRead,
	).Scan(&book.ID)

	if err != nil {
		return postgres.PgxErrorToHTTPError(err)
	}

	return nil
}
