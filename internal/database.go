package internal

import (
	"database/sql"
	"net/http"
	"strconv"
)

// InitDb starts database
func InitDb() {
	db, err = sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	createUsers, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY,
			email TEXT,
			username TEXT,
			password TEXT,
			avatar INTEGER,
			session TEXT
			)
		`)
	createUsers.Exec()

	createCategories, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY, 
			name TEXT, 
			color TEXT,
			UNIQUE(name)
		)
	`)
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

	createComments, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY,
			text TEXT, 
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
			author_id INTEGER NOT NULL, 
			post_id INTEGER NOT NULL, 
			FOREIGN KEY(author_id) REFERENCES users(id), 
			FOREIGN KEY(post_id) REFERENCES posts(id)
		)
	`)
	createComments.Exec()

	createLikes, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS likes (
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			UNIQUE(user_id, post_id),
			FOREIGN KEY(user_id) REFERENCES users(id), 
			FOREIGN KEY(post_id) REFERENCES posts(id)
		)
	`)
	createLikes.Exec()

	createDislikes, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS dislikes (
			user_id INTEGER NOT NULL,
			post_id INTEGER NOT NULL,
			UNIQUE(user_id, post_id),
			FOREIGN KEY(user_id) REFERENCES users(id), 
			FOREIGN KEY(post_id) REFERENCES posts(id)
		)
	`)
	createDislikes.Exec()

	createCommentLikes, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS commentlikes (
			user_id INTEGER NOT NULL,
			comment_id INTEGER NOT NULL,
			UNIQUE(user_id, comment_id),
			FOREIGN KEY(user_id) REFERENCES users(id), 
			FOREIGN KEY(comment_id) REFERENCES comments(id)
		)
	`)
	createCommentLikes.Exec()

	createCommentDislikes, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS commentdislikes (
			user_id INTEGER NOT NULL,
			comment_id INTEGER NOT NULL,
			UNIQUE(user_id, comment_id),
			FOREIGN KEY(user_id) REFERENCES users(id), 
			FOREIGN KEY(comment_id) REFERENCES comments(id)
		)
	`)
	createCommentDislikes.Exec()

	var categories = make(map[string]string)
	categories["Technology"] = "red"
	categories["Design"] = "blue"
	categories["Environment"] = "green"

	for category, color := range categories {
		_, err = db.Exec(`INSERT OR IGNORE INTO categories(name, color) VALUES(?, ?)`, category, color)
	}

}

func getCategories(w http.ResponseWriter) []Category {
	categoryRows, err := db.Query("SELECT * FROM categories")
	checkInternalServerError(err, w)
	var categories []Category
	var category Category
	for categoryRows.Next() {
		err = categoryRows.Scan(&category.ID, &category.Name, &category.Color)
		checkInternalServerError(err, w)
		categories = append(categories, category)
	}
	return categories
}

func getPosts(w http.ResponseWriter) []Post {
	postRows, err := db.Query("SELECT * FROM posts")
	checkInternalServerError(err, w)
	var posts []Post
	var post Post
	for postRows.Next() {
		err = postRows.Scan(&post.ID, &post.Title, &post.Content, &post.Timestamp, &post.Author, &post.Category)
		checkInternalServerError(err, w)
		posts = append(posts, post)
	}
	return posts
}

func getPostsOfCategory(category int) ([]Post, error) {

	postRows, err1 := db.Query("SELECT * FROM posts WHERE category_id=?", category)

	var posts []Post
	var post Post
	for postRows.Next() {
		err = postRows.Scan(&post.ID, &post.Title, &post.Content, &post.Timestamp, &post.Author, &post.Category)
		posts = append(posts, post)
	}
	return posts, err1
}

func getPostsOfUser(w http.ResponseWriter, user int) []Post {

	postRows, err := db.Query("SELECT * FROM posts WHERE author_id=?", user)
	checkInternalServerError(err, w)

	var posts []Post
	var post Post
	for postRows.Next() {
		err = postRows.Scan(&post.ID, &post.Title, &post.Content, &post.Timestamp, &post.Author, &post.Category)
		posts = append(posts, post)
		checkInternalServerError(err, w)
	}
	return posts
}

func getLikedPostsOfUser(w http.ResponseWriter, user int) []Post {
	likeRows, err := db.Query("SELECT post_id FROM likes WHERE user_id=?", user)
	checkInternalServerError(err, w)

	var posts []Post
	var post Post

	var postID int64

	for likeRows.Next() {
		err = likeRows.Scan(&postID)
		checkInternalServerError(err, w)

		err = db.QueryRow("SELECT * FROM posts WHERE id=?",
			postID).Scan(&post.ID, &post.Title, &post.Content, &post.Timestamp, &post.Author, &post.Category)
		checkInternalServerError(err, w)

		posts = append(posts, post)
	}

	return posts

}

func getCategoryName(w http.ResponseWriter, categoryID int) string {
	categoryName := ""
	err = db.QueryRow("SELECT name FROM categories WHERE id=?",
		categoryID).Scan(&categoryName)

	checkInternalServerError(err, w)

	return categoryName
}

func getComments(w http.ResponseWriter, postID int) []Comment {

	commentRows, err := db.Query("SELECT * FROM comments WHERE post_id=?", postID)

	checkInternalServerError(err, w)

	var comments []Comment
	var comment Comment
	for commentRows.Next() {
		err = commentRows.Scan(&comment.ID, &comment.Text, &comment.Timestamp, &comment.Author, &comment.Post)
		comments = append(comments, comment)
		checkInternalServerError(err, w)
	}
	return comments
}

func getCommentsLikes(w http.ResponseWriter, comments []Comment) []Comment {
	for i, comment := range comments {
		likes := 0
		commentLikeRows, err := db.Query("SELECT * FROM commentlikes WHERE comment_id=?", comment.ID)
		checkInternalServerError(err, w)
		for commentLikeRows.Next() {
			likes++
		}

		commentDislikesRows, err := db.Query("SELECT * FROM commentdislikes WHERE comment_id=?", comment.ID)
		checkInternalServerError(err, w)
		for commentDislikesRows.Next() {
			likes--
		}

		likesString := ""

		if likes > 0 {
			likesString = "+" + strconv.Itoa(likes)
		} else if likes < 0 {
			likesString = strconv.Itoa(likes)
		} else {
			likesString = "0"
		}

		comments[i].Likes = likesString
	}

	return comments
}

func getUserCommentsLikes(w http.ResponseWriter, comments []Comment, userID int) []Comment {
	for i, comment := range comments {
		commentLiked := false
		userCommentLikeRows, _ := db.Query("SELECT * FROM commentlikes WHERE comment_id=? AND user_id=?", comment.ID, userID)

		count := 0
		for userCommentLikeRows.Next() {
			count++
		}
		if count >= 1 {
			commentLiked = true
		}

		comments[i].UserLiked = commentLiked
	}

	return comments
}

func getUserCommentsDislikes(w http.ResponseWriter, comments []Comment, userID int) []Comment {
	for i, comment := range comments {
		commentDisliked := false
		userCommentDislikeRows, _ := db.Query("SELECT * FROM commentdislikes WHERE comment_id=? AND user_id=?", comment.ID, userID)

		count := 0
		for userCommentDislikeRows.Next() {
			count++
		}
		if count >= 1 {
			commentDisliked = true
		}

		comments[i].UserDisliked = commentDisliked
	}

	return comments
}

func getLikes(w http.ResponseWriter, postID int) int {

	likeRows, err := db.Query("SELECT * FROM likes WHERE post_id=?", postID)
	// likeRows, err := db.Query("SELECT * FROM likes")

	checkInternalServerError(err, w)
	count := 0
	for likeRows.Next() {
		count++
	}

	return count
}

func getDislikes(w http.ResponseWriter, postID int) int {

	dislikeRows, err := db.Query("SELECT * FROM dislikes WHERE post_id=?", postID)

	checkInternalServerError(err, w)
	count := 0
	for dislikeRows.Next() {
		count++
	}

	return count
}

func getUserLike(w http.ResponseWriter, postID int, userID int) bool {
	likeRows, err := db.Query("SELECT * FROM likes WHERE post_id=? AND user_id=?", postID, userID)

	if err != nil {
		return false
	}
	count := 0
	for likeRows.Next() {
		count++
	}
	if count >= 1 {
		return true
	}
	return false
}

func getUserDislike(w http.ResponseWriter, postID int, userID int) bool {
	likeRows, err := db.Query("SELECT * FROM dislikes WHERE post_id=? AND user_id=?", postID, userID)

	if err != nil {
		return false
	}
	count := 0
	for likeRows.Next() {
		count++
	}
	if count >= 1 {
		return true
	}
	return false
}

func addLike(w http.ResponseWriter, postID int, userID int) {
	// Add like
	likeRows, err := db.Query("SELECT * FROM likes WHERE post_id=? AND user_id=?", postID, userID)

	if err != nil {
		checkInternalServerError(err, w)
	}
	count := 0
	for likeRows.Next() {
		count++
	}
	if count >= 1 {
		// Remove like if exists
		_, err = db.Exec(`
			DELETE from likes WHERE user_id=? AND post_id=?
		`, userID, postID)
		checkInternalServerError(err, w)
	} else {
		// Ad like if not exists
		_, err = db.Exec(`
			INSERT OR IGNORE INTO likes (user_id, post_id) VALUES (?, ?)
		`, userID, postID)
		checkInternalServerError(err, w)
	}

	// Remove dislike
	_, err = db.Exec(`
		DELETE from dislikes WHERE user_id=? AND post_id=?
	`, userID, postID)
	checkInternalServerError(err, w)
}

func addDislike(w http.ResponseWriter, postID int, userID int) {

	// Add dislike
	dislikeRows, err := db.Query("SELECT * FROM dislikes WHERE post_id=? AND user_id=?", postID, userID)

	if err != nil {
		checkInternalServerError(err, w)
	}
	count := 0
	for dislikeRows.Next() {
		count++
	}
	if count >= 1 {
		// Remove dislike if exists
		_, err = db.Exec(`
			DELETE from dislikes WHERE user_id=? AND post_id=?
		`, userID, postID)
		checkInternalServerError(err, w)
	} else {
		// Add dislike if not exists
		_, err = db.Exec(`
			INSERT OR IGNORE INTO dislikes (user_id, post_id) VALUES (?, ?)
		`, userID, postID)
		checkInternalServerError(err, w)
	}

	// Remove like
	_, err = db.Exec(`
		DELETE from likes WHERE user_id=? AND post_id=?
	`, userID, postID)
	checkInternalServerError(err, w)
}

func addLikeToComment(w http.ResponseWriter, commentID int, userID int) {
	// Add like
	commentLikeRows, err := db.Query("SELECT * FROM commentlikes WHERE user_id=? AND comment_id=?", userID, commentID)

	if err != nil {
		checkInternalServerError(err, w)
	}
	count := 0
	for commentLikeRows.Next() {
		count++
	}
	if count >= 1 {
		// Remove like if exists
		_, err = db.Exec(`
			DELETE from commentlikes WHERE user_id=? AND comment_id=?
		`, userID, commentID)
		checkInternalServerError(err, w)
	} else {
		// Add like if not exists
		_, err = db.Exec(`
			INSERT OR IGNORE INTO commentlikes (user_id, comment_id) VALUES (?, ?)
		`, userID, commentID)
		checkInternalServerError(err, w)
	}

	// Remove dislike
	_, err = db.Exec(`
		DELETE from commentdislikes WHERE user_id=? AND comment_id=?
	`, userID, commentID)
	checkInternalServerError(err, w)
}

func addDislikeToComment(w http.ResponseWriter, commentID int, userID int) {
	// Add dislike
	commentLikeRows, err := db.Query("SELECT * FROM commentdislikes WHERE user_id=? AND comment_id=?", userID, commentID)

	if err != nil {
		checkInternalServerError(err, w)
	}
	count := 0
	for commentLikeRows.Next() {
		count++
	}
	if count >= 1 {
		// Remove dislike if exists
		_, err = db.Exec(`
			DELETE from commentdislikes WHERE user_id=? AND comment_id=?
		`, userID, commentID)
		checkInternalServerError(err, w)
	} else {
		// Add dislike if not exists
		_, err = db.Exec(`
			INSERT OR IGNORE INTO commentdislikes (user_id, comment_id) VALUES (?, ?)
		`, userID, commentID)
		checkInternalServerError(err, w)
	}

	// Remove like
	_, err = db.Exec(`
		DELETE from commentlikes WHERE user_id=? AND comment_id=?
	`, userID, commentID)
	checkInternalServerError(err, w)
}
