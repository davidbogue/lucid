package controllers

import (
	"github.com/davidbogue/lucid/models"
	"github.com/russross/blackfriday"
	"html/template"
	"net/http"
)

func EntryHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func loadEntries() []models.Entry {
	entries := make([]models.Entry, 4)
	testBody := "test \n##h2 \n markdown test"
	output := blackfriday.MarkdownBasic([]byte(testBody))

	entries[0] = models.Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	entries[1] = models.Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	entries[2] = models.Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	entries[3] = models.Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	return entries
}
