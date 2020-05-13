package internal

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var params []string

// IndexHandler handles index request
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	categories := getCategories(w)

	var indexData IndexData

	isLoggedIn, user := checkCookie(w, r)

	posts := getPosts(w)
	posts = formatPosts(w, posts)

	indexData.Categories = categories
	indexData.IndexUser = user
	indexData.LoggedIn = isLoggedIn
	indexData.Posts = posts

	t, err := template.New("index.html").ParseFiles("templates/index.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, indexData)
	checkInternalServerError(err, w)
}

// LoginHandler handles login request
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		isLoggedIn, _ := checkCookie(w, r)
		if isLoggedIn {
			http.Redirect(w, r, "/profile", 301)
		}

		t, _ := template.New("login.html").ParseFiles("templates/login.html")
		errText := ""
		if errLogin {
			errText = "Sorry, incorrect email or password!"
		}
		err = t.Execute(w, errText)
		checkInternalServerError(err, w)
		errLogin = false
		return
	}
	// grab user info from the submitted form
	email := r.FormValue("email")
	password := r.FormValue("psw")
	// query database to get match username
	var user User
	err := db.QueryRow("SELECT email, password FROM users WHERE email=?",
		email).Scan(&user.Email, &user.Password)

	if err != nil {
		errLogin = true
		http.Redirect(w, r, "/login", 301)
	}

	// validate password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		errLogin = true
		http.Redirect(w, r, "/login", 301)
	}

	createCookie(w, email)

	// authenticated = true
	http.Redirect(w, r, "/", 301)

}

// RegisterHandler handles register request
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		isLoggedIn, _ := checkCookie(w, r)
		if isLoggedIn {
			http.Redirect(w, r, "/profile", 301)
		}

		t, _ := template.New("register.html").ParseFiles("templates/register.html")
		errText := ""
		if errRegister {
			errText = "Sorry, email or username are already taken!"
		}
		errRegister = false
		t.Execute(w, errText)
		return
	}

	// grab user info
	email := r.FormValue("email")
	password := r.FormValue("password")
	username := r.FormValue("username")
	avatar := r.FormValue("avatar")

	// Check existence of user
	var user User
	err1 := db.QueryRow("SELECT email, password, username FROM users WHERE email=?",
		email).Scan(&user.Email, &user.Username, &user.Password)
	err2 := db.QueryRow("SELECT email, password, username FROM users WHERE username=?",
		username).Scan(&user.Email, &user.Username, &user.Password)
	// user is available
	if err1 == sql.ErrNoRows && err2 == sql.ErrNoRows {

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		checkInternalServerError(err, w)
		// insert to database
		_, err = db.Exec(`INSERT INTO users(email, password, username, avatar) VALUES(?, ?, ?, ?)`,
			email, hashedPassword, username, avatar)
		checkInternalServerError(err, w)

		createCookie(w, email)

		http.Redirect(w, r, "/", 301)

	} else {
		// checkInternalServerError(err1, w)
		// checkInternalServerError(err2, w)

		errRegister = true
		http.Redirect(w, r, "/register", 301)
	}

}

// LogoutHandler handles logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	deleteCookie(w)
	http.Redirect(w, r, "/", 301)
}

// ProfileHandler handles account info
func ProfileHandler(w http.ResponseWriter, r *http.Request) {

	isLoggedIn, user := checkCookie(w, r)

	var profileData ProfileData
	profileData.ProfileUser = user

	if isLoggedIn {

		if user.Avatar == 1 {
			profileData.Avatar1 = true
		} else if user.Avatar == 2 {
			profileData.Avatar2 = true
		} else {
			profileData.Avatar3 = true
		}

		posts := getPostsOfUser(w, int(user.ID))
		posts = formatPosts(w, posts)

		profileData.Posts = posts

		t, err := template.New("profile.html").ParseFiles("templates/profile.html")
		checkInternalServerError(err, w)
		err = t.Execute(w, profileData)
		checkInternalServerError(err, w)
	} else {
		http.Redirect(w, r, "/login", 301)
	}

}

// UserHandler handles public profile of user
func UserHandler(w http.ResponseWriter, r *http.Request) {
	parameters := strings.Split(r.URL.Path, "/")
	param := ""
	if len(parameters) == 3 && parameters[2] != "" {
		param = parameters[2]
	} else {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	userID, err := strconv.Atoi(param)

	if err != nil {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	var selectedUser User
	err = db.QueryRow("SELECT email, username, avatar FROM users WHERE id=?",
		userID).Scan(&selectedUser.Email, &selectedUser.Username, &selectedUser.Avatar)

	if err != nil {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	var profileData ProfileData
	profileData.ProfileUser = selectedUser

	if selectedUser.Avatar == 1 {
		profileData.Avatar1 = true
	} else if selectedUser.Avatar == 2 {
		profileData.Avatar2 = true
	} else {
		profileData.Avatar3 = true
	}

	isLoggedIn, user := checkCookie(w, r)

	posts := getPostsOfUser(w, userID)
	posts = formatPosts(w, posts)

	postsLiked := getLikedPostsOfUser(w, userID)
	postsLiked = formatPosts(w, postsLiked)

	var userData UserData
	userData.LoggedIn = isLoggedIn
	userData.ProfData = profileData
	userData.ProfileUser = user
	userData.ProfPosts = posts
	userData.LikedPosts = postsLiked

	t, err := template.New("user.html").ParseFiles("templates/user.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, userData)
	checkInternalServerError(err, w)

}

// AddPostHandler handles new post addition
func AddPostHandler(w http.ResponseWriter, r *http.Request) {
	isLoggedIn, user := checkCookie(w, r)
	if !isLoggedIn {
		http.Redirect(w, r, "/login", 301)
	}

	if r.Method != "POST" {

		categories := getCategories(w)

		var templateData IndexData
		templateData.Categories = categories
		templateData.IndexUser = user
		templateData.LoggedIn = isLoggedIn

		t, err := template.New("add_post.html").ParseFiles("templates/add_post.html")
		checkInternalServerError(err, w)
		err = t.Execute(w, templateData)
		checkInternalServerError(err, w)
		return

	}

	// // grab post info
	title := r.FormValue("title")
	content := r.FormValue("content")
	categoryID := r.FormValue("category")

	_, err = db.Exec(`INSERT INTO posts(title, content, author_id, category_id) VALUES(?, ?, ?, ?)`,
		title, content, user.ID, categoryID)

	checkInternalServerError(err, w)

	http.Redirect(w, r, "/", 301)

}

// PostHandler handles one post iwth given id
func PostHandler(w http.ResponseWriter, r *http.Request) {

	parameters := strings.Split(r.URL.Path, "/")
	postString := ""

	if len(parameters) == 3 && parameters[2] != "" {
		postString = parameters[2]
	} else {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	postID, err := strconv.Atoi(postString)

	if err != nil {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	var selectedPost Post
	err = db.QueryRow("SELECT id, title, content, timestamp, author_id, category_id FROM posts WHERE id=?",
		postID).Scan(&selectedPost.ID, &selectedPost.Title, &selectedPost.Content, &selectedPost.Timestamp, &selectedPost.Author, &selectedPost.Category)

	if err != nil {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	isLoggedIn, user := checkCookie(w, r)
	if r.Method != "POST" {

		likesCount := getLikes(w, postID)
		dislikesCount := getDislikes(w, postID)

		userLiked := false
		userDisliked := false

		if isLoggedIn {
			userLiked = getUserLike(w, postID, int(user.ID))
			userDisliked = getUserDislike(w, postID, int(user.ID))
		}

		post := formatPost(w, selectedPost)

		comments := getComments(w, postID)
		comments = formatComments(w, comments)

		var postData PostData

		postData.CurrPost = post
		postData.LoggedIn = isLoggedIn
		postData.UserData = user
		postData.Comments = comments
		postData.Likes = likesCount
		postData.Dislikes = dislikesCount
		postData.UserLiked = userLiked
		postData.UserDisliked = userDisliked

		t, err := template.New("post.html").ParseFiles("templates/post.html")
		checkInternalServerError(err, w)
		err = t.Execute(w, postData)
		checkInternalServerError(err, w)
		return
	}

	comment := r.FormValue("comment")

	if comment != "" {

		if isLoggedIn {
			_, err = db.Exec(`INSERT INTO comments(text, author_id, post_id) VALUES(?, ?, ?)`,
				comment, user.ID, postID)

			checkInternalServerError(err, w)
			http.Redirect(w, r, "/post/"+postString, 301)
		} else {
			http.Redirect(w, r, "/login", 301)
		}

	}

	likedislike := r.FormValue("likedislike")

	if isLoggedIn {
		if likedislike == "like" {
			addLike(w, postID, int(user.ID))
		} else if likedislike == "dislike" {
			addDislike(w, postID, int(user.ID))
		}
		http.Redirect(w, r, "/post/"+postString, 301)

	} else {
		http.Redirect(w, r, "/login", 301)
	}

}

// CategoryHandler handles posts of given category
func CategoryHandler(w http.ResponseWriter, r *http.Request) {
	parameters := strings.Split(r.URL.Path, "/")
	param := ""
	if len(parameters) == 3 && parameters[2] != "" {
		param = parameters[2]
	} else {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	categoryID, err := strconv.Atoi(param)

	if err != nil {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	posts, err := getPostsOfCategory(categoryID)
	category := getCategoryName(w, categoryID)

	if err != nil {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	posts = formatPosts(w, posts)
	categories := getCategories(w)
	isLoggedIn, user := checkCookie(w, r)

	var templateData IndexData

	templateData.Categories = categories
	templateData.Posts = posts
	templateData.IndexUser = user
	templateData.LoggedIn = isLoggedIn
	templateData.Category = category

	t, err := template.New("category.html").ParseFiles("templates/category.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, templateData)
	checkInternalServerError(err, w)
}
