package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/riccap/terraform-provider-uptrace/internal/client/generated"
)

func TestDashboardToState_Complete(t *testing.T) {
	ctx := context.Background()

	// Create API dashboard with all fields
	id := int64(123)
	name := "Test Dashboard"
	pinned := true
	createdAt := 1767348994143.7922
	updatedAt := 1767349994143.7922
	yamlContent := `name: Test Dashboard
gridRows:
  - items:
      - type: chart
        title: CPU Usage
        params:
          metrics:
            - name: system.cpu.utilization
              alias: cpu`

	dashboard := &generated.Dashboard{
		Id:        id,
		Name:      name,
		Pinned:    &pinned,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	}

	var state DashboardResourceModel
	diags := diag.Diagnostics{}

	// Convert to state
	dashboardToState(ctx, dashboard, yamlContent, &state, &diags)

	// Verify no errors
	require.False(t, diags.HasError(), "Conversion should not produce errors")

	// Verify all fields
	assert.Equal(t, "123", state.ID.ValueString())
	assert.Equal(t, "Test Dashboard", state.Name.ValueString())
	assert.Equal(t, yamlContent, state.YAML.ValueString())
	assert.True(t, state.Pinned.ValueBool())
	assert.Equal(t, fmt.Sprintf("%f", createdAt), state.CreatedAt.ValueString())
	assert.Equal(t, fmt.Sprintf("%f", updatedAt), state.UpdatedAt.ValueString())
}

func TestDashboardToState_MinimalFields(t *testing.T) {
	ctx := context.Background()

	// Create API dashboard with minimal fields
	id := int64(456)
	name := "Minimal Dashboard"
	yamlContent := `name: Minimal Dashboard`

	dashboard := &generated.Dashboard{
		Id:   id,
		Name: name,
		// No optional fields set
	}

	var state DashboardResourceModel
	diags := diag.Diagnostics{}

	// Convert to state
	dashboardToState(ctx, dashboard, yamlContent, &state, &diags)

	// Verify no errors
	require.False(t, diags.HasError(), "Conversion should not produce errors")

	// Verify required fields
	assert.Equal(t, "456", state.ID.ValueString())
	assert.Equal(t, "Minimal Dashboard", state.Name.ValueString())
	assert.Equal(t, yamlContent, state.YAML.ValueString())

	// Verify optional fields have defaults
	assert.False(t, state.Pinned.ValueBool(), "Pinned should default to false")
	assert.True(t, state.CreatedAt.IsNull(), "CreatedAt should be null when not provided")
	assert.True(t, state.UpdatedAt.IsNull(), "UpdatedAt should be null when not provided")
}

func TestDashboardToState_LargeYAML(t *testing.T) {
	ctx := context.Background()

	// Create API dashboard with large YAML content
	id := int64(789)
	name := "Complex Dashboard"
	yamlContent := `name: Complex Dashboard
gridRows:
  - items:
      - type: chart
        title: Request Rate
        params:
          metrics:
            - name: http_requests_total
              alias: requests
          query: rate(requests[5m])
          chartKind: line
          legend:
            show: true
      - type: table
        title: Service List
        params:
          metrics:
            - name: service_info
              alias: services
          columns:
            - name: service_name
              label: Service
  - items:
      - type: gauge
        title: Error Rate
        params:
          metrics:
            - name: http_errors_total
              alias: errors
          valueMapping:
            - value: 0-5
              color: green
            - value: 5-10
              color: yellow
            - value: 10+
              color: red
      - type: heatmap
        title: Response Time Distribution
        params:
          metrics:
            - name: http_response_time
              alias: response_time`

	dashboard := &generated.Dashboard{
		Id:   id,
		Name: name,
	}

	var state DashboardResourceModel
	diags := diag.Diagnostics{}

	// Convert to state
	dashboardToState(ctx, dashboard, yamlContent, &state, &diags)

	// Verify no errors
	require.False(t, diags.HasError(), "Conversion should not produce errors")

	// Verify YAML is preserved exactly
	assert.Equal(t, yamlContent, state.YAML.ValueString())
	assert.Greater(t, len(state.YAML.ValueString()), 100, "YAML content should be large")
}

func TestIsNotFoundError_True(t *testing.T) {
	tests := []struct {
		name  string
		error error
	}{
		{
			name:  "exact match",
			error: fmt.Errorf("not found"),
		},
		{
			name:  "with prefix",
			error: fmt.Errorf("resource not found"),
		},
		{
			name:  "with suffix",
			error: fmt.Errorf("not found: dashboard does not exist"),
		},
		{
			name:  "in middle",
			error: fmt.Errorf("error: not found in database"),
		},
		{
			name:  "uppercase",
			error: fmt.Errorf("NOT FOUND"),
		},
		{
			name:  "mixed case",
			error: fmt.Errorf("Resource Not Found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNotFoundError(tt.error)
			assert.True(t, result, "Error should be recognized as not found: %v", tt.error)
		})
	}
}

func TestIsNotFoundError_False(t *testing.T) {
	tests := []struct {
		name  string
		error error
	}{
		{
			name:  "nil error",
			error: nil,
		},
		{
			name:  "different error",
			error: fmt.Errorf("bad request"),
		},
		{
			name:  "unauthorized",
			error: fmt.Errorf("unauthorized"),
		},
		{
			name:  "internal server error",
			error: fmt.Errorf("internal server error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNotFoundError(tt.error)
			assert.False(t, result, "Error should not be recognized as not found: %v", tt.error)
		})
	}
}
