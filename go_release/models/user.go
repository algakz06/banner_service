package models

type User struct {
	Id       string `json:"id"`
	Username string `binding:"required" json:"username"`
	Password string `binding:"required" json:"password"`
	Role     string `json:"role"`
}
