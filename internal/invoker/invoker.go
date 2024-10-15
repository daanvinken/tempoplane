package invoker

import (
	"github.com/daanvinken/tempoplane/pkg/entityworkflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

type Invoker struct {
	client    client.Client
	taskQueue string
}

// NewInvoker initializes a new Invoker with the Temporal client and task queue
func NewInvoker(c client.Client, taskQueue string) *Invoker {
	return &Invoker{
		client:    c,
		taskQueue: taskQueue,
	}
}

// RegisterAndRun takes a user implementation of CRUDWorkflow, registers, and runs it
func (inv *Invoker) RegisterAndRun(workflows entityworkflow.CRUDWorkflow) {
	w := worker.New(inv.client, inv.taskQueue, worker.Options{})

	// Register the workflows from the CRUDWorkflow implementation
	w.RegisterWorkflow(workflows.CreateWorkflow)
	w.RegisterWorkflow(workflows.ReadWorkflow)
	w.RegisterWorkflow(workflows.UpdateWorkflow)
	w.RegisterWorkflow(workflows.DeleteWorkflow)

	// Start the worker to listen for workflow tasks
	go func() {
		if err := w.Start(); err != nil {
			log.Fatalf("Unable to start Temporal worker: %v", err)
		}
	}()

	log.Printf("Worker started and workflows registered on task queue: %s", inv.taskQueue)
}
