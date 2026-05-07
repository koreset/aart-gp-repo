<template>
  <v-container fluid class="pa-6">
    <v-row class="align-center mb-2">
      <v-col cols="auto">
        <v-btn variant="text" size="small" @click="goBack">
          <v-icon start>mdi-arrow-left</v-icon> Back to dashboard
        </v-btn>
      </v-col>
      <v-col>
        <h2 class="text-h5 font-weight-bold mb-0">
          {{ broker?.name || 'Broker' }}
        </h2>
        <span class="text-caption text-grey">
          Performance for {{ year }}
        </span>
      </v-col>
    </v-row>

    <v-row v-if="loadingHeader">
      <v-col><v-progress-linear indeterminate /></v-col>
    </v-row>

    <!-- Broker info card -->
    <v-row v-if="broker">
      <v-col cols="12">
        <v-card variant="outlined">
          <v-card-text>
            <v-row dense>
              <v-col cols="12" md="3">
                <div class="text-caption text-grey">Contact email</div>
                <div>{{ broker.contact_email || '—' }}</div>
              </v-col>
              <v-col cols="12" md="3">
                <div class="text-caption text-grey">Contact number</div>
                <div>{{ broker.contact_number || '—' }}</div>
              </v-col>
              <v-col cols="12" md="2">
                <div class="text-caption text-grey">FSP number</div>
                <div>{{ broker.fsp_number || '—' }}</div>
              </v-col>
              <v-col cols="12" md="2">
                <div class="text-caption text-grey">FSP category</div>
                <div>{{ broker.fsp_category || '—' }}</div>
              </v-col>
              <v-col cols="12" md="2">
                <div class="text-caption text-grey">Binder ref</div>
                <div>{{ broker.binder_agreement_ref || '—' }}</div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Summary metrics -->
    <v-row v-if="summary" class="mt-2">
      <v-col v-for="card in summaryCards" :key="card.label" cols="6" md="2">
        <v-card variant="outlined">
          <v-card-text class="py-3">
            <div class="text-caption text-grey">{{ card.label }}</div>
            <div class="text-h6 font-weight-bold">{{ card.value }}</div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Exposure tabs (lazy-loaded) -->
    <v-row class="mt-2">
      <v-col>
        <v-card variant="outlined">
          <v-tabs v-model="activeTab" color="primary" density="compact">
            <v-tab value="age">Age Band</v-tab>
            <v-tab value="sum_assured">Sum Assured</v-tab>
            <v-tab value="occupation_class">Occupation Class</v-tab>
          </v-tabs>
          <v-card-text>
            <div v-if="activeTabState.loading">
              <v-progress-linear indeterminate />
            </div>
            <div
              v-else-if="!activeTabState.buckets.length"
              class="text-grey text-center py-6"
            >
              No exposure data for this dimension.
            </div>
            <div v-else>
              <ag-charts :options="chartOptions" />
              <v-table density="compact" class="mt-4">
                <thead>
                  <tr>
                    <th class="text-left">{{ activeTabHeader }}</th>
                    <th class="text-right">Records</th>
                    <th class="text-right">Total Sum Assured</th>
                    <th class="text-right">Male</th>
                    <th class="text-right">Female</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="b in activeTabState.buckets" :key="b.label">
                    <td>{{ b.label }}</td>
                    <td class="text-right">{{ b.record_count }}</td>
                    <td class="text-right">R{{ formatNumber(b.total_sum_assured) }}</td>
                    <td class="text-right">R{{ formatNumber(b.male_sum_assured) }}</td>
                    <td class="text-right">R{{ formatNumber(b.female_sum_assured) }}</td>
                  </tr>
                </tbody>
              </v-table>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { AgCharts } from 'ag-charts-vue3'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useErrorHandler } from '@/renderer/composables/useErrorHandler'

type Dimension = 'age' | 'sum_assured' | 'occupation_class'

interface ExposureBucket {
  label: string
  sort_order: number
  record_count: number
  total_sum_assured: number
  male_sum_assured: number
  female_sum_assured: number
}

interface BrokerSummary {
  broker_id: number
  year: number
  total_quotes: number
  accepted_quotes: number
  conversion_rate: number
  total_premium: number
  in_force_premium: number
  approved_claims: number
  expected_claims: number
  expected_loss_ratio: number
  actual_loss_ratio: number
}

const props = defineProps<{
  brokerId: string | number
  year: string | number
}>()

const router = useRouter()
const { handleApiError } = useErrorHandler()

const brokerIdNum = computed(() => Number(props.brokerId))
const yearNum = computed(() => Number(props.year))

const broker = ref<any>(null)
const summary = ref<BrokerSummary | null>(null)
const loadingHeader = ref(true)

const tabState = reactive<
  Record<Dimension, { loaded: boolean; loading: boolean; buckets: ExposureBucket[] }>
>({
  age: { loaded: false, loading: false, buckets: [] },
  sum_assured: { loaded: false, loading: false, buckets: [] },
  occupation_class: { loaded: false, loading: false, buckets: [] }
})

const activeTab = ref<Dimension>('age')
const activeTabState = computed(() => tabState[activeTab.value])
const activeTabHeader = computed(() => {
  switch (activeTab.value) {
    case 'age':
      return 'Age band'
    case 'sum_assured':
      return 'Sum assured band'
    case 'occupation_class':
      return 'Occupation class'
    default:
      return ''
  }
})

const formatNumber = (num: number | null | undefined) => {
  if (num == null) return '0'
  if (num >= 1_000_000) return (num / 1_000_000).toFixed(1).replace(/\.0$/, '') + 'M'
  if (num >= 1_000) return (num / 1_000).toFixed(1).replace(/\.0$/, '') + 'K'
  return num.toString()
}

const summaryCards = computed(() => {
  if (!summary.value) return []
  return [
    { label: 'Total quotes', value: summary.value.total_quotes },
    { label: 'Accepted', value: summary.value.accepted_quotes },
    {
      label: 'Conversion',
      value: `${summary.value.conversion_rate}%`
    },
    {
      label: 'Total premium',
      value: `R${formatNumber(summary.value.total_premium)}`
    },
    {
      label: 'Expected LR',
      value: `${summary.value.expected_loss_ratio}%`
    },
    {
      label: 'Actual LR',
      value: `${summary.value.actual_loss_ratio}%`
    }
  ]
})

const chartOptions = computed<any>(() => ({
  data: activeTabState.value.buckets,
  series: [
    {
      type: 'bar',
      xKey: 'label',
      yKey: 'total_sum_assured',
      yName: 'Total sum assured',
      tooltip: {
        renderer: (params: any) =>
          `<div style="padding: 6px 8px;"><b>${params.datum.label}</b><br/>R${formatNumber(params.datum.total_sum_assured)}</div>`
      }
    }
  ],
  axes: [
    { type: 'category', position: 'bottom' },
    {
      type: 'number',
      position: 'left',
      label: { formatter: (p: any) => `R${formatNumber(p.value)}` }
    }
  ],
  background: { fill: 'transparent' },
  height: 320
}))

const loadHeader = async () => {
  loadingHeader.value = true
  try {
    const [b, s] = await Promise.all([
      GroupPricingService.getBroker(brokerIdNum.value),
      GroupPricingService.getBrokerPerformanceSummary(
        brokerIdNum.value,
        yearNum.value
      )
    ])
    broker.value = b.data
    summary.value = s.data
  } catch (err) {
    handleApiError(err)
  } finally {
    loadingHeader.value = false
  }
}

const loadDimension = async (dim: Dimension) => {
  if (tabState[dim].loaded || tabState[dim].loading) return
  tabState[dim].loading = true
  try {
    const res = await GroupPricingService.getBrokerExposureBreakdown(
      brokerIdNum.value,
      yearNum.value,
      dim
    )
    tabState[dim].buckets = res.data || []
    tabState[dim].loaded = true
  } catch (err) {
    handleApiError(err)
  } finally {
    tabState[dim].loading = false
  }
}

const goBack = () => {
  router.push({ name: 'group-pricing-dashboard' })
}

watch(activeTab, (dim) => loadDimension(dim), { immediate: false })

onMounted(async () => {
  await loadHeader()
  await loadDimension('age')
})
</script>
