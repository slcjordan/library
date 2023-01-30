package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/slcjordan/library"
	"github.com/slcjordan/library/config"
	_ "github.com/slcjordan/library/config/envvar"
	_ "github.com/slcjordan/library/log/stdlib"
	"github.com/slcjordan/library/wire/api"
)

func main() {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		// repanic
		err, ok := r.(error)
		if !ok {
			panic(r)
		}
		var liberr *library.Error
		if errors.As(err, &liberr) {
			panic(r)
		}

		log.Fatalf("the app crashed: %s", err)
	}()

	config.MustParse()
	log.Printf("listening at %s", config.HTTP.ListenAddress)
	err := http.ListenAndServe(config.HTTP.ListenAddress, api.Wire())
	if err != nil {
		panic(&library.Error{
			Actual: err,
			Desc:   "while running server",
			Type:   library.Unknown,
		})
	}
}
