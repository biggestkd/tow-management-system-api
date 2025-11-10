package utilities

import (
	"fmt"
	"tow-management-system-api/model"
)

// CalculateTowPrice determines the total cost of a tow based on pricing information and the tow request details.
// TODO: implement pricing logic.
func CalculateTowPrice(pricingInfo []*model.Price, towRequest *model.Tow) (int, error) {
	return 0, fmt.Errorf("calculate tow price not implemented")
}

// CreatePayableItem creates a payable item (e.g., invoice) for the provided total amount.
// TODO: integrate with payment processing (e.g., Stripe).
func CreatePayableItem(total int) (string, error) {
	return "", fmt.Errorf("create payable item not implemented")
}
