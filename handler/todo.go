package handler

import (
	"fmt"
	"net/http"

	"github.com/agez0s/todoGo/schema"
	"github.com/agez0s/todoGo/utils"
	"github.com/gin-gonic/gin"
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
	c.BindJSON(&r)
	u := utils.GetUsername(c)
	if err := r.ValidateCreateTodo(); err != nil {
		logger.ErrorF("validation error: %v", err.Error())
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	newTodo := schema.Todo{
		Title:       r.Title,
		Description: r.Description,
	}
}
