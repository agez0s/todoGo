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
	Todos    []Todo `gorm:"foreignKey:UserId;references:id"`
}

type Todo struct {
	gorm.Model
	UserId      int
	ID          uint
	Title       string
	DueAt       *time.Time
	Done        bool `gorm:"default:false"`
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
	DueAt       *string `json:"due_at"`
	Description string  `json:"description"`
}
