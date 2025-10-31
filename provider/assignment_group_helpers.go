package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type assignmentGroupResponse struct {
	Data struct {
		ID            int                          `json:"id"`
		Type          string                       `json:"type"`
		Attributes    assignmentGroupAttributes    `json:"attributes"`
		Relationships assignmentGroupRelationships `json:"relationships"`
	} `json:"data"`
}

type assignmentGroupAttributes struct {
	Name             string `json:"name"`
	AutoDeploy       bool   `json:"auto_deploy"`
	Type             string `json:"type"`
	InstallType      string `json:"install_type"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	DeviceCount      int    `json:"device_count"`
	GroupCount       int    `json:"group_count"`
	Priority         *int   `json:"priority,omitempty"`
	AppTrackLocation *bool  `json:"app_track_location,omitempty"`
}

type assignmentGroupRelationships struct {
	Apps         assignmentGroupRelationshipItems `json:"apps"`
	Profiles     assignmentGroupRelationshipItems `json:"profiles"`
	Devices      assignmentGroupRelationshipItems `json:"devices"`
	DeviceGroups assignmentGroupRelationshipItems `json:"device_groups"`
}

type assignmentGroupRelationshipItems struct {
	Data []assignmentGroupRelationshipItem `json:"data"`
}

type assignmentGroupRelationshipItem struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

func fetchAssignmentGroup(ctx context.Context, client *simplemdm.Client, id string) (*assignmentGroupResponse, error) {
	url := fmt.Sprintf("https://%s/api/v1/assignment_groups/%s", client.HostName, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	var assignmentGroup assignmentGroupResponse
	if err := json.Unmarshal(body, &assignmentGroup); err != nil {
		return nil, err
	}

	return &assignmentGroup, nil
}

func buildStringSetFromRelationshipItems(items []assignmentGroupRelationshipItem) types.Set {
	// Return empty set instead of null for Optional+Computed attributes
	// This prevents "was X but now null" errors when API doesn't return relationships
	values := make([]attr.Value, 0, len(items))
	for _, item := range items {
		values = append(values, types.StringValue(strconv.Itoa(item.ID)))
	}

	return types.SetValueMust(types.StringType, values)
}

type assignmentGroupUpsertRequest struct {
	Name             string
	AutoDeploy       *bool
	GroupType        *string
	InstallType      *string
	Priority         *int64
	AppTrackLocation *bool
}

func createAssignmentGroup(ctx context.Context, client *simplemdm.Client, payload assignmentGroupUpsertRequest) (*assignmentGroupResponse, error) {
	url := fmt.Sprintf("https://%s/api/v1/assignment_groups", client.HostName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("name", payload.Name)

	if payload.AutoDeploy != nil {
		q.Add("auto_deploy", strconv.FormatBool(*payload.AutoDeploy))
	}

	if payload.GroupType != nil && *payload.GroupType != "" {
		q.Add("type", *payload.GroupType)
	}

	if payload.InstallType != nil && *payload.InstallType != "" {
		q.Add("install_type", *payload.InstallType)
	}

	if payload.Priority != nil {
		q.Add("priority", strconv.FormatInt(*payload.Priority, 10))
	}

	if payload.AppTrackLocation != nil {
		q.Add("app_track_location", strconv.FormatBool(*payload.AppTrackLocation))
	}

	req.URL.RawQuery = q.Encode()

	body, err := client.RequestResponse201(req)
	if err != nil {
		return nil, err
	}

	var assignmentGroup assignmentGroupResponse
	if err := json.Unmarshal(body, &assignmentGroup); err != nil {
		return nil, err
	}

	return &assignmentGroup, nil
}

func updateAssignmentGroup(ctx context.Context, client *simplemdm.Client, id string, payload assignmentGroupUpsertRequest) error {
	url := fmt.Sprintf("https://%s/api/v1/assignment_groups/%s", client.HostName, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	if payload.Name != "" {
		q.Add("name", payload.Name)
	}

	if payload.AutoDeploy != nil {
		q.Add("auto_deploy", strconv.FormatBool(*payload.AutoDeploy))
	}

	if payload.GroupType != nil && *payload.GroupType != "" {
		q.Add("type", *payload.GroupType)
	}

	if payload.InstallType != nil && *payload.InstallType != "" {
		q.Add("install_type", *payload.InstallType)
	}

	if payload.Priority != nil {
		q.Add("priority", strconv.FormatInt(*payload.Priority, 10))
	}

	if payload.AppTrackLocation != nil {
		q.Add("app_track_location", strconv.FormatBool(*payload.AppTrackLocation))
	}

	req.URL.RawQuery = q.Encode()

	_, err = client.RequestResponse204(req)
	return err
}

func assignmentGroupAssignDevice(ctx context.Context, client *simplemdm.Client, groupID string, deviceID string, removeOthers bool) error {
	url := fmt.Sprintf("https://%s/api/v1/assignment_groups/%s/devices/%s", client.HostName, groupID, deviceID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	if removeOthers {
		q := req.URL.Query()
		q.Add("remove_others", "true")
		req.URL.RawQuery = q.Encode()
	}

	_, err = client.RequestResponse204(req)
	return err
}

func applyAssignmentGroupResponseToResourceModel(model *assignment_groupResourceModel, response *assignmentGroupResponse) {
	model.ID = types.StringValue(strconv.Itoa(response.Data.ID))
	model.Apps = buildStringSetFromRelationshipItems(response.Data.Relationships.Apps.Data)
	model.Groups = buildStringSetFromRelationshipItems(response.Data.Relationships.DeviceGroups.Data)
	model.Devices = buildStringSetFromRelationshipItems(response.Data.Relationships.Devices.Data)
	model.Profiles = buildStringSetFromRelationshipItems(response.Data.Relationships.Profiles.Data)

	model.Name = types.StringValue(response.Data.Attributes.Name)
	model.AutoDeploy = types.BoolValue(response.Data.Attributes.AutoDeploy)
	model.GroupType = types.StringValue(response.Data.Attributes.Type)

	// install_type is only returned by API for munki groups
	// For standard groups, don't set it (keep existing plan value if any)
	if response.Data.Attributes.Type == "munki" && response.Data.Attributes.InstallType != "" {
		model.InstallType = types.StringValue(response.Data.Attributes.InstallType)
	} else if response.Data.Attributes.Type != "standard" {
		// For munki groups without install_type, set to null
		model.InstallType = types.StringNull()
	}
	// For standard groups, don't modify InstallType (preserves plan value)

	if response.Data.Attributes.Priority != nil {
		model.Priority = types.Int64Value(int64(*response.Data.Attributes.Priority))
	} else {
		model.Priority = types.Int64Null()
	}

	if response.Data.Attributes.AppTrackLocation != nil {
		model.AppTrackLocation = types.BoolValue(*response.Data.Attributes.AppTrackLocation)
	} else {
		model.AppTrackLocation = types.BoolNull()
	}

	// Always set CreatedAt and UpdatedAt for resources too
	model.CreatedAt = types.StringValue(response.Data.Attributes.CreatedAt)
	model.UpdatedAt = types.StringValue(response.Data.Attributes.UpdatedAt)

	model.DeviceCount = types.Int64Value(int64(response.Data.Attributes.DeviceCount))
	model.GroupCount = types.Int64Value(int64(response.Data.Attributes.GroupCount))
}

func applyAssignmentGroupResponseToDataSourceModel(model *assignmentGroupDataSourceModel, response *assignmentGroupResponse) {
	model.ID = types.StringValue(strconv.Itoa(response.Data.ID))
	model.Name = types.StringValue(response.Data.Attributes.Name)
	model.AutoDeploy = types.BoolValue(response.Data.Attributes.AutoDeploy)
	model.GroupType = types.StringValue(response.Data.Attributes.Type)

	if response.Data.Attributes.Type == "munki" && response.Data.Attributes.InstallType != "" {
		model.InstallType = types.StringValue(response.Data.Attributes.InstallType)
	} else {
		model.InstallType = types.StringNull()
	}

	if response.Data.Attributes.Priority != nil {
		model.Priority = types.Int64Value(int64(*response.Data.Attributes.Priority))
	} else {
		model.Priority = types.Int64Null()
	}

	if response.Data.Attributes.AppTrackLocation != nil {
		model.AppTrackLocation = types.BoolValue(*response.Data.Attributes.AppTrackLocation)
	} else {
		model.AppTrackLocation = types.BoolNull()
	}

	// Always set CreatedAt and UpdatedAt, even if empty
	// This ensures Computed fields in data sources are considered "set"
	model.CreatedAt = types.StringValue(response.Data.Attributes.CreatedAt)
	model.UpdatedAt = types.StringValue(response.Data.Attributes.UpdatedAt)

	model.Apps = buildStringSetFromRelationshipItems(response.Data.Relationships.Apps.Data)
	model.Groups = buildStringSetFromRelationshipItems(response.Data.Relationships.DeviceGroups.Data)
	model.Devices = buildStringSetFromRelationshipItems(response.Data.Relationships.Devices.Data)
	model.Profiles = buildStringSetFromRelationshipItems(response.Data.Relationships.Profiles.Data)

	model.DeviceCount = types.Int64Value(int64(response.Data.Attributes.DeviceCount))
	model.GroupCount = types.Int64Value(int64(response.Data.Attributes.GroupCount))
}
