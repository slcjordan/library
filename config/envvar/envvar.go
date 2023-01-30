package envvar

import (
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/slcjordan/library"
	"github.com/slcjordan/library/config"
)

func init() {
	config.Register(0, MustParse)
}

func mustParseInt32(dest *int32, name string) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return
	}
	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		panic(&library.Error{
			Actual: err,
			Desc:   "while parsing int32 from env-var: " + name,
			Type:   library.InvalidSettings,
		})
	}
	*dest = int32(parsed)
}

func mustParseDuration(dest *time.Duration, name string) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return
	}
	var err error
	*dest, err = time.ParseDuration(value)
	if err != nil {
		panic(&library.Error{
			Actual: err,
			Desc:   "while parsing duration from env-var: " + name,
			Type:   library.InvalidSettings,
		})
	}
}

func mustMatchURL(dest *string, name string) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return
	}
	_, err := url.Parse(value)
	if err != nil {
		panic(&library.Error{
			Actual: err,
			Desc:   "while parsing url from env-var: " + name,
			Type:   library.InvalidSettings,
		})
	}
	*dest = value
}

func maybeSetString(dest *string, name string) {
	value, ok := os.LookupEnv(name)
	if !ok {
		return
	}
	*dest = value
}

// Sets config values from environment variables.
func MustParse() {
	mustParseDuration(&config.Postgres.ConnectTimeout, "LIBRARY_PG_CONNECT_TIMEOUT")
	mustMatchURL(&config.Postgres.ConnectionString, "LIBRARY_PG_CONNECTION_STRING")

	mustMatchURL(&config.HTTP.BaseURL, "LIBRARY_HTTP_BASE_URL")
	maybeSetString(&config.HTTP.ListenAddress, "LIBRARY_HTTP_LISTEN_ADDRESS")
	mustParseInt32(&config.HTTP.MaxListSize, "LIBRARY_HTTP_MAX_LIST_SIZE")
}
