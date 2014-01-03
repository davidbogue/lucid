package controllers

import (
	"net/http"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}
