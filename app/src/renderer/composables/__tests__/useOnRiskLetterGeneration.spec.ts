import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useOnRiskLetterGeneration } from '../useOnRiskLetterGeneration'

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
      output: vi.fn(() => new Blob(['fake-pdf'], { type: 'application/pdf' })),
      autoTable: vi.fn(),
      lastAutoTable: { finalY: 100 },
      splitTextToSize: vi.fn((text: string) => [text])
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

function makeTestData(overrides: Record<string, any> = {}) {
  return {
    quote: {
      id: 42,
      scheme_name: 'ABC Retirement Fund',
      scheme_contact: 'John Doe',
      scheme_email: 'john@abc.co.za',
      commencement_date: '2026-04-01',
      cover_end_date: '2027-03-31',
      industry: 'Mining',
      obligation_type: 'Compulsory',
      member_data_count: 250,
      currency: 'ZAR',
      free_cover_limit: 2000000,
      normal_retirement_age: 65,
      distribution_channel: 'broker',
      quote_broker: { broker_name: 'Test Broker' },
      ...overrides
    },
    scheme: {
      contact_person: 'John Doe',
      contact_email: 'john@abc.co.za'
    },
    insurer: {
      name: 'Test Insurer Ltd',
      address_line_1: '123 Main St',
      city: 'Johannesburg',
      province: 'Gauteng',
      post_code: '2000',
      telephone: '+27 11 000 0000',
      email: 'info@insurer.co.za',
      contact_person: 'Jane Smith',
      logo: null,
      logo_mime_type: 'image/png',
      on_risk_letter_text: null
    },
    letter: {
      letter_date: '2026-04-01',
      letter_reference: 'ORL-2026-001'
    },
    benefit_summary: [
      { benefit: 'Group Life Assurance', annual_premium: 120000 },
      { benefit: 'Permanent Total Disability', annual_premium: 50000 }
    ]
  }
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe('useOnRiskLetterGeneration', () => {
  let saveAsMock: any

  beforeEach(async () => {
    vi.clearAllMocks()
    saveAsMock = vi.mocked((await import('file-saver')).saveAs)
  })

  it('exposes isGenerating, errorMessage, and two generation functions', () => {
    const {
      isGenerating,
      errorMessage,
      generateOnRiskLetterDocx,
      generateOnRiskLetterPdf
    } = useOnRiskLetterGeneration()
    expect(isGenerating.value).toBe(false)
    expect(errorMessage.value).toBe('')
    expect(typeof generateOnRiskLetterDocx).toBe('function')
    expect(typeof generateOnRiskLetterPdf).toBe('function')
  })

  it('generates DOCX and calls saveAs', async () => {
    const { generateOnRiskLetterDocx, isGenerating } =
      useOnRiskLetterGeneration()

    await generateOnRiskLetterDocx(makeTestData())

    expect(isGenerating.value).toBe(false)
    expect(saveAsMock).toHaveBeenCalledOnce()
    const [blob, filename] = saveAsMock.mock.calls[0]
    expect(blob).toBeInstanceOf(Blob)
    expect(filename).toContain('ABC Retirement Fund')
    expect(filename).toContain('On_Risk_Letter')
    expect(filename).toMatch(/\.docx$/)
  })

  it('DOCX blob has ZIP magic bytes', async () => {
    const { generateOnRiskLetterDocx } = useOnRiskLetterGeneration()
    await generateOnRiskLetterDocx(makeTestData())

    const blob: Blob = saveAsMock.mock.calls[0][0]
    const buffer = await blob.arrayBuffer()
    const bytes = new Uint8Array(buffer)
    expect(bytes[0]).toBe(0x50)
    expect(bytes[1]).toBe(0x4b)
  })

  it('isGenerating is true during DOCX generation', async () => {
    const { generateOnRiskLetterDocx, isGenerating } =
      useOnRiskLetterGeneration()

    const states: boolean[] = [isGenerating.value]
    const promise = generateOnRiskLetterDocx(makeTestData())
    states.push(isGenerating.value)
    await promise
    states.push(isGenerating.value)

    expect(states).toEqual([false, true, false])
  })

  it('generates PDF and calls saveAs', async () => {
    const { generateOnRiskLetterPdf, isGenerating } =
      useOnRiskLetterGeneration()

    await generateOnRiskLetterPdf(makeTestData())

    expect(isGenerating.value).toBe(false)
    expect(saveAsMock).toHaveBeenCalledOnce()
    const [, filename] = saveAsMock.mock.calls[0]
    expect(filename).toContain('On_Risk_Letter')
    expect(filename).toMatch(/\.pdf$/)
  })

  it('handles data without benefit_summary', async () => {
    const { generateOnRiskLetterDocx, errorMessage } =
      useOnRiskLetterGeneration()
    const data = makeTestData()
    data.benefit_summary = []

    await generateOnRiskLetterDocx(data)
    expect(errorMessage.value).toBe('')
    expect(saveAsMock).toHaveBeenCalledOnce()
  })

  it('handles direct distribution channel (no broker)', async () => {
    const { generateOnRiskLetterDocx } = useOnRiskLetterGeneration()
    const data = makeTestData({
      distribution_channel: 'direct',
      quote_broker: null
    })

    await generateOnRiskLetterDocx(data)
    expect(saveAsMock).toHaveBeenCalledOnce()
  })

  it('handles insurer with logo', async () => {
    const { generateOnRiskLetterDocx } = useOnRiskLetterGeneration()
    const data = makeTestData()
    // Minimal valid base64 PNG (1x1 pixel)
    ;(data.insurer as any).logo =
      'iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=='

    await generateOnRiskLetterDocx(data)
    expect(saveAsMock).toHaveBeenCalledOnce()
  })

  it('sets errorMessage on DOCX generation failure', async () => {
    const { generateOnRiskLetterDocx, errorMessage, isGenerating } =
      useOnRiskLetterGeneration()

    // Pass data that will cause an error (null quote)
    await generateOnRiskLetterDocx({
      quote: null,
      scheme: {},
      insurer: {},
      letter: {},
      benefit_summary: []
    })

    expect(isGenerating.value).toBe(false)
    expect(errorMessage.value).toContain('Error generating Word document')
  })

  it('clears previous errorMessage on new generation', async () => {
    const { generateOnRiskLetterDocx, errorMessage } =
      useOnRiskLetterGeneration()

    // First call: force error
    await generateOnRiskLetterDocx({
      quote: null,
      scheme: {},
      insurer: {},
      letter: {},
      benefit_summary: []
    })
    expect(errorMessage.value).not.toBe('')

    // Second call: valid data
    await generateOnRiskLetterDocx(makeTestData())
    expect(errorMessage.value).toBe('')
  })
})
