package model

type AlertingRule struct {
	ID                         string `json:"id"`
	UpdateInterval             int    `json:"update-interval-ms"`
	MessageDestinationClientId string `json:"message-destination-client-id"`
	Alert                      *Alert
	alertDestination           *AlertDestination
}

func (ar *AlertingRule) Destination() (*AlertDestination, error) {
	if ar.alertDestination == nil {
		var err error
		ar.alertDestination, err = ar.Alert.Project.AlertDestination(ar.MessageDestinationClientId)
		if err != nil {
			return nil, err
		}
	}

	return ar.alertDestination, nil
}
