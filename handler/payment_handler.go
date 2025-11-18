package handler

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v83"
)

// PaymentService defines the contract required for payment-related operations.
type PaymentService interface {
	RetrievePaymentAccount(ctx context.Context, companyId string) (*stripe.Account, error)
	GenerateDashboardLink(ctx context.Context, companyId, returnURL, refreshURL string) (string, error)
}

// PaymentHandler handles payment-related HTTP endpoints.
type PaymentHandler struct {
	paymentService PaymentService
}

// NewPaymentHandler creates a new PaymentHandler instance.
func NewPaymentHandler(service PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: service}
}

// GetPaymentAccount GET /payments/account/:companyId
// Response: 200 Stripe Account | 400 invalid request | 404 not found
func (h *PaymentHandler) GetPaymentAccount(c *gin.Context) {
	companyId := c.Param("companyId")
	if companyId == "" {
		c.String(http.StatusBadRequest, "company id is required")
		return
	}

	account, err := h.paymentService.RetrievePaymentAccount(c.Request.Context(), companyId)
	if err != nil {
		log.Println(err.Error())
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "not found") {
			c.String(http.StatusNotFound, "company not found")
			return
		}
		c.String(http.StatusBadRequest, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, account)
}

// PostPaymentAccount POST /payments/account/:companyId
// Request Body: { "returnURL": "...", "refreshURL": "..." }
// Response: 200 { "url": "..." } | 400 invalid request | 404 not found
func (h *PaymentHandler) PostPaymentAccount(c *gin.Context) {
	companyId := c.Param("companyId")
	if companyId == "" {
		c.String(http.StatusBadRequest, "company id is required")
		return
	}

	var body struct {
		ReturnURL  string `json:"returnURL" binding:"required"`
		RefreshURL string `json:"refreshURL" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.String(http.StatusBadRequest, "invalid JSON body")
		return
	}

	url, err := h.paymentService.GenerateDashboardLink(c.Request.Context(), companyId, body.ReturnURL, body.RefreshURL)
	if err != nil {
		log.Println(err.Error())
		// Check if it's a "not found" error
		if strings.Contains(err.Error(), "not found") {
			c.String(http.StatusNotFound, "company not found")
			return
		}
		c.String(http.StatusBadRequest, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
