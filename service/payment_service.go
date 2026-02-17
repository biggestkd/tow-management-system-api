package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"tow-management-system-api/model"
	"tow-management-system-api/utilities"

	"github.com/stripe/stripe-go/v83"
)

// TowDataRepository defines the interface for tow data operations.
type TowDataRepository interface {
	Find(ctx context.Context, filterModel *model.Tow) ([]*model.Tow, error)
	Update(ctx context.Context, id string, updateData *model.Tow) error
}

// PaymentService encapsulates payment-related business logic.
type PaymentService struct {
	towDataRepository TowDataRepository
	companyRepository CompanyRepository
	stripeClient      *utilities.StripeUtility
}

// NewPaymentService constructs a PaymentService with the provided dependencies.
func NewPaymentService(towDataRepo TowDataRepository, companyRepo CompanyRepository, stripeClient *utilities.StripeUtility) *PaymentService {
	return &PaymentService{
		towDataRepository: towDataRepo,
		companyRepository: companyRepo,
		stripeClient:      stripeClient,
	}
}

// RetrievePaymentAccount retrieves the Stripe account associated with a company.
// First finds the company by ID, then retrieves the Stripe account using the company's StripeAccountId.
func (s *PaymentService) RetrievePaymentAccount(ctx context.Context, companyId string) (*stripe.Account, error) {
	if companyId == "" {
		return nil, fmt.Errorf("company id is required")
	}

	// Find the company
	companies, err := s.companyRepository.Find(ctx, &model.Company{ID: &companyId})
	if err != nil {
		return nil, fmt.Errorf("failed to find company: %w", err)
	}

	if len(companies) == 0 {
		return nil, fmt.Errorf("company not found")
	}

	company := companies[0]

	// Check if company has a Stripe account ID
	if company.StripeAccountId == nil || *company.StripeAccountId == "" {
		return nil, fmt.Errorf("company does not have a stripe account id")
	}

	// Get the Stripe account
	account, err := s.stripeClient.GetAccount(*company.StripeAccountId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve stripe account: %w", err)
	}

	return account, nil
}

// GenerateDashboardLink generates a dashboard link for a company.
// If the account details are already submitted, returns a login link.
// Otherwise, returns an account onboarding link.
func (s *PaymentService) GenerateDashboardLink(ctx context.Context, companyId, returnURL, refreshURL string) (string, error) {
	if companyId == "" {
		return "", fmt.Errorf("company id is required")
	}

	// Find the company
	companies, err := s.companyRepository.Find(ctx, &model.Company{ID: &companyId})
	if err != nil {
		return "", fmt.Errorf("failed to find company: %w", err)
	}

	if len(companies) == 0 {
		return "", fmt.Errorf("company not found")
	}

	company := companies[0]

	// Check if company has a Stripe account ID
	if company.StripeAccountId == nil || *company.StripeAccountId == "" {
		return "", fmt.Errorf("company does not have a stripe account id")
	}

	// Get the Stripe account
	account, err := s.stripeClient.GetAccount(*company.StripeAccountId)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve stripe account: %w", err)
	}

	// If details are submitted, return a login link
	if account.DetailsSubmitted {
		loginLink, err := s.stripeClient.CreateLoginLink(account.ID)
		if err != nil {
			return "", fmt.Errorf("failed to create login link: %w", err)
		}
		return loginLink, nil
	}

	// Otherwise, return an account link for onboarding
	if returnURL == "" || refreshURL == "" {
		return "", fmt.Errorf("returnURL and refreshURL are required for account onboarding")
	}

	accountLink, err := s.stripeClient.CreateAccountLink(account.ID, returnURL, refreshURL)
	if err != nil {
		return "", fmt.Errorf("failed to create account link: %w", err)
	}

	return accountLink, nil
}

// ProcessCheckoutSuccessEvent processes a Stripe checkout session success event.
// It finds the tow associated with the checkout session and updates its payment status to "paid".
func (s *PaymentService) ProcessCheckoutSuccessEvent(ctx context.Context, event *stripe.Event) error {
	if event == nil {
		return fmt.Errorf("event is required")
	}

	// Check if the event type begins with "checkout.session"
	eventType := string(event.Type)
	if !strings.HasPrefix(eventType, "checkout.session.async_payment_succeeded") {
		return fmt.Errorf("event type %s does not begin with 'checkout.session'", eventType)
	}

	// Extract the checkout session ID from event.data.object.id
	// The event.Data.Raw contains the full event data with nested object
	var checkoutSession stripe.CheckoutSession

	if err := json.Unmarshal(event.Data.Raw, &checkoutSession); err != nil {
		return fmt.Errorf("failed to unmarshal checkout session from event data: %w", err)
	}

	if checkoutSession.ID == "" {
		return fmt.Errorf("checkout session ID is empty")
	}

	// Search for the tow with PaymentReference == checkout session ID
	tow, err := s.towDataRepository.Find(ctx, &model.Tow{PaymentReference: &checkoutSession.ID})
	if err != nil {
		return fmt.Errorf("failed to find tow with payment reference %s: %w", checkoutSession.ID, err)
	}

	if len(tow) == 0 {
		return fmt.Errorf("tow not found with payment reference %s", checkoutSession.ID)
	}

	if tow[0].ID == nil || *tow[0].ID == "" {
		return fmt.Errorf("tow ID is empty")
	}

	// Update the tow to set PaymentStatus == "paid"
	paidStatus := "paid"
	if err := s.towDataRepository.Update(ctx, *tow[0].ID, &model.Tow{PaymentStatus: &paidStatus}); err != nil {
		return fmt.Errorf("failed to update tow payment status: %w", err)
	}

	return nil
}
