<template>
  <v-dialog
    :model-value="modelValue"
    persistent
    max-width="800px"
    @update:model-value="(val: boolean) => emit('update:modelValue', val)"
  >
    <v-card rounded="lg">
      <v-card-title class="text-h6 pa-4 pb-2">
        New Payment Schedule
      </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="description"
          label="Description (optional)"
          variant="outlined"
          density="compact"
          class="mb-4"
        />

        <div class="d-flex align-center justify-space-between mb-2">
          <span class="text-subtitle-2">Select Approved Claims</span>
          <v-chip size="small" color="primary" variant="tonal">
            {{ selectedClaimIDs.length }} selected
          </v-chip>
        </div>
        <v-row dense class="mb-2">
          <v-col cols="6">
            <v-text-field
              v-model="claimFilter"
              label="Filter by Claim Number"
              prepend-inner-icon="mdi-magnify"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
          <v-col cols="6">
            <v-text-field
              v-model="benefitFilter"
              label="Filter by Benefit"
              prepend-inner-icon="mdi-filter-outline"
              variant="outlined"
              density="compact"
              clearable
            />
          </v-col>
        </v-row>
        <DataGrid
          :row-data="filteredApprovedClaims"
          :column-defs="claimsColumnDefs"
          row-selection="multiple"
          density="compact"
          :pagination="false"
          @row-selection-changed="onClaimSelectionChanged"
        />
        <div class="text-subtitle-2 text-right mt-2">
          Subtotal: <strong>{{ formatCurrency(selectedTotal) }}</strong>
        </div>
      </v-card-text>
      <v-card-actions class="pa-4 pt-0">
        <v-spacer />
        <v-btn variant="text" @click="cancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="creating"
          :disabled="selectedClaimIDs.length === 0"
          @click="createSchedule"
        >
          Create Schedule
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { ColDef } from 'ag-grid-community'
import DataGrid from '@/renderer/components/tables/DataGrid.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface Claim {
  id: number
  claim_number: string
  member_name: string
  member_id_number: string
  scheme_name: string
  benefit_alias: string
  claim_amount: number
  status: string
}

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'created', schedule: any): void
  (e: 'error', message: string): void
}>()

const description = ref('')
const claimFilter = ref('')
const benefitFilter = ref('')
const approvedClaims = ref<Claim[]>([])
const selectedClaimIDs = ref<number[]>([])
const creating = ref(false)
const loadingApproved = ref(false)

const filteredApprovedClaims = computed(() => {
  let result = approvedClaims.value
  if (claimFilter.value) {
    const q = claimFilter.value.toLowerCase()
    result = result.filter((c) => c.claim_number.toLowerCase().includes(q))
  }
  if (benefitFilter.value) {
    const b = benefitFilter.value.toLowerCase()
    result = result.filter((c) => c.benefit_alias?.toLowerCase().includes(b))
  }
  return result
})

const selectedTotal = computed(() =>
  approvedClaims.value
    .filter((c) => selectedClaimIDs.value.includes(c.id))
    .reduce((sum, c) => sum + c.claim_amount, 0)
)

const claimsColumnDefs: ColDef<Claim>[] = [
  {
    checkboxSelection: true,
    headerCheckboxSelection: true,
    width: 48,
    minWidth: 48,
    maxWidth: 48,
    pinned: 'left',
    suppressMovable: true,
    resizable: false
  },
  {
    headerName: 'Claim #',
    field: 'claim_number',
    sortable: true,
    minWidth: 120
  },
  { headerName: 'Member', field: 'member_name', sortable: true, minWidth: 140 },
  { headerName: 'ID Number', field: 'member_id_number', minWidth: 120 },
  { headerName: 'Scheme', field: 'scheme_name', sortable: true, minWidth: 140 },
  { headerName: 'Benefit', field: 'benefit_alias', minWidth: 130 },
  {
    headerName: 'Amount',
    field: 'claim_amount',
    minWidth: 120,
    type: 'rightAligned',
    valueFormatter: (p) => formatCurrency(p.value)
  },
  {
    headerName: 'Status',
    field: 'status',
    minWidth: 100
  }
]

function formatCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(val ?? 0)
}

function unwrap(res: any) {
  return res?.data?.data ?? res?.data ?? res
}

function onClaimSelectionChanged(rows: Claim[]) {
  selectedClaimIDs.value = rows.map((c) => c.id)
}

async function loadApprovedClaims() {
  loadingApproved.value = true
  try {
    const res = await GroupPricingService.getClaims()
    const all: Claim[] = unwrap(res) ?? []
    approvedClaims.value = all.filter((c) => c.status === 'approved')
  } catch {
    emit('error', 'Failed to load approved claims')
  } finally {
    loadingApproved.value = false
  }
}

function cancel() {
  emit('update:modelValue', false)
}

async function createSchedule() {
  creating.value = true
  try {
    const res = await GroupPricingService.createPaymentSchedule({
      claim_ids: selectedClaimIDs.value,
      description: description.value
    })
    const schedule = unwrap(res)
    emit('created', schedule)
    emit('update:modelValue', false)
  } catch (e: any) {
    const msg =
      e?.response?.data?.error ||
      e?.response?.data ||
      'Failed to create schedule'
    emit('error', typeof msg === 'string' ? msg : 'Failed to create schedule')
  } finally {
    creating.value = false
  }
}

function reset() {
  description.value = ''
  claimFilter.value = ''
  benefitFilter.value = ''
  selectedClaimIDs.value = []
}

watch(
  () => props.modelValue,
  async (open) => {
    if (open) {
      reset()
      await loadApprovedClaims()
    }
  }
)
</script>
