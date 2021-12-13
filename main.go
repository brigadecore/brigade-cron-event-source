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

	brigadeSource, brigadeType, _ := eventConfig()
	if err != nil {
		log.Fatal(err)
	}

	qualifiers, _ := qualifiersConfig()

	labels, _ := labelsConfig()

	payload, _ := payloadConfig()

	brigadeEvent := core.Event{
		Source:     brigadeSource,
		Type:       brigadeType,
		Qualifiers: qualifiers,
		Labels:     labels,
		Payload:    payload,
	}

	eventList, err := client.Create(ctx, brigadeEvent)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(eventList)
	}

}
