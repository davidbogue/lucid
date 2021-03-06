package models

type HomePage struct {
	Profile   *Profile
	Entries   []*Entry
	NextPage  int
	MorePosts bool
	LoggedIn  bool
}

type EntryPage struct {
	Profile     *Profile
	Entry       *Entry
	NextEntryId string
	LoggedIn    bool
}

type MessagePage struct {
	Profile *Profile
	Message string
}

type Login struct {
	Email    string
	Password string
}
