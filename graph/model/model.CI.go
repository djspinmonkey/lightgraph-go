package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/djspinmonkey/lightgraph-go/restapi"
)

type CI struct {
	CIIdentifier      *CIIdentifier
	Name              string
	AssetTag          string
	SubCategory       string
	SerialNumber      string
	AssetLink         string
	AssetDisplayValue string
	AssetValue        string
	AttributesJSON    string
}

// FetchCI fetches a CI for a given className and sysID
func FetchCI(c *CIIdentifier) (*CI, error) {
	response, err := restapi.GetServiceNowResource(fmt.Sprintf("api/now/cmdb/instance/%s/%s", c.ClassName, c.SysID))
	if err != nil {
		return nil, errors.New("Failed to fetch CI: " + err.Error())
	}

	var ci CI
	err = json.NewDecoder(response.Body).Decode(&ci)
	if err != nil {
		return nil, errors.New("Failed to parse CI: " + err.Error())
	}

	return &ci, nil
}
