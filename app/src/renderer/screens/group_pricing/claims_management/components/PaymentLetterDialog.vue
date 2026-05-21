<template>
  <v-dialog
    :model-value="modelValue"
    max-width="720px"
    @update:model-value="(v) => emit('update:modelValue', v)"
  >
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-icon class="mr-2">mdi-file-document-check-outline</v-icon>
        <span>Payment Confirmation Letter</span>
        <v-spacer />
        <v-btn icon="mdi-close" variant="text" @click="close" />
      </v-card-title>
      <v-divider />

      <v-tabs v-model="tab" grow>
        <v-tab value="generate">Generate</v-tab>
        <v-tab value="history">
          History
          <v-chip
            v-if="history.length"
            size="x-small"
            class="ml-2"
            color="primary"
            variant="flat"
            >{{ history.length }}</v-chip
          >
        </v-tab>
        <v-tab value="send">Send</v-tab>
      </v-tabs>

      <v-divider />

      <v-window v-model="tab">
        <!-- Generate -->
        <v-window-item value="generate">
          <v-card-text>
            <p class="text-body-2 mb-4">
              Choose a format. The letter is generated on demand from the
              current claim data and letter settings — a history record is
              written each time so you can re-download or send it later.
            </p>
            <v-radio-group v-model="format" inline class="mb-2">
              <v-radio label="PDF" value="pdf" />
              <v-radio label="DOCX" value="docx" />
            </v-radio-group>
            <v-alert
              v-if="generateError"
              type="error"
              variant="tonal"
              density="compact"
              class="mb-3"
              >{{ generateError }}</v-alert
            >
            <v-btn
              color="primary"
              :loading="generating"
              prepend-icon="mdi-download"
              @click="generate"
              >Generate &amp; Download</v-btn
            >
          </v-card-text>
        </v-window-item>

        <!-- History -->
        <v-window-item value="history">
          <v-card-text>
            <div v-if="loadingHistory" class="text-center py-6">
              <v-progress-circular indeterminate />
            </div>
            <v-alert
              v-else-if="!history.length"
              type="info"
              variant="tonal"
              density="compact"
              >No letters generated yet for this claim.</v-alert
            >
            <v-table v-else density="compact">
              <thead>
                <tr>
                  <th>Ver</th>
                  <th>Format</th>
                  <th>Generated</th>
                  <th>By</th>
                  <th>Reference</th>
                  <th class="text-right">Actions</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in history" :key="row.id">
                  <td>v{{ row.version }}</td>
                  <td class="text-uppercase">{{ row.format }}</td>
                  <td>{{ formatDate(row.generated_at) }}</td>
                  <td>{{ row.generated_by || '—' }}</td>
                  <td class="text-caption">{{ row.letter_reference }}</td>
                  <td class="text-right">
                    <v-tooltip text="Re-download">
                      <template #activator="{ props: tipProps }">
                        <v-btn
                          v-bind="tipProps"
                          icon="mdi-download"
                          size="small"
                          variant="text"
                          :loading="redownloadingId === row.id"
                          @click="redownload(row)"
                        />
                      </template>
                    </v-tooltip>
                    <v-tooltip text="Send by email">
                      <template #activator="{ props: tipProps }">
                        <v-btn
                          v-bind="tipProps"
                          icon="mdi-email-outline"
                          size="small"
                          variant="text"
                          @click="prefillSend(row)"
                        />
                      </template>
                    </v-tooltip>
                  </td>
                </tr>
              </tbody>
            </v-table>
          </v-card-text>
        </v-window-item>

        <!-- Send -->
        <v-window-item value="send">
          <v-card-text>
            <v-alert
              v-if="!history.length"
              type="info"
              variant="tonal"
              density="compact"
              class="mb-3"
              >Generate a letter first before sending.</v-alert
            >
            <template v-else>
              <v-select
                v-model="sendLetterId"
                :items="historyOptions"
                item-title="label"
                item-value="value"
                label="Letter version to send"
                variant="outlined"
                density="compact"
                class="mb-3"
              />
              <v-select
                v-model="channel"
                :items="channelOptions"
                item-title="label"
                item-value="value"
                :item-props="channelItemProps"
                label="Channel"
                variant="outlined"
                density="compact"
                class="mb-3"
              />
              <v-text-field
                v-model="recipient"
                :label="recipientLabel"
                :hint="recipientHint"
                persistent-hint
                variant="outlined"
                density="compact"
                class="mb-3"
              />
              <v-alert
                v-if="recipientMissing"
                type="warning"
                variant="tonal"
                density="compact"
                class="mb-3"
              >
                The claim has no {{ channel }} on file. Update the claim record
                or fill in the field above to override.
              </v-alert>
              <v-alert
                v-if="sendError"
                type="error"
                variant="tonal"
                density="compact"
                class="mb-3"
                >{{ sendError }}</v-alert
              >
              <v-btn
                color="primary"
                :loading="sending"
                :disabled="
                  !sendLetterId || !channel || channelDisabled(channel)
                "
                prepend-icon="mdi-send"
                @click="send"
                >Send to claimant</v-btn
              >
            </template>
          </v-card-text>
        </v-window-item>
      </v-window>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface LetterRow {
  id: number
  claim_id: number
  version: number
  format: string
  filename: string
  letter_reference: string
  generated_at: string
  generated_by: string
}

interface Props {
  modelValue: boolean
  claimId: number | null
  claimantEmail?: string
  claimantPhone?: string
}

const props = withDefaults(defineProps<Props>(), {
  claimantEmail: '',
  claimantPhone: ''
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const tab = ref<'generate' | 'history' | 'send'>('generate')
const format = ref<'pdf' | 'docx'>('pdf')
const generating = ref(false)
const generateError = ref('')

const history = ref<LetterRow[]>([])
const loadingHistory = ref(false)
const redownloadingId = ref<number | null>(null)

const sendLetterId = ref<number | null>(null)
const channel = ref<'email' | 'sms' | 'whatsapp' | 'manual'>('email')
const recipient = ref('')
const sending = ref(false)
const sendError = ref('')

const channelOptions = [
  { label: 'Email', value: 'email', disabled: false },
  {
    label: 'SMS (coming soon)',
    value: 'sms',
    disabled: true
  },
  {
    label: 'WhatsApp (coming soon)',
    value: 'whatsapp',
    disabled: true
  },
  { label: 'Manual / In person', value: 'manual', disabled: false }
]

const channelItemProps = (item: any) => ({ disabled: item.raw?.disabled })

function channelDisabled(c: string) {
  return channelOptions.find((o) => o.value === c)?.disabled === true
}

const recipientLabel = computed(() => {
  switch (channel.value) {
    case 'email':
      return 'Email address'
    case 'sms':
    case 'whatsapp':
      return 'Mobile number'
    default:
      return 'Reference (e.g. "Hand-delivered to ...")'
  }
})

const recipientHint = computed(() => {
  switch (channel.value) {
    case 'email':
      return 'Defaults to the claim record; leave blank to use the claimant_email on file.'
    default:
      return ''
  }
})

const recipientMissing = computed(() => {
  if (recipient.value) return false
  if (channel.value === 'email' && !props.claimantEmail) return true
  if (
    (channel.value === 'sms' || channel.value === 'whatsapp') &&
    !props.claimantPhone
  )
    return true
  return false
})

const historyOptions = computed(() =>
  history.value.map((r) => ({
    value: r.id,
    label: `v${r.version} — ${r.format.toUpperCase()} — ${formatDate(r.generated_at)}`
  }))
)

watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      generateError.value = ''
      sendError.value = ''
      sendLetterId.value = null
      tab.value = 'generate'
      loadHistory()
      recipient.value = ''
    }
  }
)

watch(channel, () => {
  // Reset recipient when channel changes so the new default applies.
  recipient.value = ''
})

async function loadHistory() {
  if (!props.claimId) return
  loadingHistory.value = true
  try {
    const res = await GroupPricingService.getClaimPaymentLetterHistory(
      props.claimId
    )
    history.value = unwrap(res.data) || []
    if (history.value.length && !sendLetterId.value) {
      sendLetterId.value = history.value[0].id
    }
  } catch (err: any) {
    history.value = []
    console.error('Failed to load letter history', err)
  } finally {
    loadingHistory.value = false
  }
}

async function generate() {
  if (!props.claimId) return
  generating.value = true
  generateError.value = ''
  try {
    const call =
      format.value === 'pdf'
        ? GroupPricingService.getClaimPaymentLetterPdf
        : GroupPricingService.getClaimPaymentLetterDocx
    const res = await call(props.claimId)
    const ext = format.value
    downloadBlob(
      res.data,
      `PaymentConfirmation_claim_${props.claimId}.${ext}`,
      mimeFor(ext)
    )
    await loadHistory()
  } catch (err: any) {
    generateError.value = extractError(err)
  } finally {
    generating.value = false
  }
}

async function redownload(row: LetterRow) {
  if (!props.claimId) return
  redownloadingId.value = row.id
  try {
    const call =
      row.format === 'pdf'
        ? GroupPricingService.getClaimPaymentLetterPdf
        : GroupPricingService.getClaimPaymentLetterDocx
    const res = await call(props.claimId)
    downloadBlob(
      res.data,
      row.filename || `letter.${row.format}`,
      mimeFor(row.format)
    )
    await loadHistory()
  } catch (err: any) {
    generateError.value = extractError(err)
    tab.value = 'generate'
  } finally {
    redownloadingId.value = null
  }
}

function prefillSend(row: LetterRow) {
  sendLetterId.value = row.id
  channel.value = 'email'
  tab.value = 'send'
}

async function send() {
  if (!props.claimId || !sendLetterId.value) return
  sending.value = true
  sendError.value = ''
  try {
    await GroupPricingService.sendClaimPaymentLetter(props.claimId, {
      letter_id: sendLetterId.value,
      channel: channel.value,
      recipient: recipient.value || undefined
    })
    await loadHistory()
    tab.value = 'history'
  } catch (err: any) {
    sendError.value = extractError(err)
  } finally {
    sending.value = false
  }
}

function close() {
  emit('update:modelValue', false)
}

function downloadBlob(data: any, filename: string, type: string) {
  const blob = new Blob([data], { type })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

function mimeFor(ext: string): string {
  if (ext === 'pdf') return 'application/pdf'
  return 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
}

function formatDate(iso: string): string {
  if (!iso) return '—'
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
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
