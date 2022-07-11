package main

import (
	_ "natsuki/db"
	"natsuki/kaho"
	"natsuki/middleware"
	"natsuki/routes"
	"natsuki/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	// Create New Engine
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Set Trust Proxies
	switch utils.GetEnvDefault("NATSUKI_PROXY", "none") {
	case "cloudflare":
		// Retrieve IPs from Cloudflare
		proxyIPs, err := utils.GetCloudflareProxyIPs()
		if err != nil {
			kaho.Log(kaho.FATAL, "GIN", "Failed to get Cloudflare Proxies", err.Error())
		}

		// Set Trust Proxies & Platform
		engine.SetTrustedProxies(proxyIPs)
		engine.TrustedPlatform = gin.PlatformCloudflare
		kaho.Log(kaho.DEBUG, "GIN", "Using Cloudflare as Proxy", nil)

	case "none":
		kaho.Log(kaho.WARNING, "GIN", "Running With No Proxy", nil)
	default:
		kaho.Log(kaho.FATAL, "GIN", "Unsupported Proxy Mode, allowed values are: cloudflare, none", nil)
	}

	// Create Routes
	engine.Use(gin.Recovery(), middleware.LogRequest)
	engine.POST("/redis", middleware.VerifyToken, routes.RedisHandler)
	engine.POST("/sql", middleware.VerifyToken, routes.PostgresQuery)
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error":   0,
			"message": "404: Error Not Found",
		})
	})

	// See if we should serve via SSL
	if os.Getenv("NATSUKI_ENABLE_SSL") != "" {
		// Intialize with SSL Enabled
		kaho.Log(kaho.DEFAULT, "GIN", "Listening and Serving HTTP on :443", nil)
		if err := engine.RunTLS(":443", os.Getenv("SSL_CERT_PATH"), os.Getenv("SSL_KEY_PATH")); err != nil {
			kaho.Log(kaho.FATAL, "GIN", "Failed to start HTTP Server", err.Error())
		}
	} else {
		// Initialize without SSL
		kaho.Log(kaho.DEFAULT, "GIN", "Listening and Serving HTTP on :80", nil)
		if err := engine.Run(":80"); err != nil {
			kaho.Log(kaho.FATAL, "GIN", "Failed to start HTTP Server", err.Error())
		}
	}
}
