package service

import (
	"context"
	"fmt"
	"tow-management-system-api/model"
	"tow-management-system-api/repository"
)

type TowRepository interface {
	Create(ctx context.Context, item *model.Tow) error
	Find(ctx context.Context, filterModel *model.Tow) ([]*model.Tow, error)
}

// TowService defines business logic for the Tow entity.
type TowService struct {
	towRepository *repository.TowMongoRepository
}

// NewTowService creates a new TowService instance.
func NewTowService(towRepo *repository.TowMongoRepository) *TowService {
	return &TowService{
		towRepository: towRepo,
	}
}

// FindTowsByCompanyId retrieves all tows that belong to a specific company.
func (s *TowService) FindTowsByCompanyId(ctx context.Context, companyId string) ([]*model.Tow, error) {
	if companyId == "" {
		return nil, fmt.Errorf("company id is required")
	}

	tows, err := s.towRepository.Find(ctx, &model.Tow{
		CompanyID: &companyId,
	})

	if err != nil {
		return nil, fmt.Errorf("find tows failed: %w", err)
	}

	return tows, nil
}
