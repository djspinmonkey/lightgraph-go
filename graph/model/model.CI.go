package model

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

// FetchCI fetches a CI from the CMDB.
func FetchCI(identifier *CIIdentifier) (*CI, error) {
	// This is a stub implementation that returns a CI with the given identifier.
	return &CI{CIIdentifier: identifier, Name: "stub"}, nil
}
