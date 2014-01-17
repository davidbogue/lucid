package models

import (
	"html/template"
)

type Entry struct {
	ID       string
	Title    string
	Markdown string
	Body     template.HTML
}
