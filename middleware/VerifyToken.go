package middleware

import (
	"fmt"
	"natsuki/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var jwtSecret []byte = []byte(utils.GetEnvDefault("NATSUKI_JWT", "your-256-bit-secret"))

func VerifyToken(ctx *gin.Context) {

	// Check if Server ID was Given
	if ctx.GetHeader("RBX-Server-Id") == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Missing RBX-Server-Id Header",
			"error":   0,
		})
		return
	}

	// Check if Token was Given
	authorization := ctx.GetHeader("Authorization")
	if authorization == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Missing Authorization Header",
			"error":   0,
		})
		return
	}

	// Verify JWT Token
	_, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {

		// Verify Token Algorithm
		if token.Method.Alg() != "HS256" {
			fmt.Println("wrong method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Verify Token with Secret
		return jwtSecret, nil
	})

	// Check if parsing encountered an error
	if err != nil {
		// Return Error
		fmt.Println("parse fail")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": fmt.Sprintf("Invalid JWT Token: %s", err.Error()),
			"error":   0,
		})
	} else {
		// Continue Chain
		ctx.Next()
	}

}
