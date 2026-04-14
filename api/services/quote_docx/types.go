// Package quote_docx provides functionality to generate Group Risk Quotation documents
// as DOCX (Microsoft Word) files, mirroring the TypeScript implementation from the frontend.
// It handles document assembly, XML generation, and ZIP packaging to produce standards-compliant OOXML files.
package quote_docx

// Types for Group Risk Quotation DOCX generation
// Mirrors TypeScript interfaces from docxQuote.ts

// QuoteTotals aggregates totals across all categories
type QuoteTotals struct {
	TotalLives           int
	TotalSumAssured      float64
	TotalAnnualSalary    float64
	TotalAnnualPremium   float64
}

// LabelValueRow is a simple label-value pair for key-value tables
type LabelValueRow struct {
	Label string
	Value string
}

// PremiumSummaryRow is a row for the Premium Summary table
type PremiumSummaryRow struct {
	Category       string
	MemberCount    string
	TotalSalary    string
	TotalSumAssured string
	AnnualPremium  string
	PercentSalary  string
}

// GroupFuneralRow is a row for the Group Funeral table
type GroupFuneralRow struct {
	Category          string
	MemberCount       string
	MonthlyPremium    string
	AnnualPremium     string
	TotalAnnualPremium string
}

// PremiumBreakdownRow is per-category premium breakdown
type PremiumBreakdownRow struct {
	Benefit          string
	TotalSumAssured  string
	AnnualPremium    string
	PercentSalary    string
}

// BenefitDefinitionRow is a row for the 7-column benefits and definitions table
type BenefitDefinitionRow struct {
	Benefit          string
	SalaryMultiple   string
	BenefitStructure string
	WaitingPeriod    string
	DeferredPeriod   string
	CoverDefinition  string
	RiskType         string
}

// FuneralCoverageRow is group funeral coverage per category
type FuneralCoverageRow struct {
	Member   string
	SumAssured    interface{} // can be string or number
	MaxCovered    interface{} // can be string or number
}

// EducatorBenefitRow is educator benefits data
type EducatorBenefitRow struct {
	Level       string
	MaxTuition  interface{} // can be string or number
	MaxCoverage interface{} // can be string or number
}

// BenefitTitles are resolved benefit label names
type BenefitTitles struct {
	GlaBenefitTitle           string
	SglaBenefitTitle          string
	PtdBenefitTitle           string
	CiBenefitTitle            string
	PhiBenefitTitle           string
	TtdBenefitTitle           string
	FamilyFuneralBenefitTitle string
}

// ===== Design Constants =====

const (
	Font = "Arial"

	// Colours (hex strings)
	ColorPrimary      = "34495E"
	ColorSecondary    = "34495E"
	ColorAccent       = "E74C3C"
	ColorLightFill    = "ECF0F1"
	ColorDark         = "2C3E50"
	ColorAltRow       = "FAFAFA"
	ColorWhite        = "FFFFFF"
	ColorNavy         = "1E3A5F"
	ColorOrange       = "D47600"
	ColorLightOrange  = "FFF8F0"
	ColorLightBlue    = "F1F8FF"
	ColorMediumGray   = "586069"
	ColorInputLine    = "C8C8C8"

	// Page sizes in DXA (1 inch = 1440 DXA)
	A4Width  = 11906
	A4Height = 16838

	// Margins in DXA
	MarginTop    = 850
	MarginBottom = 850
	MarginLeft   = 1020
	MarginRight  = 1020

	// Content widths (portrait and landscape)
	ContentWidth            = A4Width - MarginLeft - MarginRight           // ~9866 DXA
	LandscapeContentWidth   = A4Height - MarginLeft - MarginRight          // ~14798 DXA

	// Font sizes in half-points (multiply pt by 2)
	SizeTitle      = 32  // 16pt
	SizeHeading    = 28  // 14pt
	SizeSubheading = 24  // 12pt
	SizeBody       = 20  // 10pt
	SizeCaption    = 18  // 9pt
	SizeSmall      = 16  // 8pt
)
