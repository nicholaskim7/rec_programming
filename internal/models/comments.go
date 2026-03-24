package models

import (
	"time"
)

type Comment struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Username string `json:"username"`
	PostID int64 `json:"post_id"`
	Comment string `json:"comment"`
	DateCreated time.Time `json:"date_created"`
}