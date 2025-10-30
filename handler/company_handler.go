package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"tow-management-system-api/model"
)

type CompanyService interface {
	CreateCompany(ctx context.Context, company *model.Company) (*model.Company, error)
	FindCompanyById(ctx context.Context, id string) (*model.Company, error)
}

type CompanyHandler struct {
	companyService CompanyService
}

func NewCompanyHandler(repo CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: repo}
}

// PostCompany POST /company
// Request Body: { "website": "...", "phone": "...", "name": "Company Name" }
// Response: 201 Company | 400 generic error text
func (h *CompanyHandler) PostCompany(ginContext *gin.Context) {
	var body model.Company
	if err := ginContext.ShouldBindJSON(&body); err != nil {
		ginContext.String(http.StatusBadRequest, "Something went wrong")
		return
	}

	created, err := h.companyService.CreateCompany(ginContext.Request.Context(), &body)
	if err != nil || created == nil {
		ginContext.String(http.StatusBadRequest, "Something went wrong")
		return
	}

	ginContext.JSON(http.StatusCreated, created)
}

// GetCompany GET /company/:id
// Response: 200 Company | 400 generic error text
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "company id is required")
		return
	}

	// Validate UUID format (your service is creating UUIDs for Company.ID)
	if _, err := uuid.Parse(id); err != nil {
		c.String(http.StatusBadRequest, "invalid company id")
		return
	}

	company, err := h.companyService.FindCompanyById(c.Request.Context(), id)

	if err != nil {
		// Service returns an error when not found; map to 404 for API clarity
		if err.Error() == "company not found" {
			c.String(http.StatusNotFound, "company not found")
			return
		}
		c.String(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, company)
}
