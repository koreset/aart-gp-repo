<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <div>
                <span class="headline">Template Manager</span>
                <p class="text-subtitle-1 text-medium-emphasis mt-2">
                  Manage bordereaux templates for different schemes and
                  customize field mappings
                </p>
              </div>
              <div class="d-flex align-center gap-2">
                <v-btn
                  color="info"
                  variant="outlined"
                  prepend-icon="mdi-download"
                  @click="exportTemplates"
                >
                  Export Templates
                </v-btn>
                <v-btn
                  color="primary"
                  size="large"
                  prepend-icon="mdi-plus"
                  @click="createTemplate"
                >
                  Create Template
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <!-- Template Summary -->
            <v-row class="mb-6">
              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Active</p
                        >
                        <p class="text-h5 font-weight-bold text-success">{{
                          templateStats.active
                        }}</p>
                        <p class="text-caption text-success">Templates</p>
                      </div>
                      <v-icon size="40" color="success">mdi-file-check</v-icon>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Draft</p
                        >
                        <p class="text-h5 font-weight-bold text-warning">{{
                          templateStats.draft
                        }}</p>
                        <p class="text-caption text-warning">In Progress</p>
                      </div>
                      <v-icon size="40" color="warning">mdi-file-edit</v-icon>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>

              <v-col cols="12" sm="6" md="3">
                <v-card variant="outlined" class="h-100">
                  <v-card-text>
                    <div class="d-flex align-center justify-space-between">
                      <div>
                        <p class="text-caption text-medium-emphasis mb-1"
                          >Usage</p
                        >
                        <p class="text-h5 font-weight-bold text-purple">{{
                          templateStats.usage
                        }}</p>
                        <p class="text-caption text-purple">This Month</p>
                      </div>
                      <v-icon size="40" color="purple">mdi-chart-line</v-icon>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Filter Controls -->
            <v-row class="mb-4">
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-h6 font-weight-bold">
                    <v-icon class="me-2">mdi-filter</v-icon>
                    Filter Templates
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" sm="6" md="3">
                        <v-select
                          v-model="filters.type"
                          :items="bordereauTypes"
                          label="Bordereaux Type"
                          variant="outlined"
                          density="compact"
                          clearable
                        />
                      </v-col>
                      <v-col cols="12" sm="6" md="3">
                        <v-select
                          v-model="filters.status"
                          :items="statusOptions"
                          label="Status"
                          variant="outlined"
                          density="compact"
                          clearable
                        />
                      </v-col>
                      <v-col cols="12" sm="6" md="3">
                        <v-text-field
                          v-model="filters.search"
                          label="Search Templates"
                          variant="outlined"
                          density="compact"
                          prepend-inner-icon="mdi-magnify"
                          clearable
                        />
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <!-- Templates Table -->
            <v-row>
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title
                    class="d-flex align-center justify-space-between"
                  >
                    <span class="text-h6 font-weight-bold"
                      >Bordereaux Templates</span
                    >
                    <v-btn
                      color="success"
                      size="small"
                      prepend-icon="mdi-refresh"
                      :loading="loading"
                      @click="refreshTemplates"
                    >
                      Refresh
                    </v-btn>
                  </v-card-title>
                  <v-card-text>
                    <v-data-table
                      :headers="headers"
                      :items="filteredTemplates"
                      :loading="loading"
                      :items-per-page="15"
                      item-key="id"
                    >
                      <template #[`item.name`]="{ item }: { item: any }">
                        <div class="d-flex align-center">
                          <v-icon
                            :color="getTypeColor(item.type)"
                            class="me-2"
                            >{{ getTypeIcon(item.type) }}</v-icon
                          >
                          <div>
                            <p class="font-weight-bold">{{ item.name }}</p>
                            <p class="text-caption text-medium-emphasis"
                              >Version {{ item.version }}</p
                            >
                          </div>
                        </div>
                      </template>

                      <template #[`item.type`]="{ item }: { item: any }">
                        <v-chip
                          :color="getTypeColor(item.type)"
                          size="small"
                          variant="tonal"
                        >
                          {{ formatType(item.type) }}
                        </v-chip>
                      </template>

                      <template #[`item.status`]="{ item }: { item: any }">
                        <v-chip
                          :color="getStatusColor(item.status)"
                          size="small"
                        >
                          {{ formatStatus(item.status) }}
                        </v-chip>
                      </template>

                      <template #[`item.format`]="{ item }: { item: any }">
                        <div class="d-flex align-center">
                          <v-icon
                            :color="getFormatColor(item.format)"
                            class="me-1"
                            size="small"
                          >
                            {{ getFormatIcon(item.format) }}
                          </v-icon>
                          {{ item.format.toUpperCase() }}
                        </div>
                      </template>

                      <template #[`item.usage_count`]="{ item }: { item: any }">
                        <div class="text-center">
                          <div class="font-weight-bold">{{
                            item.usage_count
                          }}</div>
                          <div class="text-caption text-medium-emphasis"
                            >times used</div
                          >
                        </div>
                      </template>

                      <template #[`item.last_used`]="{ item }: { item: any }">
                        {{
                          item.last_used
                            ? formatDateTime(item.last_used)
                            : 'Never'
                        }}
                      </template>

                      <template
                        #[`item.created_date`]="{ item }: { item: any }"
                      >
                        {{ formatDate(item.created_at) }}
                      </template>

                      <template #[`item.actions`]="{ item }: { item: any }">
                        <div class="d-flex align-center gap-1">
                          <v-btn
                            size="small"
                            color="info"
                            variant="tonal"
                            icon="mdi-eye"
                            @click="viewTemplate(item)"
                          />
                          <v-btn
                            size="small"
                            color="primary"
                            variant="tonal"
                            icon="mdi-pencil"
                            @click="editTemplate(item)"
                          />
                          <v-btn
                            size="small"
                            color="orange"
                            variant="tonal"
                            icon="mdi-content-copy"
                            @click="duplicateTemplate(item)"
                          />
                          <v-btn
                            size="small"
                            color="success"
                            variant="tonal"
                            icon="mdi-play"
                            :disabled="item.status !== 'active'"
                            @click="testTemplate(item)"
                          />
                          <v-menu>
                            <template #activator="{ props }">
                              <v-btn
                                size="small"
                                color="grey"
                                variant="tonal"
                                icon="mdi-dots-vertical"
                                v-bind="props"
                              />
                            </template>
                            <v-list density="compact">
                              <v-list-item @click="exportTemplate(item)">
                                <template #prepend>
                                  <v-icon>mdi-export</v-icon>
                                </template>
                                <v-list-item-title>Export</v-list-item-title>
                              </v-list-item>
                              <v-list-item @click="versionHistory(item)">
                                <template #prepend>
                                  <v-icon>mdi-history</v-icon>
                                </template>
                                <v-list-item-title
                                  >Version History</v-list-item-title
                                >
                              </v-list-item>
                              <v-divider></v-divider>
                              <v-list-item
                                v-if="
                                  item.status === 'draft' ||
                                  item.status === 'inactive'
                                "
                                :disabled="updatingStatus"
                                @click="activateTemplate(item)"
                              >
                                <template #prepend>
                                  <v-icon>mdi-check-circle</v-icon>
                                </template>
                                <v-list-item-title>Activate</v-list-item-title>
                              </v-list-item>
                              <v-list-item
                                v-if="item.status === 'active'"
                                :disabled="updatingStatus"
                                @click="deactivateTemplate(item)"
                              >
                                <template #prepend>
                                  <v-icon>mdi-pause-circle</v-icon>
                                </template>
                                <v-list-item-title
                                  >Deactivate</v-list-item-title
                                >
                              </v-list-item>
                              <v-list-item
                                v-if="item.status !== 'active'"
                                class="text-error"
                                @click="deleteTemplate(item)"
                              >
                                <template #prepend>
                                  <v-icon>mdi-delete</v-icon>
                                </template>
                                <v-list-item-title>Delete</v-list-item-title>
                              </v-list-item>
                            </v-list>
                          </v-menu>
                        </div>
                      </template>
                    </v-data-table>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Template Editor Dialog -->
    <v-dialog v-model="showEditorDialog" max-width="1200" persistent>
      <v-card v-if="editingTemplate">
        <v-card-title class="text-h6 font-weight-bold bg-primary text-white">
          <v-icon class="me-2">{{
            editingTemplate.id ? 'mdi-pencil' : 'mdi-plus'
          }}</v-icon>
          {{ editingTemplate.id ? 'Edit Template' : 'Create New Template' }}
        </v-card-title>

        <v-card-text class="pa-0">
          <v-container fluid class="pa-6">
            <v-row>
              <!-- Basic Information -->
              <v-col cols="12" md="6">
                <v-card variant="outlined" class="h-100">
                  <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
                    Basic Information
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12">
                        <v-text-field
                          v-model="editingTemplate.name"
                          label="Template Name *"
                          variant="outlined"
                          density="compact"
                          required
                        />
                      </v-col>
                      <v-col cols="12" sm="6">
                        <v-select
                          v-model="editingTemplate.type"
                          :items="bordereauTypes"
                          label="Bordereaux Type *"
                          variant="outlined"
                          density="compact"
                          required
                        />
                      </v-col>
                      <v-col cols="12" sm="6">
                        <v-select
                          v-model="editingTemplate.format"
                          :items="formatOptions"
                          label="Output Format *"
                          variant="outlined"
                          density="compact"
                          required
                        />
                      </v-col>
                      <v-col cols="12" sm="6">
                        <v-select
                          v-model="editingTemplate.status"
                          :items="statusOptions"
                          label="Status *"
                          variant="outlined"
                          density="compact"
                          required
                        />
                      </v-col>
                      <v-col cols="12">
                        <v-textarea
                          v-model="editingTemplate.description"
                          label="Description"
                          variant="outlined"
                          density="compact"
                          rows="3"
                        />
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>

              <!-- Field Mappings -->
              <v-col cols="12" md="6">
                <v-card variant="outlined" class="h-100">
                  <v-card-title
                    class="text-subtitle-1 bg-grey-lighten-4 d-flex align-center justify-space-between"
                  >
                    Field Mappings
                    <v-btn
                      size="small"
                      color="primary"
                      variant="tonal"
                      prepend-icon="mdi-plus"
                      @click="addFieldMapping"
                    >
                      Add Field
                    </v-btn>
                  </v-card-title>
                  <v-card-text
                    class="pa-0"
                    style="max-height: 400px; overflow-y: auto"
                  >
                    <v-list>
                      <v-list-item
                        v-for="(
                          mapping, index
                        ) in editingTemplate.field_mappings"
                        :key="index"
                        :value="index"
                      >
                        <div class="w-100">
                          <v-row class="mt-1">
                            <v-col cols="5">
                              <v-select
                                v-model="mapping.source_field"
                                :items="bordereauFields"
                                item-title="display_name"
                                item-value="field_name"
                                label="Source Field"
                                variant="outlined"
                                density="compact"
                                hide-details
                                :loading="loadingFields"
                                :disabled="
                                  !editingTemplate.type || loadingFields
                                "
                                clearable
                              />
                            </v-col>
                            <v-col cols="5">
                              <v-text-field
                                v-model="mapping.target_field"
                                label="Target Field"
                                variant="outlined"
                                density="compact"
                                hide-details
                              />
                            </v-col>
                            <v-col cols="2" class="d-flex align-center">
                              <v-checkbox
                                v-model="mapping.required"
                                label="*"
                                density="compact"
                                hide-details
                              />
                              <v-btn
                                size="small"
                                color="error"
                                variant="text"
                                icon="mdi-delete"
                                @click="removeFieldMapping(index)"
                              />
                            </v-col>
                          </v-row>
                        </div>
                      </v-list-item>
                    </v-list>
                  </v-card-text>
                </v-card>
              </v-col>

              <!-- Validation Rules -->
              <v-col cols="12">
                <v-card variant="outlined">
                  <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
                    Validation Rules
                  </v-card-title>
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" md="4">
                        <v-checkbox
                          v-model="
                            editingTemplate.validation_rules.validate_id_numbers
                          "
                          label="Validate SA ID Numbers"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" md="4">
                        <v-checkbox
                          v-model="
                            editingTemplate.validation_rules
                              .validate_banking_details
                          "
                          label="Validate Banking Details"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" md="4">
                        <v-checkbox
                          v-model="
                            editingTemplate.validation_rules.validate_amounts
                          "
                          label="Validate Amount Formats"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" md="4">
                        <v-checkbox
                          v-model="
                            editingTemplate.validation_rules.exclude_invalid
                          "
                          label="Exclude Invalid Records"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" md="4">
                        <v-checkbox
                          v-model="
                            editingTemplate.validation_rules
                              .require_beneficiaries
                          "
                          label="Require Beneficiaries"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" md="4">
                        <v-checkbox
                          v-model="
                            editingTemplate.validation_rules.validate_dates
                          "
                          label="Validate Date Formats"
                          density="compact"
                        />
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>

        <v-card-actions class="pa-6">
          <v-btn
            color="success"
            variant="outlined"
            prepend-icon="mdi-play"
            :disabled="!canTest"
            @click="testCurrentTemplate"
          >
            Test Template
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn color="grey" variant="outlined" @click="cancelEdit">
            Cancel
          </v-btn>
          <v-btn
            color="primary"
            variant="flat"
            :disabled="!canSave"
            :loading="saving"
            @click="saveTemplate"
          >
            {{ editingTemplate.id ? 'Update' : 'Create' }} Template
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="showPreviewDialog" max-width="900">
      <v-card v-if="previewTemplate">
        <v-card-title class="text-h6">
          {{ previewTemplate.name }}
          <v-chip size="small" class="ms-2">{{ previewTemplate.type }}</v-chip>
          <v-chip size="small" class="ms-2">{{ previewTemplate.status }}</v-chip>
        </v-card-title>
        <v-card-text>
          <p v-if="previewTemplate.description" class="mb-4">
            {{ previewTemplate.description }}
          </p>
          <h4 class="text-subtitle-1 font-weight-bold mb-2">
            Field Mappings ({{ (previewTemplate.field_mappings || []).length }})
          </h4>
          <v-data-table
            v-if="(previewTemplate.field_mappings || []).length"
            :headers="[
              { title: 'Source Field', key: 'source_field' },
              { title: 'Target Field', key: 'target_field' },
              { title: 'Required', key: 'required' }
            ]"
            :items="previewTemplate.field_mappings"
            :items-per-page="20"
            density="compact"
            class="mb-4"
          >
            <template #[`item.required`]="{ item }">
              <v-icon v-if="(item as any).required" color="success" size="small">
                mdi-check
              </v-icon>
            </template>
          </v-data-table>
          <div v-else class="text-medium-emphasis mb-4">
            No field mappings defined.
          </div>
          <h4 class="text-subtitle-1 font-weight-bold mb-2">
            Validation Rules
          </h4>
          <v-list density="compact">
            <v-list-item
              v-for="(value, key) in previewTemplate.validation_rules || {}"
              :key="key"
            >
              <v-list-item-title class="d-flex align-center">
                <v-icon
                  :color="value ? 'success' : 'grey'"
                  size="small"
                  class="me-2"
                >
                  {{ value ? 'mdi-check-circle' : 'mdi-circle-outline' }}
                </v-icon>
                {{ key }}
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="showPreviewDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="showTestDialog" max-width="1100">
      <v-card>
        <v-card-title class="text-h6">
          Template Test: {{ testTemplateName }}
        </v-card-title>
        <v-card-text>
          <div v-if="testing" class="d-flex align-center" style="gap: 12px">
            <v-progress-circular indeterminate size="20" />
            <span>Running test with a sample of live data…</span>
          </div>
          <div v-else-if="testResult">
            <v-alert
              v-if="testResult.unknown_fields?.length"
              type="warning"
              variant="tonal"
              class="mb-3"
            >
              Unknown source fields:
              {{ testResult.unknown_fields.join(', ') }}
            </v-alert>
            <v-alert
              v-if="testResult.missing_in_data?.length"
              type="info"
              variant="tonal"
              class="mb-3"
            >
              Fields with no values in the sample:
              {{ testResult.missing_in_data.join(', ') }}
            </v-alert>
            <p class="text-caption text-medium-emphasis mb-3">
              Sample size: {{ testResult.sample_size }} ·
              Source: {{ testResult.sample_source }}
            </p>
            <h4 class="text-subtitle-1 font-weight-bold mb-2">Preview Rows</h4>
            <div
              v-if="!testResult.preview_rows?.length"
              class="text-medium-emphasis"
            >
              No preview rows — the source table may be empty. Generate a
              bordereaux first to create sample data.
            </div>
            <v-data-table
              v-else
              :headers="
                Object.keys(testResult.preview_rows[0] || {}).map((k) => ({
                  title: k,
                  key: k
                }))
              "
              :items="testResult.preview_rows"
              :items-per-page="10"
              density="compact"
            />
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="showTestDialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="showDeleteDialog" max-width="480" persistent>
      <v-card>
        <v-card-title class="text-h6">Delete template?</v-card-title>
        <v-card-text>
          Are you sure you want to delete
          <strong>{{ templateToDelete?.name }}</strong
          >? This cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            variant="text"
            :disabled="deleting"
            @click="showDeleteDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn
            color="error"
            variant="flat"
            :loading="deleting"
            @click="confirmDeleteTemplate"
          >
            Delete
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useFlashStore } from '@/renderer/store/flash'

const flash = useFlashStore()
const showDeleteDialog = ref(false)
const deleting = ref(false)
const templateToDelete: any = ref(null)

// Preview + test dialog state (P2-1)
const showPreviewDialog = ref(false)
const previewTemplate: any = ref(null)
const showTestDialog = ref(false)
const testing = ref(false)
const testResult: any = ref(null)
const testTemplateName = ref('')

// Reactive data
const loading = ref(false)
const saving = ref(false)
const updatingStatus = ref(false)
const showEditorDialog = ref(false)
const editingTemplate: any = ref(null)
const loadingFields = ref(false)
const bordereauFields = ref<any[]>([])

const templateStats = ref({
  active: 12,
  draft: 3,
  insurers: 8,
  usage: 47
})

const filters = ref({
  insurer: null,
  type: '',
  status: '',
  search: ''
})

// Static data
const headers = [
  { title: 'Template Name', key: 'name', sortable: true },
  { title: 'Scheme', key: 'scheme_name', sortable: true },
  { title: 'Type', key: 'type', sortable: true },
  { title: 'Status', key: 'status', sortable: true },
  { title: 'Format', key: 'format', sortable: true },
  { title: 'Usage', key: 'usage_count', sortable: true },
  { title: 'Last Used', key: 'last_used', sortable: true },
  { title: 'Created', key: 'created_date', sortable: true },
  { title: 'Actions', key: 'actions', sortable: false, width: '250px' }
]

const bordereauTypes = [
  { title: 'Member', value: 'member' },
  { title: 'Premium', value: 'premium' },
  { title: 'Claims', value: 'claims' },
  { title: 'Benefits', value: 'benefits' },
  { title: 'Reinsurance Premium', value: 'reinsurance_premium' },
  { title: 'Reinsurance Claims', value: 'reinsurance_claims' }
]

const statusOptions = [
  { title: 'Draft', value: 'draft' },
  { title: 'Active', value: 'active' },
  { title: 'Inactive', value: 'inactive' },
  { title: 'Deprecated', value: 'deprecated' }
]

const formatOptions = [
  { title: 'Excel (.xlsx)', value: 'excel' },
  { title: 'CSV', value: 'csv' },
  { title: 'XML', value: 'xml' },
  { title: 'JSON', value: 'json' }
]

const schemes = ref([
  { id: 1, name: 'Liberty Life' },
  { id: 2, name: 'Old Mutual' },
  { id: 3, name: 'Momentum' },
  { id: 4, name: 'Discovery Life' },
  { id: 5, name: 'Sanlam' },
  { id: 6, name: 'Metropolitan' },
  { id: 7, name: 'AIG' },
  { id: 8, name: 'Hollard' },
  { id: 9, name: 'General' }
])

const templates: any = ref([])

// Computed properties
const filteredTemplates = computed(() => {
  let filtered = templates.value

  if (filters.value.type) {
    filtered = filtered.filter((t) => t.type === filters.value.type)
  }

  if (filters.value.status) {
    filtered = filtered.filter((t) => t.status === filters.value.status)
  }

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(
      (t) =>
        t.name.toLowerCase().includes(search) ||
        t.description?.toLowerCase().includes(search)
    )
  }

  return filtered
})

const canSave = computed(() => {
  return (
    !saving.value &&
    editingTemplate.value?.name &&
    editingTemplate.value?.type &&
    editingTemplate.value?.format &&
    editingTemplate.value?.field_mappings?.length > 0
  )
})

const canTest = computed(() => {
  return canSave.value && editingTemplate.value?.field_mappings?.length > 0
})

// Methods
const getTypeColor = (type: string): string => {
  const colors: Record<string, string> = {
    member: 'blue',
    premium: 'green',
    claims: 'orange',
    benefits: 'purple'
  }
  return colors[type] || 'grey'
}

const getTypeIcon = (type: string): string => {
  const icons: Record<string, string> = {
    member: 'mdi-account-group',
    premium: 'mdi-currency-usd',
    claims: 'mdi-medical-bag',
    benefits: 'mdi-gift'
  }
  return icons[type] || 'mdi-file'
}

const formatType = (type: string): string => {
  const types: Record<string, string> = {
    member: 'Member',
    premium: 'Premium',
    claims: 'Claims',
    benefits: 'Benefits'
  }
  return types[type] || type
}

const getStatusColor = (status: string): string => {
  const colors: Record<string, string> = {
    active: 'success',
    draft: 'warning',
    inactive: 'grey',
    deprecated: 'error'
  }
  return colors[status] || 'grey'
}

const formatStatus = (status: string): string => {
  return status.charAt(0).toUpperCase() + status.slice(1)
}

const getFormatColor = (format: string): string => {
  const colors: Record<string, string> = {
    excel: 'green',
    csv: 'blue',
    xml: 'orange',
    json: 'purple'
  }
  return colors[format] || 'grey'
}

const getFormatIcon = (format: string): string => {
  const icons: Record<string, string> = {
    excel: 'mdi-file-excel',
    csv: 'mdi-file-delimited',
    xml: 'mdi-file-code',
    json: 'mdi-code-json'
  }
  return icons[format] || 'mdi-file'
}

const formatDateTime = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-ZA', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const refreshTemplates = async () => {
  loading.value = true
  try {
    const templatesResponse = await GroupPricingService.getBordereauxTemplates()
    templates.value = templatesResponse.data

    // Generate template statistics based on fetched data
    const stats = {
      active: 0,
      draft: 0,
      insurers: new Set(),
      usage: 0
    }

    templates.value.forEach((template: any) => {
      // Count by status
      if (template.status === 'active') {
        stats.active++
      } else if (template.status === 'draft') {
        stats.draft++
      }

      // Track unique insurers/schemes
      if (template.insurer_id) {
        stats.insurers.add(template.insurer_id)
      }

      // Sum usage counts
      if (template.usage_count) {
        stats.usage += template.usage_count
      }
    })

    // Update template stats
    templateStats.value = {
      active: stats.active,
      draft: stats.draft,
      insurers: stats.insurers.size,
      usage: stats.usage
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to fetch templates',
      'error'
    )
    templates.value = []
  } finally {
    loading.value = false
  }
}

const fetchBordereauFields = async (bordereauType: string) => {
  if (!bordereauType) {
    bordereauFields.value = []
    return
  }

  try {
    loadingFields.value = true
    const response =
      await GroupPricingService.getBordereauxFields(bordereauType)
    console.log('Fetched bordereaux fields:', response.data)
    bordereauFields.value = response.data || []
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to load bordereaux fields',
      'error'
    )
    bordereauFields.value = []
  } finally {
    loadingFields.value = false
  }
}

const createTemplate = () => {
  editingTemplate.value = {
    name: '',
    insurer_id: null,
    type: '',
    status: 'draft',
    format: 'excel',
    description: '',
    field_mappings: [],
    validation_rules: {
      validate_id_numbers: true,
      validate_banking_details: true,
      validate_amounts: true,
      exclude_invalid: false,
      require_beneficiaries: false,
      validate_dates: true
    }
  }
  // Clear previous fields
  bordereauFields.value = []
  showEditorDialog.value = true
}

const editTemplate = (template: any) => {
  editingTemplate.value = { ...template }
  // Fetch fields for the selected type
  if (template.type) {
    fetchBordereauFields(template.type)
  }
  showEditorDialog.value = true
}

const viewTemplate = (template: any) => {
  previewTemplate.value = template
  showPreviewDialog.value = true
}

const duplicateTemplate = (template: any) => {
  editingTemplate.value = {
    ...template,
    id: null,
    name: `${template.name} (Copy)`,
    status: 'draft',
    version: '1.0',
    usage_count: 0,
    last_used: null,
    created_date: new Date().toISOString()
  }
  // Fetch fields for the duplicated template type
  if (template.type) {
    fetchBordereauFields(template.type)
  }
  showEditorDialog.value = true
}

const runTemplateTest = async (templateId: number, displayName: string) => {
  try {
    testing.value = true
    testTemplateName.value = displayName
    testResult.value = null
    showTestDialog.value = true
    const response = await GroupPricingService.testBordereauxTemplate(
      templateId,
      { sample_size: 5 }
    )
    testResult.value = response.data?.data ?? response.data
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to run template test',
      'error'
    )
    showTestDialog.value = false
  } finally {
    testing.value = false
  }
}

const testTemplate = (template: any) => {
  if (!template?.id) {
    flash.show('Save the template before testing', 'warning')
    return
  }
  runTemplateTest(template.id, template.name)
}

const testCurrentTemplate = () => {
  if (!editingTemplate.value?.id) {
    flash.show('Save the template before testing', 'warning')
    return
  }
  runTemplateTest(editingTemplate.value.id, editingTemplate.value.name)
}

const addFieldMapping = () => {
  if (!editingTemplate.value.field_mappings) {
    editingTemplate.value.field_mappings = []
  }
  editingTemplate.value.field_mappings.push({
    source_field: '',
    target_field: '',
    required: false
  })
}

const removeFieldMapping = (index: number | string) => {
  const numIndex = typeof index === 'string' ? parseInt(index, 10) : index
  editingTemplate.value.field_mappings.splice(numIndex, 1)
}

const saveTemplate = async () => {
  try {
    saving.value = true

    console.log('Saving template:', editingTemplate.value)

    const templateData = {
      name: editingTemplate.value.name,
      type: editingTemplate.value.type,
      status: editingTemplate.value.status,
      format: editingTemplate.value.format,
      description: editingTemplate.value.description,
      field_mappings: editingTemplate.value.field_mappings,
      validation_rules: editingTemplate.value.validation_rules
    }

    let response
    if (editingTemplate.value.id) {
      // Update existing template
      response = await GroupPricingService.updateBordereauxTemplate(
        editingTemplate.value.id,
        templateData
      )

      // Update local template list
      const index = templates.value.findIndex(
        (t) => t.id === editingTemplate.value.id
      )
      if (index !== -1) {
        templates.value[index] = {
          ...response.data,
          scheme_name: schemes.value.find(
            (s) => s.id === response.data.scheme_id
          )?.name
        }
      }
    } else {
      // Create new template
      response =
        await GroupPricingService.createBordereauxTemplate(templateData)

      // Add new template to local list
      const newTemplate = {
        ...response.data,
        scheme_name: schemes.value.find((s) => s.id === response.data.scheme_id)
          ?.name
      }
      templates.value.push(newTemplate)
    }

    showEditorDialog.value = false
    editingTemplate.value = null
    bordereauFields.value = []

    flash.show(`Template '${response.data.name}' saved`, 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.response?.data?.message ||
        error.message ||
        'Failed to save template',
      'error'
    )
  } finally {
    saving.value = false
  }
}

const cancelEdit = () => {
  showEditorDialog.value = false
  editingTemplate.value = null
  bordereauFields.value = []
}

const downloadJSON = (payload: unknown, filename: string) => {
  const blob = new Blob([JSON.stringify(payload, null, 2)], {
    type: 'application/json'
  })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}

const slugify = (value: string): string =>
  (value || 'template')
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '') || 'template'

const exportTemplate = (template: any) => {
  if (!template) return
  downloadJSON(template, `${slugify(template.name)}.bordereaux-template.json`)
  flash.show(`Exported '${template.name}'`, 'success')
}

const exportTemplates = () => {
  if (!templates.value.length) {
    flash.show('No templates to export', 'warning')
    return
  }
  downloadJSON(templates.value, 'bordereaux-templates.json')
  flash.show(`Exported ${templates.value.length} templates`, 'success')
}

const versionHistory = (_template: any) => {
  // Version history requires a dedicated history table + audit wiring — not yet built.
  flash.show('Template version history is not available yet', 'info')
}

const activateTemplate = async (template: any) => {
  try {
    updatingStatus.value = true
    console.log('Activating template:', template.name)

    // Prepare template data for API update
    const templateData = {
      name: template.name,
      type: template.type,
      status: 'active',
      format: template.format,
      description: template.description,
      field_mappings: template.field_mappings || [],
      validation_rules: template.validation_rules || {
        validate_id_numbers: true,
        validate_banking_details: true,
        validate_amounts: true,
        exclude_invalid: false,
        require_beneficiaries: false,
        validate_dates: true
      }
    }

    // Update template status via API
    const response = await GroupPricingService.updateBordereauxTemplate(
      template.id,
      templateData
    )

    // Update local template state only after successful API response
    const index = templates.value.findIndex((t) => t.id === template.id)
    if (index !== -1) {
      templates.value[index] = {
        ...response.data,
        scheme_name: schemes.value.find((s) => s.id === response.data.scheme_id)
          ?.name
      }
    }

    flash.show(`Template '${template.name}' activated`, 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.response?.data?.message ||
        error.message ||
        'Failed to activate template',
      'error'
    )
  } finally {
    updatingStatus.value = false
  }
}

const deactivateTemplate = async (template: any) => {
  try {
    updatingStatus.value = true
    console.log('Deactivating template:', template.name)

    // Prepare template data for API update
    const templateData = {
      name: template.name,
      type: template.type,
      status: 'inactive',
      format: template.format,
      description: template.description,
      field_mappings: template.field_mappings || [],
      validation_rules: template.validation_rules || {
        validate_id_numbers: true,
        validate_banking_details: true,
        validate_amounts: true,
        exclude_invalid: false,
        require_beneficiaries: false,
        validate_dates: true
      }
    }

    // Update template status via API
    const response = await GroupPricingService.updateBordereauxTemplate(
      template.id,
      templateData
    )

    // Update local template state only after successful API response
    const index = templates.value.findIndex((t) => t.id === template.id)
    if (index !== -1) {
      templates.value[index] = {
        ...response.data,
        scheme_name: schemes.value.find((s) => s.id === response.data.scheme_id)
          ?.name
      }
    }

    flash.show(`Template '${template.name}' deactivated`, 'success')
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.response?.data?.message ||
        error.message ||
        'Failed to deactivate template',
      'error'
    )
  } finally {
    updatingStatus.value = false
  }
}

const deleteTemplate = (template: any) => {
  templateToDelete.value = template
  showDeleteDialog.value = true
}

const confirmDeleteTemplate = async () => {
  const template = templateToDelete.value
  if (!template?.id) {
    showDeleteDialog.value = false
    return
  }
  try {
    deleting.value = true
    await GroupPricingService.deleteBordereauxTemplate(template.id)
    const index = templates.value.findIndex((t) => t.id === template.id)
    if (index !== -1) {
      templates.value.splice(index, 1)
    }
    flash.show(`Template '${template.name}' deleted`, 'success')
    showDeleteDialog.value = false
    templateToDelete.value = null
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.response?.data?.message ||
        error.message ||
        'Failed to delete template',
      'error'
    )
  } finally {
    deleting.value = false
  }
}

// Watchers
watch(
  () => editingTemplate.value?.type,
  (newType) => {
    if (newType) {
      fetchBordereauFields(newType)
    }
  }
)

onMounted(() => {
  refreshTemplates()
})
</script>

<style scoped>
.w-100 {
  width: 100%;
}
</style>
