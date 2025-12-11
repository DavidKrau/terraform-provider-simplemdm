package simplemdmext

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
)

// DeviceResponse models the payload returned from the SimpleMDM device endpoints.
type DeviceResponse struct {
	Data DeviceData `json:"data"`
}

// DeviceListResponse models the paginated device collection response.
type DeviceListResponse struct {
	Data    []DeviceData `json:"data"`
	HasMore bool         `json:"has_more"`
}

// DeviceData contains a device record with attributes and relationships.
type DeviceData struct {
	Type          string              `json:"type"`
	ID            int                 `json:"id"`
	Attributes    map[string]any      `json:"attributes"`
	Relationships DeviceRelationships `json:"relationships"`
}

// DeviceRelationships captures the relationships needed by the Terraform provider.
type DeviceRelationships struct {
	DeviceGroup struct {
		Data struct {
			Type string `json:"type"`
			ID   int    `json:"id"`
		} `json:"data"`
	} `json:"device_group"`
	Groups struct {
		Data []struct {
			Type string `json:"type"`
			ID   int    `json:"id"`
		} `json:"data"`
	} `json:"groups"`
	CustomAttributeValues struct {
		Data []struct {
			Type       string `json:"type"`
			ID         string `json:"id"`
			Attributes struct {
				Secret bool   `json:"secret"`
				Value  string `json:"value"`
			} `json:"attributes"`
		} `json:"data"`
	} `json:"custom_attribute_values"`
}

// DeviceRelatedListResponse models list payloads for related resources such as profiles.
type DeviceRelatedListResponse struct {
	Data    []DeviceRelatedItem `json:"data"`
	HasMore bool                `json:"has_more"`
}

// DeviceRelatedItem represents an entry returned from a related collection.
type DeviceRelatedItem struct {
	Type       string         `json:"type"`
	ID         jsonNumber     `json:"id"`
	Attributes map[string]any `json:"attributes"`
}

// jsonNumber wraps json.Number but keeps zero values serialisable.
type jsonNumber struct {
	raw string
}

func (n *jsonNumber) UnmarshalJSON(b []byte) error {
	n.raw = string(b)
	return nil
}

func (n jsonNumber) String() string {
	if n.raw == "" {
		return ""
	}

	// Trim surrounding quotes that indicate a JSON string representation.
	if len(n.raw) >= 2 && n.raw[0] == '"' && n.raw[len(n.raw)-1] == '"' {
		return n.raw[1 : len(n.raw)-1]
	}

	return n.raw
}

// GetDevice retrieves a single device record.
func GetDevice(ctx context.Context, client *simplemdm.Client, deviceID string, includeSecretCustomAttributes bool) (*DeviceResponse, error) {
	url := fmt.Sprintf("https://%s/api/v1/devices/%s", client.HostName, deviceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if includeSecretCustomAttributes {
		q := req.URL.Query()
		q.Add("include_secret_custom_attributes", "true")
		req.URL.RawQuery = q.Encode()
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	var resp DeviceResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListDevices retrieves all devices that satisfy the provided filters. It automatically
// walks through paginated responses using cursor-based pagination.
func ListDevices(ctx context.Context, client *simplemdm.Client, search string, includeAwaitingEnrollment, includeSecretCustomAttributes bool) ([]DeviceData, error) {
	results := make([]DeviceData, 0)
	var startingAfter string

	for {
		url := fmt.Sprintf("https://%s/api/v1/devices", client.HostName)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		q := req.URL.Query()
		q.Set("limit", "100")
		if startingAfter != "" {
			q.Set("starting_after", startingAfter)
		}
		if search != "" {
			q.Set("search", search)
		}
		if includeAwaitingEnrollment {
			q.Set("include_awaiting_enrollment", "true")
		}
		if includeSecretCustomAttributes {
			q.Set("include_secret_custom_attributes", "true")
		}
		req.URL.RawQuery = q.Encode()

		body, err := client.RequestResponse200(req)
		if err != nil {
			return nil, err
		}

		var resp DeviceListResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}

		results = append(results, resp.Data...)

		if !resp.HasMore || len(resp.Data) == 0 {
			break
		}

		// Set cursor to last item's ID for next page
		startingAfter = strconv.Itoa(resp.Data[len(resp.Data)-1].ID)
	}

	return results, nil
}

// ListDeviceProfiles fetches the profiles directly assigned to a device.
func ListDeviceProfiles(ctx context.Context, client *simplemdm.Client, deviceID string) (*DeviceRelatedListResponse, error) {
	return listRelated(ctx, client, deviceID, "profiles")
}

// ListDeviceInstalledApps fetches the installed applications for a device.
func ListDeviceInstalledApps(ctx context.Context, client *simplemdm.Client, deviceID string) (*DeviceRelatedListResponse, error) {
	return listRelated(ctx, client, deviceID, "installed_apps")
}

// ListDeviceUsers fetches the user accounts present on a device.
func ListDeviceUsers(ctx context.Context, client *simplemdm.Client, deviceID string) (*DeviceRelatedListResponse, error) {
	return listRelated(ctx, client, deviceID, "users")
}

func listRelated(ctx context.Context, client *simplemdm.Client, deviceID, endpoint string) (*DeviceRelatedListResponse, error) {
	allData := []DeviceRelatedItem{}
	var startingAfter string

	for {
		url := fmt.Sprintf("https://%s/api/v1/devices/%s/%s", client.HostName, deviceID, endpoint)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		q := req.URL.Query()
		q.Set("limit", "100")
		if startingAfter != "" {
			q.Set("starting_after", startingAfter)
		}
		req.URL.RawQuery = q.Encode()

		body, err := client.RequestResponse200(req)
		if err != nil {
			return nil, err
		}

		var resp DeviceRelatedListResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}

		allData = append(allData, resp.Data...)

		// Stop if there are no more pages or no data returned
		if !resp.HasMore || len(resp.Data) == 0 {
			break
		}

		// Set cursor to last item's ID for next page
		startingAfter = resp.Data[len(resp.Data)-1].ID.String()
	}

	return &DeviceRelatedListResponse{
		Data:    allData,
		HasMore: false,
	}, nil
}

// FlattenAttributes normalises an arbitrary map of attributes into a map of Terraform strings.
func FlattenAttributes(attrs map[string]any) map[string]string {
	flat := make(map[string]string, len(attrs))
	for key, value := range attrs {
		switch v := value.(type) {
		case nil:
			flat[key] = ""
		case string:
			flat[key] = v
		case bool:
			if v {
				flat[key] = "true"
			} else {
				flat[key] = "false"
			}
		case float64:
			// JSON numbers are floats; render without trailing zeros when possible.
			if v == float64(int64(v)) {
				flat[key] = strconv.FormatInt(int64(v), 10)
			} else {
				flat[key] = strconv.FormatFloat(v, 'f', -1, 64)
			}
		default:
			bytes, err := json.Marshal(v)
			if err != nil {
				flat[key] = fmt.Sprintf("%v", v)
				continue
			}
			flat[key] = string(bytes)
		}
	}

	return flat
}

// ConvertRelatedItems converts the list response into helper maps for Terraform.
func ConvertRelatedItems(items []DeviceRelatedItem) []map[string]string {
	results := make([]map[string]string, 0, len(items))
	for _, item := range items {
		converted := make(map[string]string)
		converted["id"] = item.ID.String()
		converted["type"] = item.Type
		if len(item.Attributes) > 0 {
			converted["attributes"] = marshalAttributes(item.Attributes)
		}
		results = append(results, converted)
	}

	return results
}

func marshalAttributes(attrs map[string]any) string {
	bytes, err := json.Marshal(attrs)
	if err != nil {
		return ""
	}

	return string(bytes)
}
