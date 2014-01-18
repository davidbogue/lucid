package controllers

import (
	"encoding/json"
	"github.com/davidbogue/lucid/models"
	"github.com/gorilla/schema"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func EntryHandler(w http.ResponseWriter, r *http.Request) {
	entryId := r.URL.Path[len("/entry/"):]
	e, err := loadEntry(entryId)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		p, _ := loadProfile()
		entryPage := &models.EntryPage{Profile: p,
			Entry:       e,
			LoggedIn:    isLoggedIn(r),
			NextEntryId: nextEntryId(entryId)}
		renderTemplate(w, "entry", entryPage)
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
	p, _ := loadProfile()
	entryPage := &models.EntryPage{Profile: p, Entry: e}
	renderTemplate(w, "editentry", entryPage)

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

func summarizeEntry(body string) template.HTML {
	if strings.Contains(body, "</p>") {
		endRange := strings.Index(body, "</p>") + 4
		return template.HTML(body[0:endRange])
	} else {
		return template.HTML(body[0:500])
	}

}

func loadEntries(page int) ([]*models.Entry, bool) {
	files, _ := ioutil.ReadDir("./data/entries/")
	files = reverseFiles(files)

	endRange := page * 4
	if endRange > len(files) {
		endRange = len(files)

	}
	morePages := endRange < len(files)
	startRange := endRange - 4
	if startRange < 0 {
		startRange = 0
	}

	filePage := files[startRange:endRange]

	entries := make([]*models.Entry, len(filePage))

	for i, f := range filePage {
		entryId := f.Name()[0 : len(f.Name())-5]
		entry, err := loadEntry(entryId)
		entry.Body = summarizeEntry(string(entry.Body))
		if err == nil {
			entries[i] = entry
		}
	}

	return entries, morePages
}

func nextEntryId(entryId string) string {
	files, _ := ioutil.ReadDir("./data/entries/")
	files = reverseFiles(files)
	for i, f := range files {
		fileEntryId := f.Name()[0 : len(f.Name())-5]
		if fileEntryId == entryId {
			if len(files) > (i + 1) {
				f2 := files[i+1]
				return f2.Name()[0 : len(f2.Name())-5]
			}
		}
	}
	return ""
}

func reverseFiles(files []os.FileInfo) []os.FileInfo {
	length := len(files)
	reverseFiles := make([]os.FileInfo, length)
	for i, f := range files {
		reverseFiles[length-(i+1)] = f
	}
	return reverseFiles
}
