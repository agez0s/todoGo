package handler

import (
	"fmt"
	"net/http"

	"github.com/agez0s/todoGo/schema"
	"github.com/agez0s/todoGo/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (r *CreateUserRequest) ValidateCreateUser() error {
	if r.Username == "" {
		return fmt.Errorf("username is required")
	}
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUserHandler(ctx *gin.Context) {
	request := CreateUserRequest{}

	ctx.BindJSON(&request)

	if err := request.ValidateCreateUser(); err != nil {
		logger.ErrorF("validation error: %v", err.Error())
		utils.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var hashedPassword *string

	hashTemp, err := hashPassword(request.Password)
	if err != nil {
		logger.ErrorF("error hashing password: %v", err.Error())
		utils.SendError(ctx, http.StatusInternalServerError, "error creating user")
		return
	}
	hashedPassword = &hashTemp

	newuser := schema.User{
		Username: request.Username,
		Password: *hashedPassword,
	}

	if err := db.Create(&newuser).Error; err != nil {
		fmt.Println("err: ", err)
		logger.ErrorF("error creating user: %v", err.Error())
		utils.SendError(ctx, http.StatusInternalServerError, "error creating user")
		return
	}
	newToken := utils.GenerateToken(newuser.Username)
	utils.SendSuccess(ctx, "create-user", gin.H{"username": newuser.Username, "token": newToken})
}

func LoginUserHandler(ctx *gin.Context) {
	request := LoginRequest{}

	ctx.BindJSON(&request)

	var user schema.User
	if err := db.Where("username = ?", request.Username).First(&user).Error; err != nil {
		logger.ErrorF("error finding user: %v", err.Error())
		utils.SendError(ctx, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if !checkPassword(request.Password, user.Password) {
		logger.ErrorF("error invalid password")
		utils.SendError(ctx, http.StatusUnauthorized, "invalid credentials")
		return
	}

	newToken := utils.GenerateToken(user.Username)
	utils.SendSuccess(ctx, "login-user", gin.H{"username": user.Username, "token": newToken})
}
