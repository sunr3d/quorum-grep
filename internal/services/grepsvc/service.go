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
func (s *grepService) ProcessChunk(ctx context.Context, task *models.Task) (*models.Result, error) {
	lines := bytes.Split(task.Data, []byte("\n"))

	pattern, err := s.makePattern(task.Options)
	if err != nil {
		return nil, fmt.Errorf("makePattern: %w", err)
	}

	matches := s.findMatches(lines, pattern, task)

	return &models.Result{
		TaskID:     task.ID,
		Matches:    matches,
		MatchCount: len(matches),
		Error:      "",
	}, nil
}

// Хелперы

// findMatches - поиск совпадений в строках.
func (s *grepService) findMatches(lines [][]byte, pattern *regexp.Regexp, task *models.Task) []models.Match {
	matches := make([]models.Match, 0, len(lines))
	collected := make(map[int]bool, len(lines))

	for i, line := range lines {
		if s.matchLine(pattern, line, task.Options) {
			contextBefore, contextAfter := s.getContext(lines, i, task.Options, collected)

			matches = append(matches, models.Match{
				LineNum:       i + 1,
				Content:       line,
				ContextBefore: contextBefore,
				ContextAfter:  contextAfter,
			})
		}
	}

	return matches
}

// getContext - получение контекста для совпадения.
func (s *grepService) getContext(lines [][]byte, pivot int, opts models.GrepOptions, collected map[int]bool) ([][]byte, [][]byte) {
	start, end := s.getContextRange(pivot, len(lines), opts)

	contextBefore := make([][]byte, 0, end-start+1)
	contextAfter := make([][]byte, 0, end-start+1)

	for i := start; i < pivot; i++ {
		if !collected[i] {
			collected[i] = true
			contextBefore = append(contextBefore, lines[i])
		}
	}

	for i := pivot + 1; i <= end; i++ {
		if !collected[i] {
			collected[i] = true
			contextAfter = append(contextAfter, lines[i])
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
