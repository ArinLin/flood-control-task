package service

import (
	"context"
	"fmt"

	"task/internal/config"
	"task/internal/store"

	"github.com/go-redis/redis/v8"
)

type (
	FloodControl interface {
		// Check возвращает false если достигнут лимит максимально разрешенного
		// кол-ва запросов согласно заданным правилам флуд контроля.
		Check(ctx context.Context, userID int64) (bool, error)
	}

	floodControlImpl struct {
		store  store.Store
		config *config.Config
	}
)

func New(store store.Store, config *config.Config) FloodControl {
	return &floodControlImpl{
		store:  store,
		config: config,
	}
}

func (fc *floodControlImpl) Check(ctx context.Context, userID int64) (bool, error) {
	key := fmt.Sprintf("user:%d:requests", userID)

	// Получаем текущее количество запросов из Redis
	count, err := fc.store.Get(ctx, key)
	if err != nil && err != redis.Nil {
		return false, err
	}
	if count >= fc.config.MaxRequests {
		return false, nil // Лимит превышен
	}

	// Увеличиваем счетчик и устанавливаем время жизни ключа, если это первый запрос, иначе просто увеличиваем счетчик
	if count == 0 {
		err = fc.store.IncrementWithExpiry(ctx, key, fc.config.Interval)
	} else {
		err = fc.store.Increment(ctx, key)
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
