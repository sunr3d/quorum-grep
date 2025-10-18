package grpchandlers

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/sunr3d/quorum-grep/models"
	pbg "github.com/sunr3d/quorum-grep/proto/grepsvc"
)

func (h *handler) ProcessChunk(ctx context.Context, req *pbg.ChunkRequest) (*pbg.ChunkResponse, error) {
	task := &models.Task{
		ID:    req.TaskId,
		Data:  req.Data,
		Index: int(req.ChunkIndex),
		Options: models.GrepOptions{
			Pattern:    req.Options.Pattern,
			After:      int(req.Options.After),
			Before:     int(req.Options.Before),
			Around:     int(req.Options.Around),
			Count:      req.Options.Count,
			IgnoreCase: req.Options.IgnoreCase,
			Invert:     req.Options.Invert,
			Fixed:      req.Options.Fixed,
			LineNum:    req.Options.LineNum,
		},
	}

	result, err := h.svc.ProcessChunk(ctx, task)
	if err != nil {
		return &pbg.ChunkResponse{
			TaskId: req.TaskId,
			Error:  err.Error(),
		}, nil
	}

	matches := make([]*pbg.Match, len(result.Matches))
	for i, match := range result.Matches {
		matches[i] = &pbg.Match{
			LineNum:       int32(match.LineNum),
			Content:       match.Content,
			ContextBefore: match.ContextBefore,
			ContextAfter:  match.ContextAfter,
		}
	}

	return &pbg.ChunkResponse{
		TaskId:     result.TaskID,
		Matches:    matches,
		MatchCount: int32(result.MatchCount),
		Error:      result.Error,
	}, nil
}

func (h *handler) HealthCheck(ctx context.Context, _ *emptypb.Empty) (*pbg.HealthResponse, error) {
	return &pbg.HealthResponse{
		Ok: true,
	}, nil
}
