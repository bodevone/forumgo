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
