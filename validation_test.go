package main

import (
	"testing"

	"github.com/brigadecore/brigade/sdk/v3"
	"github.com/stretchr/testify/require"
)

func TestValidateEvent(t *testing.T) {
	testCases := []struct {
		name       string
		event      sdk.Event
		assertions func(error)
	}{
		{
			name:  "project not specified",
			event: sdk.Event{},
			assertions: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "project ID not specified")
				require.Contains(t, err.Error(), "refusing event")
			},
		},
		{
			name: "invalid source",
			event: sdk.Event{
				ProjectID: "italian",
				Source:    "github.com/github", // Trying to emulater another source
			},
			assertions: func(err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "source is not")
				require.Contains(t, err.Error(), "refusing event")
			},
		},
		{
			name: "valid event",
			event: sdk.Event{
				ProjectID: "italian",
				Source:    source,
			},
			assertions: func(err error) {
				require.NoError(t, err)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.assertions(validateEvent(testCase.event))
		})
	}
}
