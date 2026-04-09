package services

import "sync"

// jobLocks provides a shard of mutexes keyed by Projection Job ID to minimize contention
// and prevent different jobs from blocking each other.
var jobLocks sync.Map // map[int]*sync.Mutex

// GetJobMutex returns a per-job mutex for the given projection job ID.
// It is safe for concurrent use.
func GetJobMutex(jobID int) *sync.Mutex {
	if v, ok := jobLocks.Load(jobID); ok {
		return v.(*sync.Mutex)
	}
	m := &sync.Mutex{}
	actual, _ := jobLocks.LoadOrStore(jobID, m)
	return actual.(*sync.Mutex)
}

// ReleaseJobMutex removes the mutex for the jobID from the registry.
// Only call this after you're certain no goroutines will use the lock again.
func ReleaseJobMutex(jobID int) {
	jobLocks.Delete(jobID)
}
