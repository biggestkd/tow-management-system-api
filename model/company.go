package model

type Company struct {
	ID              *string `json:"id,omitempty" bson:"_id,omitempty"`
	Website         *string `json:"website,omitempty" bson:"website,omitempty"`
	Name            *string `json:"name,omitempty" bson:"name,omitempty"`
	Status          *string `json:"status,omitempty" bson:"status,omitempty"`
	Street          *string `json:"street,omitempty" bson:"street,omitempty"`
	City            *string `json:"city,omitempty" bson:"city,omitempty"`
	ZipCode         *string `json:"zipCode,omitempty" bson:"zipCode,omitempty"`
	State           *string `json:"state,omitempty" bson:"state,omitempty"`
	PhoneNumber     *string `json:"phoneNumber,omitempty" bson:"phoneNumber,omitempty"`
	CreatedDate     int64   `json:"createdDate,omitempty" bson:"createdDate,omitempty"`
	SchedulingLink  *string `json:"schedulingLink,omitempty" bson:"schedulingLink,omitempty"`
	StripeAccountId *string `json:"stripeAccountId,omitempty" bson:"stripeAccountId,omitempty"`
}
