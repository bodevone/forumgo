package internal

import "net/http"

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

func getCategoryName(w http.ResponseWriter, categoryID int) string {
	categoryName := ""
	err = db.QueryRow("SELECT name FROM categories WHERE id=?",
		categoryID).Scan(&categoryName)

	checkInternalServerError(err, w)

	return categoryName
}
