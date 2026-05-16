<template>
  <v-card variant="outlined" rounded="lg" class="pa-3">
    <div class="d-flex align-center mb-3">
      <div class="text-subtitle-1 font-weight-medium">Finance Queries</div>
      <v-spacer />
      <v-btn
        size="small"
        variant="text"
        prepend-icon="mdi-refresh"
        :loading="loadingQueries"
        @click="loadQueries"
      >
        Refresh
      </v-btn>
    </div>

    <v-data-table
      v-if="queries.length > 0"
      :headers="headers"
      :items="queries"
      :loading="loadingQueries"
      density="compact"
      hover
    >
      <template #[`item.outcome`]="{ item }">
        <v-chip
          :color="outcomeColor(item.outcome)"
          size="x-small"
          variant="flat"
          label
        >
          {{ item.outcome }}
        </v-chip>
      </template>

      <template #[`item.raised_at`]="{ item }">
        {{ formatDate(item.raised_at) }}
      </template>

      <template #[`item.reason_code`]="{ item }">
        <span class="text-caption font-weight-medium">{{
          item.reason_code
        }}</span>
      </template>
    </v-data-table>

    <empty-state
      v-else-if="!loadingQueries"
      icon="mdi-comment-check-outline"
      title="No queries raised"
      message="Finance has not queried any line items on this schedule."
    />

    <div v-else class="text-center py-6">
      <v-progress-circular indeterminate color="primary" />
    </div>
  </v-card>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import { PAYMENT_SCHEDULE_CONTEXT } from './payment_schedule_context'

const ctx = inject(PAYMENT_SCHEDULE_CONTEXT)
if (!ctx) {
  throw new Error('Payment schedule context not provided')
}
const { queries, loadingQueries, loadQueries, formatDate } = ctx

const headers = [
  { title: 'Claim #', key: 'claim_number', sortable: true },
  { title: 'Reason code', key: 'reason_code', sortable: true },
  { title: 'Notes', key: 'notes', sortable: false },
  { title: 'Outcome', key: 'outcome', sortable: true },
  { title: 'Raised by', key: 'raised_by', sortable: true },
  { title: 'Raised at', key: 'raised_at', sortable: true }
]

function outcomeColor(outcome: string) {
  const map: Record<string, string> = {
    open: 'orange',
    queried: 'orange',
    rejected: 'error',
    resolved: 'success',
    cancelled: 'grey'
  }
  return map[outcome] ?? 'default'
}

onMounted(loadQueries)
</script>
