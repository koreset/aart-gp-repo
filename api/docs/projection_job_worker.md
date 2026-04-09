# Projection Job Worker

## Overview

This document describes the background worker implementation for processing projection jobs. The worker continuously monitors the `projection_jobs` table for jobs with "Queued" status and processes them in creation date order. If any job has an "In Progress" status, the worker waits until that job is completed before processing any new jobs.

## Implementation Details

The implementation consists of two main functions in the `services/projections.go` file:

1. `StartProjectionJobWorker()`: Starts the background worker in a separate goroutine. It uses a mutex and a boolean flag to ensure only one worker is running at a time.

2. `ProcessQueuedProjectionJobs()`: Does the actual work of processing jobs:
   - Checks if there are any "In Progress" jobs and waits if there are
   - Gets the oldest "Queued" job by creation date
   - Updates the job status to "In Progress"
   - Processes each product in the job
   - Updates the job status to "Complete" or "Failed" based on the result
   - Clears caches after job completion

The worker is started automatically when the application initializes, after database migrations have been run.

## Job Status Values

Projection jobs can have the following status values:

- "Queued": Job is waiting to be processed
- "In Progress": Job is currently being processed
- "Complete": Job has been successfully processed
- "Failed": Job processing failed (check StatusError field for details)

## Testing the Implementation

To test the background worker:

1. Create multiple projection jobs with "Queued" status:
   - Use the existing API endpoints to create projection jobs
   - Verify that the jobs are created with "Queued" status in the database

2. Verify that the worker processes jobs in creation date order:
   - Create jobs with different creation dates
   - Monitor the logs to see the order in which jobs are processed
   - Check the database to see that the oldest jobs are processed first

3. Verify that the worker respects "In Progress" jobs:
   - Manually set a job status to "In Progress" in the database
   - Create new "Queued" jobs
   - Verify that the worker does not process any new jobs until the "In Progress" job is completed
   - Change the "In Progress" job status to "Complete"
   - Verify that the worker starts processing the queued jobs

4. Verify error handling:
   - Create a job that will fail during processing (e.g., with invalid product IDs)
   - Verify that the job status is set to "Failed" and the StatusError field contains the error message
   - Verify that the worker continues to process other jobs after a failure

## Monitoring

The worker logs its activities using the application's logging system. Look for log messages with the following patterns:

- "Starting projection job worker": Indicates that the worker has started
- "Found X jobs in progress, waiting...": Indicates that the worker is waiting for in-progress jobs to complete
- "Processing queued job ID: X, Name: Y": Indicates that the worker has started processing a job
- "Completed processing job ID: X": Indicates that the worker has completed processing a job
- Error messages: Indicate that there was a problem processing a job

You can also monitor the database directly to see the status of jobs and their processing order.