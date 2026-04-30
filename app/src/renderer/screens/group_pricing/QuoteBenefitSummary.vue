<template>
  <base-card :show-actions="false">
    <template #header>
      <span class="headline">Benefits Summary</span>
    </template>
    <template #default>
      <v-tabs
        v-model="benefitsSubTab"
        color="primary"
        density="compact"
        show-arrows
      >
        <v-tab value="summary">Summary</v-tab>
        <v-tab value="breakdown">Breakdown</v-tab>
      </v-tabs>
      <v-window v-model="benefitsSubTab" class="mt-3">
        <v-window-item value="summary">
          <v-row v-if="rowData.length > 0">
            <v-col>
              <group-pricing-data-grid
                ref="dataGridRef"
                :key="`premium-grid-${rowDataKey}`"
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
        </v-window-item>
        <v-window-item value="breakdown">
          <v-row v-if="breakdownRowData.length > 0">
            <v-col>
              <group-pricing-data-grid
                :key="`breakdown-grid-${breakdownRowDataKey}`"
                :columnDefs="breakdownColumnDefs"
                :show-close-button="false"
                :rowData="breakdownRowData"
                :table-title="'Premium Breakdown By Benefit'"
                :suppressAutoSize="true"
                :density="'compact'"
                :show-export="true"
              />
            </v-col>
          </v-row>
        </v-window-item>
      </v-window>
    </template>
  </base-card>

  <base-card
    v-if="extendedFamilyCategories.length > 0"
    :show-actions="false"
    class="mt-4"
  >
    <template #header>
      <div class="d-flex align-center" style="width: 100%">
        <span class="headline">Extended Family Funeral Summary</span>
        <v-spacer />
        <v-btn
          v-if="hasAnyExtendedFamilyRates"
          size="small"
          variant="flat"
          color="white"
          class="text-primary"
          prepend-icon="mdi-download"
          @click="exportExtendedFamilyToExcel"
        >
          Export All
        </v-btn>
      </div>
    </template>
    <template #default>
      <v-tabs
        v-model="extendedFamilyActiveTab"
        color="primary"
        density="compact"
        show-arrows
      >
        <v-tab
          v-for="(rs, i) in extendedFamilyCategories"
          :key="'ef-tab-' + i"
          :value="(rs as any).category"
        >
          {{ (rs as any).category }}
        </v-tab>
      </v-tabs>
      <v-window v-model="extendedFamilyActiveTab" class="mt-3">
        <v-window-item
          v-for="(rs, i) in extendedFamilyCategories"
          :key="'ef-pane-' + i"
          :value="(rs as any).category"
        >
          <template
            v-if="((rs as any).extended_family_band_rates?.length ?? 0) > 0"
          >
            <div class="text-caption mb-2">
              Pricing method:
              {{
                (rs as any).extended_family_pricing_method === 'sum_assured'
                  ? 'Sum Assured per Band'
                  : 'Rate per 1,000'
              }}
              <span v-if="(rs as any).extended_family_age_band_type">
                &middot; Age bands:
                {{ (rs as any).extended_family_age_band_type }}
              </span>
            </div>
            <v-table density="compact">
              <thead>
                <tr>
                  <th>Age Band</th>
                  <th class="text-right">Sum Assured</th>
                  <th class="text-right">Rate per 1,000</th>
                  <th class="text-right">
                    Monthly Premium<br />(per Ext. Family Member)
                  </th>
                  <th class="text-right">
                    Annual Premium<br />(per Ext. Family Member)
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(r, j) in (rs as any).extended_family_band_rates"
                  :key="'ef-row-' + i + '-' + j"
                >
                  <td>{{ formatBandLabel(r) }}</td>
                  <td class="text-right">
                    {{
                      (rs as any).extended_family_pricing_method ===
                      'sum_assured'
                        ? roundUpToTwoDecimalsAccounting(r.sum_assured ?? 0)
                        : '—'
                    }}
                  </td>
                  <td class="text-right">
                    {{
                      (rs as any).extended_family_pricing_method ===
                      'rate_per_1000'
                        ? roundUpToTwoDecimalsAccounting(
                            (r.office_rate ?? 0) * 1000
                          )
                        : '—'
                    }}
                  </td>
                  <td class="text-right">
                    {{
                      roundUpToTwoDecimalsAccounting(
                        r.office_monthly_premium ?? 0
                      )
                    }}
                  </td>
                  <td class="text-right">
                    {{
                      roundUpToTwoDecimalsAccounting(
                        (r.office_monthly_premium ?? 0) * 12
                      )
                    }}
                  </td>
                </tr>
              </tbody>
            </v-table>
          </template>
          <v-alert v-else type="warning" density="compact" variant="tonal">
            Extended family benefit is enabled for
            <strong>{{ (rs as any).category }}</strong>
            but per-band rates could not be computed. Check that (i) age bands
            exist for the configured band type<span
              v-if="(rs as any).extended_family_age_band_type"
            >
              (<em>{{ (rs as any).extended_family_age_band_type }}</em
              >)</span
            >, and (ii) funeral rates exist for the scheme&rsquo;s risk rate
            code. Re-run the quote calculation after fixing the configuration.
          </v-alert>
        </v-window-item>
      </v-window>
    </template>
  </base-card>

  <base-card
    v-if="additionalGlaCoverCategories.length > 0"
    :show-actions="false"
    class="mt-4"
  >
    <template #header>
      <div class="d-flex align-center" style="width: 100%">
        <span class="headline"
          >{{ additionalGlaCoverBenefitTitle }} Summary</span
        >
        <v-spacer />
        <v-btn
          v-if="hasAnyAdditionalGlaCoverRates"
          size="small"
          variant="flat"
          color="white"
          class="text-primary"
          prepend-icon="mdi-download"
          @click="exportAdditionalGlaCoverToExcel"
        >
          Export All
        </v-btn>
      </div>
    </template>
    <template #default>
      <v-tabs
        v-model="additionalGlaCoverActiveTab"
        color="primary"
        density="compact"
        show-arrows
      >
        <v-tab
          v-for="(rs, i) in additionalGlaCoverCategories"
          :key="'agla-tab-' + i"
          :value="(rs as any).category"
        >
          {{ (rs as any).category }}
        </v-tab>
      </v-tabs>
      <v-window v-model="additionalGlaCoverActiveTab" class="mt-3">
        <v-window-item
          v-for="(rs, i) in additionalGlaCoverCategories"
          :key="'agla-pane-' + i"
          :value="(rs as any).category"
        >
          <template
            v-if="
              ((rs as any).additional_gla_cover_band_rates?.length ?? 0) > 0
            "
          >
            <div class="text-caption mb-2">
              <span v-if="(rs as any).additional_gla_cover_age_band_type">
                Age bands:
                <em>{{ (rs as any).additional_gla_cover_age_band_type }}</em>
              </span>
              <span
                v-if="(rs as any).additional_gla_cover_male_prop_used != null"
              >
                &middot; Male proportion used:
                {{
                  (
                    ((rs as any).additional_gla_cover_male_prop_used ?? 0) * 100
                  ).toFixed(1)
                }}%
              </span>
            </div>
            <div style="overflow-x: auto">
              <v-table density="compact" class="agla-summary-table">
                <thead>
                  <tr>
                    <th rowspan="2" class="agla-age-band">Age Band</th>
                    <th colspan="3" class="text-center agla-group-divider">
                      Risk Rate / 1,000
                    </th>
                    <th colspan="3" class="text-center agla-group-divider">
                      Binder Fee / 1,000
                    </th>
                    <th colspan="3" class="text-center agla-group-divider">
                      Outsource Fee / 1,000
                    </th>
                    <th colspan="3" class="text-center agla-group-divider">
                      Commission / 1,000
                    </th>
                    <th colspan="3" class="text-center">Office Rate / 1,000</th>
                  </tr>
                  <tr>
                    <th class="text-right agla-mf-cell">M</th>
                    <th class="text-right agla-mf-cell">F</th>
                    <th
                      class="text-right agla-combined-cell agla-group-divider"
                    >
                      Combined
                    </th>
                    <th class="text-right agla-mf-cell">M</th>
                    <th class="text-right agla-mf-cell">F</th>
                    <th
                      class="text-right agla-combined-cell agla-group-divider"
                    >
                      Combined
                    </th>
                    <th class="text-right agla-mf-cell">M</th>
                    <th class="text-right agla-mf-cell">F</th>
                    <th
                      class="text-right agla-combined-cell agla-group-divider"
                    >
                      Combined
                    </th>
                    <th class="text-right agla-mf-cell">M</th>
                    <th class="text-right agla-mf-cell">F</th>
                    <th
                      class="text-right agla-combined-cell agla-group-divider"
                    >
                      Combined
                    </th>
                    <th class="text-right agla-mf-cell">M</th>
                    <th class="text-right agla-mf-cell">F</th>
                    <th class="text-right agla-combined-cell">Combined</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(r, j) in (rs as any)
                      .additional_gla_cover_band_rates"
                    :key="'agla-row-' + i + '-' + j"
                  >
                    <td class="agla-age-band">{{ formatBandLabel(r) }}</td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.risk_rate_per1000_male ?? 0
                        )
                      }}
                    </td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.risk_rate_per1000_female ?? 0
                        )
                      }}
                    </td>
                    <td
                      class="text-right agla-combined-cell agla-group-divider"
                    >
                      {{ roundUpToTwoDecimalsAccounting(r.risk_rate_per1000) }}
                    </td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.binder_fee_per1000_male ?? 0
                        )
                      }}
                    </td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.binder_fee_per1000_female ?? 0
                        )
                      }}
                    </td>
                    <td
                      class="text-right agla-combined-cell agla-group-divider"
                    >
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.binder_fee_per1000 ?? 0
                        )
                      }}
                    </td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.outsource_fee_per1000_male ?? 0
                        )
                      }}
                    </td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.outsource_fee_per1000_female ?? 0
                        )
                      }}
                    </td>
                    <td
                      class="text-right agla-combined-cell agla-group-divider"
                    >
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.outsource_fee_per1000 ?? 0
                        )
                      }}
                    </td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.commission_per1000_male ?? 0
                        )
                      }}
                    </td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.commission_per1000_female ?? 0
                        )
                      }}
                    </td>
                    <td
                      class="text-right agla-combined-cell agla-group-divider"
                    >
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.commission_per1000 ?? 0
                        )
                      }}
                    </td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.office_rate_per1000_male ?? 0
                        )
                      }}
                    </td>
                    <td class="text-right agla-mf-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(
                          r.office_rate_per1000_female ?? 0
                        )
                      }}
                    </td>
                    <td class="text-right agla-combined-cell">
                      {{
                        roundUpToTwoDecimalsAccounting(r.office_rate_per1000)
                      }}
                    </td>
                  </tr>
                </tbody>
              </v-table>
            </div>
          </template>
          <v-alert v-else type="warning" density="compact" variant="tonal">
            {{ additionalGlaCoverBenefitTitle }} is enabled for
            <strong>{{ (rs as any).category }}</strong> but per-band rates could
            not be computed. Check that (i) age bands exist for the configured
            band type<span
              v-if="(rs as any).additional_gla_cover_age_band_type"
            >
              (<em>{{ (rs as any).additional_gla_cover_age_band_type }}</em
              >)</span
            >, and (ii) GLA rates exist for the scheme&rsquo;s risk rate code,
            waiting period and main GLA benefit type. Re-run the quote
            calculation after fixing the configuration.
          </v-alert>
        </v-window-item>
      </v-window>
    </template>
  </base-card>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import * as XLSX from 'xlsx'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import GroupPricingDataGrid from '@/renderer/components/tables/GroupPricingDataGrid.vue'
import {
  computeOfficePremium,
  officeRateFromRiskRate,
  officeProportionFromRiskProportion
} from '@/renderer/utils/quoteDataHelpers'

// import all necessary components and services

// Declare benefit title refs
const glaBenefitTitle = ref('')
const sglaBenefitTitle = ref('')
const ptdBenefitTitle = ref('')
const ciBenefitTitle = ref('')
const phiBenefitTitle = ref('')
const ttdBenefitTitle = ref('')
const familyFuneralBenefitTitle = ref('')
const additionalAccidentalGlaBenefitTitle = ref('Additional Accidental GLA')
const additionalGlaCoverBenefitTitle = ref('Additional GLA Cover')
const glaEducatorBenefitTitle = ref('GLA Educator')
const ptdEducatorBenefitTitle = ref('PTD Educator')
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
    field: 'finalAnnualPremium',
    headerName: 'Annual Premium',
    width: 170,
    minWidth: 150,
    maxWidth: 220,
    resizable: true,
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
    field: 'finalAnnualCommission',
    headerName: 'Commission',
    width: 160,
    minWidth: 140,
    maxWidth: 200,
    resizable: true,
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
// Bumped every time rowData is rebuilt so the grid child remounts. Force a
// clean AG Grid instance per rowData generation — without this, swapping
// rowData once it's been populated (e.g. benefit titles load after the first
// watcher fire, or a sibling tab becomes active and triggers a re-layout)
// leaves % of Salary / Rate per 1000 SA cells blank because AG Grid can't
// reconcile rows across the swap without a getRowId.
const rowDataKey = ref(0)

// Sub-tab inside the Benefits Summary card: 'summary' (existing per-benefit
// premiums grid) or 'breakdown' (new per-benefit decomposition into base
// premium + binder fee + outsourcing fee + commission, both pre- and
// post-discount).
const benefitsSubTab = ref<'summary' | 'breakdown'>('summary')

// Column defs for the Premium Breakdown grid. Pairs of (pre-discount, final)
// for each fee component so the eye can compare across columns. Final*Binder
// and Final*Outsourced amounts are persisted by recomputeFinalPremiumsAndCommission
// (post-discount they shrink with the office premium); on non-binder
// distribution channels both rates are 0 so those columns stay 0.
const moneyCellStyle = (params: any) => {
  if (params.data?.isSubtotal) {
    return { backgroundColor: '#f0f0f0', textAlign: 'right' }
  }
  if (params.data?.isSectionHeader) {
    return { backgroundColor: '#e3f2fd', textAlign: 'right' }
  }
  return { textAlign: 'right' }
}
const moneyValueFormatter = (params: any) => {
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
const breakdownColumnDefs: any = ref([
  { field: 'category', headerName: 'Category', rowGroup: true, hide: true },
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
    headerName: 'Before Discount',
    headerClass: 'breakdown-group-header breakdown-group-pre',
    children: [
      {
        field: 'expRiskPremium',
        headerName: 'Risk Premium',
        width: 150,
        minWidth: 130,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      },
      {
        field: 'expBasePremium',
        headerName: 'Base Premium',
        width: 150,
        minWidth: 130,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      },
      {
        field: 'expBinderFee',
        headerName: 'Binder Fee',
        width: 140,
        minWidth: 120,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      },
      {
        field: 'expOutsourcingFee',
        headerName: 'Outsourcing Fee',
        width: 160,
        minWidth: 140,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      },
      {
        field: 'expCommission',
        headerName: 'Commission',
        width: 150,
        minWidth: 130,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      },
      {
        field: 'expGrossPremium',
        headerName: 'Gross Premium',
        width: 160,
        minWidth: 140,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      }
    ]
  },
  {
    headerName: 'After Discount',
    headerClass: 'breakdown-group-header breakdown-group-post',
    children: [
      {
        field: 'finalBasePremium',
        headerName: 'Base Premium',
        width: 150,
        minWidth: 130,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      },
      {
        field: 'finalBinderFee',
        headerName: 'Binder Fee',
        width: 140,
        minWidth: 120,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      },
      {
        field: 'finalOutsourcingFee',
        headerName: 'Outsourcing Fee',
        width: 160,
        minWidth: 140,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      },
      {
        field: 'finalCommission',
        headerName: 'Commission',
        width: 150,
        minWidth: 130,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      },
      {
        field: 'finalPremium',
        headerName: 'Final Premium',
        width: 160,
        minWidth: 140,
        resizable: true,
        type: 'rightAligned',
        valueFormatter: moneyValueFormatter,
        cellStyle: moneyCellStyle
      }
    ]
  }
])
const breakdownRowData: any = ref([])
const breakdownRowDataKey = ref(0)

const roundUpToTwoDecimalsAccounting = (num) => {
  const roundedNum = Math.ceil(num * 100) / 100 // Round up to two decimal places
  return roundedNum
    .toLocaleString('en-US', {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    })
    .replace(/,/g, ' ') // Replace commas with spaces for accounting format }
}

// Labels the extended-family age-band row. Bands with max_age >= 150 are
// treated as open-ended (e.g. "70+") to match the configuration convention.
const formatBandLabel = (r: { min_age: number; max_age: number }) => {
  if (r.max_age >= 150) return `${r.min_age}+`
  return `${r.min_age}\u2013${r.max_age}`
}

// Categories with the extended-family benefit toggled on. Includes categories
// whose band rates could not be computed so the tab surfaces a warning instead
// of silently disappearing.
const extendedFamilyCategories = computed<any[]>(() =>
  (props.resultSummaries ?? []).filter(
    (rs: any) => !!rs?.extended_family_benefit
  )
)

const hasAnyExtendedFamilyRates = computed<boolean>(() =>
  extendedFamilyCategories.value.some(
    (rs: any) => (rs?.extended_family_band_rates?.length ?? 0) > 0
  )
)

const extendedFamilyActiveTab = ref<string>('')
watch(
  extendedFamilyCategories,
  (cats) => {
    if (cats.length === 0) {
      extendedFamilyActiveTab.value = ''
      return
    }
    const stillValid = cats.some(
      (rs: any) => rs?.category === extendedFamilyActiveTab.value
    )
    if (!stillValid) {
      extendedFamilyActiveTab.value = cats[0]?.category ?? ''
    }
  },
  { immediate: true }
)

// Additional GLA Cover — config + per-category band rates are mirrored onto
// each resultSummary (MRRS) row at calc time, matching the Extended Family
// pattern. Accepts bool `true`, number `1`, or string `"1"`/`"true"` so any
// DB/driver serialization quirks around BIT / TINYINT(1) don't hide the
// section.
const additionalGlaCoverCategories = computed<any[]>(() =>
  (props.resultSummaries ?? []).filter((rs: any) => {
    const v = rs?.additional_gla_cover_benefit
    return v === true || v === 1 || v === '1' || v === 'true'
  })
)

const additionalGlaCoverActiveTab = ref<string>('')
watch(
  additionalGlaCoverCategories,
  (cats) => {
    if (cats.length === 0) {
      additionalGlaCoverActiveTab.value = ''
      return
    }
    const stillValid = cats.some(
      (rs: any) => rs?.category === additionalGlaCoverActiveTab.value
    )
    if (!stillValid) {
      additionalGlaCoverActiveTab.value = cats[0]?.category ?? ''
    }
  },
  { immediate: true }
)

const hasAnyAdditionalGlaCoverRates = computed<boolean>(() =>
  additionalGlaCoverCategories.value.some(
    (rs: any) => (rs?.additional_gla_cover_band_rates?.length ?? 0) > 0
  )
)

const exportAdditionalGlaCoverToExcel = () => {
  const cats = additionalGlaCoverCategories.value.filter(
    (rs: any) => (rs?.additional_gla_cover_band_rates?.length ?? 0) > 0
  )
  if (cats.length === 0) return

  const wb = XLSX.utils.book_new()
  cats.forEach((rs: any) => {
    const aoa: any[][] = []
    aoa.push([`${additionalGlaCoverBenefitTitle.value} — ${rs.category}`])
    const headerBits: string[] = []
    if (rs.additional_gla_cover_age_band_type) {
      headerBits.push(`Age bands: ${rs.additional_gla_cover_age_band_type}`)
    }
    if (rs.additional_gla_cover_male_prop_used != null) {
      headerBits.push(
        `Male proportion used: ${(
          (rs.additional_gla_cover_male_prop_used ?? 0) * 100
        ).toFixed(1)}%`
      )
    }
    if (headerBits.length > 0) {
      aoa.push([headerBits.join('  ·  ')])
    }
    aoa.push([])
    const groupHeaderRowIdx = aoa.length
    aoa.push([
      'Age Band',
      'Risk Rate / 1,000',
      '',
      '',
      'Binder Fee / 1,000',
      '',
      '',
      'Outsource Fee / 1,000',
      '',
      '',
      'Commission / 1,000',
      '',
      '',
      'Office Rate / 1,000',
      '',
      ''
    ])
    aoa.push([
      '',
      'M',
      'F',
      'Combined',
      'M',
      'F',
      'Combined',
      'M',
      'F',
      'Combined',
      'M',
      'F',
      'Combined',
      'M',
      'F',
      'Combined'
    ])
    ;(rs.additional_gla_cover_band_rates ?? []).forEach((r: any) => {
      aoa.push([
        formatBandLabel(r),
        Number(r.risk_rate_per1000_male ?? 0),
        Number(r.risk_rate_per1000_female ?? 0),
        Number(r.risk_rate_per1000 ?? 0),
        Number(r.binder_fee_per1000_male ?? 0),
        Number(r.binder_fee_per1000_female ?? 0),
        Number(r.binder_fee_per1000 ?? 0),
        Number(r.outsource_fee_per1000_male ?? 0),
        Number(r.outsource_fee_per1000_female ?? 0),
        Number(r.outsource_fee_per1000 ?? 0),
        Number(r.commission_per1000_male ?? 0),
        Number(r.commission_per1000_female ?? 0),
        Number(r.commission_per1000 ?? 0),
        Number(r.office_rate_per1000_male ?? 0),
        Number(r.office_rate_per1000_female ?? 0),
        Number(r.office_rate_per1000 ?? 0)
      ])
    })
    const ws = XLSX.utils.aoa_to_sheet(aoa)
    ws['!merges'] = [
      // Age Band cell spans both header rows.
      {
        s: { r: groupHeaderRowIdx, c: 0 },
        e: { r: groupHeaderRowIdx + 1, c: 0 }
      },
      // Each fee group label spans its three sub-columns.
      { s: { r: groupHeaderRowIdx, c: 1 }, e: { r: groupHeaderRowIdx, c: 3 } },
      { s: { r: groupHeaderRowIdx, c: 4 }, e: { r: groupHeaderRowIdx, c: 6 } },
      { s: { r: groupHeaderRowIdx, c: 7 }, e: { r: groupHeaderRowIdx, c: 9 } },
      {
        s: { r: groupHeaderRowIdx, c: 10 },
        e: { r: groupHeaderRowIdx, c: 12 }
      },
      {
        s: { r: groupHeaderRowIdx, c: 13 },
        e: { r: groupHeaderRowIdx, c: 15 }
      }
    ]
    ws['!cols'] = [
      { wch: 18 },
      { wch: 10 },
      { wch: 10 },
      { wch: 12 },
      { wch: 10 },
      { wch: 10 },
      { wch: 12 },
      { wch: 10 },
      { wch: 10 },
      { wch: 12 },
      { wch: 10 },
      { wch: 10 },
      { wch: 12 },
      { wch: 10 },
      { wch: 10 },
      { wch: 12 }
    ]
    // Excel sheet names are capped at 31 chars and cannot contain : \ / ? * [ ]
    const sheetName = String(rs.category ?? 'Category')
      .replace(/[:\\/?*[\]]/g, '_')
      .slice(0, 31)
    XLSX.utils.book_append_sheet(wb, ws, sheetName || 'Category')
  })

  const schemeName = (props.quote as any)?.scheme_name ?? 'quote'
  const fileSlug = additionalGlaCoverBenefitTitle.value
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '_')
    .replace(/^_+|_+$/g, '')
  XLSX.writeFile(
    wb,
    `${fileSlug || 'additional_gla_cover'}_summary_${schemeName}.xlsx`
  )
}

const exportExtendedFamilyToExcel = () => {
  const cats = extendedFamilyCategories.value.filter(
    (rs: any) => (rs?.extended_family_band_rates?.length ?? 0) > 0
  )
  if (cats.length === 0) return

  const wb = XLSX.utils.book_new()
  cats.forEach((rs: any) => {
    const method = rs.extended_family_pricing_method
    const aoa: any[][] = []
    aoa.push([`Extended Family Funeral — ${rs.category}`])
    aoa.push([
      `Pricing method: ${
        method === 'sum_assured' ? 'Sum Assured per Band' : 'Rate per 1,000'
      }${
        rs.extended_family_age_band_type
          ? `  ·  Age bands: ${rs.extended_family_age_band_type}`
          : ''
      }`
    ])
    aoa.push([])
    aoa.push([
      'Age Band',
      'Sum Assured',
      'Rate per 1,000',
      'Monthly Premium (per Ext. Family Member)',
      'Annual Premium (per Ext. Family Member)'
    ])
    ;(rs.extended_family_band_rates ?? []).forEach((r: any) => {
      aoa.push([
        formatBandLabel(r),
        method === 'sum_assured' ? Number(r.sum_assured ?? 0) : '',
        method === 'rate_per_1000' ? Number((r.office_rate ?? 0) * 1000) : '',
        Number(r.office_monthly_premium ?? 0),
        Number((r.office_monthly_premium ?? 0) * 12)
      ])
    })
    const ws = XLSX.utils.aoa_to_sheet(aoa)
    ws['!cols'] = [
      { wch: 18 },
      { wch: 16 },
      { wch: 16 },
      { wch: 30 },
      { wch: 30 }
    ]
    // Excel sheet names are capped at 31 chars and cannot contain : \ / ? * [ ]
    const sheetName = String(rs.category ?? 'Category')
      .replace(/[:\\/?*[\]]/g, '_')
      .slice(0, 31)
    XLSX.utils.book_append_sheet(wb, ws, sheetName || 'Category')
  })

  const schemeName = (props.quote as any)?.scheme_name ?? 'quote'
  XLSX.writeFile(wb, `extended_family_summary_${schemeName}.xlsx`)
}

const convertExcelDataToGridData = () => {
  if (!props.resultSummaries || props.resultSummaries.length === 0) return []

  const gridData: any = []

  // Sums the 10 final commission amounts that compose the Excl. Funeral total —
  // mirrors benefitAccessors() with includeInExclFun: true on the API side.
  const sumFinalCommissionExclFuneral = (src: any): number =>
    (src?.final_gla_annual_commission_amount || 0) +
    (src?.final_additional_accidental_gla_annual_commission_amount || 0) +
    (src?.final_ptd_annual_commission_amount || 0) +
    (src?.final_ci_annual_commission_amount || 0) +
    (src?.final_sgla_annual_commission_amount || 0) +
    (src?.final_tax_saver_annual_commission_amount || 0) +
    (src?.final_ttd_annual_commission_amount || 0) +
    (src?.final_phi_annual_commission_amount || 0) +
    (src?.final_gla_educator_annual_commission_amount || 0) +
    (src?.final_ptd_educator_annual_commission_amount || 0)

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
      annualPremium: computeOfficePremium(
        resultSummary.exp_total_gla_annual_risk_premium,
        resultSummary
      ),
      finalAnnualPremium: resultSummary.final_gla_annual_office_premium,
      finalAnnualCommission: resultSummary.final_gla_annual_commission_amount,
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_gla_annual_risk_premium_salary, resultSummary) * 100)}%`,
      ratePer1000SA: officeRateFromRiskRate(
        resultSummary.exp_gla_risk_rate_per_1000_sa,
        resultSummary
      )
    })

    // Tax Saver — additive top-up on top of GLA. Salary/SA cells stay 0
    // because Tax Saver rides on the GLA salary basis (showing them would
    // double-count those columns).
    if (resultSummary.tax_saver_benefit) {
      const salary = resultSummary.total_annual_salary || 0
      const taxSaverPremium =
        computeOfficePremium(
          resultSummary.exp_total_tax_saver_annual_risk_premium,
          resultSummary
        ) || 0
      gridData.push({
        category,
        benefit: `${glaBenefitTitle.value} — Tax Saver (of GLA)`,
        annualSalary: 0,
        totalSumAssured: 0,
        annualPremium: taxSaverPremium,
        finalAnnualPremium: resultSummary.final_tax_saver_annual_office_premium,
        finalAnnualCommission:
          resultSummary.final_tax_saver_annual_commission_amount,
        percentSalary: `${roundUpToTwoDecimalsAccounting(
          (salary > 0 ? taxSaverPremium / salary : 0) * 100
        )}%`,
        ratePer1000SA: ''
      })
    }

    const schemeCategory = props.quote?.scheme_categories?.find(
      (cat: any) => cat.scheme_category === category
    )

    // Additional Accidental GLA rider (per-category opt-in).
    if (schemeCategory?.additional_accidental_gla_benefit === true) {
      gridData.push({
        category,
        benefit: additionalAccidentalGlaBenefitTitle.value,
        annualSalary: resultSummary.total_annual_salary,
        totalSumAssured:
          resultSummary.total_additional_accidental_gla_capped_sum_assured,
        annualPremium: computeOfficePremium(
          resultSummary.exp_total_additional_accidental_gla_annual_risk_premium,
          resultSummary
        ),
        finalAnnualPremium:
          resultSummary.final_additional_accidental_gla_annual_office_premium,
        finalAnnualCommission:
          resultSummary.final_additional_accidental_gla_annual_commission_amount,
        percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_additional_accidental_gla_annual_risk_premium_salary, resultSummary) * 100)}%`,
        ratePer1000SA: officeRateFromRiskRate(
          resultSummary.exp_additional_accidental_gla_risk_rate_per_1000_sa,
          resultSummary
        )
      })
    }

    // GLA Educator rider (per-category opt-in).
    if (schemeCategory?.gla_educator_benefit === 'Yes') {
      gridData.push({
        category,
        benefit: glaEducatorBenefitTitle.value,
        annualSalary: resultSummary.total_annual_salary,
        totalSumAssured: resultSummary.total_educator_sum_assured,
        annualPremium: computeOfficePremium(
          resultSummary.exp_adj_total_gla_educator_risk_premium,
          resultSummary
        ),
        finalAnnualPremium:
          resultSummary.final_gla_educator_annual_office_premium,
        finalAnnualCommission:
          resultSummary.final_gla_educator_annual_commission_amount,
        percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_adj_proportion_gla_educator_risk_premium_salary, resultSummary) * 100)}%`,
        ratePer1000SA: officeRateFromRiskRate(
          resultSummary.exp_gla_educator_risk_rate_per_1000_sa,
          resultSummary
        )
      })
    }

    gridData.push({
      category,
      benefit: ptdBenefitTitle.value,
      annualSalary: isBenefitEnabled('PTD', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_ptd_capped_sum_assured,
      annualPremium: computeOfficePremium(
        resultSummary.exp_total_ptd_annual_risk_premium,
        resultSummary
      ),
      finalAnnualPremium: resultSummary.final_ptd_annual_office_premium,
      finalAnnualCommission: resultSummary.final_ptd_annual_commission_amount,
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_ptd_annual_risk_premium_salary, resultSummary) * 100)}%`,
      ratePer1000SA: officeRateFromRiskRate(
        resultSummary.exp_ptd_risk_rate_per_1000_sa,
        resultSummary
      )
    })

    // PTD Educator rider (only for categories with the benefit enabled).
    if (schemeCategory?.ptd_educator_benefit === 'Yes') {
      gridData.push({
        category,
        benefit: ptdEducatorBenefitTitle.value,
        annualSalary: resultSummary.total_annual_salary,
        totalSumAssured: resultSummary.total_educator_sum_assured,
        annualPremium: computeOfficePremium(
          resultSummary.exp_adj_total_ptd_educator_risk_premium,
          resultSummary
        ),
        finalAnnualPremium:
          resultSummary.final_ptd_educator_annual_office_premium,
        finalAnnualCommission:
          resultSummary.final_ptd_educator_annual_commission_amount,
        percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_adj_proportion_ptd_educator_risk_premium_salary, resultSummary) * 100)}%`,
        ratePer1000SA: officeRateFromRiskRate(
          resultSummary.exp_ptd_educator_risk_rate_per_1000_sa,
          resultSummary
        )
      })
    }

    gridData.push({
      category,
      benefit: ciBenefitTitle.value,
      annualSalary: isBenefitEnabled('CI', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_ci_capped_sum_assured,
      annualPremium: computeOfficePremium(
        resultSummary.exp_total_ci_annual_risk_premium,
        resultSummary
      ),
      finalAnnualPremium: resultSummary.final_ci_annual_office_premium,
      finalAnnualCommission: resultSummary.final_ci_annual_commission_amount,
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_ci_annual_risk_premium_salary, resultSummary) * 100)}%`,
      ratePer1000SA: officeRateFromRiskRate(
        resultSummary.exp_ci_risk_rate_per_1000_sa,
        resultSummary
      )
    })

    gridData.push({
      category,
      benefit: sglaBenefitTitle.value,
      annualSalary: isBenefitEnabled('SGLA', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_sgla_capped_sum_assured,
      annualPremium: computeOfficePremium(
        resultSummary.exp_total_sgla_annual_risk_premium,
        resultSummary
      ),
      finalAnnualPremium: resultSummary.final_sgla_annual_office_premium,
      finalAnnualCommission: resultSummary.final_sgla_annual_commission_amount,
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_sgla_annual_risk_premium_salary, resultSummary) * 100)}%`,
      ratePer1000SA: officeRateFromRiskRate(
        resultSummary.exp_sgla_risk_rate_per_1000_sa,
        resultSummary
      )
    })

    gridData.push({
      category,
      benefit: phiBenefitTitle.value,
      annualSalary: isBenefitEnabled('PHI', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_phi_capped_income,
      annualPremium: computeOfficePremium(
        resultSummary.exp_total_phi_annual_risk_premium,
        resultSummary
      ),
      finalAnnualPremium: resultSummary.final_phi_annual_office_premium,
      finalAnnualCommission: resultSummary.final_phi_annual_commission_amount,
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_phi_annual_risk_premium_salary, resultSummary) * 100)}%`,
      ratePer1000SA: officeRateFromRiskRate(
        resultSummary.exp_phi_risk_rate_per_1000_sa,
        resultSummary
      )
    })

    gridData.push({
      category,
      benefit: ttdBenefitTitle.value,
      annualSalary: isBenefitEnabled('TTD', category)
        ? resultSummary.total_annual_salary
        : 0,
      totalSumAssured: resultSummary.total_ttd_capped_income,
      annualPremium: computeOfficePremium(
        resultSummary.exp_total_ttd_annual_risk_premium,
        resultSummary
      ),
      finalAnnualPremium: resultSummary.final_ttd_annual_office_premium,
      finalAnnualCommission: resultSummary.final_ttd_annual_commission_amount,
      percentSalary: `${roundUpToTwoDecimalsAccounting(officeProportionFromRiskProportion(resultSummary.exp_proportion_ttd_annual_risk_premium_salary, resultSummary) * 100)}%`,
      ratePer1000SA: officeRateFromRiskRate(
        resultSummary.exp_ttd_risk_rate_per_1000_sa,
        resultSummary
      )
    })

    // Add subtotal row. The backend's exp_total_annual_premium_excl_funeral
    // already includes the GLA TaxSaver rider plus GLA Educator and PTD Educator
    // on top of the six core benefits.
    const anyBenefitEnabled = ['GLA', 'PTD', 'CI', 'SGLA', 'PHI', 'TTD'].some(
      (benefitCode) => isBenefitEnabled(benefitCode, category)
    )
    gridData.push({
      category,
      benefit: 'Sub Total (Excl. Funeral)',
      annualSalary: anyBenefitEnabled ? resultSummary.total_annual_salary : 0,
      totalSumAssured: resultSummary.total_gla_capped_sum_assured,
      annualPremium: resultSummary.exp_total_annual_premium_excl_funeral,
      finalAnnualPremium: resultSummary.final_total_annual_premium_excl_funeral,
      finalAnnualCommission: sumFinalCommissionExclFuneral(resultSummary),
      percentSalary: `${roundUpToTwoDecimalsAccounting(
        resultSummary.total_annual_salary > 0
          ? (resultSummary.exp_total_annual_premium_excl_funeral /
              resultSummary.total_annual_salary) *
              100
          : 0
      )}%`,
      ratePer1000SA: '',
      isSubtotal: true
    })

    gridData.push({
      category,
      benefit: 'Group Funeral Annual Premium per Member',
      annualSalary: '',
      totalSumAssured: '',
      annualPremium: resultSummary.exp_total_fun_annual_premium_per_member,
      finalAnnualPremium: '',
      finalAnnualCommission: '',
      percentSalary: '',
      ratePer1000SA: ''
    })

    gridData.push({
      category,
      benefit: 'Group Funeral Annual Premium',
      annualSalary: '',
      totalSumAssured: '',
      annualPremium: computeOfficePremium(
        resultSummary.exp_total_fun_annual_risk_premium,
        resultSummary
      ),
      finalAnnualPremium: resultSummary.final_fun_annual_office_premium,
      finalAnnualCommission: resultSummary.final_fun_annual_commission_amount,
      percentSalary: '',
      ratePer1000SA: ''
    })
  })

  // Calculate totals across all categories.
  // Office premium / proportion / rate fields are accumulated by computing
  // each category's value from its own scheme loading first and then summing
  // — this preserves correctness across schemes that have different
  // expense / commission / profit loading mixes.
  if (props.resultSummaries.length > 0) {
    const totals: any = props.resultSummaries.reduce(
      (acc: any, resultSummary: any) => {
        return {
          total_gla_capped_sum_assured:
            (acc.total_gla_capped_sum_assured || 0) +
            (resultSummary.total_gla_capped_sum_assured || 0),
          exp_total_gla_annual_office_premium:
            (acc.exp_total_gla_annual_office_premium || 0) +
            (computeOfficePremium(
              resultSummary.exp_total_gla_annual_risk_premium,
              resultSummary
            ) || 0),
          exp_proportion_gla_office_premium_salary:
            (acc.exp_proportion_gla_office_premium_salary || 0) +
            (officeProportionFromRiskProportion(
              resultSummary.exp_proportion_gla_annual_risk_premium_salary,
              resultSummary
            ) || 0),

          total_ptd_capped_sum_assured:
            (acc.total_ptd_capped_sum_assured || 0) +
            (resultSummary.total_ptd_capped_sum_assured || 0),
          exp_total_ptd_annual_office_premium:
            (acc.exp_total_ptd_annual_office_premium || 0) +
            (computeOfficePremium(
              resultSummary.exp_total_ptd_annual_risk_premium,
              resultSummary
            ) || 0),
          exp_proportion_ptd_office_premium_salary:
            (acc.exp_proportion_ptd_office_premium_salary || 0) +
            (officeProportionFromRiskProportion(
              resultSummary.exp_proportion_ptd_annual_risk_premium_salary,
              resultSummary
            ) || 0),

          total_ci_capped_sum_assured:
            (acc.total_ci_capped_sum_assured || 0) +
            (resultSummary.total_ci_capped_sum_assured || 0),
          exp_total_ci_annual_office_premium:
            (acc.exp_total_ci_annual_office_premium || 0) +
            (computeOfficePremium(
              resultSummary.exp_total_ci_annual_risk_premium,
              resultSummary
            ) || 0),
          exp_proportion_ci_office_premium_salary:
            (acc.exp_proportion_ci_office_premium_salary || 0) +
            (officeProportionFromRiskProportion(
              resultSummary.exp_proportion_ci_annual_risk_premium_salary,
              resultSummary
            ) || 0),

          total_sgla_capped_sum_assured:
            (acc.total_sgla_capped_sum_assured || 0) +
            (resultSummary.total_sgla_capped_sum_assured || 0),
          exp_total_sgla_annual_office_premium:
            (acc.exp_total_sgla_annual_office_premium || 0) +
            (computeOfficePremium(
              resultSummary.exp_total_sgla_annual_risk_premium,
              resultSummary
            ) || 0),
          exp_proportion_sgla_office_premium_salary:
            (acc.exp_proportion_sgla_office_premium_salary || 0) +
            (officeProportionFromRiskProportion(
              resultSummary.exp_proportion_sgla_annual_risk_premium_salary,
              resultSummary
            ) || 0),

          total_phi_capped_income:
            (acc.total_phi_capped_income || 0) +
            (resultSummary.total_phi_capped_income || 0),
          exp_total_phi_annual_office_premium:
            (acc.exp_total_phi_annual_office_premium || 0) +
            (computeOfficePremium(
              resultSummary.exp_total_phi_annual_risk_premium,
              resultSummary
            ) || 0),
          exp_proportion_phi_office_premium_salary:
            (acc.exp_proportion_phi_office_premium_salary || 0) +
            (officeProportionFromRiskProportion(
              resultSummary.exp_proportion_phi_annual_risk_premium_salary,
              resultSummary
            ) || 0),

          total_ttd_capped_income:
            (acc.total_ttd_capped_income || 0) +
            (resultSummary.total_ttd_capped_income || 0),
          exp_total_ttd_annual_office_premium:
            (acc.exp_total_ttd_annual_office_premium || 0) +
            (computeOfficePremium(
              resultSummary.exp_total_ttd_annual_risk_premium,
              resultSummary
            ) || 0),
          exp_proportion_ttd_office_premium_salary:
            (acc.exp_proportion_ttd_office_premium_salary || 0) +
            (officeProportionFromRiskProportion(
              resultSummary.exp_proportion_ttd_annual_risk_premium_salary,
              resultSummary
            ) || 0),

          exp_total_annual_premium_excl_funeral:
            (acc.exp_total_annual_premium_excl_funeral || 0) +
            (resultSummary.exp_total_annual_premium_excl_funeral || 0),

          exp_total_fun_monthly_premium_per_member:
            (acc.exp_total_fun_monthly_premium_per_member || 0) +
            (resultSummary.exp_total_fun_monthly_premium_per_member || 0),
          exp_total_fun_annual_premium_per_member:
            (acc.exp_total_fun_annual_premium_per_member || 0) +
            (resultSummary.exp_total_fun_annual_premium_per_member || 0),
          exp_total_fun_annual_office_premium:
            (acc.exp_total_fun_annual_office_premium || 0) +
            (computeOfficePremium(
              resultSummary.exp_total_fun_annual_risk_premium,
              resultSummary
            ) || 0),

          total_educator_sum_assured:
            (acc.total_educator_sum_assured || 0) +
            (resultSummary.total_educator_sum_assured || 0),
          exp_adj_total_gla_educator_office_premium:
            (acc.exp_adj_total_gla_educator_office_premium || 0) +
            (computeOfficePremium(
              resultSummary.exp_adj_total_gla_educator_risk_premium,
              resultSummary
            ) || 0),
          exp_adj_total_ptd_educator_office_premium:
            (acc.exp_adj_total_ptd_educator_office_premium || 0) +
            (computeOfficePremium(
              resultSummary.exp_adj_total_ptd_educator_risk_premium,
              resultSummary
            ) || 0),

          total_annual_salary:
            (acc.total_annual_salary || 0) +
            (resultSummary.total_annual_salary || 0),

          // Final*OfficePremium and Final*CommissionAmount are persisted on
          // each summary by the backend (recomputeFinalPremiumsAndCommission)
          // so the totals row simply sums them — no client-side derivation.
          final_gla_annual_office_premium:
            (acc.final_gla_annual_office_premium || 0) +
            (resultSummary.final_gla_annual_office_premium || 0),
          final_gla_annual_commission_amount:
            (acc.final_gla_annual_commission_amount || 0) +
            (resultSummary.final_gla_annual_commission_amount || 0),
          final_tax_saver_annual_commission_amount:
            (acc.final_tax_saver_annual_commission_amount || 0) +
            (resultSummary.final_tax_saver_annual_commission_amount || 0),
          final_additional_accidental_gla_annual_commission_amount:
            (acc.final_additional_accidental_gla_annual_commission_amount ||
              0) +
            (resultSummary.final_additional_accidental_gla_annual_commission_amount ||
              0),
          final_ptd_annual_office_premium:
            (acc.final_ptd_annual_office_premium || 0) +
            (resultSummary.final_ptd_annual_office_premium || 0),
          final_ptd_annual_commission_amount:
            (acc.final_ptd_annual_commission_amount || 0) +
            (resultSummary.final_ptd_annual_commission_amount || 0),
          final_ci_annual_office_premium:
            (acc.final_ci_annual_office_premium || 0) +
            (resultSummary.final_ci_annual_office_premium || 0),
          final_ci_annual_commission_amount:
            (acc.final_ci_annual_commission_amount || 0) +
            (resultSummary.final_ci_annual_commission_amount || 0),
          final_sgla_annual_office_premium:
            (acc.final_sgla_annual_office_premium || 0) +
            (resultSummary.final_sgla_annual_office_premium || 0),
          final_sgla_annual_commission_amount:
            (acc.final_sgla_annual_commission_amount || 0) +
            (resultSummary.final_sgla_annual_commission_amount || 0),
          final_phi_annual_office_premium:
            (acc.final_phi_annual_office_premium || 0) +
            (resultSummary.final_phi_annual_office_premium || 0),
          final_phi_annual_commission_amount:
            (acc.final_phi_annual_commission_amount || 0) +
            (resultSummary.final_phi_annual_commission_amount || 0),
          final_ttd_annual_office_premium:
            (acc.final_ttd_annual_office_premium || 0) +
            (resultSummary.final_ttd_annual_office_premium || 0),
          final_ttd_annual_commission_amount:
            (acc.final_ttd_annual_commission_amount || 0) +
            (resultSummary.final_ttd_annual_commission_amount || 0),
          final_fun_annual_office_premium:
            (acc.final_fun_annual_office_premium || 0) +
            (resultSummary.final_fun_annual_office_premium || 0),
          final_fun_annual_commission_amount:
            (acc.final_fun_annual_commission_amount || 0) +
            (resultSummary.final_fun_annual_commission_amount || 0),
          final_gla_educator_annual_office_premium:
            (acc.final_gla_educator_annual_office_premium || 0) +
            (resultSummary.final_gla_educator_annual_office_premium || 0),
          final_gla_educator_annual_commission_amount:
            (acc.final_gla_educator_annual_commission_amount || 0) +
            (resultSummary.final_gla_educator_annual_commission_amount || 0),
          final_ptd_educator_annual_office_premium:
            (acc.final_ptd_educator_annual_office_premium || 0) +
            (resultSummary.final_ptd_educator_annual_office_premium || 0),
          final_ptd_educator_annual_commission_amount:
            (acc.final_ptd_educator_annual_commission_amount || 0) +
            (resultSummary.final_ptd_educator_annual_commission_amount || 0),
          final_total_annual_premium_excl_funeral:
            (acc.final_total_annual_premium_excl_funeral || 0) +
            (resultSummary.final_total_annual_premium_excl_funeral || 0)
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
      finalAnnualPremium: totals.final_gla_annual_office_premium,
      finalAnnualCommission: totals.final_gla_annual_commission_amount,
      percentSalary: `${roundUpToTwoDecimalsAccounting(
        (totals.exp_total_gla_annual_office_premium /
          totals.total_annual_salary || 0) * 100
      )}%`,
      ratePer1000SA: totals.total_gla_capped_sum_assured
        ? (totals.exp_total_gla_annual_office_premium * 1000) /
          totals.total_gla_capped_sum_assured
        : ''
    })

    const anyCategoryHasGlaEducator = props.quote?.scheme_categories?.some(
      (cat: any) => cat.gla_educator_benefit === 'Yes'
    )
    if (anyCategoryHasGlaEducator) {
      gridData.push({
        category: totalsCategory,
        benefit: glaEducatorBenefitTitle.value,
        annualSalary: totals.total_annual_salary,
        totalSumAssured: totals.total_educator_sum_assured,
        annualPremium: totals.exp_adj_total_gla_educator_office_premium,
        finalAnnualPremium: totals.final_gla_educator_annual_office_premium,
        finalAnnualCommission:
          totals.final_gla_educator_annual_commission_amount,
        percentSalary: `${roundUpToTwoDecimalsAccounting(
          (totals.exp_adj_total_gla_educator_office_premium /
            totals.total_annual_salary || 0) * 100
        )}%`,
        ratePer1000SA: totals.total_educator_sum_assured
          ? (totals.exp_adj_total_gla_educator_office_premium * 1000) /
            totals.total_educator_sum_assured
          : ''
      })
    }

    gridData.push({
      category: totalsCategory,
      benefit: ptdBenefitTitle.value,
      annualSalary: totals.total_annual_salary,
      totalSumAssured: totals.total_ptd_capped_sum_assured,
      annualPremium: totals.exp_total_ptd_annual_office_premium,
      finalAnnualPremium: totals.final_ptd_annual_office_premium,
      finalAnnualCommission: totals.final_ptd_annual_commission_amount,
      percentSalary: `${roundUpToTwoDecimalsAccounting((totals.exp_total_ptd_annual_office_premium / totals.total_annual_salary || 0) * 100)}%`,
      ratePer1000SA: totals.total_ptd_capped_sum_assured
        ? (totals.exp_total_ptd_annual_office_premium * 1000) /
          totals.total_ptd_capped_sum_assured
        : ''
    })

    const anyCategoryHasPtdEducator = props.quote?.scheme_categories?.some(
      (cat: any) => cat.ptd_educator_benefit === 'Yes'
    )
    if (anyCategoryHasPtdEducator) {
      gridData.push({
        category: totalsCategory,
        benefit: ptdEducatorBenefitTitle.value,
        annualSalary: totals.total_annual_salary,
        totalSumAssured: totals.total_educator_sum_assured,
        annualPremium: totals.exp_adj_total_ptd_educator_office_premium,
        finalAnnualPremium: totals.final_ptd_educator_annual_office_premium,
        finalAnnualCommission:
          totals.final_ptd_educator_annual_commission_amount,
        percentSalary: `${roundUpToTwoDecimalsAccounting(
          (totals.exp_adj_total_ptd_educator_office_premium /
            totals.total_annual_salary || 0) * 100
        )}%`,
        ratePer1000SA: totals.total_educator_sum_assured
          ? (totals.exp_adj_total_ptd_educator_office_premium * 1000) /
            totals.total_educator_sum_assured
          : ''
      })
    }

    gridData.push({
      category: totalsCategory,
      benefit: ciBenefitTitle.value,
      annualSalary: totals.total_annual_salary,
      totalSumAssured: totals.total_ci_capped_sum_assured,
      annualPremium: totals.exp_total_ci_annual_office_premium,
      finalAnnualPremium: totals.final_ci_annual_office_premium,
      finalAnnualCommission: totals.final_ci_annual_commission_amount,
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
      finalAnnualPremium: totals.final_sgla_annual_office_premium,
      finalAnnualCommission: totals.final_sgla_annual_commission_amount,
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
      finalAnnualPremium: totals.final_phi_annual_office_premium,
      finalAnnualCommission: totals.final_phi_annual_commission_amount,
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
      finalAnnualPremium: totals.final_ttd_annual_office_premium,
      finalAnnualCommission: totals.final_ttd_annual_commission_amount,
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
      finalAnnualPremium: totals.final_total_annual_premium_excl_funeral,
      finalAnnualCommission: sumFinalCommissionExclFuneral(totals),
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
      finalAnnualPremium: '',
      finalAnnualCommission: '',
      percentSalary: '',
      ratePer1000SA: ''
    })

    gridData.push({
      category: totalsCategory,
      benefit: 'Group Funeral Annual Premium',
      annualSalary: '',
      totalSumAssured: '',
      annualPremium: totals.exp_total_fun_annual_office_premium,
      finalAnnualPremium: totals.final_fun_annual_office_premium,
      finalAnnualCommission: totals.final_fun_annual_commission_amount,
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
    const additionalAccidentalGlaBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'AAGLA'
    )
    if (additionalAccidentalGlaBenefit) {
      additionalAccidentalGlaBenefitTitle.value =
        additionalAccidentalGlaBenefit.benefit_alias ||
        additionalAccidentalGlaBenefit.benefit_name
    }
    const additionalGlaCoverBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'AGLA'
    )
    if (additionalGlaCoverBenefit) {
      additionalGlaCoverBenefitTitle.value =
        additionalGlaCoverBenefit.benefit_alias ||
        additionalGlaCoverBenefit.benefit_name
    }
    const glaEducatorBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'GLA_EDU'
    )
    if (glaEducatorBenefit) {
      glaEducatorBenefitTitle.value =
        glaEducatorBenefit.benefit_alias?.trim() ||
        glaEducatorBenefit.benefit_name ||
        'GLA Educator'
    }
    const ptdEducatorBenefit: any = benefitMaps.value.find(
      (item: any) => item.benefit_code === 'PTD_EDU'
    )
    if (ptdEducatorBenefit) {
      ptdEducatorBenefitTitle.value =
        ptdEducatorBenefit.benefit_alias?.trim() ||
        ptdEducatorBenefit.benefit_name ||
        'PTD Educator'
    }
  })
})

// Build rowData only once the async benefit titles (loaded in onMounted via
// getBenefitMaps) are ready. AG Grid is configured without a getRowId, so if
// rowData is swapped out mid-lifecycle — first with empty benefit labels,
// then with real labels — it mis-maps cells across the two snapshots and
// loses values on some columns (observed: % of Salary and Rate per 1000 SA
// render as blank/"-" even though rowData has the correct strings/numbers).
// Gating on titles collapses the two rebuilds into a single setRowData so
// the grid sees one clean dataset.
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
    const titlesReady =
      glaBenefitTitle.value !== '' &&
      ptdBenefitTitle.value !== '' &&
      ciBenefitTitle.value !== '' &&
      sglaBenefitTitle.value !== '' &&
      phiBenefitTitle.value !== '' &&
      ttdBenefitTitle.value !== ''
    if (
      titlesReady &&
      props.resultSummaries &&
      props.resultSummaries.length > 0
    ) {
      rowData.value = convertExcelDataToGridData()
      rowDataKey.value++
      breakdownRowData.value = convertResultSummariesToBreakdown()
      breakdownRowDataKey.value++
    } else {
      rowData.value = []
      breakdownRowData.value = []
    }
  },
  { deep: true, immediate: true }
)

// Build per-category, per-benefit Premium Breakdown rows. Decomposes each
// benefit's office premium into base premium + binder fee + outsourcing fee
// + commission, both pre-discount (Exp*) and post-discount (Final*).
//
// expBase = ExpOfficePremium - expBinder - expOutsource - expCommission
// finalBase = FinalOfficePremium - finalBinder - finalOutsource - finalCommission
//
// Pre-discount values come straight from persisted Exp* fields. Final binder
// and outsource amounts are persisted by recomputeFinalPremiumsAndCommission
// (added 2026-04-28) — finalBinder = FinalOfficePremium * BinderFeeRate,
// 0 on non-binder distribution channels.
function convertResultSummariesToBreakdown(): any[] {
  if (!props.resultSummaries || props.resultSummaries.length === 0) return []

  const isBenefitEnabled = (benefitCode: string, categoryName: string) => {
    if (!props.quote?.scheme_categories) return true
    const schemeCategory = props.quote.scheme_categories.find(
      (cat: any) => cat.scheme_category === categoryName
    )
    if (!schemeCategory) return true
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

  const num = (v: any) => (typeof v === 'number' && Number.isFinite(v) ? v : 0)

  // Build a row from the (binder, outsource, commission) triple plus the
  // pre/post office premium.
  //
  // Pre side: expOffice is the *pre-commission* office premium grossed up by
  // the full scheme loading:
  //   expOffice = risk / (1 - (expense + profit + admin + other + binder + outsourcing))
  // Binder + outsource rates sit *inside* the denominator, so binder and
  // outsource amounts (= expOffice * rate) are already embedded in expOffice.
  // Commission is excluded from the gross-up and added on top:
  //   expBase  = expOffice - expBinder - expOutsource          (no commission subtraction)
  //   expGross = expOffice + expCommission                     = base + binder + outsource + commission
  // Invariants the user expects:
  //   gross - commission == ExpRiskPremium / (1 - SchemeTotal)
  //   base              == ExpRiskPremium / (1 - SchemeTotal) - binder - outsource
  //
  // Final side: same denominator structure as Pre side, plus discount:
  //   schemeLoadingAfterDiscount = expense + profit + admin + other + binder + outsourcing + discount
  //   finalOfficePreComm         = expRiskPremium / (1 - schemeLoadingAfterDiscount)
  // Discount is stored as a negative fraction so adding it shrinks the
  // denominator (and thus shrinks the office premium once a discount applies).
  // Backend pass 3 of recomputeFinalPremiumsAndCommission re-derives
  // finalCommission progressively against the new gross-up and persists
  //   finalOffice = finalOfficePreComm + finalCommission
  // i.e. the persisted finalOffice column already includes its commission
  // slice. So commission *is* subtracted from finalOffice when carving out
  // finalBase:
  //   finalBase    = finalOffice - finalBinder - finalOutsource - finalCommission
  //   finalPremium = finalOffice                                = finalBase + finalBinder + finalOutsource + finalCommission
  //                = expRiskPremium / (1 - schemeLoadingAfterDiscount) + finalCommission
  const buildRow = (
    category: string,
    benefit: string,
    expRiskPremium: number,
    expOffice: number,
    expBinder: number,
    expOutsource: number,
    expCommission: number,
    finalOffice: number,
    finalBinder: number,
    finalOutsource: number,
    finalCommission: number,
    extra: Record<string, any> = {}
  ) => {
    const expBase = expOffice - expBinder - expOutsource
    const finalBase =
      finalOffice - finalBinder - finalOutsource - finalCommission
    return {
      category,
      benefit,
      expRiskPremium,
      expBasePremium: expBase,
      expBinderFee: expBinder,
      expOutsourcingFee: expOutsource,
      expCommission,
      // Gross = ExpRiskPremium / (1 - schemeTotalLoading) + commission
      // (binder + outsource are already inside expOffice and cancel out).
      expGrossPremium: expOffice + expCommission,
      finalBasePremium: finalBase,
      finalBinderFee: finalBinder,
      finalOutsourcingFee: finalOutsource,
      finalCommission,
      finalPremium: finalBase + finalBinder + finalOutsource + finalCommission,
      ...extra
    }
  }

  const gridData: any[] = []
  const totals = {
    expRisk: 0,
    expBase: 0,
    expBinder: 0,
    expOutsource: 0,
    expCommission: 0,
    expGross: 0,
    finalBase: 0,
    finalBinder: 0,
    finalOutsource: 0,
    finalCommission: 0,
    finalPremium: 0
  }

  props.resultSummaries.forEach((rs: any) => {
    const category = rs.category
    const schemeCategory = props.quote?.scheme_categories?.find(
      (cat: any) => cat.scheme_category === category
    )

    let catExpRisk = 0
    let catExpBase = 0
    let catExpBinder = 0
    let catExpOutsource = 0
    let catExpCommission = 0
    let catExpGross = 0
    let catFinalBase = 0
    let catFinalBinder = 0
    let catFinalOutsource = 0
    let catFinalCommission = 0
    let catFinalPremium = 0

    const pushBenefit = (
      label: string,
      expRiskPremium: number,
      expOffice: number,
      expBinder: number,
      expOutsource: number,
      expCommission: number,
      finalOffice: number,
      finalBinder: number,
      finalOutsource: number,
      finalCommission: number,
      includeInTotals = true
    ) => {
      const row = buildRow(
        category,
        label,
        expRiskPremium,
        expOffice,
        expBinder,
        expOutsource,
        expCommission,
        finalOffice,
        finalBinder,
        finalOutsource,
        finalCommission
      )
      gridData.push(row)
      if (includeInTotals) {
        catExpRisk += expRiskPremium
        catExpBase += row.expBasePremium
        catExpBinder += expBinder
        catExpOutsource += expOutsource
        catExpCommission += expCommission
        catExpGross += row.expGrossPremium
        catFinalBase += row.finalBasePremium
        catFinalBinder += finalBinder
        catFinalOutsource += finalOutsource
        catFinalCommission += finalCommission
        catFinalPremium += row.finalPremium
      }
    }

    // GLA
    if (isBenefitEnabled('GLA', category)) {
      pushBenefit(
        glaBenefitTitle.value,
        num(rs.exp_total_gla_annual_risk_premium),
        computeOfficePremium(num(rs.exp_total_gla_annual_risk_premium), rs),
        num(rs.exp_total_gla_annual_binder_amount),
        num(rs.exp_total_gla_annual_outsourced_amount),
        num(rs.exp_total_gla_annual_commission_amount),
        num(rs.final_gla_annual_office_premium),
        num(rs.final_gla_annual_binder_amount),
        num(rs.final_gla_annual_outsourced_amount),
        num(rs.final_gla_annual_commission_amount)
      )
    }

    // Tax Saver — additive top-up on top of GLA. No Exp* binder/outsource
    // is persisted for tax saver, so those columns are 0 even on binder schemes.
    if (rs.tax_saver_benefit) {
      pushBenefit(
        `${glaBenefitTitle.value} — Tax Saver (of GLA)`,
        num(rs.exp_total_tax_saver_annual_risk_premium),
        computeOfficePremium(
          num(rs.exp_total_tax_saver_annual_risk_premium),
          rs
        ),
        0,
        0,
        num(rs.exp_total_tax_saver_annual_commission_amount),
        num(rs.final_tax_saver_annual_office_premium),
        num(rs.final_tax_saver_annual_binder_amount),
        num(rs.final_tax_saver_annual_outsourced_amount),
        num(rs.final_tax_saver_annual_commission_amount)
      )
    }

    // Additional Accidental GLA rider (per-category opt-in).
    if (schemeCategory?.additional_accidental_gla_benefit === true) {
      pushBenefit(
        additionalAccidentalGlaBenefitTitle.value,
        num(rs.exp_total_additional_accidental_gla_annual_risk_premium),
        computeOfficePremium(
          num(rs.exp_total_additional_accidental_gla_annual_risk_premium),
          rs
        ),
        num(rs.exp_total_additional_accidental_gla_annual_binder_amount),
        num(rs.exp_total_additional_accidental_gla_annual_outsourced_amount),
        num(rs.exp_total_additional_accidental_gla_annual_commission_amount),
        num(rs.final_additional_accidental_gla_annual_office_premium),
        num(rs.final_additional_accidental_gla_annual_binder_amount),
        num(rs.final_additional_accidental_gla_annual_outsourced_amount),
        num(rs.final_additional_accidental_gla_annual_commission_amount)
      )
    }

    // GLA Educator (per-category opt-in).
    if (schemeCategory?.gla_educator_benefit === 'Yes') {
      pushBenefit(
        glaEducatorBenefitTitle.value,
        num(rs.exp_adj_total_gla_educator_risk_premium),
        computeOfficePremium(
          num(rs.exp_adj_total_gla_educator_risk_premium),
          rs
        ),
        num(rs.exp_adj_total_gla_educator_binder_amount),
        num(rs.exp_adj_total_gla_educator_outsourced_amount),
        num(rs.exp_adj_total_gla_educator_commission_amount),
        num(rs.final_gla_educator_annual_office_premium),
        num(rs.final_gla_educator_annual_binder_amount),
        num(rs.final_gla_educator_annual_outsourced_amount),
        num(rs.final_gla_educator_annual_commission_amount)
      )
    }

    // PTD
    if (isBenefitEnabled('PTD', category)) {
      pushBenefit(
        ptdBenefitTitle.value,
        num(rs.exp_total_ptd_annual_risk_premium),
        computeOfficePremium(num(rs.exp_total_ptd_annual_risk_premium), rs),
        num(rs.exp_total_ptd_annual_binder_amount),
        num(rs.exp_total_ptd_annual_outsourced_amount),
        num(rs.exp_total_ptd_annual_commission_amount),
        num(rs.final_ptd_annual_office_premium),
        num(rs.final_ptd_annual_binder_amount),
        num(rs.final_ptd_annual_outsourced_amount),
        num(rs.final_ptd_annual_commission_amount)
      )
    }

    if (schemeCategory?.ptd_educator_benefit === 'Yes') {
      pushBenefit(
        ptdEducatorBenefitTitle.value,
        num(rs.exp_adj_total_ptd_educator_risk_premium),
        computeOfficePremium(
          num(rs.exp_adj_total_ptd_educator_risk_premium),
          rs
        ),
        num(rs.exp_adj_total_ptd_educator_binder_amount),
        num(rs.exp_adj_total_ptd_educator_outsourced_amount),
        num(rs.exp_adj_total_ptd_educator_commission_amount),
        num(rs.final_ptd_educator_annual_office_premium),
        num(rs.final_ptd_educator_annual_binder_amount),
        num(rs.final_ptd_educator_annual_outsourced_amount),
        num(rs.final_ptd_educator_annual_commission_amount)
      )
    }

    // CI
    if (isBenefitEnabled('CI', category)) {
      pushBenefit(
        ciBenefitTitle.value,
        num(rs.exp_total_ci_annual_risk_premium),
        computeOfficePremium(num(rs.exp_total_ci_annual_risk_premium), rs),
        num(rs.exp_total_ci_annual_binder_amount),
        num(rs.exp_total_ci_annual_outsourced_amount),
        num(rs.exp_total_ci_annual_commission_amount),
        num(rs.final_ci_annual_office_premium),
        num(rs.final_ci_annual_binder_amount),
        num(rs.final_ci_annual_outsourced_amount),
        num(rs.final_ci_annual_commission_amount)
      )
    }

    // SGLA
    if (isBenefitEnabled('SGLA', category)) {
      pushBenefit(
        sglaBenefitTitle.value,
        num(rs.exp_total_sgla_annual_risk_premium),
        computeOfficePremium(num(rs.exp_total_sgla_annual_risk_premium), rs),
        num(rs.exp_total_sgla_annual_binder_amount),
        num(rs.exp_total_sgla_annual_outsourced_amount),
        num(rs.exp_total_sgla_annual_commission_amount),
        num(rs.final_sgla_annual_office_premium),
        num(rs.final_sgla_annual_binder_amount),
        num(rs.final_sgla_annual_outsourced_amount),
        num(rs.final_sgla_annual_commission_amount)
      )
    }

    // PHI
    if (isBenefitEnabled('PHI', category)) {
      pushBenefit(
        phiBenefitTitle.value,
        num(rs.exp_total_phi_annual_risk_premium),
        computeOfficePremium(num(rs.exp_total_phi_annual_risk_premium), rs),
        num(rs.exp_total_phi_annual_binder_amount),
        num(rs.exp_total_phi_annual_outsourced_amount),
        num(rs.exp_total_phi_annual_commission_amount),
        num(rs.final_phi_annual_office_premium),
        num(rs.final_phi_annual_binder_amount),
        num(rs.final_phi_annual_outsourced_amount),
        num(rs.final_phi_annual_commission_amount)
      )
    }

    // TTD
    if (isBenefitEnabled('TTD', category)) {
      pushBenefit(
        ttdBenefitTitle.value,
        num(rs.exp_total_ttd_annual_risk_premium),
        computeOfficePremium(num(rs.exp_total_ttd_annual_risk_premium), rs),
        num(rs.exp_total_ttd_annual_binder_amount),
        num(rs.exp_total_ttd_annual_outsourced_amount),
        num(rs.exp_total_ttd_annual_commission_amount),
        num(rs.final_ttd_annual_office_premium),
        num(rs.final_ttd_annual_binder_amount),
        num(rs.final_ttd_annual_outsourced_amount),
        num(rs.final_ttd_annual_commission_amount)
      )
    }

    // Per-category sub-total (excl. funeral)
    gridData.push({
      category,
      benefit: 'Sub Total (Excl. Funeral)',
      expRiskPremium: catExpRisk,
      expBasePremium: catExpBase,
      expBinderFee: catExpBinder,
      expOutsourcingFee: catExpOutsource,
      expCommission: catExpCommission,
      expGrossPremium: catExpGross,
      finalBasePremium: catFinalBase,
      finalBinderFee: catFinalBinder,
      finalOutsourcingFee: catFinalOutsource,
      finalCommission: catFinalCommission,
      finalPremium: catFinalPremium,
      isSubtotal: true
    })

    // Group Funeral row (always shown if there's any funeral premium)
    pushBenefit(
      'Group Funeral',
      num(rs.exp_total_fun_annual_risk_premium),
      computeOfficePremium(num(rs.exp_total_fun_annual_risk_premium), rs),
      num(rs.exp_total_fun_annual_binder_amount),
      num(rs.exp_total_fun_annual_outsourced_amount),
      num(rs.exp_total_fun_annual_commission_amount),
      num(rs.final_fun_annual_office_premium),
      num(rs.final_fun_annual_binder_amount),
      num(rs.final_fun_annual_outsourced_amount),
      num(rs.final_fun_annual_commission_amount)
    )

    // Roll into scheme-wide totals (using post-pushBenefit accumulators that
    // already include Funeral via includeInTotals=true on the funeral push).
    totals.expRisk += catExpRisk
    totals.expBase += catExpBase
    totals.expBinder += catExpBinder
    totals.expOutsource += catExpOutsource
    totals.expCommission += catExpCommission
    totals.expGross += catExpGross
    totals.finalBase += catFinalBase
    totals.finalBinder += catFinalBinder
    totals.finalOutsource += catFinalOutsource
    totals.finalCommission += catFinalCommission
    totals.finalPremium += catFinalPremium
  })

  // Scheme-wide totals row (under a synthetic "Totals" group).
  gridData.push({
    category: 'Totals',
    benefit: 'Total',
    expRiskPremium: totals.expRisk,
    expBasePremium: totals.expBase,
    expBinderFee: totals.expBinder,
    expOutsourcingFee: totals.expOutsource,
    expCommission: totals.expCommission,
    expGrossPremium: totals.expGross,
    finalBasePremium: totals.finalBase,
    finalBinderFee: totals.finalBinder,
    finalOutsourcingFee: totals.finalOutsource,
    finalCommission: totals.finalCommission,
    finalPremium: totals.finalPremium,
    isSubtotal: true
  })

  return gridData
}
</script>
<style scoped>
/*
 * Style the Premium Breakdown grouped header bands ("Before Discount" /
 * "After Discount") on-brand using the app primary (#003F58). AG Grid
 * renders headers outside the scoped boundary, so :deep() is required to
 * reach into the grid DOM.
 */
:deep(.ag-theme-balham .ag-header-group-cell.breakdown-group-header) {
  background-color: #003f58;
  color: #fff;
  font-weight: 600;
  font-size: 0.8rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  border-right: 1px solid rgba(255, 255, 255, 0.25);
}
:deep(
  .ag-theme-balham
    .ag-header-group-cell.breakdown-group-header
    .ag-header-group-cell-label
) {
  justify-content: center;
}
/* Subtle divider between the two bands so the boundary reads clearly
 * without colour-coding each side. */
:deep(.ag-theme-balham .ag-header-group-cell.breakdown-group-pre) {
  border-right: 2px solid #fff;
}

/*
 * Additional GLA Cover Summary table — visual demarcation across the
 * five fee groups (Risk / Binder / Outsource / Commission / Office). Each
 * group spans M, F, Combined sub-columns, so the eye needs help finding
 * group boundaries and the headline (Combined) value within each block.
 *
 * v-table renders headers/cells outside the scoped boundary on Vuetify 3,
 * so :deep() is needed to pierce into the rendered <th>/<td>.
 */

/* Sticky Age Band column so horizontal scroll keeps the row label in view. */
:deep(.agla-summary-table th.agla-age-band),
:deep(.agla-summary-table td.agla-age-band) {
  position: sticky;
  left: 0;
  background-color: rgb(var(--v-theme-surface));
  border-right: 2px solid rgba(var(--v-border-color), var(--v-border-opacity));
  z-index: 1;
}
:deep(.agla-summary-table thead th.agla-age-band) {
  z-index: 3;
}

/* M/F sub-columns sit at slightly lower visual weight than Combined. */
:deep(.agla-summary-table .agla-mf-cell) {
  background-color: rgba(0, 63, 88, 0.04);
  color: rgba(var(--v-theme-on-surface), 0.78);
}

/* Combined column is the headline value — bold and brand-tinted. */
:deep(.agla-summary-table .agla-combined-cell) {
  font-weight: 600;
  background-color: rgba(0, 63, 88, 0.08);
}
:deep(.agla-summary-table thead th.agla-combined-cell) {
  color: #003f58;
}

/* Strong divider on the right edge of every group except the last. */
:deep(.agla-summary-table .agla-group-divider) {
  border-right: 2px solid rgba(0, 63, 88, 0.35);
}
</style>
