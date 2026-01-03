package acctest

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/riccap/tofu-uptrace-provider/internal/provider"
)

// TestAccProtoV6ProviderFactories is the provider factory for acceptance tests
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"uptrace": providerserver.NewProtocol6WithError(provider.New("test")()),
}
