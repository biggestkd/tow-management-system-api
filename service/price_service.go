package service

import (
	"context"
	"fmt"
	"tow-management-system-api/model"
	"tow-management-system-api/repository"

	"github.com/google/uuid"
)

type PriceRepository interface {
	Create(ctx context.Context, item *model.Price) error
	Find(ctx context.Context, filterModel *model.Price) ([]*model.Price, error)
	Update(ctx context.Context, id string, updateData *model.Price) error
}

// PriceService defines business logic for the Price entity.
type PriceService struct {
	priceRepository *repository.PriceMongoRepository
}

// NewPriceService creates a new PriceService instance.
func NewPriceService(priceRepo *repository.PriceMongoRepository) *PriceService {
	return &PriceService{
		priceRepository: priceRepo,
	}
}

// SetPrice creates or sets multiple prices. Returns nil if no errors.
// Uses Create if price.id is null or empty, otherwise uses Update.
func (s *PriceService) SetPrice(ctx context.Context, prices []*model.Price) error {
	if prices == nil {
		return fmt.Errorf("prices list is required")
	}

	for _, price := range prices {
		// Check if price has an ID - if nil or empty, create; otherwise update
		if price.ID == nil || *price.ID == "" {
			// Generate UUID for new prices
			id := uuid.NewString()
			price.ID = &id
			if err := s.priceRepository.Create(ctx, price); err != nil {
				return fmt.Errorf("failed to create price: %w", err)
			}
		} else {
			if err := s.priceRepository.Update(ctx, *price.ID, price); err != nil {
				return fmt.Errorf("failed to update price: %w", err)
			}
		}
	}

	return nil
}

// FindPricesByCompanyId retrieves all prices that belong to a specific company.
// Returns the prices as a slice or nil and any errors.
func (s *PriceService) FindPricesByCompanyId(ctx context.Context, companyId string) ([]*model.Price, error) {
	if companyId == "" {
		return nil, fmt.Errorf("company id is required")
	}

	prices, err := s.priceRepository.Find(ctx, &model.Price{
		CompanyID: &companyId,
	})

	if err != nil {
		return nil, fmt.Errorf("find prices failed: %w", err)
	}

	return prices, nil
}
