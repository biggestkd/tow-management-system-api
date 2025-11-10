package handler

import (
	"context"
	"log"
	"net/http"
	"tow-management-system-api/model"

	"github.com/gin-gonic/gin"
)

// PriceService defines the contract for Price-related business logic.
type PriceService interface {
	FindPricesByCompanyId(ctx context.Context, companyId string) ([]*model.Price, error)
	SetPrice(ctx context.Context, prices []*model.Price) error
}

// PriceHandler handles HTTP routes for Price-related operations.
type PriceHandler struct {
	priceService PriceService
}

// NewPriceHandler creates a new PriceHandler instance.
func NewPriceHandler(service PriceService) *PriceHandler {
	return &PriceHandler{priceService: service}
}

// GetPrices GET /prices/company/:companyId
// Retrieves all prices for a given company.
// Response: 200 [Price] | 400/404/500 generic error text
func (h *PriceHandler) GetPrices(c *gin.Context) {
	companyId := c.Param("companyId")
	if companyId == "" {
		c.String(http.StatusBadRequest, "company id is required")
		return
	}

	prices, err := h.priceService.FindPricesByCompanyId(c.Request.Context(), companyId)

	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, prices)
}

// PutPrices PUT /prices
// Creates or sets multiple prices.
// Request: [Price] (array of prices)
// Response: 204 | 400/500 generic error text
func (h *PriceHandler) PutPrices(c *gin.Context) {
	var body []*model.Price
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.priceService.SetPrice(c.Request.Context(), body); err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, "something went wrong")
		return
	}

	c.Status(http.StatusNoContent)
}
