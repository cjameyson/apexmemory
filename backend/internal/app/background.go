package app

import (
	"context"
	"sync"
	"time"
)

// BackgroundJobs manages periodic background tasks.
type BackgroundJobs struct {
	app    *Application
	wg     sync.WaitGroup
	stopCh chan struct{}
}

// NewBackgroundJobs creates a new background job manager.
func NewBackgroundJobs(app *Application) *BackgroundJobs {
	return &BackgroundJobs{
		app:    app,
		stopCh: make(chan struct{}),
	}
}

// Start begins all background jobs.
func (bg *BackgroundJobs) Start() {
	bg.wg.Add(1)
	go bg.sessionCleanupLoop()
}

// Stop gracefully stops all background jobs.
func (bg *BackgroundJobs) Stop() {
	close(bg.stopCh)
	bg.wg.Wait()
}

// sessionCleanupLoop periodically removes expired sessions.
func (bg *BackgroundJobs) sessionCleanupLoop() {
	defer bg.wg.Done()

	// Run immediately on startup, then every hour
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Initial cleanup
	bg.cleanupSessions()

	for {
		select {
		case <-ticker.C:
			bg.cleanupSessions()
		case <-bg.stopCh:
			return
		}
	}
}

func (bg *BackgroundJobs) cleanupSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := bg.app.CleanupExpiredSessions(ctx); err != nil {
		bg.app.Logger.Error("failed to cleanup expired sessions", "error", err)
	} else {
		bg.app.Logger.Debug("expired sessions cleaned up")
	}
}
