package model

// @Description Telegram user model
type User struct {
	ID       int    `json:"id"`       // User id example: 1
	Username string `json:"username"` // User username example: ivanptr21
	Fullname string `json:"fullname"` // User fullname example Ivan Petrovich
	ImageURL string `json:"imageUrl"` // User image url from firebase
}

// @Description User model
type WebUser struct {
	ID       int    `json:"id"`       // User id example: 1
	Email    string `json:"email"`    // User email example: test@test.com
	Password string `json:"password"` // user Password example: d1e8a70b5ccab1dc2f56bbf7e99f064a660c08e361a35751b9c483c88943d082
}

type UserDTO struct {
	Username string `db:"username"`
	Fullname string `db:"fullname"`
	ImageURL string `db:"imageurl"`
}
