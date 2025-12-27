package utilities

import (
	"errors"
	"fmt"
	"github.com/stripe/stripe-go/v83/loginlink"
	"os"
	"tow-management-system-api/model"

	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/account"
	"github.com/stripe/stripe-go/v83/accountlink"
	checkoutsession "github.com/stripe/stripe-go/v83/checkout/session"
)

type StripeUtility struct {
	client *stripe.Client
}

// NewStripeClient initializes stripe geoplacesClient using STRIPE_API_KEY and returns a singleton.
func NewStripeClient() (*StripeUtility, error) {
	apiKey := os.Getenv("STRIPE_API_KEY")

	if apiKey == "" {
		return nil, errors.New("STRIPE_API_KEY not set")
	}

	sc := stripe.NewClient(apiKey)

	stripe.Key = apiKey

	return &StripeUtility{client: sc}, nil
}

// CreateConnectedAccount creates a Stripe connected account and returns an onboarding URL.
// If accountID is provided and not empty, it will use the existing account; otherwise, it creates a new one.
// The returned URL allows service providers to enter their identity and banking information.
func (sc *StripeUtility) CreateConnectedAccount() (string, error) {

	params := &stripe.AccountParams{
		Country: stripe.String("US"),
		Controller: &stripe.AccountControllerParams{
			Fees: &stripe.AccountControllerFeesParams{
				Payer: stripe.String(stripe.AccountControllerFeesPayerApplication),
			},
			Losses: &stripe.AccountControllerLossesParams{
				Payments: stripe.String(stripe.AccountControllerLossesPaymentsApplication),
			},
			StripeDashboard: &stripe.AccountControllerStripeDashboardParams{
				Type: stripe.String(stripe.AccountControllerStripeDashboardTypeExpress),
			},
		},
	}

	account, err := account.New(params)

	if err != nil {
		return "", errors.New("An error occurred when calling the Stripe API to create an account: " + err.Error())
	}

	return account.ID, nil

}

// CreateLoginLink creates a login link for the connected account.
// Returns a URL that allows the account holder to access their Stripe Express dashboard.
func (sc *StripeUtility) CreateLoginLink(accountId string) (string, error) {
	if accountId == "" {
		return "", errors.New("accountId is required")
	}

	params := &stripe.LoginLinkParams{
		Account: stripe.String(accountId),
	}

	loginLink, err := loginlink.New(params)
	if err != nil {
		return "", errors.New("failed to create login link: " + err.Error())
	}

	return loginLink.URL, nil
}

// CreateAccountLink creates an account link for onboarding or updating account information.
// Returns a URL that allows the account holder to complete their account setup.
func (sc *StripeUtility) CreateAccountLink(accountId, returnURL, refreshURL string) (string, error) {
	if accountId == "" {
		return "", errors.New("accountId is required")
	}
	if returnURL == "" {
		return "", errors.New("returnURL is required")
	}
	if refreshURL == "" {
		return "", errors.New("refreshURL is required")
	}

	params := &stripe.AccountLinkParams{
		Account:    stripe.String(accountId),
		RefreshURL: stripe.String(refreshURL),
		ReturnURL:  stripe.String(returnURL),
		Type:       stripe.String(string(stripe.AccountLinkTypeAccountOnboarding)),
	}

	accountLink, err := accountlink.New(params)
	if err != nil {
		return "", errors.New("failed to create account link: " + err.Error())
	}

	return accountLink.URL, nil
}

// GetAccount retrieves a Stripe connected account by its ID.
// Returns the full account object with all account details.
func (sc *StripeUtility) GetAccount(accountId string) (*stripe.Account, error) {
	if accountId == "" {
		return nil, errors.New("accountId is required")
	}

	acct, err := account.GetByID(accountId, nil)
	if err != nil {
		return nil, errors.New("failed to retrieve account: " + err.Error())
	}

	return acct, nil
}

// CreatePayableItem creates a Stripe Checkout Session (one-time payment) and returns the hosted Checkout URL.
//
// Parameters:
// - total: total amount in cents (integer)
// - lineItems: array of (name, amount) pairs, amounts in cents
//
// Returns:
// - URL string for Stripe-hosted checkout
func (sc *StripeUtility) CreatePayableItem(total int64, lineItems []model.PayableLineItem) (string, error) {
	if total <= 0 {
		return "", errors.New("total must be greater than 0")
	}
	if len(lineItems) == 0 {
		return "", errors.New("at least one line item is required")
	}

	successURL := "https://www.google.com/"
	cancelURL := "https://www.google.com/"
	if successURL == "" {
		return "", errors.New("STRIPE_CHECKOUT_SUCCESS_URL not set")
	}
	if cancelURL == "" {
		return "", errors.New("STRIPE_CHECKOUT_CANCEL_URL not set")
	}

	sessionLineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(lineItems))

	for i, li := range lineItems {
		if li.Name == "" {
			return "", fmt.Errorf("lineItems[%d].Name is required", i)
		}
		if li.Amount <= 0 {
			return "", fmt.Errorf("lineItems[%d].Amount must be > 0", i)
		}

		sessionLineItems = append(sessionLineItems, &stripe.CheckoutSessionLineItemParams{
			Quantity: stripe.Int64(1), // the number of miles or the one per mile amount
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(string(stripe.CurrencyUSD)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(li.Name),
				},
				UnitAmount: stripe.Int64(li.Amount),
			},
		})
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		LineItems:  sessionLineItems,
		Currency:   stripe.String("USD"),
	}

	// Best-effort idempotency key to avoid duplicates if caller retries.
	// If you have a stable internal ID (e.g., service_request_id), pass it via metadata and use it here instead.
	params.SetIdempotencyKey(fmt.Sprintf("payable_%d_%d", total, len(lineItems)))

	sess, err := checkoutsession.New(params)
	if err != nil {
		return "", errors.New("failed to create checkout session: " + err.Error())
	}

	if sess.URL == "" {
		return "", errors.New("checkout session created but no URL was returned")
	}

	return sess.URL, nil
}
