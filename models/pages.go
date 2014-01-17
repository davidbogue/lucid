package models

type HomePage struct {
	Profile   *Profile
	Entries   []*Entry
	NextPage  int
	MorePosts bool
	LoggedIn  bool
}

type EntryPage struct {
	Profile  *Profile
	Entry    *Entry
	LoggedIn bool
}

type Login struct {
	Email    string
	Password string
}
