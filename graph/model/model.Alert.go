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
	Labels               []*Label
	EnableNoDataAlert    bool
	EnableNoDataDuration int
	Operand              string
	WarningThreshold     float64
	CriticalThreshold    float64
	Project              *Project
	AlertingRules        []*AlertingRule
	snoozification       *Snoozification
}

// Alerts is a collection of Alert objects. It's mostly just used for JSON parsing purposes.
type Alerts []*Alert

// JsonShapedAlerts is an intermediate representation of the JSON data returned by the API.
type JsonShapedAlerts struct {
	Data []struct {
		ID         string `json:"id"`
		Attributes struct {
			Name          string          `json:"name"`
			Description   string          `json:"description"`
			Labels        []*Label        `json:"labels"`
			AlertingRules []*AlertingRule `json:"alerting-rules"`
			Expression    struct {
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

// Snoozification represents the current snooze status of an alert.
type Snoozification struct {
	snoozed bool
	until   int64
}

// JsonShapedSnoozifications is an intermediate representation of the JSON data returned by the API.
type JsonShapedSnoozifications struct {
	Data []struct {
		Attributes struct {
			Until int64 `json:"ends-micros"`
		}
	}
}

// Destinations returns all destinations associated with the alert. This may involve fetching the destinations from the
// API if they haven't been fetched yet.
func (a *Alert) Destinations() ([]*AlertDestination, error) {
	var destinations []*AlertDestination

	for _, rule := range a.AlertingRules {
		destination, err := rule.Destination()
		if err != nil {
			return nil, err
		}

		destinations = append(destinations, destination)
	}

	return destinations, nil
}

// Snoozed returns true if the alert is snoozed, false otherwise. This will likely involve a request to the API.
func (a *Alert) Snoozed() (bool, error) {
	if a.snoozification == nil {
		s, err := a.FetchSnoozification()
		if err != nil {
			return false, err
		}

		a.snoozification = &s
	}

	return a.snoozification.snoozed, nil
}

// SnoozedUntil returns the time the alert is snoozed until, or 0 if the alert isn't snoozed. This will likely
// involve a request to the API.
func (a *Alert) SnoozedUntil() (int64, error) {
	if a.snoozification == nil {
		s, err := a.FetchSnoozification()
		if err != nil {
			return 0, err
		}

		a.snoozification = &s
	}

	return a.snoozification.until, nil
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
		alert := &Alert{
			ID:                   d.ID,
			Name:                 d.Attributes.Name,
			Description:          d.Attributes.Description,
			Labels:               d.Attributes.Labels,
			AlertingRules:        d.Attributes.AlertingRules,
			EnableNoDataAlert:    d.Attributes.Expression.EnableNoDataAlert,
			EnableNoDataDuration: d.Attributes.Expression.NoDataDuration,
			Operand:              d.Attributes.Expression.Operand,
			WarningThreshold:     d.Attributes.Expression.Thresholds.Warning,
			CriticalThreshold:    d.Attributes.Expression.Thresholds.Critical,
		}

		for _, rule := range alert.AlertingRules {
			rule.Alert = alert
		}

		*a = append(*a, alert)
	}

	return nil
}

// FetchAlerts fetches all alerts for a given project from the backing API.
func FetchAlerts(p *Project) ([]*Alert, error) {
	response, err := restapi.GetResource("/" + p.Organization.ID + "/projects/" + p.ID + "/metric_alerts")
	if err != nil {
		return nil, errors.New("Failed to fetch alerts: " + err.Error())
	}

	var alerts Alerts
	err = json.NewDecoder(response.Body).Decode(&alerts)
	if err != nil {
		return nil, errors.New("Failed to parse alerts: " + err.Error())
	}

	for _, alert := range alerts {
		alert.Project = p
		for _, rule := range alert.AlertingRules {
			rule.Alert = alert
		}
	}

	return alerts, nil
}

// FetchSnoozification fetches the Snoozification status for the alert from the backing API.
func (a *Alert) FetchSnoozification() (Snoozification, error) {
	response, err := restapi.GetResource("/" + a.Project.Organization.ID + "/projects/" + a.Project.ID + "/metric_alerts/" + a.ID + "/snoozes")
	if err != nil {
		return Snoozification{}, errors.New("Failed to fetch Snoozification: " + err.Error())
	}

	var jsonShapedSnoozifications JsonShapedSnoozifications
	err = json.NewDecoder(response.Body).Decode(&jsonShapedSnoozifications)
	if err != nil {
		return Snoozification{}, errors.New("Failed to parse Snoozification: " + err.Error())
	}

	if len(jsonShapedSnoozifications.Data) == 0 {
		return Snoozification{
			snoozed: false,
			until:   0,
		}, nil
	}

	return Snoozification{
		snoozed: true,
		until:   jsonShapedSnoozifications.Data[0].Attributes.Until,
	}, nil
}
