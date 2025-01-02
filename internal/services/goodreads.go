package services

import "goal-tracker/api/pkg/goodreads"

type GoodreadsService struct {
	profileURL string
}

func (service GoodreadsService) GetAllBooks() ([]goodreads.Book, error) {
	userID, err := goodreads.GetUserID(service.profileURL)
	if err != nil {
		return nil, err
	}

	return goodreads.GetBooks(*userID)
}
