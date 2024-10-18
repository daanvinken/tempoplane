package entityworkflow

type EntityInput struct {
	EntityID      string `json:"entityID,omitempty"`
	Kind          string `json:"kind,omitempty"`
	APIVersion    string `json:"apiVersion,omitempty"`
	Data          string `json:"data,omitempty"`
	RequesterID   string `json:"requesterID,omitempty"`
	Timestamp     int64  `json:"timestamp,omitempty"`
	CorrelationID string `json:"correlationID,omitempty"`
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
