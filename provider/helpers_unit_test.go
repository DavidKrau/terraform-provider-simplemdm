package provider

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestManagedConfigJSONParsing tests JSON parsing for managed configs
func TestManagedConfigJSONParsing(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"id": 123,
				"attributes": {
					"key": "test_key",
					"value": "test_value",
					"value_type": "string"
				}
			}
		]
	}`

	var response managedConfigListResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(response.Data) != 1 {
		t.Errorf("expected 1 config, got %d", len(response.Data))
	}

	if response.Data[0].ID != 123 {
		t.Errorf("expected ID 123, got %d", response.Data[0].ID)
	}

	if response.Data[0].Attributes.Key != "test_key" {
		t.Errorf("expected key 'test_key', got '%s'", response.Data[0].Attributes.Key)
	}
}

// TestManagedConfigSingleItemResponse tests single item response parsing
func TestManagedConfigSingleItemResponse(t *testing.T) {
	jsonData := `{
		"data": {
			"id": 456,
			"attributes": {
				"key": "another_key",
				"value": "another_value",
				"value_type": "boolean"
			}
		}
	}`

	var response managedConfigItemResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if response.Data.ID != 456 {
		t.Errorf("expected ID 456, got %d", response.Data.ID)
	}

	if response.Data.Attributes.ValueType != "boolean" {
		t.Errorf("expected value_type 'boolean', got '%s'", response.Data.Attributes.ValueType)
	}
}

// TestAssignmentGroupResponseParsing tests assignment group JSON parsing
func TestAssignmentGroupResponseParsing(t *testing.T) {
	jsonData := `{
		"data": {
			"id": 789,
			"type": "assignment_groups",
			"attributes": {
				"name": "Test Group",
				"auto_deploy": true,
				"type": "standard",
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2025-01-02T00:00:00Z",
				"device_count": 10,
				"group_count": 5,
				"priority": 7,
				"app_track_location": false
			},
			"relationships": {
				"apps": {
					"data": [
						{"id": 1, "type": "apps"},
						{"id": 2, "type": "apps"}
					]
				},
				"profiles": {
					"data": [
						{"id": 3, "type": "profiles"}
					]
				},
				"devices": {
					"data": []
				},
				"device_groups": {
					"data": []
				}
			}
		}
	}`

	var response assignmentGroupResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if response.Data.ID != 789 {
		t.Errorf("expected ID 789, got %d", response.Data.ID)
	}

	if response.Data.Attributes.Name != "Test Group" {
		t.Errorf("expected name 'Test Group', got '%s'", response.Data.Attributes.Name)
	}

	if !response.Data.Attributes.AutoDeploy {
		t.Errorf("expected auto_deploy to be true")
	}

	if response.Data.Attributes.DeviceCount != 10 {
		t.Errorf("expected device_count 10, got %d", response.Data.Attributes.DeviceCount)
	}

	if len(response.Data.Relationships.Apps.Data) != 2 {
		t.Errorf("expected 2 apps, got %d", len(response.Data.Relationships.Apps.Data))
	}

	if len(response.Data.Relationships.Profiles.Data) != 1 {
		t.Errorf("expected 1 profile, got %d", len(response.Data.Relationships.Profiles.Data))
	}
}

// TestAssignmentGroupWithNullOptionalFields tests parsing with null optional fields
func TestAssignmentGroupWithNullOptionalFields(t *testing.T) {
	jsonData := `{
		"data": {
			"id": 100,
			"type": "assignment_groups",
			"attributes": {
				"name": "Basic Group",
				"auto_deploy": false,
				"type": "standard",
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2025-01-01T00:00:00Z",
				"device_count": 0,
				"group_count": 0
			},
			"relationships": {
				"apps": {"data": []},
				"profiles": {"data": []},
				"devices": {"data": []},
				"device_groups": {"data": []}
			}
		}
	}`

	var response assignmentGroupResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if response.Data.Attributes.Priority != nil {
		t.Errorf("expected priority to be nil, got %v", *response.Data.Attributes.Priority)
	}

	if response.Data.Attributes.AppTrackLocation != nil {
		t.Errorf("expected app_track_location to be nil, got %v", *response.Data.Attributes.AppTrackLocation)
	}
}

// TestScriptJobResponseParsing tests script job JSON parsing
func TestScriptJobResponseParsing(t *testing.T) {
	jsonData := `{
		"data": {
			"id": 999,
			"type": "script_jobs",
			"attributes": {
				"status": "completed",
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2025-01-02T00:00:00Z"
			},
			"relationships": {
				"script": {
					"data": {
						"id": 111,
						"type": "scripts"
					}
				},
				"assignment_group": {
					"data": {
						"id": 222,
						"type": "assignment_groups"
					}
				}
			}
		}
	}`

	var response scriptJobResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if response.Data.ID != 999 {
		t.Errorf("expected ID 999, got %d", response.Data.ID)
	}

	if response.Data.Attributes.Status != "completed" {
		t.Errorf("expected status 'completed', got '%s'", response.Data.Attributes.Status)
	}

	if response.Data.Relationships.Script.Data == nil {
		t.Fatal("expected script relationship to be non-nil")
	}

	if response.Data.Relationships.Script.Data.ID != 111 {
		t.Errorf("expected script ID 111, got %d", response.Data.Relationships.Script.Data.ID)
	}
}

// TestScriptJobDetailsResponseParsing tests detailed script job response
func TestScriptJobDetailsResponseParsing(t *testing.T) {
	jsonData := `{
		"data": {
			"id": 500,
			"attributes": {
				"script_name": "Test Script",
				"job_name": "Test Job",
				"content": "#!/bin/bash\necho test",
				"job_id": "job-uuid-123",
				"variable_support": true,
				"status": "completed",
				"pending_count": 0,
				"success_count": 5,
				"errored_count": 1,
				"custom_attribute_regex": ".*",
				"created_by": "admin@example.com",
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2025-01-02T00:00:00Z"
			},
			"relationships": {
				"custom_attribute": {
					"data": {
						"id": "attr-123"
					}
				},
				"device": {
					"data": [
						{
							"id": 1,
							"status": "success",
							"status_code": "0",
							"response": "output"
						},
						{
							"id": 2,
							"status": "errored",
							"status_code": "1",
							"response": null
						}
					]
				}
			}
		}
	}`

	var response scriptJobDetailsResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if response.Data.ID != 500 {
		t.Errorf("expected ID 500, got %d", response.Data.ID)
	}

	if response.Data.Attributes.ScriptName != "Test Script" {
		t.Errorf("expected script_name 'Test Script', got '%s'", response.Data.Attributes.ScriptName)
	}

	if response.Data.Attributes.VariableSupport != true {
		t.Errorf("expected variable_support to be true")
	}

	if response.Data.Attributes.SuccessCount != 5 {
		t.Errorf("expected success_count 5, got %d", response.Data.Attributes.SuccessCount)
	}

	if response.Data.Attributes.ErroredCount != 1 {
		t.Errorf("expected errored_count 1, got %d", response.Data.Attributes.ErroredCount)
	}

	if len(response.Data.Relationships.Device.Data) != 2 {
		t.Errorf("expected 2 devices, got %d", len(response.Data.Relationships.Device.Data))
	}
}

// TestBuildStringSetFromRelationshipItemsEmpty tests empty relationship conversion
func TestBuildStringSetFromRelationshipItemsEmpty(t *testing.T) {
	items := []assignmentGroupRelationshipItem{}
	result := buildStringSetFromRelationshipItems(items)

	if result.IsNull() {
		t.Error("result should not be null for empty slice")
	}

	if len(result.Elements()) != 0 {
		t.Errorf("expected 0 elements, got %d", len(result.Elements()))
	}
}

// TestBuildStringSetFromRelationshipItemsMultiple tests multiple items
func TestBuildStringSetFromRelationshipItemsMultiple(t *testing.T) {
	items := []assignmentGroupRelationshipItem{
		{ID: 123, Type: "apps"},
		{ID: 456, Type: "apps"},
		{ID: 789, Type: "apps"},
	}
	result := buildStringSetFromRelationshipItems(items)

	elements := result.Elements()
	if len(elements) != 3 {
		t.Errorf("expected 3 elements, got %d", len(elements))
	}

	// Verify all IDs are present
	expectedIDs := map[string]bool{"123": false, "456": false, "789": false}
	for _, elem := range elements {
		id := elem.(types.String).ValueString()
		if _, exists := expectedIDs[id]; exists {
			expectedIDs[id] = true
		}
	}

	for id, found := range expectedIDs {
		if !found {
			t.Errorf("expected to find ID %s in set", id)
		}
	}
}

// TestFlattenScriptJobWithRelationships tests flattening with all relationships
func TestFlattenScriptJobWithRelationships(t *testing.T) {
	response := &scriptJobResponse{
		Data: scriptJobData{
			ID:   123,
			Type: "script_jobs",
			Attributes: scriptJobAttributes{
				Status:    "completed",
				CreatedAt: "2025-01-01T00:00:00Z",
				UpdatedAt: "2025-01-02T00:00:00Z",
			},
			Relationships: scriptJobRelationships{
				Script: scriptJobRelationshipItem{
					Data: &struct {
						ID   int    `json:"id"`
						Type string `json:"type"`
					}{
						ID:   456,
						Type: "scripts",
					},
				},
				AssignmentGroup: scriptJobRelationshipItem{
					Data: &struct {
						ID   int    `json:"id"`
						Type string `json:"type"`
					}{
						ID:   789,
						Type: "assignment_groups",
					},
				},
			},
		},
	}

	result := flattenScriptJob(response)

	if result.ID != 123 {
		t.Errorf("expected ID 123, got %d", result.ID)
	}

	if result.ScriptID == nil || *result.ScriptID != 456 {
		t.Errorf("expected ScriptID 456, got %v", result.ScriptID)
	}

	if result.AssignmentGroupID == nil || *result.AssignmentGroupID != 789 {
		t.Errorf("expected AssignmentGroupID 789, got %v", result.AssignmentGroupID)
	}

	if result.Status != "completed" {
		t.Errorf("expected status 'completed', got '%s'", result.Status)
	}
}

// TestFlattenScriptJobWithoutRelationships tests flattening without relationships
func TestFlattenScriptJobWithoutRelationships(t *testing.T) {
	response := &scriptJobResponse{
		Data: scriptJobData{
			ID:   999,
			Type: "script_jobs",
			Attributes: scriptJobAttributes{
				Status:    "pending",
				CreatedAt: "2025-01-01T00:00:00Z",
				UpdatedAt: "2025-01-01T00:00:00Z",
			},
			Relationships: scriptJobRelationships{
				Script: scriptJobRelationshipItem{
					Data: nil,
				},
				AssignmentGroup: scriptJobRelationshipItem{
					Data: nil,
				},
			},
		},
	}

	result := flattenScriptJob(response)

	if result.ID != 999 {
		t.Errorf("expected ID 999, got %d", result.ID)
	}

	if result.ScriptID != nil {
		t.Errorf("expected ScriptID to be nil, got %v", *result.ScriptID)
	}

	if result.AssignmentGroupID != nil {
		t.Errorf("expected AssignmentGroupID to be nil, got %v", *result.AssignmentGroupID)
	}
}
