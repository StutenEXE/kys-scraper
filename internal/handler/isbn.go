// internal/handler/isbn.go
package handler

import (
	"net/http"

	"scrapers/internal/scraper"

	"github.com/gin-gonic/gin"
)

type ISBNRequest struct {
	ISBN string `json:"isbn" binding:"required"`
}

type ISBNHandler struct {
	Scraper *scraper.ISBNScraper
}

func (h *ISBNHandler) Handle(c *gin.Context) {
	var req ISBNRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.Scraper.Scrape(c.Request.Context(), req.ISBN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
