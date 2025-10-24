package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNewAppResourceModelFromAPI_AllFields(t *testing.T) {
	ctx := context.Background()
	response := &appAPIResponse{}
	response.Data.ID = 42
	response.Data.Attributes.Name = "Marketing App"
	response.Data.Attributes.BundleIdentifier = "com.example.marketing"
	itunesID := 123456789
	response.Data.Attributes.ITunesStoreID = &itunesID
	response.Data.Attributes.AppType = "app store"
	response.Data.Attributes.InstallationChannels = []string{"assignment_groups"}
	response.Data.Attributes.PlatformSupport = "ios"
	response.Data.Attributes.ProcessingStatus = "ready"
	response.Data.Attributes.Version = "7.1"
	response.Data.Attributes.DeployTo = "outdated"
	response.Data.Attributes.Status = "deployed"
	response.Data.Attributes.CreatedAt = "2025-01-02T03:04:05Z"
	response.Data.Attributes.UpdatedAt = "2025-02-03T04:05:06Z"

	model, diags := newAppResourceModelFromAPI(ctx, response)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if model.ID.IsNull() || model.ID.ValueString() != "42" {
		t.Fatalf("expected ID to be set, got %v", model.ID)
	}
	if model.Name.ValueString() != "Marketing App" {
		t.Fatalf("expected Name to be set")
	}
	if model.AppStoreId.ValueString() != "123456789" {
		t.Fatalf("expected AppStoreId to be set")
	}
	if model.BundleId.ValueString() != "com.example.marketing" {
		t.Fatalf("expected BundleId to be set")
	}
	if model.DeployTo.ValueString() != "outdated" {
		t.Fatalf("expected DeployTo to be set")
	}
	if model.Status.ValueString() != "deployed" {
		t.Fatalf("expected Status to be set")
	}
	if model.AppType.ValueString() != "app store" {
		t.Fatalf("expected AppType to be set")
	}
	if model.Version.ValueString() != "7.1" {
		t.Fatalf("expected Version to be set")
	}
	if model.PlatformSupport.ValueString() != "ios" {
		t.Fatalf("expected PlatformSupport to be set")
	}
	if model.ProcessingStatus.ValueString() != "ready" {
		t.Fatalf("expected ProcessingStatus to be set")
	}
	if model.CreatedAt.ValueString() != "2025-01-02T03:04:05Z" {
		t.Fatalf("expected CreatedAt to be set")
	}
	if model.UpdatedAt.ValueString() != "2025-02-03T04:05:06Z" {
		t.Fatalf("expected UpdatedAt to be set")
	}
	if model.InstallationChannels.IsNull() {
		t.Fatalf("expected InstallationChannels to be set")
	}
	values := model.InstallationChannels.Elements()
	if len(values) != 1 || values[0].(types.String).ValueString() != "assignment_groups" {
		t.Fatalf("unexpected installation channels: %v", values)
	}
}

func TestNewAppResourceModelFromAPI_PartialData(t *testing.T) {
	ctx := context.Background()
	response := &appAPIResponse{}
	response.Data.ID = 99

	model, diags := newAppResourceModelFromAPI(ctx, response)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	if !model.Name.IsNull() {
		t.Fatalf("expected Name to remain null when not provided")
	}
	if !model.AppStoreId.IsNull() {
		t.Fatalf("expected AppStoreId to remain null when not provided")
	}
	if !model.BundleId.IsNull() {
		t.Fatalf("expected BundleId to remain null when not provided")
	}
	if model.DeployTo.ValueString() != "none" {
		t.Fatalf("expected DeployTo to default to none, got %q", model.DeployTo.ValueString())
	}
	if !model.InstallationChannels.IsNull() {
		t.Fatalf("expected InstallationChannels to remain null when not provided")
	}
}
