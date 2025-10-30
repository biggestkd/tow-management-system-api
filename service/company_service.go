package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strings"
	"time"
	"tow-management-system-api/model"
)

type CompanyRepository interface {
	Create(ctx context.Context, item *model.Company) error
	Find(ctx context.Context, filterModel *model.Company) ([]*model.Company, error)
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

	// Ensure ID (string pointer), schedulingLink and createdDate (int64) are set
	id := uuid.NewString()
	company.ID = &id
	company.CreatedDate = time.Now().UTC().Unix()
	company.SchedulingLink = generateSchedulingLinkSlug(company.Name)

	if err := s.companyRepository.Create(ctx, company); err != nil {
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
