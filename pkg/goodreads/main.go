package goodreads

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

//TODO add client
//TODO add logger

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

func GetBooks(userID string) ([]Book, error) {
	books := []Book{}

	page := 0
	for {
		page++

		fmt.Printf("fetching page %d\n", page)
		booksOnPage, err := getBooksFromPage(userID, page)
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

func getBooksFromPage(userID string, page int) ([]Book, error) {
	c := colly.NewCollector()

	books := []Book{}
	c.OnHTML(".bookalike.review", func(h *colly.HTMLElement) {
		book := Book{
			Title:     h.ChildText(".title .value a"),
			Author:    h.ChildText(".author .value a"),
			DatesRead: []time.Time{},
		}

		dateReadStrs := h.ChildTexts(".date_read .value span")
		for _, dateReadStr := range dateReadStrs {
			if dateReadStr == "not set" {
				continue
			}

			possibleDateFormats := []string{
				"Jan 02, 2006",
				"Jan 2006",
			}

			var dateRead time.Time
			var err error
			for _, dateFormat := range possibleDateFormats {
				dateRead, err = time.Parse(dateFormat, dateReadStr)
				if err == nil {
					break
				}
			}

			if err != nil {
				panic(err)
			}

			book.DatesRead = append(book.DatesRead, dateRead)
		}

		books = append(books, book)
	})

	err := c.Visit(
		fmt.Sprintf(
			"https://www.goodreads.com/review/list/%s?page=%d",
			userID,
			page,
		),
	)
	if err != nil {
		time.Sleep(time.Second)
		return getBooksFromPage(userID, page)
	}

	return books, nil
}
