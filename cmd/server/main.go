package main

import (
	"log"
	"net/http"
	"os"
	"scrapers/internal/dispatcher"
	"scrapers/internal/handler"
	"scrapers/internal/registry"
	"scrapers/internal/scraper"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("env file not loaded: %v", err)
	}
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	disp := dispatcher.New(registry.All())
	sh := &handler.ScrapeHandler{Dispatcher: disp}
	isbnh := &handler.ISBNHandler{Scraper: scraper.NewISBNScraper()}

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

	r.POST("/api/scraper/scrape", sh.Handle)
	r.POST("/api/scraper/isbn", isbnh.Handle)

	r.GET("/api/scraper/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":8081")
}
