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
	Name        string `json:"name"`
	AutoDeploy  bool   `json:"auto_deploy"`
	Type        string `json:"type"`
	InstallType string `json:"install_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeviceCount int    `json:"device_count"`
	GroupCount  int    `json:"group_count"`
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
	if len(items) == 0 {
		return types.SetNull(types.StringType)
	}

	values := make([]attr.Value, len(items))
	for i, item := range items {
		values[i] = types.StringValue(strconv.Itoa(item.ID))
	}

	return types.SetValueMust(types.StringType, values)
}
