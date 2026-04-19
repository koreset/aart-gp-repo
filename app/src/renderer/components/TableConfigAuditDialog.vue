<template>
  <v-dialog v-model="dialog" max-width="720px">
    <base-card>
      <template #header>
        <span class="headline">{{ tableType }} — Configuration History</span>
      </template>
      <template #default>
        <v-row v-if="latest">
          <v-col>
            <p class="text-body-2 mb-1">
              <strong>Last updated by:</strong>
              {{ latest.changed_by || '—' }}
            </p>
            <p class="text-body-2 mb-1">
              <strong>Last updated on:</strong>
              {{ formatDate(latest.changed_at) }}
            </p>
            <p class="text-body-2 mb-3">
              <strong>Currently:</strong>
              <span :class="latest.new_value ? 'text-success' : 'text-warning'">
                {{ latest.new_value ? 'Required' : 'Not required' }}
              </span>
            </p>
          </v-col>
        </v-row>
        <v-row v-if="loading">
          <v-col>
            <v-progress-linear indeterminate color="primary" />
          </v-col>
        </v-row>
        <v-row v-else-if="entries.length === 0">
          <v-col>
            <p class="text-body-2 text-medium-emphasis">
              No configuration history available for this table.
            </p>
          </v-col>
        </v-row>
        <v-row v-else>
          <v-col>
            <v-data-table
              :headers="headers"
              :items="entries"
              :items-per-page="10"
              density="compact"
              hover
            >
              <template #[`item.changed_at`]="{ item: row }">
                {{ formatDate((row as any).changed_at) }}
              </template>
              <template #[`item.old_value`]="{ item: row }">
                <v-chip
                  size="x-small"
                  :color="(row as any).old_value ? 'success' : 'grey'"
                  variant="tonal"
                  label
                >
                  {{ (row as any).old_value ? 'Required' : 'Not required' }}
                </v-chip>
              </template>
              <template #[`item.new_value`]="{ item: row }">
                <v-chip
                  size="x-small"
                  :color="(row as any).new_value ? 'success' : 'warning'"
                  variant="tonal"
                  label
                >
                  {{ (row as any).new_value ? 'Required' : 'Not required' }}
                </v-chip>
              </template>
            </v-data-table>
          </v-col>
        </v-row>
      </template>
      <template #actions>
        <v-spacer />
        <v-btn rounded variant="text" @click="dialog = false">Close</v-btn>
      </template>
    </base-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface AuditEntry {
  id: number
  table_type: string
  event_type: string
  old_value: boolean
  new_value: boolean
  changed_by: string
  changed_at: string
  details: string
}

const dialog = ref(false)
const loading = ref(false)
const entries = ref<AuditEntry[]>([])
const tableType = ref('')

const headers = [
  { title: 'When', key: 'changed_at', width: '24%' },
  { title: 'By', key: 'changed_by', width: '24%' },
  { title: 'From', key: 'old_value', width: '14%' },
  { title: 'To', key: 'new_value', width: '14%' },
  { title: 'Note', key: 'details' }
]

const latest = computed(() => entries.value[0] ?? null)

const formatDate = (s: string) => {
  if (!s) return '—'
  try {
    return new Date(s).toLocaleString()
  } catch {
    return s
  }
}

const open = async (key: string, displayName?: string) => {
  tableType.value = displayName || key
  dialog.value = true
  loading.value = true
  entries.value = []
  try {
    const res = await GroupPricingService.getTableConfigurationAudit(key)
    entries.value = res.data?.data || []
  } catch (e) {
    console.error('Failed to load table configuration audit:', e)
  } finally {
    loading.value = false
  }
}

defineExpose({ open })
</script>
