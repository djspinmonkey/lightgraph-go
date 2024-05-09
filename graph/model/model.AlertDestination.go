package model

import (
	"encoding/json"
	"errors"
	"github.com/djspinmonkey/lightgraph-go/restapi"
)

// AlertDestination represents anywhere an alert can be sent. This could be a webhook, Slack channel, PagerDuty, etc.
type AlertDestination struct {
	ID              string
	Name            string
	DestinationType string
	Url             string
	CustomHeaders   []*CustomHeader
	Channel         string
	Scope           string
	BodyTemplate    string
	IntegrationKey  string
	ServiceNowAuth  []*AuthValue
	Project         *Project
}

// AlertDestinations is a collection of AlertDestination objects. It's mostly just used for JSON parsing purposes.
type AlertDestinations []*AlertDestination

// JsonShapedAlertDestinations is an intermediate representation of the JSON data returned by the API. It's only used for
// parsing the JSON data.
type JsonShapedAlertDestinations struct {
	Data []struct {
		ID         string
		Attributes struct {
			Name            string
			DestinationType string `json:"destination_type"`
			Url             string
			CustomHeaders   map[string]string `json:"custom_headers"`
			Channel         string
			Scope           string
			BodyTemplate    string            `json:"template"`
			IntegrationKey  string            `json:"integration_key"`
			ServiceNowAuth  map[string]string `json:"auth"`
		}
	}
}

// TODO: Break out the various types of destinations into their own structs and GraphQL types.
// There should also be a shared AlertDestination interface. Here are the specific types of destinations we support and
// their associated fields (in addition to "destination_type"). This is based on observing the API responses in staging
// and may not be exhaustive.
//
// "webhook"=>["name", "url", "custom_headers", "template"],
// "bigpanda"=>["name", "url"],
// "slack"=>["channel", "scope"],
// "pagerduty"=>["name", "integration_key"],
// "servicenow"=>["name", "url", "auth"]

// Type returns the type of the alert destination. This is an alias for DestinationType.
func (ad *AlertDestination) Type() string {
	return ad.DestinationType
}

// FetchAlertDestinations fetches all alert destinations for a given project from the backing API.
func FetchAlertDestinations(project *Project) ([]*AlertDestination, error) {
	response, err := restapi.GetResource("/" + project.Organization.ID + "/projects/" + project.ID + "/destinations")
	if err != nil {
		return nil, errors.New("Failed to fetch alert destinations: " + err.Error())
	}

	var jsonShapedAlertDestinations JsonShapedAlertDestinations
	err = json.NewDecoder(response.Body).Decode(&jsonShapedAlertDestinations)
	if err != nil {
		return nil, errors.New("Failed to parse alert destinations: " + err.Error())
	}

	alertDestinations := make([]*AlertDestination, len(jsonShapedAlertDestinations.Data))
	for i, alertDestination := range jsonShapedAlertDestinations.Data {
		var authValues []*AuthValue
		for k, v := range alertDestination.Attributes.ServiceNowAuth {
			authValues = append(authValues, &AuthValue{Key: k, Value: v})
		}

		var customHeaders []*CustomHeader
		for k, v := range alertDestination.Attributes.CustomHeaders {
			customHeaders = append(customHeaders, &CustomHeader{Key: k, Value: v})
		}

		alertDestinations[i] = &AlertDestination{
			ID:              alertDestination.ID,
			Name:            alertDestination.Attributes.Name,
			DestinationType: alertDestination.Attributes.DestinationType,
			Url:             alertDestination.Attributes.Url,
			CustomHeaders:   customHeaders,
			Channel:         alertDestination.Attributes.Channel,
			Scope:           alertDestination.Attributes.Scope,
			BodyTemplate:    alertDestination.Attributes.BodyTemplate,
			IntegrationKey:  alertDestination.Attributes.IntegrationKey,
			ServiceNowAuth:  authValues,
			Project:         project,
		}
	}

	return alertDestinations, nil
}
