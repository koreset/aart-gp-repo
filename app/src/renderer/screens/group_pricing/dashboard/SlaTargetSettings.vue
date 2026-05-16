<template>
  <v-container fluid class="sla-targets pa-6">
    <v-row class="align-center mb-4" no-gutters>
      <v-col cols="auto" class="me-3">
        <v-btn
          variant="tonal"
          size="small"
          color="primary"
          prepend-icon="mdi-arrow-left"
          :to="{ name: 'group-pricing-quote-performance' }"
        >
          Back
        </v-btn>
      </v-col>
      <v-col>
        <h2 class="text-h5 font-weight-bold mb-0">SLA targets</h2>
        <span class="text-caption text-medium-emphasis">
          Per-transition turnaround targets used to score quote performance.
        </span>
      </v-col>
      <v-col cols="auto">
        <v-btn color="primary" prepend-icon="mdi-plus" @click="openNewDialog">
          New target
        </v-btn>
      </v-col>
    </v-row>

    <v-alert
      v-if="error"
      type="error"
      density="compact"
      class="mb-4"
      closable
      @click:close="error = null"
    >
      {{ error }}
    </v-alert>

    <v-card elevation="1" class="rounded-lg">
      <v-table density="compact" hover>
        <thead>
          <tr>
            <th>From</th>
            <th>To</th>
            <th>Quote type</th>
            <th class="text-right">Target (hrs)</th>
            <th class="text-right">Warning %</th>
            <th class="text-center">Active</th>
            <th>Updated by</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in targets" :key="t.id">
            <td>{{ t.from_status }}</td>
            <td>{{ t.to_status }}</td>
            <td>{{ t.quote_type || 'any' }}</td>
            <td class="text-right">{{ t.target_hours }}</td>
            <td class="text-right"
              >{{ ((t.warning_pct_of_sla ?? 0.8) * 100).toFixed(0) }}%</td
            >
            <td class="text-center">
              <v-icon v-if="t.active" color="success" size="small"
                >mdi-check-circle</v-icon
              >
              <v-icon v-else color="grey" size="small">mdi-cancel</v-icon>
            </td>
            <td class="text-medium-emphasis">{{ t.updated_by }}</td>
            <td class="text-right">
              <v-btn
                variant="text"
                size="small"
                icon="mdi-pencil-outline"
                @click="openEditDialog(t)"
              />
              <v-btn
                variant="text"
                size="small"
                color="error"
                icon="mdi-delete-outline"
                :disabled="!t.id"
                @click="confirmDelete(t)"
              />
            </td>
          </tr>
          <tr v-if="!targets.length && !loading">
            <td colspan="8" class="text-center text-medium-emphasis py-6"
              >No SLA targets configured.</td
            >
          </tr>
        </tbody>
      </v-table>
      <v-progress-linear v-if="loading" indeterminate color="primary" />
    </v-card>

    <!-- Edit / create dialog -->
    <v-dialog v-model="dialogOpen" max-width="540">
      <v-card>
        <v-card-title>
          {{ editing.id ? 'Edit SLA target' : 'New SLA target' }}
        </v-card-title>
        <v-card-text>
          <v-row dense>
            <v-col cols="12" md="6">
              <v-select
                v-model="editing.from_status"
                :items="statusOptions"
                label="From status"
                density="compact"
                hide-details
                :disabled="!!editing.id"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="editing.to_status"
                :items="statusOptions"
                label="To status"
                density="compact"
                hide-details
                :disabled="!!editing.id"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-select
                v-model="editing.quote_type"
                :items="['', 'New Business', 'Renewal']"
                label="Quote type (blank = any)"
                density="compact"
                hide-details
                :disabled="!!editing.id"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="editing.target_hours"
                type="number"
                label="Target (hours)"
                density="compact"
                hide-details
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model.number="warningPct"
                type="number"
                label="Warning threshold (%)"
                density="compact"
                hint="Default 80%"
                persistent-hint
                min="0"
                max="100"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-switch
                v-model="editing.active"
                color="primary"
                label="Active"
                density="compact"
                hide-details
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="dialogOpen = false">Cancel</v-btn>
          <v-btn color="primary" :loading="saving" @click="saveEditing"
            >Save</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Delete confirmation -->
    <v-dialog v-model="deleteDialogOpen" max-width="420">
      <v-card>
        <v-card-title>Deactivate SLA target?</v-card-title>
        <v-card-text>
          The target will be flagged inactive. Historical breach calculations
          against it remain unchanged.
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="deleteDialogOpen = false">Cancel</v-btn>
          <v-btn color="error" :loading="deleting" @click="performDelete"
            >Deactivate</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

import type { QuoteSlaTarget } from '@/renderer/api/QuoteDashboardService'
import { useSlaTargets } from '@/renderer/composables/useQuoteDashboard'

const { targets, loading, error, refresh, save, remove } = useSlaTargets()

const statusOptions = [
  'draft',
  'in_progress',
  'submitted',
  'pending_review',
  'approved',
  'rejected',
  'accepted',
  'in_force'
]

const dialogOpen = ref(false)
const deleteDialogOpen = ref(false)
const saving = ref(false)
const deleting = ref(false)
const toDelete = ref<QuoteSlaTarget | null>(null)
const editing = ref<QuoteSlaTarget>(blankTarget())

const warningPct = computed({
  get: () => Math.round((editing.value.warning_pct_of_sla ?? 0.8) * 100),
  set: (v) => {
    editing.value.warning_pct_of_sla = (Number(v) || 0) / 100
  }
})

onMounted(() => refresh())

function blankTarget(): QuoteSlaTarget {
  return {
    from_status: 'submitted',
    to_status: 'approved',
    target_hours: 24,
    warning_pct_of_sla: 0.8,
    quote_type: '',
    active: true
  }
}

function openNewDialog() {
  editing.value = blankTarget()
  dialogOpen.value = true
}

function openEditDialog(t: QuoteSlaTarget) {
  editing.value = { ...t, warning_pct_of_sla: t.warning_pct_of_sla ?? 0.8 }
  dialogOpen.value = true
}

async function saveEditing() {
  saving.value = true
  try {
    await save(editing.value)
    dialogOpen.value = false
  } catch {
    // error already surfaced via composable
  } finally {
    saving.value = false
  }
}

function confirmDelete(t: QuoteSlaTarget) {
  toDelete.value = t
  deleteDialogOpen.value = true
}

async function performDelete() {
  if (!toDelete.value?.id) {
    deleteDialogOpen.value = false
    return
  }
  deleting.value = true
  try {
    await remove(toDelete.value.id)
    deleteDialogOpen.value = false
    toDelete.value = null
  } catch {
    // error already surfaced via composable
  } finally {
    deleting.value = false
  }
}
</script>

<style scoped>
.sla-targets {
  --qp-bg: #f6f8fb;
  --qp-surface: #ffffff;
  --qp-border: #e2e8f0;
  --qp-text: #0f172a;
  --qp-text-muted: #64748b;

  background: var(--qp-bg);
  min-height: 100%;
  color: var(--qp-text);
}
.sla-targets :deep(.text-h5) {
  color: var(--qp-text);
  letter-spacing: -0.01em;
}
.sla-targets :deep(.v-card) {
  background: var(--qp-surface);
  border: 1px solid var(--qp-border);
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
}
.sla-targets :deep(.v-table) {
  background: transparent;
  --v-theme-surface: transparent;
}
.sla-targets :deep(.v-table thead th) {
  color: var(--qp-text-muted) !important;
  font-size: 0.6875rem !important;
  font-weight: 600 !important;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  background: #fafbfd;
  border-bottom: 1px solid var(--qp-border) !important;
}
.sla-targets :deep(.v-table tbody td) {
  color: var(--qp-text);
  font-size: 0.8125rem;
  border-bottom: 1px solid #f1f5f9 !important;
  font-variant-numeric: tabular-nums;
}
.sla-targets :deep(.v-table tbody tr:hover) {
  background: #f8fafc !important;
}
</style>
