package repositories

import (
	"context"
	"log"

	"github.com/ImpressionableRaccoon/urlshortener/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func ConnectPSQL() {
	poolConfig, err := pgxpool.ParseConfig(configs.DatabaseDSN)
	if err != nil {
		log.Println("Unable to parse DATABASE_URL:", err)
	}

	db, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Println("Unable to create connection pool:", err)
	}
}

func PoolPSQL() bool {
	if db.Ping(context.Background()) == nil {
		return true
	}
	return false
}
