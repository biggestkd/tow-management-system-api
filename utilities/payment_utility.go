package utilities

import (
	"github.com/google/uuid"
	"tow-management-system-api/model"
)

// CalculateTowPrice determines the total cost of a tow based on pricing information and the tow request details.
// Returns the calculated total in cents.
// TODO: implement pricing logic.
func CalculateTowPrice(pricingInfo []*model.Price, towRequest *model.Tow) (int, error) {
	//return 0, fmt.Errorf("calculate tow price not implemented")
	return 100, nil
}

// CreatePayableItem creates a payable item (e.g., invoice) for the provided total amount.
// Returns a reference to the payable item.
// TODO: integrate with payment processing (e.g., Stripe).
func CreatePayableItem(total int) (string, error) {
	//return "", fmt.Errorf("create payable item not implemented")
	return uuid.NewString(), nil
}
