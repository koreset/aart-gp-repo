<template>
  <base-card :show-actions="false">
    <template #header>
      <span class="headline">Reinsurance Premium Summary</span>
    </template>
    <template #default>
      <v-row v-if="rowData.length > 0">
        <v-col>
          <group-pricing-data-grid
            ref="dataGridRef"
            :columnDefs="columnDefs"
            :show-close-button="false"
            :rowData="rowData"
            :table-title="'Reinsurance Premium Summary By Category'"
            :suppressAutoSize="true"
            :density="'compact'"
            :show-export="true"
          />
        </v-col>
      </v-row>
      <v-row v-else>
        <v-col>
          <v-alert type="info" density="compact" variant="tonal">
            No reinsurance summary data available yet. Run the quote calculation
            and ensure the reinsurance rate/loading tables are populated.
          </v-alert>
        </v-col>
      </v-row>
    </template>
  </base-card>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import GroupPricingDataGrid from '@/renderer/components/tables/GroupPricingDataGrid.vue'

const glaBenefitTitle = ref('GLA')
const sglaBenefitTitle = ref('SGLA')
const ptdBenefitTitle = ref('PTD')
const ciBenefitTitle = ref('CI')
const phiBenefitTitle = ref('PHI')
const ttdBenefitTitle = ref('TTD')
const familyFuneralBenefitTitle = ref('Group Funeral')
const benefitMaps = ref([])

const props = defineProps({
  resultSummaries: { type: Array, required: true },
  quote: { type: Object, required: false }
})

const roundUpToTwoDecimalsAccounting = (num: any) => {
  if (num === null || num === undefined || num === '') return '-'
  const n = typeof num === 'string' ? parseFloat(num) : num
  if (Number.isNaN(n)) return '-'
  const roundedNum = Math.ceil(n * 100) / 100
  return roundedNum
    .toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    })
    .replace(/,/g, ' ')
}

const baseCellStyle = (params: any) => {
  if (params.data?.isSubtotal) {
    return { backgroundColor: '#f0f0f0', textAlign: 'right' }
  }
  if (params.data?.isSectionHeader) {
    return { backgroundColor: '#e3f2fd', textAlign: 'right' }
  }
  return { textAlign: 'right' }
}

const numericFormatter = (params: any) => {
  if (
    params.value === null ||
    params.value === undefined ||
    params.value === ''
  )
    return '-'
  return typeof params.value === 'string'
    ? params.value
    : roundUpToTwoDecimalsAccounting(params.value)
}

const columnDefs: any = ref([
  {
    field: 'category',
    headerName: 'Category',
    rowGroup: true,
    hide: true
  },
  {
    field: 'benefit',
    headerName: 'Benefit',
    width: 220,
    minWidth: 180,
    maxWidth: 320,
    flex: 1,
    resizable: true,
    cellStyle: (params: any) => {
      if (params.data?.isSubtotal) {
        return { fontWeight: 'bold', backgroundColor: '#f0f0f0' }
      }
      if (params.data?.isSectionHeader) {
        return {
          fontWeight: 'bold',
          backgroundColor: '#e3f2fd',
          fontStyle: 'italic'
        }
      }
      return { fontWeight: 'bold' }
    }
  },
  {
    field: 'annualSalary',
    headerName: 'Annual Salary',
    width: 160,
    minWidth: 140,
    maxWidth: 220,
    resizable: true,
    valueFormatter: numericFormatter,
    type: 'rightAligned',
    cellStyle: baseCellStyle
  },
  {
    field: 'totalSumAssured',
    headerName: 'Total Sum Assured',
    width: 180,
    minWidth: 150,
    maxWidth: 240,
    resizable: true,
    valueFormatter: numericFormatter,
    type: 'rightAligned',
    cellStyle: baseCellStyle
  },
  {
    field: 'annualPremium',
    headerName: 'Annual Premium',
    width: 170,
    minWidth: 140,
    maxWidth: 220,
    resizable: true,
    valueFormatter: numericFormatter,
    type: 'rightAligned',
    cellStyle: baseCellStyle
  },
  {
    field: 'totalCededSumAssured',
    headerName: 'Total Ceded Sum Assured',
    width: 210,
    minWidth: 170,
    maxWidth: 260,
    resizable: true,
    valueFormatter: numericFormatter,
    type: 'rightAligned',
    cellStyle: baseCellStyle
  },
  {
    field: 'reinsurancePremium',
    headerName: 'Reinsurance Premium',
    width: 200,
    minWidth: 170,
    maxWidth: 260,
    resizable: true,
    valueFormatter: numericFormatter,
    type: 'rightAligned',
    cellStyle: baseCellStyle
  },
  {
    field: 'reinsurancePremiumProportion',
    headerName: 'Reinsurance Premium %',
    width: 200,
    minWidth: 170,
    maxWidth: 260,
    resizable: true,
    type: 'rightAligned',
    cellStyle: baseCellStyle
  }
])
const rowData: any = ref([])

// Mirror the "is benefit enabled" guard from QuoteBenefitSummary so disabled
// benefits render as zeroes instead of phantom values.
const isBenefitEnabled = (benefitCode: string, categoryName: string) => {
  if (!props.quote?.scheme_categories) return true
  const sc = (props.quote as any).scheme_categories.find(
    (cat: any) => cat.scheme_category === categoryName
  )
  if (!sc) return true
  switch (benefitCode) {
    case 'GLA':
      return sc.gla_benefit === true
    case 'PTD':
      return sc.ptd_benefit === true
    case 'CI':
      return sc.ci_benefit === true
    case 'SGLA':
      return sc.sgla_benefit === true
    case 'PHI':
      return sc.phi_benefit === true
    case 'TTD':
      return sc.ttd_benefit === true
    case 'GFF':
      return sc.family_funeral_benefit === true
    default:
      return true
  }
}

const buildGridData = () => {
  if (!props.resultSummaries || props.resultSummaries.length === 0) return []
  const gridData: any = []

  props.resultSummaries.forEach((rs: any) => {
    const category = rs.category

    const pushRow = (
      benefitCode: string,
      benefitLabel: string,
      totalSA: number | string,
      annualPremium: number,
      cededSA: number,
      reinsurancePremium: number,
      proportion: number
    ) => {
      const enabled = isBenefitEnabled(benefitCode, category)
      gridData.push({
        category,
        benefit: benefitLabel,
        annualSalary: enabled ? rs.total_annual_salary : 0,
        totalSumAssured: totalSA,
        annualPremium,
        totalCededSumAssured: cededSA,
        reinsurancePremium,
        reinsurancePremiumProportion: `${roundUpToTwoDecimalsAccounting(
          (proportion || 0) * 100
        )}%`
      })
    }

    pushRow(
      'GLA',
      glaBenefitTitle.value,
      rs.total_gla_capped_sum_assured,
      rs.exp_total_gla_annual_office_premium,
      rs.total_gla_ceded_sum_assured,
      rs.total_gla_reinsurance_premium,
      rs.gla_reinsurance_premium_proportion
    )
    pushRow(
      'PTD',
      ptdBenefitTitle.value,
      rs.total_ptd_capped_sum_assured,
      rs.exp_total_ptd_annual_office_premium,
      rs.total_ptd_ceded_sum_assured,
      rs.total_ptd_reinsurance_premium,
      rs.ptd_reinsurance_premium_proportion
    )
    pushRow(
      'CI',
      ciBenefitTitle.value,
      rs.total_ci_capped_sum_assured,
      rs.exp_total_ci_annual_office_premium,
      rs.total_ci_ceded_sum_assured,
      rs.total_ci_reinsurance_premium,
      rs.ci_reinsurance_premium_proportion
    )
    pushRow(
      'SGLA',
      sglaBenefitTitle.value,
      rs.total_sgla_capped_sum_assured,
      rs.exp_total_sgla_annual_office_premium,
      rs.total_sgla_ceded_sum_assured,
      rs.total_sgla_reinsurance_premium,
      rs.sgla_reinsurance_premium_proportion
    )
    // TTD / PHI — the "sum assured" column for these is the capped monthly
    // income × 12 conceptually; we surface the capped income field directly
    // (consistent with QuoteBenefitSummary). Ceded is a monthly benefit.
    pushRow(
      'TTD',
      ttdBenefitTitle.value,
      rs.total_ttd_capped_income,
      rs.exp_total_ttd_annual_office_premium,
      rs.total_ttd_ceded_monthly_benefit,
      rs.total_ttd_reinsurance_premium,
      rs.ttd_reinsurance_premium_proportion
    )
    pushRow(
      'PHI',
      phiBenefitTitle.value,
      rs.total_phi_capped_income,
      rs.exp_total_phi_annual_office_premium,
      rs.total_phi_ceded_monthly_benefit,
      rs.total_phi_reinsurance_premium,
      rs.phi_reinsurance_premium_proportion
    )
    pushRow(
      'GFF',
      familyFuneralBenefitTitle.value,
      '',
      rs.exp_total_fun_annual_office_premium,
      rs.total_fun_ceded_sum_assured,
      rs.total_fun_reinsurance_premium,
      rs.fun_reinsurance_premium_proportion
    )
  })

  // Totals row across all categories
  if (props.resultSummaries.length > 0) {
    const totals: any = props.resultSummaries.reduce((acc: any, rs: any) => {
      acc.total_annual_salary =
        (acc.total_annual_salary || 0) + (rs.total_annual_salary || 0)
      acc.total_gla_capped_sum_assured =
        (acc.total_gla_capped_sum_assured || 0) +
        (rs.total_gla_capped_sum_assured || 0)
      acc.total_ptd_capped_sum_assured =
        (acc.total_ptd_capped_sum_assured || 0) +
        (rs.total_ptd_capped_sum_assured || 0)
      acc.total_ci_capped_sum_assured =
        (acc.total_ci_capped_sum_assured || 0) +
        (rs.total_ci_capped_sum_assured || 0)
      acc.total_sgla_capped_sum_assured =
        (acc.total_sgla_capped_sum_assured || 0) +
        (rs.total_sgla_capped_sum_assured || 0)
      acc.total_ttd_capped_income =
        (acc.total_ttd_capped_income || 0) + (rs.total_ttd_capped_income || 0)
      acc.total_phi_capped_income =
        (acc.total_phi_capped_income || 0) + (rs.total_phi_capped_income || 0)

      acc.exp_total_gla_annual_office_premium =
        (acc.exp_total_gla_annual_office_premium || 0) +
        (rs.exp_total_gla_annual_office_premium || 0)
      acc.exp_total_ptd_annual_office_premium =
        (acc.exp_total_ptd_annual_office_premium || 0) +
        (rs.exp_total_ptd_annual_office_premium || 0)
      acc.exp_total_ci_annual_office_premium =
        (acc.exp_total_ci_annual_office_premium || 0) +
        (rs.exp_total_ci_annual_office_premium || 0)
      acc.exp_total_sgla_annual_office_premium =
        (acc.exp_total_sgla_annual_office_premium || 0) +
        (rs.exp_total_sgla_annual_office_premium || 0)
      acc.exp_total_ttd_annual_office_premium =
        (acc.exp_total_ttd_annual_office_premium || 0) +
        (rs.exp_total_ttd_annual_office_premium || 0)
      acc.exp_total_phi_annual_office_premium =
        (acc.exp_total_phi_annual_office_premium || 0) +
        (rs.exp_total_phi_annual_office_premium || 0)
      acc.exp_total_fun_annual_office_premium =
        (acc.exp_total_fun_annual_office_premium || 0) +
        (rs.exp_total_fun_annual_office_premium || 0)

      acc.total_gla_ceded_sum_assured =
        (acc.total_gla_ceded_sum_assured || 0) +
        (rs.total_gla_ceded_sum_assured || 0)
      acc.total_ptd_ceded_sum_assured =
        (acc.total_ptd_ceded_sum_assured || 0) +
        (rs.total_ptd_ceded_sum_assured || 0)
      acc.total_ci_ceded_sum_assured =
        (acc.total_ci_ceded_sum_assured || 0) +
        (rs.total_ci_ceded_sum_assured || 0)
      acc.total_sgla_ceded_sum_assured =
        (acc.total_sgla_ceded_sum_assured || 0) +
        (rs.total_sgla_ceded_sum_assured || 0)
      acc.total_ttd_ceded_monthly_benefit =
        (acc.total_ttd_ceded_monthly_benefit || 0) +
        (rs.total_ttd_ceded_monthly_benefit || 0)
      acc.total_phi_ceded_monthly_benefit =
        (acc.total_phi_ceded_monthly_benefit || 0) +
        (rs.total_phi_ceded_monthly_benefit || 0)
      acc.total_fun_ceded_sum_assured =
        (acc.total_fun_ceded_sum_assured || 0) +
        (rs.total_fun_ceded_sum_assured || 0)

      acc.total_gla_reinsurance_premium =
        (acc.total_gla_reinsurance_premium || 0) +
        (rs.total_gla_reinsurance_premium || 0)
      acc.total_ptd_reinsurance_premium =
        (acc.total_ptd_reinsurance_premium || 0) +
        (rs.total_ptd_reinsurance_premium || 0)
      acc.total_ci_reinsurance_premium =
        (acc.total_ci_reinsurance_premium || 0) +
        (rs.total_ci_reinsurance_premium || 0)
      acc.total_sgla_reinsurance_premium =
        (acc.total_sgla_reinsurance_premium || 0) +
        (rs.total_sgla_reinsurance_premium || 0)
      acc.total_ttd_reinsurance_premium =
        (acc.total_ttd_reinsurance_premium || 0) +
        (rs.total_ttd_reinsurance_premium || 0)
      acc.total_phi_reinsurance_premium =
        (acc.total_phi_reinsurance_premium || 0) +
        (rs.total_phi_reinsurance_premium || 0)
      acc.total_fun_reinsurance_premium =
        (acc.total_fun_reinsurance_premium || 0) +
        (rs.total_fun_reinsurance_premium || 0)
      return acc
    }, {})

    // Recompute proportions from the rolled-up totals so they remain a true
    // ratio of sums, not an average of per-category ratios.
    const ratio = (num: number, den: number) => (den > 0 ? num / den : 0)
    const totalsCategory = 'Totals'
    const pushTotalRow = (
      benefitLabel: string,
      totalSA: number | string,
      annualPremium: number,
      cededSA: number,
      reinsurancePremium: number,
      proportion: number
    ) => {
      gridData.push({
        category: totalsCategory,
        benefit: benefitLabel,
        annualSalary: totals.total_annual_salary,
        totalSumAssured: totalSA,
        annualPremium,
        totalCededSumAssured: cededSA,
        reinsurancePremium,
        reinsurancePremiumProportion: `${roundUpToTwoDecimalsAccounting(
          proportion * 100
        )}%`,
        isSectionHeader: true
      })
    }

    pushTotalRow(
      glaBenefitTitle.value,
      totals.total_gla_capped_sum_assured,
      totals.exp_total_gla_annual_office_premium,
      totals.total_gla_ceded_sum_assured,
      totals.total_gla_reinsurance_premium,
      ratio(
        totals.total_gla_reinsurance_premium,
        totals.exp_total_gla_annual_office_premium
      )
    )
    pushTotalRow(
      ptdBenefitTitle.value,
      totals.total_ptd_capped_sum_assured,
      totals.exp_total_ptd_annual_office_premium,
      totals.total_ptd_ceded_sum_assured,
      totals.total_ptd_reinsurance_premium,
      ratio(
        totals.total_ptd_reinsurance_premium,
        totals.exp_total_ptd_annual_office_premium
      )
    )
    pushTotalRow(
      ciBenefitTitle.value,
      totals.total_ci_capped_sum_assured,
      totals.exp_total_ci_annual_office_premium,
      totals.total_ci_ceded_sum_assured,
      totals.total_ci_reinsurance_premium,
      ratio(
        totals.total_ci_reinsurance_premium,
        totals.exp_total_ci_annual_office_premium
      )
    )
    pushTotalRow(
      sglaBenefitTitle.value,
      totals.total_sgla_capped_sum_assured,
      totals.exp_total_sgla_annual_office_premium,
      totals.total_sgla_ceded_sum_assured,
      totals.total_sgla_reinsurance_premium,
      ratio(
        totals.total_sgla_reinsurance_premium,
        totals.exp_total_sgla_annual_office_premium
      )
    )
    pushTotalRow(
      ttdBenefitTitle.value,
      totals.total_ttd_capped_income,
      totals.exp_total_ttd_annual_office_premium,
      totals.total_ttd_ceded_monthly_benefit,
      totals.total_ttd_reinsurance_premium,
      ratio(
        totals.total_ttd_reinsurance_premium,
        totals.exp_total_ttd_annual_office_premium
      )
    )
    pushTotalRow(
      phiBenefitTitle.value,
      totals.total_phi_capped_income,
      totals.exp_total_phi_annual_office_premium,
      totals.total_phi_ceded_monthly_benefit,
      totals.total_phi_reinsurance_premium,
      ratio(
        totals.total_phi_reinsurance_premium,
        totals.exp_total_phi_annual_office_premium
      )
    )
    pushTotalRow(
      familyFuneralBenefitTitle.value,
      '',
      totals.exp_total_fun_annual_office_premium,
      totals.total_fun_ceded_sum_assured,
      totals.total_fun_reinsurance_premium,
      ratio(
        totals.total_fun_reinsurance_premium,
        totals.exp_total_fun_annual_office_premium
      )
    )
  }

  return gridData
}

// Resolve benefit display names from the benefit map endpoint, same as the
// Premiums Summary does. Falls back to the hard-coded short codes if the API
// call fails so the grid still renders.
onMounted(() => {
  GroupPricingService.getBenefitMaps()
    .then((res: any) => {
      benefitMaps.value = res.data
      const pick = (code: string, fallback: string) => {
        const b: any = (benefitMaps.value as any).find(
          (it: any) => it.benefit_code === code
        )
        if (!b) return fallback
        return b.benefit_alias && b.benefit_alias !== ''
          ? b.benefit_alias
          : b.benefit_name || fallback
      }
      glaBenefitTitle.value = pick('GLA', 'GLA')
      sglaBenefitTitle.value = pick('SGLA', 'SGLA')
      ptdBenefitTitle.value = pick('PTD', 'PTD')
      ciBenefitTitle.value = pick('CI', 'CI')
      phiBenefitTitle.value = pick('PHI', 'PHI')
      ttdBenefitTitle.value = pick('TTD', 'TTD')
      familyFuneralBenefitTitle.value = pick('GFF', 'Group Funeral')
    })
    .catch(() => {
      // Leave hard-coded defaults in place.
    })
})

watch(
  [
    () => props.resultSummaries,
    glaBenefitTitle,
    sglaBenefitTitle,
    ptdBenefitTitle,
    ciBenefitTitle,
    phiBenefitTitle,
    ttdBenefitTitle,
    familyFuneralBenefitTitle
  ],
  () => {
    if (props.resultSummaries && props.resultSummaries.length > 0) {
      rowData.value = buildGridData()
    } else {
      rowData.value = []
    }
  },
  { deep: true, immediate: true }
)
</script>
<style scoped></style>
