package models

import (
	"html/template"
)

type HomePage struct {
	Profile   *Profile
	Entries   []Entry
	MorePosts bool
	LoggedIn  bool
}

type Profile struct {
	Name     string
	Email    string
	Password string
	TagLine  string
	GitHub   string
	LinkedIn string
	Twitter  string
	Facebook string
}

type Entry struct {
	ID    string
	Title string
	Body  template.HTML
}

type Login struct {
	Email    string
	Password string
}
