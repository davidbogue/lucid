package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadProfile()
	if err != nil {
		http.Redirect(w, r, "/editprofile/", http.StatusFound)
		return
	}
	homePage := &HomePage{Profile: p, Entries: loadEntries()}
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

	profile := new(Profile)

	decoder := schema.NewDecoder()
	err = decoder.Decode(profile, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Printf("form value name is %v \n", r.FormValue("Name"))
	fmt.Printf("post form is %v \n", r.PostForm)

	// validate profile and return errors to edit screen here
	b, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Printf("profile is %+v \n", profile)

	err = ioutil.WriteFile("./data/profile.json", b, 0600)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func loadProfile() (*Profile, error) {
	filename := "./data/profile.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	profile := new(Profile)

	err = json.Unmarshal(data, profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func loadEntries() []Entry {
	entries := make([]Entry, 4)
	testBody := "test \n##h2 \n markdown test"
	output := blackfriday.MarkdownBasic([]byte(testBody))
	fmt.Printf(string(output))

	entries[0] = Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	entries[1] = Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	entries[2] = Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	entries[3] = Entry{ID: "10", Title: "This is a blog post", Body: template.HTML(output)}
	return entries
}
