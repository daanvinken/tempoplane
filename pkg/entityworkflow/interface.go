package entityworkflow

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"os"
)

// CRUDWorkflow defines the interface for CRUD operations as workflows, with standardized input and output
type CRUDWorkflow interface {
	CreateWorkflow(ctx workflow.Context, input EntityInput) (EntityOutput, error)
	ReadWorkflow(ctx workflow.Context, input EntityInput) (EntityOutput, error)
	UpdateWorkflow(ctx workflow.Context, input EntityInput) (EntityOutput, error)
	DeleteWorkflow(ctx workflow.Context, input EntityInput) (EntityOutput, error)
}

// CreateWorkflow is a placeholder for the actual implementation, which will be invoked by Temporal
func CreateWorkflow(ctx workflow.Context, input EntityInput) (EntityOutput, error) {
	// Example output to satisfy the compiler. The actual implementation will provide real logic.
	fmt.Println("This method serves purely as workflow signature, and should not be executed. Program will exit.")
	os.Exit(1)
	return EntityOutput{}, nil
}

// ReadWorkflow is a placeholder for the actual implementation, which will be invoked by Temporal
func ReadWorkflow(ctx workflow.Context, input EntityInput) (EntityOutput, error) {
	fmt.Println("This method serves purely as workflow signature, and should not be executed. Program will exit.")
	os.Exit(1)
	return EntityOutput{}, nil
}

// UpdateWorkflow is a placeholder for the actual implementation, which will be invoked by Temporal
func UpdateWorkflow(ctx workflow.Context, input EntityInput) (EntityOutput, error) {
	fmt.Println("This method serves purely as workflow signature, and should not be executed. Program will exit.")
	os.Exit(1)
	return EntityOutput{}, nil
}

// DeleteWorkflow is a placeholder for the actual implementation, which will be invoked by Temporal
func DeleteWorkflow(ctx workflow.Context, input EntityInput) (EntityOutput, error) {
	fmt.Println("This method serves purely as workflow signature, and should not be executed. Program will exit.")
	os.Exit(1)
	return EntityOutput{}, nil
}
