package acceptancetests

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/riccap/tofu-uptrace-provider/internal/client"
)

// PreCheck validates that required environment variables are set for acceptance tests.
func PreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_ACC") != "1" {
		t.Skip("TF_ACC not set, skipping acceptance test")
	}

	required := []string{
		"UPTRACE_ENDPOINT",
		"UPTRACE_TOKEN",
		"UPTRACE_PROJECT_ID",
	}

	for _, env := range required {
		if os.Getenv(env) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", env)
		}
	}
}

// GetTestProviderConfig returns HCL provider configuration for tests.
func GetTestProviderConfig() string {
	return `
provider "uptrace" {
  endpoint   = "http://localhost:14318/internal/v1"
  token      = "user1_secret_token"
  project_id = 1
}
`
}

// RandomString generates a random string for unique resource names.
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		//nolint:gosec // G404: Use of weak random generator is acceptable for test resource names
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// RandomTestName generates a random name with a prefix for test resources.
func RandomTestName(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, RandomString(8))
}

// WaitForMonitorState polls until monitor reaches expected state.
func WaitForMonitorState(ctx context.Context, client *client.Client, monitorID, expectedState string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for monitor %s to reach state %s", monitorID, expectedState)
			}

			monitor, err := client.GetMonitor(ctx, monitorID)
			if err != nil {
				// Monitor might not exist yet
				continue
			}

			if string(monitor.State) == expectedState {
				return nil
			}
		}
	}
}

// GetTestClient creates a client for testing using environment variables.
func GetTestClient() *client.Client {
	endpoint := os.Getenv("UPTRACE_ENDPOINT")
	token := os.Getenv("UPTRACE_TOKEN")
	projectID := 1 // Default for tests

	client, err := client.New(client.Config{
		Endpoint:  endpoint,
		Token:     token,
		ProjectID: int64(projectID),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create test client: %v", err))
	}

	return client
}
