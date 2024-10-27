package worker

import (
	"context"
	"sync"

	"github.com/brownei/chifunds-api/types"
	"go.uber.org/zap"
)

type Worker struct {
	ctx    context.Context
	logger *zap.SugaredLogger
}

func NewWorker(ctx context.Context, logger *zap.SugaredLogger) *Worker {
	return &Worker{
		ctx:    ctx,
		logger: logger,
	}
}

func (w *Worker) RunQueriesWithWorkerPool(jobs []types.TransferJob, workerCount int) {
	jobChan := make(chan types.TransferJob, len(jobs))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go w.worker(&wg, jobChan)
	}

	// Send jobs to the job channel
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan) // Close job channel after all jobs are sent

	// Wait for all workers to finish
	wg.Wait()
}

// worker executes queries from the job channel
func (w *Worker) worker(wg *sync.WaitGroup, jobs <-chan types.TransferJob) {
	defer wg.Done()

	for job := range jobs {
		if err := job.ExecuteQuery(w.ctx, job.Query, job.Args); err != nil {
			w.logger.Infof("Error executing job %d: %v", job.Id, err)
		} else {
			w.logger.Infof("Successfully executed job %d", job.Id)
		}
	}
}
