package main

import (
	"net/http"
	"os"
	"scrapers/internal/dispatcher"
	"scrapers/internal/handler"
	"scrapers/internal/registry"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	disp := dispatcher.New(registry.All())
	h := &handler.ScrapeHandler{Dispatcher: disp}

	r := gin.Default()

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{allowedOrigin},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type"},
	}))

	r.POST("/api/scraper/scrape", h.Handle)

	r.GET("/api/scraper/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8081")
}
