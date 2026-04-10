package services

import (
	appLog "api/log"
	"api/models"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	// maxConcurrentCalculations limits how many quote calculations run
	// simultaneously on this instance, preventing resource exhaustion.
	maxConcurrentCalculations = 2
)

var (
	calcWorkerOnce    sync.Once
	calcWorkerRunning bool
	calcWorkerMu      sync.Mutex
	calcSemaphore     chan struct{}
)

// EnqueueCalculationJob creates a queued calculation job and notifies the user
// of their position in the queue via WebSocket. Returns the job ID.
func EnqueueCalculationJob(quoteID int, basis string, credibility float64, user models.AppUser) (int, error) {
	job := models.CalculationJob{
		QuoteID:     quoteID,
		Basis:       basis,
		Credibility: credibility,
		UserEmail:   user.UserEmail,
		UserName:    user.UserName,
		Status:      "queued",
	}

	if err := DB.Create(&job).Error; err != nil {
		return 0, fmt.Errorf("failed to enqueue calculation job: %w", err)
	}

	// Count how many jobs are ahead in the queue
	var position int64
	DB.Model(&models.CalculationJob{}).
		Where("status = ? AND id <= ?", "queued", job.ID).
		Count(&position)

	sendCalculationProgress(user.UserEmail, CalculationProgress{
		QuoteID:       strconv.Itoa(quoteID),
		JobID:         job.ID,
		Phase:         "queued",
		QueuePosition: int(position),
		Progress:      0,
	})

	return job.ID, nil
}

// GetCalculationJob returns the current status of a calculation job.
func GetCalculationJob(jobID int) (models.CalculationJob, error) {
	var job models.CalculationJob
	err := DB.Where("id = ?", jobID).First(&job).Error
	return job, err
}

// StartCalculationJobWorker starts the background worker that polls for
// queued calculation jobs and processes them with bounded concurrency.
// Safe to call multiple times — only the first call starts the worker.
func StartCalculationJobWorker() {
	calcWorkerMu.Lock()
	defer calcWorkerMu.Unlock()

	if calcWorkerRunning {
		appLog.Info("Calculation job worker is already running")
		return
	}

	calcSemaphore = make(chan struct{}, maxConcurrentCalculations)
	calcWorkerRunning = true

	go func() {
		appLog.Info("Calculation job worker started")
		for {
			processNextCalculationJob()
			// Randomized sleep to avoid thundering herd across instances
			time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)
		}
	}()
}

// processNextCalculationJob finds the next queued job, claims it atomically,
// and runs the calculation. The semaphore limits concurrent calculations.
func processNextCalculationJob() {
	// Check if we have capacity before even querying
	select {
	case calcSemaphore <- struct{}{}:
		// Acquired a slot
	default:
		// All slots occupied — skip this poll cycle
		return
	}

	// We hold a semaphore slot; release it when done
	go func() {
		defer func() { <-calcSemaphore }()

		// Find the oldest queued job
		var job models.CalculationJob
		err := DB.Where("status = ?", "queued").
			Order("queued_at ASC").
			Limit(1).
			First(&job).Error
		if err != nil {
			// No queued jobs or DB error
			return
		}

		// Atomically claim the job — prevents other instances from picking it up
		result := DB.Model(&models.CalculationJob{}).
			Where("id = ? AND status = ?", job.ID, "queued").
			Updates(map[string]interface{}{
				"status":     "processing",
				"started_at": time.Now(),
			})
		if result.Error != nil || result.RowsAffected == 0 {
			// Lost the race to another instance
			return
		}

		logger := appLog.WithFields(map[string]interface{}{
			"job_id":   job.ID,
			"quote_id": job.QuoteID,
			"user":     job.UserEmail,
			"action":   "CalculationJobWorker",
		})
		logger.Info("Processing calculation job")

		// Notify user that calculation is starting
		sendCalculationProgress(job.UserEmail, CalculationProgress{
			QuoteID: strconv.Itoa(job.QuoteID),
			JobID:   job.ID,
			Phase:   "loading_data",
			Progress: 0,
		})

		// Broadcast updated queue positions to other waiting users
		broadcastQueuePositions()

		// Run the actual calculation
		user := models.AppUser{UserEmail: job.UserEmail, UserName: job.UserName}
		calcErr := CalculateGroupPricingQuote(strconv.Itoa(job.QuoteID), job.Basis, job.Credibility, user)

		now := time.Now()
		if calcErr != nil {
			logger.WithField("error", calcErr.Error()).Error("Calculation job failed")
			DB.Model(&job).Updates(map[string]interface{}{
				"status":       "failed",
				"error":        calcErr.Error(),
				"completed_at": now,
			})
			sendCalculationProgress(job.UserEmail, CalculationProgress{
				QuoteID:  strconv.Itoa(job.QuoteID),
				JobID:    job.ID,
				Phase:    "failed",
				Progress: 0,
			})
		} else {
			logger.Info("Calculation job completed successfully")
			DB.Model(&job).Updates(map[string]interface{}{
				"status":       "completed",
				"completed_at": now,
			})
			// The "completed" progress event is already sent by CalculateGroupPricingQuote
		}

		// Re-broadcast queue positions since a slot freed up
		broadcastQueuePositions()
	}()
}

// broadcastQueuePositions sends updated queue position to all users
// who currently have queued jobs, so their UI shows the correct position.
func broadcastQueuePositions() {
	var queuedJobs []models.CalculationJob
	DB.Where("status = ?", "queued").Order("queued_at ASC").Find(&queuedJobs)

	for i, job := range queuedJobs {
		sendCalculationProgress(job.UserEmail, CalculationProgress{
			QuoteID:       strconv.Itoa(job.QuoteID),
			JobID:         job.ID,
			Phase:         "queued",
			QueuePosition: i + 1,
			Progress:      0,
		})
	}
}
