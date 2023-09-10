package app

import (
	"log"
	"ml-in-kube-apiserver/internal/config"
	"ml-in-kube-apiserver/internal/delivery"

	app_redis "ml-in-kube-apiserver/internal/db/redis"

	"github.com/joho/godotenv"
)

func Run() {
	godotenv.Load() // load configs from .env file (useful for development purposes)

	cfg, err := config.ReadConfigFromShell() // load configs from shell, like a professional, this overwrites .env file
	if err != nil {
		log.Fatal("Error in reading config")
	}

	redisConn, err := app_redis.NewRedisConnection(cfg.RedisConfig)
	if err != nil {
		log.Fatal("Error in connecting to redis ", err)
	}

	httpHandler := delivery.NewHandler(cfg.HTTPServerConfig)
	httpHandler.SetImgHandler(redisConn)
	httpHandler.StartServer()
}
