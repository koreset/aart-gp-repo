import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  fmtCurrency,
  fmtDate,
  fmtPercent,
  fmtNumber,
  fmtDateAgo,
  currencyFormatter,
  dateFormatter,
  percentFormatter
} from '../formatters'
import { DateTime, Settings } from 'luxon'

// ---------------------------------------------------------------------------
// fmtCurrency
// ---------------------------------------------------------------------------

describe('fmtCurrency', () => {
  it('formats a number as ZAR currency', () => {
    const result = fmtCurrency(1234.56)
    // Intl output may vary by env; just check the number is present
    expect(result).toContain('1')
    expect(result).toContain('234')
    expect(result).toContain('56')
  })

  it('returns dash for null', () => {
    expect(fmtCurrency(null)).toBe('—')
  })

  it('returns dash for undefined', () => {
    expect(fmtCurrency(undefined)).toBe('—')
  })

  it('returns dash for empty string', () => {
    expect(fmtCurrency('')).toBe('—')
  })

  it('returns dash for non-numeric string', () => {
    expect(fmtCurrency('abc')).toBe('—')
  })

  it('formats zero', () => {
    const result = fmtCurrency(0)
    expect(result).toContain('0')
    expect(result).not.toBe('—')
  })

  it('accepts a different currency code', () => {
    const result = fmtCurrency(100, 'USD')
    expect(result).toBeTruthy()
    expect(result).not.toBe('—')
  })

  it('handles string numbers', () => {
    const result = fmtCurrency('500')
    expect(result).toContain('500')
    expect(result).not.toBe('—')
  })
})

// ---------------------------------------------------------------------------
// fmtDate
// ---------------------------------------------------------------------------

describe('fmtDate', () => {
  it('formats ISO date in medium format by default', () => {
    const result = fmtDate('2024-03-15')
    expect(result).toBe('15 Mar 2024')
  })

  it('formats in short format (MMM yyyy)', () => {
    const result = fmtDate('2024-03-15', 'short')
    expect(result).toBe('Mar 2024')
  })

  it('formats in long format (dd MMMM yyyy)', () => {
    const result = fmtDate('2024-03-15', 'long')
    expect(result).toBe('15 March 2024')
  })

  it('returns dash for falsy value', () => {
    expect(fmtDate(null)).toBe('—')
    expect(fmtDate(undefined)).toBe('—')
    expect(fmtDate('')).toBe('—')
  })

  it('returns original string for invalid date', () => {
    expect(fmtDate('not-a-date')).toBe('not-a-date')
  })

  it('handles DateTime objects', () => {
    const dt = DateTime.fromISO('2024-06-01')
    const result = fmtDate(dt)
    expect(result).toBe('01 Jun 2024')
  })
})

// ---------------------------------------------------------------------------
// fmtPercent
// ---------------------------------------------------------------------------

describe('fmtPercent', () => {
  it('formats a number as percentage with 1 decimal by default', () => {
    expect(fmtPercent(12.34)).toBe('12.3%')
  })

  it('uses custom decimal places', () => {
    expect(fmtPercent(12.345, 2)).toBe('12.35%')
  })

  it('returns dash for null', () => {
    expect(fmtPercent(null)).toBe('—')
  })

  it('returns dash for undefined', () => {
    expect(fmtPercent(undefined)).toBe('—')
  })

  it('returns dash for empty string', () => {
    expect(fmtPercent('')).toBe('—')
  })

  it('returns dash for non-numeric string', () => {
    expect(fmtPercent('abc')).toBe('—')
  })

  it('formats zero', () => {
    expect(fmtPercent(0)).toBe('0.0%')
  })
})

// ---------------------------------------------------------------------------
// fmtNumber
// ---------------------------------------------------------------------------

describe('fmtNumber', () => {
  it('formats a number with locale separators', () => {
    const result = fmtNumber(1234567)
    // locale-dependent, but should contain digits
    expect(result).not.toBe('—')
    expect(result).toContain('1')
  })

  it('returns dash for null', () => {
    expect(fmtNumber(null)).toBe('—')
  })

  it('returns dash for undefined', () => {
    expect(fmtNumber(undefined)).toBe('—')
  })

  it('returns dash for empty string', () => {
    expect(fmtNumber('')).toBe('—')
  })

  it('returns dash for non-numeric', () => {
    expect(fmtNumber('abc')).toBe('—')
  })

  it('handles string numbers', () => {
    expect(fmtNumber('42')).not.toBe('—')
  })
})

// ---------------------------------------------------------------------------
// fmtDateAgo
// ---------------------------------------------------------------------------

describe('fmtDateAgo', () => {
  it('returns dash for falsy value', () => {
    expect(fmtDateAgo(null)).toBe('—')
    expect(fmtDateAgo(undefined)).toBe('—')
    expect(fmtDateAgo('')).toBe('—')
  })

  it('returns dash for invalid date', () => {
    expect(fmtDateAgo('not-a-date')).toBe('—')
  })

  it('returns a relative string for a valid ISO date', () => {
    const recent = DateTime.now().minus({ hours: 2 }).toISO()
    const result = fmtDateAgo(recent)
    expect(result).not.toBe('—')
    expect(typeof result).toBe('string')
  })
})

// ---------------------------------------------------------------------------
// AG Grid valueFormatter wrappers
// ---------------------------------------------------------------------------

describe('AG Grid formatter wrappers', () => {
  it('currencyFormatter delegates to fmtCurrency', () => {
    const result = currencyFormatter({ value: 100 })
    expect(result).not.toBe('—')
  })

  it('currencyFormatter handles null value', () => {
    expect(currencyFormatter({ value: null })).toBe('—')
  })

  it('dateFormatter delegates to fmtDate', () => {
    const result = dateFormatter({ value: '2024-01-15' })
    expect(result).toBe('15 Jan 2024')
  })

  it('percentFormatter delegates to fmtPercent', () => {
    expect(percentFormatter({ value: 55.5 })).toBe('55.5%')
  })
})
