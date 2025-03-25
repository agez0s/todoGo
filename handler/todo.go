package handler

import (
	"fmt"
	"net/http"
	"strconv"
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
	page, err := strconv.ParseUint(c.Query("page"), 10, 0)
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	u := c.MustGet("claims").(jwt.MapClaims)
	uid, ok := u["userID"].(float64)
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "error getting user id")
		return
	}

	var todos []schema.Todo
	limit := 10
	offset := (int(page) - 1) * limit
	if err := db.Where("user_id = ?", int(uid)).Limit(limit).Offset(offset).Find(&todos).Error; err != nil {
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

	if r.Title == "" && r.Description == "" && r.DueAt == nil {
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
	//make that if dueAt is nil, it should not be updated
	if r.DueAt != nil {

		parsedTime, err := time.Parse(time.RFC3339, *r.DueAt)
		if err != nil {
			logger.ErrorF("invalid due date format: %v", err.Error())
			utils.SendError(c, http.StatusBadRequest, "invalid due date format, expected RFC3339")
			return
		}
		todo.DueAt = &parsedTime
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
	now := time.Now()
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "error getting user id")
		return
	}
	if id == "" {
		utils.SendError(c, http.StatusBadRequest, "id is required")
		return
	}
	if err := db.Model(&schema.Todo{}).Where("id = ?", id).Where("user_id = ?", uid).Updates(schema.Todo{DoneTime: &now, Done: true}).Error; err != nil {
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
