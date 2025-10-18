package models

type GrepOptions struct {
	Pattern    string
	After      int
	Before     int
	Around     int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
}
