<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <div>
                <span class="headline">Generate Bordereaux</span>
                <p class="text-subtitle-1 text-medium-emphasis mt-2">
                  Create comprehensive bordereaux for member data, premiums,
                  claims, and benefits
                </p>
              </div>
              <v-btn
                color="grey"
                variant="outlined"
                prepend-icon="mdi-arrow-left"
                @click="$router.push('/group-pricing/bordereaux-management')"
              >
                Back to Dashboard
              </v-btn>
            </div>
          </template>
          <template #default>
            <v-form ref="form" @submit.prevent="generateBordereaux">
              <v-row>
                <!-- Load Saved Configuration -->
                <v-col cols="12">
                  <v-card variant="outlined" class="mb-6">
                    <v-card-title
                      class="text-h6 font-weight-bold bg-secondary text-white"
                    >
                      <v-icon class="me-2">mdi-content-save-settings</v-icon>
                      Saved Configurations
                    </v-card-title>
                    <v-card-text class="pt-4">
                      <v-row>
                        <v-col cols="12" md="4">
                          <v-select
                            v-model="selectedConfigurationId"
                            :items="savedConfigurations"
                            item-title="name"
                            item-value="id"
                            label="Load Saved Configuration"
                            variant="outlined"
                            density="compact"
                            clearable
                            @update:model-value="loadConfiguration"
                          >
                            <template #item="{ props, item }">
                              <v-list-item v-bind="props">
                                <v-list-item-title>{{
                                  item.raw.name
                                }}</v-list-item-title>
                                <v-list-item-subtitle class="text-wrap">
                                  {{ item.raw.description }}
                                </v-list-item-subtitle>
                                <template #append>
                                  <div
                                    class="text-caption text-medium-emphasis"
                                  >
                                    Used {{ item.raw.usage_count }} times
                                  </div>
                                </template>
                              </v-list-item>
                            </template>
                          </v-select>
                        </v-col>
                        <v-col cols="12" md="4" class="d-flex align-center">
                          <v-btn
                            v-if="selectedConfigurationId"
                            color="error"
                            variant="outlined"
                            size="small"
                            @click="deleteConfiguration"
                          >
                            <v-icon class="me-1">mdi-delete</v-icon>
                            Delete
                          </v-btn>
                        </v-col>
                      </v-row>
                    </v-card-text>
                  </v-card>
                </v-col>

                <!-- Bordereaux Configuration -->
                <v-col cols="12" lg="8">
                  <v-card variant="outlined" class="mb-6">
                    <v-card-title
                      class="text-h6 font-weight-bold bg-primary text-white"
                    >
                      Bordereaux Configuration
                    </v-card-title>
                    <v-card-text class="pt-6">
                      <v-row>
                        <!-- Bordereaux Type -->
                        <v-col cols="12" md="6">
                          <v-select
                            v-model="formData.type"
                            :items="bordereauTypes"
                            item-title="label"
                            item-value="value"
                            label="Bordereaux Type *"
                            variant="outlined"
                            density="compact"
                            :rules="[rules.required]"
                            required
                            @update:model-value="onTypeChange"
                          >
                            <template #item="{ props, item }">
                              <v-list-item v-bind="props">
                                <template #prepend>
                                  <v-icon
                                    :icon="item.raw.icon"
                                    :color="item.raw.color"
                                  ></v-icon>
                                </template>
                                <v-list-item-subtitle>{{
                                  item.raw.description
                                }}</v-list-item-subtitle>
                              </v-list-item>
                            </template>
                          </v-select>
                        </v-col>

                        <!-- Scheme Selection -->
                        <v-col cols="12" md="6">
                          <v-autocomplete
                            v-model="formData.scheme_ids"
                            :items="schemes"
                            item-title="name"
                            item-value="id"
                            :label="schemeSelectionLabel"
                            variant="outlined"
                            density="compact"
                            :rules="[rules.required]"
                            multiple
                            chips
                            closable-chips
                            required
                          >
                            <template #prepend-item>
                              <v-list-item title="Select All" @click="toggle">
                                <template #prepend>
                                  <v-checkbox-btn
                                    :color="someSchemes ? 'primary' : undefined"
                                    :indeterminate="someSchemes && !allSchemes"
                                    :model-value="allSchemes"
                                  ></v-checkbox-btn>
                                </template>
                              </v-list-item>
                              <v-divider class="mt-2"></v-divider>
                            </template>
                            <template #selection="{ index }">
                              <!-- Show first few schemes as chips -->
                              <v-chip
                                v-if="index < maxDisplayedChips"
                                size="small"
                                color="primary"
                                closable
                                @click:close="
                                  removeScheme(formData.scheme_ids[index])
                                "
                              >
                                {{
                                  schemes.find(
                                    (s: any) =>
                                      s.id === formData.scheme_ids[index]
                                  )?.name
                                }}
                              </v-chip>
                              <!-- Show count of remaining items -->
                              <v-chip
                                v-else-if="
                                  index === maxDisplayedChips &&
                                  hiddenSchemesCount > 0
                                "
                                size="small"
                                color="secondary"
                                variant="outlined"
                              >
                                +{{ hiddenSchemesCount }} more
                              </v-chip>
                            </template>
                          </v-autocomplete>
                        </v-col>

                        <!-- Period Selection -->
                        <v-col cols="12" md="6">
                          <v-select
                            v-model="formData.period_type"
                            :items="periodTypes"
                            label="Period Type *"
                            variant="outlined"
                            density="compact"
                            :rules="[rules.required]"
                            required
                            @update:model-value="onPeriodTypeChange"
                          />
                        </v-col>

                        <!-- Date Range -->
                        <v-col
                          v-if="formData.period_type === 'custom'"
                          cols="12"
                          md="6"
                        >
                          <v-row>
                            <v-col cols="6">
                              <v-text-field
                                v-model="formData.start_date"
                                label="Start Date *"
                                variant="outlined"
                                density="compact"
                                type="date"
                                :rules="[rules.required]"
                                required
                              />
                            </v-col>
                            <v-col cols="6">
                              <v-text-field
                                v-model="formData.end_date"
                                label="End Date *"
                                variant="outlined"
                                density="compact"
                                type="date"
                                :rules="[rules.required]"
                                required
                              />
                            </v-col>
                          </v-row>
                        </v-col>

                        <!-- Month/Year Selection for non-custom periods -->
                        <v-col v-else cols="12" md="6">
                          <v-row>
                            <v-col cols="6">
                              <v-select
                                v-model="formData.month"
                                :items="months"
                                label="Month"
                                variant="outlined"
                                density="compact"
                                :disabled="formData.period_type === 'annual'"
                              />
                            </v-col>
                            <v-col cols="6">
                              <v-select
                                v-model="formData.year"
                                :items="years"
                                label="Year *"
                                variant="outlined"
                                density="compact"
                                :rules="[rules.required]"
                                required
                              />
                            </v-col>
                          </v-row>
                        </v-col>

                        <!-- Template Selection -->
                        <v-col cols="12" md="6">
                          <v-select
                            v-model="formData.template_id"
                            :items="availableTemplates"
                            item-title="name"
                            item-value="id"
                            label="Template *"
                            variant="outlined"
                            density="compact"
                            :rules="[rules.required]"
                            required
                          >
                            <template #item="{ props, item }">
                              <v-list-item v-bind="props">
                                <v-list-item-title>{{
                                  item.raw.name
                                }}</v-list-item-title>
                                <v-list-item-subtitle
                                  >{{ item.raw.insurer_name }} -
                                  {{ item.raw.format }}</v-list-item-subtitle
                                >
                              </v-list-item>
                            </template>
                          </v-select>
                        </v-col>

                        <!-- Output Format -->
                        <v-col cols="12" md="6">
                          <v-select
                            v-model="formData.output_format"
                            :items="outputFormats"
                            label="Output Format *"
                            variant="outlined"
                            density="compact"
                            :rules="[rules.required]"
                            required
                          />
                        </v-col>

                        <!-- Per-Scheme Generation Option -->
                        <v-col v-if="formData.scheme_ids.length > 1" cols="12">
                          <v-checkbox
                            v-model="formData.generate_per_scheme"
                            density="compact"
                            color="primary"
                          >
                            <template #label>
                              <div>
                                <div class="font-weight-medium"
                                  >Generate separate files per scheme</div
                                >
                                <div class="text-caption text-medium-emphasis">
                                  Create individual bordereaux files for each
                                  selected scheme instead of one combined file
                                </div>
                              </div>
                            </template>
                          </v-checkbox>
                        </v-col>
                      </v-row>
                    </v-card-text>
                  </v-card>

                  <!-- Advanced Options -->
                  <v-card variant="outlined" class="mb-6">
                    <v-card-title class="text-h6 font-weight-bold">
                      <v-icon class="me-2">mdi-cog</v-icon>
                      Advanced Options
                    </v-card-title>
                    <v-card-text>
                      <v-row>
                        <!-- Include Options -->
                        <v-col cols="12" md="6">
                          <p class="text-subtitle-2 font-weight-bold mb-3"
                            >Include Data</p
                          >
                          <v-checkbox
                            v-model="formData.include_terminated"
                            label="Include terminated members"
                            density="compact"
                          />
                          <v-checkbox
                            v-model="formData.include_dependants"
                            label="Include dependant data"
                            density="compact"
                          />
                          <v-checkbox
                            v-model="formData.include_beneficiaries"
                            label="Include beneficiary information"
                            density="compact"
                          />
                        </v-col>

                        <!-- Validation Options -->
                        <v-col cols="12" md="6">
                          <p class="text-subtitle-2 font-weight-bold mb-3"
                            >Validation</p
                          >
                          <v-checkbox
                            v-model="formData.validate_id_numbers"
                            label="Validate South African ID numbers"
                            density="compact"
                          />
                          <v-checkbox
                            v-model="formData.validate_banking_details"
                            label="Validate banking details"
                            density="compact"
                          />
                          <v-checkbox
                            v-model="formData.exclude_invalid"
                            label="Exclude invalid records"
                            density="compact"
                          />
                        </v-col>

                        <!-- Additional Settings -->
                        <v-col cols="12">
                          <v-text-field
                            v-model="formData.reference_number"
                            label="Reference Number (Optional)"
                            variant="outlined"
                            density="compact"
                            hint="Internal reference for tracking"
                            persistent-hint
                          />
                        </v-col>

                        <v-col cols="12">
                          <v-textarea
                            v-model="formData.notes"
                            label="Notes (Optional)"
                            variant="outlined"
                            density="compact"
                            rows="3"
                            hint="Additional notes or special instructions"
                            persistent-hint
                          />
                        </v-col>
                      </v-row>
                    </v-card-text>
                  </v-card>
                </v-col>

                <!-- Summary Panel -->
                <v-col cols="12" lg="4">
                  <v-card variant="outlined" class="sticky-card">
                    <v-card-title
                      class="text-h6 font-weight-bold bg-info text-white"
                    >
                      Generation Summary
                    </v-card-title>
                    <v-card-text class="pt-4">
                      <v-list density="compact">
                        <v-list-item>
                          <template #prepend>
                            <v-icon color="primary">mdi-file-document</v-icon>
                          </template>
                          <v-list-item-title>Type</v-list-item-title>
                          <v-list-item-subtitle>
                            {{ getSelectedType()?.label || 'Not selected' }}
                          </v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item>
                          <template #prepend>
                            <v-icon color="green">mdi-office-building</v-icon>
                          </template>
                          <v-list-item-title>Schemes</v-list-item-title>
                          <v-list-item-subtitle>
                            {{ formData.scheme_ids.length }} selected
                            <span
                              v-if="
                                formData.scheme_ids.length > 1 &&
                                formData.generate_per_scheme
                              "
                            >
                              (separate files)
                            </span>
                          </v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item>
                          <template #prepend>
                            <v-icon color="orange">mdi-calendar</v-icon>
                          </template>
                          <v-list-item-title>Period</v-list-item-title>
                          <v-list-item-subtitle>
                            {{ getPeriodSummary() }}
                          </v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item>
                          <template #prepend>
                            <v-icon color="purple">mdi-file-cog</v-icon>
                          </template>
                          <v-list-item-title>Template</v-list-item-title>
                          <v-list-item-subtitle>
                            {{ getSelectedTemplate()?.name || 'Not selected' }}
                          </v-list-item-subtitle>
                        </v-list-item>

                        <v-list-item>
                          <template #prepend>
                            <v-icon color="teal">mdi-file-export</v-icon>
                          </template>
                          <v-list-item-title>Format</v-list-item-title>
                          <v-list-item-subtitle>
                            {{
                              formData.output_format.toUpperCase() ||
                              'Not selected'
                            }}
                          </v-list-item-subtitle>
                        </v-list-item>
                      </v-list>

                      <v-divider class="my-4"></v-divider>

                      <v-alert
                        v-if="estimatedRecords > 0"
                        color="info"
                        variant="tonal"
                        density="compact"
                        class="mb-4"
                      >
                        <template #prepend>
                          <v-icon>mdi-information</v-icon>
                        </template>
                        <div>
                          Estimated
                          {{ estimatedRecords.toLocaleString() }} records
                          <div
                            v-if="
                              formData.scheme_ids.length > 1 &&
                              formData.generate_per_scheme
                            "
                            class="text-caption"
                          >
                            Will generate
                            {{ formData.scheme_ids.length }} separate files
                          </div>
                        </div>
                      </v-alert>

                      <v-btn
                        color="primary"
                        size="large"
                        block
                        :loading="loading"
                        :disabled="!canGenerate"
                        @click="generateBordereaux"
                      >
                        <v-icon class="me-2">mdi-cog</v-icon>
                        Generate Bordereaux
                      </v-btn>

                      <v-btn
                        color="secondary"
                        variant="outlined"
                        size="large"
                        block
                        class="mt-3"
                        :disabled="!canGenerate"
                        @click="openSaveConfigDialog"
                      >
                        <v-icon class="me-2">mdi-content-save</v-icon>
                        Save Configuration
                      </v-btn>

                      <v-btn
                        v-if="selectedConfigurationId"
                        color="warning"
                        variant="outlined"
                        size="large"
                        block
                        class="mt-3"
                        :disabled="!canGenerate"
                        @click="openUpdateConfigDialog"
                      >
                        <v-icon class="me-2">mdi-pencil</v-icon>
                        Update Configuration
                      </v-btn>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>
            </v-form>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Save/Update Configuration Dialog -->
    <v-dialog v-model="saveConfigDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h6 font-weight-bold">
          <v-icon class="me-2">{{
            isUpdateMode ? 'mdi-pencil' : 'mdi-content-save'
          }}</v-icon>
          {{ isUpdateMode ? 'Update Configuration' : 'Save Configuration' }}
        </v-card-title>
        <v-card-text class="pt-4">
          <v-text-field
            v-model="configName"
            label="Configuration Name *"
            variant="outlined"
            density="compact"
            :rules="[rules.required]"
            required
            hint="Enter a descriptive name for this configuration"
            persistent-hint
            :readonly="isUpdateMode"
          />
          <v-textarea
            v-model="configDescription"
            label="Description"
            variant="outlined"
            density="compact"
            rows="3"
            class="mt-4"
            hint="Optional description to help identify this configuration"
            persistent-hint
          />
        </v-card-text>
        <v-card-actions>
          <v-btn
            color="primary"
            variant="flat"
            :disabled="!configName.trim()"
            @click="isUpdateMode ? updateConfiguration() : saveConfiguration()"
          >
            <v-icon class="me-2">{{
              isUpdateMode ? 'mdi-pencil' : 'mdi-content-save'
            }}</v-icon>
            {{ isUpdateMode ? 'Update' : 'Save' }}
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn color="grey" variant="outlined" @click="closeSaveConfigDialog">
            Cancel
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Progress Dialog -->
    <v-dialog v-model="progressDialog" persistent max-width="500">
      <v-card>
        <v-card-title class="text-h6 font-weight-bold">
          Generating Bordereaux
        </v-card-title>
        <v-card-text class="py-4">
          <v-progress-linear
            :model-value="progress"
            color="primary"
            height="8"
            striped
          />
          <p class="text-center mt-3 text-body-2">
            {{ progressMessage }}
          </p>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- Success Dialog -->
    <v-dialog v-model="successDialog" max-width="600">
      <v-card>
        <v-card-title class="text-h6 font-weight-bold bg-success text-white">
          <v-icon class="me-2">mdi-check-circle</v-icon>
          Bordereaux Generated Successfully
        </v-card-title>
        <v-card-text class="pt-4">
          <v-alert color="success" variant="tonal" class="mb-4">
            <p class="font-weight-bold">Generation Complete</p>
            <p class="mb-0">
              {{ generationResult.record_count }} records processed successfully
              <span v-if="generationResult.is_multi_file">
                across {{ generationResult.file_list.length }} scheme files
              </span>
            </p>
          </v-alert>

          <v-list density="compact">
            <v-list-item>
              <template #prepend>
                <v-icon color="primary">mdi-identifier</v-icon>
              </template>
              <v-list-item-title>Bordereaux ID</v-list-item-title>
              <v-list-item-subtitle>{{
                generationResult.bordereaux_id
              }}</v-list-item-subtitle>
            </v-list-item>

            <v-list-item v-if="!generationResult.is_multi_file">
              <template #prepend>
                <v-icon color="green">mdi-file-download</v-icon>
              </template>
              <v-list-item-title>File Size</v-list-item-title>
              <v-list-item-subtitle>{{
                generationResult.file_size
              }}</v-list-item-subtitle>
            </v-list-item>

            <v-list-item v-if="generationResult.is_multi_file">
              <template #prepend>
                <v-icon color="green">mdi-file-multiple</v-icon>
              </template>
              <v-list-item-title>Files Generated</v-list-item-title>
              <v-list-item-subtitle
                >{{
                  generationResult.file_list.length
                }}
                files</v-list-item-subtitle
              >
            </v-list-item>

            <v-list-item>
              <template #prepend>
                <v-icon color="orange">mdi-clock</v-icon>
              </template>
              <v-list-item-title>Processing Time</v-list-item-title>
              <v-list-item-subtitle>{{
                generationResult.processing_time
              }}</v-list-item-subtitle>
            </v-list-item>
          </v-list>

          <!-- Multi-file details -->
          <v-expansion-panels
            v-if="generationResult.is_multi_file"
            class="mt-4"
          >
            <v-expansion-panel>
              <v-expansion-panel-title>
                <v-icon class="me-2">mdi-file-document-multiple</v-icon>
                Generated Files ({{ generationResult.file_list.length }})
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <v-list density="compact">
                  <v-list-item
                    v-for="file in generationResult.file_list"
                    :key="file.name"
                  >
                    <template #prepend>
                      <v-icon color="blue">mdi-file-document</v-icon>
                    </template>
                    <v-list-item-title>{{
                      file.scheme_name
                    }}</v-list-item-title>
                    <v-list-item-subtitle
                      >{{ file.name }} ({{ file.size }})</v-list-item-subtitle
                    >
                  </v-list-item>
                </v-list>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-card-text>
        <v-card-actions>
          <v-btn color="primary" variant="flat" @click="downloadFile">
            <v-icon class="me-2">mdi-download</v-icon>
            <span
              v-if="
                generationResult.is_multi_file &&
                generationResult.download_type === 'zip'
              "
            >
              Download ZIP File
            </span>
            <span v-else-if="generationResult.is_multi_file">
              Download All Files
            </span>
            <span v-else> Download File </span>
          </v-btn>
          <v-btn
            color="info"
            variant="outlined"
            @click="viewSubmissionTracking"
          >
            <v-icon class="me-2">mdi-send</v-icon>
            Submit to Scheme
          </v-btn>
          <v-spacer></v-spacer>
          <v-btn color="grey" variant="outlined" @click="closeSuccessDialog">
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
  <confirm-dialog ref="confirmationDialog" />
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import ConfirmDialog from '@/renderer/components/ConfirmDialog.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useFlashStore } from '@/renderer/store/flash'

const flash = useFlashStore()

interface BordereauxFormData {
  type: string
  scheme_ids: number[]
  period_type: string
  start_date: string
  end_date: string
  month: number | null
  year: number
  template_id: number | null
  output_format: string
  generate_per_scheme: boolean
  include_terminated: boolean
  include_dependants: boolean
  include_beneficiaries: boolean
  validate_id_numbers: boolean
  validate_banking_details: boolean
  exclude_invalid: boolean
  reference_number: string
  notes: string
}

interface SavedConfiguration {
  id: number
  name: string
  description: string
  config_data: BordereauxFormData
  created_by: string
  created_at: string
  last_used: string
  usage_count: number
}

// interface Props {
//   schemes: Array<any>
//   templates: Array<any>
// }

// defineProps<Props>()

const router = useRouter()
const form = ref(null)
const loading = ref(false)
const progressDialog = ref(false)
const successDialog = ref(false)
const progress = ref(0)
const progressMessage = ref('')
const estimatedRecords = ref(0)
const schemes: any = ref([])

// Saved configurations
const savedConfigurations = ref<SavedConfiguration[]>([])
const selectedConfigurationId = ref<number | null>(null)
const saveConfigDialog = ref(false)
const confirmationDialog: any = ref(null)
const configName = ref('')
const configDescription = ref('')
const isUpdateMode = ref(false)

// Job tracking properties
// const currentJobId = ref<string | null>(null)
// const jobStatus = ref<string>('')
// const jobProgress = ref(0)
const jobResult: any = ref(null)
const pollInterval: any = ref(null)

// Form data
const formData = ref<BordereauxFormData>({
  type: '',
  scheme_ids: [],
  period_type: 'monthly',
  start_date: '',
  end_date: '',
  month: new Date().getMonth() + 1,
  year: new Date().getFullYear(),
  template_id: null,
  output_format: 'excel',
  generate_per_scheme: false,
  include_terminated: false,
  include_dependants: true,
  include_beneficiaries: true,
  validate_id_numbers: true,
  validate_banking_details: true,
  exclude_invalid: false,
  reference_number: '',
  notes: ''
})

const generationResult = ref({
  bordereaux_id: '',
  record_count: 0,
  file_size: '',
  file_name: '',
  processing_time: '',
  is_multi_file: false,
  file_list: [] as Array<{ name: string; size: string; scheme_name: string }>,
  download_type: 'single' as 'single' | 'zip' | 'multiple'
})

// Validation rules
const rules = {
  required: (value: any) => !!value || 'This field is required'
}

// Static data
const bordereauTypes = [
  {
    value: 'member',
    label: 'Member Bordereaux',
    description: 'Member enrollment, changes, and terminations',
    icon: 'mdi-account-group',
    color: 'blue'
  },
  {
    value: 'premium',
    label: 'Premium Bordereaux',
    description: 'Premium calculations and collections',
    icon: 'mdi-currency-usd',
    color: 'green'
  },
  {
    value: 'claims',
    label: 'Claims Bordereaux',
    description: 'Claims submissions and settlements',
    icon: 'mdi-medical-bag',
    color: 'orange'
  },
  {
    value: 'reinsurance_premium',
    label: 'Reinsurance Premium Bordereaux',
    description: 'Reinsurance Premium calculations and collections',
    icon: 'mdi-currency-usd',
    color: 'green'
  },
  {
    value: 'reinsurance_claims',
    label: 'Reinsurance Claims Bordereaux',
    description: 'Reinsurance Claims submissions and settlements',
    icon: 'mdi-medical-bag',
    color: 'orange'
  },
  {
    value: 'benefits',
    label: 'Benefits Bordereaux',
    description: 'Benefit payments and adjustments',
    icon: 'mdi-gift',
    color: 'purple'
  }
]

const periodTypes = [
  { title: 'Monthly', value: 'monthly' },
  { title: 'Quarterly', value: 'quarterly' },
  { title: 'Annual', value: 'annual' },
  { title: 'Custom Range', value: 'custom' }
]

const outputFormats = [
  { title: 'Excel (.xlsx)', value: 'excel' },
  { title: 'CSV', value: 'csv' },
  { title: 'PDF Report', value: 'pdf' },
  { title: 'XML', value: 'xml' }
]

const months = [
  { title: 'January', value: 1 },
  { title: 'February', value: 2 },
  { title: 'March', value: 3 },
  { title: 'April', value: 4 },
  { title: 'May', value: 5 },
  { title: 'June', value: 6 },
  { title: 'July', value: 7 },
  { title: 'August', value: 8 },
  { title: 'September', value: 9 },
  { title: 'October', value: 10 },
  { title: 'November', value: 11 },
  { title: 'December', value: 12 }
]

const currentYear = new Date().getFullYear()
const years = Array.from({ length: 10 }, (_, i) => currentYear - 5 + i)

// Mock data (schemes are passed as props)

const availableTemplates = ref([
  {
    id: 1,
    name: 'Liberty Life Standard',
    insurer_name: 'Liberty Life',
    format: 'Excel',
    type: 'member'
  },
  {
    id: 2,
    name: 'Old Mutual Premium Template',
    insurer_name: 'Old Mutual',
    format: 'CSV',
    type: 'premium'
  },
  {
    id: 3,
    name: 'Momentum Claims Format',
    insurer_name: 'Momentum',
    format: 'Excel',
    type: 'claims'
  }
])

// Computed properties
const canGenerate = computed(() => {
  return (
    formData.value.type &&
    formData.value.scheme_ids.length > 0 &&
    formData.value.template_id &&
    formData.value.output_format &&
    (formData.value.period_type !== 'custom' ||
      (formData.value.start_date && formData.value.end_date))
  )
})

const schemeSelectionLabel = computed(() => {
  const count = formData.value.scheme_ids.length
  if (count === 0) {
    return 'Scheme(s) *'
  } else if (count === 1) {
    return 'Scheme(s) * (1 selected)'
  } else {
    return `Scheme(s) * (${count} selected)`
  }
})

const someSchemes = computed(() => formData.value.scheme_ids.length > 0)

const allSchemes = computed(() => {
  return (
    schemes.value.length > 0 &&
    formData.value.scheme_ids.length === schemes.value.length
  )
})

// Limit the number of chips displayed
const maxDisplayedChips = 3

const hiddenSchemesCount = computed(() => {
  const totalSelected = formData.value.scheme_ids.length
  return Math.max(0, totalSelected - maxDisplayedChips)
})

// Methods
const toggle = () => {
  if (allSchemes.value) {
    // Deselect all
    formData.value.scheme_ids = []
  } else {
    // Select all
    formData.value.scheme_ids = schemes.value.map((scheme: any) => scheme.id)
  }
}

const removeScheme = (schemeId: number) => {
  const index = formData.value.scheme_ids.indexOf(schemeId)
  if (index > -1) {
    formData.value.scheme_ids.splice(index, 1)
  }
}

const onTypeChange = () => {
  // Filter templates by type
  // Reset template selection when type changes
  formData.value.template_id = null
  updateEstimatedRecords()
}

const onPeriodTypeChange = () => {
  if (formData.value.period_type !== 'custom') {
    formData.value.start_date = ''
    formData.value.end_date = ''
  }
  updateEstimatedRecords()
}

const updateEstimatedRecords = () => {
  if (schemes.value.length === 0 || formData.value.scheme_ids.length === 0) {
    estimatedRecords.value = 0
    return
  }
  estimatedRecords.value = 0

  const selectedSchemes = schemes.value.filter((scheme) =>
    formData.value.scheme_ids.includes(scheme.id)
  )

  // get estimated records based on selected schemes and bordereaux type

  selectedSchemes.forEach((scheme: any) => {
    console.log(
      `Scheme: ${scheme.name}, Estimated Members: ${scheme.estimated_members}`
    )
    estimatedRecords.value += scheme.member_count || 0
  })
}

const getSelectedType = () => {
  return bordereauTypes.find((type) => type.value === formData.value.type)
}

const getSelectedTemplate = () => {
  return availableTemplates.value.find(
    (template) => template.id === formData.value.template_id
  )
}

const getPeriodSummary = () => {
  if (formData.value.period_type === 'custom') {
    return formData.value.start_date && formData.value.end_date
      ? `${formData.value.start_date} to ${formData.value.end_date}`
      : 'Custom range not set'
  }

  if (formData.value.period_type === 'annual') {
    return `${formData.value.year}`
  }

  const monthName =
    months.find((m) => m.value === formData.value.month)?.title || ''
  return `${monthName} ${formData.value.year}`
}

const generateBordereaux = async () => {
  loading.value = true
  progressDialog.value = true
  progress.value = 0
  progressMessage.value = 'Submitting job to server...'

  try {
    // Submit job to backend API
    console.log(
      'Submitting bordereaux generation job with data:',
      formData.value
    )
    const response = await GroupPricingService.generateBordereaux(
      formData.value
    )
    console.log('Job submission response:', response.data)
    if (response.status === 200) {
      loading.value = false
      progressDialog.value = false
      successDialog.value = true
      jobResult.value = response.data
      generationResult.value = {
        bordereaux_id: jobResult.value.bordereaux_id,
        record_count: jobResult.value.records,
        file_size: jobResult.value.file_size,
        file_name: jobResult.value.file_name,
        processing_time: jobResult.value.processing_time,
        is_multi_file: jobResult.value.is_multi_file || false,
        file_list: jobResult.value.file_list || [],
        download_type: jobResult.value.download_type || 'single'
      }
    } else {
      throw new Error(response.data.error || 'Failed to submit job')
    }
  } catch (error: any) {
    progressDialog.value = false
    loading.value = false
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to submit bordereaux generation job',
      'error'
    )
  }
}

const stopJobPolling = () => {
  if (pollInterval.value) {
    clearInterval(pollInterval.value)
    pollInterval.value = null
  }
}

const downloadFile = async () => {
  if (!jobResult.value) return

  try {
    console.log(
      'Downloading bordereaux file(s) for job ID:',
      jobResult.value.download_url
    )

    if (
      generationResult.value.is_multi_file &&
      generationResult.value.download_type === 'zip'
    ) {
      // Download zip file containing all scheme files
      const response = await GroupPricingService.downloadBordereauxFile(
        jobResult.value.download_url
      )

      const url = window.URL.createObjectURL(response.data)
      const link = document.createElement('a')
      link.href = url
      link.download = `${generationResult.value.file_name || 'bordereaux-files.zip'}`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    } else if (
      generationResult.value.is_multi_file &&
      generationResult.value.download_type === 'multiple'
    ) {
      // Download individual files sequentially
      for (const file of generationResult.value.file_list) {
        const response = await GroupPricingService.downloadBordereauxFile(
          `${jobResult.value.download_url}/${file.name}`
        )

        const url = window.URL.createObjectURL(response.data)
        const link = document.createElement('a')
        link.href = url
        link.download = file.name
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
        window.URL.revokeObjectURL(url)

        // Small delay between downloads to avoid browser blocking
        await new Promise((resolve) => setTimeout(resolve, 500))
      }
    } else {
      // Single file download (original logic)
      const response = await GroupPricingService.downloadBordereauxFile(
        jobResult.value.download_url
      )

      const url = window.URL.createObjectURL(response.data)
      const link = document.createElement('a')
      link.href = url
      link.download = `${generationResult.value.file_name}`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to download bordereaux',
      'error'
    )
  }
}

const viewSubmissionTracking = () => {
  successDialog.value = false
  router.push('/group-pricing/bordereaux-management/tracking')
}

const closeSuccessDialog = () => {
  successDialog.value = false
  // Reset form or navigate back
}

// Configuration Management Methods
const fetchSavedConfigurations = async () => {
  try {
    const response = await GroupPricingService.getBordereauxConfigurations()
    savedConfigurations.value = response.data
  } catch (error) {
    console.error('Failed to fetch saved configurations:', error)
  }
}

const loadConfiguration = (configId: number | null) => {
  if (!configId) {
    selectedConfigurationId.value = null
    return
  }

  const config = savedConfigurations.value.find((c) => c.id === configId)
  if (!config) return

  // Load configuration data into form
  const configData = config.config_data
  Object.keys(formData.value).forEach((key) => {
    if (key in configData) {
      ;(formData.value as any)[key] = (configData as any)[key]
    }
  })

  // Update usage tracking
  updateConfigurationUsage(configId)
  updateEstimatedRecords()
}

const openSaveConfigDialog = () => {
  configName.value = ''
  configDescription.value = ''
  isUpdateMode.value = false
  saveConfigDialog.value = true
}

const openUpdateConfigDialog = () => {
  if (!selectedConfigurationId.value) return

  const config = savedConfigurations.value.find(
    (c) => c.id === selectedConfigurationId.value
  )
  if (!config) return

  configName.value = config.name
  configDescription.value = config.description
  isUpdateMode.value = true
  saveConfigDialog.value = true
}

const closeSaveConfigDialog = () => {
  saveConfigDialog.value = false
  configName.value = ''
  configDescription.value = ''
  isUpdateMode.value = false
}

const saveConfiguration = async () => {
  if (!configName.value.trim()) return

  try {
    const configurationData = {
      name: configName.value.trim(),
      description: configDescription.value.trim(),
      config_data: { ...formData.value }
    }

    console.log('Saving configuration:', configurationData)

    const response =
      await GroupPricingService.saveBordereauxConfiguration(configurationData)

    if (response.status === 201) {
      await fetchSavedConfigurations()
      closeSaveConfigDialog()
      flash.show(`Configuration '${configurationData.name}' saved`, 'success')
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to save configuration',
      'error'
    )
  }
}

const updateConfiguration = async () => {
  if (!configName.value.trim() || !selectedConfigurationId.value) return

  try {
    const configurationData = {
      name: configName.value.trim(),
      description: configDescription.value.trim(),
      config_data: { ...formData.value }
    }

    console.log('Updating configuration:', configurationData)

    const response = await GroupPricingService.updateBordereauxConfiguration(
      selectedConfigurationId.value,
      configurationData
    )

    if (response.status === 200) {
      await fetchSavedConfigurations()
      closeSaveConfigDialog()
      flash.show(`Configuration '${configurationData.name}' updated`, 'success')
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to update configuration',
      'error'
    )
  }
}

const deleteConfiguration = async () => {
  if (!selectedConfigurationId.value) return

  const configName = getSelectedConfigurationName()
  const confirmed = await confirmationDialog.value?.open(
    'Confirm Deletion',
    `Are you sure you want to delete the configuration "${configName}"? This action cannot be undone.`
  )

  if (!confirmed) return

  try {
    const response = await GroupPricingService.deleteBordereauxConfiguration(
      selectedConfigurationId.value
    )

    if (response.status === 200) {
      savedConfigurations.value = savedConfigurations.value.filter(
        (c) => c.id !== selectedConfigurationId.value
      )
      selectedConfigurationId.value = null
      flash.show(`Configuration '${configName}' deleted`, 'success')
    }
  } catch (error: any) {
    flash.show(
      error.response?.data?.error ||
        error.message ||
        'Failed to delete configuration',
      'error'
    )
  }
}

const getSelectedConfigurationName = () => {
  const config = savedConfigurations.value.find(
    (c) => c.id === selectedConfigurationId.value
  )
  return config ? config.name : ''
}

const updateConfigurationUsage = async (configId: number) => {
  try {
    await GroupPricingService.updateConfigurationUsage(configId)
    // Update local usage count
    const config = savedConfigurations.value.find((c) => c.id === configId)
    if (config) {
      config.usage_count += 1
      config.last_used = new Date().toISOString()
    }
  } catch (error) {
    console.error('Failed to update configuration usage:', error)
  }
}

// Watchers
watch(
  () => formData.value.scheme_ids,
  () => {
    // Reset per-scheme option when only one scheme is selected
    if (formData.value.scheme_ids.length <= 1) {
      formData.value.generate_per_scheme = false
    }
    updateEstimatedRecords()
  }
)
watch(() => formData.value.type, updateEstimatedRecords)

onMounted(async () => {
  updateEstimatedRecords()
  // Fetch schemes from API
  try {
    const response = await GroupPricingService.getSchemesInforcev2()
    schemes.value = response.data
    const templatesResponse = await GroupPricingService.getBordereauxTemplates()
    availableTemplates.value = templatesResponse.data

    // Load saved configurations
    await fetchSavedConfigurations()
  } catch (error) {
    console.error('Failed to fetch schemes:', error)
  }
})

// Cleanup on unmount
onUnmounted(() => {
  stopJobPolling()
})
</script>

<style scoped>
.sticky-card {
  position: sticky;
  top: 24px;
}
</style>
