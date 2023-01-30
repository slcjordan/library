package library

// A Book is uniquely identified by its ISBN.
type Book struct {
	ISBN  int64
	Title string
}

// A BookList includes a next-page token for picking up at the next page.
type BookList struct {
	Books         []Book
	NextPageToken string
}
