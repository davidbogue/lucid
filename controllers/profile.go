package controllers

import (
	"encoding/json"
	"github.com/davidbogue/lucid/models"
	"github.com/gorilla/schema"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func UpdateProfilePicHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m := r.MultipartForm
	files := m.File["profilepic"]
	if len(files) > 0 {
		picFile, err := files[0].Open()
		defer picFile.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var extension = filepath.Ext(files[0].Filename)
		if strings.ToLower(extension) == ".jpeg" {
			extension = ".jpg"
		}

		//create destination file making sure the path is writeable.
		dst, err := os.Create("./web/images/profile" + extension)
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
		picFile.Close()
		err = resizeImage(extension)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// delete any old profile pics
		if extension == ".jpg" {
			os.Remove("./web/images/profile.png")
		}
	}
	http.Redirect(w, r, "/editprofile", http.StatusFound)
}

func resizeImage(ext string) error {
	file, err := os.Open("./web/images/profile" + ext)
	if err != nil {
		return err
	}
	var img image.Image
	if strings.ToLower(ext) == ".jpg" {
		img, err = jpeg.Decode(file)
		if err != nil {
			file.Close()
			return err
		}
	}
	if strings.ToLower(ext) == ".png" {
		img, err = png.Decode(file)
		if err != nil {
			file.Close()
			return err
		}
	}

	file.Close()

	// resize to width 100 using Lanczos resampling
	m := resize.Resize(100, 100, img, resize.Lanczos3)

	out, err := os.Create("./web/images/profile" + ext)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
	return nil
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

	setProfilePic(profile)

	return profile, nil
}

func setProfilePic(profile *models.Profile) {
	if _, err := os.Stat("./web/images/profile.png"); err == nil {
		profile.Picture = "profile.png"
		return
	}
	if _, err := os.Stat("./web/images/profile.jpg"); err == nil {
		profile.Picture = "profile.jpg"
		return
	}
	profile.Picture = "default.png"

}
