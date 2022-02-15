package main

import (
	"github.com/brigadecore/brigade/sdk/v3"
	"github.com/pkg/errors"
)

const source = "brigade.sh/cron"

func validateEvent(event sdk.Event) error {
	if event.ProjectID == "" {
		return errors.New("project ID not specified; refusing event")
	}
	if event.Source != source {
		return errors.Errorf(
			"source is not %q; refusing event",
			source,
		)
	}
	return nil
}
