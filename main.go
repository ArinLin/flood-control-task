package main

import (
	"context"
	"log"
	"math/rand"

	"task/internal/config"
	"task/internal/service"
	"task/internal/store"

	"github.com/go-redis/redis/v8"
)

func main() {
	// конфигурируем сервис со всеми необходимыми зависимостями
	ctx := context.Background()
	cfg, err := config.New()
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}
	log.Println(cfg)

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.StoreAddress, // адрес и порт сервера Redis
	})

	store := store.New(redisClient)
	floodControl := service.New(store, cfg)

	// проверяем что рейтлимитер отрабатывает правильно для рандомных пользователей
	for i := 0; i < cfg.MaxRequests+50; i++ {
		userID := rand.Int63n(5)
		result, err := floodControl.Check(ctx, userID)
		if err != nil {
			log.Fatal("Error checking flood control: ", err)
			return
		}

		log.Printf("Flood control for user with id: %d, check result: %t", userID, result)
	}
}
