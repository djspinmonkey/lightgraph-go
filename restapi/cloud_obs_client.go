package restapi

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

// CloudObsApiKey retrieves the API key from the environment.
func CloudObsApiKey() string {
	key := os.Getenv("LS_TOKEN")
	if key == "" {
		panic("Cannot access REST API: no API key found in $LS_TOKEN")
	}

	return key
}

// CloudObsBaseUrl retrieves the base URL for the Cloud Obs REST API from the environment.
func CloudObsBaseUrl() string {
	url := os.Getenv("LS_REST_API_URL")
	if url == "" {
		panic("Cannot access REST API: no base URL found in $LS_REST_API_URL")
	}

	return url
}

// GetCloudObsResource submits a GET request to the Cloud Obs REST API at the given path, using the configured base URL and API key.
func GetCloudObsResource(path string) (*http.Response, error) {
	url := CloudObsBaseUrl() + path
	fmt.Printf("\n******* requesting resource: %s\n", url) // debugging output

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "lightgraph-go")
	req.Header.Add("Authorization", CloudObsApiKey())

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("REST API returned status: " + resp.Status)
	}

	return resp, nil
}

// TODO: Move the various FetchFoo functions to be in the restapi package.
// However, it's not clear how to do that without creating a circular dependency.
