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
