package goodreads

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

func GetUserID(profileURL string) (*string, error) {
	c := colly.NewCollector()

	var userID string
	c.OnHTML(".profilePictureIcon", func(h *colly.HTMLElement) {
		imgURL := h.Attr("src")
		splittedSlash := strings.Split(imgURL, "/")
		userID = strings.Split(splittedSlash[len(splittedSlash)-1], ".jpg")[0]
	})

	err := c.Visit(profileURL)
	if err != nil {
		return nil, err
	}

	return &userID, nil
}

func GetBooks(userID string, shelf string) ([]Book, error) {
	books := []Book{}

	page := 0
	for {
		page++

		booksOnPage, err := getBooksFromPage(userID, shelf, page)
		if err != nil {
			return nil, err
		}

		books = append(books, booksOnPage...)

		if len(booksOnPage) == 0 {
			break
		}
	}

	return books, nil
}

func getBooksFromPage(userID string, shelf string, page int) ([]Book, error) {
	c := colly.NewCollector()

	books := []Book{}
	c.OnHTML(".bookalike.review", func(h *colly.HTMLElement) {
		book := Book{
			Title:  h.ChildText(".title .value a"),
			Author: h.ChildText(".author .value a"),
		}

		books = append(books, book)
	})

	err := c.Visit(
		fmt.Sprintf(
			"https://www.goodreads.com/review/list/%s?shelf=%s&page=%d",
			userID,
			shelf,
			page,
		),
	)
	if err != nil {
		return nil, err
	}

	return books, nil
}
