package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type (
	Store interface {
		Get(ctx context.Context, key string) (int, error)
		Increment(ctx context.Context, key string) error
		IncrementWithExpiry(ctx context.Context, key string, exp time.Duration) error
	}

	storeImpl struct {
		client *redis.Client
	}
)

func New(client *redis.Client) Store {
	return &storeImpl{
		client: client,
	}
}

func (s *storeImpl) Get(ctx context.Context, key string) (int, error) {
	return s.client.Get(ctx, key).Int()
}

func (s *storeImpl) Increment(ctx context.Context, key string) error {
	_, err := s.client.Incr(ctx, key).Result() // Просто увеличиваем значение счетчика
	return err
}

func (s *storeImpl) IncrementWithExpiry(ctx context.Context, key string, exp time.Duration) error {
	pipeline := s.client.Pipeline()
	pipeline.Incr(ctx, key)        // Увеличиваем значение счетчика
	pipeline.Expire(ctx, key, exp) // Устанавливаем время жизни ключа
	_, err := pipeline.Exec(ctx)   // Выполняем обе операции атомарно
	return err
}
