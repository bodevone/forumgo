package internal

import (
	"database/sql"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	db            *sql.DB
	err           error
	errLogin      = false
	errRegister   = false
	authenticated = false
)

// InitDb starts database
func InitDb() {
	db, err = sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}
	// defer db.Close()
	// test connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	createUsers, _ := db.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, email TEXT, username TEXT, password TEXT, avatar INTEGER, session TEXT)")
	createUsers.Exec()

	createCategories, _ := db.Prepare("CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY, name TEXT, color TEXT)")
	createCategories.Exec()

	createPosts, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY, 
			title TEXT, 
			content TEXT, 
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
			author_id INTEGER NOT NULL, 
			category_id INTEGER NOT NULL, 
			FOREIGN KEY(author_id) REFERENCES users(id), 
			FOREIGN KEY(category_id) REFERENCES categories(id)
		)
	`)
	createPosts.Exec()

	// var categories = make(map[string]string)
	// categories["Technology"] = "red"
	// categories["Design"] = "blue"
	// categories["Environment"] = "green"

	// for category, color := range categories {
	// 	_, err = db.Exec(`INSERT INTO categories(name, color) VALUES(?, ?)`, category, color)
	// }

}

func checkInternalServerError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func checkCookie(w http.ResponseWriter, r *http.Request) (bool, User) {
	c, err := r.Cookie("session_token")

	var user User

	if err != nil {
		//User is not logged in
		return false, user
	}

	sessionToken := c.Value

	err = db.QueryRow("SELECT id, email, username, avatar FROM users WHERE session=?",
		sessionToken).Scan(&user.ID, &user.Email, &user.Username, &user.Avatar)

	// checkInternalServerError(err, w)
	if err != nil {
		return false, user
	}

	return true, user
}

func createCookie(w http.ResponseWriter, email string) {

	sessionToken := uuid.Must(uuid.NewV4()).String()

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})

	addSession, err := db.Prepare(`
		UPDATE users SET session=?
		WHERE email=?
	`)
	checkInternalServerError(err, w)
	res, err := addSession.Exec(sessionToken, email)
	checkInternalServerError(err, w)
	_, err = res.RowsAffected()
	checkInternalServerError(err, w)

}

func deleteCookie(w http.ResponseWriter) {

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
	})
}
