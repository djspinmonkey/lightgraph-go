package model

import (
	"github.com/djspinmonkey/lightgraph-go/restapi"
)

type Actor struct{}

func (a Actor) BackingAPIURL() (string, error) {
	return restapi.CloudObsBaseUrl(), nil
}

func (a Actor) APIKey() (string, error) {
	return restapi.CloudObsApiKey(), nil
}

func (a Actor) Test() (string, error) {
	return "You have successfully queried this test field!", nil
}
