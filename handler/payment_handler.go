package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PaymentService defines the contract required for payment-related operations.
type PaymentService interface {
	UpdateTowPayment(ctx context.Context, towID string, paymentStatus string, paymentReference string) error
}

// PaymentHandler handles payment-related HTTP endpoints.
type PaymentHandler struct {
	paymentService PaymentService
}

// NewPaymentHandler creates a new PaymentHandler instance.
func NewPaymentHandler(service PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: service}
}

// PostPaymentWebhook POST /payment/webhook
// Receives payment provider webhook notifications and updates tow payment status.
// Expected payload includes "invoiceId" and "paymentStatus".
func (h *PaymentHandler) PostPaymentWebhook(c *gin.Context) {
	var payload struct {
		InvoiceID     string `json:"invoiceId"`
		PaymentStatus string `json:"paymentStatus"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.String(http.StatusBadRequest, "invalid JSON body")
		return
	}

	if payload.InvoiceID == "" || payload.PaymentStatus == "" {
		c.String(http.StatusBadRequest, "invoiceId and paymentStatus are required")
		return
	}

	if err := h.paymentService.UpdateTowPayment(c.Request.Context(), payload.InvoiceID, payload.PaymentStatus, payload.InvoiceID); err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, "failed to update payment")
		return
	}

	c.Status(http.StatusNoContent)
}
