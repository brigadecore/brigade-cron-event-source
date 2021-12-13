package main

// nolint: lll
import (
	"encoding/json"
	"io/ioutil"

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

func eventConfig() (string, string, error) {
	brigadeSource, err := os.GetRequiredEnvVar("BRIGADE_SOURCE")
	if err != nil {
		return brigadeSource, "", err
	}
	brigadeType, err := os.GetRequiredEnvVar("BRIGADE_TYPE")
	if err != nil {
		return brigadeSource, brigadeType, err
	}
	return brigadeSource, brigadeType, err
}

func qualifiersConfig() (map[string]string, error) {
	qualifiersBytes, err := ioutil.ReadFile("/cronjob-config/qualifiers")
	if err != nil {
		return map[string]string{}, err
	}
	qualifiersPlainText := map[string]string{}
	if err :=
		json.Unmarshal(qualifiersBytes, &qualifiersPlainText); err != nil {
		return map[string]string{}, err
	}
	return qualifiersPlainText, nil
}

func labelsConfig() (map[string]string, error) {
	labelsBytes, err := ioutil.ReadFile("/cronjob-config/labels")
	if err != nil {
		return map[string]string{}, err
	}
	labelsPlainText := map[string]string{}
	if err :=
		json.Unmarshal(labelsBytes, &labelsPlainText); err != nil {
		return map[string]string{}, err
	}
	return labelsPlainText, nil
}

func payloadConfig() (string, error) {
	payloadBytes, err := ioutil.ReadFile("/cronjob-config/payload")
	if err != nil {
		return "", err
	}
	return string(payloadBytes), nil
}
