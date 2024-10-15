package utils

import (
	"github.com/google/uuid"
	"os"
)

// generateWorkflowID creates a unique workflow ID using entityID, operationType, requesterID, and a UUID.
func GenerateWorkflowID() string {
	uniqueID := uuid.New().String()
	//return fmt.Sprintf("%s-%s-%s-%s", operationType, entityID, requesterID, uniqueID)
	return "TempoPlane-" + uniqueID + "-" + os.Getenv("USER")
}
