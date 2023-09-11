package config

import (
	"os"
	"strconv"
)

type (
	Config struct {
		RedisConfig      RedisConfig
		HTTPServerConfig HTTPServerConfig
		RabbitmqConfig   RabbitmqConfig
	}

	RedisConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		DB       int
	}

	HTTPServerConfig struct {
		Host string
		Port string
	}

	RabbitmqConfig struct {
		Host     string
		User     string
		Password string
	}
)

func ReadConfigFromShell() (Config, error) {

	redisDbField, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil && os.Getenv("REDIS_DB") != "" {
		return Config{}, err
	}

	return Config{
		RedisConfig: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			User:     os.Getenv("REDIS_USER"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDbField,
		},
		HTTPServerConfig: HTTPServerConfig{
			Host: os.Getenv("HTTP_SERVER_HOST"),
			Port: os.Getenv("HTTP_SERVER_PORT"),
		},
		RabbitmqConfig: RabbitmqConfig{
			Host:     os.Getenv("RABBITMQ_HOST"),
			User:     os.Getenv("RABBITMQ_USER"),
			Password: os.Getenv("RABBITMQ_PASSWORD"),
		},
	}, nil
}
