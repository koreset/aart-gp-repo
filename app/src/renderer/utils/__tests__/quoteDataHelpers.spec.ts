import { describe, it, expect } from 'vitest'
import {
  safeGetValue,
  roundUpToTwoDecimalsAccounting,
  dashIfEmpty,
  calculateQuoteTotals,
  hasAnyNonFuneralBenefits,
  categoryHasNonFuneralBenefits,
  buildInitialInfoRows,
  buildPremiumSummaryRows,
  buildGroupFuneralRows,
  buildPremiumBreakdownRows,
  buildGroupFuneralBreakdownRows,
  buildCategoryCommonBenefitRows,
  buildBenefitDefinitionRows,
  buildFuneralCoverageRows,
  buildEducatorBenefitRows,
  schemeTotalLoading,
  computeOfficePremium,
  officeRateFromRiskRate,
  officeProportionFromRiskProportion
} from '../quoteDataHelpers'
import type { BenefitTitles } from '@/renderer/types/docxQuote'

// ---------------------------------------------------------------------------
// Test fixtures
// ---------------------------------------------------------------------------

const mockBenefitTitles: BenefitTitles = {
  glaBenefitTitle: 'Group Life Assurance',
  sglaBenefitTitle: 'Spouse Group Life',
  ptdBenefitTitle: 'Permanent Total Disability',
  ciBenefitTitle: 'Critical Illness',
  phiBenefitTitle: 'Permanent Health Insurance',
  ttdBenefitTitle: 'Temporary Total Disability',
  familyFuneralBenefitTitle: 'Family Funeral'
}

function makeSummary(overrides: Record<string, any> = {}) {
  // Office premium is now derived from the risk premium and the scheme-level
  // loading. With expense + commission + profit = 0.20, office = risk / 0.8.
  // The values below are the equivalent risk premiums for the previous
  // office-premium fixtures (e.g. 96000 / 0.8 = 120000).
  return {
    category: 'Category A',
    member_count: 100,
    total_sum_assured: 50000000,
    total_annual_salary: 20000000,
    total_annual_premium: 300000,
    expense_loading: 0.05,
    commission_loading: 0.1,
    profit_loading: 0.05,
    total_gla_capped_sum_assured: 30000000,
    total_ptd_capped_sum_assured: 10000000,
    total_ci_capped_sum_assured: 5000000,
    total_sgla_capped_sum_assured: 3000000,
    total_phi_capped_income: 1000000,
    total_ttd_capped_income: 500000,
    exp_total_gla_annual_risk_premium: 96000,
    exp_total_ptd_annual_risk_premium: 40000,
    exp_total_ci_annual_risk_premium: 24000,
    exp_total_sgla_annual_risk_premium: 16000,
    exp_total_phi_annual_risk_premium: 12000,
    exp_total_ttd_annual_risk_premium: 8000,
    exp_total_annual_premium_excl_funeral: 245000,
    proportion_exp_total_premium_excl_funeral_salary: 0.01225,
    exp_proportion_gla_annual_risk_premium_salary: 0.0048,
    exp_proportion_ptd_annual_risk_premium_salary: 0.002,
    exp_proportion_ci_annual_risk_premium_salary: 0.0012,
    exp_proportion_sgla_annual_risk_premium_salary: 0.0008,
    exp_proportion_phi_annual_risk_premium_salary: 0.0006,
    exp_proportion_ttd_annual_risk_premium_salary: 0.0004,
    exp_total_fun_monthly_premium_per_member: 45.5,
    exp_total_fun_annual_premium_per_member: 546,
    exp_total_fun_annual_risk_premium: 43680,
    ...overrides
  }
}

function makeFuneralOnlySummary() {
  return makeSummary({
    total_gla_capped_sum_assured: 0,
    total_ptd_capped_sum_assured: 0,
    total_ci_capped_sum_assured: 0,
    total_sgla_capped_sum_assured: 0,
    total_phi_capped_income: 0,
    total_ttd_capped_income: 0
  })
}

const mockQuote = {
  quote_name: 'QTE-2026-00142',
  scheme_name: 'Acme Corp Pension Fund',
  creation_date: '2026-03-28',
  commencement_date: '2026-04-01',
  free_cover_limit: 2000000,
  use_global_salary_multiple: true
}

const mockSchemeCategory = {
  scheme_category: 'Category A',
  gla_terminal_illness_benefit: 'Yes',
  gla_educator_benefit: 'Yes',
  ptd_educator_benefit: 'No',
  phi_premium_waiver: 'Yes',
  phi_medical_aid_premium_waiver: 'No',
  gla_salary_multiple: 3,
  sgla_salary_multiple: 2,
  ptd_salary_multiple: 3,
  ci_critical_illness_salary_multiple: 2,
  gla_waiting_period: '6 months',
  sgla_waiting_period: '6 months',
  ptd_benefit_type: 'Lump Sum',
  ptd_deferred_period: '6 months',
  ptd_disability_definition: 'Own Occupation',
  ptd_risk_type: 'Level',
  ci_benefit_structure: 'Accelerated',
  ci_waiting_period: '3 months',
  ci_deferred_period: '90 days',
  ci_benefit_definition: '4 conditions',
  phi_income_replacement_percentage: 75,
  phi_waiting_period: '6 months',
  phi_deferred_period: '6 months',
  phi_disability_definition: 'Own Occupation',
  phi_risk_type: 'Level',
  ttd_income_replacement_percentage: 75,
  ttd_waiting_period: '0',
  ttd_deferred_period: '14 days',
  ttd_disability_definition: 'Own Occupation',
  ttd_risk_type: 'Level',
  family_funeral_main_member_funeral_sum_assured: 50000,
  family_funeral_spouse_funeral_sum_assured: 50000,
  family_funeral_children_funeral_sum_assured: 25000,
  family_funeral_max_number_children: 4,
  family_funeral_parent_funeral_sum_assured: 20000,
  family_funeral_parent_maximum_number_covered: 2,
  family_funeral_adult_dependant_sum_assured: 15000,
  family_funeral_max_number_adult_dependants: 2
}

// ---------------------------------------------------------------------------
// Utility function tests
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Office-premium derivation tests
// ---------------------------------------------------------------------------

describe('schemeTotalLoading', () => {
  it('returns zero when no loadings are set', () => {
    expect(schemeTotalLoading({})).toBe(0)
  })

  it('sums expense + profit + admin + other + binder + outsource', () => {
    expect(
      schemeTotalLoading({
        expense_loading: 0.05,
        profit_loading: 0.05,
        admin_loading: 0.01,
        other_loading: 0.01,
        binder_fee_rate: 0.02,
        outsource_fee_rate: 0.01
      })
    ).toBeCloseTo(0.15, 10)
  })

  it('ignores commission_loading even if present on the input shape', () => {
    expect(
      schemeTotalLoading({
        expense_loading: 0.05,
        profit_loading: 0.05,
        // @ts-expect-error commission_loading is intentionally not in the type
        commission_loading: 0.1
      })
    ).toBeCloseTo(0.1, 10)
  })

  it('treats non-finite values as zero', () => {
    expect(
      schemeTotalLoading({
        expense_loading: NaN,
        profit_loading: 0.05,
        admin_loading: Infinity
      })
    ).toBeCloseTo(0.05, 10)
  })
})

describe('computeOfficePremium', () => {
  it('returns risk premium unchanged when all loadings are zero', () => {
    expect(computeOfficePremium(100, {})).toBe(100)
  })

  it('grosses up by the full denominator (binder + outsource included)', () => {
    expect(
      computeOfficePremium(85, {
        expense_loading: 0.05,
        profit_loading: 0.05,
        admin_loading: 0.01,
        other_loading: 0.01,
        binder_fee_rate: 0.02,
        outsource_fee_rate: 0.01
      })
    ).toBeCloseTo(100, 10) // 85 / (1 - 0.15)
  })

  it('returns zero when denominator collapses to <= 0', () => {
    expect(
      computeOfficePremium(100, {
        expense_loading: 0.6,
        profit_loading: 0.5
      })
    ).toBe(0)
  })
})

describe('officeRateFromRiskRate', () => {
  it('scales risk-rate-per-1000 by 1 / (1 - schemeTotalLoading)', () => {
    expect(
      officeRateFromRiskRate(8.7, {
        expense_loading: 0.05,
        profit_loading: 0.05,
        binder_fee_rate: 0.02,
        outsource_fee_rate: 0.01
      })
    ).toBeCloseTo(10, 10) // 8.7 / (1 - 0.13)
  })
})

describe('officeProportionFromRiskProportion', () => {
  it('scales risk-proportion by 1 / (1 - schemeTotalLoading)', () => {
    expect(
      officeProportionFromRiskProportion(0.0087, {
        expense_loading: 0.05,
        profit_loading: 0.05,
        binder_fee_rate: 0.02,
        outsource_fee_rate: 0.01
      })
    ).toBeCloseTo(0.01, 10) // 0.0087 / (1 - 0.13)
  })
})

describe('safeGetValue', () => {
  it('returns nested value from a valid path', () => {
    const obj = { a: { b: { c: 42 } } }
    expect(safeGetValue(obj, 'a.b.c')).toBe(42)
  })

  it('returns defaultValue when path does not exist', () => {
    expect(safeGetValue({}, 'a.b.c', 0)).toBe(0)
  })

  it('returns defaultValue for null object', () => {
    expect(safeGetValue(null, 'a', 'fallback')).toBe('fallback')
  })

  it('returns null as default when no defaultValue specified', () => {
    expect(safeGetValue({}, 'missing')).toBeNull()
  })

  it('returns value at top level', () => {
    expect(safeGetValue({ x: 'hello' }, 'x')).toBe('hello')
  })
})

describe('roundUpToTwoDecimalsAccounting', () => {
  it('rounds up and formats with spaces', () => {
    expect(roundUpToTwoDecimalsAccounting(1234567.891)).toBe('1 234 567.90')
  })

  it('handles zero', () => {
    expect(roundUpToTwoDecimalsAccounting(0)).toBe('0.00')
  })

  it('rounds up 0.001 correctly', () => {
    expect(roundUpToTwoDecimalsAccounting(0.001)).toBe('0.01')
  })

  it('handles exact two-decimal values', () => {
    expect(roundUpToTwoDecimalsAccounting(99.99)).toBe('99.99')
  })

  it('formats large numbers with space separators', () => {
    expect(roundUpToTwoDecimalsAccounting(845320000)).toBe('845 320 000.00')
  })
})

describe('dashIfEmpty', () => {
  it('returns value when truthy', () => {
    expect(dashIfEmpty('hello')).toBe('hello')
  })

  it('returns dash for empty string', () => {
    expect(dashIfEmpty('')).toBe('-')
  })

  it('returns dash for null', () => {
    expect(dashIfEmpty(null)).toBe('-')
  })

  it('returns dash for undefined', () => {
    expect(dashIfEmpty(undefined)).toBe('-')
  })

  it('returns dash for zero', () => {
    expect(dashIfEmpty(0)).toBe('-')
  })
})

// ---------------------------------------------------------------------------
// Aggregation helper tests
// ---------------------------------------------------------------------------

describe('calculateQuoteTotals', () => {
  it('sums totals across multiple summaries', () => {
    const summaries = [
      makeSummary({
        member_count: 100,
        total_sum_assured: 5000000,
        total_annual_salary: 2000000,
        total_annual_premium: 100000
      }),
      makeSummary({
        member_count: 200,
        total_sum_assured: 10000000,
        total_annual_salary: 4000000,
        total_annual_premium: 200000
      })
    ]
    const totals = calculateQuoteTotals(summaries)
    expect(totals.totalLives).toBe(300)
    expect(totals.totalSumAssured).toBe(15000000)
    expect(totals.totalAnnualSalary).toBe(6000000)
    expect(totals.totalAnnualPremium).toBe(300000)
  })

  it('returns zeros for empty array', () => {
    const totals = calculateQuoteTotals([])
    expect(totals.totalLives).toBe(0)
    expect(totals.totalSumAssured).toBe(0)
    expect(totals.totalAnnualSalary).toBe(0)
    expect(totals.totalAnnualPremium).toBe(0)
  })

  it('handles single summary', () => {
    const totals = calculateQuoteTotals([makeSummary()])
    expect(totals.totalLives).toBe(100)
    expect(totals.totalSumAssured).toBe(50000000)
  })
})

describe('hasAnyNonFuneralBenefits', () => {
  it('returns true when GLA sum assured > 0', () => {
    expect(hasAnyNonFuneralBenefits([makeSummary()])).toBe(true)
  })

  it('returns false when all non-funeral benefits are zero', () => {
    expect(hasAnyNonFuneralBenefits([makeFuneralOnlySummary()])).toBe(false)
  })

  it('returns true if any single category has non-funeral benefits', () => {
    expect(
      hasAnyNonFuneralBenefits([makeFuneralOnlySummary(), makeSummary()])
    ).toBe(true)
  })

  it('returns false for empty array', () => {
    expect(hasAnyNonFuneralBenefits([])).toBe(false)
  })
})

describe('categoryHasNonFuneralBenefits', () => {
  it('returns true for a category with benefits', () => {
    expect(categoryHasNonFuneralBenefits(makeSummary())).toBe(true)
  })

  it('returns false for funeral-only category', () => {
    expect(categoryHasNonFuneralBenefits(makeFuneralOnlySummary())).toBe(false)
  })
})

// ---------------------------------------------------------------------------
// Table data builder tests
// ---------------------------------------------------------------------------

describe('buildInitialInfoRows', () => {
  it('returns 9 key-value rows', () => {
    const totals = calculateQuoteTotals([makeSummary()])
    const rows = buildInitialInfoRows(mockQuote, totals)
    expect(rows).toHaveLength(9)
  })

  it('includes quote name and scheme name', () => {
    const totals = calculateQuoteTotals([makeSummary()])
    const rows = buildInitialInfoRows(mockQuote, totals)
    const labels = rows.map((r) => r.label)
    expect(labels).toContain('Quote Number:')
    expect(labels).toContain('Scheme Name:')

    const quoteRow = rows.find((r) => r.label === 'Quote Number:')
    expect(quoteRow?.value).toBe('QTE-2026-00142')
  })

  it('formats total sum assured as accounting number', () => {
    const totals = calculateQuoteTotals([makeSummary()])
    const rows = buildInitialInfoRows(mockQuote, totals)
    const sumAssuredRow = rows.find((r) => r.label === 'Total Sum Assured:')
    expect(sumAssuredRow?.value).toBe('50 000 000.00')
  })
})

describe('buildPremiumSummaryRows', () => {
  it('includes one data row per category plus a totals row', () => {
    const rows = buildPremiumSummaryRows([makeSummary()])
    expect(rows).toHaveLength(2) // 1 data + 1 total
    expect(rows[rows.length - 1].category).toBe('Total')
  })

  it('skips funeral-only categories', () => {
    const rows = buildPremiumSummaryRows([
      makeFuneralOnlySummary(),
      makeSummary()
    ])
    // Only the non-funeral category + total
    expect(rows).toHaveLength(2)
  })

  it('totals row sums correctly across categories', () => {
    const summaries = [
      makeSummary({ category: 'A', member_count: 50 }),
      makeSummary({ category: 'B', member_count: 75 })
    ]
    const rows = buildPremiumSummaryRows(summaries)
    const totalRow = rows.find((r) => r.category === 'Total')
    expect(totalRow?.memberCount).toBe('125')
  })
})

describe('buildGroupFuneralRows', () => {
  it('includes one row per category plus totals', () => {
    const rows = buildGroupFuneralRows([
      makeSummary(),
      makeSummary({ category: 'B' })
    ])
    expect(rows).toHaveLength(3) // 2 data + 1 total
    expect(rows[rows.length - 1].category).toBe('Total')
  })

  it('formats monthly premium correctly', () => {
    const rows = buildGroupFuneralRows([makeSummary()])
    expect(rows[0].monthlyPremium).toBe('45.50')
  })
})

describe('buildPremiumBreakdownRows', () => {
  it('returns 6 benefit rows', () => {
    const rows = buildPremiumBreakdownRows(makeSummary(), mockBenefitTitles)
    expect(rows).toHaveLength(6)
  })

  it('uses benefit titles from the maps', () => {
    const rows = buildPremiumBreakdownRows(makeSummary(), mockBenefitTitles)
    expect(rows[0].benefit).toBe('Group Life Assurance')
    expect(rows[1].benefit).toBe('Spouse Group Life')
  })

  it('formats percentages with % suffix', () => {
    const rows = buildPremiumBreakdownRows(makeSummary(), mockBenefitTitles)
    rows.forEach((row) => {
      expect(row.percentSalary).toMatch(/%$/)
    })
  })
})

describe('buildGroupFuneralBreakdownRows', () => {
  it('returns 3 key-value rows', () => {
    const rows = buildGroupFuneralBreakdownRows(makeSummary())
    expect(rows).toHaveLength(3)
    expect(rows[0].label).toBe('Monthly Premium per Member')
    expect(rows[1].label).toBe('Annual Premium per Member')
    expect(rows[2].label).toBe('Total Annual Premium')
  })
})

describe('buildCategoryCommonBenefitRows', () => {
  it('returns 6 key-value rows', () => {
    const rows = buildCategoryCommonBenefitRows(mockSchemeCategory, mockQuote)
    expect(rows).toHaveLength(6)
  })

  it('includes terminal illness and free cover limit', () => {
    const rows = buildCategoryCommonBenefitRows(mockSchemeCategory, mockQuote)
    const labels = rows.map((r) => r.label)
    expect(labels).toContain('Terminal Illness')
    expect(labels).toContain('Free Cover Limit')
  })
})

describe('buildBenefitDefinitionRows', () => {
  it('returns 6 rows (one per benefit type)', () => {
    const rows = buildBenefitDefinitionRows(
      mockSchemeCategory,
      mockQuote,
      mockBenefitTitles
    )
    expect(rows).toHaveLength(6)
  })

  it('uses salary multiples when global flag is set', () => {
    const rows = buildBenefitDefinitionRows(
      mockSchemeCategory,
      mockQuote,
      mockBenefitTitles
    )
    expect(rows[0].salaryMultiple).toBe('3') // GLA
    expect(rows[1].salaryMultiple).toBe('2') // SGLA
  })

  it('shows "varies" when global salary multiple is off', () => {
    const quoteNoGlobal = { ...mockQuote, use_global_salary_multiple: false }
    const rows = buildBenefitDefinitionRows(
      mockSchemeCategory,
      quoteNoGlobal,
      mockBenefitTitles
    )
    expect(rows[0].salaryMultiple).toBe('varies')
    expect(rows[1].salaryMultiple).toBe('varies')
    // PHI and TTD use income replacement, not salary multiple
    expect(rows[4].salaryMultiple).toBe('0.75')
  })

  it('populates all 7 columns per row', () => {
    const rows = buildBenefitDefinitionRows(
      mockSchemeCategory,
      mockQuote,
      mockBenefitTitles
    )
    rows.forEach((row) => {
      expect(row).toHaveProperty('benefit')
      expect(row).toHaveProperty('salaryMultiple')
      expect(row).toHaveProperty('benefitStructure')
      expect(row).toHaveProperty('waitingPeriod')
      expect(row).toHaveProperty('deferredPeriod')
      expect(row).toHaveProperty('coverDefinition')
      expect(row).toHaveProperty('riskType')
    })
  })
})

describe('buildFuneralCoverageRows', () => {
  it('returns 5 member type rows', () => {
    const rows = buildFuneralCoverageRows(mockSchemeCategory)
    expect(rows).toHaveLength(5)
    expect(rows.map((r) => r.member)).toEqual([
      'Main Member',
      'Spouse',
      'Child',
      'Parent',
      'Dependant'
    ])
  })

  it('sets main member and spouse max covered to 1', () => {
    const rows = buildFuneralCoverageRows(mockSchemeCategory)
    expect(rows[0].maxCovered).toBe(1)
    expect(rows[1].maxCovered).toBe(1)
  })

  it('uses scheme category values for sum assured', () => {
    const rows = buildFuneralCoverageRows(mockSchemeCategory)
    expect(rows[0].sumAssured).toBe(50000)
    expect(rows[2].sumAssured).toBe(25000) // child
  })
})

describe('buildEducatorBenefitRows', () => {
  const mockEduBenefits = [
    {
      scheme_category: 'Category A',
      educator_benefit_structure: {
        grade0_max_tuition_per_year: 30000,
        grade0_max_coverage_years: 1,
        grade17_max_tuition_per_year: 50000,
        grade17_max_coverage_years: 7,
        grade812_max_tuition_per_year: 80000,
        grade812_max_coverage_years: 5,
        tertiary_max_tuition_per_year: 120000,
        tertiary_max_coverage_years: 4
      }
    }
  ]

  it('returns 4 education level rows when category matches', () => {
    const rows = buildEducatorBenefitRows(mockEduBenefits, 'Category A')
    expect(rows).toHaveLength(4)
    expect(rows[0].level).toBe('Grade 0')
    expect(rows[3].level).toBe('Tertiary Education')
  })

  it('returns empty array when category does not match', () => {
    const rows = buildEducatorBenefitRows(mockEduBenefits, 'Unknown Category')
    expect(rows).toHaveLength(0)
  })

  it('returns empty array for empty benefits list', () => {
    const rows = buildEducatorBenefitRows([], 'Category A')
    expect(rows).toHaveLength(0)
  })

  it('populates tuition and coverage values', () => {
    const rows = buildEducatorBenefitRows(mockEduBenefits, 'Category A')
    expect(rows[0].maxTuition).toBe(30000)
    expect(rows[0].maxCoverage).toBe(1)
    expect(rows[2].maxTuition).toBe(80000)
  })
})
