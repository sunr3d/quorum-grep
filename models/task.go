package models

type Task struct {
	ID      string
	Data    []byte
	Index   int
	Options GrepOptions
}

type Result struct {
	TaskID     string
	Matches    []Match
	MatchCount int
	Error      string
}

type Match struct {
	LineNum       int
	Content       []byte
	ContextBefore [][]byte
	ContextAfter  [][]byte
}
