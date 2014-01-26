package controllers

import (
	"github.com/davidbogue/lucid/config"
	"html/template"
	"net/http"
)

var sessionName = config.Value("session_name")

var templates = template.Must(template.ParseFiles(
	"./web/header.html",
	"./web/footer.html",
	"./web/message.html",
	"./web/index.html",
	"./web/login.html",
	"./web/entry.html",
	"./web/editentry.html",
	"./web/resetpassword.html",
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
