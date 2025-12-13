package service

import (
	"tow-management-system-api/utilities"
)

// LocationService encapsulates location-related business logic.
type LocationService struct {
	locationUtility *utilities.LocationUtility
}

// NewLocationService constructs a LocationService with the provided LocationUtility.
func NewLocationService(locationUtility *utilities.LocationUtility) *LocationService {
	return &LocationService{
		locationUtility: locationUtility,
	}
}

func (s *LocationService) GenerateAddressSuggestions(query string) []string {
	// Hard-coding NYC bias position
	biasPosition := []float64{-74.0060, 40.7128}

	return s.locationUtility.GenerateAddressSuggestions(query, biasPosition)
}
