package internal

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// IndexHandler handles index request
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.ServeFile(w, r, "templates/error.html")
		return
	}

	isLoggedIn, user := checkCookie(w, r)
	t, err := template.New("index.html").ParseFiles("templates/index.html")
	checkInternalServerError(err, w)
	if isLoggedIn {
		err = t.Execute(w, user)
	} else {
		err = t.Execute(w, nil)
	}
	checkInternalServerError(err, w)

	// rows, err := db.Query("SELECT * FROM users")
	// checkInternalServerError(err, w)
	// var users []User
	// var user User
	// for rows.Next() {
	// 	err = rows.Scan(&user.ID, &user.Email,
	// 		&user.Username, &user.Password, &user.SessionToken)
	// 	checkInternalServerError(err, w)
	// 	users = append(users, user)
	// }
	// t, err := template.New("index.html").ParseFiles("templates/index.html")
	// checkInternalServerError(err, w)
	// err = t.Execute(w, users)
	// checkInternalServerError(err, w)
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

	// // validate password
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
		checkInternalServerError(err1, w)
		checkInternalServerError(err2, w)

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

	var userData UserData
	userData.LoggedIn = isLoggedIn
	userData.ProfData = profileData
	userData.ProfileUser = user

	t, err := template.New("user.html").ParseFiles("templates/user.html")
	checkInternalServerError(err, w)
	err = t.Execute(w, userData)
	checkInternalServerError(err, w)

}
