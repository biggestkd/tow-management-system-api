package service

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tow-management-system-api/model"
)

// TowFinder is the minimal dependency MetricService needs from your repository layer.
type TowFinder interface {
	Find(ctx context.Context, filterModel *model.Tow) ([]*model.Tow, error)
}

// MetricService computes domain KPIs from tows (e.g., Active Tows).
type MetricService struct {
	towRepo TowFinder
}

// NewMetricService constructs a MetricService
func NewMetricService(towRepo TowFinder) *MetricService {
	return &MetricService{towRepo: towRepo}
}

var activeLabel = "Active Tows"
var totalLabel = "Total Tows"
var payoutLabel = "Payout Amount"

// CalculateMetrics computes metrics for a company and returns a slice of Metric documents.
func (m *MetricService) CalculateMetrics(ctx context.Context, companyID string) ([]*model.Metric, error) {
	if companyID == "" {
		return nil, fmt.Errorf("company id is required")
	}

	// Pull all tows for the company. We’re intentionally not filtering by status here
	// so we can derive multiple metrics in one pass with a single repository call.
	tows, err := m.towRepo.Find(ctx, &model.Tow{
		CompanyID: &companyID,
	})

	if err != nil {
		return nil, fmt.Errorf("find tows failed: %w", err)
	}

	activeTotal, completedTotal := calculateMetricsHelper(tows)

	// Build metrics (all fields are pointers; aligns with your Metric struct)
	activeTotalStr := fmt.Sprintf("%d", activeTotal)
	completedTotalStr := fmt.Sprintf("%d", completedTotal)
	payoutAmountStr := "0.00" // TODO: use actual values
	now := time.Now().Unix()

	metrics := []*model.Metric{
		{
			CompanyID:   &companyID,
			Type:        &activeLabel,
			Value:       &activeTotalStr,
			LastUpdated: &now,
		},
		{
			CompanyID:   &companyID,
			Type:        &totalLabel,
			Value:       &completedTotalStr,
			LastUpdated: &now,
		},
		{
			CompanyID:   &companyID,
			Type:        &payoutLabel,
			Value:       &payoutAmountStr,
			LastUpdated: &now,
		},
	}
	return metrics, nil
}

func calculateMetricsHelper(tows []*model.Tow) (int32, int32) {

	var (
		completed int32
		active    int32
	)

	terminal := map[string]struct{}{
		"COMPLETED": {},
	}

	activeStatuses := map[string]struct{}{
		"ACCEPTED":       {},
		"DISPATCHED":     {},
		"ARRIVED_PICKUP": {},
		"IN_TRANSIT":     {},
	}

	for _, t := range tows {

		if t.Status == nil {
			continue // skip if unknown status
		}

		// Normalize to uppercase to avoid case mismatches
		status := strings.ToUpper(*t.Status)

		// compute the active tows
		if _, isActive := activeStatuses[status]; isActive {
			active++
		}

		// compute the completed tows
		if _, isComplete := terminal[status]; isComplete {
			completed++
		}

	}

	return active, completed
}
