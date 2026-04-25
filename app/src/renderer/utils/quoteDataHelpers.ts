/**
 * Shared data preparation functions for quote document generation.
 *
 * These are extracted from QuoteOutput.vue so both generatePDF() and
 * generateDocxQuote() work from identical calculations.
 */

import type {
  QuoteTotals,
  LabelValueRow,
  PremiumSummaryRow,
  GroupFuneralRow,
  PremiumBreakdownRow,
  BenefitDefinitionRow,
  FuneralCoverageRow,
  EducatorBenefitRow,
  BenefitTitles
} from '@/renderer/types/docxQuote'
import formatDateString from '@/renderer/utils/helpers.js'

// ---------------------------------------------------------------------------
// Low-level utilities
// ---------------------------------------------------------------------------

/**
 * Safely traverse a nested object path, returning defaultValue on failure.
 * Extracted from QuoteOutput.vue line 332.
 */
export function safeGetValue(
  obj: any,
  path: string,
  defaultValue: any = null
): any {
  try {
    return (
      path.split('.').reduce((current, key) => current?.[key], obj) ??
      defaultValue
    )
  } catch {
    return defaultValue
  }
}

/**
 * Round up to 2 decimal places and format with space-separated thousands.
 * Extracted from QuoteOutput.vue line 343.
 */
export function roundUpToTwoDecimalsAccounting(num: number): string {
  const roundedNum = Math.ceil(num * 100) / 100
  return roundedNum
    .toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    })
    .replace(/,/g, ' ')
}

/** Return value or '-' if falsy. */
export function dashIfEmpty(value: any): string {
  return value || '-'
}

// ---------------------------------------------------------------------------
// Office-premium derivation
// ---------------------------------------------------------------------------
//
// Office premium is no longer persisted on the rating result summary. It is
// derived on the fly from the risk premium and the scheme-level loading
// (expense + commission + profit).

/**
 * Sum of the three scheme-level loadings (expense, commission, profit) that
 * make up the office-premium denominator. Each loading is a fraction
 * (e.g. 0.05 for 5%).
 */
export function schemeTotalLoading(s: {
  expense_loading?: number
  commission_loading?: number
  profit_loading?: number
}): number {
  return (
    (s?.expense_loading ?? 0) +
    (s?.commission_loading ?? 0) +
    (s?.profit_loading ?? 0)
  )
}

/**
 * Derive the office premium from a risk premium and the scheme-level
 * loading on the summary row. Guards against denom <= 0.
 */
export function computeOfficePremium(
  riskPremium: number,
  s: {
    expense_loading?: number
    commission_loading?: number
    profit_loading?: number
  }
): number {
  const denom = 1 - schemeTotalLoading(s)
  return denom <= 0 ? 0 : (riskPremium ?? 0) / denom
}

/**
 * Convert a risk-rate-per-1000-SA into the equivalent office rate by
 * scaling by 1 / (1 - schemeTotalLoading).
 */
export function officeRateFromRiskRate(
  riskRatePer1000: number,
  s: {
    expense_loading?: number
    commission_loading?: number
    profit_loading?: number
  }
): number {
  const denom = 1 - schemeTotalLoading(s)
  return denom <= 0 ? 0 : (riskRatePer1000 ?? 0) / denom
}

/**
 * Convert a risk-premium-proportion-of-salary into its office equivalent.
 */
export function officeProportionFromRiskProportion(
  riskProportion: number,
  s: {
    expense_loading?: number
    commission_loading?: number
    profit_loading?: number
  }
): number {
  const denom = 1 - schemeTotalLoading(s)
  return denom <= 0 ? 0 : (riskProportion ?? 0) / denom
}

// ---------------------------------------------------------------------------
// Aggregation helpers
// ---------------------------------------------------------------------------

/**
 * Calculate totals across all result summaries.
 * Extracted from QuoteOutput.vue lines 952-967.
 */
export function calculateQuoteTotals(resultSummaries: any[]): QuoteTotals {
  return {
    totalLives: resultSummaries.reduce(
      (sum, item) => sum + item.member_count,
      0
    ),
    totalSumAssured: resultSummaries.reduce(
      (sum, item) => sum + item.total_sum_assured,
      0
    ),
    totalAnnualSalary: resultSummaries.reduce(
      (sum, item) => sum + item.total_annual_salary,
      0
    ),
    totalAnnualPremium: resultSummaries.reduce(
      (sum, item) => sum + item.total_annual_premium,
      0
    )
  }
}

/**
 * Determine whether any category has non-funeral benefits.
 * Extracted from QuoteOutput.vue lines 1040-1049.
 */
export function hasAnyNonFuneralBenefits(resultSummaries: any[]): boolean {
  return resultSummaries.some(
    (item) =>
      safeGetValue(item, 'total_gla_capped_sum_assured', 0) > 0 ||
      safeGetValue(item, 'total_ptd_capped_sum_assured', 0) > 0 ||
      safeGetValue(item, 'total_ci_capped_sum_assured', 0) > 0 ||
      safeGetValue(item, 'total_sgla_capped_sum_assured', 0) > 0 ||
      safeGetValue(item, 'total_phi_capped_income', 0) > 0 ||
      safeGetValue(item, 'total_ttd_capped_income', 0) > 0
  )
}

/**
 * Check if a single category has non-funeral benefits.
 */
export function categoryHasNonFuneralBenefits(item: any): boolean {
  return (
    safeGetValue(item, 'total_gla_capped_sum_assured', 0) > 0 ||
    safeGetValue(item, 'total_ptd_capped_sum_assured', 0) > 0 ||
    safeGetValue(item, 'total_ci_capped_sum_assured', 0) > 0 ||
    safeGetValue(item, 'total_sgla_capped_sum_assured', 0) > 0 ||
    safeGetValue(item, 'total_phi_capped_income', 0) > 0 ||
    safeGetValue(item, 'total_ttd_capped_income', 0) > 0
  )
}

// ---------------------------------------------------------------------------
// Table data builders — Section 3 (Quote Summary)
// ---------------------------------------------------------------------------

/**
 * Build the initial info / quote summary key-value rows.
 * Extracted from QuoteOutput.vue lines 970-996.
 */
export function buildInitialInfoRows(
  quote: any,
  totals: QuoteTotals
): LabelValueRow[] {
  return [
    { label: 'Type of Policy:', value: 'Group Risk Assurance' },
    { label: 'Quote Number:', value: `${quote.quote_name}` },
    {
      label: 'Quote Date:',
      value: `${formatDateString(quote.creation_date, true, true, true)}`
    },
    { label: 'Scheme Name:', value: `${quote.scheme_name}` },
    {
      label: 'Inception Date:',
      value: `${formatDateString(quote.commencement_date, true, true, true)}`
    },
    { label: 'Number of Lives Covered:', value: `${totals.totalLives}` },
    {
      label: 'Total Sum Assured:',
      value: roundUpToTwoDecimalsAccounting(totals.totalSumAssured)
    },
    {
      label: 'Total Annual Salary:',
      value: roundUpToTwoDecimalsAccounting(totals.totalAnnualSalary)
    },
    {
      label: 'Total Annual Premium:',
      value: roundUpToTwoDecimalsAccounting(totals.totalAnnualPremium)
    }
  ]
}

// ---------------------------------------------------------------------------
// Table data builders — Section 4 (Premium Summary + Group Funeral)
// ---------------------------------------------------------------------------

/**
 * Build rows for the Premium Summary table (excluding the header row).
 * Extracted from QuoteOutput.vue lines 1085-1140.
 * Returns data rows + a totals row.
 */
export function buildPremiumSummaryRows(
  resultSummaries: any[]
): PremiumSummaryRow[] {
  const rows: PremiumSummaryRow[] = []

  resultSummaries.forEach((item) => {
    if (!categoryHasNonFuneralBenefits(item)) return

    rows.push({
      category: item.category,
      memberCount: item.member_count.toString(),
      totalSalary: roundUpToTwoDecimalsAccounting(item.total_annual_salary),
      totalSumAssured: roundUpToTwoDecimalsAccounting(item.total_sum_assured),
      annualPremium: roundUpToTwoDecimalsAccounting(
        item.exp_total_annual_premium_excl_funeral
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(item.proportion_exp_total_premium_excl_funeral_salary * 100)}%`
    })
  })

  // Totals row
  const totalLives = resultSummaries.reduce(
    (sum, item) => sum + item.member_count,
    0
  )
  const totalSalary = resultSummaries.reduce(
    (sum, item) => sum + item.total_annual_salary,
    0
  )
  const totalSumAssured = resultSummaries.reduce(
    (sum, item) => sum + item.total_sum_assured,
    0
  )
  const totalPremium = resultSummaries.reduce(
    (sum, item) => sum + item.exp_total_annual_premium_excl_funeral,
    0
  )

  rows.push({
    category: 'Total',
    memberCount: totalLives.toString(),
    totalSalary: roundUpToTwoDecimalsAccounting(totalSalary),
    totalSumAssured: roundUpToTwoDecimalsAccounting(totalSumAssured),
    annualPremium: roundUpToTwoDecimalsAccounting(totalPremium),
    percentSalary: `${roundUpToTwoDecimalsAccounting((totalPremium / totalSalary) * 100)}%`
  })

  return rows
}

/**
 * Build rows for the Group Funeral summary table.
 * Extracted from QuoteOutput.vue lines 1204-1247.
 */
export function buildGroupFuneralRows(
  resultSummaries: any[]
): GroupFuneralRow[] {
  const rows: GroupFuneralRow[] = []

  resultSummaries.forEach((item) => {
    rows.push({
      category: item.category,
      memberCount: item.member_count.toString(),
      monthlyPremium: roundUpToTwoDecimalsAccounting(
        item.exp_total_fun_monthly_premium_per_member
      ),
      annualPremium: roundUpToTwoDecimalsAccounting(
        item.exp_total_fun_annual_premium_per_member
      ),
      totalAnnualPremium: roundUpToTwoDecimalsAccounting(
        computeOfficePremium(item.exp_total_fun_annual_risk_premium, item)
      )
    })
  })

  // Totals row
  rows.push({
    category: 'Total',
    memberCount: resultSummaries
      .reduce((sum, item) => sum + item.member_count, 0)
      .toString(),
    monthlyPremium: roundUpToTwoDecimalsAccounting(
      resultSummaries.reduce(
        (sum, item) => sum + item.exp_total_fun_monthly_premium_per_member,
        0
      )
    ),
    annualPremium: roundUpToTwoDecimalsAccounting(
      resultSummaries.reduce(
        (sum, item) => sum + item.exp_total_fun_annual_premium_per_member,
        0
      )
    ),
    totalAnnualPremium: roundUpToTwoDecimalsAccounting(
      resultSummaries.reduce(
        (sum, item) =>
          sum +
          computeOfficePremium(item.exp_total_fun_annual_risk_premium, item),
        0
      )
    )
  })

  return rows
}

// ---------------------------------------------------------------------------
// Table data builders — Section 5 (Premium Breakdown per category)
// ---------------------------------------------------------------------------

/**
 * Build the per-category benefit premium breakdown rows.
 * Extracted from QuoteOutput.vue lines 1315-1354.
 */
export function buildPremiumBreakdownRows(
  item: any,
  benefitTitles: BenefitTitles
): PremiumBreakdownRow[] {
  return [
    {
      benefit: benefitTitles.glaBenefitTitle,
      totalSumAssured: roundUpToTwoDecimalsAccounting(
        item.total_gla_capped_sum_assured
      ),
      annualPremium: roundUpToTwoDecimalsAccounting(
        computeOfficePremium(item.exp_total_gla_annual_risk_premium, item)
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_gla_annual_risk_premium_salary, item) * 100)}%`
    },
    {
      benefit: benefitTitles.sglaBenefitTitle,
      totalSumAssured: roundUpToTwoDecimalsAccounting(
        item.total_sgla_capped_sum_assured
      ),
      annualPremium: roundUpToTwoDecimalsAccounting(
        computeOfficePremium(item.exp_total_sgla_annual_risk_premium, item)
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_sgla_annual_risk_premium_salary, item) * 100)}%`
    },
    {
      benefit: benefitTitles.ptdBenefitTitle,
      totalSumAssured: roundUpToTwoDecimalsAccounting(
        item.total_ptd_capped_sum_assured
      ),
      annualPremium: roundUpToTwoDecimalsAccounting(
        computeOfficePremium(item.exp_total_ptd_annual_risk_premium, item)
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_ptd_annual_risk_premium_salary, item) * 100)}%`
    },
    {
      benefit: benefitTitles.ciBenefitTitle,
      totalSumAssured: roundUpToTwoDecimalsAccounting(
        item.total_ci_capped_sum_assured
      ),
      annualPremium: roundUpToTwoDecimalsAccounting(
        computeOfficePremium(item.exp_total_ci_annual_risk_premium, item)
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_ci_annual_risk_premium_salary, item) * 100)}%`
    },
    {
      benefit: benefitTitles.phiBenefitTitle,
      totalSumAssured: roundUpToTwoDecimalsAccounting(
        item.total_phi_capped_income
      ),
      annualPremium: roundUpToTwoDecimalsAccounting(
        computeOfficePremium(item.exp_total_phi_annual_risk_premium, item)
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_phi_annual_risk_premium_salary, item) * 100)}%`
    },
    {
      benefit: benefitTitles.ttdBenefitTitle,
      totalSumAssured: roundUpToTwoDecimalsAccounting(
        item.total_ttd_capped_income
      ),
      annualPremium: roundUpToTwoDecimalsAccounting(
        computeOfficePremium(item.exp_total_ttd_annual_risk_premium, item)
      ),
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(item.exp_proportion_ttd_annual_risk_premium_salary, item) * 100)}%`
    }
  ]
}

/**
 * Build the per-category group funeral breakdown rows (key-value pairs).
 * Extracted from QuoteOutput.vue lines 1408-1425.
 */
export function buildGroupFuneralBreakdownRows(item: any): LabelValueRow[] {
  return [
    {
      label: 'Monthly Premium per Member',
      value: roundUpToTwoDecimalsAccounting(
        item.exp_total_fun_monthly_premium_per_member
      )
    },
    {
      label: 'Annual Premium per Member',
      value: roundUpToTwoDecimalsAccounting(
        item.exp_total_fun_annual_premium_per_member
      )
    },
    {
      label: 'Total Annual Premium',
      value: roundUpToTwoDecimalsAccounting(
        computeOfficePremium(item.exp_total_fun_annual_risk_premium, item)
      )
    }
  ]
}

// ---------------------------------------------------------------------------
// Table data builders — Section 6 (Benefits and Definitions)
// ---------------------------------------------------------------------------

/**
 * Look up the customised display name for a benefit code (e.g. "GLA_EDU") from
 * the benefit-map list. Falls back to the provided default when no map or
 * matching row is present. Used to resolve split-educator names ("GLA
 * Educator" / "PTD Educator") which are user-customisable in the Benefits
 * Customisation screen.
 */
export function resolveBenefitAlias(
  benefitMaps: any[] | undefined | null,
  code: string,
  fallback: string
): string {
  if (!benefitMaps) return fallback
  const match = benefitMaps.find((b: any) => b?.benefit_code === code)
  return (
    (match?.benefit_alias && String(match.benefit_alias).trim()) ||
    match?.benefit_name ||
    fallback
  )
}

/**
 * Build the common category benefits key-value rows.
 * Extracted from QuoteOutput.vue lines 1510-1520.
 */
export function buildCategoryCommonBenefitRows(
  item: any,
  quote: any,
  benefitMaps?: any[]
): LabelValueRow[] {
  const glaEduLabel = resolveBenefitAlias(
    benefitMaps,
    'GLA_EDU',
    'GLA Educator'
  )
  const ptdEduLabel = resolveBenefitAlias(
    benefitMaps,
    'PTD_EDU',
    'PTD Educator'
  )
  return [
    {
      label: 'Terminal Illness',
      value: `${item.gla_terminal_illness_benefit}`
    },
    { label: 'Free Cover Limit', value: `${quote.free_cover_limit}` },
    { label: glaEduLabel, value: `${item.gla_educator_benefit}` },
    { label: ptdEduLabel, value: `${item.ptd_educator_benefit}` },
    {
      label: 'Retirement Premium Waiver',
      value: `${item.phi_premium_waiver || 'No'}`
    },
    {
      label: 'Medical Aid Premium Waiver',
      value: `${item.phi_medical_aid_premium_waiver || 'No'}`
    }
  ]
}

/**
 * Build the 7-column benefit definition rows.
 * Extracted from QuoteOutput.vue lines 1555-1627.
 */
export function buildBenefitDefinitionRows(
  item: any,
  quote: any,
  benefitTitles: BenefitTitles
): BenefitDefinitionRow[] {
  return [
    {
      benefit: benefitTitles.glaBenefitTitle,
      salaryMultiple: quote.use_global_salary_multiple
        ? `${item.gla_salary_multiple}`
        : 'varies',
      benefitStructure: 'standalone',
      waitingPeriod: `${item.gla_waiting_period}`,
      deferredPeriod: 'n.a',
      coverDefinition: 'n.a',
      riskType: 'n.a'
    },
    {
      benefit: benefitTitles.sglaBenefitTitle,
      salaryMultiple: quote.use_global_salary_multiple
        ? `${item.sgla_salary_multiple}`
        : 'varies',
      benefitStructure: 'rider',
      waitingPeriod: `${item.sgla_waiting_period}`,
      deferredPeriod: 'n.a',
      coverDefinition: 'n.a',
      riskType: 'n.a'
    },
    {
      benefit: benefitTitles.ptdBenefitTitle,
      salaryMultiple: quote.use_global_salary_multiple
        ? `${item.ptd_salary_multiple}`
        : 'varies',
      benefitStructure: `${item.ptd_benefit_type}`,
      waitingPeriod: '0',
      deferredPeriod: `${item.ptd_deferred_period}`,
      coverDefinition: `${item.ptd_disability_definition}`,
      riskType: `${item.ptd_risk_type}`
    },
    {
      benefit: benefitTitles.ciBenefitTitle,
      salaryMultiple: quote.use_global_salary_multiple
        ? `${item.ci_critical_illness_salary_multiple}`
        : 'varies',
      benefitStructure: `${item.ci_benefit_structure}`,
      waitingPeriod: `${item.ci_waiting_period}`,
      deferredPeriod: `${item.ci_deferred_period}`,
      coverDefinition: `${item.ci_benefit_definition}`,
      riskType: 'n.a'
    },
    {
      benefit: benefitTitles.phiBenefitTitle,
      salaryMultiple: `${item.phi_income_replacement_percentage / 100}`,
      benefitStructure: 'n.a',
      waitingPeriod: `${item.phi_waiting_period}`,
      deferredPeriod: `${item.phi_deferred_period}`,
      coverDefinition: `${item.phi_disability_definition}`,
      riskType: `${item.phi_risk_type}`
    },
    {
      benefit: benefitTitles.ttdBenefitTitle,
      salaryMultiple: `${item.ttd_income_replacement_percentage / 100}`,
      benefitStructure: 'n.a',
      waitingPeriod: `${item.ttd_waiting_period}`,
      deferredPeriod: `${item.ttd_deferred_period}`,
      coverDefinition: `${item.ttd_disability_definition}`,
      riskType: `${item.ttd_risk_type}`
    }
  ]
}

/**
 * Build the group funeral coverage rows for a given scheme category.
 * Extracted from QuoteOutput.vue lines 1687-1708.
 */
export function buildFuneralCoverageRows(item: any): FuneralCoverageRow[] {
  return [
    {
      member: 'Main Member',
      sumAssured: item.family_funeral_main_member_funeral_sum_assured,
      maxCovered: 1
    },
    {
      member: 'Spouse',
      sumAssured: item.family_funeral_spouse_funeral_sum_assured,
      maxCovered: 1
    },
    {
      member: 'Child',
      sumAssured: item.family_funeral_children_funeral_sum_assured,
      maxCovered: item.family_funeral_max_number_children
    },
    {
      member: 'Parent',
      sumAssured: item.family_funeral_parent_funeral_sum_assured,
      maxCovered: item.family_funeral_parent_maximum_number_covered
    },
    {
      member: 'Dependant',
      sumAssured: item.family_funeral_adult_dependant_sum_assured,
      maxCovered: item.family_funeral_max_number_adult_dependants
    }
  ]
}

/**
 * Build educator benefit rows for a given scheme category.
 * Extracted from QuoteOutput.vue lines 1774-1802.
 */
export function buildEducatorBenefitRows(
  categoryEducatorBenefits: any[],
  schemeCategory: string
): EducatorBenefitRow[] {
  const catItem = categoryEducatorBenefits.find(
    (eb) => eb.scheme_category === schemeCategory
  )

  if (!catItem) return []

  const s = catItem.educator_benefit_structure
  return [
    {
      level: 'Grade 0',
      maxTuition: s?.grade0_max_tuition_per_year || 'n.a',
      maxCoverage: s?.grade0_max_coverage_years || 'n.a'
    },
    {
      level: 'Grade 1 - 7',
      maxTuition: s?.grade17_max_tuition_per_year || 'n.a',
      maxCoverage: s?.grade17_max_coverage_years || 'n.a'
    },
    {
      level: 'Grade 8 - 12',
      maxTuition: s?.grade812_max_tuition_per_year || 'n.a',
      maxCoverage: s?.grade812_max_coverage_years || 'n.a'
    },
    {
      level: 'Tertiary Education',
      maxTuition: s?.tertiary_max_tuition_per_year || 'n.a',
      maxCoverage: s?.tertiary_max_coverage_years || 'n.a'
    }
  ]
}
