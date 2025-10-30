package handler

import (
	"context"
	"net/http"
	"tow-management-system-api/model"

	"github.com/gin-gonic/gin"
)

// TowService defines the contract for Tow-related business logic.
type TowService interface {
	FindTowsByCompanyId(ctx context.Context, companyId string) ([]*model.Tow, error)
}

// TowHandler handles HTTP routes for Tow-related operations.
type TowHandler struct {
	towService TowService
}

// NewTowHandler creates a new TowHandler instance.
func NewTowHandler(service TowService) *TowHandler {
	return &TowHandler{towService: service}
}

// GetTowHistory GET /company/:companyId/tows
// Retrieves all tow history records for a given company.
// Response: 200 [Tow] | 400/404/500 generic error text
func (h *TowHandler) GetTowHistory(c *gin.Context) {
	companyId := c.Param("companyId")
	if companyId == "" {
		c.String(http.StatusBadRequest, "company id is required")
		return
	}

	tows, err := h.towService.FindTowsByCompanyId(c.Request.Context(), companyId)

	if err != nil {
		c.String(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, tows)
}
