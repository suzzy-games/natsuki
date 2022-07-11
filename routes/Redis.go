package routes

import (
	"context"
	"fmt"
	"natsuki/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mediocregopher/radix/v4"
)

func RedisHandler(ctx *gin.Context) {

	// Parse Request Body
	var Command []string
	if err := ctx.ShouldBindJSON(&Command); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Invalid Body: %s", err.Error()),
			"error":   0,
		})
		return
	}

	// Make Redis Request
	rctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var Response interface{}
	if err := db.Redis.Do(rctx, radix.Cmd(&Response, Command[0], Command[1:]...)); err != nil {
		// Return Error Message
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   50011,
			"message": err.Error(),
		})
		return
	}

	// Return Server Response
	ctx.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprintf("%s", Response),
	})
}
