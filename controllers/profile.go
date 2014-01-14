package controllers

import (
	"encoding/json"
	"github.com/davidbogue/lucid/models"
	"github.com/gorilla/schema"
	"io/ioutil"
	"net/http"
	"os"
)

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadProfile()
	if err == nil { // the profile exists so we need to make sure the user is logged in.
		if !isLoggedIn(r) { // if not logged in then redirect to home page
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	renderTemplate(w, "profile", p)
}

func SaveProfileHandler(w http.ResponseWriter, r *http.Request) {
	if profileExists() && !isLoggedIn(r) { // the profile exists so we need to make sure the user is logged in.
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

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

func profileExists() bool {
	_, err := os.Stat("./data/profile.json")
	return !os.IsNotExist(err)
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
