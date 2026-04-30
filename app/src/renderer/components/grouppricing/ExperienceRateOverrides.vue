<template>
  <base-card :show-actions="false">
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-table-edit" class="mr-3" />
        <span>Experience Rate Overrides</span>
        <v-spacer />
        <v-chip
          :color="entries.length > 0 ? 'success' : 'warning'"
          variant="outlined"
          size="small"
        >
          {{ entries.length }} {{ entries.length === 1 ? 'entry' : 'entries' }}
        </v-chip>
      </div>
    </template>
    <template #default>
      <p class="text-body-2 text-medium-emphasis mb-2">
        For each scheme category, choose whether each benefit uses the
        theoretical rate (no experience adjustment) or an experience-rated
        override. Experience-rated entries directly replace the
        experience-adjusted loaded rate used downstream to compute risk and
        final premiums. Rates are annual rate fractions
        (e.g.&nbsp;<code>0.0025</code>&nbsp;= 0.25%).
      </p>
      <v-alert
        v-if="isQuoteLocked"
        type="warning"
        variant="tonal"
        density="compact"
        class="mb-3"
        prepend-icon="mdi-lock-outline"
      >
        This quote is <strong>{{ lockedStatusLabel }}</strong> — experience-rate
        overrides are read-only. Edit the quote (or revert it to draft) to
        change overrides.
      </v-alert>
      <v-alert
        v-else-if="!hasResults"
        type="info"
        variant="tonal"
        density="compact"
        class="mb-3"
      >
        Run calculations once with no overrides to see the baseline
        <strong>Loaded Rate</strong> per benefit. You can then enter overrides
        anchored to those values and re-run.
      </v-alert>

      <v-row class="mt-1" align="start" no-gutters>
        <v-col cols="12" md="5" class="pe-md-2">
          <v-select
            v-model="selectedCategory"
            :items="availableCategories"
            label="Scheme Category"
            density="compact"
            variant="outlined"
            clearable
            hide-details
          />
        </v-col>
        <v-col cols="12" md="3" class="pe-md-2 mt-2 mt-md-0">
          <v-text-field
            v-model.number="credibility"
            label="Credibility (0–1)"
            type="number"
            min="0"
            max="1"
            step="0.05"
            density="compact"
            variant="outlined"
            hide-details="auto"
            :readonly="isQuoteLocked"
            :disabled="isQuoteLocked || credibilitySaving"
            :loading="credibilitySaving"
            :error-messages="credibilityError"
            @blur="saveCredibility"
            @keydown.enter="saveCredibility"
          />
          <div
            v-if="!isQuoteLocked"
            class="text-caption text-medium-emphasis mt-1"
          >
            Default credibility for new override rows. Per-row values can
            deviate.
          </div>
        </v-col>
        <v-col cols="12" md="4" class="text-md-right mt-2 mt-md-0 pt-1">
          <v-btn
            color="primary"
            rounded
            size="small"
            :disabled="!selectedCategory"
            @click="openAddDialog"
          >
            <v-icon
              start
              :icon="isQuoteLocked ? 'mdi-eye-outline' : 'mdi-plus'"
            />
            {{
              isQuoteLocked ? 'View Experience Rates' : 'Add Experience Rates'
            }}
          </v-btn>
        </v-col>
      </v-row>

      <div v-if="entries.length > 0" class="mt-4">
        <div class="d-flex align-center mb-2">
          <v-spacer />
          <v-btn
            size="x-small"
            variant="text"
            density="compact"
            @click="expandAllCategories"
          >
            Expand all
          </v-btn>
          <v-btn
            size="x-small"
            variant="text"
            density="compact"
            class="ml-1"
            @click="collapseAllCategories"
          >
            Collapse all
          </v-btn>
        </div>
        <v-expansion-panels
          v-model="openCategoryPanels"
          multiple
          variant="accordion"
          class="ero-panels"
        >
          <v-expansion-panel
            v-for="group in groupedEntries"
            :key="group.category"
            :value="group.category"
          >
            <v-expansion-panel-title>
              <span class="text-body-2 font-weight-medium">
                {{ group.category }}
              </span>
              <v-spacer />
              <v-chip size="x-small" variant="tonal" class="mr-3">
                {{ group.rows.length }}
                {{ group.rows.length === 1 ? 'benefit' : 'benefits' }}
              </v-chip>
            </v-expansion-panel-title>
            <v-expansion-panel-text class="ero-panel-text">
              <v-table density="compact" class="ero-table">
                <thead>
                  <tr>
                    <th class="text-left">Benefit</th>
                    <th class="text-left">Mode</th>
                    <th class="text-right">Override Rate</th>
                    <th class="text-right">Credibility</th>
                    <th class="text-left">Updated By</th>
                    <th class="text-left">Updated At</th>
                    <th class="text-center" style="width: 56px">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="row in group.rows"
                    :key="`${row.scheme_category}|${row.benefit}`"
                  >
                    <td>
                      <div class="d-flex align-center">
                        <v-chip
                          size="x-small"
                          variant="tonal"
                          class="ero-benefit-chip"
                        >
                          {{ aliasCodeForBenefit(row.benefit) }}
                        </v-chip>
                        <span class="text-body-2">
                          {{
                            labelForBenefit(row.scheme_category, row.benefit)
                          }}
                        </span>
                      </div>
                    </td>
                    <td>
                      <span v-if="row.mode === 'experience_rated'">
                        Experience Rated
                      </span>
                      <span v-else class="text-medium-emphasis">
                        Technical
                      </span>
                    </td>
                    <td class="text-right ero-mono">
                      <span v-if="row.mode === 'experience_rated'">
                        {{ formatRate(row.override_rate) }}
                      </span>
                      <span v-else class="text-medium-emphasis">—</span>
                    </td>
                    <td class="text-right ero-mono">
                      <span
                        v-if="row.credibility != null && row.credibility > 0"
                      >
                        {{ Number(row.credibility).toFixed(2) }}
                      </span>
                      <span v-else class="text-medium-emphasis">—</span>
                    </td>
                    <td class="text-medium-emphasis">
                      {{ row.updated_by || '—' }}
                    </td>
                    <td class="text-medium-emphasis">
                      {{ formatDate(row.updated_at) || '—' }}
                    </td>
                    <td class="text-center">
                      <v-btn
                        v-if="!isQuoteLocked"
                        size="x-small"
                        color="error"
                        variant="text"
                        :icon="true"
                        @click="
                          removeEntryByKey(row.scheme_category, row.benefit)
                        "
                      >
                        <v-icon>mdi-delete</v-icon>
                      </v-btn>
                      <v-icon
                        v-else
                        size="small"
                        color="grey-lighten-1"
                        icon="mdi-lock-outline"
                      />
                    </td>
                  </tr>
                </tbody>
              </v-table>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </div>

      <p v-else class="text-body-2 text-medium-emphasis mt-4 mb-0">
        No overrides yet. Pick a category and click
        <strong>+ Add Experience Rates</strong> to start.
      </p>

      <div v-if="entries.length > 0 && !isQuoteLocked" class="mt-4 d-flex">
        <v-btn
          color="primary"
          rounded
          size="small"
          :loading="saving"
          @click="saveAll"
        >
          Save Overrides
        </v-btn>
        <v-btn
          color="red"
          rounded
          size="small"
          variant="outlined"
          class="ml-3"
          :disabled="saving"
          @click="deleteAll"
        >
          Delete All
        </v-btn>
      </div>

      <v-dialog v-model="dialogOpen" max-width="780" persistent scrollable>
        <v-card>
          <v-card-title class="d-flex align-center">
            <span>Experience Rates — {{ dialogCategory }}</span>
            <v-spacer />
            <v-chip
              v-if="!hasResults"
              color="warning"
              size="small"
              variant="tonal"
            >
              Run calculations to see Loaded Rates
            </v-chip>
            <v-chip v-else color="success" size="small" variant="tonal">
              Baseline from latest calculation
            </v-chip>
          </v-card-title>
          <v-divider />
          <v-card-text>
            <p class="text-body-2 text-medium-emphasis mb-4">
              Technical leaves the experience-adjusted rate equal to the loaded
              rate (no adjustment). Experience-rated replaces it with your
              supplied annual fraction.
            </p>

            <!-- Aligned grid: each benefit row uses the same column widths -->
            <div class="ero-grid-header text-caption text-medium-emphasis mb-1">
              <div>Benefit</div>
              <div>Mode</div>
              <div class="text-right">Loaded Rate (latest run)</div>
              <div>Override Rate</div>
              <div>Credibility</div>
            </div>
            <v-divider class="mb-2" />

            <div
              v-for="(b, i) in dialogBenefits"
              :key="b.code"
              class="ero-grid-row"
              :class="{ 'ero-grid-row-alt': i % 2 === 1 }"
            >
              <div class="d-flex align-center">
                <v-chip size="small" variant="tonal" class="ero-benefit-chip">
                  {{ aliasCodeForBenefit(b.code) }}
                </v-chip>
                <span class="text-body-2">{{ b.label }}</span>
              </div>

              <v-radio-group
                v-model="b.mode"
                inline
                density="compact"
                hide-details
                class="ero-radio"
                :disabled="isQuoteLocked"
              >
                <v-radio label="Technical" value="theoretical" />
                <v-radio label="Experience-rated" value="experience_rated" />
              </v-radio-group>

              <div class="text-right ero-mono">
                <template v-if="b.loadedRate !== null">
                  <div>{{ formatRate(b.loadedRate) }}</div>
                  <div class="text-caption text-medium-emphasis">
                    {{ loadedRateCaption(b.code, b.loadedRate) }}
                  </div>
                </template>
                <span v-else class="text-medium-emphasis text-caption">
                  Run calc to view
                </span>
              </div>

              <div>
                <v-text-field
                  v-if="b.mode === 'experience_rated'"
                  v-model.number="b.override_rate"
                  type="number"
                  step="0.0001"
                  min="0"
                  density="compact"
                  variant="outlined"
                  hide-details
                  placeholder="0.0025"
                  :readonly="isQuoteLocked"
                  :disabled="isQuoteLocked"
                />
                <span v-else class="text-medium-emphasis text-caption">—</span>
              </div>

              <div>
                <v-text-field
                  v-model.number="b.credibility"
                  type="number"
                  step="0.05"
                  min="0"
                  max="1"
                  density="compact"
                  variant="outlined"
                  hide-details
                  placeholder="0.00"
                  :readonly="isQuoteLocked"
                  :disabled="isQuoteLocked"
                />
              </div>
            </div>

            <div
              v-if="dialogAuditEntries.length > 0"
              class="mt-4 pa-3 ero-audit"
            >
              <div class="text-caption text-medium-emphasis mb-1">
                Audit trail
              </div>
              <div
                v-for="a in dialogAuditEntries"
                :key="`${a.benefit}-${a.updated_at}`"
                class="text-caption"
              >
                <strong>{{ a.benefit }}</strong> — last updated by
                <span class="text-medium-emphasis">{{
                  a.updated_by || 'system'
                }}</span>
                on
                <span class="text-medium-emphasis">{{
                  formatDate(a.updated_at) || '—'
                }}</span>
              </div>
            </div>
          </v-card-text>
          <v-divider />
          <v-card-actions>
            <v-spacer />
            <v-btn variant="text" @click="dialogOpen = false">
              {{ isQuoteLocked ? 'Close' : 'Cancel' }}
            </v-btn>
            <v-btn v-if="!isQuoteLocked" color="primary" @click="commitDialog">
              Add / Update
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </template>
  </base-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import { useFlashStore } from '@/renderer/store/flash'

type Mode = 'experience_rated' | 'theoretical'
type Benefit = 'GLA' | 'AAGLA' | 'SGLA' | 'PTD' | 'TTD' | 'PHI' | 'CI' | 'FUN'

interface OverrideRow {
  id?: number
  quote_id: number
  scheme_category: string
  benefit: Benefit
  mode: Mode
  override_rate: number
  // Per-row credibility (0-1). Recorded on every calc run into
  // HistoricalCredibilityData's per-benefit credibility columns.
  credibility?: number
  created_at?: string
  created_by?: string
  updated_at?: string
  updated_by?: string
}

interface DialogBenefit {
  code: Benefit
  label: string
  mode: Mode
  override_rate: number
  credibility: number
  loadedRate: number | null
  updated_at?: string
  updated_by?: string
}

const props = withDefaults(
  defineProps<{
    quote: Record<string, any>
    resultSummaries?: any[]
  }>(),
  { resultSummaries: () => [] }
)

const emit = defineEmits<{
  (e: 'overrides-updated', count: number): void
}>()

const flash = useFlashStore()

const entries = ref<OverrideRow[]>([])
const selectedCategory = ref<string | null>(null)
const dialogOpen = ref(false)
const dialogCategory = ref<string>('')
const dialogBenefits = ref<DialogBenefit[]>([])
const saving = ref(false)
// Which category panels are currently open. Default to all expanded after
// each load; user can collapse individually or use the bulk buttons.
const openCategoryPanels = ref<string[]>([])

// Manually-entered credibility (0-1) the actuary supplies alongside the
// override rates. Loaded from the quote on mount, saved on blur or Enter.
// Persisted on the quote so it survives recalcs and is recorded against
// HistoricalCredibilityData on each calc run.
const credibility = ref<number | null>(null)
const credibilitySaving = ref(false)
const credibilityError = ref<string | string[] | undefined>(undefined)
// Mirror of the last value we successfully persisted to the server. Used by
// saveCredibility() to skip redundant writes without mutating the (read-only)
// `quote` prop directly.
const lastPersistedCredibility = ref<number>(0)

// System-wide benefit configuration (GroupBenefitMapper) — supplies the
// customised short code (`benefit_alias_code`) and long-form alias
// (`benefit_alias`) for each canonical code (GLA, SGLA, PTD, TTD, PHI, CI,
// FUN). Loaded once on mount; falls back to the canonical code / hardcoded
// label if the map is missing or has empty fields.
interface BenefitMap {
  benefit_code: string
  benefit_alias: string
  benefit_alias_code: string
  benefit_name: string
}
const benefitMaps = ref<BenefitMap[]>([])

const benefitFlagMap: Record<Benefit, string> = {
  GLA: 'gla_benefit',
  AAGLA: 'additional_accidental_gla_benefit',
  SGLA: 'sgla_benefit',
  PTD: 'ptd_benefit',
  TTD: 'ttd_benefit',
  PHI: 'phi_benefit',
  CI: 'ci_benefit',
  FUN: 'family_funeral_benefit'
}

// Fallback labels used when the scheme category has no customised alias for
// a benefit. The actual UI label resolves to:
//   <category>.<b>_alias  (if non-empty)
//   benefitLabel[code]    (otherwise)
//   code                  (final fallback)
const benefitLabel: Record<Benefit, string> = {
  GLA: 'Group Life Assurance',
  AAGLA: 'Additional Accidental GLA',
  SGLA: 'Spouse Group Life',
  PTD: 'Permanent Total Disability',
  TTD: 'Temporary Total Disability',
  PHI: 'Permanent Health Insurance',
  CI: 'Critical Illness',
  FUN: 'Family Funeral'
}

// Field on SchemeCategory that holds the per-benefit customised alias.
const benefitAliasFieldMap: Record<Benefit, string> = {
  GLA: 'gla_alias',
  // No per-category alias for AAGLA on SchemeCategory; falls through to
  // global benefit-config alias / hardcoded label.
  AAGLA: '',
  SGLA: 'sgla_alias',
  PTD: 'ptd_alias',
  TTD: 'ttd_alias',
  PHI: 'phi_alias',
  CI: 'ci_alias',
  FUN: 'family_funeral_alias'
}

const findBenefitMap = (code: Benefit): BenefitMap | undefined =>
  benefitMaps.value.find((m) => (m.benefit_code || '').toUpperCase() === code)

// Long-form label resolution priority:
//   1. Per-category alias on SchemeCategory (e.g. `gla_alias`)
//   2. Global Benefit Configuration `benefit_alias` (GroupBenefitMapper)
//   3. Hardcoded fallback (`benefitLabel`)
//   4. Canonical code (final fallback)
const labelForBenefit = (category: string, code: Benefit): string => {
  const record = findCategoryRecord(category) as Record<string, any> | undefined
  const categoryAlias = record?.[benefitAliasFieldMap[code]]
  if (typeof categoryAlias === 'string' && categoryAlias.trim().length > 0) {
    return categoryAlias.trim()
  }
  const globalAlias = findBenefitMap(code)?.benefit_alias
  if (typeof globalAlias === 'string' && globalAlias.trim().length > 0) {
    return globalAlias.trim()
  }
  return benefitLabel[code] ?? code
}

// Short chip code resolution priority:
//   1. Global Benefit Configuration `benefit_alias_code` (e.g. "SI" for CI)
//   2. Canonical code (GLA, SGLA, PTD, TTD, PHI, CI, FUN)
const aliasCodeForBenefit = (code: Benefit): string => {
  const aliasCode = findBenefitMap(code)?.benefit_alias_code
  if (typeof aliasCode === 'string' && aliasCode.trim().length > 0) {
    return aliasCode.trim()
  }
  return code
}

// Formatting helpers — declared as `const` arrow functions (not hoisted
// `function` declarations) so Vite HMR rebinds them reliably on every accept.
// Hoisted function declarations sometimes fail to land on the `$setup` proxy
// after incremental edits, requiring a hard reload to resurface.
const formatRate = (value: number | null | undefined): string => {
  if (value === null || value === undefined || !isFinite(value)) return '—'
  return Number(value).toFixed(6)
}

const formatDate = (value: string | null | undefined): string => {
  if (!value) return ''
  try {
    const d = new Date(value)
    if (isNaN(d.getTime())) return ''
    return (
      d.toLocaleDateString() +
      ' ' +
      d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    )
  } catch {
    return ''
  }
}

const availableCategories = computed<string[]>(() => {
  const cats = props.quote?.selected_scheme_categories ?? []
  return Array.isArray(cats) ? [...new Set(cats.filter(Boolean))] : []
})

const hasResults = computed<boolean>(
  () => Array.isArray(props.resultSummaries) && props.resultSummaries.length > 0
)

// Once a quote is approved, accepted, or in-force it represents a settled
// commercial position. Editing experience-rate overrides at that point would
// silently change the priced premium without a fresh approval cycle, so the
// UI (and the backend save/delete endpoints) reject mutations from this
// state forwards.
const LOCKED_STATUSES = new Set(['approved', 'accepted', 'in_force'])
const isQuoteLocked = computed<boolean>(() => {
  const status = String(props.quote?.status || '').toLowerCase()
  return LOCKED_STATUSES.has(status)
})
const lockedStatusLabel = computed<string>(() => {
  const raw = String(props.quote?.status || '')
  if (!raw) return ''
  return raw.replace(/_/g, ' ')
})

const dialogAuditEntries = computed(() => {
  return dialogBenefits.value
    .filter((b) => b.updated_at || b.updated_by)
    .map((b) => ({
      benefit: b.code,
      updated_at: b.updated_at,
      updated_by: b.updated_by
    }))
})

// Group entries by scheme_category for the collapsible category panels.
// Categories preserve first-seen order; rows within a category preserve
// insertion order from the master `entries` array.
const groupedEntries = computed<
  Array<{ category: string; rows: OverrideRow[] }>
>(() => {
  const map = new Map<string, OverrideRow[]>()
  for (const row of entries.value) {
    const key = row.scheme_category
    if (!map.has(key)) map.set(key, [])
    map.get(key)!.push(row)
  }
  return Array.from(map, ([category, rows]) => ({ category, rows }))
})

const expandAllCategories = () => {
  openCategoryPanels.value = groupedEntries.value.map((g) => g.category)
}

const collapseAllCategories = () => {
  openCategoryPanels.value = []
}

function findCategoryRecord(category: string) {
  return (props.quote?.scheme_categories ?? []).find(
    (c: any) => c?.scheme_category === category
  )
}

function findSummaryForCategory(category: string) {
  if (!Array.isArray(props.resultSummaries)) return null
  return (
    props.resultSummaries.find(
      (s: any) => s?.category === category || s?.Category === category
    ) ?? null
  )
}

// Weighted average of loaded_<b>_rate per category. Mathematically:
//   weighted_avg = Σ(LoadedRate_i × weight_i) / Σ(weight_i)
// For SA-based benefits the weight is capped sum assured, so
//   total_<b>_annual_risk_premium = Σ(LoadedRate × CappedSA)
//   weighted_avg = total_premium / total_capped_sum_assured
// For PHI/TTD the weight is covered income (capped income aggregate). The
// TTD premium has an extra NumberMonthlyPayments factor baked in
// (premium = LoadedRate × MonthlyIncome × Months), so we divide it out to
// recover the underlying loaded_ttd_rate. The fields used here are exactly
// the ones the rating service actually populates — TotalPhiMonthlyBenefit
// and the *_risk_rate_per_1000_sa fields for PHI/TTD are declared on the
// model but never assigned, so we avoid them.
const TTD_MONTHLY_PAYMENTS = 12

function loadedRateForBenefit(category: string, code: Benefit): number | null {
  const summary = findSummaryForCategory(category) as any
  if (!summary) return null

  let premium = 0
  let denom = 0

  switch (code) {
    case 'GLA':
      premium = Number(summary.total_gla_annual_risk_premium) || 0
      denom = Number(summary.total_gla_capped_sum_assured) || 0
      break
    case 'AAGLA':
      premium =
        Number(summary.total_additional_accidental_gla_annual_risk_premium) || 0
      denom =
        Number(summary.total_additional_accidental_gla_capped_sum_assured) || 0
      break
    case 'SGLA':
      premium = Number(summary.total_sgla_annual_risk_premium) || 0
      denom = Number(summary.total_sgla_capped_sum_assured) || 0
      break
    case 'PTD':
      premium = Number(summary.total_ptd_annual_risk_premium) || 0
      denom = Number(summary.total_ptd_capped_sum_assured) || 0
      break
    case 'CI':
      premium = Number(summary.total_ci_annual_risk_premium) || 0
      denom = Number(summary.total_ci_capped_sum_assured) || 0
      break
    case 'FUN':
      premium = Number(summary.total_fun_annual_risk_premium) || 0
      denom = Number(summary.total_family_funeral_sum_assured) || 0
      break
    case 'PHI':
      // premium = LoadedPhiRate × MonthlyBenefit. TotalPhiCappedIncome is
      // the monthly capped-income aggregate, which is the override's
      // multiplier proxy and the same weight basis as MonthlyBenefit.
      premium = Number(summary.total_phi_annual_risk_premium) || 0
      denom = Number(summary.total_phi_capped_income) || 0
      break
    case 'TTD':
      // premium = LoadedTtdRate × MonthlyIncome × NumberMonthlyPayments.
      // Divide by NumberMonthlyPayments to recover the LoadedTtdRate basis,
      // then weight by monthly capped-income aggregate.
      premium =
        (Number(summary.total_ttd_annual_risk_premium) || 0) /
        TTD_MONTHLY_PAYMENTS
      denom = Number(summary.total_ttd_capped_income) || 0
      break
  }

  if (denom <= 0 || !isFinite(premium) || premium <= 0) return null
  return premium / denom
}

// Caption shown under the loaded-rate value. SA-based benefits use the
// familiar "/ 1000 SA" representation; PHI/TTD show "/ 1000 covered income"
// because their premium is income-driven, not SA-driven.
function loadedRateCaption(code: Benefit, value: number): string {
  if (code === 'PHI' || code === 'TTD') {
    return `(${(value * 1000).toFixed(4)} / 1000 covered income)`
  }
  return `(${(value * 1000).toFixed(4)} / 1000 SA)`
}

function benefitsForCategory(category: string): Benefit[] {
  const record = findCategoryRecord(category)
  if (!record) {
    return Object.keys(benefitFlagMap) as Benefit[]
  }
  return (Object.keys(benefitFlagMap) as Benefit[]).filter(
    (b) => !!record[benefitFlagMap[b]]
  )
}

function openAddDialog() {
  if (!selectedCategory.value) return
  const cat = selectedCategory.value
  const enabled = benefitsForCategory(cat)
  if (enabled.length === 0) {
    flash.show(`Category "${cat}" has no benefits selected.`, 'warning')
    return
  }
  dialogCategory.value = cat
  // Default credibility for new rows = the quote-level value the user
  // entered above. Existing rows keep whatever they were previously saved
  // with so per-row deviations survive re-opens.
  const defaultCredibility =
    Number(props.quote?.experience_override_credibility) || 0
  dialogBenefits.value = enabled.map<DialogBenefit>((code) => {
    const existing = entries.value.find(
      (r) => r.scheme_category === cat && r.benefit === code
    )
    return {
      code,
      label: labelForBenefit(cat, code),
      mode: existing?.mode ?? 'theoretical',
      override_rate: existing?.override_rate ?? 0,
      credibility:
        existing?.credibility !== undefined && existing?.credibility !== null
          ? Number(existing.credibility)
          : defaultCredibility,
      loadedRate: loadedRateForBenefit(cat, code),
      updated_at: existing?.updated_at,
      updated_by: existing?.updated_by
    }
  })
  dialogOpen.value = true
}

function commitDialog() {
  if (isQuoteLocked.value) {
    flash.show(
      `Quote is ${lockedStatusLabel.value} — overrides cannot be edited.`,
      'warning'
    )
    return
  }
  const cat = dialogCategory.value
  for (const b of dialogBenefits.value) {
    if (
      b.mode === 'experience_rated' &&
      (b.override_rate === null ||
        b.override_rate === undefined ||
        Number(b.override_rate) < 0)
    ) {
      flash.show(
        `${b.code} requires a non-negative override rate when experience-rated.`,
        'warning'
      )
      return
    }
    const credValue = Number(b.credibility)
    if (!isFinite(credValue) || credValue < 0 || credValue > 1) {
      flash.show(`${b.code} credibility must be between 0 and 1.`, 'warning')
      return
    }
  }
  for (const b of dialogBenefits.value) {
    const idx = entries.value.findIndex(
      (r) => r.scheme_category === cat && r.benefit === b.code
    )
    const row: OverrideRow = {
      quote_id: props.quote.id,
      scheme_category: cat,
      benefit: b.code,
      mode: b.mode,
      override_rate:
        b.mode === 'experience_rated' ? Number(b.override_rate) || 0 : 0,
      credibility: Number(b.credibility) || 0,
      // Preserve audit trail until next save round-trips it from the server
      created_at: entries.value[idx]?.created_at,
      created_by: entries.value[idx]?.created_by,
      updated_at: entries.value[idx]?.updated_at,
      updated_by: entries.value[idx]?.updated_by
    }
    if (idx >= 0) entries.value[idx] = row
    else entries.value.push(row)
  }
  dialogOpen.value = false
  flash.show(
    `Updated overrides for ${cat}. Click Save Overrides to persist.`,
    'success'
  )
}

function removeEntryByKey(category: string, benefit: Benefit) {
  if (isQuoteLocked.value) return
  const idx = entries.value.findIndex(
    (r) => r.scheme_category === category && r.benefit === benefit
  )
  if (idx >= 0) entries.value.splice(idx, 1)
}

async function load() {
  if (!props.quote?.id) return
  try {
    const res = await GroupPricingService.getExperienceRateOverrides(
      props.quote.id
    )
    entries.value = Array.isArray(res?.data) ? res.data : []
    expandAllCategories()
  } catch (err) {
    flash.show('Failed to load experience-rate overrides.', 'error')
  }
}

// Read the persisted credibility off the quote. The field defaults to 0 on
// the backend; we surface null in the input so the placeholder shows when
// nothing has been entered yet. Also seeds the local "last persisted" mirror
// so saveCredibility() can skip writes when the value is unchanged.
function syncCredibilityFromQuote() {
  const raw = Number(props.quote?.experience_override_credibility)
  const numeric = isFinite(raw) ? raw : 0
  credibility.value = numeric > 0 ? numeric : null
  lastPersistedCredibility.value = numeric
}

async function saveCredibility() {
  if (isQuoteLocked.value) return
  if (!props.quote?.id) return
  const value =
    credibility.value === null || credibility.value === undefined
      ? 0
      : Number(credibility.value)
  if (!isFinite(value) || value < 0 || value > 1) {
    credibilityError.value = 'Must be between 0 and 1'
    return
  }
  credibilityError.value = undefined
  // Skip the round-trip when nothing actually changed (typical onBlur after
  // tabbing through without editing).
  if (Math.abs(lastPersistedCredibility.value - value) < 1e-9) return
  credibilitySaving.value = true
  try {
    await GroupPricingService.updateExperienceOverrideCredibility(
      props.quote.id,
      value
    )
    // Track the new value locally (without mutating the read-only `quote`
    // prop) so the next blur sees it as the persisted baseline and skips
    // the redundant write.
    lastPersistedCredibility.value = value
    flash.show('Credibility saved.', 'success')
  } catch (err: any) {
    credibilityError.value =
      err?.response?.data?.error || 'Failed to save credibility.'
    flash.show(credibilityError.value as string, 'error')
  } finally {
    credibilitySaving.value = false
  }
}

async function loadBenefitMaps() {
  try {
    const res = await GroupPricingService.getBenefitMaps()
    benefitMaps.value = Array.isArray(res?.data) ? res.data : []
  } catch (err) {
    // Non-fatal — fall back to canonical codes / hardcoded labels.
    benefitMaps.value = []
  }
}

async function saveAll() {
  if (isQuoteLocked.value) {
    flash.show(
      `Quote is ${lockedStatusLabel.value} — overrides cannot be saved.`,
      'warning'
    )
    return
  }
  if (entries.value.length === 0) {
    flash.show('Add at least one entry before saving.', 'warning')
    return
  }
  saving.value = true
  try {
    const payload = entries.value.map((r) => ({
      quote_id: props.quote.id,
      scheme_category: r.scheme_category,
      benefit: r.benefit,
      mode: r.mode,
      override_rate: r.mode === 'experience_rated' ? r.override_rate : 0,
      credibility: Number(r.credibility) || 0
    }))
    const res = await GroupPricingService.saveExperienceRateOverrides(payload)
    entries.value = Array.isArray(res?.data) ? res.data : entries.value
    expandAllCategories()
    emit('overrides-updated', entries.value.length)
    flash.show('Experience-rate overrides saved.', 'success')
  } catch (err) {
    flash.show('Failed to save experience-rate overrides.', 'error')
  } finally {
    saving.value = false
  }
}

async function deleteAll() {
  if (isQuoteLocked.value) {
    flash.show(
      `Quote is ${lockedStatusLabel.value} — overrides cannot be deleted.`,
      'warning'
    )
    return
  }
  if (!confirm('Delete every experience-rate override on this quote?')) return
  try {
    await GroupPricingService.deleteExperienceRateOverrides(props.quote.id)
    entries.value = []
    emit('overrides-updated', 0)
    flash.show('All overrides deleted.', 'success')
  } catch (err) {
    flash.show('Failed to delete overrides.', 'error')
  }
}

watch(
  () => props.quote?.id,
  (id) => {
    if (id) {
      load()
      syncCredibilityFromQuote()
    }
  }
)

watch(
  () => props.quote?.experience_override_credibility,
  () => syncCredibilityFromQuote()
)

onMounted(() => {
  loadBenefitMaps()
  load()
  syncCredibilityFromQuote()
})
</script>

<style scoped>
.ero-table {
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 4px;
}

/* Fixed-width chip for the saved-entries table so every benefit name's first
 * letter starts at the same X regardless of how short or long the benefit
 * code (e.g. SI vs FUN vs SGL). */
.ero-benefit-chip {
  min-width: 56px;
  justify-content: center;
  margin-right: 12px;
}

.ero-benefit-chip :deep(.v-chip__content) {
  width: 100%;
  justify-content: center;
}

.ero-panels {
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 4px;
}

.ero-panel-text :deep(.v-expansion-panel-text__wrapper) {
  padding: 0 8px 8px;
}

.ero-mono {
  font-variant-numeric: tabular-nums;
  font-family:
    'JetBrains Mono', 'Fira Code', Consolas, 'Courier New', monospace;
}

/* 5-column aligned grid for the dialog:
 * benefit | mode | loaded rate | override rate | credibility
 */
.ero-grid-header,
.ero-grid-row {
  display: grid;
  grid-template-columns:
    minmax(180px, 1.5fr) minmax(200px, 1.3fr) minmax(140px, 0.9fr)
    minmax(120px, 0.8fr) minmax(110px, 0.7fr);
  align-items: center;
  column-gap: 16px;
}

.ero-grid-row {
  padding: 8px 8px;
  border-radius: 4px;
}

.ero-grid-row-alt {
  background-color: rgba(0, 0, 0, 0.02);
}

.ero-grid-row + .ero-grid-row {
  margin-top: 4px;
}

.ero-radio :deep(.v-selection-control) {
  margin-inline-end: 12px;
}

.ero-radio :deep(.v-label) {
  font-size: 0.85rem;
}

.ero-audit {
  background-color: rgba(0, 0, 0, 0.03);
  border-radius: 4px;
}
</style>
