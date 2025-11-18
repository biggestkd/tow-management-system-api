package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
	"tow-management-system-api/model"
	"tow-management-system-api/utilities"
)

type TowRepository interface {
	Create(ctx context.Context, item *model.Tow) error
	Find(ctx context.Context, filterModel *model.Tow) ([]*model.Tow, error)
	Update(ctx context.Context, id string, updateData *model.Tow) error
}

type PriceRepositoryForTowService interface {
	Find(ctx context.Context, filterModel *model.Price) ([]*model.Price, error)
}

// TowService defines business logic for the Tow entity.
type TowService struct {
	towRepository   TowRepository
	priceRepository PriceRepositoryForTowService
}

// NewTowService creates a new TowService instance.
func NewTowService(towRepo TowRepository, priceRepo PriceRepositoryForTowService) *TowService {
	return &TowService{
		towRepository:   towRepo,
		priceRepository: priceRepo,
	}
}

// ScheduleTow calculates pricing, creates a payable invoice, persists the tow, and returns the saved entity.
func (s *TowService) ScheduleTow(ctx context.Context, towRequest *model.Tow) (*model.Tow, error) {
	if towRequest == nil {
		return nil, fmt.Errorf("tow request is required")
	}

	if towRequest.CompanyID == nil || *towRequest.CompanyID == "" {
		return nil, fmt.Errorf("company id is required")
	}

	pricingInfo, err := s.priceRepository.Find(ctx, &model.Price{
		CompanyID: towRequest.CompanyID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load pricing information: %w", err)
	}

	total, err := utilities.CalculateTowPrice(pricingInfo, towRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to calculate tow price: %w", err)
	}

	towRequest.Price = &total

	invoiceID, err := utilities.CreatePayableItem(total)

	if err != nil {
		return nil, fmt.Errorf("failed to create payable item: %w", err)
	}

	paymentStatus := "unpaid"
	towRequest.PaymentStatus = &paymentStatus
	towRequest.PaymentReference = &invoiceID
	now := time.Now().UTC().Unix()
	towRequest.CreatedAt = &now
	id := uuid.NewString()
	towRequest.ID = &id

	if err := s.towRepository.Create(ctx, towRequest); err != nil {
		return nil, fmt.Errorf("failed to save tow: %w", err)
	}

	return towRequest, nil
}

// FindTowsByCompanyId retrieves all tows that belong to a specific company.
func (s *TowService) FindTowsByCompanyId(ctx context.Context, companyId string) ([]*model.Tow, error) {
	if companyId == "" {
		return nil, fmt.Errorf("company id is required")
	}

	tows, err := s.towRepository.Find(ctx, &model.Tow{
		CompanyID: &companyId,
	})

	if err != nil {
		return nil, fmt.Errorf("find tows failed: %w", err)
	}

	return tows, nil
}

// UpdateTow updates a tow by its ID with the provided partial fields.
func (s *TowService) UpdateTow(ctx context.Context, towId string, update *model.Tow) error {
	if towId == "" {
		return fmt.Errorf("tow id is required")
	}
	if update == nil {
		return fmt.Errorf("update body is required")
	}

	if err := s.towRepository.Update(ctx, towId, update); err != nil {
		return fmt.Errorf("update tow failed: %w", err)
	}
	return nil
}
