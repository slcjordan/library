package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/slcjordan/library"
	"github.com/slcjordan/library/config"
	"github.com/slcjordan/library/log"
)

//go:generate go run github.com/golang/mock/mockgen -package=http -destination=../test/mocks/http/http.go -source=http.go

type ListBooksController interface {
	ListBooks(ctx context.Context, PageToken string, TotalSize int32) (library.BookList, error)
}

type BookCRUDController interface {
	CreateBook(ctx context.Context, book library.Book) error
	DeleteBook(ctx context.Context, isbn int64) error
	GetBook(ctx context.Context, isbn int64) (library.Book, error)
	UpdateBook(ctx context.Context, book library.Book) error
}

func fromPtr[V any](input *V, otherwise V) V {
	if input == nil {
		return otherwise
	}
	return *input
}

// Server handles incoming requests.
type Server struct {
	ListBooksController ListBooksController
	BookCRUDController  BookCRUDController
}

func (s *Server) serialize(ctx context.Context, w http.ResponseWriter, data any) {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		log.Errorf(ctx, "while encoding %T: %s", data, err)
	}
}

func (s *Server) reportError(ctx context.Context, w http.ResponseWriter, err error) {
	var libErr *library.Error
	if errors.As(err, &libErr) {
		switch libErr.Type {
		case library.BadInput:
			w.WriteHeader(http.StatusBadRequest)
			s.serialize(ctx, w, Error{
				Code:    int(libErr.Type),
				Message: fmt.Sprintf("Bad Input: %s", err),
			})
			return
		case library.Timeout:
			w.WriteHeader(http.StatusGatewayTimeout)
			s.serialize(ctx, w, Error{
				Code:    int(libErr.Type),
				Message: "request timed out",
			})
			return
		default:
			// fall through
		}
	}
	log.Errorf(ctx, "unknown error during request handling: %s", err)
	w.WriteHeader(
		http.StatusInternalServerError,
	)
	s.serialize(ctx, w, Error{
		Code:    int(library.Unknown),
		Message: "Internal error. Check logs for details.",
	})
}

// UserErrorHandler handles unexpected errors.
func (s *Server) UserErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	s.reportError(r.Context(), w, &library.Error{
		Type:   library.BadInput,
		Actual: err,
		Desc:   "while serving request",
	})
}

// ListBooks returns a result of books in the library.
func (s *Server) ListBooks(w http.ResponseWriter, r *http.Request, params ListBooksParams) {
	ctx := r.Context()
	totalSize := fromPtr(params.TotalSize, config.HTTP.MaxListSize)
	if totalSize < 0 || totalSize > config.HTTP.MaxListSize {
		s.reportError(ctx, w, &library.Error{
			Type:   library.BadInput,
			Desc:   "while checking parameter bounds",
			Actual: fmt.Errorf("total size should be between %d and %d but got %d", 0, config.HTTP.MaxListSize, totalSize),
		})
		return
	}
	bookList, err := s.ListBooksController.ListBooks(ctx, fromPtr(params.PageToken, ""), totalSize)
	if err != nil {
		s.reportError(ctx, w, err)
		return
	}
	result := BookList{
		Items:         make([]Book, 0, len(bookList.Books)),
		NextPageToken: bookList.NextPageToken,
	}
	for _, b := range bookList.Books {
		result.Items = append(result.Items, Book{
			Isbn:  b.ISBN,
			Title: b.Title,
		})
	}
	s.serialize(ctx, w, result)
}

// CreateBook adds a single book to the library.
func (s *Server) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var book Book
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		s.reportError(ctx, w, &library.Error{
			Type:   library.BadInput,
			Desc:   "while parsing request body",
			Actual: err,
		})
		return
	}
	err = s.BookCRUDController.CreateBook(ctx, library.Book{
		Title: book.Title,
		ISBN:  book.Isbn,
	})
	if err != nil {
		s.reportError(ctx, w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// DeleteBook handles deleting a book
func (s *Server) DeleteBook(w http.ResponseWriter, r *http.Request, isbn Isbn) {
	ctx := r.Context()
	err := s.BookCRUDController.DeleteBook(ctx, isbn)
	if err != nil {
		s.reportError(ctx, w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// FetchBook handles fetching a book
func (s *Server) FetchBook(w http.ResponseWriter, r *http.Request, isbn Isbn) {
	ctx := r.Context()
	book, err := s.BookCRUDController.GetBook(ctx, isbn)
	if err != nil {
		s.reportError(ctx, w, err)
		return
	}

	result := Book{
		Isbn:  book.ISBN,
		Title: book.Title,
	}
	s.serialize(ctx, w, result)
	w.WriteHeader(http.StatusNoContent)
}

// UpdateBook handles updating a book
func (s *Server) UpdateBook(w http.ResponseWriter, r *http.Request, isbn Isbn) {
	ctx := r.Context()
	var book BookPartial
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&book)
	if err != nil {
		s.reportError(ctx, w, &library.Error{
			Type:   library.BadInput,
			Desc:   "while parsing request body",
			Actual: err,
		})
		return
	}
	err = s.BookCRUDController.UpdateBook(ctx, library.Book{
		ISBN:  isbn,
		Title: book.Title,
	})
	if err != nil {
		s.reportError(ctx, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
