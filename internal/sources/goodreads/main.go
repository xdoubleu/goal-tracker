package goodreads

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// todo move to repositories
func GetUserId(profileUrl string) string {
	c := colly.NewCollector()

	var userId string
	c.OnHTML(".profilePictureIcon", func(h *colly.HTMLElement) {
		imgUrl := h.Attr("src")
		splittedSlash := strings.Split(imgUrl, "/")
		userId = strings.Split(splittedSlash[len(splittedSlash)-1], ".jpg")[0]
	})

	c.Visit(profileUrl)

	return userId
}

func GetBooks(userId string, shelf string) []Book {
	books := []Book{}

	page := 0
	for {
		page++

		booksOnPage := getBooksFromPage(userId, shelf, page)
		books = append(books, booksOnPage...)

		fmt.Printf("%d, %d\n", page, len(booksOnPage))
		if len(booksOnPage) == 0 {
			break
		}
	}

	return books
}

func getBooksFromPage(userId string, shelf string, page int) []Book {
	c := colly.NewCollector()

	books := []Book{}
	c.OnHTML(".bookalike.review", func(h *colly.HTMLElement) {
		book := Book{
			Title:  h.ChildText(".title .value a"),
			Author: h.ChildText(".author .value a"),
		}

		books = append(books, book)
	})

	c.Visit(fmt.Sprintf("https://www.goodreads.com/review/list/%s?shelf=%s&page=%d", userId, shelf, page))

	return books
}
