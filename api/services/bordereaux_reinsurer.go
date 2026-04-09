package services

import (
	"api/models"
	"fmt"
	"strconv"
	"time"
)

// CreateReinsurerAcceptance creates a new acceptance record for a submitted bordereaux
func CreateReinsurerAcceptance(req models.CreateReinsurerAcceptanceRequest, user models.AppUser) (models.ReinsurerAcceptance, error) {
	acc := models.ReinsurerAcceptance{
		GeneratedBordereauxID: req.GeneratedBordereauxID,
		ReinsurerName:         req.ReinsurerName,
		ReinsurerCode:         req.ReinsurerCode,
		SubmittedAmount:       req.SubmittedAmount,
		DueDate:               req.DueDate,
		Notes:                 req.Notes,
		Status:                "pending",
		CreatedBy:             user.UserName,
	}
	if err := DB.Create(&acc).Error; err != nil {
		return acc, fmt.Errorf("failed to create reinsurer acceptance: %w", err)
	}
	return acc, nil
}

// GetReinsurerAcceptances returns acceptance records filtered by bordereaux ID and/or status
func GetReinsurerAcceptances(generatedID, status string) ([]models.ReinsurerAcceptance, error) {
	var records []models.ReinsurerAcceptance
	q := DB.Order("created_at DESC")
	if generatedID != "" {
		q = q.Where("generated_bordereaux_id = ?", generatedID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// UpdateReinsurerAcceptance updates the status and financial fields of an acceptance record
func UpdateReinsurerAcceptance(id int, req models.UpdateReinsurerAcceptanceRequest, user models.AppUser) (models.ReinsurerAcceptance, error) {
	var acc models.ReinsurerAcceptance
	if err := DB.First(&acc, id).Error; err != nil {
		return acc, fmt.Errorf("acceptance record %d not found: %w", id, err)
	}

	if req.Status != "" {
		acc.Status = req.Status
	}
	if req.AcceptedAmount > 0 {
		acc.AcceptedAmount = req.AcceptedAmount
		acc.Variance = acc.SubmittedAmount - req.AcceptedAmount
	}
	if req.QueryDetails != "" {
		acc.QueryDetails = req.QueryDetails
	}
	if req.Notes != "" {
		acc.Notes = req.Notes
	}
	if req.ReceivedDate != nil {
		acc.ReceivedDate = req.ReceivedDate
	}
	acc.UpdatedBy = user.UserName
	acc.UpdatedAt = time.Now()

	if err := DB.Save(&acc).Error; err != nil {
		return acc, fmt.Errorf("failed to update acceptance record: %w", err)
	}
	return acc, nil
}

// GetAcceptanceStats returns a summary of acceptance statuses for a given bordereaux
func GetAcceptanceStats(generatedID string) (models.AcceptanceStats, error) {
	var records []models.ReinsurerAcceptance
	q := DB.Where("generated_bordereaux_id = ?", generatedID)
	if err := q.Find(&records).Error; err != nil {
		return models.AcceptanceStats{}, err
	}

	stats := models.AcceptanceStats{Total: len(records)}
	for _, r := range records {
		switch r.Status {
		case "pending":
			stats.Pending++
		case "accepted":
			stats.Accepted++
		case "queried":
			stats.Queried++
		case "rejected":
			stats.Rejected++
		}
		stats.TotalSubmitted += r.SubmittedAmount
		stats.TotalAccepted += r.AcceptedAmount
		stats.TotalVariance += r.Variance
	}
	return stats, nil
}

// CreateReinsurerRecovery records a claim recovery received from a reinsurer
func CreateReinsurerRecovery(req models.CreateReinsurerRecoveryRequest, user models.AppUser) (models.ReinsurerRecovery, error) {
	pct := 0.0
	if req.ClaimAmount > 0 {
		pct = (req.RecoveredAmount / req.ClaimAmount) * 100
	}

	status := "pending"
	if req.RecoveredAmount >= req.ClaimAmount && req.ClaimAmount > 0 {
		status = "full"
	} else if req.RecoveredAmount > 0 {
		status = "partial"
	}

	rec := models.ReinsurerRecovery{
		GeneratedBordereauxID: req.GeneratedBordereauxID,
		ClaimReference:        req.ClaimReference,
		ReinsurerName:         req.ReinsurerName,
		ReinsurerCode:         req.ReinsurerCode,
		ClaimAmount:           req.ClaimAmount,
		RecoveredAmount:       req.RecoveredAmount,
		RecoveryPercentage:    pct,
		Status:                status,
		Notes:                 req.Notes,
		RecordedBy:            user.UserName,
	}
	if err := DB.Create(&rec).Error; err != nil {
		return rec, fmt.Errorf("failed to create recovery record: %w", err)
	}
	return rec, nil
}

// GetReinsurerRecoveries returns recovery records filtered by bordereaux ID, claim ref, and/or status
func GetReinsurerRecoveries(generatedID, claimRef, status string) ([]models.ReinsurerRecovery, error) {
	var records []models.ReinsurerRecovery
	q := DB.Order("created_at DESC")
	if generatedID != "" {
		q = q.Where("generated_bordereaux_id = ?", generatedID)
	}
	if claimRef != "" {
		q = q.Where("claim_reference LIKE ?", "%"+claimRef+"%")
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// UpdateReinsurerRecovery updates a recovery record's amounts and status
func UpdateReinsurerRecovery(id int, req models.UpdateReinsurerRecoveryRequest, user models.AppUser) (models.ReinsurerRecovery, error) {
	var rec models.ReinsurerRecovery
	if err := DB.First(&rec, id).Error; err != nil {
		return rec, fmt.Errorf("recovery record %d not found: %w", id, err)
	}

	if req.RecoveredAmount > 0 {
		rec.RecoveredAmount = req.RecoveredAmount
		if rec.ClaimAmount > 0 {
			rec.RecoveryPercentage = (req.RecoveredAmount / rec.ClaimAmount) * 100
		}
	}
	if req.Status != "" {
		rec.Status = req.Status
	}
	if req.Notes != "" {
		rec.Notes = req.Notes
	}
	if req.ReceivedDate != nil {
		rec.ReceivedDate = req.ReceivedDate
	}
	rec.UpdatedAt = time.Now()

	if err := DB.Save(&rec).Error; err != nil {
		return rec, fmt.Errorf("failed to update recovery record: %w", err)
	}
	return rec, nil
}

// GetRecoveryStats returns summary recovery statistics for a given bordereaux
func GetRecoveryStats(generatedID string) (models.RecoveryStats, error) {
	var records []models.ReinsurerRecovery
	q := DB.Where("generated_bordereaux_id = ?", generatedID)
	if err := q.Find(&records).Error; err != nil {
		return models.RecoveryStats{}, err
	}

	stats := models.RecoveryStats{Total: len(records)}
	for _, r := range records {
		switch r.Status {
		case "pending":
			stats.Pending++
		case "partial":
			stats.Partial++
		case "full":
			stats.Full++
		case "disputed":
			stats.Disputed++
		}
		stats.TotalClaimAmount += r.ClaimAmount
		stats.TotalRecovered += r.RecoveredAmount
	}
	stats.TotalOutstanding = stats.TotalClaimAmount - stats.TotalRecovered
	return stats, nil
}

// GenerateClaimRecovery calculates the ceded portion of an approved/paid claim using the
// scheme's active treaty and creates a ReinsurerRecovery record for tracking. This is the
// downstream link between claim notifications and actual recovery amounts.
func GenerateClaimRecovery(claim models.GroupSchemeClaim, user models.AppUser) (*models.ReinsurerRecovery, error) {
	if claim.ClaimAmount <= 0 {
		return nil, nil // nothing to recover
	}

	// Find active treaty for the claim's scheme
	treaties, err := GetActiveTreatiesForScheme(claim.SchemeId)
	if err != nil || len(treaties) == 0 {
		return nil, nil // no treaty — fully retained
	}
	treaty := treaties[0]

	// Calculate ceded portion
	ceded, _, belowRetention := CalculateClaimCession(claim.ClaimAmount, treaty)
	if belowRetention || ceded <= 0 {
		return nil, nil // fully retained
	}

	// Check for existing recovery on this claim to avoid duplicates
	var existing int64
	DB.Model(&models.ReinsurerRecovery{}).
		Where("claim_reference = ?", claim.ClaimNumber).
		Count(&existing)
	if existing > 0 {
		return nil, nil // already tracked
	}

	pct := (ceded / claim.ClaimAmount) * 100

	rec := models.ReinsurerRecovery{
		ClaimReference:     claim.ClaimNumber,
		ReinsurerName:      treaty.ReinsurerName,
		ReinsurerCode:      treaty.ReinsurerCode,
		ClaimAmount:        claim.ClaimAmount,
		RecoveredAmount:    0, // not yet recovered — pending from reinsurer
		RecoveryPercentage: pct,
		Status:             "pending",
		Notes:              fmt.Sprintf("Auto-generated on claim %s (status: %s). Ceded %.2f%%.", claim.ClaimNumber, claim.Status, pct),
		RecordedBy:         user.UserName,
	}
	if err := DB.Create(&rec).Error; err != nil {
		return nil, fmt.Errorf("failed to create recovery record: %w", err)
	}

	// Write audit trail
	_ = writeAudit(DB, AuditContext{
		Area:      "group-pricing",
		Entity:    "reinsurer_recoveries",
		EntityID:  strconv.Itoa(rec.ID),
		Action:    "CREATE",
		ChangedBy: user.UserName,
	}, struct{}{}, rec)

	return &rec, nil
}
