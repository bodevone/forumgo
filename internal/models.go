package internal

// type Cost struct {
// 	Id             int64   `json:"id"`
// 	ElectricAmount int64   `json:"electric_amount"`
// 	ElectricPrice  float64 `json:"electric_price"`
// 	WaterAmount    int64   `json:"water_amount"`
// 	WaterPrice     float64 `json:"water_price"`
// 	CheckedDate    string  `json:"checked_date"`
// }

// User stores user data
type User struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Avatar       int    `json:"avatar"`
	SessionToken string `json:"sessionToken"`
}

// Post stores post data
type Post struct {
	ID           int64  `json:"id"`
	Title        string `json:"email"`
	Content      string `json:"username"`
	Timestamp    string `json:"timestamp"`
	Author       int64  `json:"autor"`
	AuthorName   string `json:"author_name"`
	Category     int64  `json:"category"`
	CategoryName string `json:"category_name"`
}

// Category displays category of post
type Category struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Comment displays comment of user in post
type Comment struct {
	ID           int64  `json:"id"`
	Text         string `json:"name"`
	Timestamp    string `json:"timestamp"`
	Author       int64  `json:"autor"`
	AuthorName   string `json:"author_name"`
	Post         int64  `json:"post"`
	Likes        string `json:"likes"`
	UserLiked    bool   `json:"userliked"`
	UserDisliked bool   `json:"userdisliked"`
}

// Like on post from user
type Like struct {
	UserID int64 `json:"userid"`
	PostID int64 `json:"postid"`
}

// Dislike on post from user
type Dislike struct {
	UserID int64 `json:"userid"`
	PostID int64 `json:"postid"`
}

// CommentLike shows likes on comment
type CommentLike struct {
	UserID    int64 `json:"userid"`
	CommentID int64 `json:"commentid"`
}
