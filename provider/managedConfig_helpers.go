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
	ID         int                     `json:"id"`
	Attributes managedConfigAttributes `json:"attributes"`
}

type managedConfigListResponse struct {
	Data []managedConfigAPIResource `json:"data"`
}

type managedConfigItemResponse struct {
	Data managedConfigAPIResource `json:"data"`
}

var errManagedConfigNotFound = errors.New("managed config not found")

func fetchManagedConfig(ctx context.Context, client *simplemdm.Client, appID, configID string) (*managedConfigAPIResource, error) {
	if client == nil {
		return nil, errors.New("simplemdm client is not configured")
	}

	url := fmt.Sprintf("https://%s/api/v1/apps/%s/managed_configs", client.HostName, appID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	var response managedConfigListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	for _, item := range response.Data {
		if strconv.Itoa(item.ID) == configID {
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

	type requester func(*http.Request) ([]byte, error)

	attempts := []struct {
		fn       requester
		expected string
	}{
		{fn: client.RequestResponse202, expected: "202"},
		{fn: client.RequestResponse200, expected: "200"},
		{fn: client.RequestResponse204, expected: "204"},
	}

	for _, attempt := range attempts {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
		if err != nil {
			return err
		}

		if _, err = attempt.fn(req); err != nil {
			if strings.Contains(err.Error(), "non "+attempt.expected+" status code") {
				continue
			}

			return err
		}

		return nil
	}

	return fmt.Errorf("failed to push managed config updates for app %s", appID)
}
