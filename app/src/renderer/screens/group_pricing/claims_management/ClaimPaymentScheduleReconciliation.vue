<template>
  <v-progress-linear
    v-if="loadingRecon"
    indeterminate
    color="deep-purple"
    class="mb-2"
  />

  <div v-if="reconSummary" class="d-flex gap-2 mb-3 flex-wrap">
    <v-chip color="default" variant="tonal">
      Total: {{ reconSummary.total_transactions }}
    </v-chip>
    <v-chip color="success" variant="tonal">
      Paid: {{ reconSummary.paid }} ({{
        formatCurrency(reconSummary.total_paid)
      }})
    </v-chip>
    <v-chip color="error" variant="tonal">
      Failed: {{ reconSummary.failed }} ({{
        formatCurrency(reconSummary.total_failed)
      }})
    </v-chip>
    <v-chip color="orange" variant="tonal">
      Unmatched: {{ reconSummary.unmatched }}
    </v-chip>
  </div>

  <empty-state
    v-if="!loadingRecon && reconResults.length === 0"
    icon="mdi-scale-balance"
    title="No reconciliation data yet"
    message="Upload a bank response file to reconcile this payment run."
  />
  <v-table v-else density="compact" class="border rounded mb-3">
    <thead>
      <tr>
        <th>Claim #</th>
        <th>Account</th>
        <th class="text-right">Amount</th>
        <th>Status</th>
        <th>Failure Reason</th>
        <th>Bank Ref</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="r in reconResults" :key="r.id">
        <td>{{ r.claim_number || '—' }}</td>
        <td>{{ r.account_number }}</td>
        <td class="text-right">{{ formatCurrency(r.amount) }}</td>
        <td>
          <v-chip
            :color="reconStatusColor(r.status)"
            size="x-small"
            label
            variant="flat"
          >
            {{ r.status }}
          </v-chip>
        </td>
        <td class="text-caption">{{ r.failure_reason || '—' }}</td>
        <td class="text-caption">{{ r.bank_reference || '—' }}</td>
      </tr>
    </tbody>
  </v-table>

  <v-btn
    v-if="
      hasPermission('claims_pay:retry_failed') &&
      reconResults.some((r) => r.status === 'failed')
    "
    color="orange"
    size="small"
    variant="outlined"
    prepend-icon="mdi-refresh"
    :loading="retrying"
    @click="retryFailed"
  >
    Retry Failed Payments
  </v-btn>
</template>

<script setup lang="ts">
import { inject, onMounted } from 'vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import { PAYMENT_SCHEDULE_CONTEXT } from './payment_schedule_context'

const ctx = inject(PAYMENT_SCHEDULE_CONTEXT)
if (!ctx) {
  throw new Error(
    'ClaimPaymentScheduleReconciliation must be rendered inside ClaimPaymentScheduleLayout'
  )
}

const {
  reconResults,
  reconSummary,
  loadingRecon,
  retrying,
  loadReconData,
  retryFailed,
  reconStatusColor,
  formatCurrency,
  hasPermission,
  schedule
} = ctx

onMounted(async () => {
  if (schedule.value && reconResults.value.length === 0) {
    await loadReconData()
  }
})
</script>
