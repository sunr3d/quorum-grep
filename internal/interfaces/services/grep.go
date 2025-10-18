package services

import (
	"context"

	"github.com/sunr3d/quorum-grep/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=GrepService --output=../../../mocks --filename=mock_grep_service.go --with-expecter
type GrepService interface {
	ProcessChunk(ctx context.Context, task *models.Task) (*models.Result, error)
}
