package main

import (
	"github.com/davidbogue/lucid/controllers"
	"net/http"
)

func main() {
	http.HandleFunc("/", controllers.HomePageHandler)
	http.HandleFunc("/editprofile/", controllers.EditProfileHandler)
	http.HandleFunc("/saveprofile/", controllers.SaveProfileHandler)
	http.HandleFunc("/login/", controllers.LoginPageHandler)
	http.HandleFunc("/loginform/", controllers.LoginFormHandler)
	http.HandleFunc("/logout/", controllers.LogoutHandler)
	http.HandleFunc("/entry/", controllers.EntryHandler)
	http.HandleFunc("/editentry/", controllers.EditEntryHandler)
	http.HandleFunc("/saveentry/", controllers.SaveEntryHandler)
	http.HandleFunc("/updateprofilepic/", controllers.UpdateProfilePicHandler)

	// static files
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("./web/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./web/images/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./web/fonts/"))))

	http.ListenAndServe(":8080", nil)
}
