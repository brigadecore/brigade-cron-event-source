package main

import (
	"io/ioutil"

	"github.com/brigadecore/brigade-foundations/os"
	"github.com/brigadecore/brigade/sdk/v3"
	"github.com/brigadecore/brigade/sdk/v3/restmachinery"
	"github.com/ghodss/yaml"
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

func event() (sdk.Event, error) {
	event := sdk.Event{}
	eventBytes, err := ioutil.ReadFile("/app/config/event.yaml")
	if err != nil {
		return event, err
	}
	err = yaml.Unmarshal(eventBytes, &event)
	return event, err
}
