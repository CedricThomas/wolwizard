package jobs

import (
	"context"
	"log"

	"github.com/CedricThomas/console/internal/config"
	"github.com/CedricThomas/console/internal/controller"
	"github.com/CedricThomas/console/internal/input/cron"
)

// RegisterPCAgent registers periodic jobs
func RegisterPCAgent(
	ctx context.Context,
	scheduler cron.Scheduler,
	pcAgentContoller controller.PCAgent,
	cfg *config.Config,
) error {
	job := &cron.Job{
		Name:     "metrics-collection",
		Schedule: cfg.MetricsReportingSchedule,
		Runner:   pcAgentContoller.SendCurrentHostAsyncMetrics,
	}

	jobID, err := scheduler.RegisterJob(ctx, job)
	if err != nil {
		return err
	}

	log.Printf("Successfully registered job %s with ID %d", job.Name, jobID)
	return nil
}
