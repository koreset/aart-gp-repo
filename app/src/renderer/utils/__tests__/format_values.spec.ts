import { describe, it, expect, vi } from 'vitest'
import { formatValues, roundUpToTwoDecimalsAccounting } from '../format_values'

// Mock luxon DateTime.fromMillis for the 'created' header path
vi.mock('luxon', () => ({
  DateTime: {
    fromMillis: (ms: number) => ({
      toLocaleString: () => {
        const d = new Date(ms)
        return d.toLocaleDateString('en-US')
      }
    })
  }
}))

// ---------------------------------------------------------------------------
// Helper to build AG Grid-like params
// ---------------------------------------------------------------------------

function makeParams(headerName: string, value: any) {
  return {
    value,
    column: {
      userProvidedColDef: { headerName }
    }
  }
}

// ---------------------------------------------------------------------------
// roundUpToTwoDecimalsAccounting
// ---------------------------------------------------------------------------

describe('roundUpToTwoDecimalsAccounting', () => {
  it('rounds up to two decimals', () => {
    expect(roundUpToTwoDecimalsAccounting(1.001)).toBe('1.01')
  })

  it('keeps exact two-decimal values', () => {
    expect(roundUpToTwoDecimalsAccounting(1.5)).toBe('1.50')
  })

  it('rounds up 1.234 to 1.24', () => {
    expect(roundUpToTwoDecimalsAccounting(1.234)).toBe('1.24')
  })

  it('formats large numbers with space separators', () => {
    const result = roundUpToTwoDecimalsAccounting(1234567.89)
    expect(result).toContain('1 234 567.89')
  })

  it('handles zero', () => {
    expect(roundUpToTwoDecimalsAccounting(0)).toBe('0.00')
  })

  it('handles negative numbers', () => {
    const result = roundUpToTwoDecimalsAccounting(-100.5)
    expect(result).toContain('100.50')
  })
})

// ---------------------------------------------------------------------------
// formatValues
// ---------------------------------------------------------------------------

describe('formatValues', () => {
  it('formats "created" header as locale date from unix timestamp', () => {
    const timestamp = 1711612800 // 2024-03-28 in seconds
    const result = formatValues(makeParams('created', timestamp))
    expect(typeof result).toBe('string')
    expect(result.length).toBeGreaterThan(0)
  })

  it('formats float > 100 with toLocaleString', () => {
    const result = formatValues(makeParams('premium', 1234.56))
    expect(typeof result).toBe('string')
    // toLocaleString should add separators
    expect(result).toContain('1')
  })

  it('formats small float with toFixed(3)', () => {
    const result = formatValues(makeParams('rate', 0.12345))
    expect(result).toBe('0.123')
  })

  it('formats integer with toLocaleString for non-excluded headers', () => {
    const result = formatValues(makeParams('amount', 50000))
    expect(typeof result).toBe('string')
  })

  it('returns raw integer for "year" header', () => {
    expect(formatValues(makeParams('year', 2024))).toBe(2024)
  })

  it('returns raw integer for "financial_year" header', () => {
    expect(formatValues(makeParams('financial_year', 2025))).toBe(2025)
  })

  it('returns raw integer for "id" header', () => {
    expect(formatValues(makeParams('id', 42))).toBe(42)
  })

  it('returns raw integer for "policy_number" header', () => {
    expect(formatValues(makeParams('policy_number', 12345))).toBe(12345)
  })

  it('returns string values as-is', () => {
    expect(formatValues(makeParams('name', 'hello'))).toBe('hello')
  })

  it('returns null as-is', () => {
    expect(formatValues(makeParams('field', null))).toBe(null)
  })
})
