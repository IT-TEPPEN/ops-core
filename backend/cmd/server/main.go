package main

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Allow all origins (for development)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API Route
	api := r.Group("/api")
	{
		api.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Hello from Go backend!",
			})
		})
	}

	// Static file serving (Frontend)
	r.StaticFile("/", "./frontend/dist/index.html")
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")
	r.StaticFile("/vite.svg", "./frontend/dist/vite.svg")
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			filePath := "./frontend/dist" + c.Request.URL.Path
			absPath, _ := filepath.Abs(filePath)
			if _, err := http.Dir(filepath.Dir(absPath)).Open(filepath.Base(absPath)); err != nil {
				c.File("./frontend/dist/index.html")
				return
			}
			c.File(filePath)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
		}
	})

	port := "8080"
	println("Backend server listening at http://localhost:" + port)
	r.Run(":" + port)
}
