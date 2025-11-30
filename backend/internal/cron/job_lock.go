package cron

import (
	"sync"
	"time"
)

// JobLock prevents overlapping cron job executions
type JobLock struct {
	mu      sync.Mutex
	running bool
	lastRun time.Time
}

// NewJobLock creates a new job lock
func NewJobLock() *JobLock {
	return &JobLock{}
}

// TryLock attempts to acquire the lock
// Returns true if lock was acquired, false if already running
func (jl *JobLock) TryLock() bool {
	jl.mu.Lock()
	defer jl.mu.Unlock()
	
	if jl.running {
		return false // Job already running
	}
	
	jl.running = true
	jl.lastRun = time.Now()
	return true
}

// Unlock releases the lock
func (jl *JobLock) Unlock() {
	jl.mu.Lock()
	defer jl.mu.Unlock()
	jl.running = false
}

// IsRunning checks if job is currently running
func (jl *JobLock) IsRunning() bool {
	jl.mu.Lock()
	defer jl.mu.Unlock()
	return jl.running
}

// LastRun returns the last run time
func (jl *JobLock) LastRun() time.Time {
	jl.mu.Lock()
	defer jl.mu.Unlock()
	return jl.lastRun
}