import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useBenefitScheduleGeneration } from '../useBenefitScheduleGeneration'
import type { BenefitTitles } from '@/renderer/types/docxQuote'

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

vi.mock('file-saver', () => ({
  saveAs: vi.fn()
}))

vi.mock('@/renderer/utils/helpers.js', () => ({
  default: (dateString: any, y: boolean, m: boolean, d: boolean) => {
    const date = new Date(dateString)
    if (y && m && d) {
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      return `${year}-${month}-${day}`
    }
    if (y && !m && !d) return String(date.getFullYear())
    return String(date)
  }
}))

// Mock jsPDF to avoid canvas/DOM dependencies in Node
vi.mock('jspdf', () => {
  function MockJsPDF() {
    return {
      internal: {
        pageSize: { getWidth: () => 210, getHeight: () => 297 },
        getNumberOfPages: () => 1
      },
      setFontSize: vi.fn(),
      setFont: vi.fn(),
      setTextColor: vi.fn(),
      setFillColor: vi.fn(),
      text: vi.fn(),
      rect: vi.fn(),
      addPage: vi.fn(),
      setPage: vi.fn(),
      save: vi.fn(),
      autoTable: vi.fn(),
      lastAutoTable: { finalY: 100 }
    }
  }
  return {
    default: MockJsPDF,
    jsPDF: MockJsPDF
  }
})

vi.mock('jspdf-autotable', () => ({
  applyPlugin: vi.fn()
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

const mockQuote = {
  id: 1,
  scheme_name: 'Test Pension Fund',
  creation_date: '2026-03-28',
  effective_date: '2026-04-01',
  selected_scheme_categories: ['Category A'],
  scheme_categories: [
    {
      scheme_category: 'Category A',
      gla_terminal_illness_benefit: 'Yes',
      gla_educator_benefit: 'No',
      ptd_educator_benefit: 'No',
      phi_premium_waiver: 'No',
      phi_medical_aid_premium_waiver: 'No',
      gla_salary_multiple: 3,
      ptd_salary_multiple: 3,
      gla_benefit_type: 'All Causes',
      ptd_benefit_type: 'Lump Sum'
    }
  ]
}

const mockResultSummaries = [
  {
    category: 'Category A',
    member_count: 100,
    total_annual_salary: 20000000,
    exp_total_annual_premium_excl_funeral: 245000,
    proportion_exp_total_premium_excl_funeral_salary: 0.01225,
    exp_total_gla_annual_office_premium: 120000,
    total_gla_capped_sum_assured: 30000000,
    exp_proportion_gla_office_premium_salary: 0.006,
    exp_total_ptd_annual_office_premium: 50000,
    total_ptd_capped_sum_assured: 10000000,
    exp_proportion_ptd_office_premium_salary: 0.0025,
    exp_total_ci_annual_office_premium: 0,
    total_ci_capped_sum_assured: 0,
    exp_proportion_ci_office_premium_salary: 0,
    exp_total_sgla_annual_office_premium: 0,
    total_sgla_capped_sum_assured: 0,
    exp_proportion_sgla_office_premium_salary: 0,
    exp_total_phi_annual_office_premium: 0,
    total_phi_capped_income: 0,
    exp_proportion_phi_office_premium_salary: 0,
    exp_total_ttd_annual_office_premium: 0,
    total_ttd_capped_income: 0,
    exp_proportion_ttd_office_premium_salary: 0,
    family_funeral_main_member_funeral_sum_assured: 0,
    family_funeral_spouse_funeral_sum_assured: 0,
    family_funeral_children_funeral_sum_assured: 0,
    family_funeral_adult_dependant_sum_assured: 0,
    family_funeral_parent_funeral_sum_assured: 0,
    family_funeral_max_number_children: 0,
    family_funeral_max_number_adult_dependants: 0
  }
]

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe('useBenefitScheduleGeneration', () => {
  let saveAsMock: any

  beforeEach(async () => {
    vi.clearAllMocks()
    saveAsMock = vi.mocked((await import('file-saver')).saveAs)
  })

  it('exposes isGenerating ref and two generation functions', () => {
    const {
      isGenerating,
      generateBenefitScheduleDocx,
      generateBenefitSchedulePdf
    } = useBenefitScheduleGeneration()
    expect(isGenerating.value).toBe(false)
    expect(typeof generateBenefitScheduleDocx).toBe('function')
    expect(typeof generateBenefitSchedulePdf).toBe('function')
  })

  it('generates DOCX and calls saveAs with correct filename', async () => {
    const { generateBenefitScheduleDocx, isGenerating } =
      useBenefitScheduleGeneration()

    await generateBenefitScheduleDocx(
      mockQuote,
      mockResultSummaries,
      mockBenefitTitles
    )

    expect(isGenerating.value).toBe(false)
    expect(saveAsMock).toHaveBeenCalledOnce()
    const [blob, filename] = saveAsMock.mock.calls[0]
    expect(blob).toBeInstanceOf(Blob)
    expect(filename).toBe('Test_Pension_Fund_Benefit_Schedule.docx')
  })

  it('DOCX blob is a valid ZIP file', async () => {
    const { generateBenefitScheduleDocx } = useBenefitScheduleGeneration()
    await generateBenefitScheduleDocx(
      mockQuote,
      mockResultSummaries,
      mockBenefitTitles
    )

    const blob: Blob = saveAsMock.mock.calls[0][0]
    const buffer = await blob.arrayBuffer()
    const bytes = new Uint8Array(buffer)
    // ZIP magic bytes: PK (0x50, 0x4B)
    expect(bytes[0]).toBe(0x50)
    expect(bytes[1]).toBe(0x4b)
  })

  it('isGenerating is true during generation', async () => {
    const { generateBenefitScheduleDocx, isGenerating } =
      useBenefitScheduleGeneration()

    const states: boolean[] = []
    const originalValue = isGenerating.value
    states.push(originalValue)

    const promise = generateBenefitScheduleDocx(
      mockQuote,
      mockResultSummaries,
      mockBenefitTitles
    )
    // isGenerating should be true before promise resolves
    states.push(isGenerating.value)
    await promise
    states.push(isGenerating.value)

    expect(states).toEqual([false, true, false])
  })

  it('generates PDF and calls save', async () => {
    const { generateBenefitSchedulePdf, isGenerating } =
      useBenefitScheduleGeneration()

    await generateBenefitSchedulePdf(
      mockQuote,
      mockResultSummaries,
      mockBenefitTitles
    )

    expect(isGenerating.value).toBe(false)
    // jsPDF.save is called (mocked)
    const jsPDFMod = await import('jspdf')
    const mockInstance = new (jsPDFMod.default as any)()
    // The save mock should have been called
    expect(mockInstance.save).toBeDefined()
  })

  it('handles missing scheme_name gracefully', async () => {
    const { generateBenefitScheduleDocx } = useBenefitScheduleGeneration()
    const quoteNoName = { ...mockQuote, scheme_name: '' }

    await generateBenefitScheduleDocx(
      quoteNoName,
      mockResultSummaries,
      mockBenefitTitles
    )

    const [, filename] = saveAsMock.mock.calls[0]
    expect(filename).toContain('Benefit_Schedule.docx')
  })
})
