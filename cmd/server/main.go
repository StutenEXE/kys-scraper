package main

import (
	"scrapers/internal/dispatcher"
	"scrapers/internal/handler"
	"scrapers/internal/registry"

	"github.com/gin-gonic/gin"
)

func main() {
	disp := dispatcher.New(registry.All())
	h := &handler.ScrapeHandler{Dispatcher: disp}

	r := gin.Default()
	r.POST("/scrape", h.Handle)
	r.Run(":8081")
}
