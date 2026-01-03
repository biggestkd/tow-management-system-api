package service

import (
	"context"
	"fmt"
	"math"
	"os"
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
	towRepository     TowRepository
	priceRepository   PriceRepositoryForTowService
	companyRepository CompanyRepository
	locationUtility   *utilities.LocationUtility
	stripeClient      *utilities.StripeUtility
	emailUtility      *utilities.AmazonSesUtility
}

// NewTowService creates a new TowService instance.
func NewTowService(towRepo TowRepository, priceRepo PriceRepositoryForTowService, companyRepo CompanyRepository, locationUtility *utilities.LocationUtility, stripeClient *utilities.StripeUtility, emailUtility *utilities.AmazonSesUtility) *TowService {
	return &TowService{
		towRepository:     towRepo,
		priceRepository:   priceRepo,
		companyRepository: companyRepo,
		locationUtility:   locationUtility,
		stripeClient:      stripeClient,
		emailUtility:      emailUtility,
	}
}

// ScheduleTow calculates pricing, creates a payable invoice, persists the tow, and returns the saved entity.
func (s *TowService) ScheduleTow(ctx context.Context, towRequest *model.Tow, schedulingLink string) (*model.Tow, error) {
	if towRequest == nil {
		return nil, fmt.Errorf("tow request is required")
	}

	if schedulingLink == "" {
		return nil, fmt.Errorf("schedulingLink is required")
	}

	// Fetch company information
	companies, err := s.companyRepository.Find(ctx, &model.Company{
		SchedulingLink: &schedulingLink,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch company: %w", err)
	}
	if len(companies) == 0 {
		return nil, fmt.Errorf("company not found")
	}

	pricingInfo, err := s.priceRepository.Find(ctx, &model.Price{
		CompanyID: companies[0].ID,
	})

	towRequest.CompanyID = companies[0].ID

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

	checkoutSessionId, checkoutURL, err := s.stripeClient.CreatePayableItem(int64(total), lineItems)

	if err != nil {
		return nil, fmt.Errorf("failed to create payable item: %w", err)
	}

	paymentStatus := "unpaid"
	towRequest.PaymentStatus = &paymentStatus
	towRequest.PaymentReference = &checkoutSessionId
	towRequest.CheckoutUrl = &checkoutURL
	now := time.Now().UTC().Unix()
	towRequest.CreatedAt = &now
	id := uuid.NewString()
	towRequest.ID = &id
	// TODO: update the status to check if the requestor was a driver or company
	status := "ACCEPTED"
	towRequest.Status = &status

	if err := s.towRepository.Create(ctx, towRequest); err != nil {
		return nil, fmt.Errorf("failed to save tow: %w", err)
	}

	emailContent, err := s.formatPaymentEmail(ctx, towRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to format email: %w", err)
	}

	subject := "Service Confirmation – Complete Your Payment"
	err = s.emailUtility.SendEmail(ctx, *towRequest.PrimaryContact.Email, subject, emailContent)

	if err != nil {
		return nil, err
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
func (s *TowService) GetEstimate(ctx context.Context, companySchedulingLink string, pickup string, dropoff string) (int64, error) {
	if companySchedulingLink == "" {
		return 0, fmt.Errorf("company id is required")
	}
	if pickup == "" {
		return 0, fmt.Errorf("pickup is required")
	}
	if dropoff == "" {
		return 0, fmt.Errorf("dropoff is required")
	}

	// Fetch company information
	companies, err := s.companyRepository.Find(ctx, &model.Company{
		SchedulingLink: &companySchedulingLink,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to fetch company: %w", err)
	}
	if len(companies) == 0 {
		return 0, fmt.Errorf("company not found")
	}

	// load prices
	pricingInfo, err := s.priceRepository.Find(ctx, &model.Price{
		CompanyID: companies[0].ID,
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

// formatPaymentEmail formats the payment confirmation email content using company and tow information.
func (s *TowService) formatPaymentEmail(ctx context.Context, towRequest *model.Tow) (string, error) {
	if towRequest == nil {
		return "", fmt.Errorf("tow request is required")
	}
	if towRequest.CompanyID == nil || *towRequest.CompanyID == "" {
		return "", fmt.Errorf("company id is required")
	}
	if towRequest.PrimaryContact == nil || towRequest.PrimaryContact.Email == nil || *towRequest.PrimaryContact.Email == "" {
		return "", fmt.Errorf("primary contact is required")
	}
	if towRequest.CheckoutUrl == nil || *towRequest.CheckoutUrl == "" {
		return "", fmt.Errorf("checkout URL is required")
	}

	// Fetch company information
	companies, err := s.companyRepository.Find(ctx, &model.Company{
		ID: towRequest.CompanyID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch company: %w", err)
	}
	if len(companies) == 0 {
		return "", fmt.Errorf("company not found")
	}

	company := companies[0]
	companyName := "the service provider"
	if company.Name != nil && *company.Name != "" {
		companyName = *company.Name
	}

	companyPhoneNumber := ""
	if company.PhoneNumber != nil && *company.PhoneNumber != "" {
		companyPhoneNumber = *company.PhoneNumber
	}

	// Get platform information from environment variables
	platformName := os.Getenv("PLATFORM_NAME")
	if platformName == "" {
		platformName = "Tow Management Platform"
	}

	platformSupportEmail := os.Getenv("PLATFORM_SUPPORT_EMAIL")
	if platformSupportEmail == "" {
		platformSupportEmail = "support@towmanagementplatform.com"
	}

	platformWebsite := os.Getenv("PLATFORM_WEBSITE")
	if platformWebsite == "" {
		platformWebsite = "https://towmanagementplatform.com"
	}

	// Format the email content
	emailContent := fmt.Sprintf(`Thank you for choosing %s for your service.

To complete your transaction securely, please use the link below:

%s

This payment is processed through our platform on behalf of %s to ensure a safe and reliable experience.

If you have any questions regarding your service or payment, please contact %s directly`, companyName, *towRequest.CheckoutUrl, companyName, companyName)

	if companyPhoneNumber != "" {
		emailContent += fmt.Sprintf(" at %s", companyPhoneNumber)
	}

	emailContent += fmt.Sprintf(`, or reply to this email for additional support.

Thank you for your business.

Best regards,
%s
Customer Support Team
%s
`, platformName, platformWebsite)

	return emailContent, nil
}
