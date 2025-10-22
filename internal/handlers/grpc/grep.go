package grpchandlers

import (
	"context"

	"github.com/wb-go/wbf/zlog"

	"github.com/sunr3d/quorum-grep/models"
	pbg "github.com/sunr3d/quorum-grep/proto/grepsvc"
)

// ProcessChunk - ручка gRPC для обработки куска данных.
func (h *handler) ProcessChunk(ctx context.Context, req *pbg.ChunkRequest) (*pbg.ChunkResponse, error) {
	zlog.Logger.Info().
		Str("task_id", req.TaskId).
		Int("chunk_index", int(req.ChunkIndex)).
		Int("data_size", len(req.Data)).
		Str("pattern", req.Options.Pattern).
		Msg("Получен запрос на обработку куска данных")

	task := &models.Task{
		Data:        req.Data,
		Index:       int(req.ChunkIndex),
		LineNumbers: req.LineNumbers,
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
		zlog.Logger.Error().
			Err(err).
			Str("task_id", req.TaskId).
			Err(err).
			Msg("Ошибка при обработке куска данных")
		return &pbg.ChunkResponse{
			TaskId: req.TaskId,
		}, nil
	}

	matches := make([]*pbg.Match, len(result.Matches))
	for i, match := range result.Matches {
		matches[i] = &pbg.Match{
			Content:    match.Content,
			LineNumber: match.LineNumber,
		}
	}

	zlog.Logger.Info().
		Str("task_id", req.TaskId).
		Int("matches_count", len(result.Matches)).
		Msg("Кусок данных обработан")

	return &pbg.ChunkResponse{
		TaskId:     req.TaskId,
		Matches:    matches,
		MatchCount: int64(len(result.Matches)),
	}, nil
}
