<template>
  <v-card variant="outlined" rounded="lg" class="pa-3">
    <div class="d-flex align-center mb-3">
      <div class="text-subtitle-1 font-weight-medium">Queries & follow-ups</div>
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
        <span class="text-caption font-weight-medium">
          {{ item.reason_code }}
        </span>
      </template>

      <template #[`item.notes`]="{ item }">
        <div class="text-body-2">{{ item.notes }}</div>
        <div
          v-if="item.resolution_notes"
          class="text-caption text-success mt-1"
        >
          <v-icon size="12" icon="mdi-reply" class="mr-1" />
          <strong>{{ item.resolved_by }}:</strong>
          {{ item.resolution_notes }}
        </div>
      </template>

      <template #[`item.actions`]="{ item }">
        <v-btn
          v-if="
            item.outcome === 'open' &&
            hasPermission('claims_pay:finance_review')
          "
          size="x-small"
          variant="text"
          color="primary"
          prepend-icon="mdi-reply"
          @click="openRespond(item)"
        >
          Respond
        </v-btn>
        <span
          v-else-if="item.resolved_at"
          class="text-caption text-medium-emphasis"
        >
          {{ formatDate(item.resolved_at) }}
        </span>
      </template>
    </v-data-table>

    <empty-state
      v-else-if="!loadingQueries"
      icon="mdi-comment-check-outline"
      title="No queries raised"
      message="No queries or follow-ups have been raised on this schedule."
    />

    <div v-else class="text-center py-6">
      <v-progress-circular indeterminate color="primary" />
    </div>

    <!-- Respond dialog -->
    <v-dialog v-model="respondDialog" max-width="520px">
      <v-card rounded="lg">
        <v-card-title class="text-h6 pa-4 pb-2">Respond to query</v-card-title>
        <v-card-text>
          <div v-if="responding" class="mb-3">
            <div class="text-caption text-medium-emphasis">
              Raised by {{ responding.raised_by }} ·
              {{ formatDate(responding.raised_at) }}
            </div>
            <div
              v-if="responding.claim_number"
              class="text-caption text-medium-emphasis"
            >
              Claim {{ responding.claim_number }} ·
              {{ responding.reason_code }}
            </div>
            <v-alert
              variant="tonal"
              color="grey"
              density="compact"
              class="mt-2"
            >
              {{ responding.notes }}
            </v-alert>
          </div>
          <v-textarea
            v-model="responseText"
            variant="outlined"
            density="compact"
            rows="4"
            placeholder="Type your response to claims..."
            :rules="[(v) => !!v?.trim() || 'Response is required']"
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0">
          <v-spacer />
          <v-btn variant="text" @click="respondDialog = false">Cancel</v-btn>
          <v-btn
            color="primary"
            :loading="resolving"
            :disabled="!responseText.trim()"
            @click="submitResponse"
          >
            Send response
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-snackbar v-model="snackbar" :color="snackbarColor" timeout="3500">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-card>
</template>

<script setup lang="ts">
import { inject, onMounted, ref } from 'vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'
import { PAYMENT_SCHEDULE_CONTEXT } from './payment_schedule_context'

interface ScheduleQuery {
  id: number
  claim_number: string
  reason_code: string
  notes: string
  outcome: string
  raised_by: string
  raised_at: string
  resolution_notes?: string
  resolved_by?: string
  resolved_at?: string | null
}

const ctx = inject(PAYMENT_SCHEDULE_CONTEXT)
if (!ctx) {
  throw new Error('Payment schedule context not provided')
}
const { queries, loadingQueries, loadQueries, formatDate } = ctx
const { hasPermission } = usePermissionCheck()

const respondDialog = ref(false)
const responding = ref<ScheduleQuery | null>(null)
const responseText = ref('')
const resolving = ref(false)

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const headers = [
  { title: 'Claim #', key: 'claim_number', sortable: true },
  { title: 'Reason code', key: 'reason_code', sortable: true },
  { title: 'Notes', key: 'notes', sortable: false },
  { title: 'Outcome', key: 'outcome', sortable: true },
  { title: 'Raised by', key: 'raised_by', sortable: true },
  { title: 'Raised at', key: 'raised_at', sortable: true },
  { title: '', key: 'actions', sortable: false, align: 'end' as const }
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

function openRespond(item: ScheduleQuery) {
  responding.value = item
  responseText.value = ''
  respondDialog.value = true
}

function notify(msg: string, color = 'success') {
  snackbarMessage.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

async function submitResponse() {
  if (!responding.value || !responseText.value.trim()) return
  resolving.value = true
  try {
    await GroupPricingService.resolveScheduleQuery(
      responding.value.id,
      responseText.value.trim()
    )
    notify('Response sent. Query marked resolved.')
    respondDialog.value = false
    responding.value = null
    responseText.value = ''
    await loadQueries()
  } catch (e: any) {
    const msg = e?.response?.data?.error || 'Failed to send response'
    notify(typeof msg === 'string' ? msg : 'Failed to send response', 'error')
  } finally {
    resolving.value = false
  }
}

onMounted(loadQueries)
</script>
