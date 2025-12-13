package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LocationService defines the contract required for location-related operations.
type LocationService interface {
	GenerateAddressSuggestions(query string) []string
}

// LocationHandler handles location-related HTTP endpoints.
type LocationHandler struct {
	locationService LocationService
}

// NewLocationHandler creates a new LocationHandler instance.
func NewLocationHandler(service LocationService) *LocationHandler {
	return &LocationHandler{locationService: service}
}

// SuggestLocations GET /locations/suggest?query=...
// Response: 200 { "suggestions": [...] } | 400 invalid request
func (h *LocationHandler) SuggestLocations(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.String(http.StatusBadRequest, "query parameter is required")
		return
	}

	suggestions := h.locationService.GenerateAddressSuggestions(query)
	if suggestions == nil {
		log.Println("failed to generate address suggestions")
		c.String(http.StatusBadRequest, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}
