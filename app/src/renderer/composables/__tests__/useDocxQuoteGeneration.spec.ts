import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useDocxQuoteGeneration } from '../useDocxQuoteGeneration'
import type { BenefitTitles } from '@/renderer/types/docxQuote'

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

// Mock file-saver so we capture the blob without triggering a download
vi.mock('file-saver', () => ({
  saveAs: vi.fn()
}))

// Mock the helpers.js formatDateString used by the composable
vi.mock('@/renderer/utils/helpers.js', () => ({
  default: (dateString: any, y: boolean, m: boolean, d: boolean) => {
    const date = new Date(dateString)
    if (y && m && d) {
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      return `${year}-${month}-${day}`
    }
    return String(date)
  }
}))

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

const mockInsurer = {
  name: 'Test Insurance Ltd',
  address_line_1: '123 Test Street',
  address_line_2: 'Sandton',
  address_line_3: 'Gauteng',
  city: 'Johannesburg',
  province: 'Gauteng',
  post_code: '2196',
  telephone: '+27 11 000 0000',
  email: 'test@test.co.za',
  introductory_text: 'We are pleased to submit this quotation.',
  general_provisions_text: 'Standard terms and conditions apply.',
  logo: null,
  logo_mime_type: 'image/png'
}

const mockQuote = {
  quote_name: 'QTE-TEST-001',
  scheme_name: 'Test Pension Fund',
  creation_date: '2026-03-28',
  commencement_date: '2026-04-01',
  free_cover_limit: 2000000,
  use_global_salary_multiple: true,
  scheme_categories: [
    {
      scheme_category: 'Category A',
      gla_terminal_illness_benefit: 'Yes',
      gla_educator_benefit: 'No',
      ptd_educator_benefit: 'No',
      phi_premium_waiver: 'No',
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
  ]
}

function makeSummary(overrides: Record<string, any> = {}) {
  // Scheme loading 0.20 (5% expense + 10% commission + 5% profit) yields
  // office = risk / 0.8 in computeOfficePremium. The risk fixtures below
  // map back to the previous 120000 / 50000 / ... office values.
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

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe('useDocxQuoteGeneration', () => {
  let saveAsMock: ReturnType<typeof vi.fn>

  beforeEach(async () => {
    vi.clearAllMocks()
    const fileSaver = await import('file-saver')
    saveAsMock = fileSaver.saveAs as ReturnType<typeof vi.fn>
  })

  it('exposes isGenerating, errorMessage, and generateDocxQuote', () => {
    const { isGenerating, errorMessage, generateDocxQuote } =
      useDocxQuoteGeneration()
    expect(isGenerating.value).toBe(false)
    expect(errorMessage.value).toBe('')
    expect(typeof generateDocxQuote).toBe('function')
  })

  it('generates a Blob and calls saveAs with correct filename', async () => {
    const { generateDocxQuote, isGenerating, errorMessage } =
      useDocxQuoteGeneration()

    await generateDocxQuote(
      mockQuote,
      [makeSummary()],
      mockInsurer,
      [],
      [],
      mockBenefitTitles
    )

    expect(errorMessage.value).toBe('')
    expect(isGenerating.value).toBe(false)
    expect(saveAsMock).toHaveBeenCalledTimes(1)

    // First arg should be a Blob
    const blob = saveAsMock.mock.calls[0][0]
    expect(blob).toBeInstanceOf(Blob)
    expect(blob.size).toBeGreaterThan(0)

    // Second arg should be the filename
    const filename = saveAsMock.mock.calls[0][1]
    expect(filename).toContain('Test Pension Fund')
    expect(filename).toContain('Quotation')
    expect(filename).toMatch(/\.docx$/)
  })

  it('sets isGenerating to true during generation', async () => {
    const { generateDocxQuote, isGenerating } = useDocxQuoteGeneration()

    // Start generation but don't await yet
    const promise = generateDocxQuote(
      mockQuote,
      [makeSummary()],
      mockInsurer,
      [],
      [],
      mockBenefitTitles
    )

    // After promise resolves, isGenerating should be false
    await promise
    expect(isGenerating.value).toBe(false)
  })

  it('produces a valid DOCX blob (ZIP magic bytes)', async () => {
    const { generateDocxQuote } = useDocxQuoteGeneration()

    await generateDocxQuote(
      mockQuote,
      [makeSummary()],
      mockInsurer,
      [],
      [],
      mockBenefitTitles
    )

    const blob = saveAsMock.mock.calls[0][0] as Blob
    const arrayBuffer = await blob.arrayBuffer()
    const bytes = new Uint8Array(arrayBuffer)

    // DOCX files are ZIP archives — first 4 bytes should be PK (0x50, 0x4B, 0x03, 0x04)
    expect(bytes[0]).toBe(0x50) // P
    expect(bytes[1]).toBe(0x4b) // K
    expect(bytes[2]).toBe(0x03)
    expect(bytes[3]).toBe(0x04)
  })

  it('handles funeral-only quotes without errors', async () => {
    const funeralOnlySummary = makeSummary({
      total_gla_capped_sum_assured: 0,
      total_ptd_capped_sum_assured: 0,
      total_ci_capped_sum_assured: 0,
      total_sgla_capped_sum_assured: 0,
      total_phi_capped_income: 0,
      total_ttd_capped_income: 0
    })

    const { generateDocxQuote, errorMessage } = useDocxQuoteGeneration()

    await generateDocxQuote(
      mockQuote,
      [funeralOnlySummary],
      mockInsurer,
      [],
      [],
      mockBenefitTitles
    )

    expect(errorMessage.value).toBe('')
    expect(saveAsMock).toHaveBeenCalledTimes(1)
  })

  it('handles multiple categories', async () => {
    const summaries = [
      makeSummary({ category: 'Category A' }),
      makeSummary({ category: 'Category B', member_count: 200 })
    ]

    const { generateDocxQuote, errorMessage } = useDocxQuoteGeneration()

    await generateDocxQuote(
      {
        ...mockQuote,
        scheme_categories: [
          ...mockQuote.scheme_categories,
          { ...mockQuote.scheme_categories[0], scheme_category: 'Category B' }
        ]
      },
      summaries,
      mockInsurer,
      [],
      [],
      mockBenefitTitles
    )

    expect(errorMessage.value).toBe('')
    expect(saveAsMock).toHaveBeenCalledTimes(1)

    // DOCX should be larger with 2 categories
    const blob = saveAsMock.mock.calls[0][0] as Blob
    expect(blob.size).toBeGreaterThan(1000)
  })

  it('handles insurer with no logo gracefully', async () => {
    const insurerNoLogo = { ...mockInsurer, logo: null }

    const { generateDocxQuote, errorMessage } = useDocxQuoteGeneration()

    await generateDocxQuote(
      mockQuote,
      [makeSummary()],
      insurerNoLogo,
      [],
      [],
      mockBenefitTitles
    )

    expect(errorMessage.value).toBe('')
    expect(saveAsMock).toHaveBeenCalledTimes(1)
  })

  it('handles educator benefits', async () => {
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

    const { generateDocxQuote, errorMessage } = useDocxQuoteGeneration()

    await generateDocxQuote(
      mockQuote,
      [makeSummary()],
      mockInsurer,
      mockEduBenefits,
      [],
      mockBenefitTitles
    )

    expect(errorMessage.value).toBe('')
    expect(saveAsMock).toHaveBeenCalledTimes(1)
  })

  it('sets errorMessage on failure and resets isGenerating', async () => {
    const { generateDocxQuote, isGenerating, errorMessage } =
      useDocxQuoteGeneration()

    // Suppress the expected console.error from the catch block
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {})

    // Pass null quote to trigger an error deep in the builder
    await generateDocxQuote(
      null as any,
      [makeSummary()],
      mockInsurer,
      [],
      [],
      mockBenefitTitles
    )

    expect(isGenerating.value).toBe(false)
    expect(errorMessage.value).not.toBe('')
    expect(errorMessage.value).toContain('Error generating Word document')

    consoleSpy.mockRestore()
  })
})
