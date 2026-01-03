package model

type Vehicle struct {
	Year        *string `json:"year,omitempty" bson:"year,omitempty"`
	Make        *string `json:"make,omitempty" bson:"make,omitempty"`
	Model       *string `json:"model,omitempty" bson:"model,omitempty"`
	State       *string `json:"state,omitempty" bson:"state,omitempty"`
	PlateNumber *string `json:"plateNumber,omitempty" bson:"plateNumber,omitempty"`
}

type PrimaryContact struct {
	LastName  *string `json:"lastName,omitempty" bson:"lastName,omitempty"`
	FirstName *string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	Email     *string `json:"email,omitempty" bson:"email,omitempty"`
	Phone     *string `json:"phone,omitempty" bson:"phone,omitempty"`
}

type Tow struct {
	ID               *string         `json:"id,omitempty" bson:"_id,omitempty"`
	Destination      *string         `json:"destination,omitempty" bson:"destination,omitempty"`
	Pickup           *string         `json:"pickup,omitempty" bson:"pickup,omitempty"`
	Vehicle          *Vehicle        `json:"vehicle,omitempty" bson:"vehicle,omitempty"`
	PrimaryContact   *PrimaryContact `json:"primaryContact,omitempty" bson:"primaryContact,omitempty"`
	Attachments      []string        `json:"attachments,omitempty" bson:"attachments,omitempty"`
	Notes            *string         `json:"notes,omitempty" bson:"notes,omitempty"`
	History          []string        `json:"history,omitempty" bson:"history,omitempty"`
	Status           *string         `json:"status,omitempty" bson:"status,omitempty"`                     // pending, accepted, dispatched, arrived_pickup, in_transit, completed, cancelled
	PaymentStatus    *string         `json:"paymentStatus,omitempty" bson:"paymentStatus,omitempty"`       // unpaid, paid
	PaymentReference *string         `json:"paymentReference,omitempty" bson:"paymentReference,omitempty"` // payment reference id from stripe
	CheckoutUrl      *string         `json:"checkoutUrl,omitempty" bson:"checkoutUrl,omitempty"`           // Stripe checkout session URL
	CompanyID        *string         `json:"companyId,omitempty" bson:"companyId,omitempty"`
	CreatedAt        *int64          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Price            *int            `json:"price,omitempty" bson:"price,omitempty"`
}
