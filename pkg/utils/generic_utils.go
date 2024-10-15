package utils

import (
	"github.com/google/uuid"
)

// generateWorkflowID creates a unique workflow ID using entityID, operationType, requesterID, and a UUID.
func GenerateWorkflowID(requester string) string {
	uniqueID := uuid.New().String()
	//return fmt.Sprintf("%s-%s-%s-%s", operationType, entityID, requesterID, uniqueID)
	return "TempoPlane-" + uniqueID + "-" + requester
}
