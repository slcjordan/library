package db

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"

	"github.com/slcjordan/library"
	"github.com/slcjordan/library/db/sqlc"
	"github.com/slcjordan/library/log"
)

type Queryer struct {
	DBTX DBTX
}

// ListBooks returns a list of books.
func (q *Queryer) ListBooks(ctx context.Context, PageToken string, TotalSize int32) (library.BookList, error) {
	params := sqlc.ListBooksParams{
		PageToken: PageToken,
		TotalSize: TotalSize,
	}
	books, err := sqlc.New(q.DBTX).ListBooks(ctx, params)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return library.BookList{}, &library.Error{
				Type:   library.Timeout,
				Actual: err,
				Desc:   "while retrieving a list of books",
			}
		}
		return library.BookList{}, &library.Error{
			Type:   library.DatabaseError,
			Actual: err,
			Desc:   "while retrieving a list of books",
		}
	}
	return toBookList(books), nil
}

func toBookList(books []sqlc.Book) library.BookList {
	var result library.BookList
	for _, b := range books {
		result.Books = append(result.Books, library.Book{
			ISBN:  b.Isbn,
			Title: b.Title,
		})
		result.NextPageToken = b.Title // next result is after the last item
	}
	return result
}

// CreateBook creates a single book.
func (q *Queryer) CreateBook(ctx context.Context, book library.Book) error {
	params := sqlc.CreateBookParams{
		Isbn:  book.ISBN,
		Title: book.Title,
	}

	err := sqlc.New(q.DBTX).CreateBook(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				log.Infof(ctx, "while creating book: %s", err)
				return &library.Error{
					Type:   library.BadInput,
					Actual: err,
					Desc:   "a book with that isbn already exists",
				}
			}
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return &library.Error{
				Type:   library.Timeout,
				Actual: err,
				Desc:   "while creating a book",
			}
		}
		return &library.Error{
			Type:   library.DatabaseError,
			Actual: err,
			Desc:   "while creating a book",
		}
	}
	return nil
}

// DeleteBook deletes a single book.
func (q *Queryer) DeleteBook(ctx context.Context, isbn int64) error {
	err := sqlc.New(q.DBTX).DeleteBook(ctx, isbn)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &library.Error{
				Type:   library.Timeout,
				Actual: err,
				Desc:   "while deleting a book",
			}
		}
		return &library.Error{
			Type:   library.DatabaseError,
			Actual: err,
			Desc:   "while deleting a book",
		}
	}
	return nil
}

// GetBook fetches a single book.
func (q *Queryer) GetBook(ctx context.Context, isbn int64) (library.Book, error) {
	book, err := sqlc.New(q.DBTX).GetBook(ctx, isbn)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return library.Book{}, &library.Error{
				Type:   library.Timeout,
				Actual: err,
				Desc:   "while fetching a book",
			}
		}
		return library.Book{}, &library.Error{
			Type:   library.DatabaseError,
			Actual: err,
			Desc:   "while fetching a book",
		}
	}
	return library.Book{
		Title: book.Title,
		ISBN:  book.Isbn,
	}, nil
}

// UpdateBook updates a single book.
func (q *Queryer) UpdateBook(ctx context.Context, book library.Book) error {
	params := sqlc.UpdateBookParams{
		Isbn:  book.ISBN,
		Title: book.Title,
	}

	err := sqlc.New(q.DBTX).UpdateBook(ctx, params)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &library.Error{
				Type:   library.Timeout,
				Actual: err,
				Desc:   "while updating a book",
			}
		}
		return &library.Error{
			Type:   library.DatabaseError,
			Actual: err,
			Desc:   "while updating a book",
		}
	}
	return nil
}
