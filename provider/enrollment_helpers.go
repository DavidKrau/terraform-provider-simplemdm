package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"

	"github.com/DavidKrau/simplemdm-go-client"
)

type enrollmentResponse struct {
	Data enrollmentData `json:"data"`
}

type enrollmentData struct {
	ID            int                     `json:"id"`
	Type          string                  `json:"type"`
	Attributes    enrollmentAttributes    `json:"attributes"`
	Relationships enrollmentRelationships `json:"relationships"`
}

type enrollmentAttributes struct {
	URL            *string `json:"url"`
	UserEnrollment bool    `json:"user_enrollment"`
	WelcomeScreen  bool    `json:"welcome_screen"`
	Authentication bool    `json:"authentication"`
}

type enrollmentRelationships struct {
	DeviceGroup enrollmentRelationshipItem  `json:"device_group"`
	Device      *enrollmentRelationshipItem `json:"device,omitempty"`
}

type enrollmentRelationshipItem struct {
	Data *struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
}

type enrollmentUpsertRequest struct {
	DeviceGroupID  string
	UserEnrollment *bool
	WelcomeScreen  *bool
	Authentication *bool
}

type enrollmentFlat struct {
	ID             int
	URL            *string
	UserEnrollment bool
	WelcomeScreen  bool
	Authentication bool
	DeviceGroupID  *int
	DeviceID       *int
}

func fetchEnrollment(ctx context.Context, client *simplemdm.Client, id string) (*enrollmentResponse, error) {
	url := fmt.Sprintf("https://%s/api/v1/enrollments/%s", client.HostName, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, err
	}

	var enrollment enrollmentResponse
	if err := json.Unmarshal(body, &enrollment); err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func createEnrollment(ctx context.Context, client *simplemdm.Client, payload enrollmentUpsertRequest) (*enrollmentResponse, error) {
	endpoint := fmt.Sprintf("https://%s/api/v1/enrollments", client.HostName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("device_group_id", payload.DeviceGroupID)

	if payload.UserEnrollment != nil {
		q.Add("user_enrollment", strconv.FormatBool(*payload.UserEnrollment))
	}

	if payload.WelcomeScreen != nil {
		q.Add("welcome_screen", strconv.FormatBool(*payload.WelcomeScreen))
	}

	if payload.Authentication != nil {
		q.Add("authentication", strconv.FormatBool(*payload.Authentication))
	}

	req.URL.RawQuery = q.Encode()

	body, err := client.RequestResponse201(req)
	if err != nil {
		return nil, err
	}

	var enrollment enrollmentResponse
	if err := json.Unmarshal(body, &enrollment); err != nil {
		return nil, err
	}

	return &enrollment, nil
}

func deleteEnrollment(ctx context.Context, client *simplemdm.Client, id string) error {
	endpoint := fmt.Sprintf("https://%s/api/v1/enrollments/%s", client.HostName, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	_, err = client.RequestResponse204(req)
	return err
}

func sendEnrollmentInvitation(ctx context.Context, client *simplemdm.Client, id string, contact string) error {
	endpoint := fmt.Sprintf("https://%s/api/v1/enrollments/%s/invitations", client.HostName, id)

	form := neturl.Values{}
	form.Set("contact", contact)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.RequestResponse200(req)
	return err
}

func flattenEnrollment(response *enrollmentResponse) enrollmentFlat {
	var deviceGroupID *int
	if response.Data.Relationships.DeviceGroup.Data != nil {
		value := response.Data.Relationships.DeviceGroup.Data.ID
		deviceGroupID = &value
	}

	var deviceID *int
	if response.Data.Relationships.Device != nil && response.Data.Relationships.Device.Data != nil {
		value := response.Data.Relationships.Device.Data.ID
		deviceID = &value
	}

	return enrollmentFlat{
		ID:             response.Data.ID,
		URL:            response.Data.Attributes.URL,
		UserEnrollment: response.Data.Attributes.UserEnrollment,
		WelcomeScreen:  response.Data.Attributes.WelcomeScreen,
		Authentication: response.Data.Attributes.Authentication,
		DeviceGroupID:  deviceGroupID,
		DeviceID:       deviceID,
	}
}
