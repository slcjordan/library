package config

import "time"

var Postgres struct {
	ConnectTimeout   time.Duration
	ConnectionString string
}

var HTTP struct {
	BaseURL       string
	ListenAddress string
	MaxListSize   int32
}
