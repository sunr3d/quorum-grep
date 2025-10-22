package server

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/quorum-grep/internal/config"
)

const (
	ShutdownTimeout = 30 * time.Second
)

type Server struct {
	addr       string
	grpcServer *grpc.Server
}

// New - создает новый сервер gRPC.
func New(cfg *config.GRPCServerConfig) *Server {
	grpcServer := grpc.NewServer()

	return &Server{
		addr:       fmt.Sprintf(":%d", cfg.Port),
		grpcServer: grpcServer,
	}
}

// Run - запускает сервер gRPC с graceful shutdown.
func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("net.Listen %s: %w", s.addr, err)
	}

	srvErr := make(chan error, 1)
	go func() {
		zlog.Logger.Info().
			Str("addr", s.addr).
			Msg("Запуск gRPC сервера...")
		if err := s.grpcServer.Serve(listener); err != nil {
			srvErr <- fmt.Errorf("grpcServer.Serve: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		zlog.Logger.Info().
			Msg("Получен сигнал о завершении работы, инициализация graceful shutdown...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		done := make(chan struct{})
		go func() {
			defer close(done)
			s.grpcServer.GracefulStop()
		}()

		select {
		case <-done:
			zlog.Logger.Info().
				Msg("gRPC сервер успешно остановлен")
			return nil
		case <-shutdownCtx.Done():
			zlog.Logger.Warn().
				Msg("gRPC сервер остановлен по таймауту")
			s.grpcServer.Stop()
			return fmt.Errorf("graceful shutdown таймаут")
		}

	case err := <-srvErr:
		zlog.Logger.Error().
			Err(err).
			Msg("gRPC сервер остановлен с ошибкой")
		return fmt.Errorf("gRPC сервер остановлен с ошибкой: %w", err)
	}
}

func (s *Server) GetGRPCServer() *grpc.Server {
	return s.grpcServer
}
