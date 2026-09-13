package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedisLineStore struct {
	client *redis.Client
}

func NewRedisLineStore(ctx context.Context, options redis.Options) (*RedisLineStore, error) {
	client := newRedisClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	return &RedisLineStore{client: client}, nil
}

func redisSessionKey(sessionID string, suffix string) string {
	var builder strings.Builder
	builder.Grow(len("sharedcanvas:sessions::") + len(sessionID) + len(suffix))
	builder.WriteString("sharedcanvas:sessions:")
	builder.WriteString(sessionID)
	builder.WriteByte(':')
	builder.WriteString(suffix)

	return builder.String()
}

func redisLinesKey(sessionID string) string {
	return redisSessionKey(sessionID, "lines")
}

func redisNextIDKey(sessionID string) string {
	return redisSessionKey(sessionID, "next_id")
}

func (store *RedisLineStore) NextLineID(ctx context.Context, sessionID string) (int, error) {
	id, err := store.client.Incr(ctx, redisNextIDKey(sessionID)).Result()
	if err != nil {
		return 0, fmt.Errorf("increment next id: %w", err)
	}

	return int(id), nil
}

func (store *RedisLineStore) SaveLine(ctx context.Context, sessionID string, line *Line) error {
	id, err := store.NextLineID(ctx, sessionID)
	if err != nil {
		return err
	}

	line.ID = id

	payload, err := json.Marshal(line)
	if err != nil {
		return fmt.Errorf("marshal line: %w", err)
	}

	if err := store.client.RPush(ctx, redisLinesKey(sessionID), payload).Err(); err != nil {
		return fmt.Errorf("save line: %w", err)
	}

	return nil
}

func (store *RedisLineStore) ListLines(ctx context.Context, sessionID string) ([]Line, error) {
	items, err := store.client.LRange(ctx, redisLinesKey(sessionID), 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("list lines: %w", err)
	}

	lines := make([]Line, 0, len(items))
	for _, item := range items {
		var line Line
		if err := json.Unmarshal([]byte(item), &line); err != nil {
			return nil, fmt.Errorf("decode line: %w", err)
		}
		lines = append(lines, line)
	}

	return lines, nil
}

func (store *RedisLineStore) ListLinesSince(ctx context.Context, sessionID string, sinceID int) ([]Line, error) {
	lines, err := store.ListLines(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	filtered := make([]Line, 0, len(lines))
	for _, line := range lines {
		if line.ID > sinceID {
			filtered = append(filtered, line)
		}
	}

	return filtered, nil
}

func (store *RedisLineStore) ClearLines(ctx context.Context, sessionID string) error {
	if err := store.client.Del(ctx, redisLinesKey(sessionID)).Err(); err != nil {
		return fmt.Errorf("clear lines: %w", err)
	}

	return nil
}
