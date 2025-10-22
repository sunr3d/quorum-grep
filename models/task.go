package models

type Task struct {
	Data        []byte
	Index       int
	LineNumbers []int64
	Options     GrepOptions
}

type Result struct {
	Matches    []Match
	MatchCount int
	Error      string
	TaskIndex  int
}

type Match struct {
	Content    []byte
	LineNumber int64
}
