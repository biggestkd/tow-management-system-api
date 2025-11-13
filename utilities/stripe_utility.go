package utilities

import (
	"errors"
	"os"
	"time"

	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/account"
	"github.com/stripe/stripe-go/v83/accountlink"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/invoice"
	"github.com/stripe/stripe-go/v83/invoiceitem"
)

type StripeClient struct {
	client *stripe.Client
}

// InvoiceItem represents an item to be added to an invoice
type InvoiceItem struct {
	Description string
	Amount      int64 // Amount in cents
	Quantity    int64 // Quantity of the item
}

// NewStripeClient initializes stripe client using STRIPE_API_KEY and returns a singleton.
func NewStripeClient() (*StripeClient, error) {
	apiKey := os.Getenv("STRIPE_API_KEY")

	if apiKey == "" {
		return nil, errors.New("STRIPE_API_KEY not set")
	}

	sc := stripe.NewClient(apiKey)

	return &StripeClient{client: sc}, nil
}

// CreateConnectedAccount creates a Stripe connected account and returns an onboarding URL.
// If accountID is provided and not empty, it will use the existing account; otherwise, it creates a new one.
// The returned URL allows service providers to enter their identity and banking information.
func (sc *StripeClient) CreateConnectedAccount(returnURL string, refreshURL string) (string, string, error) {

	// Create a new connected account
	params := &stripe.AccountParams{
		Type:    stripe.String(string(stripe.AccountTypeExpress)),
		Country: stripe.String("US"), // Default to US, can be made configurable
	}

	acc, err := account.New(params)
	if err != nil {
		return "", "", errors.New("failed to create connected account: " + err.Error())
	}
	connectedAccountID := acc.ID

	// Create an account link for onboarding
	linkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(connectedAccountID),
		RefreshURL: stripe.String(refreshURL),
		ReturnURL:  stripe.String(returnURL),
		Type:       stripe.String(string(stripe.AccountLinkTypeAccountOnboarding)),
	}

	link, err := accountlink.New(linkParams)
	if err != nil {
		return "", "", errors.New("failed to create account link: " + err.Error())
	}

	return link.URL, connectedAccountID, nil
}

// UpdateConnectedAccount creates a URL from Stripe that service providers can access to update their identity/banking info.
// The stripeAccountId is required as this function is for updating existing connected accounts.
func (sc *StripeClient) UpdateConnectedAccount(stripeAccountId string, returnURL string, refreshURL string) (string, error) {
	if stripeAccountId == "" {
		return "", errors.New("stripeAccountId is required for updating connected account")
	}

	// Verify the account exists
	_, err := account.GetByID(stripeAccountId, nil)
	if err != nil {
		return "", errors.New("failed to retrieve connected account: " + err.Error())
	}

	// Create an account link for updating account information
	linkParams := &stripe.AccountLinkParams{
		Account:    stripe.String(stripeAccountId),
		RefreshURL: stripe.String(refreshURL),
		ReturnURL:  stripe.String(returnURL),
		Type:       stripe.String(string(stripe.AccountLinkTypeAccountUpdate)),
	}

	link, err := accountlink.New(linkParams)
	if err != nil {
		return "", errors.New("failed to create account link: " + err.Error())
	}

	return link.URL, nil
}

// CreateCustomerOnConnectedAccount creates a customer attached to the connected account.
// Takes the required customer information and stripe account id, and returns a customer reference
// that can be used to create an invoice or a charge.
func (sc *StripeClient) CreateCustomerOnConnectedAccount(stripeAccountId string, email string, name string, phone string) (string, error) {
	if stripeAccountId == "" {
		return "", errors.New("stripeAccountId is required")
	}
	if email == "" {
		return "", errors.New("email is required")
	}

	// Verify the connected account exists
	_, err := account.GetByID(stripeAccountId, nil)
	if err != nil {
		return "", errors.New("failed to retrieve connected account: " + err.Error())
	}

	// Create customer parameters
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}

	if name != "" {
		params.Name = stripe.String(name)
	}

	if phone != "" {
		params.Phone = stripe.String(phone)
	}

	// Set the account to create the customer on the connected account
	params.SetStripeAccount(stripeAccountId)

	// Create the customer on the connected account
	customer, err := customer.New(params)
	if err != nil {
		return "", errors.New("failed to create customer on connected account: " + err.Error())
	}

	return customer.ID, nil
}

// CreateInvoice creates an invoice with the given total, items, and customer reference.
// Generates invoice items and returns the URL to pay the invoice.
func (sc *StripeClient) CreateInvoice(stripeAccountId string, customerId string, total int64, items []InvoiceItem) (string, error) {

	// Create the invoice

	now := time.Now().Unix()
	delay := 30 * time.Second
	finalizeAt := now + delay.Microseconds()

	invoiceParams := &stripe.InvoiceParams{
		Customer:                 stripe.String(customerId),
		AutoAdvance:              stripe.Bool(false), // Don't auto-finalize, we'll do it manually
		AutomaticallyFinalizesAt: stripe.Int64(finalizeAt),
	}

	// Attach invoice to connected account
	invoiceParams.SetStripeAccount(stripeAccountId)

	invoice, err := invoice.New(invoiceParams)

	if err != nil {
		return "", errors.New("failed to create invoice: " + err.Error())
	}

	// Add items to invoice
	for _, item := range items {

		itemParams := &stripe.InvoiceItemParams{
			Customer:    stripe.String(customerId),
			Amount:      stripe.Int64(item.Amount),
			Description: stripe.String(item.Description),
			Invoice:     stripe.String(invoice.ID),
			Quantity:    stripe.Int64(item.Quantity),
		}

		// Set the account to create the invoice item on the connected account
		itemParams.SetStripeAccount(stripeAccountId)

		_, err := invoiceitem.New(itemParams)

		if err != nil {
			return "", errors.New("failed to create invoice item: " + err.Error())
		}
	}

	if err != nil {
		return "", errors.New("failed to finalize invoice: " + err.Error())
	}

	// Return the hosted invoice URL for payment
	if invoice.HostedInvoiceURL == "" {
		return "", errors.New("invoice created but no payment URL available")
	}

	return invoice.HostedInvoiceURL, nil
}
