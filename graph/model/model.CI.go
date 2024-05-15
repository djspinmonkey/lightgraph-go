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
}

// JsonShapedCI is an intermediate representation of the JSON data returned by the API.
type JsonShapedCI struct {
	Result struct {
		Attributes struct {
			Name              string `json:"name"`
			AssetTag          string `json:"asset_tag"`
			SubCategory       string `json:"sub_category"`
			SerialNumber      string `json:"serial_number"`
			AssetLink         string `json:"asset_link"`
			AssetDisplayValue string `json:"asset_display_value"`
			AssetValue        string `json:"asset_value"`
			SysID             string `json:"sys_id"`
			ClassName         string `json:"class_name"`
		}
	}
}

// FetchCI fetches a CI for a given className and sysID
func FetchCI(c *CIIdentifier) (*CI, error) {
	response, err := restapi.GetServiceNowResource(fmt.Sprintf("/api/now/cmdb/instance/%s/%s", c.ClassName, c.SysID))
	if err != nil {
		return nil, errors.New("Failed to fetch CI: " + err.Error())
	}

	ciJSON := JsonShapedCI{}
	err = json.NewDecoder(response.Body).Decode(&ciJSON)
	if err != nil {
		return nil, errors.New("Failed to parse CI: " + err.Error())
	}

	ci := CI{
		CIIdentifier: &CIIdentifier{
			SysID:     ciJSON.Result.Attributes.SysID,
			ClassName: ciJSON.Result.Attributes.ClassName,
		},
		Name:              ciJSON.Result.Attributes.Name,
		AssetTag:          ciJSON.Result.Attributes.AssetTag,
		SubCategory:       ciJSON.Result.Attributes.SubCategory,
		SerialNumber:      ciJSON.Result.Attributes.SerialNumber,
		AssetLink:         ciJSON.Result.Attributes.AssetLink,
		AssetDisplayValue: ciJSON.Result.Attributes.AssetDisplayValue,
		AssetValue:        ciJSON.Result.Attributes.AssetValue,
	}

	return &ci, nil
}
