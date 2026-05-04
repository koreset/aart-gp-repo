package services

import (
	"fmt"
	"time"

	"api/models"
)

// ApplyEffectiveExitDeactivations flips in-force members whose effective_exit_date
// has been reached to status = 'INACTIVE'. The frontend stores the chosen exit date
// as-is and leaves status = 'ACTIVE' for future-dated exits; this sweeper transitions
// those members on or after the chosen date.
func ApplyEffectiveExitDeactivations() {
	now := time.Now()
	res := DB.Model(&models.GPricingMemberDataInForce{}).
		Where("effective_exit_date IS NOT NULL").
		Where("effective_exit_date <= ?", now).
		Where("UPPER(COALESCE(status, '')) <> ?", "INACTIVE").
		Update("status", "INACTIVE")
	if res.Error != nil {
		fmt.Printf("[scheduler] member exit-date sweep failed: %v\n", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		fmt.Printf("[scheduler] flipped %d member(s) to INACTIVE based on effective_exit_date\n", res.RowsAffected)
	}
}

// StartMemberExitDateSweeper runs ApplyEffectiveExitDeactivations once at boot
// (so a server restart picks up any exits that elapsed during downtime) and then
// once every 24 hours.
func StartMemberExitDateSweeper() {
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			fmt.Println("Sweeping members past their effective exit date...")
			ApplyEffectiveExitDeactivations()
			<-ticker.C
		}
	}()
}
