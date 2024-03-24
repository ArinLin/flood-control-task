package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	StoreAddress string        `env:"STORE_ADDR" envDefault:"localhost:6379"` // Адрес хранилища
	MaxRequests  int           `env:"MAX_REQUESTS" envDefault:"10"`           // Максимальное количество запросов
	Interval     time.Duration `env:"INTERVAL" envDefault:"1m"`               // Интервал времени для максимального количества запросов
}

func New() (*Config, error) {
	var c Config

	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &c, nil
}
