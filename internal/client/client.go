package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/riccap/tofu-uptrace-provider/internal/client/generated"
)

// Client wraps the generated Uptrace API client with higher-level operations.
type Client struct {
	client    *generated.ClientWithResponses
	projectID int64
}

// Config holds the configuration for creating a new Uptrace client.
type Config struct {
	// Endpoint is the base URL of the Uptrace API
	Endpoint string
	// Token is the authentication bearer token
	Token string
	// ProjectID is the default project ID for operations
	ProjectID int64
	// HTTPClient is an optional custom HTTP client
	HTTPClient *http.Client
}

// New creates a new Uptrace API client.
func New(cfg Config) (*Client, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("token is required")
	}
	if cfg.ProjectID <= 0 {
		return nil, fmt.Errorf("projectID must be greater than 0")
	}

	// Ensure endpoint doesn't have trailing slash
	endpoint := strings.TrimSuffix(cfg.Endpoint, "/")

	// Create request editor to add authentication header
	authEditor := func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
		return nil
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	client, err := generated.NewClientWithResponses(
		endpoint,
		generated.WithHTTPClient(httpClient),
		generated.WithRequestEditorFn(authEditor),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	return &Client{
		client:    client,
		projectID: cfg.ProjectID,
	}, nil
}

// ListMonitors retrieves all monitors for the project.
func (c *Client) ListMonitors(ctx context.Context) ([]generated.Monitor, error) {
	resp, err := c.client.ListMonitorsWithResponse(ctx, c.projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list monitors: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, c.handleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil {
		return []generated.Monitor{}, nil
	}

	return resp.JSON200.Monitors, nil
}

// GetMonitor retrieves a specific monitor by ID.
func (c *Client) GetMonitor(ctx context.Context, monitorID string) (*generated.Monitor, error) {
	resp, err := c.client.GetMonitorWithResponse(ctx, c.projectID, monitorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get monitor: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, c.handleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected empty response")
	}

	return &resp.JSON200.Monitor, nil
}

// CreateMonitor creates a new monitor.
func (c *Client) CreateMonitor(ctx context.Context, input generated.MonitorInput) (*generated.Monitor, error) {
	resp, err := c.client.CreateMonitorWithResponse(ctx, c.projectID, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitor: %w", err)
	}

	// Uptrace API returns 200 for successful creation
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, c.handleErrorResponse(resp.StatusCode(), resp.Body)
	}

	// Try 200 first (most common), then 201
	if resp.JSON200 != nil {
		return &resp.JSON200.Monitor, nil
	}
	if resp.JSON201 != nil {
		return &resp.JSON201.Monitor, nil
	}

	return nil, fmt.Errorf("unexpected empty response")
}

// UpdateMonitor updates an existing monitor.
func (c *Client) UpdateMonitor(ctx context.Context, monitorID string, input generated.MonitorInput) (*generated.Monitor, error) {
	resp, err := c.client.UpdateMonitorWithResponse(ctx, c.projectID, monitorID, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update monitor: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, c.handleErrorResponse(resp.StatusCode(), resp.Body)
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("unexpected empty response")
	}

	return &resp.JSON200.Monitor, nil
}

// DeleteMonitor deletes a monitor by ID.
func (c *Client) DeleteMonitor(ctx context.Context, monitorID string) error {
	resp, err := c.client.DeleteMonitorWithResponse(ctx, c.projectID, monitorID)
	if err != nil {
		return fmt.Errorf("failed to delete monitor: %w", err)
	}

	// Uptrace API returns 200 for successful deletion
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusNoContent {
		return c.handleErrorResponse(resp.StatusCode(), resp.Body)
	}

	return nil
}

// handleErrorResponse processes error responses from the API.
func (c *Client) handleErrorResponse(statusCode int, body []byte) error {
	switch statusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("bad request: %s", string(body))
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized: invalid or missing authentication token")
	case http.StatusForbidden:
		return fmt.Errorf("forbidden: insufficient permissions")
	case http.StatusNotFound:
		return fmt.Errorf("not found: resource does not exist")
	case http.StatusInternalServerError:
		return fmt.Errorf("internal server error: %s", string(body))
	default:
		return fmt.Errorf("unexpected status code %d: %s", statusCode, string(body))
	}
}
