package goodreads

import "time"

type Book struct {
	Title     string
	Author    string
	DatesRead []time.Time
}
