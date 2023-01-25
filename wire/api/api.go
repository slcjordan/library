package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/slcjordan/library/config"
	"github.com/slcjordan/library/db"
	libhttp "github.com/slcjordan/library/http"
)

// Wire wires up admin dependencies.
func Wire() http.Handler {
	router := chi.NewRouter()
	conn := db.MustConnect()
	queryer := &db.Queryer{
		DBTX: conn,
	}
	server := &libhttp.Server{
		ListBooksController: queryer,
		BookCRUDController: queryer,
	}
	options := libhttp.ChiServerOptions{
		BaseURL:    config.HTTP.BaseURL,
		BaseRouter: router,
		Middlewares: []libhttp.MiddlewareFunc{
			// TODO add more, including throttling
			middleware.Logger,
			middleware.Recoverer,
			middleware.Timeout(4 * time.Second),
		},
		ErrorHandlerFunc: server.UserErrorHandler,
	}
	handler := libhttp.HandlerWithOptions(server, options)
	return handler
}
