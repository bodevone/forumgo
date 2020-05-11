package internal

// ProfileData stores data for profile page
type ProfileData struct {
	ProfileUser User
	Avatar1     bool
	Avatar2     bool
	Avatar3     bool
	Posts       []Post
}

// UserData stores data for user page
type UserData struct {
	LoggedIn    bool
	ProfileUser User
	ProfData    ProfileData
	ProfPosts   []Post
}

// IndexData stores data for index page
type IndexData struct {
	IndexUser  User
	LoggedIn   bool
	Categories []Category
	Posts      []Post
	Category   string
}

// PostData stores data for post page
type PostData struct {
	LoggedIn bool
	UserData User
	CurrPost Post
	Comments []Comment
}
