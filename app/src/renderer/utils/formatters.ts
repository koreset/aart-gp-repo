import { DateTime } from 'luxon'

export function fmtCurrency(val: unknown, currency = 'ZAR'): string {
  if (val === null || val === undefined || val === '') return '—'
  const num = Number(val)
  if (isNaN(num)) return '—'
  return new Intl.NumberFormat('en-ZA', { style: 'currency', currency }).format(
    num
  )
}

export function fmtDate(
  val: unknown,
  format: 'short' | 'medium' | 'long' = 'medium'
): string {
  if (!val) return '—'
  const dt = DateTime.isDateTime(val) ? val : DateTime.fromISO(String(val))
  if (!dt.isValid) return String(val)
  switch (format) {
    case 'short':
      return dt.toFormat('MMM yyyy')
    case 'long':
      return dt.toFormat('dd MMMM yyyy')
    default:
      return dt.toFormat('dd MMM yyyy')
  }
}

export function fmtPercent(val: unknown, decimals = 1): string {
  if (val === null || val === undefined || val === '') return '—'
  const num = Number(val)
  if (isNaN(num)) return '—'
  return `${num.toFixed(decimals)}%`
}

export function fmtNumber(val: unknown): string {
  if (val === null || val === undefined || val === '') return '—'
  const num = Number(val)
  if (isNaN(num)) return '—'
  return num.toLocaleString('en-ZA')
}

export function fmtDateAgo(val: unknown): string {
  if (!val) return '—'
  const dt = DateTime.fromISO(String(val))
  if (!dt.isValid) return '—'
  return dt.toRelative() ?? '—'
}

// AG Grid valueFormatter wrappers
export const currencyFormatter = (p: { value: unknown }) => fmtCurrency(p.value)
export const dateFormatter = (p: { value: unknown }) => fmtDate(p.value)
export const percentFormatter = (p: { value: unknown }) => fmtPercent(p.value)
