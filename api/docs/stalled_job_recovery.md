# Stalled Projection Job Recovery

This document describes the mechanism for detecting and recovering stalled projection jobs that may have been interrupted by a server shutdown.

## Overview

When a projection job is running and the server shuts down unexpectedly, the job will be left with a status of "In Progress" in the database. When the server restarts, the stalled job recovery mechanism will detect these jobs and mark them as "Failed" so they can be requeued if needed.

## Implementation Details

The stalled job recovery mechanism consists of two parts:

1. **Startup Recovery**: When the server starts up, the `RecoverStalledJobs()` function is called to find any jobs with "In Progress" status and mark them as "Failed".

2. **Continuous Monitoring**: During normal operation, the `ProcessQueuedProjectionJobs()` function checks for any jobs with "In Progress" status and marks them as "Failed" if they are found. This provides an additional layer of protection in case any stalled jobs were missed during startup.

## Testing the Solution

To test the stalled job recovery mechanism, follow these steps:

1. **Create a Test Job**:
   - Create a new projection job through the normal interface
   - Manually update its status to "In Progress" in the database:
     ```sql
     UPDATE projection_jobs SET status = 'In Progress' WHERE id = <job_id>;
     ```

2. **Restart the Server**:
   - Stop the server
   - Start the server again

3. **Verify Recovery**:
   - Check the server logs for messages like:
     ```
     Found <n> stalled jobs during server startup. Recovering...
     Recovering stalled job ID: <job_id>, Name: <job_name>. Marking as failed.
     ```
   - Verify that the job's status has been updated to "Failed" in the database:
     ```sql
     SELECT id, run_name, status, status_error FROM projection_jobs WHERE id = <job_id>;
     ```
   - The `status_error` field should contain the message "Job was interrupted by server shutdown"

## Monitoring and Maintenance

If stalled jobs are occurring frequently, it might indicate other issues that need to be addressed. Consider monitoring the frequency of stalled jobs and investigating any patterns.

## Future Enhancements

For a more robust solution, consider implementing a heartbeat mechanism where running jobs periodically update a "last_updated" timestamp. This would allow you to distinguish between truly stalled jobs and jobs that are just taking a long time to process.