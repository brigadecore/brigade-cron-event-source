package main

import (
	"log"

	"github.com/brigadecore/brigade-foundations/signals"
	"github.com/brigadecore/brigade/sdk/v2/core"
	"github.com/brigadecore/brigade/sdk/v2/restmachinery"
)

func main() {

	ctx := signals.Context()

	client := core.NewEventsClient(
		"https://localhost:8444",
		"e01f2b82a1d042889396889ad741e9f2EH1yCgpIk13P3piCcJnSGKE1JhyqBCc30gKXNFsyVIPSgDAYnxS1FW2JW9FgjGTKLq6oqKioNsLEvuJ6gq57kXkdgr2MnBhWFTmFJshPb3MosKPvge8ppw8snk5rcxrjuLf9CckMBaP9pSMhm6uwSQ3MJxVYTIMhGCJz1X95fNYB6oewE2tHZkdVluDiZvrnwFiyUE8MMJ5KKdGVi2gfUMEGiN06x3OJ",
		&restmachinery.APIClientOptions{
			AllowInsecureConnections: true,
		},
	)

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
