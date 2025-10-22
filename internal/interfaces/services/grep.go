package services

import (
	"context"

	"github.com/sunr3d/quorum-grep/models"
)

type GrepService interface {
	ProcessChunk(ctx context.Context, task *models.Task) (*models.Result, error)
}
