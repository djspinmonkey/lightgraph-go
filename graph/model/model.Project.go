package model

import (
	"encoding/json"
	"errors"

	"github.com/djspinmonkey/lightgraph-go/restapi"
)

type Project struct {
	ID                string
	Name              string
	Organization      *Organization
	alerts            []*Alert
	alertDestinations []*AlertDestination
}

type JsonShapedProject struct {
	Data struct {
		Attributes struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
	}
}

// Alerts returns all alerts for the project. It caches the alerts after the first request.
func (p *Project) Alerts() ([]*Alert, error) {
	if p.alerts == nil {
		var err error
		p.alerts, err = FetchAlerts(p)
		if err != nil {
			return nil, err
		}
	}

	return p.alerts, nil
}

// Alert returns the alert with the given ID, or nil if it doesn't exist or isn't associated with this project.
func (p *Project) Alert(id string) (*Alert, error) {
	alerts, err := p.Alerts()
	if err != nil {
		return nil, err
	}

	for _, alert := range alerts {
		if alert.ID == id {
			return alert, nil
		}
	}

	return nil, nil
}

// AlertDestinations returns all alert destinations for the project. It caches the destinations after the first request.
func (p *Project) AlertDestinations() ([]*AlertDestination, error) {
	if p.alertDestinations == nil {
		var err error
		p.alertDestinations, err = FetchAlertDestinations(p)
		if err != nil {
			return nil, err
		}
	}

	return p.alertDestinations, nil
}

// AlertDestination returns the alert destination with the given ID,
// or nil if it doesn't exist or isn't associated with this project.
func (p *Project) AlertDestination(id string) (*AlertDestination, error) {
	alertDestinations, err := p.AlertDestinations()
	if err != nil {
		return nil, err
	}

	for _, alertDestination := range alertDestinations {
		if alertDestination.ID == id {
			return alertDestination, nil
		}
	}

	return nil, nil
}

// Below this are all the Fetch* functions that fetch data from the backing API. I'd like to move these to the
// restapi package, but I'm not sure how to do that without creating a circular dependency.

// FetchProject submits a GET request to the REST API for the project with the given org and project IDs.
func FetchProject(org *Organization, projectID string) (*Project, error) {
	response, err := restapi.GetCloudObsResource("/" + org.ID + "/projects/" + projectID)
	if err != nil {
		return nil, errors.New("Failed to fetch project: " + err.Error())
	}

	var jsonShapedProject JsonShapedProject
	err = json.NewDecoder(response.Body).Decode(&jsonShapedProject)
	if err != nil {
		return nil, errors.New("Failed to parse project: " + err.Error())
	}

	project := &Project{
		ID:                jsonShapedProject.Data.Attributes.ID,
		Name:              jsonShapedProject.Data.Attributes.Name,
		Organization:      org,
		alerts:            []*Alert{},
		alertDestinations: []*AlertDestination{},
	}

	return project, nil
}
