package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
	"tow-management-system-api/model"
)

type CompanyRepository interface {
	Create(ctx context.Context, item *model.Company) error
}

type CompanyService struct {
	companyRepository CompanyRepository
}

func NewCompanyService(companyRepo CompanyRepository) *CompanyService {
	return &CompanyService{
		companyRepository: companyRepo,
	}
}

// CreateCompany returns success/failure (bool) per spec.
func (s *CompanyService) CreateCompany(ctx context.Context, company *model.Company) (*model.Company, error) {
	if company == nil {
		return nil, fmt.Errorf("company payload is nil")
	}

	// Ensure ID (string pointer) and createdDate (int64) are set
	id := uuid.NewString()
	company.ID = &id
	company.CreatedDate = time.Now().UTC().Unix()

	if err := s.companyRepository.Create(ctx, company); err != nil {
		return nil, fmt.Errorf("create company failed: %w", err)
	}

	return company, nil
}
