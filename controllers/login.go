package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/davidbogue/lucid/config"
	"github.com/davidbogue/lucid/models"
	"github.com/gorilla/schema"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadProfile()
	if err != nil {
		http.Redirect(w, r, "/editprofile/", http.StatusFound)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	login := new(models.Login)

	decoder := schema.NewDecoder()
	err = decoder.Decode(login, r.PostForm)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	passwordHash := getMD5HashWithSalt(login.Password)

	if login.Email == p.Email && passwordHash == p.Password {
		session, _ := SessionStore.Get(r, sessionName)
		session.Values["logged-in"] = true
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		renderTemplate(w, "login", "Wrong answer... try again.")
	}

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := SessionStore.Get(r, sessionName)
	session.Values["logged-in"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "resetpassword", nil)
}

func ResetPasswordFormHandler(w http.ResponseWriter, r *http.Request) {
	profile, err := loadProfile()
	if err != nil {
		http.Redirect(w, r, "/editprofile/", http.StatusFound)
		return
	}

	err = r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	email := r.FormValue("Email")
	if email == profile.Email {
		err = sendResetPasswordEmail(profile)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		renderTemplate(w, "resetpassword", "That's not correct.")
	}

	messagePage := &models.MessagePage{Profile: profile,
		Message: "Check your email for your new password."}
	renderTemplate(w, "message", messagePage)
}

func sendResetPasswordEmail(profile *models.Profile) error {
	newPassword := generateRandomPassword()
	profile.Password = getMD5HashWithSalt(newPassword)
	err := writeProfile(profile)
	if err != nil {
		return err
	}

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		config.Value("email_from"),
		config.Value("email_pass"),
		config.Value("email_smtp"),
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	sub := "Subject: Lucid Blog Password Reset\r\n\r\n"
	body := "Your new password is " + newPassword

	err = smtp.SendMail(
		config.Value("email_smtp")+":"+config.Value("email_port"),
		auth,
		config.Value("email_from"),
		[]string{profile.Email},
		[]byte(sub+body),
	)

	return err
}

func generateRandomPassword() string {
	seed := time.Now().Unix()
	rand.Seed(seed)
	randomNumber := strconv.Itoa(rand.Intn(100000))
	hasher := md5.New()
	hasher.Write([]byte(randomNumber))
	return hex.EncodeToString(hasher.Sum(nil))[0:8]
}

func getMD5HashWithSalt(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text + config.Value("password_salt")))
	return hex.EncodeToString(hasher.Sum(nil))
}
