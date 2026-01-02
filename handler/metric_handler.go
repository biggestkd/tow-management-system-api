package handler

import (
	"context"
	"log"
	"net/http"
	"tow-management-system-api/model"

	"github.com/gin-gonic/gin"
)

// MetricService defines the minimal contract for computing metrics.
type MetricService interface {
	CalculateMetrics(ctx context.Context, companyId string) ([]*model.Metric, error)
}

// MetricHandler handles HTTP routes for metric-related operations.
type MetricHandler struct {
	metricService MetricService
}

// NewMetricHandler creates a new MetricHandler instance.
func NewMetricHandler(service MetricService) *MetricHandler {
	return &MetricHandler{
		metricService: service,
	}
}

// GetCompanyMetrics GET /metrics/:companyId
// Retrieves computed metrics (e.g., Active Tows, Completed Tows) for a given company.
//
// Response:
//
//	200 [Metric] - JSON array of metrics
//	400 - "company id is required"
//	500 - "something went wrong"
func (h *MetricHandler) GetCompanyMetrics(c *gin.Context) {
	companyId := c.Param("companyId")
	if companyId == "" {
		c.String(http.StatusBadRequest, "company id is required")
		return
	}

	metrics, err := h.metricService.CalculateMetrics(c.Request.Context(), companyId)
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusInternalServerError, "something went wrong")
		return
	}

	c.JSON(http.StatusOK, metrics)
}
