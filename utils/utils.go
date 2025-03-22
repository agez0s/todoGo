package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/agez0s/todoGo/config"
	"github.com/agez0s/todoGo/schema"
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

// func GenerateToken(u string) string {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"username": u,
// 		"exp":      time.Now().Add(time.Hour * 24).Unix(),
// 	})
// 	tokenString, _ := token.SignedString([]byte(config.JWT_SECRET))
// 	return tokenString
// }

type Claims struct {
	Username string `json:"username"`
	userID   string `json:"userID"`
	exp      *int64 `json:"exp"`
}

func GenerateToken(d schema.User) (string, error) {

	claims := jwt.MapClaims{}
	claims["username"] = d.Username
	claims["userID"] = d.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWT_SECRET))

}

// func GetUsername(ctx *gin.Context) string {
// 	t, ex := ctx.Get("claims")
// 	if !ex {
// 		SendError(ctx, http.StatusUnauthorized, "invalid token")
// 		return ""
// 	}
// 	claims := t.(jwt.MapClaims)
// 	return claims["username"].(string)

// }
