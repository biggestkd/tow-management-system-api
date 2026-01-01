package handler

import (
	"context"
	"log"
	"net/http"
	"tow-management-system-api/model"

	"github.com/gin-gonic/gin"
)

// TowService defines the contract for Tow-related business logic.
type TowService interface {
	ScheduleTow(ctx context.Context, towRequest *model.Tow, schedulingLink string) (*model.Tow, error)
	FindTowsByCompanyId(ctx context.Context, companyId string) ([]*model.Tow, error)
	UpdateTow(ctx context.Context, towId string, update *model.Tow) error
	GetEstimate(ctx context.Context, companyId string, pickup string, dropoff string) (int64, error)
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

// PostTow POST /tows/:companyId
// Create a new tow request for the given company.
// Request: Tow payload in JSON body
// Response: 201 [Tow] | 400/404/500 generic error text
func (h *TowHandler) PostTow(c *gin.Context) {
	schedulingLink := c.Param("schedulingLink")
	if schedulingLink == "" {
		message := "schedulingLink is required"
		log.Println(message)
		c.String(http.StatusBadRequest, message)
		return
	}

	var towBody model.Tow
	if err := c.ShouldBindJSON(&towBody); err != nil {
		message := "invalid JSON towBody"
		log.Println(message)
		c.String(http.StatusBadRequest, message)
		return
	}

	tow, err := h.towService.ScheduleTow(c.Request.Context(), &towBody, schedulingLink)
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, "something went wrong")
		return
	}

	c.JSON(http.StatusCreated, tow)
}

// PutUpdateTow PUT /tows/:towId
// Partially updates a tow by ID.
// Request: partial Tow fields in JSON body
// Response: 204 | 400/404/500 generic error text
func (h *TowHandler) PutUpdateTow(c *gin.Context) {
	towId := c.Param("towId")
	if towId == "" {
		c.String(http.StatusBadRequest, "tow id is required")
		return
	}

	var body model.Tow
	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.towService.UpdateTow(c.Request.Context(), towId, &body); err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, "something went wrong")
		return
	}

	c.Status(http.StatusNoContent)
}

// GetEstimate GET /tow/estimates?pickup=123&dropoff=456&company=0009990
// Calculates and returns a price estimate for a tow request.
// Query parameters: pickup (required), dropoff (required), company (required)
// Response: 200 { "estimate": int } | 400/404/500 generic error text
func (h *TowHandler) GetEstimate(c *gin.Context) {
	pickup := c.Query("pickup")
	dropoff := c.Query("dropoff")
	company := c.Query("company")

	if pickup == "" {
		c.String(http.StatusBadRequest, "pickup is required")
		return
	}
	if dropoff == "" {
		c.String(http.StatusBadRequest, "dropoff is required")
		return
	}
	if company == "" {
		c.String(http.StatusBadRequest, "company is required")
		return
	}

	estimate, err := h.towService.GetEstimate(c.Request.Context(), company, pickup, dropoff)
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, gin.H{"estimate": estimate})
}
