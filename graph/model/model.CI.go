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
			Name         string `json:"name"`
			SubCategory  string `json:"sub_category"`
			AssetTag     string `json:"asset_tag"`
			SerialNumber string `json:"serial_number"`
			SysID        string `json:"sys_id"`
			ClassName    string `json:"class_name"`
			Asset        struct {
				AssetLink         string `json:"link"`
				AssetDisplayValue string `json:"display_value"`
				AssetValue        string `json:"value"`
			}
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
		CIIdentifier:      c,
		Name:              ciJSON.Result.Attributes.Name,
		AssetTag:          ciJSON.Result.Attributes.AssetTag,
		SubCategory:       ciJSON.Result.Attributes.SubCategory,
		SerialNumber:      ciJSON.Result.Attributes.SerialNumber,
		AssetLink:         ciJSON.Result.Attributes.Asset.AssetLink,
		AssetDisplayValue: ciJSON.Result.Attributes.Asset.AssetDisplayValue,
		AssetValue:        ciJSON.Result.Attributes.Asset.AssetValue,
	}

	return &ci, nil
}
