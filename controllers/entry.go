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
	testBody := "When we look at social software services like Facebook and Twitter, we are really talking about systems who’s whole purpose is to get us to form narratives through them. We form these narratives by stringing together syntagms (fragments of text) into sequential interwoven dialogs that together form stories/narratives amongst others.\n\nThe feeds and reverse chronological way these are presented is the easiest distillation of what we contribute. We are story tellers through the medium..."
	output := blackfriday.MarkdownBasic([]byte(testBody))

	entries[0] = models.Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	entries[1] = models.Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	entries[2] = models.Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	entries[3] = models.Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	return entries
}
