package service

import (
	"context"
	"fmt"

	"tow-management-system-api/model"
)

// PaymentService encapsulates payment-related business logic.
type PaymentService struct {
	towRepository TowRepository
}

// NewPaymentService constructs a PaymentService with the provided tow repository dependency.
func NewPaymentService(towRepo TowRepository) *PaymentService {
	return &PaymentService{
		towRepository: towRepo,
	}
}

// UpdateTowPayment updates the payment status and reference for a tow using the underlying repository.
func (s *PaymentService) UpdateTowPayment(ctx context.Context, towID string, paymentStatus string, paymentReference string) error {
	if towID == "" {
		return fmt.Errorf("tow id is required")
	}
	if paymentStatus == "" {
		return fmt.Errorf("payment status is required")
	}

	status := paymentStatus
	update := &model.Tow{
		PaymentStatus: &status,
	}

	if paymentReference != "" {
		reference := paymentReference
		update.PaymentReference = &reference
	}

	if err := s.towRepository.Update(ctx, towID, update); err != nil {
		return fmt.Errorf("failed to update tow payment: %w", err)
	}

	return nil
}
