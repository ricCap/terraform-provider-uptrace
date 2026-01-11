package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/riccap/terraform-provider-uptrace/internal/client"
	"github.com/riccap/terraform-provider-uptrace/internal/client/generated"
)

func newTestClient(server *httptest.Server) *client.Client {
	c, err := client.New(client.Config{
		Endpoint:  server.URL,
		Token:     "test-token",
		ProjectID: 1,
	})
	if err != nil {
		panic(err)
	}
	return c
}

// TestPinDashboard tests the PinDashboard client method
func TestPinDashboard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/metrics/1/dashboards/123/pinned" {
			t.Errorf("Expected path /metrics/1/dashboards/123/pinned, got %s", r.URL.Path)
		}

		// Return success
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.PinDashboard(context.Background(), 123)
	if err != nil {
		t.Fatalf("PinDashboard failed: %v", err)
	}
}

// TestUnpinDashboard tests the UnpinDashboard client method
func TestUnpinDashboard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/metrics/1/dashboards/123/unpinned" {
			t.Errorf("Expected path /metrics/1/dashboards/123/unpinned, got %s", r.URL.Path)
		}

		// Return success
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.UnpinDashboard(context.Background(), 123)
	if err != nil {
		t.Fatalf("UnpinDashboard failed: %v", err)
	}
}

// TestCloneDashboard tests the CloneDashboard client method
func TestCloneDashboard(t *testing.T) {
	expectedDashboard := generated.Dashboard{
		Id:   456,
		Name: "Cloned Dashboard",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/metrics/1/dashboards/123/clone" {
			t.Errorf("Expected path /metrics/1/dashboards/123/clone, got %s", r.URL.Path)
		}

		// Return cloned dashboard
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]generated.Dashboard{
			"dashboard": expectedDashboard,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := newTestClient(server)
	dashboard, err := c.CloneDashboard(context.Background(), 123)
	if err != nil {
		t.Fatalf("CloneDashboard failed: %v", err)
	}

	if dashboard.Id != expectedDashboard.Id {
		t.Errorf("Expected dashboard ID %d, got %d", expectedDashboard.Id, dashboard.Id)
	}
	if dashboard.Name != expectedDashboard.Name {
		t.Errorf("Expected dashboard name %s, got %s", expectedDashboard.Name, dashboard.Name)
	}
}

// TestResetDashboard tests the ResetDashboard client method
func TestResetDashboard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/metrics/1/dashboards/123/reset" {
			t.Errorf("Expected path /metrics/1/dashboards/123/reset, got %s", r.URL.Path)
		}

		// Return success
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.ResetDashboard(context.Background(), 123)
	if err != nil {
		t.Fatalf("ResetDashboard failed: %v", err)
	}
}

// TestPinDashboard_Error tests error handling in PinDashboard
func TestPinDashboard_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		response := map[string]interface{}{
			"error": map[string]string{
				"code":    "not_found",
				"message": "Dashboard not found",
			},
			"statusCode": 404,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := newTestClient(server)
	err := c.PinDashboard(context.Background(), 123)
	if err == nil {
		t.Fatal("Expected error for 404 response, got nil")
	}
}

// TestCloneDashboard_Error tests error handling in CloneDashboard
func TestCloneDashboard_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"error": map[string]string{
				"code":    "invalid_request",
				"message": "Cannot clone dashboard",
			},
			"statusCode": 400,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("Failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := newTestClient(server)
	_, err := c.CloneDashboard(context.Background(), 123)
	if err == nil {
		t.Fatal("Expected error for 400 response, got nil")
	}
}
