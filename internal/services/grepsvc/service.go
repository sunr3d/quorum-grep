package grepsvc

import (
	"bytes"
	"context"

	"github.com/sunr3d/quorum-grep/internal/interfaces/services"
	"github.com/sunr3d/quorum-grep/models"
)

var _ services.GrepService = (*grepService)(nil)

type grepService struct{}

func New() services.GrepService {
	return &grepService{}
}

func (s *grepService) ProcessChunk(ctx context.Context, task *models.Task) (*models.Result, error) {
	lines := bytes.Split(task.Data, []byte("\n"))

	return &models.Result{
		TaskID:     task.ID,
		MatchCount: len(lines),
		Error:      "",
	}, nil
}
