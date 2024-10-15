package invoker

import (
	"github.com/daanvinken/tempoplane/pkg/entityworkflow"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"os"
)

type Invoker struct {
	client    client.Client
	taskQueue string
}

// NewInvoker initializes a new Invoker
func NewInvoker(c client.Client, taskQueue string) *Invoker {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return &Invoker{
		client:    c,
		taskQueue: taskQueue,
	}
}

// RegisterAndRun takes CRUDWorkflow and RegisterActivities functions to register workflows and activities
func (inv *Invoker) RegisterAndRun(workflows entityworkflow.CRUDWorkflow, registerActivities func(worker.Worker)) {
	w := worker.New(inv.client, inv.taskQueue, worker.Options{})

	// Register the CRUD workflows
	w.RegisterWorkflow(workflows.CreateWorkflow)
	w.RegisterWorkflow(workflows.ReadWorkflow)
	w.RegisterWorkflow(workflows.UpdateWorkflow)
	w.RegisterWorkflow(workflows.DeleteWorkflow)

	// Call the user-provided function to register activities
	registerActivities(w)

	go func() {
		log.Info().Str("taskQueue", inv.taskQueue).Msg("Starting worker")
		if err := w.Start(); err != nil {
			log.Fatal().Err(err).Msg("Unable to start Temporal worker")
		}
	}()
}
