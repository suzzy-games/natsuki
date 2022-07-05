package routes

import (
	"context"
	"fmt"
	"natsuki/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SQLQuery struct {
	Query     string   `json:"query"`
	Fields    []string `json:"fields"`
	Arguments []any    `json:"args"`
}

func PostgresQuery(ctx *gin.Context) {

	// Parse Request Body
	var Body SQLQuery
	if err := ctx.ShouldBindJSON(&Body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Invalid Body: %s", err.Error()),
			"error":   0,
		})
		return
	}

	// Run Postgres Query
	rctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := db.Postgres.Query(rctx, Body.Query, Body.Arguments...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   50011,
		})
		return
	}
	defer rows.Close()

	// Return Rows
	rowData := make([][]any, 0)
	for rows.Next() {
		// Return Values
		val, err := rows.Values()
		rowData = append(rowData, val)

		// Return Error
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
				"code":    50012,
			})
			return
		}
	}

	// Return Error (if any)
	if rows.Err() != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": rows.Err().Error(),
			"error":   50011,
		})
		return
	}

	// Return Rows
	ctx.JSON(http.StatusOK, gin.H{
		"result": rowData,
	})
}
