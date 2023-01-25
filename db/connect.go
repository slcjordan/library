package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/slcjordan/library/config"
	"github.com/slcjordan/library/db/sqlc"
)

//go:generate go run github.com/golang/mock/mockgen -package=db -destination=../test/mocks/db/connect.go -source=connect.go

type DBTX interface {
	sqlc.DBTX
}

func MustConnect() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), config.Postgres.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, config.Postgres.ConnectionString)
	if err != nil {
		panic(err)
	}
	return pool
}
