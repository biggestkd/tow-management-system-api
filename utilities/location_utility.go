package utilities

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/geoplaces"
)

// LocationUtility provides methods for interacting with Amazon Location Service.
type LocationUtility struct {
	client    *geoplaces.Client
	indexName string
}

// NewLocationUtility initializes AWS Location Service client using default AWS config and returns a utility instance.
func NewLocationUtility() (*LocationUtility, error) {

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))

	if err != nil {
		return nil, errors.New("failed to load AWS config: " + err.Error())
	}

	client := geoplaces.NewFromConfig(cfg)

	return &LocationUtility{
		client: client,
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

	// generate suggestions using ALS client
	suggestOutput, err := a.client.Suggest(context.Background(), params)

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
