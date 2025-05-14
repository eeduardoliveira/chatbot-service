package redis

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
    rediscontext "chatbot-service/app/domain/context"
)

type RedisContextRepository struct {
    client *redis.Client
}

func NewRedisContextRepository(client *redis.Client) *RedisContextRepository {
    return &RedisContextRepository{client: client}
}

func (r *RedisContextRepository) LoadContext(ctx context.Context, sessionID string) (rediscontext.ChatContext, error) {
    key := fmt.Sprintf("context:session:%s", sessionID)
    val, err := r.client.Get(ctx, key).Result()
    if err == redis.Nil {
        return rediscontext.ChatContext{}, nil
    }
    if err != nil {
        return rediscontext.ChatContext{}, err
    }

    var chatCtx rediscontext.ChatContext
    json.Unmarshal([]byte(val), &chatCtx)
    return chatCtx, nil
}

func (r *RedisContextRepository) SaveContext(ctx context.Context, sessionID string, chatCtx rediscontext.ChatContext) error {
    key := fmt.Sprintf("context:session:%s", sessionID)
    bytes, _ := json.Marshal(chatCtx)
    return r.client.Set(ctx, key, bytes, 24*time.Hour).Err()
}