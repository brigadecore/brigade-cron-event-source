package main

// nolint: lll
import (
	"github.com/brigadecore/brigade-foundations/os"
	"github.com/brigadecore/brigade/sdk/v2/restmachinery"
)

// apiClientConfig populates the Brigade SDK's APIClientOptions from
// environment variables.
func apiClientConfig() (string, string, restmachinery.APIClientOptions, error) {
	opts := restmachinery.APIClientOptions{}
	address, err := os.GetRequiredEnvVar("API_ADDRESS")
	if err != nil {
		return address, "", opts, err
	}
	token, err := os.GetRequiredEnvVar("API_TOKEN")
	if err != nil {
		return address, token, opts, err
	}
	opts.AllowInsecureConnections, err =
		os.GetBoolFromEnvVar("API_IGNORE_CERT_WARNINGS", false)
	return address, token, opts, err
}
