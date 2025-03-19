package handler

import (
	"fmt"
	"net/http"

	"github.com/agez0s/todoGo/config"
	"github.com/agez0s/todoGo/schema"
	"github.com/agez0s/todoGo/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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

func (r *LoginRequest) ValidateLoginRequest() error {
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
	r := CreateUserRequest{}

	ctx.BindJSON(&r)

	if err := r.ValidateCreateUser(); err != nil {
		logger.ErrorF("validation error: %v", err.Error())
		utils.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var hashedPassword *string

	hashTemp, err := hashPassword(r.Password)
	if err != nil {
		logger.ErrorF("error hashing password: %v", err.Error())
		utils.SendError(ctx, http.StatusInternalServerError, "error creating user")
		return
	}
	hashedPassword = &hashTemp

	newuser := schema.User{
		Username: r.Username,
		Password: *hashedPassword,
	}

	if err := db.Create(&newuser).Error; err != nil {
		fmt.Println("err: ", err)
		logger.ErrorF("error creating user: %v", err.Error())
		utils.SendError(ctx, http.StatusInternalServerError, "error creating user")
		return
	}
	newToken, err1 := utils.GenerateToken(newuser)
	if err1 != nil {
		logger.ErrorF("error generating token: %v", err1.Error())
		utils.SendError(ctx, http.StatusInternalServerError, "error generating token")
		return
	}
	utils.SendSuccess(ctx, "create-user", gin.H{"username": newuser.Username, "token": newToken})
}

func LoginUserHandler(ctx *gin.Context) {
	r := LoginRequest{}
	ctx.BindJSON(&r)

	if err := r.ValidateLoginRequest(); err != nil {
		fmt.Println(r)
		logger.ErrorF("validation error: %v", err.Error())
		utils.SendError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var user schema.User
	if err := db.Where("username = ?", r.Username).First(&user).Error; err != nil {
		logger.ErrorF("error finding user: %v", err.Error())
		utils.SendError(ctx, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if !checkPassword(r.Password, user.Password) {
		logger.ErrorF("error invalid password")
		utils.SendError(ctx, http.StatusUnauthorized, "invalid credentials")
		return
	}

	newToken, err := utils.GenerateToken(user)
	if err != nil {
		logger.ErrorF("error generating token: %v", err.Error())
		utils.SendError(ctx, http.StatusInternalServerError, "error generating token")
		return
	}
	utils.SendSuccess(ctx, "login-user", gin.H{"username": user.Username, "token": newToken})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendError(ctx, http.StatusUnauthorized, "missing token")
			ctx.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		const bearerPrefix = "Bearer "
		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			utils.SendError(ctx, http.StatusUnauthorized, "invalid token format")
			ctx.Abort()
			return
		}
		tokenString := authHeader[len(bearerPrefix):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.JWT_SECRET), nil
		})
		if err != nil || !token.Valid {
			utils.SendError(ctx, http.StatusUnauthorized, "invalid token")
			ctx.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Set("claims", claims)
			fmt.Println("setou:", claims)
		} else {
			utils.SendError(ctx, http.StatusUnauthorized, "invalid token claims")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
