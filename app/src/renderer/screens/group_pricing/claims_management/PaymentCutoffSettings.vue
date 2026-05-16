<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="8">
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center">
              <v-btn
                icon="mdi-arrow-left"
                variant="text"
                class="mr-2"
                @click="goBack"
              />
              <span class="headline">Payment Cut-off Settings</span>
            </div>
          </template>

          <template #default>
            <v-alert
              type="info"
              variant="tonal"
              density="compact"
              icon="mdi-information-outline"
              class="mb-4"
            >
              <div class="text-body-2">
                Configure when the system automatically bundles approved claims
                into a payment schedule. The scheduler ticks every minute and
                runs once per HH:MM cut-off per day. The
                <strong>daily payment limit</strong> is checked at the first
                finance authorisation — schedules whose net would push today's
                authorised total above the limit are blocked.
              </div>
            </v-alert>

            <v-switch
              v-model="form.enabled"
              color="success"
              label="Enable automatic cut-off generation"
              density="compact"
              class="mb-3"
              hide-details
            />

            <v-text-field
              v-model="form.cutoff_times"
              label="Cut-off times (HH:MM, comma-separated)"
              variant="outlined"
              density="compact"
              placeholder="11:00, 15:00"
              hint="At each time, all currently approved claims are bundled into a draft schedule."
              persistent-hint
              class="mb-3"
            />

            <v-text-field
              v-model="form.timezone"
              label="Timezone (IANA)"
              variant="outlined"
              density="compact"
              hint="e.g. Africa/Johannesburg. Cut-off times are interpreted in this timezone."
              persistent-hint
              class="mb-3"
            />

            <v-text-field
              v-model.number="form.daily_payment_limit"
              type="number"
              label="Daily payment limit (ZAR) — 0 for no limit"
              variant="outlined"
              density="compact"
              hint="Sum of NetTotal across schedules that received first finance auth today must not exceed this."
              persistent-hint
              class="mb-4"
            />

            <div class="d-flex gap-2">
              <v-btn
                color="primary"
                :loading="saving"
                :disabled="!canSave"
                @click="save"
              >
                Save settings
              </v-btn>
              <v-btn variant="text" @click="loadConfig"> Reset </v-btn>
            </div>
          </template>
        </base-card>

        <!-- Recent runs -->
        <base-card class="mt-4" :show-actions="false">
          <template #header>
            <div class="d-flex align-center">
              <span class="headline">Recent cut-off runs</span>
              <v-spacer />
              <v-btn
                size="small"
                variant="text"
                prepend-icon="mdi-refresh"
                @click="loadRuns"
              >
                Refresh
              </v-btn>
            </div>
          </template>
          <template #default>
            <v-data-table
              :headers="runHeaders"
              :items="runs"
              :loading="loadingRuns"
              density="compact"
              hover
            >
              <template #[`item.status`]="{ item }">
                <v-chip
                  size="x-small"
                  variant="flat"
                  :color="runColor(item.status)"
                  >{{ item.status }}</v-chip
                >
              </template>
              <template #[`item.scheduled_at`]="{ item }">
                {{ formatDateTime(item.scheduled_at) }}
              </template>
              <template #[`item.total_amount`]="{ item }">
                {{ formatCurrency(item.total_amount) }}
              </template>
            </v-data-table>
          </template>
        </base-card>
      </v-col>

      <v-col cols="12" md="4">
        <base-card :show-actions="false">
          <template #header>
            <span class="headline">Next cut-off</span>
          </template>
          <template #default>
            <div v-if="next" class="text-center py-4">
              <div class="text-overline text-medium-emphasis">Scheduled</div>
              <div class="text-h5 font-weight-bold">{{ next }}</div>
              <v-btn
                v-if="hasPermission('claims_pay:run_cutoff')"
                color="indigo"
                class="mt-4"
                variant="outlined"
                size="small"
                prepend-icon="mdi-play-circle-outline"
                :loading="running"
                @click="runNow"
              >
                Run now
              </v-btn>
            </div>
            <empty-state
              v-else
              icon="mdi-clock-alert-outline"
              title="No automatic cut-off configured"
              message="Enable automatic generation and supply at least one HH:MM time."
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="4000">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const router = useRouter()
const { hasPermission } = usePermissionCheck()

const form = reactive({
  enabled: true,
  cutoff_times: '',
  timezone: 'Africa/Johannesburg',
  daily_payment_limit: 0
})
const saving = ref(false)
const loading = ref(false)
const next = ref('')
const running = ref(false)

const runs = ref<any[]>([])
const loadingRuns = ref(false)

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const runHeaders = [
  { title: 'Scheduled at', key: 'scheduled_at', sortable: true },
  { title: 'Trigger', key: 'trigger_type', sortable: true },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Claims', key: 'claims_count', align: 'end' as const },
  { title: 'Total', key: 'total_amount', align: 'end' as const },
  { title: 'Triggered by', key: 'triggered_by' }
]

const canSave = computed(() => {
  if (!form.enabled) return true
  return form.cutoff_times.trim().length > 0
})

function unwrap(res: any) {
  const body = res?.data
  if (body && typeof body === 'object' && 'success' in body && 'data' in body) {
    return body.data
  }
  return body
}

function notify(msg: string, color = 'success') {
  snackbarMessage.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

function formatCurrency(val: number) {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(val ?? 0)
}

function formatDateTime(val: string) {
  if (!val) return '—'
  return new Date(val).toLocaleString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function runColor(status: string) {
  const map: Record<string, string> = {
    ok: 'success',
    running: 'info',
    no_claims: 'grey',
    error: 'error'
  }
  return map[status] ?? 'default'
}

function goBack() {
  router.push({ name: 'group-pricing-claim-payment-schedules' })
}

async function loadConfig() {
  loading.value = true
  try {
    const res = await GroupPricingService.getPaymentCutoffConfig()
    const cfg = unwrap(res) ?? {}
    form.enabled = Boolean(cfg.enabled)
    form.cutoff_times = cfg.cutoff_times ?? ''
    form.timezone = cfg.timezone || 'Africa/Johannesburg'
    form.daily_payment_limit = cfg.daily_payment_limit ?? 0
  } catch (e: any) {
    notify('Failed to load cut-off settings', 'error')
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    await GroupPricingService.savePaymentCutoffConfig(form)
    notify('Cut-off settings saved.')
    await Promise.all([loadConfig(), loadNext()])
  } catch (e: any) {
    notify(
      e?.response?.data?.message ??
        e?.response?.data ??
        'Failed to save settings',
      'error'
    )
  } finally {
    saving.value = false
  }
}

async function loadNext() {
  try {
    const res = await GroupPricingService.getNextPaymentCutoff()
    const body = unwrap(res)
    if (body?.configured && body?.scheduled_at) {
      next.value = formatDateTime(body.scheduled_at)
    } else {
      next.value = ''
    }
  } catch {
    next.value = ''
  }
}

async function loadRuns() {
  loadingRuns.value = true
  try {
    const res = await GroupPricingService.listPaymentCutoffRuns(30)
    runs.value = unwrap(res) ?? []
  } catch {
    runs.value = []
  } finally {
    loadingRuns.value = false
  }
}

async function runNow() {
  running.value = true
  try {
    const res = await GroupPricingService.runPaymentCutoffNow()
    const run = unwrap(res)
    notify(
      run?.status === 'no_claims'
        ? 'No approved claims at this time — nothing to schedule.'
        : run?.status === 'ok'
          ? `Schedule created with ${run.claims_count} claims.`
          : (run?.error_message ?? 'Cut-off failed.'),
      run?.status === 'ok'
        ? 'success'
        : run?.status === 'error'
          ? 'error'
          : 'info'
    )
    await loadRuns()
  } catch (e: any) {
    notify(
      e?.response?.data?.message ??
        e?.response?.data ??
        'Failed to run cut-off',
      'error'
    )
  } finally {
    running.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadConfig(), loadNext(), loadRuns()])
})
</script>
