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
	taskQueue := "TempoPlane-HostLocalFile"

	// Configure workflow execution options
	workflowOptions := client.StartWorkflowOptions{
		TaskQueue: taskQueue,
		ID:        utils.GenerateWorkflowID(),
	}

	entityInput := entityworkflow.EntityInput{
		EntityID:      "entity-123",
		Kind:          "HostLocalFile",
		APIVersion:    "0.0.1",
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
