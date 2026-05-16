<template>
  <v-dialog
    :model-value="modelValue"
    max-width="540"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-card v-if="mode === 'open'">
      <v-card-title>Flag {{ targetName || 'user' }}</v-card-title>
      <v-card-subtitle class="pb-2">
        Raise a coaching or capacity concern. The note is internal — the flagged
        user only sees a neutral prompt to follow up.
      </v-card-subtitle>
      <v-card-text>
        <v-row dense>
          <v-col cols="12">
            <v-select
              v-model="form.flag_reason"
              :items="reasonOptions"
              label="Reason"
              density="compact"
              hide-details
            />
          </v-col>
          <v-col cols="12">
            <v-textarea
              v-model="form.note"
              label="Internal note"
              hint="What's prompting this flag? (Minimum 10 characters.)"
              persistent-hint
              rows="4"
              auto-grow
              :counter="500"
              :rules="[noteRule]"
            />
          </v-col>
        </v-row>
        <v-alert
          v-if="error"
          type="error"
          density="compact"
          variant="tonal"
          class="mt-2"
        >
          {{ error }}
        </v-alert>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="cancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="saving"
          :disabled="!canSubmit"
          @click="submitOpen"
        >
          Open flag
        </v-btn>
      </v-card-actions>
    </v-card>

    <v-card v-else-if="mode === 'resolve' && flagToResolve">
      <v-card-title>Resolve flag</v-card-title>
      <v-card-subtitle class="pb-2">
        Closing this flag preserves the history. The flagged user is not
        notified on resolve.
      </v-card-subtitle>
      <v-card-text>
        <div class="resolve-meta mb-3">
          <div>
            <span class="meta-label">User</span>
            <span class="meta-value">{{ flagToResolve.user_name }}</span>
          </div>
          <div>
            <span class="meta-label">Reason</span>
            <v-chip
              size="x-small"
              :color="reasonColor(flagToResolve.flag_reason)"
              variant="tonal"
              class="ms-1"
            >
              {{ reasonLabel(flagToResolve.flag_reason) }}
            </v-chip>
          </div>
          <div>
            <span class="meta-label">Opened</span>
            <span class="meta-value">
              {{ flagToResolve.opened_by_name }} ·
              {{ formatDate(flagToResolve.opened_at) }}
            </span>
          </div>
          <div v-if="flagToResolve.note">
            <span class="meta-label">Original note</span>
            <span class="meta-value note-value">{{ flagToResolve.note }}</span>
          </div>
        </div>
        <v-textarea
          v-model="resolutionNote"
          label="Resolution note"
          hint="How was this resolved? (Minimum 10 characters.)"
          persistent-hint
          rows="4"
          auto-grow
          :counter="500"
          :rules="[noteRule]"
        />
        <v-alert
          v-if="error"
          type="error"
          density="compact"
          variant="tonal"
          class="mt-2"
        >
          {{ error }}
        </v-alert>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn variant="text" @click="cancel">Cancel</v-btn>
        <v-btn
          color="primary"
          :loading="saving"
          :disabled="!resolveCanSubmit"
          @click="submitResolve"
        >
          Resolve
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'

import type {
  OpenUserFlagBody,
  QuoteUserFlag
} from '@/renderer/api/QuoteDashboardService'

const props = defineProps<{
  modelValue: boolean
  mode: 'open' | 'resolve'
  targetName?: string
  targetEmail?: string
  flagToResolve?: QuoteUserFlag | null
  saving?: boolean
  error?: string | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', open: boolean): void
  (e: 'submit-open', body: OpenUserFlagBody): void
  (e: 'submit-resolve', id: number, note: string): void
}>()

const reasonOptions = [
  { title: 'Coaching — skill or process concern', value: 'coaching' },
  { title: 'Capacity — workload or availability', value: 'capacity' }
]

const form = ref<{
  flag_reason: 'coaching' | 'capacity'
  note: string
}>({
  flag_reason: 'coaching',
  note: ''
})

const resolutionNote = ref('')

// Reset form when the dialog is opened with a new target.
watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      form.value = { flag_reason: 'coaching', note: '' }
      resolutionNote.value = ''
    }
  }
)

const noteRule = (v: string) =>
  (v && v.trim().length >= 10) ||
  'Please write at least 10 characters of context.'

const canSubmit = computed(() => form.value.note.trim().length >= 10)
const resolveCanSubmit = computed(
  () => resolutionNote.value.trim().length >= 10
)

function cancel() {
  emit('update:modelValue', false)
}

function submitOpen() {
  if (!props.targetName || !canSubmit.value) return
  emit('submit-open', {
    user_name: props.targetName,
    user_email: props.targetEmail,
    flag_reason: form.value.flag_reason,
    note: form.value.note.trim()
  })
}

function submitResolve() {
  if (!props.flagToResolve || !resolveCanSubmit.value) return
  emit('submit-resolve', props.flagToResolve.id, resolutionNote.value.trim())
}

function reasonLabel(r: string): string {
  if (r === 'coaching') return 'Coaching'
  if (r === 'capacity') return 'Capacity'
  return r
}

function reasonColor(r: string): string {
  if (r === 'coaching') return 'warning'
  if (r === 'capacity') return 'info'
  return 'primary'
}

function formatDate(iso?: string | null): string {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
}
</script>

<style scoped>
.resolve-meta {
  font-size: 0.8125rem;
  display: grid;
  gap: 6px;
}
.meta-label {
  color: var(--v-theme-on-surface-variant, #64748b);
  font-weight: 600;
  font-size: 0.6875rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  margin-right: 6px;
}
.note-value {
  display: block;
  margin-top: 2px;
  padding: 8px 10px;
  background: #f8fafc;
  border-left: 3px solid #cbd5e1;
  border-radius: 4px;
  white-space: pre-wrap;
}
</style>
