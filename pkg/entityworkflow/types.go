package entityworkflow

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// EntityInput represents the input data required for CRUD operations on an entity.
type EntityInput struct {
	EntityID      string               `json:"entityID,omitempty"`
	Kind          string               `json:"kind,omitempty"`
	RequesterID   string               `json:"requesterID,omitempty"`
	Timestamp     int64                `json:"timestamp,omitempty"`
	CorrelationID string               `json:"correlationID,omitempty"`
	Metadata      runtime.RawExtension `json:"metadata,omitempty"` // Flexible field for additional data
}

// Status represents the status of a workflow operation.
type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
	StatusUnknown Status = "unknown"
)

type EntityOutput struct {
	Status  Status `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}
