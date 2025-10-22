package grepsvc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sunr3d/quorum-grep/models"
)

// Тест для основного метода сервиса.
func TestGrepService_ProcessChunk(t *testing.T) {
	svc := New()

	tests := []struct {
		name     string
		task     *models.Task
		expected *models.Result
		wantErr  bool
	}{
		{
			name: "базовый поиск",
			task: &models.Task{
				Data:        []byte("line1\npattern found\nline3"),
				Index:       0,
				LineNumbers: []int64{1, 2, 3},
				Options: models.GrepOptions{
					Pattern: "pattern",
				},
			},
			expected: &models.Result{
				Matches:    []models.Match{{Content: []byte("pattern found"), LineNumber: 2}},
				MatchCount: 1,
				Error:      "",
				TaskIndex:  0,
			},
			wantErr: false,
		},
		{
			name: "поиск с контекстом -A 1",
			task: &models.Task{
				Data:        []byte("line1\npattern found\nline3\nline4"),
				Index:       0,
				LineNumbers: []int64{1, 2, 3, 4},
				Options: models.GrepOptions{
					Pattern: "pattern",
					After:   1,
				},
			},
			expected: &models.Result{
				Matches:    []models.Match{{Content: []byte("pattern found"), LineNumber: 2}, {Content: []byte("line3"), LineNumber: 3}},
				MatchCount: 2,
				Error:      "",
				TaskIndex:  0,
			},
			wantErr: false,
		},
		{
			name: "поиск с контекстом -B 1",
			task: &models.Task{
				Data:        []byte("line1\npattern found\nline3"),
				Index:       0,
				LineNumbers: []int64{1, 2, 3},
				Options: models.GrepOptions{
					Pattern: "pattern",
					Before:  1,
				},
			},
			expected: &models.Result{
				Matches:    []models.Match{{Content: []byte("line1"), LineNumber: 1}, {Content: []byte("pattern found"), LineNumber: 2}},
				MatchCount: 2,
				Error:      "",
				TaskIndex:  0,
			},
			wantErr: false,
		},
		{
			name: "игнорирование регистра -i",
			task: &models.Task{
				Data:        []byte("line1\nPATTERN found\nline3"),
				Index:       0,
				LineNumbers: []int64{1, 2, 3},
				Options: models.GrepOptions{
					Pattern:    "pattern",
					IgnoreCase: true,
				},
			},
			expected: &models.Result{
				Matches:    []models.Match{{Content: []byte("PATTERN found"), LineNumber: 2}},
				MatchCount: 1,
				Error:      "",
				TaskIndex:  0,
			},
			wantErr: false,
		},
		{
			name: "инвертированный поиск -v",
			task: &models.Task{
				Data:        []byte("line1\npattern found\nline3"),
				Index:       0,
				LineNumbers: []int64{1, 2, 3},
				Options: models.GrepOptions{
					Pattern: "pattern",
					Invert:  true,
				},
			},
			expected: &models.Result{
				Matches:    []models.Match{{Content: []byte("line1"), LineNumber: 1}, {Content: []byte("line3"), LineNumber: 3}},
				MatchCount: 2,
				Error:      "",
				TaskIndex:  0,
			},
			wantErr: false,
		},
		{
			name: "фиксированная строка -F",
			task: &models.Task{
				Data:        []byte("line1\npattern.found\nline3"),
				Index:       0,
				LineNumbers: []int64{1, 2, 3},
				Options: models.GrepOptions{
					Pattern: "pattern.found",
					Fixed:   true,
				},
			},
			expected: &models.Result{
				Matches:    []models.Match{{Content: []byte("pattern.found"), LineNumber: 2}},
				MatchCount: 1,
				Error:      "",
				TaskIndex:  0,
			},
			wantErr: false,
		},
		{
			name: "пустые данные",
			task: &models.Task{
				Data:        []byte(""),
				Index:       0,
				LineNumbers: []int64{},
				Options: models.GrepOptions{
					Pattern: "pattern",
				},
			},
			expected: &models.Result{
				Matches:    []models.Match{},
				MatchCount: 0,
				Error:      "",
				TaskIndex:  0,
			},
			wantErr: false,
		},
		{
			name: "невалидный regex",
			task: &models.Task{
				Data:        []byte("line1\nline2"),
				Index:       0,
				LineNumbers: []int64{1, 2},
				Options: models.GrepOptions{
					Pattern: "[invalid",
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := svc.ProcessChunk(context.Background(), tt.task)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.MatchCount, result.MatchCount)
			assert.Equal(t, tt.expected.Error, result.Error)
			assert.Equal(t, len(tt.expected.Matches), len(result.Matches))

			for i, expectedMatch := range tt.expected.Matches {
				if i < len(result.Matches) {
					assert.Equal(t, string(expectedMatch.Content), string(result.Matches[i].Content))
					assert.Equal(t, expectedMatch.LineNumber, result.Matches[i].LineNumber)
				}
			}
		})
	}
}

// Тесты для хелперов.
func TestGrepService_makePattern(t *testing.T) {
	svc := &grepService{}

	tests := []struct {
		name     string
		opts     models.GrepOptions
		expected string
		wantErr  bool
	}{
		{
			name: "обычный паттерн",
			opts: models.GrepOptions{
				Pattern: "test",
			},
			expected: "test",
			wantErr:  false,
		},
		{
			name: "игнорирование регистра",
			opts: models.GrepOptions{
				Pattern:    "test",
				IgnoreCase: true,
			},
			expected: "(?i)test",
			wantErr:  false,
		},
		{
			name: "фиксированная строка",
			opts: models.GrepOptions{
				Pattern: "test.pattern",
				Fixed:   true,
			},
			expected: "test\\.pattern",
			wantErr:  false,
		},
		{
			name: "игнорирование регистра + фиксированная строка",
			opts: models.GrepOptions{
				Pattern:    "test.pattern",
				IgnoreCase: true,
				Fixed:      true,
			},
			expected: "(?i)test\\.pattern",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, err := svc.makePattern(tt.opts)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, pattern.String())
		})
	}
}

func TestGrepService_getContextRange(t *testing.T) {
	svc := &grepService{}

	tests := []struct {
		name     string
		pivot    int
		linesLen int
		opts     models.GrepOptions
		expected [2]int
	}{
		{
			name:     "обычный контекст",
			pivot:    2,
			linesLen: 5,
			opts: models.GrepOptions{
				Before: 1,
				After:  1,
			},
			expected: [2]int{1, 3},
		},
		{
			name:     "контекст в начале",
			pivot:    0,
			linesLen: 5,
			opts: models.GrepOptions{
				Before: 2,
				After:  1,
			},
			expected: [2]int{0, 1},
		},
		{
			name:     "контекст в конце",
			pivot:    4,
			linesLen: 5,
			opts: models.GrepOptions{
				Before: 1,
				After:  2,
			},
			expected: [2]int{3, 4},
		},
		{
			name:     "флаг -C (around)",
			pivot:    2,
			linesLen: 5,
			opts: models.GrepOptions{
				Around: 1,
			},
			expected: [2]int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := svc.getContextRange(tt.pivot, tt.linesLen, tt.opts)
			assert.Equal(t, tt.expected[0], start)
			assert.Equal(t, tt.expected[1], end)
		})
	}
}
