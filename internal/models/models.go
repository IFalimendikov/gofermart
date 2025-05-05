package models 

type User struct {
	ID string `json:"user_id"`
	Login string `json:"login"`
	Password string `json:"password"`
}