package config

import (
	"fmt"

	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/zlog"
)

func GetConfig() (*Config, error) {
	cfg := config.New()
	if err := cfg.Load("config.yaml", ".env", ""); err != nil {
		zlog.Logger.Warn().Msgf("config.Load(): %v. Продолжаем с дефолтными значениями...", err)
	}

	setDefaults(cfg)

	var c Config
	if err := cfg.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("cfg.Unmarshal: %w", err)
	}

	return &c, nil
}

func setDefaults(cfg *config.Config) {
	cfg.SetDefault("LOG_LEVEL", "info")

	cfg.SetDefault("CLIENT.SERVER_LIST", []string{"localhost:50051", "localhost:50052", "localhost:50053"})
	cfg.SetDefault("CLIENT.TIMEOUT", "30s")
	cfg.SetDefault("CLIENT.CHUNK_SIZE", 1024)
}
