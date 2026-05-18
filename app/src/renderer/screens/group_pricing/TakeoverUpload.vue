<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center w-100">
              <v-btn
                size="small"
                variant="text"
                prepend-icon="mdi-arrow-left"
                @click="router.go(-1)"
                >Back</v-btn
              >
              <span class="headline ml-4"
                >Prior insurer schedule — Quote #{{ quoteId }}</span
              >
              <v-spacer />
              <v-chip
                v-if="schedule"
                size="small"
                color="success"
                variant="tonal"
                >Loaded: {{ schedule.member_count }} members ({{
                  schedule.in_force_count
                }}
                in force)</v-chip
              >
            </div>
          </template>

          <template #default>
            <v-card variant="outlined" rounded="lg" class="mb-4">
              <v-card-title class="font-weight-bold"
                >Upload schedule CSV</v-card-title
              >
              <v-card-text>
                <p class="text-caption mb-2"
                  >Canonical CSV header:
                  <code
                    >member_id_number, member_name, date_of_birth,
                    gla_sum_assured, ptd_sum_assured, ci_sum_assured,
                    prior_loadings, prior_exclusions, in_force</code
                  >. `prior_loadings` is pipe-separated `benefit:percent` (e.g.
                  `gla:25|ptd:10`). `prior_exclusions` is pipe-separated codes
                  (e.g. `smoker|diabetes`).</p
                >
                <v-row dense>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.insurer_name"
                      label="Outgoing insurer name"
                      density="compact"
                      variant="outlined"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.certificate_number"
                      label="Certificate number"
                      density="compact"
                      variant="outlined"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.effective_date"
                      label="Effective date (YYYY-MM-DD)"
                      density="compact"
                      variant="outlined"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.expiry_date"
                      label="Expiry date (YYYY-MM-DD)"
                      density="compact"
                      variant="outlined"
                    />
                  </v-col>
                  <v-col cols="12">
                    <v-textarea
                      v-model="form.notes"
                      label="Notes"
                      rows="2"
                      density="compact"
                      variant="outlined"
                    />
                  </v-col>
                  <v-col cols="12">
                    <v-file-input
                      v-model="form.file"
                      label="Schedule CSV"
                      accept=".csv"
                      density="compact"
                      variant="outlined"
                      show-size
                    />
                  </v-col>
                </v-row>
                <v-alert
                  v-if="error"
                  type="error"
                  variant="tonal"
                  density="compact"
                  class="mb-2"
                  >{{ error }}</v-alert
                >
                <v-btn
                  color="primary"
                  :loading="busy"
                  :disabled="!form.file"
                  @click="upload"
                  >Upload &amp; match</v-btn
                >
              </v-card-text>
            </v-card>

            <v-card v-if="summary" variant="outlined" rounded="lg" class="mb-4">
              <v-card-title class="d-flex align-center font-weight-bold">
                <span>Match preview</span>
                <v-spacer />
                <v-btn
                  size="small"
                  variant="text"
                  :loading="busy"
                  prepend-icon="mdi-refresh"
                  @click="rematch"
                  >Re-match</v-btn
                >
                <v-btn
                  v-if="
                    (summary.continuation_no_evidence ||
                      summary.continuation_with_loading) > 0
                  "
                  size="small"
                  color="primary"
                  variant="tonal"
                  class="ml-2"
                  :loading="busy"
                  prepend-icon="mdi-check-all"
                  @click="applyTerms"
                  >Apply takeover terms to cases</v-btn
                >
              </v-card-title>
              <v-card-text>
                <v-row dense>
                  <v-col cols="6" md="3">
                    <stat-card
                      title="Continuation, no evidence"
                      :value="summary.continuation_no_evidence"
                      color="success"
                      icon="mdi-shield-check-outline"
                    />
                  </v-col>
                  <v-col cols="6" md="3">
                    <stat-card
                      title="Continuation with loading"
                      :value="summary.continuation_with_loading"
                      color="warning"
                      icon="mdi-shield-alert-outline"
                    />
                  </v-col>
                  <v-col cols="6" md="3">
                    <stat-card
                      title="New evidence required"
                      :value="summary.new_evidence_required"
                      color="info"
                      icon="mdi-clipboard-search-outline"
                    />
                  </v-col>
                  <v-col cols="6" md="3">
                    <stat-card
                      title="Unmatched"
                      :value="summary.unmatched"
                      color="grey"
                      icon="mdi-help-circle-outline"
                    />
                  </v-col>
                </v-row>
                <p class="text-caption text-grey mt-2"
                  >Matched by ID: <strong>{{ summary.matched_by_id }}</strong>
                  · Matched by name+DOB:
                  <strong>{{ summary.matched_by_name_and_dob }}</strong> ·
                  Total: <strong>{{ summary.total }}</strong></p
                >
              </v-card-text>
            </v-card>

            <v-card v-if="schedule" variant="outlined" rounded="lg">
              <v-card-title class="font-weight-bold"
                >Prior members</v-card-title
              >
              <v-card-text>
                <v-table density="compact">
                  <thead>
                    <tr>
                      <th>ID number</th>
                      <th>Name</th>
                      <th>DOB</th>
                      <th class="text-right">GLA SA</th>
                      <th class="text-right">PTD SA</th>
                      <th class="text-right">CI SA</th>
                      <th>Prior loadings</th>
                      <th>Prior exclusions</th>
                      <th>In force</th>
                      <th>Match</th>
                      <th>Outcome</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="m in schedule.members || []" :key="m.id">
                      <td>{{ m.member_id_number || '—' }}</td>
                      <td>{{ m.member_name }}</td>
                      <td>{{ formatDate(m.date_of_birth) }}</td>
                      <td class="text-right">{{
                        formatMoney(m.gla_sum_assured)
                      }}</td>
                      <td class="text-right">{{
                        formatMoney(m.ptd_sum_assured)
                      }}</td>
                      <td class="text-right">{{
                        formatMoney(m.ci_sum_assured)
                      }}</td>
                      <td class="text-caption">{{ m.prior_loadings }}</td>
                      <td class="text-caption">{{ m.prior_exclusions }}</td>
                      <td>
                        <v-chip
                          :color="m.in_force ? 'success' : 'grey'"
                          size="x-small"
                          variant="tonal"
                          >{{ m.in_force ? 'Yes' : 'No' }}</v-chip
                        >
                      </td>
                      <td class="text-caption">
                        {{
                          m.matched_member_name
                            ? `${m.matched_member_name} (${m.matched_category})`
                            : '—'
                        }}
                      </td>
                      <td>
                        <v-chip
                          :color="outcomeColor(m.takeover_outcome)"
                          size="x-small"
                          variant="tonal"
                          >{{ outcomeLabel(m.takeover_outcome) }}</v-chip
                        >
                      </td>
                    </tr>
                    <tr v-if="!(schedule.members || []).length">
                      <td colspan="11" class="text-grey text-center"
                        >No members on this schedule yet.</td
                      >
                    </tr>
                  </tbody>
                </v-table>
              </v-card-text>
            </v-card>

            <empty-state
              v-if="!schedule && !busy"
              icon="mdi-clipboard-arrow-left-outline"
              title="No prior schedule"
              message="Upload a CSV to start the takeover matching workflow."
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import StatCard from '@/renderer/components/StatCard.vue'
import EmptyState from '@/renderer/components/EmptyState.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const route = useRoute()
const router = useRouter()
const quoteId = computed(() => String(route.params.quoteId))

interface PriorMember {
  id: number
  member_id_number: string
  member_name: string
  date_of_birth: string
  gla_sum_assured: number
  ptd_sum_assured: number
  ci_sum_assured: number
  prior_loadings: string
  prior_exclusions: string
  in_force: boolean
  matched_member_name: string
  matched_category: string
  matched_case_id: number
  takeover_outcome: string
}
interface Schedule {
  id: number
  quote_id: number
  insurer_name: string
  certificate_number: string
  member_count: number
  in_force_count: number
  members?: PriorMember[]
}
interface MatchSummary {
  schedule_id: number
  total: number
  in_force: number
  matched_by_id: number
  matched_by_name_and_dob: number
  unmatched: number
  continuation_no_evidence: number
  continuation_with_loading: number
  new_evidence_required: number
}

const schedule = ref<Schedule | null>(null)
const summary = ref<MatchSummary | null>(null)
const busy = ref(false)
const error = ref('')
const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const form = ref({
  insurer_name: '',
  certificate_number: '',
  effective_date: '',
  expiry_date: '',
  notes: '',
  file: [] as File[]
})

const outcomeLabel = (o: string) => {
  if (o === 'continuation_no_evidence') return 'Auto-continue'
  if (o === 'continuation_with_loading') return 'Continue + loading'
  if (o === 'new_evidence_required') return 'Need evidence'
  return 'Unmatched'
}
const outcomeColor = (o: string) => {
  if (o === 'continuation_no_evidence') return 'success'
  if (o === 'continuation_with_loading') return 'warning'
  if (o === 'new_evidence_required') return 'info'
  return 'grey'
}
const formatDate = (s: string) => (s ? new Date(s).toLocaleDateString() : '—')
const formatMoney = (n: number) =>
  n ? Number(n).toLocaleString(undefined, { maximumFractionDigits: 0 }) : '—'

const flash = (msg: string, color = 'success') => {
  snackbarMessage.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

const load = async () => {
  try {
    const res = await GroupPricingService.getPriorInsurerScheduleForQuote(
      quoteId.value
    )
    schedule.value = res.data
  } catch (err: any) {
    if (err?.response?.status !== 404) console.warn(err)
    schedule.value = null
  }
}

const upload = async () => {
  if (!form.value.file?.length) return
  error.value = ''
  busy.value = true
  try {
    const fd = new FormData()
    const file = Array.isArray(form.value.file)
      ? form.value.file[0]
      : form.value.file
    fd.append('file', file)
    fd.append('quote_id', quoteId.value)
    fd.append('insurer_name', form.value.insurer_name)
    fd.append('certificate_number', form.value.certificate_number)
    fd.append('effective_date', form.value.effective_date)
    fd.append('expiry_date', form.value.expiry_date)
    fd.append('notes', form.value.notes)
    const res = await GroupPricingService.uploadPriorInsurerSchedule(fd)
    schedule.value = res.data?.schedule || null
    summary.value = res.data?.match_summary || null
    await load() // reload to pick up members
    flash('Schedule uploaded and matched')
  } catch (err: any) {
    error.value = err?.response?.data || err?.message || 'Upload failed'
  } finally {
    busy.value = false
  }
}

const rematch = async () => {
  if (!schedule.value) return
  busy.value = true
  try {
    const res = await GroupPricingService.rematchPriorInsurerSchedule(
      schedule.value.id
    )
    summary.value = res.data
    await load()
    flash('Rematch complete')
  } catch (err: any) {
    flash(err?.response?.data || 'Rematch failed', 'error')
  } finally {
    busy.value = false
  }
}

const applyTerms = async () => {
  if (!schedule.value) return
  busy.value = true
  try {
    const res = await GroupPricingService.applyTakeoverTermsToCases(
      schedule.value.id
    )
    flash(`Takeover terms applied to ${res.data?.touched ?? 0} cases`)
  } catch (err: any) {
    flash(err?.response?.data || 'Apply failed', 'error')
  } finally {
    busy.value = false
  }
}

onMounted(load)
</script>
