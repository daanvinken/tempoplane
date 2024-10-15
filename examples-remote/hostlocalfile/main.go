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

	// Execute the CreateWorkflow
	workflowExecution, err := c.ExecuteWorkflow(context.Background(), workflowOptions, entityworkflow.CreateWorkflow, entityInput)
	if err != nil {
		log.Fatal().Msgf("Failed to start CreateWorkflow: %v", err)
	}
	log.Info().Msgf("Started CreateWorkflow with WorkflowID: %s and RunID: %s", workflowExecution.GetID(), workflowExecution.GetRunID())

	// TODO Optional: Wait for CreateWorkflow to complete and retrieve the result
	var createResult entityworkflow.EntityOutput
	err = workflowExecution.Get(context.Background(), &createResult)
	if err != nil {
		log.Fatal().Msgf("Failed to get CreateWorkflow result: %v", err)
	}
	fmt.Printf("CreateWorkflow completed successfully with result: %s\n", createResult.Message)

	// Wait briefly before starting the DeleteWorkflow
	time.Sleep(2 * time.Second)

	// Prepare and start the DeleteWorkflow execution
	deleteWorkflowExecution, err := c.ExecuteWorkflow(context.Background(), workflowOptions, entityworkflow.DeleteWorkflow, entityInput)
	if err != nil {
		log.Fatal().Msgf("Failed to start DeleteWorkflow: %v", err)
	}
	log.Info().Msgf("Started DeleteWorkflow with WorkflowID: %s and RunID: %s", deleteWorkflowExecution.GetID(), deleteWorkflowExecution.GetRunID())

	// TODO Optional: Wait for the DeleteWorkflow to complete and retrieve the result
	var deleteResult entityworkflow.EntityOutput
	err = deleteWorkflowExecution.Get(context.Background(), &deleteResult)
	if err != nil {
		log.Fatal().Msgf("Failed to get DeleteWorkflow result: %v", err)
	}
	fmt.Printf("DeleteWorkflow completed successfully with result: %s\n", deleteResult.Message)

}
