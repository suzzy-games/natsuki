package db

import (
	"context"
	"log"
	"natsuki/utils"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

var Postgres *pgxpool.Pool = createPostgresClientPool()

func createPostgresClientPool() *pgxpool.Pool {
	// Retrieve Options from Environment
	postgUrl := utils.GetEnvDefault("NATSUKI_POSTGRES_URL", "postgres://postgres:password@localhost:5432/")
	poolsize := utils.GetEnvDefaultInt("NATSUKI_POSTGRES_POOL_SIZE", 5)

	// Configure Database
	connConfig, err := pgxpool.ParseConfig(postgUrl)
	if err != nil {
		log.Fatalln("unable to parse postgres url", err)
	}

	connConfig.MaxConns = int32(poolsize)
	connConfig.MinConns = int32(poolsize)

	// Create Database Pool
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	pool, err := pgxpool.ConnectConfig(ctx, connConfig)
	if err != nil {
		log.Fatalln("unable to create postgres pool", err.Error())
	}

	// Conduct Ping Test
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalln("failed to verify postgres ping", err.Error())
	}

	// Successfully created MySQL Pool
	log.Printf("[SQL][INFO] Created Pool with %v Client(s)", poolsize)
	return pool
}
