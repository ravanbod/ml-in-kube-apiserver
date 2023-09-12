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
		S3Config         S3Config
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

	S3Config struct {
		Host      string
		AccessKey string
		SecretKey string
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
		S3Config: S3Config{
			Host:      os.Getenv("S3_URL"),
			AccessKey: os.Getenv("S3_ACCESS_KEY"),
			SecretKey: os.Getenv("S3_SECRET_KEY"),
		},
	}, nil
}
