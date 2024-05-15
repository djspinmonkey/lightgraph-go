package restapi

import (
	"errors"
	"fmt"
	"net/http"
)

type LoginCredentials struct {
	Username string
	Password string
}

// TODO: use configured parameters or env variables for login and BaseURL

const testSNOWInstanceURL = "https://demoallwf60414.service-now.com"

// ServiceNowLoginCredentials returns a hardcoded username and password for accessing the demoallwf60414.service-now instance
func ServiceNowLoginCredentials() LoginCredentials {
	return LoginCredentials{Username: "admin", Password: "TX)K%9#:(VuS9Tn=S;+1"}
}

// ServiceNowBaseURL retrieves the base URL for the ServiceNow API
func ServiceNowBaseURL() string {
	return testSNOWInstanceURL
}

// GetServiceNowResource submits a GET request to the ServiceNow API at the given path, using the (currently hardcoded) base URL and API key.
func GetServiceNowResource(path string) (*http.Response, error) {
	url := ServiceNowBaseURL() + path
	fmt.Printf("\n******* requesting resource: %s\n", url) // debugging output

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "lightgraph-go")
	req.Header.Add("Accept", "application/json")

	creds := ServiceNowLoginCredentials()
	req.SetBasicAuth(creds.Username, creds.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("ServiceNow API returned status: " + resp.Status)
	}

	return resp, nil
}

// TODO: Move the various FetchFoo functions to be in the restapi package.
// However, it's not clear how to do that without creating a circular dependency.
