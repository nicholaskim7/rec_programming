package models

import (
	"time"
)

type Post struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	Username string `json:"username"`
	Title string `json:"title"`
	Body string `json:"body"`
	Tag string `json:"tag"`
	DateCreated time.Time `json:"date_created"`
}