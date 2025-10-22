package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"tow-management-system-api/model"
)

type CompanyService interface {
	CreateCompany(ctx context.Context, company *model.Company) (*model.Company, error)
}

type CompanyHandler struct {
	companyService CompanyService
}

func NewCompanyHandler(repo CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService: repo}
}

// PostCompany POST /company
// Request: { "website": "...", "phone": "...", "name": "Company Name" }
// Requires userId (for now via header X-User-Id; replace with Cognito JWT later)
// Response: 201 Company | 400 generic error text
func (h *CompanyHandler) PostCompany(context *gin.Context) {
	var body model.Company
	if err := context.ShouldBindJSON(&body); err != nil {
		context.String(http.StatusBadRequest, "Something went wrong")
		return
	}
	
	created, err := h.companyService.CreateCompany(context, &body)
	if err != nil || created == nil {
		context.String(http.StatusBadRequest, "Something went wrong")
		return
	}

	context.JSON(http.StatusCreated, created)
}
