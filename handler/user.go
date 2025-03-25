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

func CreateUserHandler(c *gin.Context) {
	r := CreateUserRequest{}

	c.BindJSON(&r)

	if err := r.ValidateCreateUser(); err != nil {
		logger.ErrorF("validation error: %v", err.Error())
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}
	var hashedPassword *string

	hashTemp, err := hashPassword(r.Password)
	if err != nil {
		logger.ErrorF("error hashing password: %v", err.Error())
		utils.SendError(c, http.StatusInternalServerError, "error creating user")
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
		utils.SendError(c, http.StatusInternalServerError, "error creating user")
		return
	}
	newToken, err1 := utils.GenerateToken(newuser)
	if err1 != nil {
		logger.ErrorF("error generating token: %v", err1.Error())
		utils.SendError(c, http.StatusInternalServerError, "error generating token")
		return
	}
	utils.SendSuccess(c, "create-user", gin.H{"username": newuser.Username, "token": newToken})
}

func LoginUserHandler(c *gin.Context) {
	r := LoginRequest{}
	c.BindJSON(&r)

	if err := r.ValidateLoginRequest(); err != nil {
		fmt.Println(r)
		logger.ErrorF("validation error: %v", err.Error())
		utils.SendError(c, http.StatusBadRequest, err.Error())
		return
	}

	var user schema.User
	if err := db.Where("username = ?", r.Username).First(&user).Error; err != nil {
		logger.ErrorF("error finding user: %v", err.Error())
		utils.SendError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if !checkPassword(r.Password, user.Password) {
		logger.ErrorF("error invalid password")
		utils.SendError(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	newToken, err := utils.GenerateToken(user)
	if err != nil {
		logger.ErrorF("error generating token: %v", err.Error())
		utils.SendError(c, http.StatusInternalServerError, "error generating token")
		return
	}
	utils.SendSuccess(c, "login-user", gin.H{"username": user.Username, "token": newToken})
}

func GetProfileHandler(c *gin.Context) {
	u := c.MustGet("claims").(jwt.MapClaims)
	uid, ok := u["userID"].(float64)
	if !ok {
		utils.SendError(c, http.StatusInternalServerError, "error getting user id")
		return
	}

	var user schema.User
	if err := db.Where("id = ?", int(uid)).First(&user).Error; err != nil {
		logger.ErrorF("error finding user: %v", err.Error())
		utils.SendError(c, http.StatusUnauthorized, "invalid user")
		return
	}

	utils.SendSuccess(c, "get-profile", gin.H{"username": user.Username})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendError(c, http.StatusUnauthorized, "missing token")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		const bearerPrefix = "Bearer "
		if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
			utils.SendError(c, http.StatusUnauthorized, "invalid token format")
			c.Abort()
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
			utils.SendError(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("claims", claims)
		} else {
			utils.SendError(c, http.StatusUnauthorized, "invalid token claims")
			c.Abort()
			return
		}

		c.Next()
	}
}
