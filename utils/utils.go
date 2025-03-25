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

func SendError(c *gin.Context, code int, msg string) {
	c.Header("Content-Type", "application/json")
	c.JSON(code, gin.H{
		"message":   msg,
		"errorCode": code,
	})
}

func SendSuccess(c *gin.Context, op string, data interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("operation %s successful", op),
		"data":    data,
	})
}

// type Claims struct {
// 	Username string `json:"username"`
// 	userID   string `json:"userID"`
// 	exp      *int64 `json:"exp"`
// }

func GenerateToken(d schema.User) (string, error) {

	claims := jwt.MapClaims{}
	claims["username"] = d.Username
	claims["userID"] = d.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWT_SECRET))

}
