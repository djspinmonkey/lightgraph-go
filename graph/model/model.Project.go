package model

import (
	"encoding/json"
	"errors"
	"github.com/djspinmonkey/lightgraph-go/restapi"
)

type Project struct {
	ID   string
	Name string
}

// FetchProject submits a GET request to the REST API for the project with the given org and project IDs.
//
// I'd like this to be in the restapi package, but it's not clear how to do that without creating a circular dependency.
func FetchProject(orgID string, projectID string) (*Project, error) {
	response, err := restapi.GetResource("/" + orgID + "/projects/" + projectID)
	if err != nil {
		return nil, errors.New("Failed to fetch project: " + err.Error())
	}

	// This intermediate representation seems to be required in Go in order to get the deeply nested JSON data into a
	// flat Go struct (ie, Project). I'd love to find a better way to do this!
	var tempProject struct {
		Data struct {
			Attributes struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
		}
	}
	err = json.NewDecoder(response.Body).Decode(&tempProject)
	if err != nil {
		return nil, errors.New("Failed to parse project: " + err.Error())
	}

	project := &Project{
		ID:   tempProject.Data.Attributes.ID,
		Name: tempProject.Data.Attributes.Name,
	}

	return project, nil
}
