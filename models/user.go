package models

type User struct {
  Id       string `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
