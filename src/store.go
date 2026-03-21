package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type LineStore interface {
	SaveLine(ctx context.Context, line Line) error
	ListLines(ctx context.Context) ([]Line, error)
}

func NewLineStoreFromEnv(ctx context.Context) (LineStore, error) {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("LINE_STORE"))) {
	case "", "memory":
		return NewMemoryLineStore(), nil
	case "redis":
		return NewRedisLineStore(ctx, redis.Options{
			Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       envIntOrDefault("REDIS_DB", 0),
		})
	default:
		return nil, fmt.Errorf("unsupported LINE_STORE value %q", os.Getenv("LINE_STORE"))
	}
}

func envOrDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func envIntOrDefault(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func newRedisClient(options redis.Options) *redis.Client {
	return redis.NewClient(&options)
}
