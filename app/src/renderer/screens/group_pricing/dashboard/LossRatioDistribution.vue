<template>
  <v-row>
    <v-col cols="12" md="7">
      <v-card variant="outlined" class="pa-3" style="height: 100%">
        <div class="d-flex align-center mb-2">
          <div class="text-subtitle-1 font-weight-bold flex-grow-1">
            Loss-Ratio Distribution
          </div>
          <v-tooltip location="top" max-width="320">
            <template #activator="{ props: tipProps }">
              <v-icon
                v-bind="tipProps"
                size="small"
                icon="mdi-information-outline"
                class="ml-1 mr-1"
                style="cursor: help"
              />
            </template>
            <span>
              Schemes bucketed by Inception-to-date Actual Loss Ratio (claims ÷
              earned premium). Schemes shifting right are consuming more premium
              in claims than originally priced.
            </span>
          </v-tooltip>
          <ChartMenu
            :chart-ref="histogramChart"
            title="Loss-Ratio Distribution"
            :data="buckets"
          />
        </div>
        <div class="text-caption text-medium-emphasis mb-2">
          Number of in-force schemes in each ITD ALR band.
        </div>
        <ag-charts
          v-if="histogramOptions"
          ref="histogramChart"
          :options="histogramOptions"
        />
      </v-card>
    </v-col>
    <v-col cols="12" md="5">
      <v-card variant="outlined" class="pa-3" style="height: 100%">
        <div class="d-flex align-center mb-2">
          <div class="text-subtitle-1 font-weight-bold flex-grow-1">
            Top 10 Worst-Performing Schemes
          </div>
          <v-tooltip location="top" max-width="320">
            <template #activator="{ props: tipProps }">
              <v-icon
                v-bind="tipProps"
                size="small"
                icon="mdi-information-outline"
                class="ml-1"
                style="cursor: help"
              />
            </template>
            <span>
              Top 10 schemes ranked by ITD ALR (highest first). Ties broken by
              ITD claims paid. Click a row to drill through to the scheme.
            </span>
          </v-tooltip>
        </div>
        <v-table density="compact">
          <thead>
            <tr>
              <th>Scheme</th>
              <th class="text-right">ALR %</th>
              <th class="text-right">ITD Claims</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in worst"
              :key="row.scheme_id"
              style="cursor: pointer"
              @click="goToScheme(row.scheme_id)"
            >
              <td>{{ row.scheme_name }}</td>
              <td class="text-right">
                <span :style="alrStyle(row.actual_loss_ratio)">
                  {{ fmtPct(row.actual_loss_ratio) }}
                </span>
              </td>
              <td class="text-right">{{ fmtCurrency(row.itd_claims_paid) }}</td>
            </tr>
            <tr v-if="worst.length === 0">
              <td colspan="3" class="text-center text-medium-emphasis">
                No data
              </td>
            </tr>
          </tbody>
        </v-table>
      </v-card>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { AgCharts } from 'ag-charts-vue3'
import ChartMenu from '@/renderer/components/charts/ChartMenu.vue'

const histogramChart: any = ref(null)

interface LossRatioBucket {
  label: string
  scheme_count: number
  premium: number
}
interface WorstRow {
  scheme_id: number
  scheme_name: string
  actual_loss_ratio: number | null
  itd_claims_paid: number
}

const props = defineProps<{
  buckets: LossRatioBucket[]
  worst: WorstRow[]
}>()

const router = useRouter()

const fmtCurrency = (v: number | null | undefined) =>
  v == null
    ? '—'
    : new Intl.NumberFormat('en-ZA', {
        style: 'currency',
        currency: 'ZAR',
        maximumFractionDigits: 0
      }).format(v)
const fmtPct = (v: number | null | undefined) =>
  v == null ? '—' : `${v.toFixed(1)}%`

const bandColor = (label: string) => {
  switch (label) {
    case '<50%':
      return '#2e7d32'
    case '50-80%':
      return '#9ccc65'
    case '80-100%':
      return '#fbc02d'
    case '100-150%':
      return '#ef6c00'
    case '>150%':
      return '#c62828'
    default:
      return '#90a4ae'
  }
}

const alrStyle = (v: number | null) => {
  if (v == null) return ''
  if (v > 100) return 'color:#c62828;font-weight:600'
  if (v > 80) return 'color:#ef6c00;font-weight:600'
  return ''
}

const goToScheme = (id: number) =>
  router.push({ name: 'group-pricing-schemes-detail', params: { id } })

const histogramOptions = computed<any>(() => ({
  data: props.buckets.map((b) => ({
    label: b.label,
    count: b.scheme_count,
    color: bandColor(b.label)
  })),
  background: { fill: 'transparent' },
  height: 280,
  series: [
    {
      type: 'bar',
      xKey: 'label',
      yKey: 'count',
      yName: 'Schemes',
      itemStyler: (p: any) => ({ fill: p.datum.color })
    }
  ],
  axes: [
    {
      type: 'category',
      position: 'bottom',
      title: { text: 'ITD ALR band' }
    },
    { type: 'number', position: 'left', title: { text: 'Schemes' } }
  ],
  legend: { enabled: false }
}))
</script>
