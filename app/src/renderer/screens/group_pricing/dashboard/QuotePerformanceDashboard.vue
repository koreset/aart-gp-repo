<template>
  <v-container fluid class="quote-performance pa-6">
    <!-- Header -->
    <v-row class="align-center mb-4" no-gutters>
      <v-col>
        <h2 class="text-h5 font-weight-bold mb-0">Quote Performance</h2>
        <span class="text-caption text-medium-emphasis">
          Per-user productivity, SLA compliance, conversion and pipeline value.
        </span>
      </v-col>
      <v-col cols="auto">
        <v-btn
          color="primary"
          variant="tonal"
          prepend-icon="mdi-file-table-outline"
          :to="{ name: 'group-pricing-quote-extract' }"
          class="me-2"
        >
          Open extract
        </v-btn>
        <v-btn
          v-if="canManageSla"
          color="primary"
          variant="tonal"
          prepend-icon="mdi-timer-cog-outline"
          :to="{ name: 'group-pricing-sla-targets' }"
        >
          SLA targets
        </v-btn>
      </v-col>
    </v-row>

    <!-- Filters -->
    <v-card elevation="1" class="rounded-lg pa-4 mb-4">
      <v-row dense align="end">
        <v-col cols="12" sm="6" md="3">
          <v-text-field
            v-model="filters.from"
            type="date"
            label="From"
            density="compact"
            hide-details
            clearable
            @update:model-value="onFilterChange"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-text-field
            v-model="filters.to"
            type="date"
            label="To"
            density="compact"
            hide-details
            clearable
            @update:model-value="onFilterChange"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-combobox
            v-model="filters.users"
            label="Users"
            density="compact"
            chips
            closable-chips
            multiple
            hide-details
            clearable
            @update:model-value="onFilterChange"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-combobox
            v-model="filters.region"
            label="Regions"
            density="compact"
            chips
            closable-chips
            multiple
            hide-details
            clearable
            @update:model-value="onFilterChange"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-select
            v-model="filters.quote_type"
            :items="['New Business', 'Renewal']"
            label="Quote type"
            density="compact"
            chips
            closable-chips
            multiple
            hide-details
            clearable
            @update:model-value="onFilterChange"
          />
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-select
            v-model="filters.distribution_channel"
            :items="['broker', 'direct', 'binder', 'tied_agent']"
            label="Channel"
            density="compact"
            chips
            closable-chips
            multiple
            hide-details
            clearable
            @update:model-value="onFilterChange"
          />
        </v-col>
        <v-col cols="auto" class="ms-auto">
          <v-btn
            variant="text"
            density="compact"
            size="small"
            @click="resetAndRefresh"
          >
            Reset filters
          </v-btn>
          <v-btn
            color="primary"
            variant="tonal"
            density="compact"
            size="small"
            prepend-icon="mdi-refresh"
            :loading="loading"
            @click="refreshAll"
          >
            Refresh
          </v-btn>
        </v-col>
      </v-row>
    </v-card>

    <!-- Error banner -->
    <v-alert
      v-if="error"
      type="error"
      density="compact"
      class="mb-4"
      closable
      @click:close="error = null"
    >
      {{ error }}
    </v-alert>

    <!-- Backfill caveat -->
    <div class="backfill-note mb-4">
      <v-icon size="small" class="me-2">mdi-information-outline</v-icon>
      <span>
        Cycle-time and SLA metrics for quotes created before this feature was
        installed are estimates derived from existing timestamps.
      </span>
    </div>

    <!-- KPI tiles -->
    <v-row dense class="mb-4">
      <v-col v-for="card in summaryCards" :key="card.label" cols="6" md="3">
        <div class="metric-tile" :class="`metric-tile--${card.color}`">
          <div class="metric-tile__head">
            <v-icon size="18" class="metric-tile__icon">
              {{ card.icon }}
            </v-icon>
            <span class="metric-tile__label">{{ card.label }}</span>
          </div>
          <div class="metric-tile__value">{{ card.value }}</div>
          <div v-if="card.sub" class="metric-tile__sub">{{ card.sub }}</div>
        </div>
      </v-col>
    </v-row>

    <!-- Charts -->
    <v-row class="mb-4">
      <v-col cols="12" md="6">
        <v-card elevation="1" class="rounded-lg pa-4">
          <div class="d-flex align-center mb-3">
            <v-icon color="primary" class="me-2">mdi-filter-variant</v-icon>
            <div class="text-subtitle-1 font-weight-bold">Quote funnel</div>
          </div>
          <ag-charts :options="funnelChartOptions" :style="{ height: '280px' }" />
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card elevation="1" class="rounded-lg pa-4">
          <div class="d-flex align-center mb-3">
            <v-icon color="primary" class="me-2">mdi-chart-line</v-icon>
            <div class="text-subtitle-1 font-weight-bold">Trend</div>
            <v-spacer />
            <v-btn-toggle
              v-model="trendBucket"
              density="compact"
              mandatory
              variant="outlined"
              divided
              @update:model-value="refreshAll"
            >
              <v-btn value="day" size="x-small">Day</v-btn>
              <v-btn value="week" size="x-small">Week</v-btn>
              <v-btn value="month" size="x-small">Month</v-btn>
            </v-btn-toggle>
          </div>
          <ag-charts :options="trendChartOptions" :style="{ height: '280px' }" />
        </v-card>
      </v-col>
    </v-row>

    <!-- SLA breaches -->
    <v-row class="mb-4">
      <v-col cols="12" md="6">
        <v-card elevation="1" class="rounded-lg pa-4">
          <div class="d-flex align-center mb-3">
            <v-icon color="warning" class="me-2">mdi-timer-alert-outline</v-icon>
            <div class="text-subtitle-1 font-weight-bold">
              SLA breaches by transition
            </div>
          </div>
          <v-table density="compact" hover>
            <thead>
              <tr>
                <th>Transition</th>
                <th class="text-right">Target (hrs)</th>
                <th class="text-right">Transitions</th>
                <th class="text-right">Breaches</th>
                <th class="text-right">Compliance</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="row in sla.breaches_by_transition"
                :key="row.from_status + '|' + row.to_status"
              >
                <td
                  >{{ row.from_status }} → {{ row.to_status }}</td
                >
                <td class="text-right">{{ row.target_hours }}</td>
                <td class="text-right">{{ row.transition_count }}</td>
                <td class="text-right">
                  <v-chip
                    size="x-small"
                    :color="row.breach_count > 0 ? 'warning' : 'success'"
                    variant="tonal"
                    >{{ row.breach_count }}</v-chip
                  >
                </td>
                <td class="text-right">
                  {{ compliancePct(row.breach_count, row.transition_count) }}
                </td>
              </tr>
              <tr v-if="!sla.breaches_by_transition.length">
                <td colspan="5" class="text-center text-medium-emphasis py-4"
                  >No transitions recorded for this filter.</td
                >
              </tr>
            </tbody>
          </v-table>
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card elevation="1" class="rounded-lg pa-4">
          <div class="d-flex align-center mb-3">
            <v-icon color="primary" class="me-2"
              >mdi-account-multiple-outline</v-icon
            >
            <div class="text-subtitle-1 font-weight-bold">
              Top users by SLA breaches
            </div>
          </div>
          <v-table density="compact" hover>
            <thead>
              <tr>
                <th>User</th>
                <th class="text-right">Transitions</th>
                <th class="text-right">Breaches</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in sla.breaches_by_user" :key="row.user_name">
                <td>{{ row.user_name }}</td>
                <td class="text-right">{{ row.transition_count }}</td>
                <td class="text-right">
                  <v-chip
                    size="x-small"
                    :color="row.breach_count > 0 ? 'warning' : 'success'"
                    variant="tonal"
                    >{{ row.breach_count }}</v-chip
                  >
                </td>
              </tr>
              <tr v-if="!sla.breaches_by_user.length">
                <td colspan="3" class="text-center text-medium-emphasis py-4"
                  >No transitions recorded for this filter.</td
                >
              </tr>
            </tbody>
          </v-table>
        </v-card>
      </v-col>
    </v-row>

    <!-- User leaderboard -->
    <v-card elevation="1" class="rounded-lg pa-4">
      <div class="d-flex align-center mb-3">
        <v-icon color="primary" class="me-2">mdi-trophy-outline</v-icon>
        <div class="text-subtitle-1 font-weight-bold">User leaderboard</div>
        <v-spacer />
        <span class="text-caption text-medium-emphasis"
          >{{ kpis.length }} users in scope</span
        >
      </div>
      <v-table density="compact" hover class="leaderboard">
        <thead>
          <tr>
            <th>User</th>
            <th class="text-right">Total</th>
            <th class="text-right">Submitted</th>
            <th class="text-right">Approved</th>
            <th class="text-right">Accepted</th>
            <th class="text-right">In force</th>
            <th class="text-right">Approval %</th>
            <th class="text-right">Conv. %</th>
            <th class="text-right">Avg cycle (hrs)</th>
            <th class="text-right">SLA %</th>
            <th class="text-right">Pipeline premium</th>
            <th class="text-right">Won premium</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in kpis" :key="row.user_name">
            <td class="font-weight-medium">{{ row.user_name }}</td>
            <td class="text-right">{{ row.total_quotes }}</td>
            <td class="text-right">{{ row.submitted_count }}</td>
            <td class="text-right">{{ row.approved_count }}</td>
            <td class="text-right">{{ row.accepted_count }}</td>
            <td class="text-right">{{ row.in_force_count }}</td>
            <td class="text-right">{{ pct(row.approval_rate) }}</td>
            <td class="text-right">{{ pct(row.conversion_rate) }}</td>
            <td class="text-right">{{ hours(row.avg_total_cycle_hours) }}</td>
            <td class="text-right">
              <v-chip
                v-if="row.sla_transition_count > 0"
                size="x-small"
                :color="slaColor(row)"
                variant="tonal"
                >{{ pct(row.sla_compliance_pct) }}</v-chip
              >
              <span v-else class="text-medium-emphasis">—</span>
            </td>
            <td class="text-right"
              >{{ formatCurrency(row.pipeline_annual_premium) }}</td
            >
            <td class="text-right"
              >{{ formatCurrency(row.total_annual_premium) }}</td
            >
          </tr>
          <tr v-if="!kpis.length">
            <td colspan="12" class="text-center text-medium-emphasis py-6"
              >No quotes match the active filter.</td
            >
          </tr>
        </tbody>
      </v-table>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { AgCharts } from 'ag-charts-vue3'

import {
  useQuoteDashboard,
  type QuotePerformanceKpis
} from '@/renderer/composables/useQuoteDashboard'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

const {
  filters,
  resetFilters,
  kpis,
  funnel,
  trend,
  sla,
  loading,
  error,
  trendBucket,
  refreshAll
} = useQuoteDashboard()

const { hasPermission } = usePermissionCheck()
const canManageSla = computed(() => hasPermission('quote:manage_sla_targets'))

onMounted(() => {
  refreshAll()
})

function onFilterChange() {
  // Debounce-ish: refresh on next tick so combobox commits.
  setTimeout(refreshAll, 0)
}

function resetAndRefresh() {
  resetFilters()
  refreshAll()
}

const totals = computed(() => {
  const init = {
    quotes: 0,
    submitted: 0,
    approved: 0,
    accepted: 0,
    pipeline: 0,
    won: 0,
    breaches: 0,
    transitions: 0,
    cycle_sum: 0,
    cycle_count: 0
  }
  for (const row of kpis.value) {
    init.quotes += row.total_quotes
    init.submitted += row.submitted_count
    init.approved += row.approved_count
    init.accepted += row.accepted_count
    init.pipeline += row.pipeline_annual_premium
    init.won += row.total_annual_premium
    init.breaches += row.sla_breach_count
    init.transitions += row.sla_transition_count
    if (row.avg_total_cycle_hours > 0) {
      init.cycle_sum += row.avg_total_cycle_hours * row.accepted_count
      init.cycle_count += row.accepted_count
    }
  }
  return init
})

const summaryCards = computed(() => {
  const t = totals.value
  const approvalRate = t.submitted > 0 ? t.approved / t.submitted : 0
  const conversionRate = t.submitted > 0 ? t.accepted / t.submitted : 0
  const avgCycle = t.cycle_count > 0 ? t.cycle_sum / t.cycle_count : 0
  const slaPct = t.transitions > 0 ? 1 - t.breaches / t.transitions : 0
  return [
    {
      icon: 'mdi-file-document-multiple-outline',
      color: 'primary',
      label: 'Total quotes',
      value: t.quotes.toLocaleString(),
      sub: `${t.submitted.toLocaleString()} submitted`
    },
    {
      icon: 'mdi-check-decagram-outline',
      color: 'success',
      label: 'Approval rate',
      value: pct(approvalRate),
      sub: `${t.approved.toLocaleString()} approved`
    },
    {
      icon: 'mdi-handshake-outline',
      color: 'primary',
      label: 'Conversion rate',
      value: pct(conversionRate),
      sub: `${t.accepted.toLocaleString()} accepted`
    },
    {
      icon: 'mdi-timer-outline',
      color: 'info',
      label: 'Avg cycle time',
      value: hours(avgCycle),
      sub: 'submit → accept'
    },
    {
      icon: 'mdi-cash-multiple',
      color: 'primary',
      label: 'Won annual premium',
      value: formatCurrency(t.won),
      sub: 'accepted + in_force'
    },
    {
      icon: 'mdi-cash-clock',
      color: 'warning',
      label: 'Pipeline premium',
      value: formatCurrency(t.pipeline),
      sub: 'submitted + approved'
    },
    {
      icon: 'mdi-timer-alert-outline',
      color: t.breaches > 0 ? 'warning' : 'success',
      label: 'SLA compliance',
      value: t.transitions > 0 ? pct(slaPct) : '—',
      sub: `${t.breaches.toLocaleString()} breaches`
    },
    {
      icon: 'mdi-account-multiple-outline',
      color: 'primary',
      label: 'Active users',
      value: kpis.value.length.toLocaleString(),
      sub: 'with at least one quote'
    }
  ]
})

// Series colours align with the .metric-tile palette so the dashboard
// reads as one cohesive design rather than ag-charts' default rainbow.
const PALETTE = {
  primary: '#2563eb',
  success: '#059669',
  info: '#0891b2',
  warning: '#d97706',
  danger: '#dc2626',
  text: '#0f172a',
  textMuted: '#64748b'
}

const funnelChartOptions = computed<any>(() => ({
  background: { fill: 'transparent' },
  data: funnel.value.map((s) => ({
    stage: s.stage,
    count: s.count,
    dwell: s.avg_dwell_hours
  })),
  series: [
    {
      type: 'bar',
      direction: 'horizontal',
      xKey: 'stage',
      yKey: 'count',
      yName: 'Quotes',
      fill: PALETTE.primary,
      cornerRadius: 4,
      tooltip: {
        renderer: ({ datum }: any) =>
          `<b>${datum.stage}</b><br>${datum.count.toLocaleString()} quotes${
            datum.dwell > 0
              ? `<br>Avg dwell: ${datum.dwell.toFixed(1)} hrs`
              : ''
          }`
      }
    }
  ],
  axes: [
    {
      type: 'category',
      position: 'left',
      label: { color: PALETTE.textMuted, fontSize: 12 },
      line: { enabled: false },
      tick: { enabled: false }
    },
    {
      type: 'number',
      position: 'bottom',
      label: { color: PALETTE.textMuted, fontSize: 11 },
      gridLine: { style: [{ stroke: '#e2e8f0' }] }
    }
  ],
  legend: { enabled: false }
}))

const trendChartOptions = computed<any>(() => ({
  background: { fill: 'transparent' },
  data: trend.value,
  series: [
    {
      type: 'line',
      xKey: 'bucket',
      yKey: 'submitted',
      yName: 'Submitted',
      stroke: PALETTE.primary,
      marker: { enabled: true, fill: PALETTE.primary, stroke: PALETTE.primary }
    },
    {
      type: 'line',
      xKey: 'bucket',
      yKey: 'approved',
      yName: 'Approved',
      stroke: PALETTE.success,
      marker: { enabled: true, fill: PALETTE.success, stroke: PALETTE.success }
    },
    {
      type: 'line',
      xKey: 'bucket',
      yKey: 'accepted',
      yName: 'Accepted',
      stroke: PALETTE.info,
      marker: { enabled: true, fill: PALETTE.info, stroke: PALETTE.info }
    },
    {
      type: 'line',
      xKey: 'bucket',
      yKey: 'rejected',
      yName: 'Rejected',
      stroke: PALETTE.danger,
      marker: { enabled: true, fill: PALETTE.danger, stroke: PALETTE.danger }
    }
  ],
  axes: [
    {
      type: 'category',
      position: 'bottom',
      label: { color: PALETTE.textMuted, fontSize: 11 },
      line: { stroke: '#e2e8f0' },
      tick: { stroke: '#e2e8f0' }
    },
    {
      type: 'number',
      position: 'left',
      label: { color: PALETTE.textMuted, fontSize: 11 },
      gridLine: { style: [{ stroke: '#e2e8f0' }] }
    }
  ],
  legend: {
    position: 'bottom',
    item: { label: { color: PALETTE.textMuted, fontSize: 12 } }
  }
}))

function pct(value: number): string {
  if (!isFinite(value)) return '—'
  return (value * 100).toFixed(1) + '%'
}

function hours(value: number): string {
  if (!value || !isFinite(value)) return '—'
  if (value >= 48) return (value / 24).toFixed(1) + ' d'
  return value.toFixed(1)
}

function formatCurrency(value: number): string {
  if (!value || !isFinite(value)) return 'R0'
  if (value >= 1_000_000) return 'R' + (value / 1_000_000).toFixed(2) + 'M'
  if (value >= 1_000) return 'R' + (value / 1_000).toFixed(1) + 'K'
  return 'R' + Math.round(value).toLocaleString()
}

function compliancePct(breach: number, total: number): string {
  if (total === 0) return '—'
  return pct(1 - breach / total)
}

function slaColor(row: QuotePerformanceKpis): string {
  if (row.sla_compliance_pct >= 0.95) return 'success'
  if (row.sla_compliance_pct >= 0.8) return 'warning'
  return 'error'
}
</script>

<style scoped>
/*
 * Local palette — kept inside this component so it doesn't fight whatever
 * theme the rest of the app applies. Slate-ish neutrals + a restrained
 * accent set, tuned for a management dashboard where the data should be
 * the loud thing on screen, not the chrome.
 */
.quote-performance {
  --qp-bg: #f6f8fb;
  --qp-surface: #ffffff;
  --qp-border: #e2e8f0;
  --qp-text: #0f172a;
  --qp-text-muted: #64748b;
  --qp-text-subtle: #94a3b8;
  --qp-primary: #2563eb;
  --qp-success: #059669;
  --qp-warning: #d97706;
  --qp-danger: #dc2626;
  --qp-info: #0891b2;

  background: var(--qp-bg);
  min-height: 100%;
  color: var(--qp-text);
}

/* Section title (h2) */
.quote-performance :deep(.text-h5) {
  color: var(--qp-text);
  letter-spacing: -0.01em;
}

/* Every v-card on this page gets the same clean surface treatment. */
.quote-performance :deep(.v-card) {
  background: var(--qp-surface);
  border: 1px solid var(--qp-border);
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
}

/* Subtitle inside cards */
.quote-performance :deep(.text-subtitle-1) {
  color: var(--qp-text);
  font-weight: 600;
}

/* Backfill caveat — much lighter than a full v-alert */
.backfill-note {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 8px 12px;
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 8px;
  color: #075985;
  font-size: 0.8125rem;
}
.backfill-note .v-icon {
  color: var(--qp-info);
}

/* ──────────── KPI tiles ──────────── */
.metric-tile {
  position: relative;
  background: var(--qp-surface);
  border: 1px solid var(--qp-border);
  border-radius: 10px;
  padding: 14px 16px 12px;
  overflow: hidden;
  box-shadow: 0 1px 2px rgba(15, 23, 42, 0.04);
  transition: box-shadow 120ms ease, transform 120ms ease;
}
.metric-tile::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--tile-accent, var(--qp-primary));
}
.metric-tile:hover {
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.06);
}
.metric-tile__head {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
}
.metric-tile__icon {
  color: var(--tile-accent, var(--qp-primary));
}
.metric-tile__label {
  color: var(--qp-text-muted);
  font-size: 0.6875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}
.metric-tile__value {
  color: var(--qp-text);
  font-size: 1.5rem;
  font-weight: 700;
  line-height: 1.2;
  font-variant-numeric: tabular-nums;
}
.metric-tile__sub {
  color: var(--qp-text-subtle);
  font-size: 0.75rem;
  margin-top: 2px;
}

/* Accent palette — keyed off the card.color string from the summaryCards
 * computed in script. */
.metric-tile--primary { --tile-accent: var(--qp-primary); }
.metric-tile--success { --tile-accent: var(--qp-success); }
.metric-tile--warning { --tile-accent: var(--qp-warning); }
.metric-tile--danger  { --tile-accent: var(--qp-danger);  }
.metric-tile--info    { --tile-accent: var(--qp-info);    }

/* ──────────── Tables ──────────── */
.quote-performance :deep(.v-table) {
  background: transparent;
  --v-theme-surface: transparent;
}
.quote-performance :deep(.v-table thead th) {
  color: var(--qp-text-muted) !important;
  font-size: 0.6875rem !important;
  font-weight: 600 !important;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  border-bottom: 1px solid var(--qp-border) !important;
  background: #fafbfd;
}
.quote-performance :deep(.v-table tbody td) {
  color: var(--qp-text);
  font-size: 0.8125rem;
  border-bottom: 1px solid #f1f5f9 !important;
}
.quote-performance :deep(.v-table tbody tr:hover) {
  background: #f8fafc !important;
}

.leaderboard td,
.leaderboard th {
  white-space: nowrap;
}
.leaderboard :deep(td:not(:first-child)),
.leaderboard :deep(th:not(:first-child)) {
  font-variant-numeric: tabular-nums;
}

/* Restrained chip styling — Vuetify defaults are too saturated. */
.quote-performance :deep(.v-chip.v-chip--variant-tonal) {
  font-weight: 600;
  font-size: 0.6875rem;
  letter-spacing: 0.02em;
}

/* Filter card — slightly less elevation than data cards so the eye
 * lands on the data first. */
.quote-performance > .v-card:first-of-type {
  box-shadow: none;
}

/* Tone down v-btn outlined/tonal so they don't compete with data. */
.quote-performance :deep(.v-btn--variant-tonal) {
  font-weight: 600;
  letter-spacing: 0;
}
</style>
