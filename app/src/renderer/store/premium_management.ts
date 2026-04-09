import { defineStore } from 'pinia'

export interface PremiumScheduleSummary {
  id: number
  scheme_id: number
  scheme_name: string
  month: number
  year: number
  member_count: number
  gross_premium: number
  net_payable: number
  status: 'draft' | 'finalized' | 'invoiced'
  generated_date: string
}

export interface InvoiceSummary {
  id: number
  invoice_number: string
  scheme_id: number
  scheme_name: string
  month: number
  year: number
  issue_date: string
  due_date: string
  gross_amount: number
  net_payable: number
  paid_amount: number
  balance: number
  status: 'draft' | 'sent' | 'partial' | 'paid' | 'overdue'
}

export interface PaymentSummary {
  id: number
  scheme_id: number
  scheme_name: string
  invoice_id: number | null
  invoice_number: string
  payment_date: string
  method: 'eft' | 'debit_order' | 'cheque'
  bank_reference: string
  amount: number
  status: 'matched' | 'unmatched' | 'partial'
}

export interface ArrearsRecord {
  scheme_id: number
  scheme_name: string
  days_0_to_30: number
  days_31_to_60: number
  days_61_to_90: number
  days_over_90: number
  total_outstanding: number
  grace_period_expiry: string | null
  status: 'current' | 'overdue' | 'grace_period' | 'suspended' | 'lapsed'
}

export interface PremiumDashboardKPIs {
  due_this_month: number
  collected: number
  collection_rate: number
  outstanding: number
  overdue: number
  overdue_scheme_count: number
}

export const usePremiumManagementStore = defineStore('premiumManagement', {
  state: () => ({
    // Filters
    selectedSchemeId: null as number | null,
    selectedPeriod: {
      month: null as number | null,
      year: null as number | null
    },

    // Data
    schedules: [] as PremiumScheduleSummary[],
    invoices: [] as InvoiceSummary[],
    payments: [] as PaymentSummary[],
    arrearsData: [] as ArrearsRecord[],
    dashboardKpis: null as PremiumDashboardKPIs | null,

    // UI state
    isLoading: false,
    selectedYear: new Date().getFullYear(),

    // Lookup data
    paymentMethods: ['EFT', 'Debit Order', 'Cheque'],
    scheduleStatuses: ['draft', 'finalized', 'invoiced'],
    invoiceStatuses: ['draft', 'sent', 'partial', 'paid', 'overdue'],
    paymentStatuses: ['matched', 'unmatched', 'partial'],
    arrearsStatuses: [
      'current',
      'overdue',
      'grace_period',
      'suspended',
      'lapsed'
    ],
    contributionTypes: ['Employer Only', 'Split']
  }),

  actions: {
    setSelectedScheme(schemeId: number | null) {
      this.selectedSchemeId = schemeId
    },
    setPeriod(month: number | null, year: number | null) {
      this.selectedPeriod = { month, year }
    },
    resetFilters() {
      this.selectedSchemeId = null
      this.selectedPeriod = { month: null, year: null }
    },
    setSchedules(schedules: PremiumScheduleSummary[]) {
      this.schedules = schedules
    },
    setInvoices(invoices: InvoiceSummary[]) {
      this.invoices = invoices
    },
    setPayments(payments: PaymentSummary[]) {
      this.payments = payments
    },
    setArrearsData(arrears: ArrearsRecord[]) {
      this.arrearsData = arrears
    },
    setDashboardKpis(kpis: PremiumDashboardKPIs) {
      this.dashboardKpis = kpis
    },
    setLoading(loading: boolean) {
      this.isLoading = loading
    }
  }
})
