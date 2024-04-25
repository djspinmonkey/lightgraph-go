package model

import (
	"github.com/djspinmonkey/lightgraph-go/restapi"
)

type Actor struct{}

func (a Actor) BackingAPIURL() (string, error) {
	return restapi.BaseUrl(), nil
}

func (a Actor) APIKey() (string, error) {
	return restapi.ApiKey(), nil
}

func (a Actor) Test() (string, error) {
	return "You have successfully queried this test field!", nil
}
