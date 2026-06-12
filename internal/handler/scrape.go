package handler

import (
	"net/http"

	"scrapers/internal/dispatcher"

	"github.com/gin-gonic/gin"
)

type ScrapeRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type ScrapeHandler struct {
	Dispatcher *dispatcher.Dispatcher
}

func (h *ScrapeHandler) Handle(c *gin.Context) {
	var req ScrapeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sc, err := h.Dispatcher.For(req.URL)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	result, err := sc.Scrape(c.Request.Context(), req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
