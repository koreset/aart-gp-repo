import { describe, it, expect } from 'vitest'
import { statusCellRenderer } from '../statusCellRenderer'

describe('statusCellRenderer', () => {
  it('renders a known status with correct color and label', () => {
    const html = statusCellRenderer('draft')
    expect(html).toContain('Draft')
    expect(html).toContain('#9E9E9E')
    expect(html).toContain('<span')
  })

  it('renders "approved" status', () => {
    const html = statusCellRenderer('approved')
    expect(html).toContain('Approved')
    expect(html).toContain('#8BC34A')
  })

  it('renders "overdue" status', () => {
    const html = statusCellRenderer('overdue')
    expect(html).toContain('Overdue')
    expect(html).toContain('#F44336')
  })

  it('handles status with spaces by converting to underscored key', () => {
    const html = statusCellRenderer('in force')
    expect(html).toContain('In Force')
    expect(html).toContain('#4CAF50')
  })

  it('is case-insensitive', () => {
    const html = statusCellRenderer('DRAFT')
    expect(html).toContain('Draft')
    expect(html).toContain('#9E9E9E')
  })

  it('falls back to gray for unknown status', () => {
    const html = statusCellRenderer('unknown_status')
    expect(html).toContain('#9E9E9E')
    expect(html).toContain('unknown_status')
  })

  it('handles null/undefined with dash', () => {
    const html = statusCellRenderer(null as any)
    expect(html).toContain('—')
    expect(html).toContain('#9E9E9E')
  })

  it('handles empty string with dash', () => {
    const html = statusCellRenderer('')
    expect(html).toContain('—')
  })

  it('renders as a styled span element', () => {
    const html = statusCellRenderer('active')
    expect(html).toMatch(/^<span style="/)
    expect(html).toMatch(/<\/span>$/)
    expect(html).toContain('border-radius')
    expect(html).toContain('font-weight')
  })
})
