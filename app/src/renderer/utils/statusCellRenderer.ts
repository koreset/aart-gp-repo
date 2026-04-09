import { STATUS_COLORS } from '@/renderer/constants/designTokens'

export function statusCellRenderer(status: string): string {
  const key = (status ?? '').toLowerCase().replace(/\s+/g, '_')
  const token = STATUS_COLORS[key] ?? { hex: '#9E9E9E', label: status || '—' }
  return `<span style="
    padding: 2px 10px;
    border-radius: 12px;
    font-size: 11px;
    font-weight: 600;
    background: ${token.hex}22;
    color: ${token.hex};
    border: 1px solid ${token.hex}55;
    white-space: nowrap;
  ">${token.label}</span>`
}
