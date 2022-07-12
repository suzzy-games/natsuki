package middleware

import (
	"fmt"
	"natsuki/kaho"
	"time"

	"github.com/gin-gonic/gin"
)

type HTTPLogPayload struct {
	RemoteAddr string `json:"remoteAddr"`
	Latency    int64  `json:"latency"`
	Status     int    `json:"status"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Query      string `json:"query"`
	Token      string `json:"token"`
	ServerId   string `json:"serverId"`
}

func LogRequest(ctx *gin.Context) {
	// Start Tracking
	startTime := time.Now()
	ctx.Next()
	endTime := time.Now()

	// Create HTTP Message
	// 2003/06/26 12:34:56 [HTTP] [INFO] {STATUS} {METHOD} {PATH}{QUERY} {RESPONSE_TIME} {REMOTE_IP}
	logMessage := fmt.Sprintf(
		"%d %s %s%s %s %s",
		ctx.Writer.Status(),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		ctx.Request.URL.RawQuery,
		endTime.Sub(startTime).String(),
		ctx.RemoteIP(),
	)

	// Log Request
	kaho.KahoLogRaw(kaho.KahoLogEntry{
		Severity:  kaho.INFO,
		Timestamp: startTime,
		Service:   "HTTP",
		Message:   logMessage,
		Payload: &HTTPLogPayload{
			RemoteAddr: ctx.RemoteIP(),
			Latency:    time.Now().UnixNano() - startTime.UnixNano(),
			Status:     ctx.Writer.Status(),
			Method:     ctx.Request.Method,
			Path:       ctx.Request.URL.Path,
			Query:      ctx.Request.URL.RawQuery,
			ServerId:   ctx.GetHeader("RBX-Server-Id"),
			Token:      ctx.GetHeader("Authorization"),
		},
	})
}
