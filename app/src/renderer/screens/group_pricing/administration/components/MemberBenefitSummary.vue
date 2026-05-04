<template>
  <v-dialog v-model="dialog" max-width="1200" persistent scrollable>
    <v-card>
      <v-card-title
        class="d-flex justify-space-between align-center bg-primary text-white"
      >
        <div class="d-flex align-center">
          <v-icon class="mr-2">mdi-chart-line</v-icon>
          <span>Benefit Summary - {{ memberName }}</span>
        </div>
        <v-btn
          size="small"
          rounded
          icon="mdi-close"
          variant="text"
          @click="closeDialog"
        />
      </v-card-title>

      <v-card-text class="pa-0">
        <v-container v-if="loading" class="text-center py-8">
          <v-progress-circular indeterminate color="primary" size="64" />
          <div class="mt-4 text-h6">Loading benefit summary...</div>
        </v-container>

        <v-container v-else-if="error" class="text-center py-8">
          <v-icon color="error" size="64">mdi-alert-circle</v-icon>
          <div class="mt-4 text-h6 text-error">{{ error }}</div>
          <v-btn
            size="small"
            rounded
            class="mt-4"
            color="primary"
            @click="loadBenefitSummary"
            >Try Again</v-btn
          >
        </v-container>

        <v-container v-else-if="benefitSummary" class="py-4">
          <!-- Member Information Header -->
          <v-row class="mb-4">
            <v-col>
              <v-card variant="outlined">
                <v-card-text>
                  <v-row>
                    <v-col cols="12" md="2">
                      <div class="text-caption text-grey">Member ID</div>
                      <div class="text-body-1 font-weight-medium">{{
                        benefitSummary.member_id_number
                      }}</div>
                    </v-col>
                    <v-col cols="12" md="2">
                      <div class="text-caption text-grey">Annual Salary</div>
                      <div class="text-body-1 font-weight-medium">{{
                        formatCurrency(benefitSummary.annual_salary)
                      }}</div>
                    </v-col>
                    <v-col cols="12" md="3">
                      <div class="text-caption text-grey">Scheme</div>
                      <div class="text-body-1 font-weight-medium">{{
                        benefitSummary.scheme_name
                      }}</div>
                    </v-col>
                    <v-col cols="12" md="3">
                      <div class="text-caption text-grey">Scheme Category</div>
                      <div class="text-body-1 font-weight-medium">{{
                        benefitSummary.scheme_category || '-'
                      }}</div>
                    </v-col>
                    <v-col cols="12" md="2">
                      <div class="text-caption text-grey">Status</div>
                      <v-chip
                        :color="getStatusColor(benefitSummary.status)"
                        size="small"
                      >
                        {{ (benefitSummary.status || 'ACTIVE').toUpperCase() }}
                      </v-chip>
                    </v-col>
                  </v-row>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Benefits Summary Cards -->
          <v-row>
            <v-col
              v-for="benefit in activeBenefits"
              :key="benefit.benefit_code"
              cols="12"
              md="6"
              lg="4"
            >
              <v-card
                variant="outlined"
                class="benefit-card"
                :class="{ 'benefit-card--active': benefit.is_active }"
              >
                <v-card-title class="d-flex justify-space-between align-center">
                  <span class="text-subtitle-1">{{ benefit.name }}</span>
                  <v-chip
                    :color="benefit.is_active ? 'success' : 'grey'"
                    size="small"
                    variant="tonal"
                  >
                    {{ benefit.is_active ? 'Active' : 'Inactive' }}
                  </v-chip>
                </v-card-title>
                <v-card-text>
                  <div class="benefit-details">
                    <div class="benefit-row">
                      <span class="text-caption text-grey">Sum Assured</span>
                      <span class="text-body-2 font-weight-medium">
                        {{ formatCurrency(benefit.covered_sum_assured) }}
                      </span>
                    </div>
                    <div class="benefit-row">
                      <span class="text-caption text-grey"
                        >Multiple of Salary</span
                      >
                      <span class="text-body-2 font-weight-medium">
                        {{ benefit.salary_multiple }}x
                      </span>
                    </div>
                    <div v-if="benefit.waiting_period" class="benefit-row">
                      <span class="text-caption text-grey">Waiting Period</span>
                      <span class="text-body-2 font-weight-medium">
                        {{ benefit.waiting_period }} months
                      </span>
                    </div>
                  </div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Premium Summary -->
          <v-row class="mt-4">
            <v-col>
              <v-card variant="outlined" class="premium-summary">
                <v-card-title class="bg-info text-white">
                  <v-icon class="mr-2">mdi-calculator</v-icon>
                  Premium Summary
                </v-card-title>
                <v-card-text>
                  <v-row class="mt-1">
                    <v-col cols="12" md="6">
                      <v-card variant="tonal" color="success">
                        <v-card-text class="text-center">
                          <div class="text-h4 font-weight-bold">
                            {{ formatCurrency(totalMonthlyPremium) }}
                          </div>
                          <div class="text-subtitle-1"
                            >Total Monthly Premium</div
                          >
                        </v-card-text>
                      </v-card>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-card variant="tonal" color="primary">
                        <v-card-text class="text-center">
                          <div class="text-h4 font-weight-bold">
                            {{ formatCurrency(totalAnnualPremium) }}
                          </div>
                          <div class="text-subtitle-1"
                            >Total Annual Premium</div
                          >
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                  <v-row class="mt-4">
                    <v-col cols="12">
                      <v-card variant="tonal" color="warning">
                        <v-card-text class="text-center">
                          <div class="text-h4 font-weight-bold">
                            {{ premiumAsPercentageOfSalary }}%
                          </div>
                          <div class="text-subtitle-1"
                            >Premium as % of Annual Salary</div
                          >
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                  <v-row class="mt-4">
                    <v-col cols="12" md="6">
                      <v-card variant="tonal" color="info">
                        <v-card-text class="text-center">
                          <div class="text-h4 font-weight-bold">
                            {{ formatCurrency(funeralMonthlyPremium) }}
                          </div>
                          <div class="text-subtitle-1"
                            >Funeral Monthly Premium</div
                          >
                        </v-card-text>
                      </v-card>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-card variant="tonal" color="secondary">
                        <v-card-text class="text-center">
                          <div class="text-h4 font-weight-bold">
                            {{ formatCurrency(funeralAnnualPremium) }}
                          </div>
                          <div class="text-subtitle-1"
                            >Funeral Annual Premium</div
                          >
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Benefits Coverage Chart -->
          <v-row class="mt-4">
            <v-col>
              <v-card variant="outlined">
                <v-card-title class="bg-secondary text-white">
                  <v-icon class="mr-2">mdi-chart-pie</v-icon>
                  Benefits Coverage Distribution
                </v-card-title>
                <v-card-text>
                  <div class="coverage-chart">
                    <div
                      v-for="benefit in activeBenefits.filter(
                        (b) => b.is_active
                      )"
                      :key="benefit.benefit_code"
                      class="coverage-item"
                    >
                      <div
                        class="d-flex justify-space-between align-center mb-2"
                      >
                        <span class="text-body-2">{{ benefit.name }}</span>
                        <span class="font-weight-bold">{{
                          formatCurrency(benefit.covered_sum_assured)
                        }}</span>
                      </div>
                      <v-progress-linear
                        :model-value="
                          (benefit.covered_sum_assured / totalSumAssured) * 100
                        "
                        :color="getBenefitColor(benefit.benefit_code)"
                        height="8"
                        rounded
                      />
                    </div>
                  </div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Group Family Funeral Benefits -->
          <v-row v-if="benefitSummary.gff" class="mt-4">
            <v-col>
              <v-card variant="outlined">
                <v-card-title class="bg-brown text-white">
                  <v-icon class="mr-2">mdi-home-heart</v-icon>
                  Group Family Funeral Benefits
                </v-card-title>
                <v-card-text>
                  <v-row>
                    <v-col
                      v-for="(amount, member) in benefitSummary.gff"
                      :key="member"
                      cols="12"
                      sm="6"
                      md="4"
                      lg="2.4"
                      class="d-flex align-stretch"
                    >
                      <v-card
                        v-if="
                          member !== 'currency' &&
                          member !== 'children_count' &&
                          member !== 'dependants_count'
                        "
                        variant="tonal"
                        color="brown"
                        class="w-100 gff-benefit-card"
                      >
                        <v-card-text class="text-center">
                          <div class="gff-member-type">
                            {{ formatMemberType(member as string) }}
                            <span
                              v-if="getGffCount(member as string)"
                              class="gff-count"
                            >
                              ({{ getGffCount(member as string) }})
                            </span>
                          </div>
                          <div class="gff-amount">
                            {{ formatCurrency(amount as number) }}
                          </div>
                        </v-card-text>
                      </v-card>
                    </v-col>
                  </v-row>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>

      <v-card-actions class="px-4 py-3">
        <v-btn
          size="small"
          rounded
          color="success"
          variant="outlined"
          @click="downloadSummary"
        >
          <v-icon left>mdi-download</v-icon>
          Download PDF
        </v-btn>
        <v-btn
          size="small"
          rounded
          color="info"
          variant="outlined"
          @click="printSummary"
        >
          <v-icon left>mdi-printer</v-icon>
          Print
        </v-btn>
        <v-spacer />
        <v-btn
          size="small"
          rounded
          color="grey"
          variant="outlined"
          @click="closeDialog"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

import jsPDF from 'jspdf'
import { applyPlugin } from 'jspdf-autotable'

applyPlugin(jsPDF)
// Extend jsPDF type to include autoTable
// declare module 'jspdf' {
//   interface jsPDF {
//     autoTable: (options: any) => jsPDF
//     lastAutoTable: {
//       finalY: number
//     }
//   }
// }

interface BenefitDetail {
  benefit_code: string
  name: string
  is_active: boolean
  covered_sum_assured: number
  salary_multiple: number
  annual_premium: number
  monthly_premium: number
  waiting_period?: number
}

interface GroupFamilyFuneral {
  currency: string
  main_member: number
  spouse: number
  children: number
  parents: number
  dependants: number
  children_count?: number
  dependants_count?: number
}

interface BenefitSummary {
  member_id: number
  member_name: string
  member_id_number: string
  annual_salary: number
  annual_premium: number
  monthly_premium: number
  funeral_annual_premium: number
  funeral_monthly_premium: number
  premium_salary_prop: number
  scheme_name: string
  scheme_category: string
  status: string
  benefits: BenefitDetail[]
  gff?: GroupFamilyFuneral
}

interface Props {
  modelValue: boolean
  memberId: number | null
  memberName: string
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// State
const loading = ref(false)
const error = ref('')
const benefitSummary = ref<BenefitSummary | null>(null)

// Computed
const dialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const activeBenefits = computed(() => {
  return benefitSummary.value?.benefits || []
})

const totalMonthlyPremium = computed(() => {
  return benefitSummary.value?.monthly_premium || 0
})

const totalAnnualPremium = computed(() => {
  return benefitSummary.value?.annual_premium || 0
})

const funeralAnnualPremium = computed(() => {
  return benefitSummary.value?.funeral_annual_premium || 0
})

const funeralMonthlyPremium = computed(() => {
  return benefitSummary.value?.funeral_monthly_premium || 0
})

const totalSumAssured = computed(() => {
  return activeBenefits.value
    .filter((b) => b.is_active)
    .reduce((total, benefit) => total + (benefit.covered_sum_assured || 0), 0)
})

const premiumAsPercentageOfSalary = computed(() => {
  const salary = benefitSummary.value?.annual_salary || 0
  if (salary === 0) return '0.00'
  return ((totalAnnualPremium.value / salary) * 100).toFixed(2)
})

// Watch for dialog opening
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue && props.memberId) {
      loadBenefitSummary()
    }
  }
)

// Methods
const loadBenefitSummary = async () => {
  if (!props.memberId) return

  loading.value = true
  error.value = ''

  try {
    const response = await GroupPricingService.getMemberBenefitSummary(
      props.memberId
    )
    benefitSummary.value = response.data
    console.log('Benefit summary data fetched:', response.data)
    console.log('Loaded benefit summary:', benefitSummary.value)
  } catch (err: any) {
    error.value =
      err.response?.data?.message || 'Failed to load benefit summary'
    console.error('Error loading benefit summary:', err)
  } finally {
    loading.value = false
  }
}

const closeDialog = () => {
  dialog.value = false
  benefitSummary.value = null
  error.value = ''
}

const formatCurrency = (amount: number | null | undefined) => {
  if (!amount) return 'R 0.00'
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(amount)
}

const formatMemberType = (memberType: string) => {
  const typeMap: { [key: string]: string } = {
    main_member: 'Main Member',
    spouse: 'Spouse',
    children: 'Children',
    parents: 'Parents',
    dependants: 'Dependants'
  }
  return typeMap[memberType] || memberType
}

const getGffCount = (memberType: string) => {
  if (!benefitSummary.value?.gff) return null

  if (memberType === 'children' && benefitSummary.value.gff.children_count) {
    return benefitSummary.value.gff.children_count
  }

  if (
    memberType === 'dependants' &&
    benefitSummary.value.gff.dependants_count
  ) {
    return benefitSummary.value.gff.dependants_count
  }

  return null
}

const getStatusColor = (status: string | undefined) => {
  switch ((status || 'active').toLowerCase()) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'error'
    case 'pending':
      return 'warning'
    case 'suspended':
      return 'orange'
    default:
      return 'grey'
  }
}

const getBenefitColor = (benefitCode: string) => {
  const colors: { [key: string]: string } = {
    GLA: 'blue',
    SGLA: 'green',
    PTD: 'orange',
    CI: 'red',
    TTD: 'purple',
    PHI: 'teal',
    GFF: 'brown'
  }
  return colors[benefitCode] || 'grey'
}

const downloadSummary = () => {
  if (!benefitSummary.value) return

  try {
    // Create new PDF document
    // eslint-disable-next-line new-cap
    const doc: any = new jsPDF()
    const pageWidth = doc.internal.pageSize.width
    const margin = 20

    // Header
    doc.setFontSize(20)
    doc.setFont('helvetica', 'bold')
    doc.text('Member Benefit Summary', pageWidth / 2, 25, { align: 'center' })

    // Member Information
    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text('Member Information', margin, 45)

    doc.setFontSize(10)
    doc.setFont('helvetica', 'normal')
    let yPos = 55

    const memberInfo = [
      ['Member Name:', benefitSummary.value.member_name],
      ['Member ID:', benefitSummary.value.member_id_number],
      ['Annual Salary:', formatCurrency(benefitSummary.value.annual_salary)],
      ['Scheme:', benefitSummary.value.scheme_name],
      ['Scheme Category:', benefitSummary.value.scheme_category || '-'],
      ['Status:', (benefitSummary.value.status || 'ACTIVE').toUpperCase()]
    ]

    memberInfo.forEach(([label, value]) => {
      doc.text(label, margin, yPos)
      doc.text(value, margin + 60, yPos)
      yPos += 8
    })

    yPos += 10

    // Benefits Table
    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text('Benefits Details', margin, yPos)
    yPos += 10

    const benefitsData = activeBenefits.value.map((benefit) => [
      benefit.name,
      benefit.is_active ? 'Active' : 'Inactive',
      formatCurrency(benefit.covered_sum_assured),
      `${benefit.salary_multiple}x`,
      formatCurrency(benefit.monthly_premium || 0),
      formatCurrency(benefit.annual_premium || 0),
      benefit.waiting_period ? `${benefit.waiting_period} months` : 'N/A'
    ])

    doc.autoTable({
      startY: yPos,
      head: [
        [
          'Benefit',
          'Status',
          'Sum Assured',
          'Salary Multiple',
          'Monthly Premium',
          'Annual Premium',
          'Waiting Period'
        ]
      ],
      body: benefitsData,
      theme: 'striped',
      headStyles: {
        fillColor: [63, 81, 181],
        textColor: 255,
        fontStyle: 'bold'
      },
      alternateRowStyles: {
        fillColor: [245, 245, 245]
      },
      styles: {
        fontSize: 8,
        cellPadding: 3
      },
      columnStyles: {
        0: { cellWidth: 35 },
        1: { cellWidth: 20, halign: 'center' },
        2: { cellWidth: 25, halign: 'right' },
        3: { cellWidth: 20, halign: 'center' },
        4: { cellWidth: 25, halign: 'right' },
        5: { cellWidth: 25, halign: 'right' },
        6: { cellWidth: 20, halign: 'center' }
      }
    })

    yPos = doc.lastAutoTable.finalY + 15

    // Group Family Funeral Benefits (if available)
    if (benefitSummary.value.gff) {
      doc.setFontSize(14)
      doc.setFont('helvetica', 'bold')
      doc.text('Group Family Funeral Benefits', margin, yPos)
      yPos += 10

      const gffData = [
        ['Main Member:', formatCurrency(benefitSummary.value.gff.main_member)],
        ['Spouse:', formatCurrency(benefitSummary.value.gff.spouse)],
        [
          `Children${benefitSummary.value.gff.children_count ? ` (${benefitSummary.value.gff.children_count})` : ''}:`,
          formatCurrency(benefitSummary.value.gff.children)
        ],
        ['Parents:', formatCurrency(benefitSummary.value.gff.parents)],
        [
          `Dependants${benefitSummary.value.gff.dependants_count ? ` (${benefitSummary.value.gff.dependants_count})` : ''}:`,
          formatCurrency(benefitSummary.value.gff.dependants)
        ]
      ]

      doc.autoTable({
        startY: yPos,
        body: gffData,
        theme: 'striped',
        headStyles: {
          fillColor: [121, 85, 72],
          textColor: 255
        },
        alternateRowStyles: {
          fillColor: [245, 245, 245]
        },
        styles: {
          fontSize: 10,
          cellPadding: 5
        },
        columnStyles: {
          0: { cellWidth: 60, fontStyle: 'bold' },
          1: { cellWidth: 40, halign: 'right' }
        }
      })

      yPos = doc.lastAutoTable.finalY + 15
    }

    // Premium Summary
    doc.setFontSize(14)
    doc.setFont('helvetica', 'bold')
    doc.text('Premium Summary', margin, yPos)
    yPos += 10

    const premiumData = [
      ['Total Monthly Premium:', formatCurrency(totalMonthlyPremium.value)],
      ['Total Annual Premium:', formatCurrency(totalAnnualPremium.value)],
      ['Premium as % of Salary:', premiumAsPercentageOfSalary.value + '%']
    ]

    doc.autoTable({
      startY: yPos,
      body: premiumData,
      theme: 'plain',
      styles: {
        fontSize: 10,
        cellPadding: 5
      },
      columnStyles: {
        0: { cellWidth: 60, fontStyle: 'bold' },
        1: { cellWidth: 40, halign: 'right' }
      }
    })

    // Footer
    const pageHeight = doc.internal.pageSize.height
    doc.setFontSize(8)
    doc.setFont('helvetica', 'italic')
    doc.text(
      'Generated on: ' + new Date().toLocaleDateString(),
      margin,
      pageHeight - 15
    )
    doc.text('Page 1', pageWidth - margin - 20, pageHeight - 15)

    // Save the PDF
    const fileName =
      'Benefit_Summary_' +
      benefitSummary.value.member_name.replace(/\s+/g, '_') +
      '_' +
      new Date().toISOString().split('T')[0] +
      '.pdf'
    doc.save(fileName)

    console.log('Benefit summary PDF downloaded for member:', props.memberName)
  } catch (error) {
    console.error('Error generating PDF:', error)
    // You might want to show a user-friendly error message here
    alert('Failed to generate PDF. Please try again.')
  }
}

const escapeHtml = (value: unknown): string => {
  return String(value ?? '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

const buildPrintHtml = (): string => {
  const data = benefitSummary.value
  if (!data) return ''

  const status = (data.status || 'ACTIVE').toUpperCase()

  const benefitCards = activeBenefits.value
    .map((b) => {
      const waitingRow = b.waiting_period
        ? `<div class="row"><span class="muted">Waiting Period</span><span class="value">${b.waiting_period} months</span></div>`
        : ''
      return `
        <div class="card ${b.is_active ? 'active' : 'inactive'}">
          <div class="card-head">
            <span class="card-title">${escapeHtml(b.name)}</span>
            <span class="badge ${b.is_active ? 'badge-active' : 'badge-inactive'}">${b.is_active ? 'Active' : 'Inactive'}</span>
          </div>
          <div class="row"><span class="muted">Sum Assured</span><span class="value">${formatCurrency(b.covered_sum_assured)}</span></div>
          <div class="row"><span class="muted">Multiple of Salary</span><span class="value">${b.salary_multiple}x</span></div>
          ${waitingRow}
        </div>
      `
    })
    .join('')

  const total = totalSumAssured.value || 1
  const coverageItems = activeBenefits.value
    .filter((b) => b.is_active)
    .map((b) => {
      const pct = Math.max(
        0,
        Math.min(100, ((b.covered_sum_assured || 0) / total) * 100)
      )
      const palette: { [k: string]: string } = {
        GLA: '#1976d2',
        SGLA: '#43a047',
        PTD: '#fb8c00',
        CI: '#e53935',
        TTD: '#8e24aa',
        PHI: '#00897b',
        GFF: '#795548'
      }
      const color = palette[b.benefit_code] || '#757575'
      return `
        <div class="cov-item">
          <div class="cov-row">
            <span>${escapeHtml(b.name)}</span>
            <span class="bold">${formatCurrency(b.covered_sum_assured)}</span>
          </div>
          <div class="bar-track"><div class="bar-fill" style="width:${pct.toFixed(2)}%; background:${color};"></div></div>
        </div>
      `
    })
    .join('')

  let gffSection = ''
  if (data.gff) {
    const order: Array<keyof GroupFamilyFuneral> = [
      'main_member',
      'spouse',
      'children',
      'parents',
      'dependants'
    ]
    const items = order
      .filter((key) => key in data.gff!)
      .map((key) => {
        const amount = data.gff![key] as number
        const count = getGffCount(key as string)
        const countSuffix = count ? ` (${count})` : ''
        return `
          <div class="gff-item">
            <div class="gff-label">${formatMemberType(key as string)}${countSuffix}</div>
            <div class="gff-value">${formatCurrency(amount)}</div>
          </div>
        `
      })
      .join('')
    gffSection = `
      <section class="section">
        <h2 class="section-title brown">Group Family Funeral Benefits</h2>
        <div class="gff-grid">${items}</div>
      </section>
    `
  }

  const generated = new Date().toLocaleString('en-ZA')

  return `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<title>Benefit Summary - ${escapeHtml(props.memberName)}</title>
<style>
  *, *::before, *::after { box-sizing: border-box; }
  html, body { margin: 0; padding: 0; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    color: #212121;
    font-size: 11pt;
    -webkit-print-color-adjust: exact;
    print-color-adjust: exact;
  }
  .page { padding: 16mm 14mm; }
  .doc-header { display: flex; justify-content: space-between; align-items: flex-end; border-bottom: 2px solid #1976d2; padding-bottom: 8pt; margin-bottom: 14pt; }
  .doc-title { font-size: 18pt; font-weight: 700; color: #1976d2; margin: 0; }
  .doc-sub { font-size: 10pt; color: #555; margin-top: 2pt; }
  .doc-meta { font-size: 9pt; color: #666; text-align: right; }
  .section { margin-bottom: 14pt; page-break-inside: avoid; }
  .section-title { font-size: 12pt; font-weight: 700; color: #fff; background: #1976d2; padding: 6pt 10pt; margin: 0 0 8pt 0; border-radius: 3pt; }
  .section-title.info { background: #0288d1; }
  .section-title.brown { background: #795548; }
  .section-title.secondary { background: #455a64; }

  .info-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8pt; border: 1px solid #e0e0e0; padding: 10pt; border-radius: 3pt; }
  .info-cell .lbl { font-size: 8pt; color: #757575; text-transform: uppercase; letter-spacing: 0.5pt; }
  .info-cell .val { font-size: 10pt; font-weight: 600; margin-top: 2pt; word-break: break-word; }
  .status-pill { display: inline-block; padding: 2pt 8pt; border-radius: 10pt; font-size: 8pt; font-weight: 700; color: #fff; }
  .status-active { background: #2e7d32; }
  .status-inactive { background: #c62828; }
  .status-pending { background: #f9a825; }
  .status-other { background: #757575; }

  .card-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8pt; }
  .card { border: 1px solid #e0e0e0; border-radius: 3pt; padding: 8pt 10pt; page-break-inside: avoid; }
  .card.active { border-color: #81c784; background: #f1f8f1; }
  .card.inactive { background: #fafafa; }
  .card-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 6pt; }
  .card-title { font-weight: 700; font-size: 10pt; }
  .badge { font-size: 7pt; padding: 1pt 6pt; border-radius: 8pt; font-weight: 700; }
  .badge-active { background: #c8e6c9; color: #1b5e20; }
  .badge-inactive { background: #eeeeee; color: #555; }
  .row { display: flex; justify-content: space-between; padding: 3pt 0; border-bottom: 1px dashed #eee; font-size: 9.5pt; }
  .row:last-child { border-bottom: none; }
  .muted { color: #757575; }
  .value { font-weight: 600; }
  .bold { font-weight: 700; }

  .premium-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 8pt; }
  .premium-grid.full { grid-template-columns: 1fr; }
  .premium-tile { border: 1px solid #e0e0e0; border-radius: 3pt; padding: 10pt; text-align: center; }
  .premium-tile .amount { font-size: 16pt; font-weight: 700; }
  .premium-tile .lbl { font-size: 9pt; color: #555; margin-top: 2pt; }
  .tile-success { background: #e8f5e9; border-color: #a5d6a7; color: #1b5e20; }
  .tile-primary { background: #e3f2fd; border-color: #90caf9; color: #0d47a1; }
  .tile-warning { background: #fff8e1; border-color: #ffe082; color: #e65100; }
  .tile-info    { background: #e1f5fe; border-color: #81d4fa; color: #01579b; }
  .tile-secondary { background: #eceff1; border-color: #b0bec5; color: #263238; }

  .cov-item { margin-bottom: 6pt; }
  .cov-row { display: flex; justify-content: space-between; font-size: 9.5pt; margin-bottom: 2pt; }
  .bar-track { height: 6pt; background: #eee; border-radius: 3pt; overflow: hidden; }
  .bar-fill { height: 100%; }

  .gff-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 6pt; }
  .gff-item { border: 1px solid #d7ccc8; background: #efebe9; border-radius: 3pt; padding: 8pt; text-align: center; }
  .gff-label { font-size: 9pt; color: #5d4037; font-weight: 600; margin-bottom: 4pt; }
  .gff-value { font-size: 11pt; font-weight: 700; color: #3e2723; }

  .footer { margin-top: 16pt; padding-top: 6pt; border-top: 1px solid #e0e0e0; font-size: 8pt; color: #888; display: flex; justify-content: space-between; }

  @page { size: A4 portrait; margin: 10mm; }
</style>
</head>
<body>
  <div class="page">
    <div class="doc-header">
      <div>
        <h1 class="doc-title">Member Benefit Summary</h1>
        <div class="doc-sub">${escapeHtml(props.memberName)}</div>
      </div>
      <div class="doc-meta">Generated: ${escapeHtml(generated)}</div>
    </div>

    <section class="section">
      <div class="info-grid">
        <div class="info-cell"><div class="lbl">Member ID</div><div class="val">${escapeHtml(data.member_id_number)}</div></div>
        <div class="info-cell"><div class="lbl">Annual Salary</div><div class="val">${formatCurrency(data.annual_salary)}</div></div>
        <div class="info-cell"><div class="lbl">Scheme</div><div class="val">${escapeHtml(data.scheme_name)}</div></div>
        <div class="info-cell"><div class="lbl">Scheme Category</div><div class="val">${escapeHtml(data.scheme_category || '-')}</div></div>
        <div class="info-cell"><div class="lbl">Status</div><div class="val"><span class="status-pill ${
          status === 'ACTIVE'
            ? 'status-active'
            : status === 'INACTIVE'
              ? 'status-inactive'
              : status === 'PENDING'
                ? 'status-pending'
                : 'status-other'
        }">${escapeHtml(status)}</span></div></div>
      </div>
    </section>

    <section class="section">
      <h2 class="section-title">Benefits</h2>
      <div class="card-grid">${benefitCards}</div>
    </section>

    <section class="section">
      <h2 class="section-title info">Premium Summary</h2>
      <div class="premium-grid">
        <div class="premium-tile tile-success">
          <div class="amount">${formatCurrency(totalMonthlyPremium.value)}</div>
          <div class="lbl">Total Monthly Premium</div>
        </div>
        <div class="premium-tile tile-primary">
          <div class="amount">${formatCurrency(totalAnnualPremium.value)}</div>
          <div class="lbl">Total Annual Premium</div>
        </div>
      </div>
      <div class="premium-grid full" style="margin-top:8pt;">
        <div class="premium-tile tile-warning">
          <div class="amount">${escapeHtml(premiumAsPercentageOfSalary.value)}%</div>
          <div class="lbl">Premium as % of Annual Salary</div>
        </div>
      </div>
      <div class="premium-grid" style="margin-top:8pt;">
        <div class="premium-tile tile-info">
          <div class="amount">${formatCurrency(funeralMonthlyPremium.value)}</div>
          <div class="lbl">Funeral Monthly Premium</div>
        </div>
        <div class="premium-tile tile-secondary">
          <div class="amount">${formatCurrency(funeralAnnualPremium.value)}</div>
          <div class="lbl">Funeral Annual Premium</div>
        </div>
      </div>
    </section>

    ${
      coverageItems
        ? `<section class="section">
      <h2 class="section-title secondary">Benefits Coverage Distribution</h2>
      ${coverageItems}
    </section>`
        : ''
    }

    ${gffSection}

    <div class="footer">
      <span>${escapeHtml(props.memberName)} · ${escapeHtml(data.member_id_number)}</span>
      <span>AART Group Pricing</span>
    </div>
  </div>
</body>
</html>`
}

const printSummary = () => {
  if (!benefitSummary.value) return
  const html = buildPrintHtml()

  const iframe = document.createElement('iframe')
  iframe.setAttribute('aria-hidden', 'true')
  iframe.style.position = 'fixed'
  iframe.style.right = '0'
  iframe.style.bottom = '0'
  iframe.style.width = '0'
  iframe.style.height = '0'
  iframe.style.border = '0'
  iframe.style.opacity = '0'
  document.body.appendChild(iframe)

  const cleanup = () => {
    if (iframe.parentNode) iframe.parentNode.removeChild(iframe)
  }

  iframe.onload = () => {
    try {
      const win = iframe.contentWindow
      if (!win) {
        cleanup()
        return
      }
      win.focus()
      win.print()
      win.onafterprint = () => cleanup()
      // Fallback in case onafterprint never fires
      setTimeout(cleanup, 60000)
    } catch (err) {
      console.error('Failed to print benefit summary:', err)
      cleanup()
    }
  }

  iframe.srcdoc = html
}
</script>

<style scoped>
.benefit-card {
  height: 100%;
  transition: all 0.3s ease;
}

.benefit-card--active {
  border-color: rgba(76, 175, 80, 0.5);
  background-color: rgba(76, 175, 80, 0.02);
}

.benefit-details {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.benefit-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.benefit-row:last-child {
  border-bottom: none;
}

.premium-summary .v-card-title {
  padding: 12px 16px;
}

.coverage-chart {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.coverage-item {
  padding: 8px 0;
}

.v-progress-linear {
  margin-top: 4px;
}

/* Group Family Funeral Styles */
.bg-brown {
  background-color: #795548 !important;
}

.gff-benefit-card {
  min-height: 100px;
  display: flex;
  align-items: center;
}

.gff-member-type {
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 8px;
  color: #5d4037;
  text-transform: capitalize;
}

.gff-amount {
  font-size: 1.1rem;
  font-weight: bold;
  color: #3e2723;
}

.gff-count {
  font-size: 0.75rem;
  font-weight: 400;
  color: #6d4c41;
  margin-left: 4px;
}

@media (max-width: 960px) {
  .benefit-details {
    gap: 8px;
  }

  .benefit-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
    text-align: left;
  }

  .gff-benefit-card {
    min-height: 80px;
  }

  .gff-member-type {
    font-size: 0.8rem;
  }

  .gff-amount {
    font-size: 1rem;
  }
}
</style>
