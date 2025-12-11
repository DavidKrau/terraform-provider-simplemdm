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
	DeviceGroup      *enrollmentRelationshipItem `json:"device_group,omitempty"`
	AssignmentGroup  *enrollmentRelationshipItem `json:"assignment_group,omitempty"`
	Device           *enrollmentRelationshipItem `json:"device,omitempty"`
}

type enrollmentRelationshipItem struct {
	Data *struct {
		ID   int    `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
}

type enrollmentUpsertRequest struct {
	DeviceGroupID     string
	AssignmentGroupID string
	UserEnrollment    *bool
	WelcomeScreen     *bool
	Authentication    *bool
}

type enrollmentFlat struct {
	ID                int
	URL               *string
	UserEnrollment    bool
	WelcomeScreen     bool
	Authentication    bool
	DeviceGroupID     *int
	AssignmentGroupID *int
	DeviceID          *int
}

func fetchEnrollment(ctx context.Context, client *simplemdm.Client, id string) (*enrollmentResponse, error) {
	url := fmt.Sprintf("https://%s/api/v1/enrollments/%s", client.HostName, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for enrollment %s: %w", id, err)
	}

	body, err := client.RequestResponse200(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch enrollment %s: %w", id, err)
	}

	var enrollment enrollmentResponse
	if err := json.Unmarshal(body, &enrollment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal enrollment %s response: %w", id, err)
	}

	return &enrollment, nil
}

// createEnrollment creates a new enrollment using the SimpleMDM API.
// NOTE: This endpoint (POST /api/v1/enrollments) is not documented in the official
// SimpleMDM API specification but is supported by the API. This implementation may
// be subject to change without notice if the API changes.
func createEnrollment(ctx context.Context, client *simplemdm.Client, payload enrollmentUpsertRequest) (*enrollmentResponse, error) {
	endpoint := fmt.Sprintf("https://%s/api/v1/enrollments", client.HostName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create enrollment request: %w", err)
	}

	q := req.URL.Query()
	
	// Support both legacy device_group_id and modern assignment_group_id
	if payload.DeviceGroupID != "" {
		q.Add("device_group_id", payload.DeviceGroupID)
	}
	
	if payload.AssignmentGroupID != "" {
		q.Add("assignment_group_id", payload.AssignmentGroupID)
	}

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
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}

	var enrollment enrollmentResponse
	if err := json.Unmarshal(body, &enrollment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal enrollment creation response: %w", err)
	}

	return &enrollment, nil
}

func deleteEnrollment(ctx context.Context, client *simplemdm.Client, id string) error {
	endpoint := fmt.Sprintf("https://%s/api/v1/enrollments/%s", client.HostName, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request for enrollment %s: %w", id, err)
	}

	_, err = client.RequestResponse204(req)
	if err != nil {
		return fmt.Errorf("failed to delete enrollment %s: %w", id, err)
	}
	return nil
}

func sendEnrollmentInvitation(ctx context.Context, client *simplemdm.Client, id string, contact string) error {
	endpoint := fmt.Sprintf("https://%s/api/v1/enrollments/%s/invitations", client.HostName, id)

	form := neturl.Values{}
	form.Set("contact", contact)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create invitation request for enrollment %s: %w", id, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.RequestResponse200(req)
	if err != nil {
		return fmt.Errorf("failed to send invitation for enrollment %s: %w", id, err)
	}
	return nil
}

func flattenEnrollment(response *enrollmentResponse) enrollmentFlat {
	var deviceGroupID *int
	if response.Data.Relationships.DeviceGroup != nil &&
	   response.Data.Relationships.DeviceGroup.Data != nil {
		value := response.Data.Relationships.DeviceGroup.Data.ID
		deviceGroupID = &value
	}

	var assignmentGroupID *int
	if response.Data.Relationships.AssignmentGroup != nil &&
	   response.Data.Relationships.AssignmentGroup.Data != nil {
		value := response.Data.Relationships.AssignmentGroup.Data.ID
		assignmentGroupID = &value
	}

	var deviceID *int
	if response.Data.Relationships.Device != nil && response.Data.Relationships.Device.Data != nil {
		value := response.Data.Relationships.Device.Data.ID
		deviceID = &value
	}

	return enrollmentFlat{
		ID:                response.Data.ID,
		URL:               response.Data.Attributes.URL,
		UserEnrollment:    response.Data.Attributes.UserEnrollment,
		WelcomeScreen:     response.Data.Attributes.WelcomeScreen,
		Authentication:    response.Data.Attributes.Authentication,
		DeviceGroupID:     deviceGroupID,
		AssignmentGroupID: assignmentGroupID,
		DeviceID:          deviceID,
	}
}

type enrollmentsListResponse struct {
	Data    []enrollmentData `json:"data"`
	HasMore bool             `json:"has_more"`
}

func listEnrollments(ctx context.Context, client *simplemdm.Client, startingAfter int) ([]enrollmentResponse, error) {
	var allEnrollments []enrollmentResponse
	limit := 100

	for {
		url := fmt.Sprintf("https://%s/api/v1/enrollments?limit=%d", client.HostName, limit)
		if startingAfter > 0 {
			url += fmt.Sprintf("&starting_after=%d", startingAfter)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create list enrollments request: %w", err)
		}

		body, err := client.RequestResponse200(req)
		if err != nil {
			return nil, fmt.Errorf("failed to list enrollments: %w", err)
		}

		var response enrollmentsListResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal enrollments list response: %w", err)
		}

		for _, data := range response.Data {
			allEnrollments = append(allEnrollments, enrollmentResponse{Data: data})
		}

		if !response.HasMore {
			break
		}

		if len(response.Data) > 0 {
			startingAfter = response.Data[len(response.Data)-1].ID
		} else {
			break
		}
	}

	return allEnrollments, nil
}
