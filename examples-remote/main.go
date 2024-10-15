// example-remote/main.go

package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"os"
)

// Define the workflow function signature, which matches the one from CRUDWorkflow's CreateWorkflow
// You don't need the full CRUDWorkflow interface here, just the signature for invocation
// TODO simply import
func CreateWorkflow(entityID string, data string) (string, error) {
	return "", nil // This line is just to satisfy the compiler. The actual implementation will be picked up by Temporal.
}

func main() {
	// Set up the Temporal client to interact with the Temporal service
	c, err := client.Dial(client.Options{
		HostPort:  os.Getenv("TEMPORAL_ADDRESS"),
		Namespace: os.Getenv("TEMPORAL_NS"),
	})
	defer c.Close()

	// Define the task queue that the workflow worker is listening on
	taskQueue := "my-task-queue"

	// Configure workflow execution options
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: taskQueue,
	}

	// Define input parameters for the CreateWorkflow
	entityID := "entity-123"
	entityData := "Sample data for entity 123"

	// Start the CreateWorkflow execution remotely
	workflowExecution, err := c.ExecuteWorkflow(context.Background(), workflowOptions, CreateWorkflow, entityID, entityData)
	if err != nil {
		log.Fatal().Msgf("Failed to start CreateWorkflow: %v", err)
	}
	log.Printf("Started CreateWorkflow with WorkflowID: %s and RunID: %s", workflowExecution.GetID(), workflowExecution.GetRunID())

	// Optional: Wait for the workflow to complete and retrieve the result
	var result string
	err = workflowExecution.Get(context.Background(), &result)
	if err != nil {
		log.Fatal().Msgf("Failed to get CreateWorkflow result: %v", err)
	}
	fmt.Printf("CreateWorkflow completed successfully with result: %s\n", result)
}
