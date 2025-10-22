package grpchandlers

import (
	"github.com/sunr3d/quorum-grep/internal/interfaces/services"
	pbg "github.com/sunr3d/quorum-grep/proto/grepsvc"
)

var _ pbg.GrepServiceServer = (*handler)(nil)

type handler struct {
	pbg.UnimplementedGrepServiceServer
	svc services.GrepService
}

// New - конструктор handler.
func New(svc services.GrepService) pbg.GrepServiceServer {
	return &handler{
		svc: svc,
	}
}
