package controllers

import (
	"encoding/json"
	"github.com/davidbogue/lucid/models"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"io/ioutil"
	"net/http"
)

var SessionStore = sessions.NewCookieStore([]byte("d8e2f09c-6e37-44a8-a3ec-7a5608b54383"))

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadProfile()
	if err != nil {
		http.Redirect(w, r, "/editprofile/", http.StatusFound)
		return
	}
	homePage := &models.HomePage{Profile: p, Entries: loadEntries(), LoggedIn: isLoggedIn(r)}
	renderTemplate(w, "index", homePage)
}

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	p, _ := loadProfile()
	renderTemplate(w, "profile", p)
}

func SaveProfileHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	profile := new(models.Profile)

	decoder := schema.NewDecoder()
	err = decoder.Decode(profile, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// validate profile and return errors to edit screen here
	profile.Password = getMD5HashWithSalt(profile.Password)
	b, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ioutil.WriteFile("./data/profile.json", b, 0600)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func loadProfile() (*models.Profile, error) {
	filename := "./data/profile.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	profile := new(models.Profile)

	err = json.Unmarshal(data, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
