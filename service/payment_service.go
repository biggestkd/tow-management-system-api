package service

import (
	"context"
	"fmt"

	"tow-management-system-api/model"
	"tow-management-system-api/utilities"

	"github.com/stripe/stripe-go/v83"
)

// PaymentService encapsulates payment-related business logic.
type PaymentService struct {
	towRepository     TowRepository
	companyRepository CompanyRepository
	stripeClient      *utilities.StripeUtility
}

// NewPaymentService constructs a PaymentService with the provided dependencies.
func NewPaymentService(towRepo TowRepository, companyRepo CompanyRepository, stripeClient *utilities.StripeUtility) *PaymentService {
	return &PaymentService{
		towRepository:     towRepo,
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

//
