package main

import (
	"context"
	"flag"
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

	port := flag.Int("port", 50051, "порт для запуска сервера")
	flag.Parse()

	cfg := &config.GRPCServerConfig{
		Port: *port,
	}

	zlog.Logger.Info().Msgf("cfg: %+v", cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := entrypoint.RunServer(ctx, cfg); err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Msg("entrypoint.RunServer")
	}
}
