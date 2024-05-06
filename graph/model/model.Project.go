package model

import (
	"encoding/json"
	"errors"
	"github.com/djspinmonkey/lightgraph-go/restapi"
)

type Project struct {
	ID           string
	Name         string
	Organization *Organization
	alerts       []*Alert
}

// FetchProject submits a GET request to the REST API for the project with the given org and project IDs.
//
// I'd like this to be in the restapi package, but it's not clear how to do that without creating a circular dependency.
func FetchProject(org *Organization, projectID string) (*Project, error) {
	response, err := restapi.GetResource("/" + org.ID + "/projects/" + projectID)
	if err != nil {
		return nil, errors.New("Failed to fetch project: " + err.Error())
	}

	// This intermediate representation seems to be required in Go in order to get the deeply nested JSON data into a
	// flat Go struct (ie, Project). I'd love to find a better way to do this!
	var jsonShapedProject struct {
		Data struct {
			Attributes struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
		}
	}
	err = json.NewDecoder(response.Body).Decode(&jsonShapedProject)
	if err != nil {
		return nil, errors.New("Failed to parse project: " + err.Error())
	}

	project := &Project{
		ID:           jsonShapedProject.Data.Attributes.ID,
		Name:         jsonShapedProject.Data.Attributes.Name,
		Organization: org,
	}

	return project, nil
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

// Alert retrieves the alert with the given ID.
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
