<template>
  <v-container fluid class="user-flags pa-6">
    <v-row class="align-center mb-4" no-gutters>
      <v-col cols="auto" class="me-3">
        <v-btn
          variant="tonal"
          size="small"
          color="primary"
          prepend-icon="mdi-arrow-left"
          :to="{ name: 'group-pricing-quote-performance' }"
        >
          Back
        </v-btn>
      </v-col>
      <v-col>
        <h2 class="text-h5 font-weight-bold mb-0">User flags</h2>
        <span class="text-caption text-medium-emphasis">
          Coaching and capacity flags raised against users from the Quote
          Performance Dashboard.
        </span>
      </v-col>
    </v-row>

    <v-card elevation="1" class="rounded-lg mb-4 pa-2">
      <v-tabs v-model="activeTab" density="compact" color="primary">
        <v-tab value="open">
          Open
          <v-chip
            v-if="openCount > 0"
            size="x-small"
            color="warning"
            variant="tonal"
            class="ms-2"
            >{{ openCount }}</v-chip
          >
        </v-tab>
        <v-tab value="resolved">Resolved</v-tab>
        <v-tab value="all">All</v-tab>
      </v-tabs>
    </v-card>

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

    <v-card elevation="1" class="rounded-lg">
      <v-table density="compact" hover>
        <thead>
          <tr>
            <th>User</th>
            <th>Reason</th>
            <th>Opened by</th>
            <th>Opened at</th>
            <th>Note</th>
            <th v-if="activeTab !== 'open'">Resolved by</th>
            <th v-if="activeTab !== 'open'">Resolved at</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="f in flags" :key="f.id">
            <td>{{ f.user_name }}</td>
            <td>
              <v-chip
                size="x-small"
                :color="f.flag_reason === 'coaching' ? 'warning' : 'info'"
                variant="tonal"
              >
                {{ f.flag_reason === 'coaching' ? 'Coaching' : 'Capacity' }}
              </v-chip>
            </td>
            <td>{{ f.opened_by_name }}</td>
            <td class="text-no-wrap">{{ formatDate(f.opened_at) }}</td>
            <td class="note-cell">
              <span v-if="f.note">{{ truncate(f.note, 80) }}</span>
              <span v-else class="text-medium-emphasis">—</span>
              <v-tooltip
                v-if="f.note && f.note.length > 80"
                activator="parent"
                location="top"
                max-width="400"
                >{{ f.note }}</v-tooltip
              >
            </td>
            <td v-if="activeTab !== 'open'">
              {{ f.resolved_by_name || '—' }}
            </td>
            <td v-if="activeTab !== 'open'" class="text-no-wrap">
              {{ f.resolved_at ? formatDate(f.resolved_at) : '—' }}
            </td>
            <td class="text-right">
              <v-btn
                v-if="!f.resolved_at"
                variant="text"
                size="small"
                color="primary"
                prepend-icon="mdi-check-circle-outline"
                @click="openResolveDialog(f)"
              >
                Resolve
              </v-btn>
              <span v-else class="text-medium-emphasis text-caption"
                >Resolved</span
              >
            </td>
          </tr>
          <tr v-if="!flags.length && !loading">
            <td
              :colspan="activeTab === 'open' ? 6 : 8"
              class="text-center text-medium-emphasis py-6"
              >No flags match this view.</td
            >
          </tr>
        </tbody>
      </v-table>
      <v-progress-linear v-if="loading" indeterminate color="primary" />
    </v-card>

    <UserFlagDialog
      v-model="dialogOpen"
      mode="resolve"
      :flag-to-resolve="dialogFlag"
      :saving="saving"
      :error="error"
      @submit-resolve="handleResolve"
    />
  </v-container>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'

import type { QuoteUserFlag } from '@/renderer/api/QuoteDashboardService'
import { useUserFlags } from '@/renderer/composables/useQuoteDashboard'
import UserFlagDialog from '@/renderer/screens/group_pricing/dashboard/UserFlagDialog.vue'

const { flags, loading, saving, error, filter, refresh, resolveFlag } =
  useUserFlags()

const activeTab = ref<'open' | 'resolved' | 'all'>('open')
const dialogOpen = ref(false)
const dialogFlag = ref<QuoteUserFlag | null>(null)

const openCount = computed(
  () => flags.value.filter((f) => !f.resolved_at).length
)

watch(activeTab, (tab) => {
  filter.value.status = tab
  refresh()
})

onMounted(() => {
  filter.value.status = 'open'
  refresh()
})

function openResolveDialog(flag: QuoteUserFlag) {
  dialogFlag.value = flag
  dialogOpen.value = true
}

async function handleResolve(id: number, note: string) {
  try {
    await resolveFlag(id, note)
    dialogOpen.value = false
    await refresh()
  } catch {
    // error surfaces via the error ref on the dialog
  }
}

function formatDate(iso?: string | null): string {
  if (!iso) return ''
  try {
    return new Date(iso).toLocaleString()
  } catch {
    return iso
  }
}

function truncate(s: string, n: number): string {
  if (!s) return ''
  return s.length > n ? s.slice(0, n - 1) + '…' : s
}
</script>

<style scoped>
.user-flags {
  background: #f6f8fb;
  min-height: 100%;
}
.user-flags :deep(.v-card) {
  background: #ffffff;
  border: 1px solid #e2e8f0;
}
.user-flags :deep(.v-table thead th) {
  color: #64748b !important;
  font-size: 0.6875rem !important;
  font-weight: 600 !important;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  border-bottom: 1px solid #e2e8f0 !important;
  background: #fafbfd;
}
.note-cell {
  max-width: 360px;
}
</style>
