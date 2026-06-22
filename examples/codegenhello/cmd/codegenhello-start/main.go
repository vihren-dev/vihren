// Command codegenhello-start starts HelloWorkflow through the generated typed
// client and prints the result. Pass a name as the first argument.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.temporal.io/sdk/client"

	"github.com/vihren-dev/vihren/examples/codegenhello"
)

func main() {
	name := "Ada"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	defer temporalClient.Close()

	out, err := codegenhello.NewClient(temporalClient).HelloWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			ID:        "vihren-codegenhello-" + name,
			TaskQueue: codegenhello.DefaultTaskQueue,
		},
		codegenhello.GreetingInput{Name: name},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.Message)
}
