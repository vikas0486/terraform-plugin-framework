package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// ProviderConfig stores configuration supplied by the
// Terraform provider block.
//
// Example:
//
//	provider "thales" {
//	  endpoint = "http://localhost:8080"
//	}
//
// In future versions, this struct can be extended to include
// authentication, TLS settings, proxy configuration, timeouts,
// retry policies, etc.
type ProviderConfig struct {
	Endpoint types.String `tfsdk:"endpoint"`
}
