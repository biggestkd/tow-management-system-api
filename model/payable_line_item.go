package model

// PayableLineItem represents a single item the customer is paying for.
// Amount is in the smallest currency unit (e.g. USD cents).
type PayableLineItem struct {
	Name     string
	Amount   int64
	Quantity int64
}
