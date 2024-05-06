package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/djspinmonkey/lightgraph-go/restapi"
)

// Alert represents a single metric alert.
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

// Alerts is a collection of Alert objects. It's mostly just used for JSON parsing purposes.
type Alerts []*Alert

// JsonShapedAlerts is an intermediate representation of the JSON data returned by the API.
//
// I would love to parse the JSON straight into the desired struct without this intermediate representation, but
// this seems to be The Go Way™.
type JsonShapedAlerts struct {
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

// FetchAlerts fetches all alerts for a given project from the backing API.
//
// I'd like this to be in the restapi package, but it's not clear how to do that without creating a circular dependency.
func FetchAlerts(project *Project) ([]*Alert, error) {
	response, err := restapi.GetResource("/" + project.Organization.ID + "/projects/" + project.ID + "/metric_alerts")
	if err != nil {
		return nil, errors.New("Failed to fetch alerts: " + err.Error())
	}

	var alerts Alerts
	err = json.NewDecoder(response.Body).Decode(&alerts)
	if err != nil {
		return nil, errors.New("Failed to parse alerts: " + err.Error())
	}

	for _, alert := range alerts {
		alert.Project = project
	}

	return alerts, nil
}

// UnmarshalJSON parses the JSON data returned by the API into a collection of Alert structs. Note that the receiver
// is a pointer to Alerts, not a single Alert.
func (a *Alerts) UnmarshalJSON(rawJson []byte) error {
	var parsedJson JsonShapedAlerts
	jsonReader := bytes.NewReader(rawJson)
	err := json.NewDecoder(jsonReader).Decode(&parsedJson)
	if err != nil {
		return err
	}

	for _, d := range parsedJson.Data {
		*a = append(*a, &Alert{
			ID:                   d.ID,
			Name:                 d.Attributes.Name,
			Description:          d.Attributes.Description,
			Labels:               d.Attributes.Labels,
			EnableNoDataAlert:    d.Attributes.Expression.EnableNoDataAlert,
			EnableNoDataDuration: d.Attributes.Expression.NoDataDuration,
			Operand:              d.Attributes.Expression.Operand,
			WarningThreshold:     d.Attributes.Expression.Thresholds.Warning,
			CriticalThreshold:    d.Attributes.Expression.Thresholds.Critical,
		})
	}

	return nil
}
