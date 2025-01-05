package mocks

import (
	"goal-tracker/api/pkg/goodreads"
)

type MockGoodreadsClient struct {
}

func NewMockGoodreadsClient() goodreads.Client {
	return MockGoodreadsClient{}
}

func (m MockGoodreadsClient) GetBooks(_ string) ([]goodreads.Book, error) {
	return []goodreads.Book{}, nil
}

func (m MockGoodreadsClient) GetUserID(_ string) (*string, error) {
	value := "userId"
	return &value, nil
}
