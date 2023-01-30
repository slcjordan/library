package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/slcjordan/library/config"
	_ "github.com/slcjordan/library/config/envvar"
	_ "github.com/slcjordan/library/log/stdlib"
	"github.com/slcjordan/library/wire/api"
	"github.com/slcjordan/oops"
)

type Action func(t *testing.T)

func TestIntegration(t *testing.T) {
	config.MustParse()
	config.Postgres.ConnectTimeout = 1 * time.Second
	dbURL, err := url.Parse(config.Postgres.ConnectionString)

	if err != nil {
		panic(fmt.Errorf("could not parse postgres connection string %w", err))
	}

	for _, test := range []struct {
		Desc              string
		NetworkConditions []oops.Condition
		Action            Action
	}{
		{
			Desc: "list books timeout",
			NetworkConditions: []oops.Condition{oops.ReadLatency(
				time.Millisecond, // connect handshake
				time.Millisecond,
				time.Millisecond,
				10*time.Second, // request
			)},
			Action: Do(httptest.NewRequest(
				http.MethodGet, "http://"+config.HTTP.ListenAddress+"/api/v1/books", nil,
			),
				LatencyLessThan(6*time.Second),
				StatusShouldBe(http.StatusGatewayTimeout),
			),
		},
		{
			Desc: "delete nonexisting book",
			Action: Do(httptest.NewRequest(
				http.MethodDelete, "http://"+config.HTTP.ListenAddress+"/api/v1/books/1234567890123", nil,
			), StatusShouldBe(http.StatusNoContent)),
		},
		{
			Desc: "create book happy path",
			Action: Do(httptest.NewRequest(
				http.MethodPost, "http://"+config.HTTP.ListenAddress+"/api/v1/books", strings.NewReader(`
				{
					"isbn": 1234567890123,
					"title": "Domain Driven Design"
				}
				`),
			), StatusShouldBe(http.StatusCreated)),
		},
		{
			Desc: "create book again is error",
			Action: Do(httptest.NewRequest(
				http.MethodPost, "http://"+config.HTTP.ListenAddress+"/api/v1/books", strings.NewReader(`
				{
					"isbn": 1234567890123,
					"title": "Domain Driven Design"
				}
				`),
			), StatusShouldBe(http.StatusBadRequest)),
		},
		{
			Desc: "list books happy path",
			Action: Do(httptest.NewRequest(
				http.MethodGet, "http://"+config.HTTP.ListenAddress+"/api/v1/books", nil,
			), StatusShouldBe(http.StatusOK)),
		},
		{
			Desc: "get book happy path",
			Action: Do(httptest.NewRequest(
				http.MethodGet, "http://"+config.HTTP.ListenAddress+"/api/v1/books/1234567890123", nil,
			), StatusShouldBe(http.StatusOK)),
		},
		{
			Desc: "update book happy path",
			Action: Do(httptest.NewRequest(
				http.MethodPut, "http://"+config.HTTP.ListenAddress+"/api/v1/books/1234567890123", strings.NewReader(`
				{
					"isbn": 1234567890123,
					"title": "Clean Code"
				}
				`),
			), StatusShouldBe(http.StatusOK)),
		},
		/*
		{
			Desc: "delete book cleanup",
			Action: Do(httptest.NewRequest(
				http.MethodDelete, "http://"+config.HTTP.ListenAddress+"/api/v1/books/1234567890123", nil,
			), StatusShouldBe(http.StatusNoContent)),
		},
		*/
	} {
		func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			proxy := NewTestTCPProxy(dbURL.Host, test.NetworkConditions)
			curr := *dbURL
			curr.Host = proxy.Addr
			config.Postgres.ConnectionString = curr.String()
			//nolint:errcheck
			go proxy.Run(ctx)
			t.Run(test.Desc, test.Action)
		}()
	}
}

type Validator func(*httptest.ResponseRecorder) error

func Do(r *http.Request, validators ...Validator) Action {
	return func(t *testing.T) {
		handler := api.Wire()
		resp := httptest.NewRecorder()
		handler.ServeHTTP(resp, r)

		for _, v := range validators {
			err := v(resp)
			if err != nil {
				t.Fatalf("while validating response: %s", err)
			}
		}
	}
}

func StatusShouldBe(status int) Validator {
	return func(resp *httptest.ResponseRecorder) error {
		if resp.Code != status {
			return fmt.Errorf("expected response status %d but got %d", status, resp.Code)
		}
		return nil
	}
}

func LatencyLessThan(timeout time.Duration) Validator {
	start := time.Now()
	return func(resp *httptest.ResponseRecorder) error {
		latency := time.Since(start)
		if latency > timeout {
			return fmt.Errorf("latency should have been less than %s but was %s", timeout, latency)
		}
		return nil
	}
}
