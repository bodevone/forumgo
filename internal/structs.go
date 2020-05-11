package internal

// ProfileData stores data for profile handler
type ProfileData struct {
	ProfileUser User
	Avatar1     bool
	Avatar2     bool
	Avatar3     bool
	Posts       []Post
}

// UserData stores data for profile handler
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
