package main

import (
	"encoding/json"
	"fmt"
	"github.com/daanvinken/tempoplane/internal/invoker"
	ew "github.com/daanvinken/tempoplane/pkg/entityworkflow"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"os"
	"time"
)

// User's implementation of the CRUDWorkflow interface
type MyEntityWorkflow struct{}

func (w *MyEntityWorkflow) CreateWorkflow(ctx workflow.Context, entityInput ew.EntityInput) (ew.EntityOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Unpack Metadata to get specific fields like "filePath" for CreateFileActivity
	var metadata map[string]interface{}
	if err := json.Unmarshal(entityInput.Metadata.Raw, &metadata); err != nil {
		log.Error().Err(err).Str("entityID", entityInput.EntityID).Msg("Failed to parse metadata")
		return ew.EntityOutput{Status: ew.StatusError, Message: "Failed to parse metadata"}, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Extract the "filePath" from metadata
	filePath, ok := metadata["filePath"].(string)
	if !ok || filePath == "" {
		log.Error().Str("entityID", entityInput.EntityID).Msg("File path missing in metadata")
		return ew.EntityOutput{Status: ew.StatusError, Message: "File path missing in metadata"}, fmt.Errorf("file path missing in metadata")
	}

	// Execute CreateFileActivity with extracted filePath
	var fileResult string
	err := workflow.ExecuteActivity(ctx, CreateFileActivity, filePath).Get(ctx, &fileResult)
	if err != nil {
		log.Error().Err(err).Str("entityID", entityInput.EntityID).Msg("Failed to create file")
		return ew.EntityOutput{Status: ew.StatusError, Message: "Failed to create file"}, fmt.Errorf("failed to create file: %w", err)
	}
	log.Info().Str("entityID", entityInput.EntityID).Msg("File created successfully")

	// Execute SendSlackNotificationActivity using the entity ID and metadata content
	var slackResult string
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	err = workflow.ExecuteActivity(ctx, SendSlackNotificationActivity, slackWebhookURL, fmt.Sprintf("Entity %s created with file path %s", entityInput.EntityID, filePath)).Get(ctx, &slackResult)
	if err != nil {
		log.Error().Err(err).Str("entityID", entityInput.EntityID).Msg("Failed to send Slack notification")
		return ew.EntityOutput{Status: ew.StatusError, Message: "Failed to send Slack notification"}, fmt.Errorf("failed to send Slack notification: %w", err)
	}
	log.Info().Str("entityID", entityInput.EntityID).Msg("Slack notification sent successfully")

	// Return combined result
	finalResult := fmt.Sprintf("%s; %s", fileResult, slackResult)
	return ew.EntityOutput{Status: ew.StatusSuccess, Message: finalResult}, nil
}

func (w *MyEntityWorkflow) ReadWorkflow(ctx workflow.Context, entityInput ew.EntityInput) (ew.EntityOutput, error) {
	workflow.GetLogger(ctx).Info("Running ReadWorkflow", "entityID", entityInput.EntityID)
	message := fmt.Sprintf("Read entity %s", entityInput.EntityID)
	return ew.EntityOutput{Status: ew.StatusSuccess, Message: message}, nil
}

func (w *MyEntityWorkflow) UpdateWorkflow(ctx workflow.Context, entityInput ew.EntityInput) (ew.EntityOutput, error) {
	workflow.GetLogger(ctx).Info("Running UpdateWorkflow", "entityID", entityInput.EntityID)
	return ew.EntityOutput{Status: ew.StatusSuccess, Message: "did nothing, is update"}, nil
}

func (w *MyEntityWorkflow) DeleteWorkflow(ctx workflow.Context, entityInput ew.EntityInput) (ew.EntityOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Unpack Metadata to get specific fields like "filePath" for DeleteFileActivity
	var metadata map[string]interface{}
	if err := json.Unmarshal(entityInput.Metadata.Raw, &metadata); err != nil {
		log.Error().Err(err).Str("entityID", entityInput.EntityID).Msg("Failed to parse metadata")
		return ew.EntityOutput{Status: ew.StatusError, Message: "Failed to parse metadata"}, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Extract the "filePath" from metadata
	filePath, ok := metadata["filePath"].(string)
	if !ok || filePath == "" {
		log.Error().Str("entityID", entityInput.EntityID).Msg("File path missing in metadata")
		return ew.EntityOutput{Status: ew.StatusError, Message: "File path missing in metadata"}, fmt.Errorf("file path missing in metadata")
	}

	// Execute DeleteFileActivity with extracted filePath
	var result string
	err := workflow.ExecuteActivity(ctx, DeleteFileActivity, filePath).Get(ctx, &result)
	if err != nil {
		log.Error().Err(err).Str("entityID", entityInput.EntityID).Msg("Failed to delete file")
		return ew.EntityOutput{Status: ew.StatusError, Message: "Failed to delete file"}, fmt.Errorf("failed to delete file: %w", err)
	}
	log.Info().Str("entityID", entityInput.EntityID).Msg("File deleted successfully")

	return ew.EntityOutput{Status: ew.StatusSuccess, Message: result}, nil
}

// RegisterActivities registers the CreateFileActivity and SendSlackNotificationActivity with the worker.
func RegisterActivities(w worker.Worker) {
	w.RegisterActivity(CreateFileActivity)
	w.RegisterActivity(SendSlackNotificationActivity)
	w.RegisterActivity(DeleteFileActivity)
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
	taskQueue := "TempoPlane-HostLocalFile"
	myInvoker := invoker.NewInvoker(c, taskQueue)

	// Create an instance of the user's CRUD workflow implementation
	myWorkflows := &MyEntityWorkflow{}

	// Register and run the workflows through the Invoker
	myInvoker.RegisterAndRun(myWorkflows, RegisterActivities)

	// Block indefinitely to keep the worker running
	select {}
}
