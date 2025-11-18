package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"tow-management-system-api/model"
	"tow-management-system-api/utilities"

	"github.com/google/uuid"
)

type CompanyRepository interface {
	Create(ctx context.Context, item *model.Company) error
	Find(ctx context.Context, filterModel *model.Company) ([]*model.Company, error)
	Update(ctx context.Context, id string, updateData *model.Company) error
}

type CompanyService struct {
	companyRepository CompanyRepository
	stripeClient      *utilities.StripeUtility
}

func NewCompanyService(companyRepo CompanyRepository, stripeClient *utilities.StripeUtility) *CompanyService {
	return &CompanyService{
		companyRepository: companyRepo,
		stripeClient:      stripeClient,
	}
}

// CreateCompany returns success/failure (bool) per spec.
func (s *CompanyService) CreateCompany(ctx context.Context, company *model.Company) (*model.Company, error) {
	if company == nil {
		return nil, fmt.Errorf("company payload is nil")
	}

	// Ensure ID (string pointer), schedulingLink and createdDate (int64) are set
	id := uuid.NewString()
	company.ID = &id
	company.CreatedDate = time.Now().UTC().Unix()
	company.SchedulingLink = generateSchedulingLinkSlug(company.Name)

	account, err := s.stripeClient.CreateConnectedAccount()

	company.StripeAccountId = &account

	if err != nil {
		return nil, fmt.Errorf("create company failed: %w", err)
	}

	if err = s.companyRepository.Create(ctx, company); err != nil {
		return nil, fmt.Errorf("create company failed: %w", err)
	}

	return company, nil
}

func (s *CompanyService) FindCompanyById(ctx context.Context, id string) (*model.Company, error) {
	if id == "" {
		return nil, fmt.Errorf("company id is required")
	}

	company, err := s.companyRepository.Find(ctx, &model.Company{ID: &id})

	if err != nil {
		return nil, fmt.Errorf("find company failed: %w", err)
	}

	if company == nil {
		return nil, fmt.Errorf("company not found")
	}

	return company[0], nil
}

// UpdateCompany updates a company by its ID with the provided partial fields.
func (s *CompanyService) UpdateCompany(ctx context.Context, companyId string, update *model.Company) error {
	if companyId == "" {
		return fmt.Errorf("company id is required")
	}
	if update == nil {
		return fmt.Errorf("update body is required")
	}

	if err := s.companyRepository.Update(ctx, companyId, update); err != nil {
		return fmt.Errorf("update company failed: %w", err)
	}
	return nil
}

func generateSchedulingLinkSlug(companyName *string) *string {
	// Convert to lowercase
	slug := strings.ToLower(*companyName)

	// Replace any non-alphanumeric characters or spaces with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")

	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	return &slug
}
