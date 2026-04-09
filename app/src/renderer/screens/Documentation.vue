<template>
  <div class="docs-root">
    <!-- ── Top bar ──────────────────────────────────────────────── -->
    <div class="docs-topbar">
      <div class="docs-topbar-left">
        <span class="docs-title">Group Pricing — Data Dictionary</span>
        <span class="docs-subtitle">
          {{ totalVisible }} variable{{ totalVisible !== 1 ? 's' : '' }} across
          {{ visibleSections.length }} section{{
            visibleSections.length !== 1 ? 's' : ''
          }}
        </span>
      </div>

      <div class="docs-topbar-right">
        <v-text-field
          v-model="search"
          variant="outlined"
          density="compact"
          prepend-inner-icon="mdi-magnify"
          placeholder="Search variables…"
          clearable
          hide-details
          class="docs-search"
        />
        <v-chip-group
          v-model="sourceFilter"
          mandatory
          class="docs-filter-group"
        >
          <v-chip value="all" size="small" filter variant="tonal">All</v-chip>
          <v-chip
            value="User Input"
            size="small"
            filter
            variant="tonal"
            color="blue"
            >User Input</v-chip
          >
          <v-chip
            value="Calculation Engine"
            size="small"
            filter
            variant="tonal"
            color="deep-purple"
            >Formulas</v-chip
          >
        </v-chip-group>
      </div>
    </div>

    <!-- ── Body: nav + content ──────────────────────────────────── -->
    <div class="docs-body">
      <!-- Left nav -->
      <nav class="docs-nav">
        <div
          v-for="section in visibleSections"
          :key="section.id"
          class="docs-nav-item"
          :class="{ 'docs-nav-item--active': activeSection === section.id }"
          @click="scrollToSection(section.id)"
        >
          <span class="docs-nav-label">{{ section.label }}</span>
          <span class="docs-nav-count">{{ section.visibleItems.length }}</span>
        </div>
        <div v-if="visibleSections.length === 0" class="docs-nav-empty">
          No matches
        </div>
      </nav>

      <!-- Scrollable content -->
      <main ref="contentEl" class="docs-content" @scroll.passive="onScroll">
        <template v-for="section in visibleSections" :key="section.id">
          <section :id="`section-${section.id}`" class="docs-section">
            <div class="docs-section-header">
              <span class="docs-section-title">{{ section.label }}</span>
              <span class="docs-section-count">
                {{ section.visibleItems.length }} var{{
                  section.visibleItems.length !== 1 ? 's' : ''
                }}
              </span>
            </div>

            <div
              v-for="item in section.visibleItems"
              :key="item.data_variable"
              class="var-card"
            >
              <!-- Human name + badges -->
              <div class="var-card-header">
                <span class="var-human-name">{{
                  toTitle(item.data_variable)
                }}</span>
                <div class="var-badges">
                  <v-chip size="x-small" variant="outlined" class="mr-1">{{
                    item.data_type
                  }}</v-chip>
                  <v-chip
                    size="x-small"
                    variant="tonal"
                    :color="
                      item.data_source_type === 'User Input'
                        ? 'blue'
                        : 'deep-purple'
                    "
                  >
                    {{
                      item.data_source_type === 'Calculation Engine'
                        ? 'Formula'
                        : 'User Input'
                    }}
                  </v-chip>
                </div>
              </div>

              <!-- Identifier + copy button -->
              <div class="var-identifier-row">
                <code class="var-identifier">{{ item.data_variable }}</code>
                <v-btn
                  :icon="true"
                  size="x-small"
                  variant="text"
                  :color="copied === item.data_variable ? 'success' : undefined"
                  class="ml-1"
                  @click="copyVar(item.data_variable)"
                >
                  <v-icon size="13">
                    {{
                      copied === item.data_variable
                        ? 'mdi-check'
                        : 'mdi-content-copy'
                    }}
                  </v-icon>
                </v-btn>
              </div>

              <!-- Description -->
              <p class="var-description">{{ item.data_description }}</p>

              <!-- Formula block — Calculation Engine -->
              <div
                v-if="item.data_source_type === 'Calculation Engine'"
                class="formula-block"
              >
                <div class="formula-header">
                  <v-icon size="12" class="mr-1">mdi-function-variant</v-icon>
                  Formula
                </div>
                <pre
                  class="formula-pre"
                ><code>{{ item.data_source }}</code></pre>
              </div>

              <!-- Source ref — User Input -->
              <div v-else class="source-ref">
                <v-icon size="13" color="blue-darken-1" class="mr-1"
                  >mdi-database-import-outline</v-icon
                >
                <span>{{ item.data_source.trim() }}</span>
              </div>
            </div>
          </section>
        </template>

        <!-- Empty state -->
        <div v-if="visibleSections.length === 0" class="docs-empty">
          <v-icon size="52" color="grey-lighten-1" class="mb-3"
            >mdi-text-search-variant</v-icon
          >
          <div class="docs-empty-title">No variables match "{{ search }}"</div>
          <div class="docs-empty-sub"
            >Try a different term or clear the filters below</div
          >
          <v-btn
            variant="tonal"
            size="small"
            class="mt-4"
            @click="clearFilters"
          >
            Clear filters
          </v-btn>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { groupPricing } from '../data/group_pricing'

// ── State ─────────────────────────────────────────────────────────

const search = ref('')
const sourceFilter = ref('all')
const activeSection = ref('')
const copied = ref<string | null>(null)
const contentEl = ref<HTMLElement | null>(null)

// ── Section definitions ───────────────────────────────────────────

const SECTION_DEFS = [
  {
    id: 'free-cover',
    label: 'Free Cover Limit',
    variables: [
      'member_distribution_free_cover_limit',
      'free_cover_limit_scaling_factor',
      'free_cover_limit_percentile',
      'free_cover_limit_nearest_multiple',
      'free_cover_limit'
    ]
  },
  {
    id: 'dates',
    label: 'Scheme Dates',
    variables: ['commencement_date', 'entry_date', 'exit_date']
  },
  {
    id: 'educator',
    label: 'Educator Benefit',
    variables: [
      'education_level',
      'max_coverage_period',
      'max_tuition_per_year',
      'max_book_allowance_proportion',
      'max_book_allowance_amount',
      'max_book_allowance',
      'max_accommodation_allowance_proportion',
      'max_accommodation_allowance_amount',
      'max_accommodation_allowance',
      'grade0_sum_assured',
      'grade1_7_sum_assured',
      'grade8_12_sum_assured',
      'tertiary_sum_assured',
      'educator_risk_premium'
    ]
  },
  {
    id: 'demographics',
    label: 'Member Demographics',
    variables: ['spouse_age_gap', 'min_age', 'max_age']
  },
  {
    id: 'gla',
    label: 'GLA & Reinsurance',
    variables: [
      'gla_terminal_illness_loading_rate',
      'is_lumpsum_reins_gla_dependent'
    ]
  },
  {
    id: 'pricing',
    label: 'Pricing Parameters',
    variables: [
      'premium_rates_guaranteed_period_months',
      'quote_validity_period_months',
      'annual_expense_amount'
    ]
  },
  {
    id: 'experience',
    label: 'Experience Rating',
    variables: [
      'full_credibility_threshold',
      'credibility',
      'blended_gla_rate',
      'annual_experience_weighted_rate',
      'gla_experience_adjustment',
      'gla_theoretical_rate',
      'mannually_added_credibility'
    ]
  }
]

const itemMap = Object.fromEntries(
  groupPricing.map((item: any) => [item.data_variable, item])
)

// ── Filtering ─────────────────────────────────────────────────────

const matchesFilters = (item: any): boolean => {
  if (!item) return false
  if (
    sourceFilter.value !== 'all' &&
    item.data_source_type !== sourceFilter.value
  )
    return false
  const q = search.value.toLowerCase().trim()
  if (!q) return true
  return (
    item.data_variable.toLowerCase().includes(q) ||
    item.data_description.toLowerCase().includes(q) ||
    item.data_source.toLowerCase().includes(q)
  )
}

const visibleSections = computed(() =>
  SECTION_DEFS.map((def) => ({
    ...def,
    visibleItems: def.variables.map((v) => itemMap[v]).filter(matchesFilters)
  })).filter((s) => s.visibleItems.length > 0)
)

const totalVisible = computed(() =>
  visibleSections.value.reduce((n, s) => n + s.visibleItems.length, 0)
)

// Reset scroll position and active section whenever filters change
watch([search, sourceFilter], async () => {
  await nextTick()
  if (contentEl.value) contentEl.value.scrollTop = 0
  activeSection.value = visibleSections.value[0]?.id ?? ''
})

// ── Scroll spy ────────────────────────────────────────────────────

const onScroll = () => {
  if (!contentEl.value) return
  const BUFFER = 80
  const scrollTop = contentEl.value.scrollTop
  let current = visibleSections.value[0]?.id ?? ''

  for (const section of visibleSections.value) {
    const el = contentEl.value.querySelector(
      `#section-${section.id}`
    ) as HTMLElement | null
    if (el && el.offsetTop - BUFFER <= scrollTop) {
      current = section.id
    }
  }
  activeSection.value = current
}

const scrollToSection = async (id: string) => {
  activeSection.value = id
  await nextTick()
  if (!contentEl.value) return
  const el = contentEl.value.querySelector(
    `#section-${id}`
  ) as HTMLElement | null
  if (el)
    contentEl.value.scrollTo({ top: el.offsetTop - 16, behavior: 'smooth' })
}

// ── Utilities ─────────────────────────────────────────────────────

const toTitle = (snake: string) =>
  snake.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase())

const copyVar = async (varName: string) => {
  try {
    await navigator.clipboard.writeText(varName)
    copied.value = varName
    setTimeout(() => {
      copied.value = null
    }, 1500)
  } catch {
    // clipboard unavailable (e.g. in some Electron sandboxes)
  }
}

const clearFilters = () => {
  search.value = ''
  sourceFilter.value = 'all'
}

onMounted(() => {
  activeSection.value = visibleSections.value[0]?.id ?? ''
})
</script>

<style scoped>
/* ── Root ───────────────────────────────────────────────────────── */
.docs-root {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 48px); /* 48px = compact v-app-bar */
  overflow: hidden;
}

/* ── Top bar ────────────────────────────────────────────────────── */
.docs-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  background: var(--color-card-header);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
  gap: 16px;
}

.docs-topbar-left {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.docs-title {
  font-size: 14.5px;
  font-weight: 600;
  color: white;
  letter-spacing: 0.01em;
}

.docs-subtitle {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
}

.docs-topbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.docs-search {
  width: 260px;
}

.docs-filter-group {
  margin-left: 8px;
}

/* ── Search field: white on dark topbar ─────────────────────────── */
.docs-topbar :deep(.v-field__outline) {
  --v-field-border-opacity: 1;
  color: rgba(255, 255, 255, 0.5) !important;
}
.docs-topbar :deep(.v-field:hover .v-field__outline),
.docs-topbar :deep(.v-field--focused .v-field__outline) {
  color: white !important;
}
.docs-topbar :deep(.v-field__input),
.docs-topbar :deep(.v-field__input::placeholder),
.docs-topbar :deep(.v-label) {
  color: white !important;
  opacity: 1 !important;
  caret-color: white !important;
}
.docs-topbar :deep(.v-field__prepend-inner .v-icon),
.docs-topbar :deep(.v-field__clearable .v-icon) {
  color: rgba(255, 255, 255, 0.6) !important;
}

/* ── Filter chips: white on dark topbar ─────────────────────────── */
.docs-topbar :deep(.v-chip) {
  color: rgba(255, 255, 255, 0.75) !important;
  border-color: rgba(255, 255, 255, 0.35) !important;
  background: transparent !important;
}
.docs-topbar :deep(.v-chip--variant-tonal) {
  background: rgba(255, 255, 255, 0.08) !important;
}
.docs-topbar :deep(.v-chip.v-chip--selected),
.docs-topbar :deep(.v-chip--selected) {
  background: white !important;
  color: var(--color-card-header) !important;
  border-color: white !important;
  font-weight: 600;
}
.docs-topbar :deep(.v-chip--selected .v-icon) {
  color: var(--color-card-header) !important;
}

/* ── Body ───────────────────────────────────────────────────────── */
.docs-body {
  display: flex;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

/* ── Left nav ───────────────────────────────────────────────────── */
.docs-nav {
  width: 200px;
  flex-shrink: 0;
  overflow-y: auto;
  background: #1c3545;
  border-right: 1px solid rgba(255, 255, 255, 0.07);
  padding: 10px 0;
}

.docs-nav-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 14px 8px 16px;
  cursor: pointer;
  border-left: 3px solid transparent;
  transition:
    background 0.12s,
    border-color 0.12s;
  gap: 8px;
  user-select: none;
}

.docs-nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.docs-nav-item--active {
  background: rgba(79, 195, 247, 0.08);
  border-left-color: #4fc3f7;
}

.docs-nav-item--active .docs-nav-label {
  color: #4fc3f7;
  font-weight: 600;
}

.docs-nav-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.7);
  flex: 1;
  line-height: 1.4;
}

.docs-nav-count {
  font-size: 10px;
  color: rgba(255, 255, 255, 0.3);
  background: rgba(255, 255, 255, 0.07);
  border-radius: 10px;
  padding: 1px 6px;
  flex-shrink: 0;
}

.docs-nav-empty {
  padding: 16px;
  font-size: 11.5px;
  color: rgba(255, 255, 255, 0.3);
  text-align: center;
}

/* ── Main content ───────────────────────────────────────────────── */
.docs-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px 28px;
  background: #f4f6f9;
}

/* ── Section ────────────────────────────────────────────────────── */
.docs-section {
  margin-bottom: 36px;
}

.docs-section-header {
  display: flex;
  align-items: baseline;
  gap: 10px;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 2px solid #003f58;
}

.docs-section-title {
  font-size: 15px;
  font-weight: 700;
  color: #003f58;
  letter-spacing: 0.01em;
}

.docs-section-count {
  font-size: 11px;
  color: #778;
}

/* ── Variable card ──────────────────────────────────────────────── */
.var-card {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 13px 16px;
  margin-bottom: 8px;
  transition: box-shadow 0.15s;
}

.var-card:hover {
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.07);
}

.var-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 4px;
}

.var-human-name {
  font-size: 13px;
  font-weight: 600;
  color: #1a2e3b;
  flex: 1;
}

.var-badges {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.var-identifier-row {
  display: flex;
  align-items: center;
  margin-bottom: 9px;
}

.var-identifier {
  font-family: 'Courier New', Courier, monospace;
  font-size: 11px;
  color: #556;
  background: #f1f5f9;
  padding: 2px 7px;
  border-radius: 4px;
  border: 1px solid #e2e8f0;
}

.var-description {
  font-size: 12.5px;
  color: #444;
  line-height: 1.65;
  margin: 0 0 10px 0;
}

/* ── Formula block ──────────────────────────────────────────────── */
.formula-block {
  background: #1e293b;
  border-radius: 6px;
  overflow: hidden;
}

.formula-header {
  display: flex;
  align-items: center;
  padding: 5px 12px;
  background: rgba(255, 255, 255, 0.04);
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.formula-pre {
  margin: 0;
  padding: 10px 14px 12px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 11.5px;
  color: #a8d8f0;
  line-height: 1.65;
  white-space: pre-wrap;
  word-break: break-word;
}

/* ── Source ref ─────────────────────────────────────────────────── */
.source-ref {
  display: inline-flex;
  align-items: center;
  font-size: 11.5px;
  color: #1565c0;
  background: #e3f2fd;
  border-radius: 4px;
  padding: 3px 8px;
}

/* ── Empty state ────────────────────────────────────────────────── */
.docs-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px 20px;
  text-align: center;
}

.docs-empty-title {
  font-size: 14.5px;
  font-weight: 500;
  color: #555;
  margin-bottom: 6px;
}

.docs-empty-sub {
  font-size: 12px;
  color: #888;
}
</style>
