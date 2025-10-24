package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	simplemdm "github.com/DavidKrau/simplemdm-go-client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type scriptJobDeviceDetail struct {
	ID         string
	Status     string
	StatusCode *string
	Response   *string
}

type scriptJobDetailsData struct {
	ID                   string
	ScriptName           string
	JobName              string
	JobIdentifier        string
	Status               string
	PendingCount         int64
	SuccessCount         int64
	ErroredCount         int64
	Content              string
	VariableSupport      bool
	CreatedBy            string
	CreatedAt            string
	UpdatedAt            string
	CustomAttribute      string
	CustomAttributeRegex string
	Devices              []scriptJobDeviceDetail
}

var scriptJobDeviceAttrTypes = map[string]attr.Type{
	"id":          types.StringType,
	"status":      types.StringType,
	"status_code": types.StringType,
	"response":    types.StringType,
}

type scriptJobDetailsResponse struct {
	Data struct {
		ID         int `json:"id"`
		Attributes struct {
			ScriptName           string `json:"script_name"`
			JobName              string `json:"job_name"`
			Content              string `json:"content"`
			JobID                string `json:"job_id"`
			VariableSupport      bool   `json:"variable_support"`
			Status               string `json:"status"`
			PendingCount         int    `json:"pending_count"`
			SuccessCount         int    `json:"success_count"`
			ErroredCount         int    `json:"errored_count"`
			CustomAttributeRegex string `json:"custom_attribute_regex"`
			CreatedBy            string `json:"created_by"`
			CreatedAt            string `json:"created_at"`
			UpdatedAt            string `json:"updated_at"`
		} `json:"attributes"`
		Relationships struct {
			CustomAttribute struct {
				Data *struct {
					ID string `json:"id"`
				} `json:"data"`
			} `json:"custom_attribute"`
			Device struct {
				Data []struct {
					ID         int     `json:"id"`
					Status     string  `json:"status"`
					StatusCode *string `json:"status_code"`
					Response   *string `json:"response"`
				} `json:"data"`
			} `json:"device"`
		} `json:"relationships"`
	} `json:"data"`
}

func fetchScriptJobDetails(ctx context.Context, client *simplemdm.Client, id string) (*scriptJobDetailsData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://%s/api/v1/script_jobs/%s", client.HostName, id), nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	var payload scriptJobDetailsResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	details := &scriptJobDetailsData{
		ID:                   strconv.Itoa(payload.Data.ID),
		ScriptName:           payload.Data.Attributes.ScriptName,
		JobName:              payload.Data.Attributes.JobName,
		JobIdentifier:        payload.Data.Attributes.JobID,
		Status:               payload.Data.Attributes.Status,
		PendingCount:         int64(payload.Data.Attributes.PendingCount),
		SuccessCount:         int64(payload.Data.Attributes.SuccessCount),
		ErroredCount:         int64(payload.Data.Attributes.ErroredCount),
		Content:              payload.Data.Attributes.Content,
		VariableSupport:      payload.Data.Attributes.VariableSupport,
		CreatedBy:            payload.Data.Attributes.CreatedBy,
		CreatedAt:            payload.Data.Attributes.CreatedAt,
		UpdatedAt:            payload.Data.Attributes.UpdatedAt,
		CustomAttributeRegex: payload.Data.Attributes.CustomAttributeRegex,
	}

	if payload.Data.Relationships.CustomAttribute.Data != nil {
		details.CustomAttribute = payload.Data.Relationships.CustomAttribute.Data.ID
	}

	for _, device := range payload.Data.Relationships.Device.Data {
		deviceDetail := scriptJobDeviceDetail{
			ID:     strconv.Itoa(device.ID),
			Status: device.Status,
		}

		if device.StatusCode != nil && *device.StatusCode != "" {
			statusCode := *device.StatusCode
			deviceDetail.StatusCode = &statusCode
		}

		if device.Response != nil && *device.Response != "" {
			response := *device.Response
			deviceDetail.Response = &response
		}

		details.Devices = append(details.Devices, deviceDetail)
	}

	return details, nil
}

func scriptJobDevicesListValue(ctx context.Context, devices []scriptJobDeviceDetail) (types.List, diag.Diagnostics) {
	if len(devices) == 0 {
		return types.ListValue(types.ObjectType{AttrTypes: scriptJobDeviceAttrTypes}, []attr.Value{})
	}

	values := make([]attr.Value, 0, len(devices))
	var diags diag.Diagnostics

	for _, device := range devices {
		attrs := map[string]attr.Value{
			"id":     types.StringValue(device.ID),
			"status": types.StringValue(device.Status),
		}

		if device.StatusCode != nil {
			attrs["status_code"] = types.StringValue(*device.StatusCode)
		} else {
			attrs["status_code"] = types.StringNull()
		}

		if device.Response != nil {
			attrs["response"] = types.StringValue(*device.Response)
		} else {
			attrs["response"] = types.StringNull()
		}

		obj, d := types.ObjectValue(scriptJobDeviceAttrTypes, attrs)
		diags.Append(d...)
		values = append(values, obj)
	}

	list, d := types.ListValue(types.ObjectType{AttrTypes: scriptJobDeviceAttrTypes}, values)
	diags.Append(d...)

	return list, diags
}

func isNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "404")
}
