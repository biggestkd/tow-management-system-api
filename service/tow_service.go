package service

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"
	"tow-management-system-api/model"
	"tow-management-system-api/utilities"

	"github.com/google/uuid"
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
	stripeClient    *utilities.StripeUtility
}

// NewTowService creates a new TowService instance.
func NewTowService(towRepo TowRepository, priceRepo PriceRepositoryForTowService, locationUtility *utilities.LocationUtility, stripeClient *utilities.StripeUtility) *TowService {
	return &TowService{
		towRepository:   towRepo,
		priceRepository: priceRepo,
		locationUtility: locationUtility,
		stripeClient:    stripeClient,
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

	log.Println(*towRequest.CompanyID)

	if err != nil {
		return nil, fmt.Errorf("failed to load pricing information: %w", err)
	}

	// Calculate total miles
	pickupAddress, err := s.locationUtility.ParseGeocodeFromAddress(*towRequest.Pickup)

	if err != nil {
		return nil, err
	}

	destinationAddress, err := s.locationUtility.ParseGeocodeFromAddress(*towRequest.Destination)

	if err != nil {
		return nil, err
	}

	totalMiles, err := s.locationUtility.CalculateDistanceBetweenCoordinates(pickupAddress, destinationAddress)

	if err != nil {
		return nil, err
	}

	// generate the payable line item for the hookup fee
	var lineItems []model.PayableLineItem

	lineItems = append(lineItems, model.PayableLineItem{
		Name:     "Hook Up Fee",
		Amount:   int64(*pricingInfo[0].Amount),
		Quantity: 1,
	})

	mileageCost := int64(*pricingInfo[1].Amount) * int64(math.Ceil(totalMiles))

	lineItems = append(lineItems, model.PayableLineItem{
		Name:     fmt.Sprintf("%f miles at $%f per mile", totalMiles, float64(*pricingInfo[1].Amount/100)), // must divide by 100 to get it into dollar amounts
		Amount:   mileageCost,
		Quantity: 1,
	})

	totalFloat := calculatePriceFromMiles(pricingInfo, totalMiles)
	total := int(math.Round(totalFloat))

	towRequest.Price = &total

	checkoutURL, err := s.stripeClient.CreatePayableItem(int64(total), lineItems)

	if err != nil {
		return nil, fmt.Errorf("failed to create payable item: %w", err)
	}

	paymentStatus := "unpaid"
	towRequest.PaymentStatus = &paymentStatus
	towRequest.PaymentReference = &checkoutURL
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
