<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <div class="d-flex align-center">
                <v-btn
                  class="mr-3"
                  size="small"
                  variant="text"
                  prepend-icon="mdi-arrow-left"
                  @click="$router.back()"
                >
                  Back
                </v-btn>
                <span class="headline">Email Outbox</span>
              </div>
              <v-btn
                size="small"
                variant="outlined"
                rounded
                prepend-icon="mdi-refresh"
                :loading="loading"
                @click="load"
              >
                Refresh
              </v-btn>
            </div>
          </template>
          <template #default>
            <v-row class="mb-3" align="center">
              <v-col cols="12" md="4">
                <v-select
                  v-model="statusFilter"
                  :items="statusOptions"
                  label="Status"
                  variant="outlined"
                  density="compact"
                  hide-details
                  clearable
                  @update:model-value="
                    () => {
                      page = 1
                      load()
                    }
                  "
                />
              </v-col>
            </v-row>

            <v-data-table-server
              :headers="headers"
              :items="items"
              :items-length="total"
              :loading="loading"
              :items-per-page="pageSize"
              :page="page"
              density="compact"
              @update:options="onOptions"
            >
              <template #[`item.status`]="{ item }: { item: any }">
                <v-chip
                  size="x-small"
                  :color="statusColor(item.status)"
                  variant="tonal"
                >
                  {{ item.status }}
                </v-chip>
              </template>
              <template #[`item.to_recipients`]="{ item }: { item: any }">
                <span class="text-caption">{{
                  formatRecipients(item.to_recipients)
                }}</span>
              </template>
              <template #[`item.actions`]="{ item }: { item: any }">
                <v-btn
                  size="x-small"
                  variant="text"
                  icon="mdi-eye-outline"
                  @click="openDetail(item)"
                />
                <v-btn
                  v-if="item.status === 'failed' || item.status === 'pending'"
                  size="x-small"
                  variant="text"
                  icon="mdi-refresh"
                  :loading="retryingId === item.id"
                  @click="retry(item)"
                />
              </template>
            </v-data-table-server>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-dialog v-model="detailDialog" max-width="800">
      <v-card v-if="detail">
        <v-card-title>
          Outbox #{{ detail.id }}
          <v-chip
            class="ml-2"
            size="small"
            :color="statusColor(detail.status)"
            variant="tonal"
          >
            {{ detail.status }}
          </v-chip>
        </v-card-title>
        <v-card-text>
          <p><strong>Template:</strong> {{ detail.template_code }}</p>
          <p
            ><strong>To:</strong>
            {{ formatRecipients(detail.to_recipients) }}</p
          >
          <p v-if="formatRecipients(detail.cc_recipients)">
            <strong>Cc:</strong> {{ formatRecipients(detail.cc_recipients) }}
          </p>
          <p><strong>Subject:</strong> {{ detail.subject }}</p>
          <p>
            <strong>Attempts:</strong> {{ detail.attempts }} /
            {{ detail.max_attempts }}
          </p>
          <p v-if="detail.last_error" class="text-error">
            <strong>Last error:</strong> {{ detail.last_error }}
          </p>
          <v-divider class="my-3" />
          <div v-html="detail.body" />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn variant="text" @click="detailDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useFlashStore } from '@/renderer/store/flash'

const flash = useFlashStore()

interface OutboxRow {
  id: number
  template_code: string
  to_recipients: string
  cc_recipients: string
  subject: string
  status: string
  attempts: number
  max_attempts: number
  last_error: string
  created_at: string
  sent_at: string | null
  body: string
}

const items = ref<OutboxRow[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(25)
const loading = ref(false)
const statusFilter = ref('')
const statusOptions = [
  { title: 'All', value: '' },
  { title: 'Pending', value: 'pending' },
  { title: 'Sending', value: 'sending' },
  { title: 'Sent', value: 'sent' },
  { title: 'Failed', value: 'failed' }
]

const detailDialog = ref(false)
const detail = ref<OutboxRow | null>(null)
const retryingId = ref<number | null>(null)

const headers = [
  { title: 'ID', key: 'id', width: '80px' },
  { title: 'Template', key: 'template_code' },
  { title: 'To', key: 'to_recipients' },
  { title: 'Subject', key: 'subject' },
  { title: 'Status', key: 'status', width: '120px' },
  { title: 'Attempts', key: 'attempts', width: '100px' },
  { title: 'Created', key: 'created_at', width: '180px' },
  { title: '', key: 'actions', sortable: false, width: '100px' }
]

const statusColor = (s: string) => {
  switch (s) {
    case 'sent':
      return 'success'
    case 'sending':
      return 'info'
    case 'failed':
      return 'error'
    default:
      return 'warning'
  }
}

const formatRecipients = (json: string): string => {
  if (!json) return ''
  try {
    const arr = JSON.parse(json)
    return Array.isArray(arr) ? arr.join(', ') : json
  } catch {
    return json
  }
}

const load = async () => {
  loading.value = true
  try {
    const { data } = await GroupPricingService.listEmailOutbox({
      status: statusFilter.value || undefined,
      page: page.value,
      page_size: pageSize.value
    })
    items.value = data.items || []
    total.value = data.total || 0
  } catch (err: any) {
    flash.show(err?.response?.data?.error || 'Failed to load outbox', 'error')
  } finally {
    loading.value = false
  }
}

const onOptions = (opts: { page: number; itemsPerPage: number }) => {
  page.value = opts.page
  pageSize.value = opts.itemsPerPage
  load()
}

const openDetail = async (row: OutboxRow) => {
  try {
    const { data } = await GroupPricingService.getEmailOutboxItem(row.id)
    detail.value = data
    detailDialog.value = true
  } catch (err: any) {
    flash.show(err?.response?.data?.error || 'Failed to load row', 'error')
  }
}

const retry = async (row: OutboxRow) => {
  retryingId.value = row.id
  try {
    await GroupPricingService.retryEmailOutbox(row.id)
    flash.show(`Requeued #${row.id}`, 'success')
    await load()
  } catch (err: any) {
    flash.show(err?.response?.data?.error || 'Retry failed', 'error')
  } finally {
    retryingId.value = null
  }
}

onMounted(load)
</script>
