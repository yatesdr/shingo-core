package engine

import (
	"log"
	"time"

	"shingoedge/store"
)

// HourlyTracker accumulates counter deltas into hourly buckets in the database.
type HourlyTracker struct {
	db *store.DB
}

// NewHourlyTracker creates a new HourlyTracker.
func NewHourlyTracker(db *store.DB) *HourlyTracker {
	return &HourlyTracker{db: db}
}

// HandleDelta records a counter delta into the current date/hour bucket.
func (ht *HourlyTracker) HandleDelta(delta CounterDeltaEvent) {
	if delta.LineID == 0 || delta.JobStyleID == 0 {
		return
	}

	now := time.Now()
	countDate := now.Format("2006-01-02")
	hour := now.Hour()

	if err := ht.db.UpsertHourlyCount(delta.LineID, delta.JobStyleID, countDate, hour, delta.Delta); err != nil {
		log.Printf("hourly tracker upsert: %v", err)
	}
}
