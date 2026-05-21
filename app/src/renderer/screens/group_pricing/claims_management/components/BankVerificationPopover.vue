<template>
  <v-menu
    v-model="open"
    :close-on-content-click="false"
    location="bottom start"
    offset="8"
    min-width="340"
  >
    <template #activator="{ props: menuProps }">
      <button
        v-bind="menuProps"
        type="button"
        class="bav-trigger"
        :title="`Bank verification — ${chipLabel}`"
      >
        <v-chip
          :color="chipColor"
          size="x-small"
          variant="flat"
          :prepend-icon="chipIcon"
        >
          {{ chipLabel }}
        </v-chip>
        <div class="bav-trigger__meta">
          <div class="bav-trigger__bank">{{
            bankName || 'No bank on file'
          }}</div>
          <div v-if="maskedAccount" class="bav-trigger__account">{{
            maskedAccount
          }}</div>
        </div>
      </button>
    </template>

    <v-card max-width="380" rounded="lg">
      <v-card-title class="text-subtitle-2 pa-3 pb-2 d-flex align-center">
        <v-icon class="mr-2" size="small">mdi-bank-check</v-icon>
        <span>Bank account verification</span>
        <v-spacer />
        <v-btn
          icon="mdi-close"
          variant="text"
          size="x-small"
          density="compact"
          @click="open = false"
        />
      </v-card-title>
      <v-divider />
      <v-card-text class="pa-3">
        <div v-if="loading" class="text-center py-3">
          <v-progress-circular indeterminate size="24" />
        </div>
        <template v-else>
          <div
            v-if="!status?.has_result"
            class="text-body-2 text-medium-emphasis"
          >
            No verification has been recorded for this claim yet.
          </div>
          <template v-else>
            <div class="d-flex align-center mb-2">
              <v-chip
                :color="chipColor"
                size="small"
                variant="flat"
                class="mr-2"
                >{{ chipLabel }}</v-chip
              >
              <span
                v-if="status.verified_at"
                class="text-caption text-medium-emphasis"
              >
                {{ relativeAge(status.verified_at) }}
              </span>
            </div>
            <div
              v-if="status.stale && status.stale_reason"
              class="text-body-2 text-warning mb-2"
              >{{ status.stale_reason }}</div
            >
            <div
              v-if="status.provider_request_id"
              class="text-caption text-medium-emphasis"
            >
              Provider ref: {{ status.provider_request_id }}
            </div>
            <div class="text-caption text-medium-emphasis">
              Attempts on this claim: {{ status.last_attempt }}
            </div>
          </template>

          <v-divider class="my-3" />
          <v-alert
            v-if="reverifyError"
            type="error"
            density="compact"
            variant="tonal"
            class="mb-2"
            closable
            @click:close="reverifyError = ''"
            >{{ reverifyError }}</v-alert
          >
          <v-tooltip
            :disabled="!reverifyDisabledReason"
            location="top"
            max-width="280"
          >
            <template #activator="{ props: tipProps }">
              <span v-bind="tipProps">
                <v-btn
                  color="primary"
                  size="small"
                  block
                  prepend-icon="mdi-refresh"
                  :loading="reverifying"
                  :disabled="!!reverifyDisabledReason"
                  @click="reverify"
                  >Re-verify now</v-btn
                >
              </span>
            </template>
            {{ reverifyDisabledReason }}
          </v-tooltip>
        </template>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import type { BankVerificationStatus } from '../payment_schedule_context'

interface Props {
  scheduleId: number
  itemId: number
  claimId: number
  bankName?: string
  accountNumber?: string
  bankingChanged?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  bankName: '',
  accountNumber: '',
  bankingChanged: false
})

const emit = defineEmits<{ (e: 'reverified'): void }>()

const open = ref(false)
const status = ref<BankVerificationStatus | null>(null)
const loading = ref(false)
const reverifying = ref(false)
const reverifyError = ref('')

const maskedAccount = computed(() => {
  const acc = (props.accountNumber || '').trim()
  if (!acc) return ''
  if (acc.length <= 4) return acc
  return '••••' + acc.slice(-4)
})

const chipLabel = computed(() => {
  if (!status.value || !status.value.has_result) return 'Not verified'
  if (status.value.status === 'failed') return 'Failed'
  if (status.value.status === 'pending') return 'Pending'
  if (status.value.stale) return 'Verified · stale'
  return 'Verified'
})

const chipColor = computed(() => {
  if (!status.value || !status.value.has_result) return 'grey'
  if (status.value.status === 'failed') return 'error'
  if (status.value.status === 'pending') return 'info'
  if (status.value.stale) return 'warning'
  return 'success'
})

const chipIcon = computed(() => {
  if (!status.value || !status.value.has_result)
    return 'mdi-help-circle-outline'
  if (status.value.status === 'failed') return 'mdi-close-circle-outline'
  if (status.value.status === 'pending') return 'mdi-clock-outline'
  if (status.value.stale) return 'mdi-alert-outline'
  return 'mdi-check-circle-outline'
})

const reverifyDisabledReason = computed(() => {
  if (!status.value?.verified_at) return ''
  const ageMs = Date.now() - new Date(status.value.verified_at).getTime()
  if (ageMs < 60 * 60 * 1000) {
    return 'A verification was made in the last hour — wait before trying again to avoid duplicate provider charges.'
  }
  return ''
})

watch(open, async (val) => {
  if (val && !status.value) {
    await load()
  }
})

async function load() {
  loading.value = true
  reverifyError.value = ''
  try {
    const res = await GroupPricingService.getLineBankVerification(
      props.scheduleId,
      props.itemId
    )
    status.value = unwrap<BankVerificationStatus>(res.data)
  } catch (err: any) {
    reverifyError.value = extractError(err)
  } finally {
    loading.value = false
  }
}

async function reverify() {
  reverifying.value = true
  reverifyError.value = ''
  try {
    await GroupPricingService.reverifyLineBankAccount(
      props.scheduleId,
      props.itemId
    )
    await load()
    emit('reverified')
  } catch (err: any) {
    reverifyError.value = extractError(err)
  } finally {
    reverifying.value = false
  }
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

<style scoped>
.bav-trigger {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
  background: transparent;
  border: none;
  padding: 0;
  cursor: pointer;
  text-align: left;
}
.bav-trigger__meta {
  font-size: 11px;
  line-height: 1.2;
  color: var(--v-medium-emphasis-opacity, rgba(0, 0, 0, 0.6));
}
.bav-trigger__bank {
  font-weight: 500;
}
.bav-trigger__account {
  font-family: ui-monospace, SFMono-Regular, monospace;
  font-size: 10px;
}
</style>
