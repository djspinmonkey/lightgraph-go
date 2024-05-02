package model

import (
	"encoding/json"
	"errors"
	"github.com/djspinmonkey/lightgraph-go/restapi"
)

type Alert struct {
	ID                   string
	Name                 string
	Description          string
	Labels               []Label
	EnableNoDataAlert    bool
	EnableNoDataDuration int
	Operand              string
	WarningThreshold     float64
	CriticalThreshold    float64
	Project              *Project
}

// FetchAlerts fetches all alerts for a given project from the backing API.
//
// I'd like this to be in the restapi package, but it's not clear how to do that without creating a circular dependency.
func FetchAlerts(project *Project) ([]*Alert, error) {
	response, err := restapi.GetResource("/" + project.Organization.ID + "/projects/" + project.ID + "/metric_alerts")
	if err != nil {
		return nil, errors.New("Failed to fetch alerts: " + err.Error())
	}

	// This intermediate representation seems to be required in Go in order to get the deeply nested JSON data into a
	// flat Go struct (ie, Project). I'd love to find a better way to do this!
	var jsonShapedAlerts struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Labels      []Label `json:"labels"`
				Expression  struct {
					Operand           string `json:"operand"`
					EnableNoDataAlert bool   `json:"enable-no-data-alert"`
					NoDataDuration    int    `json:"no-data-duration-ms"`
					Thresholds        struct {
						Warning  float64 `json:"warning"`
						Critical float64 `json:"critical"`
					}
				}
			}
		}
	}
	err = json.NewDecoder(response.Body).Decode(&jsonShapedAlerts)
	if err != nil {
		return nil, errors.New("Failed to parse alerts: " + err.Error())
	}

	alerts := make([]*Alert, len(jsonShapedAlerts.Data))
	for i, jsonAlert := range jsonShapedAlerts.Data {
		alerts[i] = &Alert{
			ID:                   jsonAlert.ID,
			Name:                 jsonAlert.Attributes.Name,
			Description:          jsonAlert.Attributes.Description,
			Labels:               jsonAlert.Attributes.Labels,
			EnableNoDataAlert:    jsonAlert.Attributes.Expression.EnableNoDataAlert,
			EnableNoDataDuration: jsonAlert.Attributes.Expression.NoDataDuration,
			Operand:              jsonAlert.Attributes.Expression.Operand,
			WarningThreshold:     jsonAlert.Attributes.Expression.Thresholds.Warning,
			CriticalThreshold:    jsonAlert.Attributes.Expression.Thresholds.Critical,
			Project:              project,
		}
	}

	return alerts, nil
}
