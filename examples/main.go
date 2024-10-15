// examples/main.go

package main

import (
	"fmt"
	"github.com/daanvinken/tempoplane/internal/invoker"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"os"
)

// User's implementation of the CRUDWorkflow interface
type MyEntityWorkflow struct{}

func (w *MyEntityWorkflow) CreateWorkflow(ctx workflow.Context, entityID string, data string) (string, error) {
	// User's custom activities within the Create workflow
	workflow.GetLogger(ctx).Info("Running CreateWorkflow", "entityID", entityID, "data", data)
	return fmt.Sprintf("Created entity %s with data: %s", entityID, data), nil
}

func (w *MyEntityWorkflow) ReadWorkflow(ctx workflow.Context, entityID string) (string, error) {
	// User's custom activities within the Read workflow
	workflow.GetLogger(ctx).Info("Running ReadWorkflow", "entityID", entityID)
	return fmt.Sprintf("Read entity %s", entityID), nil
}

func (w *MyEntityWorkflow) UpdateWorkflow(ctx workflow.Context, entityID string, data string) (string, error) {
	// User's custom activities within the Update workflow
	workflow.GetLogger(ctx).Info("Running UpdateWorkflow", "entityID", entityID, "data", data)
	return fmt.Sprintf("Updated entity %s with data: %s", entityID, data), nil
}

func (w *MyEntityWorkflow) DeleteWorkflow(ctx workflow.Context, entityID string) (string, error) {
	// User's custom activities within the Delete workflow
	workflow.GetLogger(ctx).Info("Running DeleteWorkflow", "entityID", entityID)
	fmt.Println("test")
	return fmt.Sprintf("Deleted entity %s", entityID), nil
}

func main() {
	// Create a Temporal client
	c, err := client.Dial(client.Options{
		HostPort:  os.Getenv("TEMPORAL_ADDRESS"),
		Namespace: os.Getenv("TEMPORAL_NS"),
	})
	if err != nil {
		log.Fatal().Msgf("Unable to create Temporal client: %v", err)
	}
	defer c.Close()

	// Initialize the Invoker with the client and task queue name
	taskQueue := "my-task-queue"
	myInvoker := invoker.NewInvoker(c, taskQueue)

	// Create an instance of the user's CRUD workflow implementation
	myWorkflows := &MyEntityWorkflow{}

	// Register and run the workflows through the Invoker
	myInvoker.RegisterAndRun(myWorkflows)

	// Block indefinitely to keep the worker running
	select {}
}
