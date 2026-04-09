package services

import (
	"fmt"
	"sync"
	"time"
)

// UploadStatus represents the status of a long-running upload/import task
// This is an in-memory tracker meant to improve UX by allowing clients to poll progress
// Note: In-memory means status will reset on process restart; for durability use Redis/DB.

type UploadState string

const (
	UploadStatePending    UploadState = "PENDING"
	UploadStateProcessing UploadState = "PROCESSING"
	UploadStateCompleted  UploadState = "COMPLETED"
	UploadStateFailed     UploadState = "FAILED"
)

type UploadStatus struct {
	ID         string      `json:"id"`
	Filename   string      `json:"filename"`
	ProductID  int         `json:"product_id"`
	Year       int         `json:"year"`
	MpVersion  string      `json:"mp_version"`
	TotalBytes int64       `json:"total_bytes"`
	BytesDone  int64       `json:"bytes_done"`
	Progress   float64     `json:"progress"` // Percentage (0.0 to 100.0)
	State      UploadState `json:"state"`
	Error      string      `json:"error,omitempty"`
	StartedAt  time.Time   `json:"started_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

var (
	uploadStatusMu sync.RWMutex
	uploadStatuses = map[string]*UploadStatus{}
	uploadTTL      = 24 * time.Hour
)

func uploadRedisKey(id string) string {
	return fmt.Sprintf("ads:v1:upload:%s", id)
}

func CreateUploadStatus(id string, filename string, productID int, year int, mpVersion string, totalBytes int64) *UploadStatus {
	uploadStatusMu.Lock()
	defer uploadStatusMu.Unlock()
	st := &UploadStatus{
		ID:         id,
		Filename:   filename,
		ProductID:  productID,
		Year:       year,
		MpVersion:  mpVersion,
		TotalBytes: totalBytes,
		BytesDone:  0,
		State:      UploadStatePending,
		StartedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	uploadStatuses[id] = st
	// Persist to Redis if available so other instances can read the status
	if RedisAvailable() {
		RedisSetJSON(uploadRedisKey(id), st, uploadTTL)
	}
	return st
}

func GetUploadStatus(id string) (*UploadStatus, bool) {
	uploadStatusMu.RLock()
	st, ok := uploadStatuses[id]
	uploadStatusMu.RUnlock()
	if ok {
		return st, true
	}
	// Fallback to Redis for cross-instance visibility
	if RedisAvailable() {
		var rst UploadStatus
		if RedisGetJSON(uploadRedisKey(id), &rst) {
			return &rst, true
		}
	}
	return nil, false
}

func SetUploadProcessing(id string) {
	uploadStatusMu.Lock()
	if st, ok := uploadStatuses[id]; ok {
		st.State = UploadStateProcessing
		st.UpdatedAt = time.Now()
		// Persist to Redis
		if RedisAvailable() {
			RedisSetJSON(uploadRedisKey(id), st, uploadTTL)
		}
	}
	uploadStatusMu.Unlock()
}

func SetUploadCompleted(id string) {
	uploadStatusMu.Lock()
	if st, ok := uploadStatuses[id]; ok {
		st.State = UploadStateCompleted
		st.BytesDone = st.TotalBytes
		st.Progress = 100.0
		st.UpdatedAt = time.Now()
		if RedisAvailable() {
			RedisSetJSON(uploadRedisKey(id), st, uploadTTL)
		}
	}
	uploadStatusMu.Unlock()
}

func SetUploadFailed(id string, err error) {
	uploadStatusMu.Lock()
	if st, ok := uploadStatuses[id]; ok {
		st.State = UploadStateFailed
		st.Error = err.Error()
		st.UpdatedAt = time.Now()
		if RedisAvailable() {
			RedisSetJSON(uploadRedisKey(id), st, uploadTTL)
		}
	}
	uploadStatusMu.Unlock()
}

func AddUploadProgress(id string, n int64) {
	uploadStatusMu.Lock()
	if st, ok := uploadStatuses[id]; ok {
		st.BytesDone += n
		if st.TotalBytes > 0 {
			st.Progress = float64(st.BytesDone) / float64(st.TotalBytes) * 100.0
			if st.Progress > 100.0 {
				st.Progress = 100.0
			}
		}
		st.UpdatedAt = time.Now()
		if RedisAvailable() {
			RedisSetJSON(uploadRedisKey(id), st, uploadTTL)
		}
	}
	uploadStatusMu.Unlock()
}
