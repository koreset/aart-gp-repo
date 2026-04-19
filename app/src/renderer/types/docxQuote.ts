/**
 * TypeScript interfaces for the DOCX quote generation feature.
 * Used by both quoteDataHelpers.ts and useDocxQuoteGeneration.ts.
 */

/** Colour palette matching the PDF design constants */
export interface DocxColorConfig {
  primary: string // Dark blue-gray header fills, headings
  secondary: string // Secondary text colour
  accent: string // Red accent (sparingly used)
  lightFill: string // Light grey (label cells, banners)
  dark: string // Body text
  altRow: string // Alternating row shading
  white: string
}

/** Font configuration */
export interface DocxFontConfig {
  primary: string
  sizes: {
    title: number // half-points (e.g. 32 = 16pt)
    heading: number
    subheading: number
    body: number
    caption: number
  }
}

/** Full document configuration */
export interface DocxQuoteConfig {
  colors: DocxColorConfig
  fonts: DocxFontConfig
  margins: {
    top: number // DXA units
    bottom: number
    left: number
    right: number
  }
  pageSize: {
    width: number // DXA units
    height: number
  }
  contentWidth: number // computed: pageSize.width - margins.left - margins.right
}

/** Benefit title labels resolved from benefit maps */
export interface BenefitTitles {
  glaBenefitTitle: string
  sglaBenefitTitle: string
  ptdBenefitTitle: string
  ciBenefitTitle: string
  phiBenefitTitle: string
  ttdBenefitTitle: string
  familyFuneralBenefitTitle: string
  additionalAccidentalGlaBenefitTitle?: string
  additionalGlaCoverBenefitTitle?: string
}

/** Aggregated totals across all categories */
export interface QuoteTotals {
  totalLives: number
  totalSumAssured: number
  totalAnnualSalary: number
  totalAnnualPremium: number
}

/** A simple label-value pair for key-value tables */
export interface LabelValueRow {
  label: string
  value: string
}

/** Row data for the Premium Summary table */
export interface PremiumSummaryRow {
  category: string
  memberCount: string
  totalSalary: string
  totalSumAssured: string
  annualPremium: string
  percentSalary: string
}

/** Row data for the Group Funeral table */
export interface GroupFuneralRow {
  category: string
  memberCount: string
  monthlyPremium: string
  annualPremium: string
  totalAnnualPremium: string
}

/** Row data for per-category premium breakdown */
export interface PremiumBreakdownRow {
  benefit: string
  totalSumAssured: string
  annualPremium: string
  percentSalary: string
}

/** Row data for the 7-column benefits and definitions table */
export interface BenefitDefinitionRow {
  benefit: string
  salaryMultiple: string
  benefitStructure: string
  waitingPeriod: string
  deferredPeriod: string
  coverDefinition: string
  riskType: string
}

/** Row data for group funeral coverage per category */
export interface FuneralCoverageRow {
  member: string
  sumAssured: string | number
  maxCovered: string | number
}

/** Row data for educator benefits */
export interface EducatorBenefitRow {
  level: string
  maxTuition: string | number
  maxCoverage: string | number
}
