import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { usePremiumManagementStore } from '../premium_management'
import type {
  PremiumScheduleSummary,
  InvoiceSummary
} from '../premium_management'

describe('usePremiumManagementStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('initialises with default state', () => {
    const store = usePremiumManagementStore()
    expect(store.selectedSchemeId).toBeNull()
    expect(store.selectedPeriod).toEqual({ month: null, year: null })
    expect(store.schedules).toEqual([])
    expect(store.invoices).toEqual([])
    expect(store.payments).toEqual([])
    expect(store.arrearsData).toEqual([])
    expect(store.dashboardKpis).toBeNull()
    expect(store.isLoading).toBe(false)
    expect(store.selectedYear).toBe(new Date().getFullYear())
  })

  it('setSelectedScheme updates filter', () => {
    const store = usePremiumManagementStore()
    store.setSelectedScheme(42)
    expect(store.selectedSchemeId).toBe(42)
  })

  it('setPeriod updates month and year', () => {
    const store = usePremiumManagementStore()
    store.setPeriod(3, 2024)
    expect(store.selectedPeriod).toEqual({ month: 3, year: 2024 })
  })

  it('resetFilters clears scheme and period', () => {
    const store = usePremiumManagementStore()
    store.setSelectedScheme(10)
    store.setPeriod(6, 2025)
    store.resetFilters()
    expect(store.selectedSchemeId).toBeNull()
    expect(store.selectedPeriod).toEqual({ month: null, year: null })
  })

  it('setSchedules stores schedule data', () => {
    const store = usePremiumManagementStore()
    const schedules: PremiumScheduleSummary[] = [
      {
        id: 1,
        scheme_id: 1,
        scheme_name: 'Test',
        month: 3,
        year: 2024,
        member_count: 100,
        gross_premium: 50000,
        net_payable: 48000,
        status: 'draft',
        generated_date: '2024-03-01'
      }
    ]
    store.setSchedules(schedules)
    expect(store.schedules).toHaveLength(1)
    expect(store.schedules[0].scheme_name).toBe('Test')
  })

  it('setInvoices stores invoice data', () => {
    const store = usePremiumManagementStore()
    const invoices: InvoiceSummary[] = [
      {
        id: 1,
        invoice_number: 'INV-001',
        scheme_id: 1,
        scheme_name: 'Test',
        month: 3,
        year: 2024,
        issue_date: '2024-03-01',
        due_date: '2024-03-31',
        gross_amount: 50000,
        net_payable: 48000,
        paid_amount: 0,
        balance: 48000,
        status: 'sent'
      }
    ]
    store.setInvoices(invoices)
    expect(store.invoices).toHaveLength(1)
  })

  it('setDashboardKpis stores KPI data', () => {
    const store = usePremiumManagementStore()
    const kpis = {
      due_this_month: 100000,
      collected: 80000,
      collection_rate: 80,
      outstanding: 20000,
      overdue: 5000,
      overdue_scheme_count: 2
    }
    store.setDashboardKpis(kpis)
    expect(store.dashboardKpis).toEqual(kpis)
  })

  it('setLoading toggles loading state', () => {
    const store = usePremiumManagementStore()
    store.setLoading(true)
    expect(store.isLoading).toBe(true)
    store.setLoading(false)
    expect(store.isLoading).toBe(false)
  })

  it('has correct lookup arrays', () => {
    const store = usePremiumManagementStore()
    expect(store.paymentMethods).toContain('EFT')
    expect(store.scheduleStatuses).toContain('draft')
    expect(store.invoiceStatuses).toContain('overdue')
    expect(store.paymentStatuses).toContain('matched')
    expect(store.arrearsStatuses).toContain('lapsed')
    expect(store.contributionTypes).toContain('Employer Only')
  })
})
