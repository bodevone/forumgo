package main

import (
	"log"
	"net/http"

	internal "forum/internal"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	internal.InitDb()

	//os.Setenv("PORT", "8898")
	// port := os.Getenv("PORT")
	port := "8080"
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	// route
	http.HandleFunc("/", internal.IndexHandler)
	http.HandleFunc("/login", internal.LoginHandler)
	http.HandleFunc("/register", internal.RegisterHandler)
	http.HandleFunc("/logout", internal.LogoutHandler)
	http.HandleFunc("/profile", internal.ProfileHandler)
	// http.HandleFunc("/users", internal.Handler)

	http.Handle("/statics/",
		http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))),
	)
	http.ListenAndServe(":"+port, nil)
}
