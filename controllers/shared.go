package controllers

import (
	"html/template"
	"net/http"
)

var sessionName = "lucid-session"

var templates = template.Must(template.ParseFiles(
	"./web/header.html",
	"./web/footer.html",
	"./web/index.html",
	"./web/login.html",
	"./web/profile.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, "header.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = templates.ExecuteTemplate(w, "footer.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func isLoggedIn(r *http.Request) bool {
	session, _ := SessionStore.Get(r, sessionName)
	if session.Values["logged-in"] != nil {
		return session.Values["logged-in"].(bool)
	}
	return false
}
