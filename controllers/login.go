package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/davidbogue/lucid/models"
	"github.com/gorilla/schema"
	"net/http"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadProfile()
	if err != nil {
		http.Redirect(w, r, "/editprofile/", http.StatusFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	login := new(models.Login)

	decoder := schema.NewDecoder()
	err = decoder.Decode(login, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	passwordHash := getMD5HashWithSalt(login.Password)

	if login.Email == p.Email && passwordHash == p.Password {
		session, _ := SessionStore.Get(r, sessionName)
		session.Values["logged-in"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		renderTemplate(w, "login", "Wha wha wha ... try again.")
	}

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := SessionStore.Get(r, sessionName)
	session.Values["logged-in"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func getMD5HashWithSalt(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text + "lucid"))
	return hex.EncodeToString(hasher.Sum(nil))
}
