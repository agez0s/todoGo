package handler

// auth
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//todos

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueAt       string `json:"dueAt"`
}

type UpdateTodoRequest struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	DueAt       *string `json:"dueAt"`
}
