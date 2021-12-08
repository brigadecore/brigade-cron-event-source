package main

import (
	"log"

	"github.com/brigadecore/brigade-foundations/signals"
	"github.com/brigadecore/brigade/sdk/v2/core"
)

func main() {

	ctx := signals.Context()

	address, token, opts, err := apiClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := core.NewEventsClient(address, token, &opts)

	brigadeEvent := core.Event{
		Source: "cronsource",
		Type:   "cron",
		// Qualifiers: map[string]string{
		// 	"source": "cronsource",
		// 	"type":   "cron",
		// },
		Payload: string("test"),
	}

	test, err := client.Create(ctx, brigadeEvent)

	log.Println(
		test, err,
	)
}
