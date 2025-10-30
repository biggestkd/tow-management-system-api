package model

type Metric struct {
	CompanyID   *string `json:"companyId,omitempty" bson:"companyId,omitempty"`
	Type        *string `json:"type,omitempty" bson:"type,omitempty"`
	Value       *string `json:"value,omitempty" bson:"value,omitempty"`
	LastUpdated *int64  `json:"lastUpdated,omitempty" bson:"lastUpdated,omitempty"`
}
