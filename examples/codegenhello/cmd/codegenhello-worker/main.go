// Command codegenhello-worker runs a Temporal worker for the codegenhello
// example. The generated Register call is the entire registration surface.
package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/vihren-dev/vihren/examples/codegenhello"
)

func main() {
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	defer temporalClient.Close()

	w := worker.New(temporalClient, codegenhello.DefaultTaskQueue, worker.Options{})
	codegenhello.Register(w, &codegenhello.GreetingActivities{Prefix: "Hello"})

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatal(err)
	}
}
