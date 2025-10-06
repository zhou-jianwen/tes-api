package serviceinfo

import "code.byted.org/epscp/vetes-api/pkg/version"

// DefaultServiceInfo ...
var DefaultServiceInfo = ServiceInfo{
	ID:   "com.volcengine.bioos.veTES",
	Name: "veTES-api",
	Type: TypeInfo{
		Group:    "com.volcengine.bioos",
		Artifact: "TES",
		Version:  version.Get().Version,
	},
	Organization: OrganizationInfo{
		Name: "volcengine",
		URL:  "https://volcengine.com",
	},
	Version: version.Get().Version,
}

// ServiceInfo ...
type ServiceInfo struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Type             TypeInfo         `json:"type"`
	Description      string           `json:"description,omitempty"`
	Organization     OrganizationInfo `json:"organization"`
	ContactURL       string           `json:"contactURL,omitempty"`
	DocumentationURL string           `json:"documentationURL,omitempty"`
	CreatedAt        string           `json:"createdAt,omitempty"`
	UpdatedAt        string           `json:"updatedAt,omitempty"`
	Environment      string           `json:"environment,omitempty"`
	Version          string           `json:"version"`
	Storage          []string         `json:"storage,omitempty"`
}

// TypeInfo ...
type TypeInfo struct {
	Group    string `json:"group"`
	Artifact string `json:"artifact"`
	Version  string `json:"version"`
}

// OrganizationInfo ...
type OrganizationInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
