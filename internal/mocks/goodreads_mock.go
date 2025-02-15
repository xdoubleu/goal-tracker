package mocks

import (
	"goal-tracker/api/pkg/goodreads"
	"time"
)

type MockGoodreadsClient struct {
}

func NewMockGoodreadsClient() goodreads.Client {
	return MockGoodreadsClient{}
}

func (m MockGoodreadsClient) GetBooks(_ string) ([]goodreads.Book, error) {
	return []goodreads.Book{
		{
			ID:        1,
			Shelf:     "shelf",
			Tags:      []string{"tag1"},
			Title:     "Title",
			Author:    "Author",
			DatesRead: []time.Time{time.Now().AddDate(-1, 0, 0)},
		},
		{
			ID:        2,
			Shelf:     "shelf",
			Tags:      []string{"tag1"},
			Title:     "Title2",
			Author:    "Author",
			DatesRead: []time.Time{time.Now()},
		},
	}, nil
}

func (m MockGoodreadsClient) GetUserID(_ string) (*string, error) {
	value := "userId"
	return &value, nil
}
