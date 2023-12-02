package decoder

import "time"

// Page represents a page from Wikipedia.
type Page struct {
	ID        int64     `xml:"-"`
	UpdatedAt time.Time `xml:"revision>timestamp"`
	Title     string    `xml:"title"`
	Text      string    `xml:"revision>text"`
}
