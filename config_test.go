package main

// nolint: lll
import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/brigadecore/brigade/sdk/v3"
	"github.com/brigadecore/brigade/sdk/v3/restmachinery"
	"github.com/stretchr/testify/require"
)

// Note that unit testing in Go does NOT clear environment variables between
// tests, which can sometimes be a pain, but it's fine here-- so each of these
// test functions uses a series of test cases that cumulatively build upon one
// another.

func TestAPIClientConfig(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func()
		assertions func(
			address string,
			token string,
			opts restmachinery.APIClientOptions,
			err error,
		)
	}{
		{
			name:  "API_ADDRESS not set",
			setup: func() {},
			assertions: func(
				_ string,
				_ string,
				_ restmachinery.APIClientOptions,
				err error,
			) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value not found for")
				require.Contains(t, err.Error(), "API_ADDRESS")
			},
		},
		{
			name: "API_TOKEN not set",
			setup: func() {
				t.Setenv("API_ADDRESS", "foo")
			},
			assertions: func(
				_ string,
				_ string,
				_ restmachinery.APIClientOptions,
				err error,
			) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value not found for")
				require.Contains(t, err.Error(), "API_TOKEN")
			},
		},
		{
			name: "SUCCESS not set",
			setup: func() {
				t.Setenv("API_TOKEN", "bar")
				t.Setenv("API_IGNORE_CERT_WARNINGS", "true")
			},
			assertions: func(
				address string,
				token string,
				opts restmachinery.APIClientOptions,
				err error,
			) {
				require.NoError(t, err)
				require.Equal(t, "foo", address)
				require.Equal(t, "bar", token)
				require.True(t, opts.AllowInsecureConnections)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.setup()
			address, token, opts, err := apiClientConfig()
			testCase.assertions(address, token, opts, err)
		})
	}
}

func TestEvent(t *testing.T) {
	const testProject = "italian"
	const testSource = "brigade.sh/cron"
	const testType = "foo"
	testCases := []struct {
		name       string
		setup      func()
		assertions func(sdk.Event, error)
	}{
		{
			name: "EVENT_PATH not set",
			assertions: func(_ sdk.Event, err error) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "value not found for")
				require.Contains(t, err.Error(), "EVENT_PATH")
			},
		},
		{
			name: "EVENT_PATH path does not exist",
			setup: func() {
				t.Setenv("EVENT_PATH", "/completely/bogus/path")
			},
			assertions: func(_ sdk.Event, err error) {
				require.Error(t, err)
				require.Contains(
					t,
					err.Error(),
					"file /completely/bogus/path does not exist",
				)
			},
		},
		{
			name: "EVENT_PATH does not contain valid json",
			setup: func() {
				eventFile, err := ioutil.TempFile("", "event.json")
				require.NoError(t, err)
				defer eventFile.Close()
				_, err = eventFile.Write([]byte("this is not json"))
				require.NoError(t, err)
				t.Setenv("EVENT_PATH", eventFile.Name())
			},
			assertions: func(_ sdk.Event, err error) {
				require.Error(t, err)
				require.Contains(
					t, err.Error(), "invalid character",
				)
			},
		},
		{
			name: "success",
			setup: func() {
				eventFile, err := ioutil.TempFile("", "event.json")
				require.NoError(t, err)
				defer eventFile.Close()
				_, err =
					eventFile.Write(
						[]byte(
							fmt.Sprintf(
								`{"projectID":%q,"source":%q,"type":%q}`,
								testProject,
								testSource,
								testType,
							),
						),
					)
				require.NoError(t, err)
				t.Setenv("EVENT_PATH", eventFile.Name())
			},
			assertions: func(event sdk.Event, err error) {
				require.NoError(t, err)
				require.Equal(t, testProject, event.ProjectID)
				require.Equal(t, testSource, event.Source)
				require.Equal(t, testType, event.Type)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setup != nil {
				testCase.setup()
			}
			testCase.assertions(event())
		})
	}
}
