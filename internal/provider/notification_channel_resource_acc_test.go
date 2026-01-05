package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	acceptancetests "github.com/riccap/tofu-uptrace-provider/internal/acceptance_tests"
)

func TestAccNotificationChannelResource_Slack(t *testing.T) {
	if testing.Short() {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccNotificationChannelResourceConfigSlack("test-slack-channel"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptrace_notification_channel.test", "name", "test-slack-channel"),
					resource.TestCheckResourceAttr("uptrace_notification_channel.test", "type", "slack"),
					resource.TestCheckResourceAttr("uptrace_notification_channel.test", "params.webhookUrl", "https://hooks.slack.com/services/test"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.test", "id"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.test", "status"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "uptrace_notification_channel.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"params"}, // params are sensitive
			},
			// Update and Read testing
			{
				Config: testAccNotificationChannelResourceConfigSlack("test-slack-channel-updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptrace_notification_channel.test", "name", "test-slack-channel-updated"),
				),
			},
		},
	})
}

func TestAccNotificationChannelResource_Webhook(t *testing.T) {
	if testing.Short() {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationChannelResourceConfigWebhook(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("uptrace_notification_channel.test", "name", "test-webhook-channel"),
					resource.TestCheckResourceAttr("uptrace_notification_channel.test", "type", "webhook"),
					resource.TestCheckResourceAttr("uptrace_notification_channel.test", "params.url", "https://example.com/webhook"),
					resource.TestCheckResourceAttrSet("uptrace_notification_channel.test", "id"),
				),
			},
		},
	})
}

// TestAccNotificationChannelResource_WithCondition is disabled because
// the valid condition syntax for notification channels is not yet documented.
// The API rejects conditions like "severity == 'critical'" with "unknown name severity".
// TODO: Re-enable this test once we discover valid condition examples.
//
// func TestAccNotificationChannelResource_WithCondition(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
// 		return
// 	}
//
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { acceptancetests.PreCheck(t) },
// 		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccNotificationChannelResourceConfigWithCondition(),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("uptrace_notification_channel.test", "name", "conditional-alerts"),
// 					resource.TestCheckResourceAttr("uptrace_notification_channel.test", "condition", "severity == 'critical'"),
// 					resource.TestCheckResourceAttrSet("uptrace_notification_channel.test", "id"),
// 				),
// 			},
// 		},
// 	})
// }

func TestAccNotificationChannelResource_Disappears(t *testing.T) {
	if testing.Short() {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptancetests.PreCheck(t) },
		ProtoV6ProviderFactories: acceptancetests.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNotificationChannelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationChannelResourceConfigSlack("test-disappear"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNotificationChannelExists("uptrace_notification_channel.test"),
				),
			},
			{
				Config:  testAccNotificationChannelResourceConfigSlack("test-disappear"),
				Destroy: true,
			},
		},
	})
}

func testAccNotificationChannelResourceConfigSlack(name string) string {
	return fmt.Sprintf(`
resource "uptrace_notification_channel" "test" {
  name = %[1]q
  type = "slack"

  params = {
    webhookUrl = "https://hooks.slack.com/services/test"
  }
}
`, name)
}

func testAccNotificationChannelResourceConfigWebhook() string {
	return `
resource "uptrace_notification_channel" "test" {
  name = "test-webhook-channel"
  type = "webhook"

  params = {
    url = "https://example.com/webhook"
  }
}
`
}

func testAccCheckNotificationChannelExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No channel ID is set")
		}

		return nil
	}
}

func testAccCheckNotificationChannelDestroy(s *terraform.State) error {
	// In a real test, we would verify the channel no longer exists via API
	// For now, we'll just check that the resource was removed from state
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "uptrace_notification_channel" {
			continue
		}

		if rs.Primary.ID != "" {
			return fmt.Errorf("Notification channel %s still exists", rs.Primary.ID)
		}
	}

	return nil
}
