<template>
  <v-dialog :model-value="modelValue" max-width="900" @update:model-value="emit('update:modelValue', $event)">
    <v-card>
      <v-card-title>New Manual Journal Entry</v-card-title>
      <v-card-text>
        <v-row dense>
          <v-col cols="12" md="8">
            <v-text-field
              v-model="description"
              label="Description"
              density="compact"
              variant="outlined"
            />
          </v-col>
          <v-col cols="12" md="4">
            <v-select
              v-model="periodId"
              :items="periods"
              item-title="name"
              item-value="id"
              label="Period (optional — defaults to today)"
              density="compact"
              variant="outlined"
              clearable
            />
          </v-col>
        </v-row>

        <v-table density="compact">
          <thead>
            <tr>
              <th>Account</th>
              <th class="text-right">Debit</th>
              <th class="text-right">Credit</th>
              <th>Description</th>
              <th />
            </tr>
          </thead>
          <tbody>
            <tr v-for="(line, idx) in lines" :key="idx">
              <td style="min-width: 220px">
                <v-autocomplete
                  v-model="line.account_id"
                  :items="accounts"
                  :item-title="accountLabel"
                  item-value="id"
                  density="compact"
                  variant="outlined"
                  hide-details
                />
              </td>
              <td>
                <v-text-field
                  v-model.number="line.debit"
                  type="number"
                  density="compact"
                  variant="outlined"
                  hide-details
                  reverse
                />
              </td>
              <td>
                <v-text-field
                  v-model.number="line.credit"
                  type="number"
                  density="compact"
                  variant="outlined"
                  hide-details
                  reverse
                />
              </td>
              <td>
                <v-text-field
                  v-model="line.description"
                  density="compact"
                  variant="outlined"
                  hide-details
                />
              </td>
              <td>
                <v-btn
                  v-if="lines.length > 2"
                  icon="mdi-delete"
                  variant="text"
                  size="small"
                  @click="lines.splice(idx, 1)"
                />
              </td>
            </tr>
          </tbody>
        </v-table>

        <div class="d-flex justify-space-between align-center mt-3">
          <v-btn
            prepend-icon="mdi-plus"
            variant="text"
            size="small"
            @click="addLine"
            >Add line</v-btn
          >
          <div class="text-subtitle-2">
            <span class="mr-6">Debit: <strong>{{ format(totalDr) }}</strong></span>
            <span class="mr-6">Credit: <strong>{{ format(totalCr) }}</strong></span>
            <v-chip
              :color="balanced ? 'success' : 'warning'"
              size="small"
              variant="tonal"
              >{{ balanced ? 'Balanced' : 'Unbalanced' }}</v-chip
            >
          </div>
        </div>

        <v-alert
          v-if="error"
          color="error"
          variant="tonal"
          class="mt-3"
          icon="mdi-alert"
          >{{ error }}</v-alert
        >
      </v-card-text>
      <v-alert
        v-if="info"
        color="info"
        variant="tonal"
        class="ma-3"
        icon="mdi-information-outline"
        density="compact"
      >
        {{ info }}
      </v-alert>
      <v-card-actions>
        <v-spacer />
        <v-btn @click="close">Cancel</v-btn>
        <v-btn
          variant="text"
          :disabled="!balanced || !description.trim() || saving"
          :loading="saving && action === 'draft'"
          @click="saveAsDraft"
          >Save Draft</v-btn
        >
        <v-btn
          color="primary"
          :disabled="!balanced || !description.trim() || saving"
          :loading="saving && action === 'submit'"
          @click="submitForApproval"
          >Submit for Approval</v-btn
        >
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import GLService, {
  AccountingPeriod,
  GLAccount,
  ManualJournalLineInput
} from '@/renderer/api/GeneralLedgerService'

const props = defineProps<{
  modelValue: boolean
  accounts: GLAccount[]
  periods: AccountingPeriod[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: boolean): void
  (e: 'posted'): void
}>()

const description = ref('')
const periodId = ref<number | undefined>(undefined)
const lines = reactive<ManualJournalLineInput[]>([
  { account_id: 0, debit: 0, credit: 0, description: '' },
  { account_id: 0, debit: 0, credit: 0, description: '' }
])
const saving = ref(false)
const error = ref('')
const info = ref('')
const action = ref<'' | 'draft' | 'submit'>('')

const accountLabel = (a: GLAccount) => `${a.code} — ${a.name}`

const totalDr = computed(() =>
  lines.reduce((sum, l) => sum + (Number(l.debit) || 0), 0)
)
const totalCr = computed(() =>
  lines.reduce((sum, l) => sum + (Number(l.credit) || 0), 0)
)
const balanced = computed(
  () =>
    Math.abs(totalDr.value - totalCr.value) < 0.005 &&
    totalDr.value > 0 &&
    lines.every((l) => l.account_id > 0)
)

const format = (n: number) =>
  new Intl.NumberFormat(undefined, { minimumFractionDigits: 2 }).format(n || 0)

const addLine = () => {
  lines.push({ account_id: 0, debit: 0, credit: 0, description: '' })
}

const reset = () => {
  description.value = ''
  periodId.value = undefined
  lines.splice(0, lines.length)
  lines.push({ account_id: 0, debit: 0, credit: 0, description: '' })
  lines.push({ account_id: 0, debit: 0, credit: 0, description: '' })
  error.value = ''
  info.value = ''
  action.value = ''
}

const close = () => {
  emit('update:modelValue', false)
}

watch(
  () => props.modelValue,
  (open) => {
    if (open) reset()
  }
)

const buildPayload = () => ({
  description: description.value.trim(),
  period_id: periodId.value || undefined,
  lines: lines.map((l) => ({
    account_id: l.account_id,
    debit: Number(l.debit) || 0,
    credit: Number(l.credit) || 0,
    description: l.description
  }))
})

const saveAsDraft = async () => {
  error.value = ''
  info.value = ''
  saving.value = true
  action.value = 'draft'
  try {
    await GLService.draftManualJournal(buildPayload())
    info.value = 'Saved as draft. Submit when ready for approval.'
    emit('posted')
    close()
  } catch (e: any) {
    error.value = e?.response?.data?.error || e?.message || 'Failed to save draft'
  } finally {
    saving.value = false
    action.value = ''
  }
}

const submitForApproval = async () => {
  error.value = ''
  info.value = ''
  saving.value = true
  action.value = 'submit'
  try {
    const draft = await GLService.draftManualJournal(buildPayload())
    if (draft?.id) {
      await GLService.submitManualJournal(draft.id)
    }
    info.value = 'Submitted for approval. A different user must approve before posting.'
    emit('posted')
    close()
  } catch (e: any) {
    error.value =
      e?.response?.data?.error || e?.message || 'Failed to submit for approval'
  } finally {
    saving.value = false
    action.value = ''
  }
}
</script>
