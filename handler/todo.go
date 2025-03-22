package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/agez0s/todoGo/schema"
	"github.com/agez0s/todoGo/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (r *CreateTodoRequest) ValidateCreateTodo() error {
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	return nil
}

func CreateTodoHandler(c *gin.Context) {
	r := CreateTodoRequest{}
	u := c.MustGet("claims").(jwt.MapClaims)
	uid, ok := u["userID"].(float64)
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "error getting user id")
		return
	}

	fmt.Println(u)

	c.BindJSON(&r)

	if err := r.ValidateCreateTodo(); err != nil {
		logger.ErrorF("validation error: %v", err.Error())
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	var dueAt *time.Time
	if r.DueAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, r.DueAt)
		if err != nil {
			logger.ErrorF("invalid due date format: %v", err.Error())
			utils.SendError(c, http.StatusBadRequest, "invalid due date format, expected RFC3339")
			return
		}
		dueAt = &parsedTime
	}

	newTodo := schema.Todo{
		Title:       r.Title,
		Description: r.Description,
		DueAt:       dueAt,
		UserId:      int(uid),
	}

	if err := db.Create(&newTodo).Error; err != nil {
		logger.ErrorF("error creating todo: %v", err.Error())
		utils.SendError(c, http.StatusInternalServerError, "error creating todo")
		return
	}
	utils.SendSuccess(c, "create-todo", newTodo)
}

func ListTodosHandler(c *gin.Context) {
	u := c.MustGet("claims").(jwt.MapClaims)
	uid, ok := u["userID"].(float64)
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "error getting user id")
		return
	}

	var todos []schema.Todo
	if err := db.Where("user_id = ?", int(uid)).Find(&todos).Error; err != nil {
		logger.ErrorF("error getting todos: %v", err.Error())
		utils.SendError(c, http.StatusInternalServerError, "error getting todos")
		return
	}
	utils.SendSuccess(c, "list-todos", todos)
}

func UpdateTodoHandler(c *gin.Context) {
	r := UpdateTodoRequest{}
	u := c.MustGet("claims").(jwt.MapClaims)
	uid, ok := u["userID"].(float64)
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "error getting user id")
		return
	}

	c.BindJSON(&r)

	if r.Title == "" && r.Description == "" && r.DueAt == "" {
		utils.SendError(c, http.StatusBadRequest, "at least one field is required")
		return
	}

	var todo schema.Todo
	if err := db.Where("id = ? AND user_id = ?", r.ID, int(uid)).First(&todo).Error; err != nil {
		logger.ErrorF("error getting todo: %v", err.Error())
		utils.SendError(c, http.StatusInternalServerError, "error getting todo")
		return
	}

	todo.Title = r.Title
	todo.Description = r.Description
	parsedTime, err := time.Parse(time.RFC3339, r.DueAt)
	if err != nil {
		logger.ErrorF("invalid due date format: %v", err.Error())
		utils.SendError(c, http.StatusBadRequest, "invalid due date format, expected RFC3339")
		return
	}
	todo.DueAt = &parsedTime
	todo.Done = r.Done
	if r.Done {
		now := time.Now()
		todo.DoneTime = &now
	}
	if err := db.Save(&todo).Error; err != nil {
		logger.ErrorF("error updating todo: %v", err.Error())
		utils.SendError(c, http.StatusInternalServerError, "error updating todo")
		return
	}
	utils.SendSuccess(c, "update-todo", todo)
}

func MarkDoneHandler(c *gin.Context) {
	id := c.Query("id")
	u := c.MustGet("claims").(jwt.MapClaims)
	uid, ok := u["userID"].(float64)
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "error getting user id")
		return
	}
	if id == "" {
		utils.SendError(c, http.StatusBadRequest, "id is required")
		return
	}
	if err := db.Model(&schema.Todo{}).Where("id = ?", id).Where("user_id = ?", uid).Update("done", true).Error; err != nil {
		logger.ErrorF("error marking todo as done: %v", err.Error())
		utils.SendError(c, http.StatusInternalServerError, "error marking todo as done")
		return
	}

}

func DeleteTodoHandler(c *gin.Context) {
	id := c.Query("id")
	u := c.MustGet("claims").(jwt.MapClaims)
	uid, ok := u["userID"].(float64)
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "error getting user id")
		return
	}
	if id == "" {
		utils.SendError(c, http.StatusBadRequest, "id is required")
		return
	}
	if err := db.Where("id = ?", id).Where("user_id = ?", uid).Delete(&schema.Todo{}).Error; err != nil {
		logger.ErrorF("error deleting todo: %v", err.Error())
		utils.SendError(c, http.StatusInternalServerError, "error deleting todo")
		return
	}
}
