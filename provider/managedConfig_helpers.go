package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
)

type managedConfigAttributes struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ValueType string `json:"value_type"`
}

type managedConfigAPIResource struct {
	Type       string                  `json:"type"`
	ID         int                     `json:"id"`
	Attributes managedConfigAttributes `json:"attributes"`
}

type managedConfigListResponse struct {
	Data    []managedConfigAPIResource `json:"data"`
	HasMore bool                       `json:"has_more"`
}

type managedConfigItemResponse struct {
	Data managedConfigAPIResource `json:"data"`
}

var errManagedConfigNotFound = errors.New("managed config not found")

func fetchManagedConfig(ctx context.Context, client *simplemdm.Client, appID, configID string) (*managedConfigAPIResource, error) {
	if client == nil {
		return nil, errors.New("simplemdm client is not configured")
	}

	// Fetch all configs with pagination support
	allConfigs, err := fetchAllManagedConfigs(ctx, client, appID)
	if err != nil {
		return nil, err
	}

	// Find the specific config
	for _, item := range allConfigs {
		if strconv.Itoa(item.ID) == configID {
			// Validate type field
			if item.Type != "" && item.Type != "managed_config" {
				return nil, fmt.Errorf("unexpected resource type: %s (expected managed_config)", item.Type)
			}
			return &item, nil
		}
	}

	return nil, errManagedConfigNotFound
}

func createManagedConfig(ctx context.Context, client *simplemdm.Client, appID, key, value, valueType string) (*managedConfigAPIResource, error) {
	url := fmt.Sprintf("https://%s/api/v1/apps/%s/managed_configs", client.HostName, appID)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	fields := map[string]string{
		"key":        key,
		"value":      value,
		"value_type": valueType,
	}

	for name, fieldValue := range fields {
		if err := writer.WriteField(name, fieldValue); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	body, err := client.RequestResponse201(req)
	if err != nil {
		return nil, err
	}

	var response managedConfigItemResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	// Validate type field
	if response.Data.Type != "" && response.Data.Type != "managed_config" {
		return nil, fmt.Errorf("unexpected resource type: %s (expected managed_config)", response.Data.Type)
	}

	return &response.Data, nil
}

func deleteManagedConfig(ctx context.Context, client *simplemdm.Client, appID, configID string) error {
	url := fmt.Sprintf("https://%s/api/v1/apps/%s/managed_configs/%s", client.HostName, appID, configID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	_, err = client.RequestResponse204(req)
	return err
}

func pushManagedConfigUpdates(ctx context.Context, client *simplemdm.Client, appID string) error {
	url := fmt.Sprintf("https://%s/api/v1/apps/%s/managed_configs/push", client.HostName, appID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create push request: %w", err)
	}

	// API returns 202 Accepted for async push operations
	if _, err = client.RequestResponse202(req); err != nil {
		return fmt.Errorf("failed to push managed config updates: %w", err)
	}

	return nil
}
