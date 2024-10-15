package entityworkflow

type EntityInput struct {
	EntityID      string
	Data          string
	RequesterID   string
	DC            string
	Env           string
	Timestamp     int64
	CorrelationID string
}

// Status represents the status of a workflow operation.
type Status string

const (
	StatusSuccess Status = "success"
	StatusError   Status = "error"
	StatusUnknown Status = "unknown"
)

type EntityOutput struct {
	Status  Status
	Message string
}
