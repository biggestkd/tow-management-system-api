package utilities

import (
	"errors"
	"github.com/stripe/stripe-go/v83/loginlink"
	"os"

	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/account"
	"github.com/stripe/stripe-go/v83/accountlink"
)

type StripeUtility struct {
	client *stripe.Client
}

// NewStripeClient initializes stripe client using STRIPE_API_KEY and returns a singleton.
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
