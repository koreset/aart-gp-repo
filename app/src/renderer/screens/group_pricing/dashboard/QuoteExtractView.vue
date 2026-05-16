<template>
  <v-container fluid class="quote-extract pa-6">
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
        <h2 class="text-h5 font-weight-bold mb-0">Quote extract</h2>
        <span class="text-caption text-medium-emphasis"
          >Filter, browse and export the full quote book.</span
        >
      </v-col>
      <v-col cols="auto">
        <v-btn
          color="primary"
          variant="tonal"
          prepend-icon="mdi-file-excel-outline"
          :loading="exporting"
          :disabled="!total"
          @click="downloadXlsx"
        >
          Export Excel
        </v-btn>
      </v-col>
    </v-row>

    <!-- Filters -->
    <v-card elevation="1" class="rounded-lg pa-4 mb-4">
      <v-row dense align="end">
        <v-col cols="12" sm="6" md="3">
          <v-combobox
            v-model="filter.created_by"
            label="Created by"
            density="compact"
            multiple
            chips
            closable-chips
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-combobox
            v-model="filter.reviewer"
            label="Reviewer"
            density="compact"
            multiple
            chips
            closable-chips
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-select
            v-model="filter.status"
            :items="statusOptions"
            label="Status"
            density="compact"
            multiple
            chips
            closable-chips
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-combobox
            v-model="filter.region"
            label="Region"
            density="compact"
            multiple
            chips
            closable-chips
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-select
            v-model="filter.quote_type"
            :items="['New Business', 'Renewal']"
            label="Quote type"
            density="compact"
            multiple
            chips
            closable-chips
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-select
            v-model="filter.distribution_channel"
            :items="['broker', 'direct', 'binder', 'tied_agent']"
            label="Channel"
            density="compact"
            multiple
            chips
            closable-chips
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6" md="2">
          <v-text-field
            v-model.number="filter.min_annual_premium"
            type="number"
            label="Min premium"
            density="compact"
            hide-details
          />
        </v-col>
        <v-col cols="12" sm="6" md="2">
          <v-text-field
            v-model.number="filter.max_annual_premium"
            type="number"
            label="Max premium"
            density="compact"
            hide-details
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-text-field
            v-model="filter.from"
            type="date"
            label="Created from"
            density="compact"
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-text-field
            v-model="filter.to"
            type="date"
            label="Created to"
            density="compact"
            hide-details
            clearable
          />
        </v-col>
        <v-col cols="auto" class="ms-auto">
          <v-btn
            color="primary"
            density="compact"
            size="small"
            prepend-icon="mdi-magnify"
            :loading="loading"
            @click="applyFilter"
          >
            Apply filter
          </v-btn>
        </v-col>
      </v-row>
    </v-card>

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
      <div class="d-flex align-center pa-3">
        <span class="text-subtitle-2 font-weight-medium me-3"
          >{{ total.toLocaleString() }} quotes</span
        >
        <v-spacer />
        <v-pagination
          v-model="filter.page"
          :length="pageCount"
          :total-visible="5"
          density="compact"
          @update:model-value="refresh"
        />
      </div>
      <v-divider />
      <v-table density="compact" hover class="extract-table">
        <thead>
          <tr>
            <th>Quote</th>
            <th>Type</th>
            <th>Scheme</th>
            <th>Industry</th>
            <th>Regions</th>
            <th>Status</th>
            <th>Created by</th>
            <th>Reviewer</th>
            <th>Approved by</th>
            <th class="text-right">Created</th>
            <th class="text-right">Submitted</th>
            <th class="text-right">Approved</th>
            <th class="text-right">Accepted</th>
            <th class="text-right">Annual premium</th>
            <th class="text-right">Members</th>
            <th class="text-right">Cycle (hrs)</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in rows" :key="row.id">
            <td class="font-weight-medium">{{ row.quote_name }}</td>
            <td>{{ row.quote_type }}</td>
            <td>{{ row.scheme_name }}</td>
            <td>{{ row.industry }}</td>
            <td>{{ row.regions }}</td>
            <td>
              <v-chip
                size="x-small"
                :color="statusColor(row.status)"
                variant="tonal"
                >{{ row.status }}</v-chip
              >
            </td>
            <td>{{ row.created_by }}</td>
            <td>{{ row.reviewer }}</td>
            <td>{{ row.approved_by }}</td>
            <td class="text-right">{{ formatDate(row.creation_date) }}</td>
            <td class="text-right">{{ formatDate(row.submitted_at) }}</td>
            <td class="text-right">{{ formatDate(row.approved_at) }}</td>
            <td class="text-right">{{ formatDate(row.accepted_at) }}</td>
            <td class="text-right">{{ formatCurrency(row.annual_premium) }}</td>
            <td class="text-right">{{ row.member_count }}</td>
            <td class="text-right">{{
              row.cycle_hours != null ? row.cycle_hours.toFixed(1) : '—'
            }}</td>
          </tr>
          <tr v-if="!rows.length && !loading">
            <td colspan="16" class="text-center text-medium-emphasis py-6"
              >No quotes match the active filter.</td
            >
          </tr>
        </tbody>
      </v-table>
      <v-progress-linear v-if="loading" indeterminate color="primary" />
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'

import { useQuoteExtract } from '@/renderer/composables/useQuoteDashboard'

const {
  filter,
  rows,
  total,
  loading,
  error,
  exporting,
  refresh,
  downloadXlsx
} = useQuoteExtract()

const statusOptions = [
  'draft',
  'in_progress',
  'submitted',
  'pending_review',
  'approved',
  'rejected',
  'accepted',
  'in_force',
  'out_of_force',
  'cancelled',
  'expired',
  'lapsed',
  'not_taken_up'
]

onMounted(() => refresh())

function applyFilter() {
  filter.value.page = 1
  refresh()
}

const pageCount = computed(() => {
  const size = filter.value.page_size || 50
  return Math.max(1, Math.ceil(total.value / size))
})

function statusColor(status: string): string {
  switch (status) {
    case 'accepted':
    case 'in_force':
    case 'approved':
      return 'success'
    case 'rejected':
    case 'cancelled':
    case 'lapsed':
    case 'expired':
      return 'error'
    case 'submitted':
    case 'pending_review':
      return 'warning'
    default:
      return 'default'
  }
}

function formatDate(value?: string | null): string {
  if (!value) return '—'
  const d = new Date(value)
  if (isNaN(d.getTime())) return '—'
  return d.toLocaleDateString()
}

function formatCurrency(value: number | null | undefined): string {
  if (!value || !isFinite(value)) return 'R0'
  if (value >= 1_000_000) return 'R' + (value / 1_000_000).toFixed(2) + 'M'
  if (value >= 1_000) return 'R' + (value / 1_000).toFixed(1) + 'K'
  return 'R' + Math.round(value).toLocaleString()
}
</script>

<style scoped>
.quote-extract {
  --qp-bg: #f6f8fb;
  --qp-surface: #ffffff;
  --qp-border: #e2e8f0;
  --qp-text: #0f172a;
  --qp-text-muted: #64748b;
  --qp-primary: #2563eb;

  background: var(--qp-bg);
  min-height: 100%;
  color: var(--qp-text);
}
.quote-extract :deep(.text-h5) {
  color: var(--qp-text);
  letter-spacing: -0.01em;
}
.quote-extract :deep(.v-card) {
  background: var(--qp-surface);
  border: 1px solid var(--qp-border);
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
}
.quote-extract :deep(.v-table) {
  background: transparent;
  --v-theme-surface: transparent;
}
.quote-extract :deep(.v-table thead th) {
  color: var(--qp-text-muted) !important;
  font-size: 0.6875rem !important;
  font-weight: 600 !important;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  background: #fafbfd;
  border-bottom: 1px solid var(--qp-border) !important;
}
.quote-extract :deep(.v-table tbody td) {
  color: var(--qp-text);
  font-size: 0.8125rem;
  border-bottom: 1px solid #f1f5f9 !important;
}
.quote-extract :deep(.v-table tbody tr:hover) {
  background: #f8fafc !important;
}
.quote-extract :deep(.v-chip.v-chip--variant-tonal) {
  font-weight: 600;
  font-size: 0.6875rem;
  letter-spacing: 0.02em;
}
.extract-table {
  overflow-x: auto;
}
.extract-table :deep(td:nth-child(n + 10)),
.extract-table :deep(th:nth-child(n + 10)) {
  font-variant-numeric: tabular-nums;
}
.extract-table td,
.extract-table th {
  white-space: nowrap;
}
</style>
