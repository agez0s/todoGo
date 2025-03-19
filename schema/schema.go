package schema

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint
	Username string `gorm:"unique"`
	Password string
}

type Todo struct {
	gorm.Model
	ID          uint
	UserID      uint
	Title       string
	Done        bool
	DoneTime    *time.Time
	Description string
}

//responses

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type TodoResponse struct {
	ID          uint    `json:"id"`
	UserID      uint    `json:"user_id"`
	Title       string  `json:"title"`
	Done        bool    `json:"done"`
	DoneTime    *string `json:"done_time"`
	Description string  `json:"description"`
}
