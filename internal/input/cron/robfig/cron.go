package robfig

import (
	"context"
	"log"

	icron "github.com/CedricThomas/console/internal/input/cron"
	robfigcron "github.com/robfig/cron/v3"
)

// robfigScheduler implements icron.Scheduler using the robfig/cron library
type robfigScheduler struct {
	cron *robfigcron.Cron
}

// NewRobfigScheduler creates a new icron scheduler backed by robfig/cron
func NewRobfigScheduler() icron.Scheduler {
	return &robfigScheduler{
		cron: robfigcron.New(robfigcron.WithSeconds()),
	}
}

// Start begins executing scheduled jobs
func (r *robfigScheduler) Start() {
	r.cron.Start()
	log.Println("Cron scheduler started")
}

// Stop halts all scheduled jobs gracefully
func (r *robfigScheduler) Stop() error {
	log.Println("Stopping cron scheduler...")
	done := make(chan struct{})
	go func() {
		r.cron.Stop()
		close(done)
	}()
	<-done
	log.Println("Cron scheduler stopped")
	return nil
}

// RegisterJob adds a new job to the scheduler
// Returns a job ID that can be used to remove the job later
func (r *robfigScheduler) RegisterJob(ctx context.Context, job *icron.Job) (int, error) {
	if job.Runner == nil {
		return 0, nil
	}

	jobID, err := r.cron.AddFunc(job.Schedule, func() {
		// TODO fix
		if err := job.Runner(ctx); err != nil {
			log.Printf("error in job %q: %v", job.Name, err)
		}
	})
	if err != nil {
		return 0, err
	}

	log.Printf("registered job %q with schedule %s (ID: %d)", job.Name, job.Schedule, jobID)
	return int(jobID), nil
}

// RemoveJob removes a job by its ID
func (r *robfigScheduler) RemoveJob(id int) error {
	log.Printf("removing job with ID %d", id)
	r.cron.Remove(robfigcron.EntryID(id))
	return nil
}
