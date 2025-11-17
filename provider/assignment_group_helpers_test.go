package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSetElementsToStringSliceHandlesNullAndUnknown(t *testing.T) {
	nullSet := types.SetNull(types.StringType)
	if result := setElementsToStringSlice(nullSet); len(result) != 0 {
		t.Fatalf("expected empty slice for null set, got %v", result)
	}

	unknownSet := types.SetUnknown(types.StringType)
	if result := setElementsToStringSlice(unknownSet); len(result) != 0 {
		t.Fatalf("expected empty slice for unknown set, got %v", result)
	}
}

func TestSetElementsToStringSliceFiltersInvalidElements(t *testing.T) {
	set := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("alpha"),
		types.StringUnknown(),
		types.StringNull(),
		types.StringValue("beta"),
	})

	result := setElementsToStringSlice(set)
	if len(result) != 2 {
		t.Fatalf("expected two valid values, got %d", len(result))
	}

	if result[0] != "alpha" || result[1] != "beta" {
		t.Fatalf("unexpected values returned: %v", result)
	}
}

func TestUpdateAssignmentGroupObjectsSkipsWhenNoData(t *testing.T) {
	err := updateAssignmentGroupObjects(
		context.Background(),
		nil,
		"group-id",
		types.SetUnknown(types.StringType),
		types.SetNull(types.StringType),
		"apps",
		false,
	)
	if err != nil {
		t.Fatalf("expected no error when sets are empty/unknown, got %v", err)
	}
}

func TestUpdateAssignmentGroupObjectsSkipsWhenPlanUnknown(t *testing.T) {
	state := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("123"),
	})

	err := updateAssignmentGroupObjects(
		context.Background(),
		nil,
		"group-id",
		state,
		types.SetUnknown(types.StringType),
		"apps",
		false,
	)
	if err != nil {
		t.Fatalf("expected no error when plan set is unknown, got %v", err)
	}
}
