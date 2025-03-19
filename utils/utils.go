package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/agez0s/todoGo/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

func SendError(ctx *gin.Context, code int, msg string) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(code, gin.H{
		"message":   msg,
		"errorCode": code,
	})
}

func SendSuccess(ctx *gin.Context, op string, data interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("operation %s successful", op),
		"data":    data,
	})
}

func GenerateToken(u string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": u,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(config.JWT_SECRET))
	return tokenString
}
