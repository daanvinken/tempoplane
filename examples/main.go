// examples/main.go

package main

import (
	"fmt"
	"github.com/daanvinken/tempoplane/internal/invoker"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"os"
	"time"
)

// User's implementation of the CRUDWorkflow interface
type MyEntityWorkflow struct{}

func (w *MyEntityWorkflow) CreateWorkflow(ctx workflow.Context, entityID string, data string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute CreateFileActivity
	var fileResult string
	err := workflow.ExecuteActivity(ctx, CreateFileActivity, fmt.Sprintf("/tmp/%s.txt", entityID), data).Get(ctx, &fileResult)
	if err != nil {
		log.Error().Err(err).Str("entityID", entityID).Msg("Failed to create file")
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	log.Info().Str("entityID", entityID).Msg("File created successfully")

	// Execute SendSlackNotificationActivity
	var slackResult string
	slackWebhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	err = workflow.ExecuteActivity(ctx, SendSlackNotificationActivity, slackWebhookURL, fmt.Sprintf("Entity %s created with data %s", entityID, data)).Get(ctx, &slackResult)
	if err != nil {
		log.Error().Err(err).Str("entityID", entityID).Msg("Failed to send Slack notification")
		return "", fmt.Errorf("failed to send Slack notification: %w", err)
	}
	log.Info().Str("entityID", entityID).Msg("Slack notification sent successfully")

	// Return combined result
	finalResult := fmt.Sprintf("%s; %s", fileResult, slackResult)
	return finalResult, nil
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

// RegisterActivities registers the CreateFileActivity and SendSlackNotificationActivity with the worker.
func RegisterActivities(w worker.Worker) {
	w.RegisterActivity(CreateFileActivity)
	w.RegisterActivity(SendSlackNotificationActivity)
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
	myInvoker.RegisterAndRun(myWorkflows, RegisterActivities)

	// Block indefinitely to keep the worker running
	select {}
}
