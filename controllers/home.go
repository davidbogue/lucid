package controllers

import (
	"github.com/davidbogue/lucid/models"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

var SessionStore = sessions.NewCookieStore([]byte("d8e2f09c-6e37-44a8-a3ec-7a5608b54383"))

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadProfile()
	if err != nil {
		http.Redirect(w, r, "/editprofile/", http.StatusFound)
		return
	}

	pageNumber := getPageNumber(r)
	entries, morePages := loadEntries(pageNumber)
	homePage := &models.HomePage{Profile: p,
		Entries:   entries,
		MorePosts: morePages,
		NextPage:  pageNumber + 1,
		LoggedIn:  isLoggedIn(r)}
	renderTemplate(w, "index", homePage)
}

func getPageNumber(r *http.Request) int {
	page := r.URL.Path[len("/"):]
	// string to int
	i, err := strconv.Atoi(page)
	if err != nil {
		return 1
	}
	return i
}
