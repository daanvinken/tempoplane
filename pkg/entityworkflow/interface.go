package entityworkflow

import (
	"go.temporal.io/sdk/workflow"
)

// CRUDWorkflow defines the interface for CRUD operations as workflows
type CRUDWorkflow interface {
	CreateWorkflow(ctx workflow.Context, entityID string, data string) (string, error)
	ReadWorkflow(ctx workflow.Context, entityID string) (string, error)
	UpdateWorkflow(ctx workflow.Context, entityID string, data string) (string, error)
	DeleteWorkflow(ctx workflow.Context, entityID string) (string, error)
}
