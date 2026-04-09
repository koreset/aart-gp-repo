<template>
  <base-card :show-actions="false">
    <template #header>
      <span class="headline">Benefits Summary</span>
    </template>
    <template #default>
      <v-row v-if="rowData.length > 0">
        <v-col>
          <group-pricing-data-grid
            ref="dataGridRef"
            :columnDefs="columnDefs"
            :show-close-button="false"
            :rowData="rowData"
            :table-title="'Premiums Summary By Category'"
            :suppressAutoSize="true"
            :density="'compact'"
            :show-export="true"
          />
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

// import all necessary components and services

// Declare benefit title refs
const glaBenefitTitle = ref('')
const sglaBenefitTitle = ref('')
const ptdBenefitTitle = ref('')
const ciBenefitTitle = ref('')
const phiBenefitTitle = ref('')
const ttdBenefitTitle = ref('')
const familyFuneralBenefitTitle = ref('')
const benefitMaps = ref([])

const props = defineProps({
  resultSummaries: { type: Array, required: true },
  quote: { type: Object, required: false }
})

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
    width: 250,
    minWidth: 200,
    maxWidth: 400,
    flex: 1,
    resizable: true,
    suppressAutoSize: false,
    suppressSizeToFit: false,
    cellStyle: (params) => {
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
    minWidth: 150,
    maxWidth: 220,
    resizable: true,
    suppressAutoSize: false,
    suppressSizeToFit: false,
    valueFormatter: (params) => {
      if (
        params.value === null ||
        params.value === undefined ||
        params.value === ''
      )
        return '-'
      return typeof params.value === 'string'
        ? params.value
        : roundUpToTwoDecimalsAccounting(params.value)
    },
    type: 'rightAligned',
    cellStyle: (params) => {
      if (params.data?.isSubtotal) {
        return { backgroundColor: '#f0f0f0', textAlign: 'right' }
      }
      if (params.data?.isSectionHeader) {
        return { backgroundColor: '#e3f2fd', textAlign: 'right' }
      }
      return { textAlign: 'right' }
    }
  },
  {
    field: 'totalSumAssured',
    headerName: 'Total Sum Assured',
    width: 180,
    minWidth: 150,
    maxWidth: 250,
    resizable: true,
    suppressAutoSize: false,
    suppressSizeToFit: false,
    valueFormatter: (params) => {
      if (
        params.value === null ||
        params.value === undefined ||
        params.value === ''
      )
        return '-'
      return typeof params.value === 'string'
        ? params.value
        : roundUpToTwoDecimalsAccounting(params.value)
    },
    type: 'rightAligned',
    cellStyle: (params) => {
      if (params.data?.isSubtotal) {
        return { backgroundColor: '#f0f0f0', textAlign: 'right' }
      }
      if (params.data?.isSectionHeader) {
        return { backgroundColor: '#e3f2fd', textAlign: 'right' }
      }
      return { textAlign: 'right' }
    }
  },
  {
    field: 'annualPremium',
    headerName: 'Annual Premium',
    width: 160,
    minWidth: 150,
    maxWidth: 220,
    resizable: true,
    suppressAutoSize: false,
    suppressSizeToFit: false,
    valueFormatter: (params) => {
      if (
        params.value === null ||
        params.value === undefined ||
        params.value === ''
      )
        return '-'
      return typeof params.value === 'string'
        ? params.value
        : roundUpToTwoDecimalsAccounting(params.value)
    },
    type: 'rightAligned',
    cellStyle: (params) => {
      if (params.data?.isSubtotal) {
        return { backgroundColor: '#f0f0f0', textAlign: 'right' }
      }
      if (params.data?.isSectionHeader) {
        return { backgroundColor: '#e3f2fd', textAlign: 'right' }
      }
      return { textAlign: 'right' }
    }
  },
  {
    field: 'percentSalary',
    headerName: '% of Salary',
    width: 130,
    minWidth: 120,
    maxWidth: 160,
    resizable: true,
    suppressAutoSize: false,
    suppressSizeToFit: false,
    type: 'rightAligned',
    cellStyle: (params) => {
      if (params.data?.isSubtotal) {
        return { backgroundColor: '#f0f0f0', textAlign: 'right' }
      }
      if (params.data?.isSectionHeader) {
        return { backgroundColor: '#e3f2fd', textAlign: 'right' }
      }
      return { textAlign: 'right' }
    }
  },
  {
    field: 'ratePer1000SA',
    headerName: 'Rate per 1000 SA',
    width: 170,
    minWidth: 150,
    maxWidth: 220,
    resizable: true,
    suppressAutoSize: false,
    suppressSizeToFit: false,
    valueFormatter: (params) => {
      if (
        params.value === null ||
        params.value === undefined ||
        params.value === ''
      )
        return '-'
      return typeof params.value === 'string'
        ? params.value
        : roundUpToTwoDecimalsAccounting(params.value)
    },
    type: 'rightAligned',
    cellStyle: (params) => {
      if (params.data?.isSubtotal) {
        return { backgroundColor: '#f0f0f0', textAlign: 'right' }
      }
      if (params.data?.isSectionHeader) {
        return { backgroundColor: '#e3f2fd', textAlign: 'right' }
      }
      return { textAlign: 'right' }
    }
  }
])
const rowData: any = ref([])

const roundUpToTwoDecimalsAccounting = (num) => {
  const roundedNum = Math.ceil(num * 100) / 100 // Round up to two decimal places
  return roundedNum
    .toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    })
    .replace(/,/g, ' ') // Replace commas with spaces for accounting format }
}

const convertExcelDataToGridData = () => {
  if (!props.resultSummaries || props.resultSummaries.length === 0) return []

  const gridData: any = []

  // Helper function to check if a benefit is enabled
  const isBenefitEnabled = (benefitCode: string, categoryName: string) => {
    if (!props.quote?.scheme_categories) return true // Default to enabled if no quote data

    const schemeCategory = props.quote.scheme_categories.find(
      (cat: any) => cat.scheme_category === categoryName
    )

    if (!schemeCategory) return true // Default to enabled if category not found

    // Check if the benefit is active based on the benefit code
    switch (benefitCode) {
      case 'GLA':
        return schemeCategory.gla_benefit === true
      case 'PTD':
        return schemeCategory.ptd_benefit === true
      case 'CI':
        return schemeCategory.ci_benefit === true
      case 'SGLA':
        return schemeCategory.sgla_benefit === true
      case 'PHI':
        return schemeCategory.phi_benefit === true
      case 'TTD':
        return schemeCategory.ttd_benefit === true
      default:
        return true
    }
  }

  props.resultSummaries.forEach((resultSummary: any) => {
    const category = resultSummary.category

    // Add benefit rows
    gridData.push({
      category,
      benefit: glaBenefitTitle.value,
      annualSalary: isBenefitEnabled('GLA', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_gla_capped_sum_assured,
      annualPremium: resultSummary.exp_total_gla_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(resultSummary.exp_proportion_gla_office_premium_salary * 100)}%`,
      ratePer1000SA: resultSummary.exp_gla_office_rate_per_1000_sa
    })

    gridData.push({
      category,
      benefit: ptdBenefitTitle.value,
      annualSalary: isBenefitEnabled('PTD', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_ptd_capped_sum_assured,
      annualPremium: resultSummary.exp_total_ptd_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(resultSummary.exp_proportion_ptd_office_premium_salary * 100)}%`,
      ratePer1000SA: resultSummary.exp_ptd_office_rate_per_1000_sa
    })

    gridData.push({
      category,
      benefit: ciBenefitTitle.value,
      annualSalary: isBenefitEnabled('CI', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_ci_capped_sum_assured,
      annualPremium: resultSummary.exp_total_ci_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(resultSummary.exp_proportion_ci_office_premium_salary * 100)}%`,
      ratePer1000SA: resultSummary.exp_ci_office_rate_per_1000_sa
    })

    gridData.push({
      category,
      benefit: sglaBenefitTitle.value,
      annualSalary: isBenefitEnabled('SGLA', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_sgla_capped_sum_assured,
      annualPremium: resultSummary.exp_total_sgla_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(resultSummary.exp_proportion_sgla_office_premium_salary * 100)}%`,
      ratePer1000SA: resultSummary.exp_sgla_office_rate_per_1000_sa
    })

    gridData.push({
      category,
      benefit: phiBenefitTitle.value,
      annualSalary: isBenefitEnabled('PHI', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_phi_capped_income,
      annualPremium: resultSummary.exp_total_phi_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(resultSummary.exp_proportion_phi_office_premium_salary * 100)}%`,
      ratePer1000SA: resultSummary.exp_phi_office_rate_per_1000_sa
    })

    gridData.push({
      category,
      benefit: ttdBenefitTitle.value,
      annualSalary: isBenefitEnabled('TTD', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_ttd_capped_income,
      annualPremium: resultSummary.exp_total_ttd_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(resultSummary.exp_proportion_ttd_office_premium_salary * 100)}%`,
      ratePer1000SA: resultSummary.exp_ttd_office_rate_per_1000_sa
    })

    // Add subtotal row
    const anyBenefitEnabled = ['GLA', 'PTD', 'CI', 'SGLA', 'PHI', 'TTD'].some(
      (benefitCode) => isBenefitEnabled(benefitCode, category)
    )
    gridData.push({
      category,
      benefit: 'Sub Total (Excl. Funeral)',
      annualSalary: anyBenefitEnabled ? resultSummary.total_annual_salary : 0,
      totalSumAssured: resultSummary.total_gla_capped_sum_assured,
      annualPremium: resultSummary.exp_total_annual_premium_excl_funeral,
      percentSalary: `${roundUpToTwoDecimalsAccounting(resultSummary.proportion_exp_total_premium_excl_funeral_salary * 100)}%`,
      ratePer1000SA: '',
      isSubtotal: true
    })

    gridData.push({
      category,
      benefit: 'Group Funeral Annual Premium per Member',
      annualSalary: '',
      totalSumAssured: '',
      annualPremium: resultSummary.exp_total_fun_annual_premium_per_member,
      percentSalary: '',
      ratePer1000SA: ''
    })

    gridData.push({
      category,
      benefit: 'Group Funeral Annual Premium',
      annualSalary: '',
      totalSumAssured: '',
      annualPremium: resultSummary.exp_total_fun_annual_office_premium,
      percentSalary: '',
      ratePer1000SA: ''
    })
  })

  // Calculate totals across all categories
  if (props.resultSummaries.length > 0) {
    const totals: any = props.resultSummaries.reduce(
      (acc: any, resultSummary: any) => {
        return {
          total_gla_capped_sum_assured:
            (acc.total_gla_capped_sum_assured || 0) +
            (resultSummary.total_gla_capped_sum_assured || 0),
          exp_total_gla_annual_office_premium:
            (acc.exp_total_gla_annual_office_premium || 0) +
            (resultSummary.exp_total_gla_annual_office_premium || 0),
          exp_proportion_gla_office_premium_salary:
            (acc.exp_proportion_gla_office_premium_salary || 0) +
            (resultSummary.exp_proportion_gla_office_premium_salary || 0),

          total_ptd_capped_sum_assured:
            (acc.total_ptd_capped_sum_assured || 0) +
            (resultSummary.total_ptd_capped_sum_assured || 0),
          exp_total_ptd_annual_office_premium:
            (acc.exp_total_ptd_annual_office_premium || 0) +
            (resultSummary.exp_total_ptd_annual_office_premium || 0),
          exp_proportion_ptd_office_premium_salary:
            (acc.exp_proportion_ptd_office_premium_salary || 0) +
            (resultSummary.exp_proportion_ptd_office_premium_salary || 0),

          total_ci_capped_sum_assured:
            (acc.total_ci_capped_sum_assured || 0) +
            (resultSummary.total_ci_capped_sum_assured || 0),
          exp_total_ci_annual_office_premium:
            (acc.exp_total_ci_annual_office_premium || 0) +
            (resultSummary.exp_total_ci_annual_office_premium || 0),
          exp_proportion_ci_office_premium_salary:
            (acc.exp_proportion_ci_office_premium_salary || 0) +
            (resultSummary.exp_proportion_ci_office_premium_salary || 0),

          total_sgla_capped_sum_assured:
            (acc.total_sgla_capped_sum_assured || 0) +
            (resultSummary.total_sgla_capped_sum_assured || 0),
          exp_total_sgla_annual_office_premium:
            (acc.exp_total_sgla_annual_office_premium || 0) +
            (resultSummary.exp_total_sgla_annual_office_premium || 0),
          exp_proportion_sgla_office_premium_salary:
            (acc.exp_proportion_sgla_office_premium_salary || 0) +
            (resultSummary.exp_proportion_sgla_office_premium_salary || 0),

          total_phi_capped_income:
            (acc.total_phi_capped_income || 0) +
            (resultSummary.total_phi_capped_income || 0),
          exp_total_phi_annual_office_premium:
            (acc.exp_total_phi_annual_office_premium || 0) +
            (resultSummary.exp_total_phi_annual_office_premium || 0),
          exp_proportion_phi_office_premium_salary:
            (acc.exp_proportion_phi_office_premium_salary || 0) +
            (resultSummary.exp_proportion_phi_office_premium_salary || 0),

          total_ttd_capped_income:
            (acc.total_ttd_capped_income || 0) +
            (resultSummary.total_ttd_capped_income || 0),
          exp_total_ttd_annual_office_premium:
            (acc.exp_total_ttd_annual_office_premium || 0) +
            (resultSummary.exp_total_ttd_annual_office_premium || 0),
          exp_proportion_ttd_office_premium_salary:
            (acc.exp_proportion_ttd_office_premium_salary || 0) +
            (resultSummary.exp_proportion_ttd_office_premium_salary || 0),

          exp_total_annual_premium_excl_funeral:
            (acc.exp_total_annual_premium_excl_funeral || 0) +
            (resultSummary.exp_total_annual_premium_excl_funeral || 0),
          proportion_exp_total_premium_excl_funeral_salary:
            (acc.proportion_exp_total_premium_excl_funeral_salary || 0) +
            (resultSummary.proportion_exp_total_premium_excl_funeral_salary ||
              0),

          exp_total_fun_monthly_premium_per_member:
            (acc.exp_total_fun_monthly_premium_per_member || 0) +
            (resultSummary.exp_total_fun_monthly_premium_per_member || 0),
          exp_total_fun_annual_premium_per_member:
            (acc.exp_total_fun_annual_premium_per_member || 0) +
            (resultSummary.exp_total_fun_annual_premium_per_member || 0),
          exp_total_fun_annual_office_premium:
            (acc.exp_total_fun_annual_office_premium || 0) +
            (resultSummary.exp_total_fun_annual_office_premium || 0),

          total_annual_salary:
            (acc.total_annual_salary || 0) +
            (resultSummary.total_annual_salary || 0)
        }
      },
      {}
    )

    // Add totals category rows
    const totalsCategory = 'Totals'

    gridData.push({
      category: totalsCategory,
      benefit: glaBenefitTitle.value,
      annualSalary: totals.total_annual_salary,
      totalSumAssured: totals.total_gla_capped_sum_assured,
      annualPremium: totals.exp_total_gla_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(
        (totals.exp_total_gla_annual_office_premium /
          totals.total_annual_salary || 0) * 100
      )}%`,
      ratePer1000SA: totals.total_gla_capped_sum_assured
        ? (totals.exp_total_gla_annual_office_premium * 1000) /
          totals.total_gla_capped_sum_assured
        : ''
    })

    gridData.push({
      category: totalsCategory,
      benefit: ptdBenefitTitle.value,
      annualSalary: totals.total_annual_salary,
      totalSumAssured: totals.total_ptd_capped_sum_assured,
      annualPremium: totals.exp_total_ptd_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting((totals.exp_total_ptd_annual_office_premium / totals.total_annual_salary || 0) * 100)}%`,
      ratePer1000SA: totals.total_ptd_capped_sum_assured
        ? (totals.exp_total_ptd_annual_office_premium * 1000) /
          totals.total_ptd_capped_sum_assured
        : ''
    })

    gridData.push({
      category: totalsCategory,
      benefit: ciBenefitTitle.value,
      annualSalary: totals.total_annual_salary,
      totalSumAssured: totals.total_ci_capped_sum_assured,
      annualPremium: totals.exp_total_ci_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(
        (totals.exp_total_ci_annual_office_premium /
          totals.total_annual_salary || 0) * 100
      )}%`,
      ratePer1000SA: totals.total_ci_capped_sum_assured
        ? (totals.exp_total_ci_annual_office_premium * 1000) /
          totals.total_ci_capped_sum_assured
        : ''
    })

    gridData.push({
      category: totalsCategory,
      benefit: sglaBenefitTitle.value,
      annualSalary: totals.total_annual_salary,
      totalSumAssured: totals.total_sgla_capped_sum_assured,
      annualPremium: totals.exp_total_sgla_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(
        (totals.exp_total_sgla_annual_office_premium /
          totals.total_annual_salary || 0) * 100
      )}%`,
      ratePer1000SA: totals.total_sgla_capped_sum_assured
        ? (totals.exp_total_sgla_annual_office_premium * 1000) /
          totals.total_sgla_capped_sum_assured
        : ''
    })

    gridData.push({
      category: totalsCategory,
      benefit: phiBenefitTitle.value,
      annualSalary: totals.total_annual_salary,
      totalSumAssured: totals.total_phi_capped_income,
      annualPremium: totals.exp_total_phi_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(
        (totals.exp_total_phi_annual_office_premium /
          totals.total_annual_salary || 0) * 100
      )}%`,
      ratePer1000SA: totals.total_phi_capped_income
        ? (totals.exp_total_phi_annual_office_premium * 1000) /
          totals.total_phi_capped_income
        : ''
    })

    gridData.push({
      category: totalsCategory,
      benefit: ttdBenefitTitle.value,
      annualSalary: totals.total_annual_salary,
      totalSumAssured: totals.total_ttd_capped_income,
      annualPremium: totals.exp_total_ttd_annual_office_premium,
      percentSalary: `${roundUpToTwoDecimalsAccounting(
        (totals.exp_total_ttd_annual_office_premium /
          totals.total_annual_salary || 0) * 100
      )}%`,
      ratePer1000SA: totals.total_ttd_capped_income
        ? (totals.exp_total_ttd_annual_office_premium * 1000) /
          totals.total_ttd_capped_income
        : ''
    })

    // Add totals subtotal row
    // Calculate total annual salary only for categories where any benefit is enabled
    const totalAnnualSalaryForEnabledBenefits = props.resultSummaries.reduce(
      (acc: number, resultSummary: any) => {
        const categoryHasEnabledBenefits = [
          'GLA',
          'PTD',
          'CI',
          'SGLA',
          'PHI',
          'TTD'
        ].some((benefitCode) =>
          isBenefitEnabled(benefitCode, resultSummary.category)
        )
        return (
          acc +
          (categoryHasEnabledBenefits ? resultSummary.total_annual_salary : 0)
        )
      },
      0
    )

    gridData.push({
      category: totalsCategory,
      benefit: 'Sub Total (Excl. Funeral)',
      annualSalary: totalAnnualSalaryForEnabledBenefits,
      totalSumAssured: totals.total_gla_capped_sum_assured,
      annualPremium: totals.exp_total_annual_premium_excl_funeral,
      percentSalary: `${roundUpToTwoDecimalsAccounting(
        (totals.exp_total_annual_premium_excl_funeral /
          totalAnnualSalaryForEnabledBenefits || 0) * 100
      )}%`,
      ratePer1000SA: '',
      isSubtotal: true
    })

    gridData.push({
      category: totalsCategory,
      benefit: 'Group Funeral Annual Premium per Member',
      annualSalary: '',
      totalSumAssured: '',
      annualPremium: totals.exp_total_fun_annual_premium_per_member,
      percentSalary: '',
      ratePer1000SA: ''
    })

    gridData.push({
      category: totalsCategory,
      benefit: 'Group Funeral Annual Premium',
      annualSalary: '',
      totalSumAssured: '',
      annualPremium: totals.exp_total_fun_annual_office_premium,
      percentSalary: '',
      ratePer1000SA: ''
    })
  }

  return gridData
}

onMounted(() => {
  GroupPricingService.getBenefitMaps().then((res) => {
    benefitMaps.value = res.data
    const glaBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'GLA'
    )
    if (glaBenefit.benefit_alias !== '') {
      glaBenefitTitle.value = glaBenefit.benefit_alias
    } else {
      glaBenefitTitle.value = glaBenefit.benefit_name
    }
    const sglaBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'SGLA'
    )
    if (sglaBenefit.benefit_alias !== '') {
      sglaBenefitTitle.value = sglaBenefit.benefit_alias
    } else {
      sglaBenefitTitle.value = sglaBenefit.benefit_name
    }
    const ptdBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'PTD'
    )
    if (ptdBenefit.benefit_alias !== '') {
      ptdBenefitTitle.value = ptdBenefit.benefit_alias
    } else {
      ptdBenefitTitle.value = ptdBenefit.benefit_name
    }
    const ciBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'CI'
    )
    if (ciBenefit.benefit_alias !== '') {
      ciBenefitTitle.value = ciBenefit.benefit_alias
    } else {
      ciBenefitTitle.value = ciBenefit.benefit_name
    }
    const phiBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'PHI'
    )
    if (phiBenefit.benefit_alias !== '') {
      phiBenefitTitle.value = phiBenefit.benefit_alias
    } else {
      phiBenefitTitle.value = phiBenefit.benefit_name
    }
    const ttdBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'TTD'
    )
    if (ttdBenefit.benefit_alias !== '') {
      ttdBenefitTitle.value = ttdBenefit.benefit_alias
    } else {
      ttdBenefitTitle.value = ttdBenefit.benefit_name
    }
    const familyFuneralBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'GFF'
    )
    if (familyFuneralBenefit.benefit_alias !== '') {
      familyFuneralBenefitTitle.value = familyFuneralBenefit.benefit_alias
    } else {
      familyFuneralBenefitTitle.value = familyFuneralBenefit.benefit_name
    }
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
    ttdBenefitTitle
  ],
  () => {
    if (props.resultSummaries && props.resultSummaries.length > 0) {
      rowData.value = convertExcelDataToGridData()
    } else {
      rowData.value = []
    }
  },
  { deep: true, immediate: true }
)
</script>
<style scoped></style>
