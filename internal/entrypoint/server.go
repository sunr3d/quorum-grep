package entrypoint

import (
	"context"

	"github.com/sunr3d/quorum-grep/internal/config"
	grpchandlers "github.com/sunr3d/quorum-grep/internal/handlers/grpc"
	"github.com/sunr3d/quorum-grep/internal/server"
	"github.com/sunr3d/quorum-grep/internal/services/grepsvc"
	pbg "github.com/sunr3d/quorum-grep/proto/grepsvc"
)

func RunServer(ctx context.Context, cfg *config.Config) error {
	svc := grepsvc.New()

	handler := grpchandlers.New(svc)

	srv := server.New(cfg)

	pbg.RegisterGrepServiceServer(srv.GetGRPCServer(), handler)

	return srv.Run(ctx)
}
