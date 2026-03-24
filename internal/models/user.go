package models
import (
	"time"
)

type User struct {
	ID int64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Salt string `json:"salt"`
	Username string `json:"username"`
	DateCreated time.Time `json:"date_created"`
}

type UserResponse struct {
	ID int64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Username string `json:"username"`
	DateCreated time.Time `json:"date_created"`
}

type UserLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
	User User `json:"user"`
}