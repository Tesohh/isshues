package db_complex

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// https://github.com/sqlc-dev/sqlc/issues/2116

func GetDbPool(connectionString string) (*pgxpool.Pool, error) {
	// Set up a new pool with the custom types and the config.
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), config)

	// Collect the custom data types once, store them in memory, and register them for every future connection.
	customTypes, err := GetCustomDataTypes(context.Background(), dbpool)
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		for _, t := range customTypes {
			conn.TypeMap().RegisterType(t)
		}
		return nil
	}
	// Immediately close the old pool and open a new one with the new config.
	dbpool.Close()
	dbpool, err = pgxpool.NewWithConfig(context.Background(), config)
	return dbpool, err
}
