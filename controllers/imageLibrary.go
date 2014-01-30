package controllers

import (
	"net/http"
)

func ImageLibraryHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "imagelibrary", nil)
}
