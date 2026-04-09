package services

import (
	appLog "api/log"
	"api/models"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

// winFeatureNames defines the fixed feature vector order used for training and scoring.
var winFeatureNames = []string{
	"is_renewal",
	"distribution_channel",
	"discount_pct",
	"commission_pct",
	"total_loading_pct",
	"member_count",
	"avg_age",
	"avg_income",
	"gender_ratio",
	"has_experience_data",
	"expected_loss_ratio",
	"broker_historic_rate",
	"days_to_commencement",
	"benefit_count",
}

// ── Logistic regression primitives ───────────────────────────────────────────

func sigmoid(z float64) float64 {
	return 1.0 / (1.0 + math.Exp(-z))
}

// zscore normalises a feature vector in-place given pre-computed means and std devs.
func zscoreNorm(features, means, stddevs []float64) []float64 {
	norm := make([]float64, len(features))
	for i, v := range features {
		sd := stddevs[i]
		if sd == 0 {
			norm[i] = 0
		} else {
			norm[i] = (v - means[i]) / sd
		}
	}
	return norm
}

// computeMeansStdDevs returns column-wise means and sample std devs for an n×p matrix.
func computeMeansStdDevs(X [][]float64) (means, stddevs []float64) {
	if len(X) == 0 {
		return nil, nil
	}
	p := len(X[0])
	n := float64(len(X))
	means = make([]float64, p)
	stddevs = make([]float64, p)

	for _, row := range X {
		for j, v := range row {
			means[j] += v
		}
	}
	for j := range means {
		means[j] /= n
	}
	for _, row := range X {
		for j, v := range row {
			diff := v - means[j]
			stddevs[j] += diff * diff
		}
	}
	for j := range stddevs {
		if n > 1 {
			stddevs[j] = math.Sqrt(stddevs[j] / (n - 1))
		}
	}
	return means, stddevs
}

// trainLogisticRegression fits a binary logistic regression using gradient descent with L2 reg.
// X is already z-score normalised. Returns weight vector including bias (w[0] = bias).
func trainLogisticRegression(X [][]float64, y []float64, epochs int, lr, lambda float64) []float64 {
	if len(X) == 0 {
		return nil
	}
	p := len(X[0])
	// weights[0] = bias; weights[1..p] = feature weights
	weights := make([]float64, p+1)

	for epoch := 0; epoch < epochs; epoch++ {
		grad := make([]float64, p+1)
		for i, row := range X {
			z := weights[0]
			for j, xv := range row {
				z += weights[j+1] * xv
			}
			pred := sigmoid(z)
			err := pred - y[i]
			grad[0] += err
			for j, xv := range row {
				grad[j+1] += err * xv
			}
		}
		n := float64(len(X))
		weights[0] -= lr * grad[0] / n
		for j := 1; j <= p; j++ {
			// L2 regularisation (skip bias)
			weights[j] -= lr * (grad[j]/n + lambda*weights[j])
		}
		_ = epoch
	}
	return weights
}

// predictProb returns P(win) for a single z-score-normalised feature vector.
func predictProb(weights, normFeatures []float64) float64 {
	z := weights[0]
	for j, xv := range normFeatures {
		z += weights[j+1] * xv
	}
	return sigmoid(z)
}

// evaluateModel returns accuracy and trapezoidal AUC on a held-out set.
// X is already z-score normalised using training means/stddevs.
func evaluateModel(weights []float64, X [][]float64, y []float64) (accuracy, auc float64) {
	if len(X) == 0 {
		return 0, 0
	}
	type pair struct{ prob, label float64 }
	pairs := make([]pair, len(X))
	correct := 0
	for i, row := range X {
		p := predictProb(weights, row)
		pairs[i] = pair{p, y[i]}
		predicted := 0.0
		if p >= 0.5 {
			predicted = 1.0
		}
		if predicted == y[i] {
			correct++
		}
	}
	accuracy = float64(correct) / float64(len(X))

	// Trapezoidal AUC
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].prob > pairs[j].prob })
	var tp, fp, prevTp, prevFp float64
	var pos, neg float64
	for _, pr := range pairs {
		if pr.label == 1 {
			pos++
		} else {
			neg++
		}
	}
	if pos == 0 || neg == 0 {
		return accuracy, 0.5
	}
	for _, pr := range pairs {
		if pr.label == 1 {
			tp++
		} else {
			fp++
		}
		auc += (fp/neg - prevFp/neg) * (tp/pos + prevTp/pos) / 2
		prevTp, prevFp = tp, fp
	}
	return accuracy, auc
}

// ── Feature extraction ────────────────────────────────────────────────────────

type winFeatureRow struct {
	quoteID  int
	features []float64
	label    float64 // 1=win, 0=loss; -1=unlabelled (non-terminal)
}

func extractFeatures(
	quote models.GroupPricingQuote,
	summary models.MemberRatingResultSummary,
	brokerHistoricRate float64,
) []float64 {
	// 1. is_renewal
	isRenewal := 0.0
	if quote.QuoteType == "Renewal" {
		isRenewal = 1.0
	}

	// 2. distribution_channel (broker=0, direct=1, binder=2, tied_agent=3)
	channelMap := map[models.DistributionChannel]float64{
		models.ChannelBroker:    0,
		models.ChannelDirect:    1,
		models.ChannelBinder:    2,
		models.ChannelTiedAgent: 3,
	}
	channel := channelMap[quote.DistributionChannel]

	// 3. discount_pct
	discountPct := quote.Loadings.Discount

	// 4. commission_pct (stored as decimal, convert to %)
	commissionPct := quote.Loadings.CommissionLoading * 100

	// 5. total_loading_pct (sum of all loadings minus discount)
	totalLoading := quote.Loadings.CommissionLoading*100 +
		quote.Loadings.ProfitLoading*100 +
		quote.Loadings.ExpenseLoading*100 +
		quote.Loadings.AdminLoading*100 +
		quote.Loadings.ContingencyLoading*100 +
		quote.Loadings.OtherLoading*100 -
		quote.Loadings.Discount

	// 6. member_count (log-scaled)
	memberCount := math.Log1p(summary.MemberCount)

	// 7. avg_age
	avgAge := float64(quote.MemberAverageAge)

	// 8. avg_income (log-scaled)
	avgIncome := math.Log1p(quote.MemberAverageIncome)

	// 9. gender_ratio
	genderRatio := quote.MemberMaleFemaleDistribution

	// 10. has_experience_data
	hasExperience := 0.0
	if quote.ClaimsExperienceCount > 0 {
		hasExperience = 1.0
	}

	// 11. expected_loss_ratio (capped at 2.0)
	elr := 0.0
	if summary.TotalAnnualPremium > 0 {
		elr = summary.TotalExpectedClaims / summary.TotalAnnualPremium
		if elr > 2.0 {
			elr = 2.0
		}
	}

	// 12. broker_historic_rate
	brokerRate := brokerHistoricRate

	// 13. days_to_commencement
	daysToCommencement := 0.0
	if !quote.CommencementDate.IsZero() && !quote.CreationDate.IsZero() {
		diff := quote.CommencementDate.Sub(quote.CreationDate)
		daysToCommencement = diff.Hours() / 24
	}

	// 14. benefit_count
	benefitCount := float64(len(quote.SelectedSchemeCategories))
	if benefitCount < 1 {
		benefitCount = 1
	}

	return []float64{
		isRenewal,
		channel,
		discountPct,
		commissionPct,
		totalLoading,
		memberCount,
		avgAge,
		avgIncome,
		genderRatio,
		hasExperience,
		elr,
		brokerRate,
		daysToCommencement,
		benefitCount,
	}
}

// loadBrokerHistoricRates builds a map[brokerID]acceptanceRate for all brokers
// in the training set. Returns 0.5 for brokers with fewer than 5 quotes.
func loadBrokerHistoricRates() map[int]float64 {
	type row struct {
		BrokerID      int   `gorm:"column:broker_id"`
		TotalCount    int64 `gorm:"column:total_count"`
		AcceptedCount int64 `gorm:"column:accepted_count"`
	}
	var rows []row
	DB.Table("group_pricing_quotes").
		Select(`broker_id,
			COUNT(*) as total_count,
			SUM(CASE WHEN status IN ('accepted','in_force') THEN 1 ELSE 0 END) as accepted_count`).
		Group("broker_id").
		Scan(&rows)

	result := make(map[int]float64, len(rows))
	for _, r := range rows {
		if r.TotalCount < 5 {
			result[r.BrokerID] = 0.5
		} else {
			result[r.BrokerID] = float64(r.AcceptedCount) / float64(r.TotalCount)
		}
	}
	return result
}

// ── Training ─────────────────────────────────────────────────────────────────

// TrainWinProbabilityModel fetches terminal quotes, trains a logistic regression,
// evaluates it on a held-out set, and persists the model. A cold-start guard
// prevents training when fewer than 20 terminal quotes are available.
func TrainWinProbabilityModel(triggeredBy string) error {
	logger := appLog.WithFields(map[string]interface{}{
		"action":       "TrainWinProbabilityModel",
		"triggered_by": triggeredBy,
	})
	logger.Info("Starting win probability model training")

	// Terminal statuses: wins and losses
	winStatuses := []string{string(models.StatusAccepted), string(models.StatusInForce)}
	lossStatuses := []string{
		string(models.StatusRejected),
		string(models.StatusNotTakenUp),
		string(models.StatusExpired),
		string(models.StatusOutOfForce),
		string(models.StatusCancelled),
	}
	terminalStatuses := append(winStatuses, lossStatuses...)

	var terminalQuotes []models.GroupPricingQuote
	DB.Where("status IN ?", terminalStatuses).
		Preload("SchemeCategories").
		Order("creation_date ASC").
		Find(&terminalQuotes)

	if len(terminalQuotes) < 20 {
		logger.Info("Insufficient training data — skipping model training")
		return nil
	}

	// Build broker historic rates map
	brokerRates := loadBrokerHistoricRates()

	// Collect feature rows
	var rows []winFeatureRow
	for _, q := range terminalQuotes {
		var summary models.MemberRatingResultSummary
		DB.Where("quote_id = ?", q.ID).First(&summary)

		// Skip quotes with missing data
		if summary.TotalAnnualPremium == 0 || len(q.SelectedSchemeCategories) < 1 {
			continue
		}

		brokerRate, ok := brokerRates[q.QuoteBroker.ID]
		if !ok {
			brokerRate = 0.5
		}

		features := extractFeatures(q, summary, brokerRate)

		label := 0.0
		for _, ws := range winStatuses {
			if string(q.Status) == ws {
				label = 1.0
				break
			}
		}

		rows = append(rows, winFeatureRow{quoteID: q.ID, features: features, label: label})
	}

	if len(rows) < 20 {
		logger.Info("Insufficient usable training rows — skipping model training")
		return nil
	}

	// Chronological 80/20 split
	splitIdx := int(float64(len(rows)) * 0.8)
	trainRows := rows[:splitIdx]
	testRows := rows[splitIdx:]

	// Build raw X matrices for normalisation
	trainX := make([][]float64, len(trainRows))
	trainY := make([]float64, len(trainRows))
	for i, r := range trainRows {
		trainX[i] = r.features
		trainY[i] = r.label
	}

	means, stddevs := computeMeansStdDevs(trainX)

	// Normalise training set
	normTrainX := make([][]float64, len(trainX))
	for i, row := range trainX {
		normTrainX[i] = zscoreNorm(row, means, stddevs)
	}

	weights := trainLogisticRegression(normTrainX, trainY, 500, 0.01, 0.001)

	// Normalise test set and evaluate
	normTestX := make([][]float64, len(testRows))
	testY := make([]float64, len(testRows))
	for i, r := range testRows {
		normTestX[i] = zscoreNorm(r.features, means, stddevs)
		testY[i] = r.label
	}
	accuracy, auc := evaluateModel(weights, normTestX, testY)

	logger.WithFields(map[string]interface{}{
		"training_size": len(trainRows),
		"test_size":     len(testRows),
		"accuracy":      FloatPrecision(accuracy*100, 2),
		"auc":           FloatPrecision(auc, 4),
	}).Info("Win probability model trained")

	// Marshal arrays to JSON for storage
	weightsJSON, _ := json.Marshal(weights)
	meansJSON, _ := json.Marshal(means)
	stddevsJSON, _ := json.Marshal(stddevs)
	featureNamesJSON, _ := json.Marshal(winFeatureNames)

	model := models.WinProbabilityModel{
		Weights:        string(weightsJSON),
		FeatureNames:   string(featureNamesJSON),
		FeatureMeans:   string(meansJSON),
		FeatureStdDevs: string(stddevsJSON),
		TrainingSize:   len(trainRows),
		Accuracy:       FloatPrecision(accuracy*100, 2),
		AUC:            FloatPrecision(auc, 4),
		TrainedAt:      time.Now(),
		CreatedBy:      triggeredBy,
	}
	if err := DB.Create(&model).Error; err != nil {
		logger.WithField("error", err.Error()).Error("Failed to persist win probability model")
		return err
	}

	logger.WithField("model_id", model.ID).Info("Win probability model persisted — triggering rescore")
	go func() {
		if err := RescoreAllQuotes(); err != nil {
			appLog.WithField("error", err.Error()).Error("RescoreAllQuotes failed after retraining")
		}
	}()
	return nil
}

// ── Scoring ───────────────────────────────────────────────────────────────────

// ScoreQuote computes and persists the win probability for a single quote.
// Returns nil (no error) if no model is available yet (cold start).
func ScoreQuote(quoteID int) (*models.QuoteWinProbability, error) {
	cacheKey := fmt.Sprintf("win-prob-quote-%d", quoteID)

	// Check in-memory cache (1 h)
	if cached, ok := Cache.Get(cacheKey); ok {
		if wp, ok2 := cached.(*models.QuoteWinProbability); ok2 {
			return wp, nil
		}
	}

	// Load latest model
	var model models.WinProbabilityModel
	if err := DB.Order("id DESC").First(&model).Error; err != nil {
		// No model yet — cold start
		return nil, nil
	}

	// Parse stored arrays
	var weights, means, stddevs []float64
	if err := json.Unmarshal([]byte(model.Weights), &weights); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(model.FeatureMeans), &means); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(model.FeatureStdDevs), &stddevs); err != nil {
		return nil, err
	}

	// Load quote
	var quote models.GroupPricingQuote
	if err := DB.Where("id = ?", quoteID).Preload("SchemeCategories").First(&quote).Error; err != nil {
		return nil, err
	}

	var summary models.MemberRatingResultSummary
	DB.Where("quote_id = ?", quoteID).First(&summary)

	brokerRates := loadBrokerHistoricRates()
	brokerRate, ok := brokerRates[quote.QuoteBroker.ID]
	if !ok {
		brokerRate = 0.5
	}

	features := extractFeatures(quote, summary, brokerRate)
	normFeatures := zscoreNorm(features, means, stddevs)
	score := predictProb(weights, normFeatures)
	scorePct := FloatPrecision(score*100, 1)

	// Band assignment
	band := "low"
	switch {
	case scorePct >= 75:
		band = "very_high"
	case scorePct >= 50:
		band = "high"
	case scorePct >= 25:
		band = "medium"
	}

	// All feature contributions sorted by absolute impact (|weight × normalised value|, signed).
	// Storing all 14 allows the frontend explainer to render a full SHAP-style breakdown.
	type featureContrib struct {
		Name         string  `json:"name"`
		Contribution float64 `json:"contribution"`
		Weight       float64 `json:"weight"`
	}
	contribs := make([]featureContrib, len(winFeatureNames))
	for i, name := range winFeatureNames {
		contribution := 0.0
		w := 0.0
		if i+1 < len(weights) {
			w = weights[i+1]
			contribution = w * normFeatures[i]
		}
		contribs[i] = featureContrib{Name: name, Contribution: FloatPrecision(contribution, 4), Weight: FloatPrecision(w, 4)}
	}
	sort.Slice(contribs, func(i, j int) bool {
		return math.Abs(contribs[i].Contribution) > math.Abs(contribs[j].Contribution)
	})
	topFeaturesJSON, _ := json.Marshal(contribs)

	wp := &models.QuoteWinProbability{
		QuoteID:     quoteID,
		Score:       FloatPrecision(score, 4),
		ScorePct:    scorePct,
		Band:        band,
		TopFeatures: string(topFeaturesJSON),
		ModelID:     model.ID,
		ScoredAt:    time.Now(),
	}

	// Upsert — explicit find-then-save so wp retains all freshly computed fields.
	var existing models.QuoteWinProbability
	if DB.Where("quote_id = ?", quoteID).First(&existing).Error == nil {
		wp.ID = existing.ID
		DB.Save(wp)
	} else {
		DB.Create(wp)
	}

	// Cache in-memory (1 h)
	Cache.Set(cacheKey, wp, 1)
	// Cache in Redis (24 h)
	RedisSetJSON(cacheKey, wp, 24*time.Hour)

	return wp, nil
}

// RescoreAllQuotes rescores all non-terminal quotes in parallel (goroutines).
func RescoreAllQuotes() error {
	activeStatuses := []string{
		string(models.StatusInProgress),
		string(models.StatusPendingReview),
		string(models.StatusApproved),
	}
	var activeQuotes []models.GroupPricingQuote
	DB.Where("status IN ?", activeStatuses).Select("id").Find(&activeQuotes)

	for _, q := range activeQuotes {
		id := q.ID
		go func(qid int) {
			if _, err := ScoreQuote(qid); err != nil {
				appLog.WithField("quote_id", qid).WithField("error", err.Error()).
					Warn("ScoreQuote failed during batch rescore")
			}
		}(id)
	}
	appLog.WithField("count", len(activeQuotes)).Info("Batch rescore dispatched")
	return nil
}

// GetWinProbabilityModelInfo returns the latest model metadata.
func GetWinProbabilityModelInfo() (*models.WinProbabilityModel, error) {
	var model models.WinProbabilityModel
	if err := DB.Order("id DESC").First(&model).Error; err != nil {
		return nil, nil // no model yet
	}
	// Don't expose raw weight arrays in the info response
	model.Weights = ""
	model.FeatureMeans = ""
	model.FeatureStdDevs = ""
	return &model, nil
}

// GetQuoteWinProbabilityScore returns the stored score for a quote, or nil if not yet scored.
func GetQuoteWinProbabilityScore(quoteID int) (*models.QuoteWinProbability, error) {
	// Try scoring on demand if not yet stored
	wp, err := ScoreQuote(quoteID)
	return wp, err
}

// StartWinProbabilityRetrainingJob starts a weekly background retraining job.
func StartWinProbabilityRetrainingJob() {
	go func() {
		for {
			now := time.Now()
			// Next Sunday at 02:00
			daysUntilSunday := (7 - int(now.Weekday())) % 7
			if daysUntilSunday == 0 {
				daysUntilSunday = 7
			}
			nextRun := time.Date(now.Year(), now.Month(), now.Day()+daysUntilSunday, 2, 0, 0, 0, now.Location())
			wait := time.Until(nextRun)
			appLog.WithField("next_retrain", nextRun.Format(time.RFC3339)).
				Info("Win probability retraining job scheduled")
			time.Sleep(wait)
			appLog.Info("Running scheduled win probability model retraining")
			if err := TrainWinProbabilityModel("scheduled"); err != nil {
				appLog.WithField("error", err.Error()).Error("Scheduled win probability training failed")
			}
		}
	}()
}

// QuoteWinProbabilityBandCounts returns count of active quotes per probability band for the dashboard.
func QuoteWinProbabilityBandCounts() map[string]int64 {
	bands := []string{"0-20", "20-40", "40-60", "60-80", "80-100"}
	result := map[string]int64{
		"0-20":   0,
		"20-40":  0,
		"40-60":  0,
		"60-80":  0,
		"80-100": 0,
	}

	activeStatuses := []string{
		string(models.StatusInProgress),
		string(models.StatusPendingReview),
		string(models.StatusApproved),
	}

	var scores []models.QuoteWinProbability
	DB.Table("quote_win_probabilities").
		Joins("JOIN group_pricing_quotes q ON q.id = quote_win_probabilities.quote_id").
		Where("q.status IN ?", activeStatuses).
		Select("score_pct").
		Find(&scores)

	for _, s := range scores {
		switch {
		case s.ScorePct < 20:
			result["0-20"]++
		case s.ScorePct < 40:
			result["20-40"]++
		case s.ScorePct < 60:
			result["40-60"]++
		case s.ScorePct < 80:
			result["60-80"]++
		default:
			result["80-100"]++
		}
	}

	_ = bands
	return result
}

// intQuoteID parses a string quoteID to int for goroutine hooks.
func intQuoteID(s string) int {
	id, _ := strconv.Atoi(s)
	return id
}
