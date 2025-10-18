package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/quorum-grep/internal/config"
	"github.com/sunr3d/quorum-grep/internal/entrypoint"
)

func main() {
	zlog.Init()
	zlog.Logger.Info().Msg("Запуск сервера grep...")

	cfg, err := config.GetConfig()
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Msg("config.GetConfig")
	}

	zlog.SetLevel(cfg.LogLevel)
	zlog.Logger.Info().Msgf("cfg: %+v", cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := entrypoint.RunServer(ctx, cfg); err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Msg("entrypoint.RunServer")
	}

}
