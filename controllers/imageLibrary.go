package controllers

import (
	"io"
	"net/http"
	"os"
)

func ImageLibraryHandler(w http.ResponseWriter, r *http.Request) {

	renderTemplate(w, "imagelibrary", nil)
}

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m := r.MultipartForm
	files := m.File["photo"]
	if len(files) > 0 {
		picFile, err := files[0].Open()
		defer picFile.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//create destination file making sure the path is writeable.
		dst, err := os.Create("./web/images/library/" + files[0].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, picFile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("/images/library/" + files[0].Filename))
	}

}
