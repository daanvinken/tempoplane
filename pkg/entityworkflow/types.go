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

type EntityOutput struct {
	Status  string
	Message string
}
