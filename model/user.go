package model

type User struct {
	ID          *string `json:"id,omitempty" bson:"_id,omitempty"`
	CompanyID   *string `json:"companyId,omitempty" bson:"companyId,omitempty"`
	CreatedDate int64   `json:"createdDate,omitempty" bson:"createdDate,omitempty"`
	FirstName   *string `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName    *string `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Phone       *string `json:"phone,omitempty" bson:"phone,omitempty"`
	Email       *string `json:"email,omitempty" bson:"email,omitempty"`
}
