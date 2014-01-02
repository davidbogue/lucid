package controllers

import (
	"net/http"
)

func EntryHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}
