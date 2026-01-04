package acceptancetests

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TestAccProtoV6ProviderFactories is populated by the provider package to avoid import cycles.
var TestAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
