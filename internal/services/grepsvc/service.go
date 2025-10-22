package grepsvc

import (
	"bytes"
	"context"
	"fmt"
	"regexp"

	"github.com/sunr3d/quorum-grep/internal/interfaces/services"
	"github.com/sunr3d/quorum-grep/models"
)

var _ services.GrepService = (*grepService)(nil)

type grepService struct{}

// New - конструктор grepService.
func New() services.GrepService {
	return &grepService{}
}

// ProcessChunk - метод для обработки кусочка данных.
func (s *grepService) ProcessChunk(_ context.Context, task *models.Task) (*models.Result, error) {
	lines := bytes.Split(task.Data, []byte("\n"))
	lineLen := len(lines)
	if lineLen > 0 && len(lines[lineLen-1]) == 0 {
		lines = lines[:lineLen-1]
	}

	pattern, err := s.makePattern(task.Options)
	if err != nil {
		return nil, fmt.Errorf("makePattern: %w", err)
	}

	matches := s.findMatches(lines, pattern, task)

	return &models.Result{
		Matches:    matches,
		MatchCount: len(matches),
	}, nil
}

// Хелперы

// findMatches - поиск совпадений в строках.
func (s *grepService) findMatches(lines [][]byte, pattern *regexp.Regexp, task *models.Task) []models.Match {
	matches := make([]models.Match, 0, len(lines)*2)
	collected := make(map[int]struct{}, len(lines))

	for i, line := range lines {
		if s.matchLine(pattern, line, task.Options) {
			contextBefore, contextAfter := s.getContext(lines, i, task, collected)

			for _, ctxLine := range contextBefore {
				matches = append(matches, models.Match{
					Content:    ctxLine.Content,
					LineNumber: ctxLine.LineNumber,
				})
			}

			if _, ok := collected[i]; !ok {
				collected[i] = struct{}{}
				matches = append(matches, models.Match{
					Content:    line,
					LineNumber: task.LineNumbers[i],
				})
			}

			for _, ctxLine := range contextAfter {
				matches = append(matches, models.Match{
					Content:    ctxLine.Content,
					LineNumber: ctxLine.LineNumber,
				})
			}
		}
	}

	return matches
}

// getContext - получение контекста для совпадения.
func (s *grepService) getContext(
	lines [][]byte,
	pivot int,
	task *models.Task,
	collected map[int]struct{},
) ([]models.Match, []models.Match) {
	start, end := s.getContextRange(pivot, len(lines), task.Options)

	contextBefore := make([]models.Match, 0, end-start+1)
	contextAfter := make([]models.Match, 0, end-start+1)

	for i := start; i < pivot; i++ {
		if _, ok := collected[i]; !ok {
			collected[i] = struct{}{}
			contextBefore = append(contextBefore, models.Match{
				Content:    lines[i],
				LineNumber: task.LineNumbers[i],
			})
		}
	}

	for i := pivot + 1; i <= end; i++ {
		if _, ok := collected[i]; !ok {
			collected[i] = struct{}{}
			contextAfter = append(contextAfter, models.Match{
				Content:    lines[i],
				LineNumber: task.LineNumbers[i],
			})
		}
	}

	return contextBefore, contextAfter
}

// makePattern - создание регулярного выражения для поиска из паттерна и опций.
func (s *grepService) makePattern(opts models.GrepOptions) (*regexp.Regexp, error) {
	pattern := opts.Pattern

	if opts.Fixed {
		pattern = regexp.QuoteMeta(pattern)
	}

	if opts.IgnoreCase {
		pattern = "(?i)" + pattern
	}

	return regexp.Compile(pattern)
}

// matchLine - проверка совпадения строки с регулярным выражением.
func (s *grepService) matchLine(pattern *regexp.Regexp, line []byte, opts models.GrepOptions) bool {
	match := pattern.Match(line)
	return match != opts.Invert
}

// getContextRange - получение диапазона строк для контекста.
func (s *grepService) getContextRange(pivot int, linesLen int, opts models.GrepOptions) (int, int) {
	before := opts.Before
	after := opts.After

	if opts.Around > 0 {
		before = opts.Around
		after = opts.Around
	}

	start := pivot - before
	if start < 0 {
		start = 0
	}

	end := pivot + after
	if end >= linesLen {
		end = linesLen - 1
	}

	return start, end
}
