// example-remote/main.go

package main

import (
	"context"
	"fmt"
	"github.com/daanvinken/tempoplane/pkg/entityworkflow"
	"github.com/daanvinken/tempoplane/pkg/utils"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"os"
	"time"
)

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
		ID:        utils.GenerateWorkflowID(),
	}

	// Define input parameters for the CreateWorkflow
	entityID := "entity-123"

	entityInput := entityworkflow.EntityInput{
		EntityID:      entityID,
		Data:          "SomeTestDataCouldBeJSON",
		RequesterID:   os.Getenv("USER"),
		DC:            "EIN1",
		Env:           "beta",
		Timestamp:     time.Now().Unix(),
		CorrelationID: "420",
	}

	// Start the CreateWorkflow execution remotely
	workflowExecution, err := c.ExecuteWorkflow(context.Background(), workflowOptions, entityworkflow.CreateWorkflow, entityInput)
	if err != nil {
		log.Fatal().Msgf("Failed to start CreateWorkflow: %v", err)
	}
	log.Printf("Started CreateWorkflow with WorkflowID: %s and RunID: %s", workflowExecution.GetID(), workflowExecution.GetRunID())

	// Optional: Wait for the workflow to complete and retrieve the result
	var result entityworkflow.EntityOutput
	err = workflowExecution.Get(context.Background(), &result)
	if err != nil {
		log.Fatal().Msgf("Failed to get CreateWorkflow result: %v", err)
	}
	fmt.Printf("CreateWorkflow completed successfully with result: %s\n", result)
}
