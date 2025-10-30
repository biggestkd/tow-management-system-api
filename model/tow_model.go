package model

type Tow struct {
	ID             *string  `json:"id,omitempty" bson:"_id,omitempty"`
	Destination    *string  `json:"destination,omitempty" bson:"destination,omitempty"`
	Pickup         *string  `json:"pickup,omitempty" bson:"pickup,omitempty"`
	Vehicle        *string  `json:"vehicle,omitempty" bson:"vehicle,omitempty"`
	PrimaryContact *string  `json:"primaryContact,omitempty" bson:"primaryContact,omitempty"`
	Attachments    []string `json:"attachments,omitempty" bson:"attachments,omitempty"`
	Notes          *string  `json:"notes,omitempty" bson:"notes,omitempty"`
	History        []string `json:"history,omitempty" bson:"history,omitempty"`
	Status         *string  `json:"status,omitempty" bson:"status,omitempty"` // pending, accepted, dispatched, arrived_pickup, in_transit, completed, cancelled
	CompanyID      *string  `json:"companyId,omitempty" bson:"companyId,omitempty"`
	CreatedAt      *int64   `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Price          *int     `json:"price,omitempty" bson:"price,omitempty"`
}
