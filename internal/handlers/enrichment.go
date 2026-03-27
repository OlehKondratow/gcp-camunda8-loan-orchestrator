package handlers

import (
	"context"
	"fmt"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"go.uber.org/zap"
)

const enrichmentVariablesJSON = `{"clientCreditScore":750}`

// EnrichmentHandler runs the data-enrichment service task: logs job metadata, simulates
// enrichment with a fixed credit score, and completes the job in Zeebe.
type EnrichmentHandler struct {
	log *zap.Logger
}

func NewEnrichmentHandler(log *zap.Logger) *EnrichmentHandler {
	return &EnrichmentHandler{log: log}
}

// Handle implements worker.JobHandler.
func (h *EnrichmentHandler) Handle(client worker.JobClient, job entities.Job) {
	ctx := context.Background()

	h.log.Info("enrichment job activated",
		zap.Int64("job_key", job.GetKey()),
		zap.Int64("process_instance_key", job.GetProcessInstanceKey()),
		zap.String("job_type", job.GetType()),
	)

	dispatch, err := client.NewCompleteJobCommand().
		JobKey(job.GetKey()).
		VariablesFromString(enrichmentVariablesJSON)
	if err != nil {
		h.failJob(ctx, client, job, fmt.Errorf("complete job variables: %w", err))
		return
	}

	if _, err := dispatch.Send(ctx); err != nil {
		h.failJob(ctx, client, job, fmt.Errorf("complete job: %w", err))
		return
	}

	h.log.Info("enrichment job completed",
		zap.Int64("job_key", job.GetKey()),
		zap.Int64("process_instance_key", job.GetProcessInstanceKey()),
	)
}

func (h *EnrichmentHandler) failJob(ctx context.Context, client worker.JobClient, job entities.Job, jobErr error) {
	h.log.Error("enrichment job failed", zap.Error(jobErr),
		zap.Int64("job_key", job.GetKey()),
		zap.Int64("process_instance_key", job.GetProcessInstanceKey()),
	)

	retries := job.GetRetries() - 1
	if retries < 0 {
		retries = 0
	}

	_, err := client.NewFailJobCommand().
		JobKey(job.GetKey()).
		Retries(retries).
		ErrorMessage(jobErr.Error()).
		Send(ctx)
	if err != nil {
		h.log.Error("fail job command failed", zap.Error(err),
			zap.Int64("job_key", job.GetKey()),
		)
	}
}
