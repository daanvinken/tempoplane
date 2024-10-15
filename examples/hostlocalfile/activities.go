// internal/workflows/activities.go

package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// CreateFileActivity creates a file on the filesystem with the specified content.
func CreateFileActivity(ctx context.Context, filePath, content string) (string, error) {
	log.Info().Str("filePath", filePath).Msg("Creating file")
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		log.Error().Err(err).Str("filePath", filePath).Msg("Failed to create file")
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	log.Info().Str("filePath", filePath).Msg("File created successfully")
	return fmt.Sprintf("File created at %s", filePath), nil
}

// SendSlackNotificationActivity sends a notification to Slack via a webhook URL.
func SendSlackNotificationActivity(ctx context.Context, webhookURL, message string) (string, error) {
	log.Info().Msg("Sending Slack notification")

	// Construct the Slack message payload
	payload := fmt.Sprintf(`{"text": "%s", "channel": "%s"}`, message, os.Getenv("SLACK_CHANNEL"))
	req, err := http.NewRequest("POST", webhookURL, strings.NewReader(payload))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create HTTP request for Slack")
		return "", fmt.Errorf("failed to create Slack request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send Slack notification")
		return "", fmt.Errorf("failed to send Slack notification: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("Failed to read response body")
		} else {
			bodyContent := string(bodyBytes)
			log.Error().
				Int("statusCode", resp.StatusCode).
				Str("body", bodyContent).
				Msg("Non-OK response from Slack")
		}

		return "", fmt.Errorf("non-OK response from Slack: %d", resp.StatusCode)
	}
	return "Slack notification sent successfully", nil
}

// DeleteFileActivity deletes a file on the filesystem at the specified path.
func DeleteFileActivity(ctx context.Context, filePath string) (string, error) {
	log.Info().Str("filePath", filePath).Msg("Deleting file")
	err := os.Remove(filePath)
	if err != nil {
		log.Error().Err(err).Str("filePath", filePath).Msg("Failed to delete file")
		return "", fmt.Errorf("failed to delete file: %w", err)
	}
	log.Info().Str("filePath", filePath).Msg("File deleted successfully")
	return fmt.Sprintf("File deleted at %s", filePath), nil
}
