package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	// Return 200 OK
	ctx.JSON(http.StatusOK, gin.H{
		"online": true,
	})
}
