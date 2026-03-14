package cron

import "context"

// Job represents a scheduled job with optional metadata
type Job struct {
	Name     string
	Schedule string
	Runner   func(ctx context.Context) error
}

type Scheduler interface {
	// Start begins executing scheduled jobs
	Start()

	// Stop halts all scheduled jobs gracefully
	Stop() error

	// RegisterJob adds a new job to the scheduler
	// Returns a job ID that can be used to remove the job later
	RegisterJob(ctx context.Context, job *Job) (int, error)

	// RemoveJob removes a job by its ID
	RemoveJob(id int) error
}
