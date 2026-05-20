<template>
  <v-navigation-drawer
    :model-value="modelValue"
    location="right"
    temporary
    width="420"
    @update:model-value="(v) => emit('update:modelValue', v)"
  >
    <v-toolbar density="compact" color="primary">
      <v-icon class="mr-2">mdi-shield-search-outline</v-icon>
      <v-toolbar-title class="text-subtitle-1">
        Pre-authorisation review
      </v-toolbar-title>
      <v-spacer />
      <v-btn icon="mdi-close" variant="text" @click="close" />
    </v-toolbar>

    <div v-if="item" class="pa-4">
      <div class="text-caption text-medium-emphasis mb-1">Claim</div>
      <div class="text-subtitle-2">{{ item.claim_number }}</div>
      <div class="text-body-2 mb-4">{{ item.beneficiary_name || item.member_name }}</div>

      <!-- Banking -->
      <v-card variant="outlined" rounded="lg" class="mb-3">
        <v-card-title class="text-subtitle-2 pa-3 pb-1 d-flex align-center">
          <v-icon size="small" class="mr-2">mdi-bank-check</v-icon>
          Banking verification
        </v-card-title>
        <v-card-text class="pa-3 pt-1">
          <div class="text-body-2 mb-1">
            <strong>{{ item.bank_name || '—' }}</strong>
            <span class="text-medium-emphasis">·
              {{ maskedAccount }}</span>
          </div>
          <div v-if="loadingBank" class="text-center py-2">
            <v-progress-circular indeterminate size="20" />
          </div>
          <template v-else-if="bankStatus">
            <div v-if="!bankStatus.has_result" class="text-body-2 text-medium-emphasis">
              No verification recorded yet.
            </div>
            <template v-else>
              <v-chip
                :color="bankChipColor"
                size="small"
                variant="flat"
                class="mb-1"
              >
                {{ bankChipLabel }}
              </v-chip>
              <div v-if="bankStatus.verified_at" class="text-caption text-medium-emphasis">
                {{ relativeAge(bankStatus.verified_at) }}
              </div>
              <div v-if="bankStatus.stale_reason" class="text-caption text-warning mt-1">
                {{ bankStatus.stale_reason }}
              </div>
            </template>
          </template>

          <v-alert
            v-if="bankError"
            type="error"
            density="compact"
            variant="tonal"
            class="mt-2"
          >{{ bankError }}</v-alert>

          <v-btn
            color="primary"
            size="small"
            variant="flat"
            block
            class="mt-3"
            prepend-icon="mdi-refresh"
            :loading="reverifying"
            :disabled="!!reverifyDisabledReason"
            :title="reverifyDisabledReason"
            @click="reverify"
          >Re-verify now</v-btn>
        </v-card-text>
      </v-card>

      <!-- Amount drift -->
      <v-card variant="outlined" rounded="lg" class="mb-3">
        <v-card-title class="text-subtitle-2 pa-3 pb-1 d-flex align-center">
          <v-icon size="small" class="mr-2">mdi-scale-balance</v-icon>
          Amount review
        </v-card-title>
        <v-card-text class="pa-3 pt-1">
          <div class="d-flex justify-space-between text-body-2">
            <span class="text-medium-emphasis">Approved</span>
            <strong>{{ formatCurrency(item.approved_amount_snapshot ?? 0) }}</strong>
          </div>
          <div class="d-flex justify-space-between text-body-2">
            <span class="text-medium-emphasis">Scheduled</span>
            <strong>{{ formatCurrency(item.gross_amount ?? item.claim_amount ?? 0) }}</strong>
          </div>
          <v-divider class="my-2" />
          <div class="d-flex justify-space-between text-body-2">
            <span class="text-medium-emphasis">Drift</span>
            <strong :class="driftDeltaClass">{{ driftDeltaLabel }}</strong>
          </div>
          <div v-if="!hasDrift" class="text-caption text-success mt-2">
            <v-icon size="x-small">mdi-check-circle-outline</v-icon>
            Amount matches the approved figure.
          </div>
          <div v-else-if="item.amount_drift_resolved" class="text-caption text-medium-emphasis mt-2">
            Drift acknowledged by {{ item.amount_drift_resolved_by || 'finance' }}.
          </div>
          <div v-else class="text-caption text-warning mt-2">
            Drift outstanding — acknowledge or query the line before authorising.
          </div>
          <div v-if="hasDrift && !item.amount_drift_resolved" class="mt-3">
            <v-btn
              color="primary"
              size="small"
              variant="flat"
              :loading="acknowledging"
              @click="acknowledge"
            >Acknowledge drift</v-btn>
          </div>
        </v-card-text>
      </v-card>

      <!-- Cross-claim signals -->
      <v-card variant="outlined" rounded="lg" class="mb-3">
        <v-card-title class="text-subtitle-2 pa-3 pb-1 d-flex align-center">
          <v-icon size="small" class="mr-2">mdi-account-multiple-outline</v-icon>
          Cross-claim duplicates
        </v-card-title>
        <v-card-text class="pa-3 pt-1">
          <div v-if="!idHits.length && !accountHits.length" class="text-caption text-success">
            <v-icon size="x-small">mdi-check-circle-outline</v-icon>
            No prior claims found for this ID or bank account.
          </div>
          <template v-else>
            <div v-if="idHits.length" class="mb-3">
              <div class="text-caption text-warning mb-1">Same claimant ID:</div>
              <div v-for="r in idHits" :key="`id-${r}`" class="text-caption">
                • {{ r }}
              </div>
            </div>
            <div v-if="accountHits.length">
              <div class="text-caption text-warning mb-1">Same bank account:</div>
              <div v-for="r in accountHits" :key="`acc-${r}`" class="text-caption">
                • {{ r }}
              </div>
            </div>
          </template>
        </v-card-text>
      </v-card>

      <!-- Other flags summary -->
      <v-card variant="outlined" rounded="lg">
        <v-card-title class="text-subtitle-2 pa-3 pb-1 d-flex align-center">
          <v-icon size="small" class="mr-2">mdi-flag-outline</v-icon>
          Other risk flags
        </v-card-title>
        <v-card-text class="pa-3 pt-1">
          <div v-if="!otherFlags.length" class="text-caption text-medium-emphasis">
            No other risk flags on this line.
          </div>
          <div v-else class="d-flex flex-wrap gap-1">
            <v-chip
              v-for="flag in otherFlags"
              :key="flag"
              size="x-small"
              color="warning"
              variant="flat"
            >{{ flag }}</v-chip>
          </div>
        </v-card-text>
      </v-card>
    </div>
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import type {
  BankVerificationStatus,
  RiskFlags,
  ScheduleItem
} from '../payment_schedule_context'

interface Props {
  modelValue: boolean
  scheduleId: number
  item: ScheduleItem | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'acknowledged', itemId: number): void
  (e: 'reverified'): void
}>()

const bankStatus = ref<BankVerificationStatus | null>(null)
const loadingBank = ref(false)
const bankError = ref('')
const reverifying = ref(false)
const acknowledging = ref(false)

const flags = computed<RiskFlags>(() => {
  if (!props.item) return {}
  return parseRiskFlags(props.item)
})

const idHits = computed<string[]>(() => flags.value.id_paid_before_refs ?? [])
const accountHits = computed<string[]>(
  () => flags.value.account_used_before_refs ?? []
)

const otherFlags = computed<string[]>(() => {
  const f = flags.value
  const out: string[] = []
  if (f.banking_change_30d) out.push('Banking changed (30d)')
  if (f.contestable) out.push('Contestable')
  if (f.recent_reinstatement) out.push('Reinstated')
  if (f.fraud_risk_level) out.push(`Fraud: ${f.fraud_risk_level}`)
  if (
    props.item?.duplicate_beneficiary_flag &&
    !props.item?.duplicate_beneficiary_cleared
  ) {
    out.push('Duplicate beneficiary (in-schedule)')
  }
  return out
})

const maskedAccount = computed(() => {
  const acc = (props.item?.bank_account_number || '').trim()
  if (!acc) return '—'
  if (acc.length <= 4) return acc
  return '••••' + acc.slice(-4)
})

const bankChipLabel = computed(() => {
  if (!bankStatus.value || !bankStatus.value.has_result) return 'Not verified'
  if (bankStatus.value.status === 'failed') return 'Failed'
  if (bankStatus.value.status === 'pending') return 'Pending'
  if (bankStatus.value.stale) return 'Verified · stale'
  return 'Verified'
})

const bankChipColor = computed(() => {
  if (!bankStatus.value || !bankStatus.value.has_result) return 'grey'
  if (bankStatus.value.status === 'failed') return 'error'
  if (bankStatus.value.status === 'pending') return 'info'
  if (bankStatus.value.stale) return 'warning'
  return 'success'
})

const reverifyDisabledReason = computed(() => {
  if (!bankStatus.value?.verified_at) return ''
  const ageMs = Date.now() - new Date(bankStatus.value.verified_at).getTime()
  if (ageMs < 60 * 60 * 1000) {
    return 'A verification was made in the last hour — wait before re-trying to avoid duplicate provider charges.'
  }
  return ''
})

const hasDrift = computed(() => {
  if (!props.item) return false
  const approved = props.item.approved_amount_snapshot ?? 0
  if (approved === 0) return false
  const gross = props.item.gross_amount ?? props.item.claim_amount ?? 0
  return Math.abs(gross - approved) > 1
})

const driftDelta = computed(() => {
  if (!props.item) return 0
  const approved = props.item.approved_amount_snapshot ?? 0
  const gross = props.item.gross_amount ?? props.item.claim_amount ?? 0
  return gross - approved
})

const driftDeltaLabel = computed(() => {
  const d = driftDelta.value
  const sign = d > 0 ? '+' : d < 0 ? '−' : ''
  return `${sign}${formatCurrency(Math.abs(d))}`
})

const driftDeltaClass = computed(() => {
  if (driftDelta.value > 0) return 'text-error'
  if (driftDelta.value < 0) return 'text-warning'
  return ''
})

watch(
  () => [props.modelValue, props.item?.id],
  ([open]) => {
    if (open && props.item) {
      loadBank()
    }
  }
)

async function loadBank() {
  if (!props.item) return
  loadingBank.value = true
  bankError.value = ''
  try {
    const res = await GroupPricingService.getLineBankVerification(
      props.scheduleId,
      props.item.id
    )
    bankStatus.value = unwrap<BankVerificationStatus>(res.data)
  } catch (err: any) {
    bankError.value = extractError(err)
    bankStatus.value = null
  } finally {
    loadingBank.value = false
  }
}

async function reverify() {
  if (!props.item) return
  reverifying.value = true
  bankError.value = ''
  try {
    await GroupPricingService.reverifyLineBankAccount(
      props.scheduleId,
      props.item.id
    )
    await loadBank()
    emit('reverified')
  } catch (err: any) {
    bankError.value = extractError(err)
  } finally {
    reverifying.value = false
  }
}

async function acknowledge() {
  if (!props.item) return
  acknowledging.value = true
  try {
    await GroupPricingService.acknowledgeAmountDrift(
      props.scheduleId,
      props.item.id
    )
    // Parent owns the schedule items array; emitting causes it to update
    // the item's resolved flag, which re-renders this drawer via props.
    emit('acknowledged', props.item.id)
  } catch (err: any) {
    bankError.value = extractError(err)
  } finally {
    acknowledging.value = false
  }
}

function close() {
  emit('update:modelValue', false)
}

function parseRiskFlags(item: ScheduleItem): RiskFlags {
  const raw = item.risk_flags
  if (!raw) return {}
  if (typeof raw === 'string') {
    try {
      return JSON.parse(raw) as RiskFlags
    } catch {
      return {}
    }
  }
  return raw as RiskFlags
}

function formatCurrency(n: number): string {
  const num = Number.isFinite(n) ? n : 0
  return num.toLocaleString('en-ZA', {
    style: 'currency',
    currency: 'ZAR',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  })
}

function relativeAge(iso: string): string {
  if (!iso) return ''
  const ms = Date.now() - new Date(iso).getTime()
  if (ms < 60_000) return 'just now'
  const minutes = Math.round(ms / 60_000)
  if (minutes < 60) return `${minutes} min ago`
  const hours = Math.round(minutes / 60)
  if (hours < 24) return `${hours} hour${hours === 1 ? '' : 's'} ago`
  const days = Math.round(hours / 24)
  return `${days} day${days === 1 ? '' : 's'} ago`
}

function unwrap<T>(payload: any): T {
  if (payload && typeof payload === 'object' && 'data' in payload) {
    return payload.data as T
  }
  return payload as T
}

function extractError(err: any): string {
  return (
    err?.response?.data?.error ||
    err?.response?.data?.message ||
    err?.message ||
    'Unexpected error'
  )
}
</script>
