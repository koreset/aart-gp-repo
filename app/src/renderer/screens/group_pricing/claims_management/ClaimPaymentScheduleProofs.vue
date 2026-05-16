<template>
  <empty-state
    v-if="
      !schedule?.proof_of_payments || schedule.proof_of_payments.length === 0
    "
    icon="mdi-receipt-text-outline"
    title="No proof of payment uploaded yet"
    message="Upload proof of payment to confirm this schedule and mark all claims as Paid."
  />
  <v-list v-else density="compact" class="border rounded">
    <v-list-item
      v-for="proof in schedule.proof_of_payments"
      :key="proof.id"
      :subtitle="`Uploaded by ${proof.uploaded_by} on ${formatDate(proof.uploaded_at)}${proof.notes ? ' — ' + proof.notes : ''}`"
      :title="proof.file_name"
    >
      <template #prepend>
        <v-icon color="success">mdi-receipt-text-check</v-icon>
      </template>
      <template #append>
        <v-btn
          size="x-small"
          variant="text"
          color="primary"
          icon="mdi-download"
          :loading="downloadingProof === proof.id"
          @click="downloadProof(proof)"
        />
      </template>
    </v-list-item>
  </v-list>
</template>

<script setup lang="ts">
import { inject } from 'vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import { PAYMENT_SCHEDULE_CONTEXT } from './payment_schedule_context'

const ctx = inject(PAYMENT_SCHEDULE_CONTEXT)
if (!ctx) {
  throw new Error(
    'ClaimPaymentScheduleProofs must be rendered inside ClaimPaymentScheduleLayout'
  )
}

const { schedule, downloadingProof, downloadProof, formatDate } = ctx
</script>
