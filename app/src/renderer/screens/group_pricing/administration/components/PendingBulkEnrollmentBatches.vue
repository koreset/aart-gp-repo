<template>
  <div>
    <v-row class="mb-3" align="center">
      <v-col cols="12" md="6">
        <v-select
          v-model="selectedScheme"
          :items="schemes"
          item-title="name"
          item-value="id"
          label="Filter by Scheme"
          variant="outlined"
          density="compact"
          clearable
          @update:model-value="reload"
        />
      </v-col>
      <v-col cols="12" md="6" class="d-flex justify-end gap-2">
        <v-btn
          rounded
          size="small"
          color="primary"
          variant="outlined"
          :loading="loading"
          @click="reload"
        >
          Refresh
        </v-btn>
        <v-btn
          rounded
          size="small"
          color="primary"
          prepend-icon="mdi-file-upload"
          @click="goToUpload"
        >
          New Bulk Enrollment
        </v-btn>
      </v-col>
    </v-row>

    <div :style="{ height: gridHeight, width: '100%' }">
      <data-grid
        :column-defs="columnDefs"
        :row-data="batches"
        :pagination="false"
        :loading="loading"
        style="height: 100%; width: 100%"
        @row-double-clicked="openBatch"
      />
    </div>

    <div
      v-if="!loading && batches.length === 0"
      class="text-center py-6 text-grey"
    >
      No bulk enrollment batches pending approval.
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useGridHeight } from '@/renderer/composables/useGridHeight'

interface Props {
  status?: string
}
const props = withDefaults(defineProps<Props>(), {
  status: 'pending_approval'
})

const emit = defineEmits<{
  (e: 'count-changed', count: number): void
}>()

const router = useRouter()
const gridHeight = useGridHeight(420)
const loading = ref(false)
const batches = ref<any[]>([])
const schemes = ref<Array<{ id: number; name: string }>>([])
const selectedScheme = ref<number | null>(null)

const formatDate = (v: any) => {
  if (!v) return ''
  const d = new Date(v)
  return isNaN(d.getTime()) ? v : d.toLocaleString()
}

const columnDefs = [
  {
    headerName: 'Uploaded At',
    field: 'uploaded_at',
    valueFormatter: (p: any) => formatDate(p.value),
    minWidth: 180,
    sort: 'desc'
  },
  { headerName: 'Uploaded By', field: 'uploaded_by', minWidth: 160 },
  { headerName: 'File', field: 'file_name', minWidth: 220 },
  {
    headerName: 'Members',
    field: 'member_count',
    minWidth: 100,
    type: 'numericColumn'
  },
  {
    headerName: 'Valid',
    field: 'valid_count',
    minWidth: 90,
    type: 'numericColumn',
    cellStyle: { color: '#2e7d32' }
  },
  {
    headerName: 'Blocking',
    field: 'blocking_count',
    minWidth: 100,
    type: 'numericColumn',
    cellStyle: (p: any) =>
      p.value > 0 ? { color: '#c62828', fontWeight: 600 } : null
  },
  {
    headerName: 'Soft',
    field: 'soft_error_count',
    minWidth: 90,
    type: 'numericColumn',
    cellStyle: { color: '#ef6c00' }
  },
  {
    headerName: 'Status',
    field: 'status',
    minWidth: 140,
    cellRenderer: (p: any) => {
      const v = p.value || ''
      const colors: Record<string, string> = {
        pending_approval: '#1976d2',
        approved: '#2e7d32',
        rejected: '#c62828',
        cancelled: '#616161'
      }
      const c = colors[v] || '#616161'
      return `<span style="background:${c};color:#fff;padding:2px 8px;border-radius:10px;font-size:11px;">${v.replace('_', ' ')}</span>`
    }
  }
]

const reload = async () => {
  loading.value = true
  try {
    const accumulator: any[] = []
    const schemeIds = selectedScheme.value
      ? [selectedScheme.value]
      : schemes.value.map((s) => s.id)
    for (const id of schemeIds) {
      const res = await GroupPricingService.listBulkEnrollmentBatches(
        id,
        props.status
      )
      const items = res?.data?.batches ?? []
      accumulator.push(...items)
    }
    batches.value = accumulator
    emit('count-changed', accumulator.length)
  } catch (err) {
    console.error('Failed to load batches', err)
    batches.value = []
    emit('count-changed', 0)
  } finally {
    loading.value = false
  }
}

const openBatch = (event: any) => {
  const batchId = event?.data?.id
  if (!batchId) return
  router.push({
    name: 'group-pricing-bulk-enrollment-batch',
    params: { batchId: String(batchId) }
  })
}

const goToUpload = () => {
  router.push({ name: 'group-pricing-bulk-enrollment' })
}

onMounted(async () => {
  try {
    const res = await GroupPricingService.getSchemes()
    schemes.value = res?.data ?? []
  } catch (err) {
    console.error('Failed to load schemes for filter', err)
  }
  reload()
})

defineExpose({ reload })
</script>
