<template>
  <v-dialog v-model="dialog" max-width="1400" persistent scrollable>
    <v-card>
      <v-card-title
        class="d-flex justify-space-between align-center bg-primary text-white"
      >
        <div class="d-flex align-center">
          <v-icon class="mr-2">mdi-history</v-icon>
          <span>Member History - {{ memberName }}</span>
        </div>
        <v-btn icon="mdi-close" variant="text" @click="closeDialog" />
      </v-card-title>

      <v-card-text class="pa-0">
        <v-container v-if="loading" class="text-center py-8">
          <v-progress-circular indeterminate color="primary" size="64" />
          <div class="mt-4 text-h6">Loading member history...</div>
        </v-container>

        <v-container v-else class="py-4">
          <!-- Filter Controls -->
          <v-row class="mb-4">
            <v-col cols="12" md="4">
              <v-select
                v-model="selectedEventType"
                :items="eventTypeOptions"
                label="Filter by Event Type"
                prepend-inner-icon="mdi-filter"
                clearable
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="4">
              <v-text-field
                v-model="dateRangeStart"
                type="date"
                label="From Date"
                prepend-inner-icon="mdi-calendar-start"
                variant="outlined"
                density="compact"
              />
            </v-col>
            <v-col cols="12" md="4">
              <v-text-field
                v-model="dateRangeEnd"
                type="date"
                label="To Date"
                prepend-inner-icon="mdi-calendar-end"
                variant="outlined"
                density="compact"
              />
            </v-col>
          </v-row>

          <!-- Statistics Cards -->
          <v-row class="mb-4">
            <v-col cols="6" md="3">
              <v-card variant="tonal" color="info">
                <v-card-text class="text-center">
                  <div class="text-h5 font-weight-bold">{{
                    historyStats.totalEvents
                  }}</div>
                  <div class="text-caption">Total Events</div>
                </v-card-text>
              </v-card>
            </v-col>
            <v-col cols="6" md="3">
              <v-card variant="tonal" color="success">
                <v-card-text class="text-center">
                  <div class="text-h5 font-weight-bold">{{
                    historyStats.policyChanges
                  }}</div>
                  <div class="text-caption">Policy Changes</div>
                </v-card-text>
              </v-card>
            </v-col>
            <v-col cols="6" md="3">
              <v-card variant="tonal" color="warning">
                <v-card-text class="text-center">
                  <div class="text-h5 font-weight-bold">{{
                    historyStats.claims
                  }}</div>
                  <div class="text-caption">Claims</div>
                </v-card-text>
              </v-card>
            </v-col>
            <v-col cols="6" md="3">
              <v-card variant="tonal" color="error">
                <v-card-text class="text-center">
                  <div class="text-h5 font-weight-bold">{{
                    historyStats.statusChanges
                  }}</div>
                  <div class="text-caption">Status Changes</div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>

          <!-- Timeline -->
          <v-row>
            <v-col>
              <v-card variant="outlined">
                <v-card-title class="bg-secondary text-white">
                  <v-icon class="mr-2">mdi-timeline</v-icon>
                  Event Timeline
                </v-card-title>
                <v-card-text class="pa-0">
                  <div class="timeline-container">
                    <v-timeline side="end" density="compact">
                      <v-timeline-item
                        v-for="(event, index) in filteredHistory"
                        :key="index"
                        :dot-color="getEventColor(event.type)"
                        size="small"
                        class="timeline-item"
                      >
                        <template #icon>
                          <v-icon size="small">{{
                            getEventIcon(event.type)
                          }}</v-icon>
                        </template>
                        <v-card
                          :color="getEventColor(event.type)"
                          variant="tonal"
                          class="mb-2"
                        >
                          <v-card-text class="pa-3">
                            <div
                              class="d-flex justify-space-between align-center mb-2"
                            >
                              <div class="text-subtitle-2 font-weight-bold">
                                {{ event.title }}
                              </div>
                              <div class="d-flex align-center gap-2">
                                <div class="text-caption text-grey">
                                  {{ formatDateTime(event.timestamp) }}
                                </div>
                                <v-chip
                                  :color="getEventColor(event.type)"
                                  size="small"
                                >
                                  {{ formatEventType(event.type) }}
                                </v-chip>
                              </div>
                            </div>
                            <div class="text-body-2 mb-2">{{
                              event.description
                            }}</div>

                            <!-- Event-specific details -->
                            <div v-if="event.details" class="event-details">
                              <!-- Salary Change Details -->
                              <div
                                v-if="event.type === 'salary_change'"
                                class="details-grid"
                              >
                                <div class="detail-item">
                                  <span class="detail-label">Previous:</span>
                                  <span class="detail-value">{{
                                    formatCurrency(event.details.previousValue)
                                  }}</span>
                                </div>
                                <div class="detail-item">
                                  <span class="detail-label">New:</span>
                                  <span class="detail-value">{{
                                    formatCurrency(event.details.newValue)
                                  }}</span>
                                </div>
                                <div class="detail-item">
                                  <span class="detail-label">Change:</span>
                                  <span
                                    :class="
                                      event.details.newValue >
                                      event.details.previousValue
                                        ? 'text-success'
                                        : 'text-error'
                                    "
                                    class="detail-value"
                                  >
                                    {{
                                      formatCurrency(
                                        Math.abs(
                                          event.details.newValue -
                                            event.details.previousValue
                                        )
                                      )
                                    }}
                                    {{
                                      event.details.newValue >
                                      event.details.previousValue
                                        ? '↑'
                                        : '↓'
                                    }}
                                  </span>
                                </div>
                              </div>

                              <!-- Benefit Change Details -->
                              <div
                                v-else-if="event.type === 'benefit_change'"
                                class="details-grid"
                              >
                                <div class="detail-item">
                                  <span class="detail-label">Benefit:</span>
                                  <span class="detail-value">{{
                                    event.details.benefitName
                                  }}</span>
                                </div>
                                <div class="detail-item">
                                  <span class="detail-label">Action:</span>
                                  <span class="detail-value">{{
                                    event.details.action
                                  }}</span>
                                </div>
                                <div
                                  v-if="event.details.amount"
                                  class="detail-item"
                                >
                                  <span class="detail-label">Amount:</span>
                                  <span class="detail-value">{{
                                    formatCurrency(event.details.amount)
                                  }}</span>
                                </div>
                              </div>

                              <!-- Claim Details -->
                              <div
                                v-else-if="event.type === 'claim'"
                                class="details-grid"
                              >
                                <div class="detail-item">
                                  <span class="detail-label">Claim #:</span>
                                  <span class="detail-value">{{
                                    event.details.claimNumber
                                  }}</span>
                                </div>
                                <div class="detail-item">
                                  <span class="detail-label">Status:</span>
                                  <v-chip
                                    :color="
                                      getClaimStatusColor(event.details.status)
                                    "
                                    size="x-small"
                                  >
                                    {{ event.details.status }}
                                  </v-chip>
                                </div>
                                <div class="detail-item">
                                  <span class="detail-label">Amount:</span>
                                  <span class="detail-value">{{
                                    formatCurrency(event.details.amount)
                                  }}</span>
                                </div>
                              </div>

                              <!-- Beneficiary Change Details -->
                              <div
                                v-else-if="event.type === 'beneficiary_change'"
                                class="details-grid"
                              >
                                <div class="detail-item">
                                  <span class="detail-label">Beneficiary:</span>
                                  <span class="detail-value">{{
                                    event.details.beneficiaryName
                                  }}</span>
                                </div>
                                <div class="detail-item">
                                  <span class="detail-label">Action:</span>
                                  <span class="detail-value">{{
                                    event.details.action
                                  }}</span>
                                </div>
                                <div
                                  v-if="event.details.allocation"
                                  class="detail-item"
                                >
                                  <span class="detail-label">Allocation:</span>
                                  <span class="detail-value"
                                    >{{ event.details.allocation }}%</span
                                  >
                                </div>
                              </div>

                              <!-- Generic Details -->
                              <div
                                v-else-if="typeof event.details === 'object'"
                                class="details-grid"
                              >
                                <div
                                  v-for="(value, key) in event.details"
                                  :key="key"
                                  class="detail-item"
                                >
                                  <span class="detail-label"
                                    >{{ formatDetailKey(key) }}:</span
                                  >
                                  <span class="detail-value">{{ value }}</span>
                                </div>
                              </div>
                            </div>

                            <div
                              v-if="event.performedBy"
                              class="text-caption text-grey mt-2"
                            >
                              Performed by: {{ event.performedBy }}
                            </div>
                          </v-card-text>
                        </v-card>
                      </v-timeline-item>
                    </v-timeline>

                    <div
                      v-if="filteredHistory.length === 0"
                      class="text-center text-grey py-8"
                    >
                      <v-icon size="64" class="mb-4">mdi-calendar-clock</v-icon>
                      <div class="text-h6">No events found</div>
                      <div class="text-body-2"
                        >Try adjusting your filter criteria</div
                      >
                    </div>
                  </div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>

      <v-card-actions class="px-4 py-3">
        <v-btn
          size="small"
          rounded
          color="success"
          variant="outlined"
          @click="exportHistory"
        >
          <v-icon left>mdi-download</v-icon>
          Export History
        </v-btn>
        <v-btn
          size="small"
          rounded
          color="info"
          variant="outlined"
          @click="printHistory"
        >
          <v-icon left>mdi-printer</v-icon>
          Print
        </v-btn>
        <v-spacer />
        <v-btn
          size="small"
          rounded
          color="grey"
          variant="outlined"
          @click="closeDialog"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface HistoryEvent {
  timestamp: Date
  type: string
  title: string
  description: string
  details?: any
  performedBy?: string
}

interface Props {
  modelValue: boolean
  memberId: number | null
  memberName: string
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// State
const loading = ref(false)
const memberHistory = ref<HistoryEvent[]>([])
const selectedEventType = ref('')
const dateRangeStart = ref('')
const dateRangeEnd = ref('')

// Event type options
const eventTypeOptions = [
  { value: 'enrollment', title: 'Enrollment' },
  { value: 'salary_change', title: 'Salary Changes' },
  { value: 'benefit_change', title: 'Benefit Changes' },
  { value: 'status_change', title: 'Status Changes' },
  { value: 'claim', title: 'Claims' },
  { value: 'beneficiary_change', title: 'Beneficiary Changes' },
  { value: 'contact_update', title: 'Contact Updates' },
  { value: 'document', title: 'Document Events' },
  { value: 'payment', title: 'Payments' }
]

// Computed
const dialog = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const filteredHistory = computed(() => {
  let filtered = [...memberHistory.value]

  if (selectedEventType.value) {
    filtered = filtered.filter(
      (event) => event.type === selectedEventType.value
    )
  }

  if (dateRangeStart.value) {
    const startDate = new Date(dateRangeStart.value)
    filtered = filtered.filter(
      (event) => new Date(event.timestamp) >= startDate
    )
  }

  if (dateRangeEnd.value) {
    const endDate = new Date(dateRangeEnd.value)
    endDate.setHours(23, 59, 59, 999) // End of day
    filtered = filtered.filter((event) => new Date(event.timestamp) <= endDate)
  }

  return filtered.sort(
    (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  )
})

const historyStats = computed(() => {
  return {
    totalEvents: memberHistory.value.length,
    policyChanges: memberHistory.value.filter((e) =>
      ['benefit_change', 'enrollment'].includes(e.type)
    ).length,
    claims: memberHistory.value.filter((e) => e.type === 'claim').length,
    statusChanges: memberHistory.value.filter((e) => e.type === 'status_change')
      .length
  }
})

// Watch for dialog opening
watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue && props.memberId) {
      loadMemberHistory()
    }
  }
)

// Methods

const loadMemberHistory = async () => {
  if (!props.memberId) return

  loading.value = true

  try {
    // Simulate API call delay
    await new Promise((resolve) => setTimeout(resolve, 1000))

    const response = await GroupPricingService.getMemberHistory(props.memberId)

    // Generate mock data
    memberHistory.value = response.data
  } catch (err: any) {
    console.error('Error loading member history:', err)
  } finally {
    loading.value = false
  }
}

const closeDialog = () => {
  dialog.value = false
  memberHistory.value = []
  selectedEventType.value = ''
  dateRangeStart.value = ''
  dateRangeEnd.value = ''
}

const formatDateTime = (date: Date) => {
  return new Date(date).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(amount)
}

const formatEventType = (type: string) => {
  return type
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

const formatDetailKey = (key: string | number) => {
  const keyStr = String(key)
  return keyStr
    .split(/[_-]/)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

const getEventColor = (type: string) => {
  const colorMap: { [key: string]: string } = {
    enrollment: 'success',
    salary_change: 'info',
    benefit_change: 'primary',
    status_change: 'warning',
    claim: 'error',
    beneficiary_change: 'purple',
    contact_update: 'teal',
    document: 'orange',
    payment: 'green'
  }
  return colorMap[type] || 'grey'
}

const getEventIcon = (type: string) => {
  const iconMap: { [key: string]: string } = {
    enrollment: 'mdi-account-plus',
    salary_change: 'mdi-cash',
    benefit_change: 'mdi-shield-account',
    status_change: 'mdi-account-switch',
    claim: 'mdi-file-document-alert',
    beneficiary_change: 'mdi-account-heart',
    contact_update: 'mdi-card-account-details',
    document: 'mdi-file-document',
    payment: 'mdi-credit-card'
  }
  return iconMap[type] || 'mdi-information'
}

const getClaimStatusColor = (status: string) => {
  const statusColors: { [key: string]: string } = {
    Submitted: 'info',
    'Under Review': 'warning',
    Approved: 'success',
    Paid: 'success',
    Rejected: 'error',
    Pending: 'orange'
  }
  return statusColors[status] || 'grey'
}

const exportHistory = () => {
  // Implementation for exporting history to Excel/PDF
}

const printHistory = () => {
  // Implementation for printing
  window.print()
}
</script>

<style scoped>
.timeline-container {
  max-height: 800px;
  overflow-y: auto;
}

.timeline-item {
  margin-bottom: 16px;
}

.event-details {
  margin-top: 8px;
  padding: 8px;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.details-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 8px;
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}

.detail-label {
  font-size: 0.75rem;
  font-weight: 500;
  color: rgba(0, 0, 0, 0.6);
}

.detail-value {
  font-size: 0.75rem;
  font-weight: 600;
  text-align: right;
}

@media (max-width: 960px) {
  .details-grid {
    grid-template-columns: 1fr;
  }

  .detail-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .detail-value {
    text-align: left;
  }
}
</style>
