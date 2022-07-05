package db

import (
	"context"
	"log"
	"natsuki/utils"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mediocregopher/radix/v4"
)

var Redis radix.Client = createRedisClientPool()

func createRedisClientPool() radix.Client {

	// Retrieve Options from Environment
	var redisPass = utils.GetEnvDefault("NATSUKI_REDIS_PASS", "")
	var poolsize = utils.GetEnvDefaultInt("NATSUKI_REDIS_POOL_SIZE", 5)
	var redisAddr = utils.GetEnvDefault("NATSUKI_REDIS_ADDR", "127.0.0.1:6379")
	PoolConfig := radix.PoolConfig{
		Size: poolsize,
		Dialer: radix.Dialer{
			AuthPass: redisPass,
		},
	}

	// Create Redis Pool
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	client, err := PoolConfig.New(ctx, "tcp", redisAddr)
	if err != nil {
		log.Fatalln("unable to create redis pool", err.Error())
	}

	// Conduct Ping Test
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Do(ctx, radix.Cmd(nil, "PING")); err != nil {
		log.Fatalln("failed to verify redis ping", err.Error())
	}

	// Successfully created Redis Pool
	log.Printf("[RDB][INFO] Created Pool with %v Client(s)", poolsize)
	return client
}
