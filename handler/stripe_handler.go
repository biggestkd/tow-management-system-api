package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v83"
)

// StripePaymentService defines the contract required for Stripe webhook processing.
type StripePaymentService interface {
	ProcessCheckoutSuccessEvent(ctx context.Context, event *stripe.Event) error
}

// StripeHandler handles Stripe webhook endpoints.
type StripeHandler struct {
	paymentService StripePaymentService
}

// NewStripeHandler creates a new StripeHandler instance.
func NewStripeHandler(service StripePaymentService) *StripeHandler {
	return &StripeHandler{
		paymentService: service,
	}
}

// PostWebhook POST /webhooks/stripe
// Handles Stripe webhook events, processes checkout session success events.
// Response: 200 OK | 400 Bad Request | 500 Internal Server Error
func (h *StripeHandler) PostWebhook(c *gin.Context) {
	// Read the request body
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read webhook body: %v\n", err)
		c.String(http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Unmarshal the payload into a stripe.Event
	event := stripe.Event{}
	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("Failed to parse webhook body json: %v\n", err)
		c.String(http.StatusBadRequest, "Failed to parse webhook body")
		return
	}

	// Process the checkout success event
	if err := h.paymentService.ProcessCheckoutSuccessEvent(c.Request.Context(), &event); err != nil {
		log.Printf("Failed to process checkout success event: %v\n", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to process event: %v", err))
		return
	}

	c.String(http.StatusOK, "Event processed successfully")
}
