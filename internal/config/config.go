package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type APIConfig struct {
	Debug  bool   `envconfig:"DEBUG" default:"true"`
	Host   string `envconfig:"HOST" default:"0.0.0.0"`
	Port   string `envconfig:"PORT" default:"8000"`
	Domain string `envconfig:"DOMAIN" default:""`

	GasGraphTtl     time.Duration `envconfig:"GAS_GRAPH_CACHE_TTL" default:"30s"`
	TransactionsUrl string        `envconfig:"TRANSACTIONS_URL" default:"https://github.com/CryptoRStar/GasPriceTestTask/raw/main/gas_price.json"`
}

type Config struct {
	API *APIConfig
}

var (
	once   sync.Once
	config *Config
)

// GetConfig Загружает конфиг из .env файла и возвращает объект конфигурации
// В случае, если не передать параметр `envfiles`, берется `.env` файл из корня проекта.
func GetConfig(envfiles ...string) (*Config, error) {
	var err error
	once.Do(func() {
		_ = godotenv.Load(envfiles...)

		var c Config
		err = envconfig.Process("", &c)
		if err != nil {
			err = fmt.Errorf("error parse config from env variables: %w", err)
			return
		}
		config = &c
	})

	return config, err
}
