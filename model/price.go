package model

type Price struct {
	ID        *string `json:"id,omitempty" bson:"_id,omitempty"`
	ItemName  *string `json:"itemName,omitempty" bson:"itemName,omitempty"`
	Amount    *int    `json:"amount,omitempty" bson:"amount,omitempty"`
	CompanyID *string `json:"companyId,omitempty" bson:"companyId,omitempty"`
}
