package utilities

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/geoplaces"
	"github.com/aws/aws-sdk-go-v2/service/georoutes"
	"github.com/aws/aws-sdk-go-v2/service/georoutes/types"
)

// LocationUtility provides methods for interacting with Amazon Location Service.
type LocationUtility struct {
	geoplacesClient *geoplaces.Client
	georoutesClient *georoutes.Client
	indexName       string
}

// NewLocationUtility initializes AWS Location Service geoplacesClient using default AWS config and returns a utility instance.
func NewLocationUtility() (*LocationUtility, error) {

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))

	if err != nil {
		return nil, errors.New("failed to load AWS config: " + err.Error())
	}

	placesClient := geoplaces.NewFromConfig(cfg)
	georoutesClient := georoutes.NewFromConfig(cfg)

	return &LocationUtility{
		geoplacesClient: placesClient,
		georoutesClient: georoutesClient,
	}, nil
}

// GenerateAddressSuggestions takes a query string and returns an array of address suggestions.
// The query parameter should contain the partial address or location text to search for.
func (a *LocationUtility) GenerateAddressSuggestions(query string, biasLocation []float64) []string {

	// Create the input suggestion using provided query
	params := &geoplaces.SuggestInput{
		QueryText:    aws.String(query),
		MaxResults:   aws.Int32(7),
		BiasPosition: biasLocation,
	}

	// generate suggestions using ALS geoplacesClient
	suggestOutput, err := a.geoplacesClient.Suggest(context.Background(), params)

	if err != nil {
		return nil
	}

	// Convert suggested responses to string array
	var results []string

	for _, item := range suggestOutput.ResultItems {
		results = append(results, *item.Title)
	}

	return results
}

// ParseGeocodeFromAddress takes an address string and returns
// a longitude and latitude nearest the location of the address.
func (a *LocationUtility) ParseGeocodeFromAddress(address string) ([]float64, error) {

	// Create the input geocode using provided address
	params := &geoplaces.GeocodeInput{
		QueryText: aws.String(address),
	}

	// generate longitude and latitude using ALS geoplacesClient
	geocodeOutput, err := a.geoplacesClient.Geocode(context.Background(), params)

	if err != nil {
		return []float64{}, err
	}

	// return in long/lat position
	return []float64{geocodeOutput.ResultItems[0].Position[0], geocodeOutput.ResultItems[0].Position[1]}, nil
}

// CalculateDistanceBetweenCoordinates calculates the distance between two coordinates in miles.
func (a *LocationUtility) CalculateDistanceBetweenCoordinates(coordinates1, coordinates2 []float64) (float64, error) {

	params := &georoutes.CalculateRoutesInput{
		Origin:                        coordinates1,
		Destination:                   coordinates2,
		InstructionsMeasurementSystem: types.MeasurementSystemMetric,
	}

	// Calculate the distance in meters between the two coordinates using the Haversine formula
	routesOutput, err := a.georoutesClient.CalculateRoutes(context.TODO(), params)

	if err != nil {
		return 0, err
	}

	// Convert the distance to miles
	return float64(routesOutput.Routes[0].Summary.Distance) / 1609.00, nil
}
