package controllers

import (
	"api/log"
	"api/models"
	"api/services"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// parseDashboardQuery extracts the shared QuotePerformanceQuery filter
// envelope from a gin context. From/To accept RFC3339 (YYYY-MM-DDTHH:mm:ss)
// or YYYY-MM-DD; everything else is multi-value query strings.
func parseDashboardQuery(c *gin.Context) models.QuotePerformanceQuery {
	q := models.QuotePerformanceQuery{
		Users:               c.QueryArray("users"),
		Region:              c.QueryArray("region"),
		QuoteType:           c.QueryArray("quote_type"),
		DistributionChannel: c.QueryArray("distribution_channel"),
	}
	if v := c.Query("from"); v != "" {
		if t, err := parseDashboardDate(v); err == nil {
			q.From = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := parseDashboardDate(v); err == nil {
			q.To = &t
		}
	}
	return q
}

func parseDashboardDate(v string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, v); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse date %q (expected RFC3339 or YYYY-MM-DD)", v)
}

// dashboardContext builds a request-scoped logger context that mirrors
// the pattern used elsewhere in this controller file (see GenerateGroupPricingQuote).
func dashboardContext(c *gin.Context) context.Context {
	requestID, exists := c.Get("requestID")
	var ctx context.Context
	if exists {
		ctx = context.WithValue(context.Background(), log.RequestIDKey, requestID.(string))
	} else {
		ctx = context.Background()
	}
	userEmail, emailExists := c.Get("userEmail")
	userName, nameExists := c.Get("userName")
	if emailExists && nameExists {
		ctx = log.ContextWithUserInfo(ctx, userEmail.(string), userName.(string))
	}
	return ctx
}

// GetQuotePerformanceKpis is GET /group-pricing/dashboard/kpis.
func GetQuotePerformanceKpis(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)
	q := parseDashboardQuery(c)
	rows, err := services.GetQuotePerformanceKpis(q)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to load quote performance KPIs")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// GetQuotePerformanceFunnel is GET /group-pricing/dashboard/funnel.
func GetQuotePerformanceFunnel(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)
	q := parseDashboardQuery(c)
	stages, err := services.GetQuoteFunnel(q)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to load quote funnel")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stages)
}

// GetQuotePerformanceTrend is GET /group-pricing/dashboard/trend?bucket=...
func GetQuotePerformanceTrend(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)
	q := parseDashboardQuery(c)
	bucket := c.DefaultQuery("bucket", "day")
	rows, err := services.GetQuoteTrend(q, bucket)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to load quote trend")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// GetQuotePerformanceSlaBreaches is GET /group-pricing/dashboard/sla-breaches.
func GetQuotePerformanceSlaBreaches(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)
	q := parseDashboardQuery(c)
	summary, err := services.GetQuoteSlaBreaches(q)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to compute SLA breaches")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

// PostQuoteExtract is POST /group-pricing/dashboard/extract. Accepts a
// full QuoteExtractFilter in the body (multi-select arrays can grow past
// URL length limits, hence POST not GET).
func PostQuoteExtract(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)
	var f models.QuoteExtractFilter
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rows, total, err := services.ExtractQuotes(f)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to extract quotes")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"rows":      rows,
		"total":     total,
		"page":      f.Page,
		"page_size": f.PageSize,
	})
}

// GetQuoteExtractXlsx is GET /group-pricing/dashboard/extract.xlsx. The
// filter is query-encoded so the browser can hit it via a normal <a
// download> tag. Caps at 50k rows to keep excelize's in-memory model
// bounded.
func GetQuoteExtractXlsx(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)

	f := models.QuoteExtractFilter{
		CreatedBy:           c.QueryArray("created_by"),
		Reviewer:            c.QueryArray("reviewer"),
		Region:              c.QueryArray("region"),
		QuoteType:           c.QueryArray("quote_type"),
		Industry:            c.QueryArray("industry"),
		DistributionChannel: c.QueryArray("distribution_channel"),
		OrderBy:             c.Query("order_by"),
	}
	for _, s := range c.QueryArray("status") {
		f.Status = append(f.Status, models.Status(s))
	}
	if v := c.Query("min_annual_premium"); v != "" {
		if x, err := strconv.ParseFloat(v, 64); err == nil {
			f.MinAnnualPremium = &x
		}
	}
	if v := c.Query("max_annual_premium"); v != "" {
		if x, err := strconv.ParseFloat(v, 64); err == nil {
			f.MaxAnnualPremium = &x
		}
	}
	if v := c.Query("from"); v != "" {
		if t, err := parseDashboardDate(v); err == nil {
			f.From = &t
		}
	}
	if v := c.Query("to"); v != "" {
		if t, err := parseDashboardDate(v); err == nil {
			f.To = &t
		}
	}
	// Page through up to the export cap.
	f.Page = 1
	f.PageSize = 50_000

	rows, total, err := services.ExtractQuotes(f)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to fetch quotes for xlsx export")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if total > 50_000 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "Extract exceeds 50,000 rows. Please narrow the filter and try again.",
			"total": total,
		})
		return
	}

	data, err := services.ExportQuoteExtractXlsx(rows)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to render quote extract xlsx")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	filename := "quote-performance-extract-" + time.Now().UTC().Format("20060102-150405") + ".xlsx"
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

// GetGroupPricingQuoteStatusHistory is GET /group-pricing/quotes/:id/status-history.
func GetGroupPricingQuoteStatusHistory(c *gin.Context) {
	quoteId := c.Param("id")
	audits, err := services.GetGroupPricingQuoteStatusAudit(quoteId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, audits)
}

// GetQuoteSlaTargets is GET /group-pricing/dashboard/sla-targets.
func GetQuoteSlaTargets(c *gin.Context) {
	targets, err := services.ListQuoteSlaTargets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, targets)
}

// PostQuoteSlaTarget handles both create (no ID) and update (with ID) via
// the same payload — the service does an upsert keyed on the natural
// (from_status, to_status, quote_type) tuple.
func PostQuoteSlaTarget(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)
	var t models.QuoteSlaTarget
	if err := c.BindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	saved, err := services.UpsertQuoteSlaTarget(t, user)
	if err != nil {
		logger.WithField("error", err.Error()).Warn("Failed to upsert SLA target")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, saved)
}

// DeleteQuoteSlaTarget is DELETE /group-pricing/dashboard/sla-targets/:id.
// Soft-disables the target rather than hard-deleting so historical
// breach calculations remain reproducible.
func DeleteQuoteSlaTarget(c *gin.Context) {
	ctx := dashboardContext(c)
	logger := log.WithContext(ctx)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	user := c.MustGet("user").(models.AppUser)
	if err := services.DeactivateQuoteSlaTarget(id, user); err != nil {
		logger.WithField("error", err.Error()).Error("Failed to deactivate SLA target")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
