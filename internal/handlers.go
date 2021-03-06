package internal

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var params []string

// IndexHandler handles index request
func IndexHandler(w http.ResponseWriter, r *http.Request) {

	methods := []string{"GET"}
	checkAllowedMethods(methods, w, r)

	if r.URL.Path != "/" {
		pageNotFound(w, r)
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

	methods := []string{"GET", "POST"}
	checkAllowedMethods(methods, w, r)

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

	email := r.FormValue("email")
	password := r.FormValue("psw")

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

	http.Redirect(w, r, "/", 301)

}

// RegisterHandler handles register request
func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	methods := []string{"GET", "POST"}
	checkAllowedMethods(methods, w, r)

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
		errRegister = true
		http.Redirect(w, r, "/register", 301)
	}

}

// ProfileHandler handles account info
func ProfileHandler(w http.ResponseWriter, r *http.Request) {

	methods := []string{"GET", "POST"}
	checkAllowedMethods(methods, w, r)

	isLoggedIn, user := checkCookie(w, r)

	var profileData ProfileData
	profileData.ProfileUser = user

	if r.Method == "POST" {
		logout := r.FormValue("logout")

		if logout != "" {
			if isLoggedIn {
				deleteCookie(w, int(user.ID))
			} else {
				http.Redirect(w, r, "/login", 301)
			}
			http.Redirect(w, r, "/", 301)

		} else {
			errPost := errors.New("No POST request Data")
			checkInternalServerError(errPost, w)
		}
	}

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

	methods := []string{"GET"}
	checkAllowedMethods(methods, w, r)

	parameters := strings.Split(r.URL.Path, "/")
	param := ""
	if len(parameters) == 3 && parameters[2] != "" {
		param = parameters[2]
	} else {
		pageNotFound(w, r)
		return
	}

	userID, err := strconv.Atoi(param)

	if err != nil {
		pageNotFound(w, r)
		return
	}

	var selectedUser User
	err = db.QueryRow("SELECT email, username, avatar FROM users WHERE id=?",
		userID).Scan(&selectedUser.Email, &selectedUser.Username, &selectedUser.Avatar)

	if err != nil {
		pageNotFound(w, r)
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

	methods := []string{"GET", "POST"}
	checkAllowedMethods(methods, w, r)

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

	// grab post info
	title := r.FormValue("title")
	content := r.FormValue("content")

	categories := getCategories(w)
	categoryIDs := []string{}
	for _, category := range categories {
		id := strconv.Itoa(int(category.ID))
		if r.FormValue(id) != "" {
			categoryIDs = append(categoryIDs, id)
		}
	}

	ins, err := db.Exec(`INSERT INTO posts(title, content, author_id) VALUES(?, ?, ?)`,
		title, content, user.ID)
	checkInternalServerError(err, w)

	postID, _ := ins.LastInsertId()

	for _, categoryID := range categoryIDs {
		_, err = db.Exec(`INSERT OR IGNORE INTO postcategories(post_id, category_id) VALUES(?, ?)`,
			postID, categoryID)
	}

	// categoryIDsFormString := ""
	// notFirst := false
	// for _, categoryID := range categoryIDs {
	// 	if notFirst {
	// 		categoryIDsFormString += ","
	// 	}
	// 	notFirst = true
	// 	categoryIDsFormString += categoryID
	// }

	http.Redirect(w, r, "/", 301)

}

// PostHandler handles one post iwth given id
func PostHandler(w http.ResponseWriter, r *http.Request) {

	methods := []string{"GET", "POST"}
	checkAllowedMethods(methods, w, r)

	parameters := strings.Split(r.URL.Path, "/")
	postString := ""

	if len(parameters) == 3 && parameters[2] != "" {
		postString = parameters[2]
	} else {
		pageNotFound(w, r)
		return
	}

	postID, err := strconv.Atoi(postString)

	if err != nil {
		pageNotFound(w, r)
		return
	}

	var selectedPost Post
	err = db.QueryRow("SELECT id, title, content, timestamp, author_id FROM posts WHERE id=?",
		postID).Scan(&selectedPost.ID, &selectedPost.Title, &selectedPost.Content, &selectedPost.Timestamp, &selectedPost.Author)

	if err != nil {
		pageNotFound(w, r)
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
		comments = getCommentsLikes(w, comments)
		comments = getUserCommentsLikes(w, comments, int(user.ID))
		comments = getUserCommentsDislikes(w, comments, int(user.ID))

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

	if likedislike != "" {
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

	likedislikecomment := r.FormValue("likedislikecomment")

	commentID, err := strconv.Atoi(r.FormValue("commentid"))

	if isLoggedIn {
		if likedislikecomment == "like" {
			addLikeToComment(w, commentID, int(user.ID))
		} else if likedislikecomment == "dislike" {
			addDislikeToComment(w, commentID, int(user.ID))
		}
		http.Redirect(w, r, "/post/"+postString, 301)
	} else {
		http.Redirect(w, r, "/login", 301)
	}
}

// CategoryHandler handles posts of given category
func CategoryHandler(w http.ResponseWriter, r *http.Request) {

	methods := []string{"GET"}
	checkAllowedMethods(methods, w, r)

	parameters := strings.Split(r.URL.Path, "/")
	param := ""
	if len(parameters) == 3 && parameters[2] != "" {
		param = parameters[2]
	} else {
		pageNotFound(w, r)
		return
	}

	categoryID, err := strconv.Atoi(param)

	if err != nil {
		pageNotFound(w, r)
		return
	}

	posts, err := getPostsOfCategory(w, categoryID)

	if err != nil {
		pageNotFound(w, r)
		return
	}

	category, err := getCategoryName(w, categoryID)

	if err != nil {
		pageNotFound(w, r)
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

	t, _ := template.New("category.html").ParseFiles("templates/category.html")
	t.Execute(w, templateData)
}
