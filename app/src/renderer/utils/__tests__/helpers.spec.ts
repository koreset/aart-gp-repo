import { describe, it, expect } from 'vitest'
import { formatDateString, toMinutes } from '../helpers'

// ---------------------------------------------------------------------------
// formatDateString
// ---------------------------------------------------------------------------

describe('formatDateString', () => {
  it('returns year only when getyear=true, getmonth=false, getday=false', () => {
    expect(formatDateString('2024-03-15', true, false, false)).toBe(2024)
  })

  it('returns YYYY-MM when getyear=true, getmonth=true, getday=false', () => {
    expect(formatDateString('2024-03-15', true, true, false)).toBe('2024-03')
  })

  it('returns YYYY-MM-DD when all flags are true', () => {
    expect(formatDateString('2024-03-15', true, true, true)).toBe('2024-03-15')
  })

  it('zero-pads single-digit months', () => {
    expect(formatDateString('2024-01-15', true, true, false)).toBe('2024-01')
  })

  it('zero-pads single-digit days', () => {
    expect(formatDateString('2024-03-05', true, true, true)).toBe('2024-03-05')
  })

  it('handles double-digit months and days without padding', () => {
    expect(formatDateString('2024-12-25', true, true, true)).toBe('2024-12-25')
  })

  it('returns undefined when no flags are true', () => {
    expect(formatDateString('2024-03-15', false, false, false)).toBeUndefined()
  })
})

// ---------------------------------------------------------------------------
// toMinutes
// ---------------------------------------------------------------------------

describe('toMinutes', () => {
  it('converts 1 to "1 m, 0 s"', () => {
    expect(toMinutes(1)).toBe('1 m, 0 s')
  })

  it('converts 0 to "0 m, 0 s"', () => {
    expect(toMinutes(0)).toBe('0 m, 0 s')
  })

  it('converts 0.5 to "0 m, 18 s"', () => {
    // 0.5 * 60 = 30; minutes=0; seconds = ((30%60)/100)*60 = 18
    expect(toMinutes(0.5)).toBe('0 m, 18 s')
  })

  it('converts 2.5 to "2 m, 18 s"', () => {
    // 2.5 * 60 = 150; minutes=2; seconds = ((150%60)/100)*60 = 18
    expect(toMinutes(2.5)).toBe('2 m, 18 s')
  })
})
