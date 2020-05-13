package main

import (
	"net/http"
	"os"

	internal "forum/internal"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	internal.InitDb()

	//os.Setenv("PORT", "8898")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// routes
	http.HandleFunc("/", internal.IndexHandler)
	http.HandleFunc("/login", internal.LoginHandler)
	http.HandleFunc("/register", internal.RegisterHandler)
	http.HandleFunc("/logout", internal.LogoutHandler)
	http.HandleFunc("/profile", internal.ProfileHandler)
	http.HandleFunc("/user/", internal.UserHandler)
	http.HandleFunc("/add-post", internal.AddPostHandler)
	// http.HandleFunc("/posts", internal.PostsHandler)
	http.HandleFunc("/post/", internal.PostHandler)
	http.HandleFunc("/category/", internal.CategoryHandler)

	http.Handle("/statics/",
		http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))),
	)
	http.ListenAndServe(":"+port, nil)
}
