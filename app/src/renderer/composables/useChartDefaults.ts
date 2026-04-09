import { CHART_COLORS } from '@/renderer/constants/designTokens'

const baseTheme = {
  palette: { fills: CHART_COLORS, strokes: CHART_COLORS }
}

export function barChartOptions(
  overrides: Record<string, unknown> = {}
): Record<string, unknown> {
  return {
    background: { fill: 'transparent' },
    theme: baseTheme,
    legend: { enabled: true, position: 'bottom' },
    height: 260,
    ...overrides
  }
}

export function pieChartOptions(
  overrides: Record<string, unknown> = {}
): Record<string, unknown> {
  return {
    background: { fill: 'transparent' },
    theme: baseTheme,
    legend: { enabled: true, position: 'bottom' },
    height: 260,
    ...overrides
  }
}

export function lineChartOptions(
  overrides: Record<string, unknown> = {}
): Record<string, unknown> {
  return {
    background: { fill: 'transparent' },
    theme: baseTheme,
    legend: { enabled: true, position: 'bottom' },
    height: 260,
    ...overrides
  }
}
