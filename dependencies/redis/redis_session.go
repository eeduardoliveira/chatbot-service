package redis

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/redis/go-redis/v9"
)

type RedisSessionRepository struct {
    client *redis.Client
}

func NewRedisSessionRepository(client *redis.Client) *RedisSessionRepository {
    return &RedisSessionRepository{client: client}
}

func (r *RedisSessionRepository) GetSessionID(ctx context.Context, phoneNumber string) (string, error) {
    key := fmt.Sprintf("session:phone:%s", phoneNumber)
    return r.client.Get(ctx, key).Result()
}

func (r *RedisSessionRepository) SaveSessionID(ctx context.Context, phoneNumber, sessionID string) error {
    key := fmt.Sprintf("session:phone:%s", phoneNumber)
    return r.client.Set(ctx, key, sessionID, 24*time.Hour).Err()
}

func (r *RedisSessionRepository) GetOrCreateSessionID(ctx context.Context, phoneNumber string) (string, error) {
    sessionID, err := r.GetSessionID(ctx, phoneNumber)
    if err == redis.Nil {
        newSessionID := uuid.New().String()
        err := r.SaveSessionID(ctx, phoneNumber, newSessionID)
        return newSessionID, err
    } else if err != nil {
        return "", err
    }
    return sessionID, nil
}