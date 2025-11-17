package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

	req.URL.RawQuery = buildAssignmentGroupQuery(payload, true).Encode()

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

	req.URL.RawQuery = buildAssignmentGroupQuery(payload, false).Encode()

	_, err = client.RequestResponse204(req)
	return err
}

// buildAssignmentGroupQuery constructs the query parameters shared by the create and update operations.
// When includeName is true the "name" parameter is always sent, mirroring the API requirement for creation requests.
// For updates, the name is only provided when it has a non-empty value so partial updates remain possible.
func buildAssignmentGroupQuery(payload assignmentGroupUpsertRequest, includeName bool) url.Values {
	values := url.Values{}

	if includeName || payload.Name != "" {
		values.Set("name", payload.Name)
	}

	setOptionalBool(values, "auto_deploy", payload.AutoDeploy)
	setOptionalString(values, "type", payload.GroupType)
	setOptionalString(values, "install_type", payload.InstallType)

	if payload.Priority != nil {
		values.Set("priority", strconv.FormatInt(*payload.Priority, 10))
	}

	setOptionalBool(values, "app_track_location", payload.AppTrackLocation)

	return values
}

// setOptionalBool adds the given key to the query values when the pointer contains a value.
func setOptionalBool(values url.Values, key string, value *bool) {
	if value != nil {
		values.Set(key, strconv.FormatBool(*value))
	}
}

// setOptionalString adds the given key to the query values when the pointer contains a non-empty string.
func setOptionalString(values url.Values, key string, value *string) {
	if value != nil && *value != "" {
		values.Set(key, *value)
	}
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
	// For standard groups, set to null since API doesn't return it
	if response.Data.Attributes.Type == "munki" {
		if response.Data.Attributes.InstallType != "" {
			model.InstallType = types.StringValue(response.Data.Attributes.InstallType)
		} else {
			model.InstallType = types.StringNull()
		}
	} else {
		// For standard groups, always set to null
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

// setElementsToStringSlice converts a types.Set to a []string slice
func setElementsToStringSlice(set types.Set) []string {
	if set.IsNull() || set.IsUnknown() {
		return []string{}
	}

	elements := set.Elements()
	result := make([]string, 0, len(elements))
	for _, element := range elements {
		stringElement, ok := element.(types.String)
		if !ok || stringElement.IsNull() || stringElement.IsUnknown() {
			continue
		}

		result = append(result, stringElement.ValueString())
	}
	return result
}

// assignObjectsToGroup assigns multiple objects to an assignment group
// Used during Create operations to assign apps, profiles, groups, or devices
func assignObjectsToGroup(
	ctx context.Context,
	client *simplemdm.Client,
	groupID string,
	objects types.Set,
	objectType string,
	removeOthers bool,
) error {
	if objects.IsNull() || objects.IsUnknown() {
		return nil
	}

	for _, objectID := range objects.Elements() {
		idString := objectID.(types.String).ValueString()

		var err error
		if objectType == "devices" {
			// Devices use special assignment function with removeOthers parameter
			err = assignmentGroupAssignDevice(ctx, client, groupID, idString, removeOthers)
		} else {
			// Apps, profiles, and device_groups use standard assignment
			err = client.AssignmentGroupAssignObject(groupID, idString, objectType)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

// updateAssignmentGroupObjects updates assignments by computing diff and applying changes
// Used during Update operations to sync state with plan
func updateAssignmentGroupObjects(
	ctx context.Context,
	client *simplemdm.Client,
	groupID string,
	stateObjects types.Set,
	planObjects types.Set,
	objectType string,
	removeOthers bool,
) error {
	// When the plan set is null or unknown, Terraform cannot determine a desired
	// target state for the relationship. In these cases we must leave the
	// existing assignments untouched and bail out early regardless of what the
	// current state contains.
	if planObjects.IsNull() || planObjects.IsUnknown() {
		return nil
	}

	if stateObjects.IsNull() || stateObjects.IsUnknown() {
		stateObjects = types.SetNull(types.StringType)
	}

	// Convert sets to string slices
	stateSlice := setElementsToStringSlice(stateObjects)
	planSlice := setElementsToStringSlice(planObjects)

	// Compute diff
	toAdd, toRemove := diffFunction(stateSlice, planSlice)

	// Add new objects
	for _, objectID := range toAdd {
		var err error
		if objectType == "devices" {
			err = assignmentGroupAssignDevice(ctx, client, groupID, objectID, removeOthers)
		} else {
			err = client.AssignmentGroupAssignObject(groupID, objectID, objectType)
		}
		if err != nil {
			return err
		}
	}

	// Remove old objects
	for _, objectID := range toRemove {
		err := client.AssignmentGroupUnAssignObject(groupID, objectID, objectType)
		if err != nil {
			return err
		}
	}

	return nil
}

// diffFunction computes the difference between state and plan lists
// Returns items to add and items to remove
// Optimized to O(n) complexity using map lookups instead of nested loops
func diffFunction(state []string, plan []string) (add []string, remove []string) {
	// Create map of state items for O(1) lookups
	stateMap := make(map[string]bool, len(state))
	for _, s := range state {
		stateMap[s] = true
	}

	// Create map of plan items and identify additions
	planMap := make(map[string]bool, len(plan))
	for _, p := range plan {
		planMap[p] = true
		if !stateMap[p] {
			add = append(add, p)
		}
	}

	// Identify removals
	for _, s := range state {
		if !planMap[s] {
			remove = append(remove, s)
		}
	}

	return add, remove
}

// preservePlannedRelationships handles eventual consistency by preserving planned values
// when the API doesn't immediately return assigned relationships
func preservePlannedRelationships(
	model *assignment_groupResourceModel,
	plannedApps, plannedProfiles, plannedGroups, plannedDevices types.Set,
	apiReturnedApps, apiReturnedProfiles, apiReturnedGroups, apiReturnedDevices bool,
) {
	// Restore planned relationship values if they were set but API returned empty
	// This prevents "planned X but got Y" errors due to API eventual consistency
	if !plannedApps.IsNull() && !plannedApps.IsUnknown() && !apiReturnedApps {
		model.Apps = plannedApps
	}
	if !plannedProfiles.IsNull() && !plannedProfiles.IsUnknown() && !apiReturnedProfiles {
		model.Profiles = plannedProfiles
	}
	if !plannedGroups.IsNull() && !plannedGroups.IsUnknown() && !apiReturnedGroups {
		model.Groups = plannedGroups
	}
	if !plannedDevices.IsNull() && !plannedDevices.IsUnknown() && !apiReturnedDevices {
		model.Devices = plannedDevices
	}
}
