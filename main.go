package main

import (
	"context"
	"log"

	"github.com/brigadecore/brigade/sdk/v3"
)

func main() {
	address, token, opts, err := apiClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := sdk.NewEventsClient(address, token, &opts)

	event, err := event()
	if err != nil {
		log.Fatal(err)
	}

	var eventList sdk.EventList
	if eventList, err =
		client.Create(context.Background(), event, nil); err != nil {
		log.Fatalf("error creating event: %s", err)
	}

	log.Println("Created events: ")
	for _, event = range eventList.Items {
		log.Println(event.ID)
	}
}
