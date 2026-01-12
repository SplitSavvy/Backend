package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(connStr string) (*pgxpool.Pool, error){
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil{
		return nil, err
	}
	err = pool.Ping(context.Background())
	if err != nil{
		return nil, err
	}
	DB = pool
	return pool, nil
}