package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/brigadecore/brigade-foundations/file"
	"github.com/brigadecore/brigade-foundations/os"
	"github.com/brigadecore/brigade/sdk/v3"
	"github.com/brigadecore/brigade/sdk/v3/restmachinery"
	"github.com/pkg/errors"
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
	eventPath, err := os.GetRequiredEnvVar("EVENT_PATH")
	if err != nil {
		return event, err
	}
	var exists bool
	if exists, err = file.Exists(eventPath); err != nil {
		return event, err
	}
	if !exists {
		return event, errors.Errorf("file %s does not exist", eventPath)
	}
	eventBytes, err := ioutil.ReadFile(eventPath)
	if err != nil {
		return event, err
	}
	err = json.Unmarshal(eventBytes, &event)
	return event, err
}
