package controllers

import (
	"encoding/json"
	"github.com/davidbogue/lucid/models"
	"github.com/gorilla/schema"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"
)

func EntryHandler(w http.ResponseWriter, r *http.Request) {
	entryId := r.URL.Path[len("/entry/"):]
	e, err := loadEntry(entryId)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		renderTemplate(w, "entry", e)
	}
}

func EditEntryHandler(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	e := new(models.Entry)
	entryId := r.URL.Path[len("/editentry/"):]
	if entryId != "" {
		var err error
		e, err = loadEntry(entryId)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
	renderTemplate(w, "editentry", e)

}

func SaveEntryHandler(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	entry := new(models.Entry)
	decoder := schema.NewDecoder()
	err = decoder.Decode(entry, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if entry.ID == "" {
		entry.ID = getNextEntryId()
	}

	entry.Body = template.HTML(blackfriday.MarkdownBasic([]byte(entry.Markdown)))
	entryJson, err := json.Marshal(entry)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ioutil.WriteFile("./data/entries/"+entry.ID+".json", entryJson, 0600)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/entry/"+entry.ID, http.StatusFound)

}

func getNextEntryId() string {
	t := time.Now()
	return t.Format("20060102150405")
}

func loadEntry(id string) (*models.Entry, error) {
	filename := "./data/entries/" + id + ".json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	entry := new(models.Entry)

	err = json.Unmarshal(data, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func loadEntries(page int) []*models.Entry {
	files, _ := ioutil.ReadDir("./data/entries/")

	endRange := page * 4
	if endRange > len(files) {
		endRange = len(files)
	}
	startRange := endRange - 4
	if startRange < 0 {
		startRange = 0
	}

	filePage := files[startRange:endRange]

	entries := make([]*models.Entry, len(filePage))

	// testBody := "When we look at social software services like Facebook and Twitter, we are really talking about systems who’s whole purpose is to get us to form narratives through them. We form these narratives by stringing together syntagms (fragments of text) into sequential interwoven dialogs that together form stories/narratives amongst others.\n\nThe feeds and reverse chronological way these are presented is the easiest distillation of what we contribute. We are story tellers through the medium..."
	// output := blackfriday.MarkdownBasic([]byte(testBody))

	for i, f := range filePage {
		entryId := f.Name()[0 : len(f.Name())-5]
		entry, err := loadEntry(entryId)
		if err == nil {
			entries[i] = entry
		}
	}

	return entries
}
