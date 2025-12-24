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
	locationUtility *utilities.LocationUtility
}

// NewTowService creates a new TowService instance.
func NewTowService(towRepo TowRepository, priceRepo PriceRepositoryForTowService, locationUtility *utilities.LocationUtility) *TowService {
	return &TowService{
		towRepository:   towRepo,
		priceRepository: priceRepo,
		locationUtility: locationUtility,
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

// GetEstimate calculates and returns a price estimate for a tow request without creating a tow or payment reference.
func (s *TowService) GetEstimate(ctx context.Context, companyId string, pickup string, dropoff string) (int64, error) {
	if companyId == "" {
		return 0, fmt.Errorf("company id is required")
	}
	if pickup == "" {
		return 0, fmt.Errorf("pickup is required")
	}
	if dropoff == "" {
		return 0, fmt.Errorf("dropoff is required")
	}

	// load prices
	pricingInfo, err := s.priceRepository.Find(ctx, &model.Price{
		CompanyID: &companyId,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to load pricing information: %w", err)
	}

	// Convert the addresses to geo positions
	pickupCoordinates, err := s.locationUtility.ParseGeocodeFromAddress(pickup)
	if err != nil {
		return 0, fmt.Errorf("failed to parse pickup location: %w", err)
	}

	dropoffCoordinates, err := s.locationUtility.ParseGeocodeFromAddress(dropoff)
	if err != nil {
		return 0, fmt.Errorf("failed to parse dropoff location: %w", err)
	}

	totalMiles, err := s.locationUtility.CalculateDistanceBetweenCoordinates(pickupCoordinates, dropoffCoordinates)

	totalAmount := calculatePriceFromMiles(pricingInfo, totalMiles)

	if err != nil {
		return 0, fmt.Errorf("failed to calculate tow price: %w", err)
	}

	return int64(totalAmount), nil
}

// calculatePriceFromMiles calculates the total price in cents based on pricing information and total miles.
// Prices array always contains an item with "Hook Up Fee" and "Per Mile Amount".
func calculatePriceFromMiles(prices []*model.Price, totalMiles float64) float64 {
	total := 0.0

	for _, price := range prices {
		if price.ItemName == nil || price.Amount == nil {
			continue
		}

		if *price.ItemName == "Hook Up Fee" {
			total += float64(*price.Amount)
		} else if *price.ItemName == "Per Mile Amount" {
			total += float64(*price.Amount) * totalMiles
		}
	}

	return total
}
