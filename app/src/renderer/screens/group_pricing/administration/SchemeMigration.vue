<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12">
        <v-card>
          <v-card-title
            class="d-flex align-center"
            style="background-color: #223f54; color: white"
          >
            <v-icon class="mr-2" color="white"
              >mdi-database-import-outline</v-icon
            >
            Scheme Migration
          </v-card-title>

          <v-card-text class="pa-4">
            <v-alert type="info" variant="tonal" class="mb-4">
              <div class="text-subtitle-2">Bulk Scheme & Member Migration</div>
              <div class="text-body-2 mt-1">
                Upload CSV files to create schemes, categories, members,
                beneficiaries, and claims experience in a single atomic
                operation. Steps 1-3 are required. Steps 4-5 are optional.
              </div>
            </v-alert>

            <v-stepper v-model="currentStep" alt-labels>
              <v-stepper-header>
                <v-stepper-item title="Scheme Setup" :value="1" />
                <v-divider />
                <v-stepper-item title="Categories" :value="2" />
                <v-divider />
                <v-stepper-item title="Member Data" :value="3" />
                <v-divider />
                <v-stepper-item title="Beneficiaries" :value="4" />
                <v-divider />
                <v-stepper-item title="Claims" :value="5" />
                <v-divider />
                <v-stepper-item title="Execute" :value="6" />
              </v-stepper-header>

              <v-stepper-window>
                <!-- Step 1: Scheme Setup -->
                <v-stepper-window-item :value="1">
                  <v-card flat class="pa-4">
                    <div class="d-flex align-center mb-3">
                      <div class="text-h6">Step 1: Scheme Setup</div>
                      <v-spacer />
                      <v-btn
                        size="small"
                        variant="outlined"
                        prepend-icon="mdi-download"
                        @click="downloadTemplate('scheme_setup')"
                      >
                        Download Template
                      </v-btn>
                    </div>
                    <v-file-input
                      v-model="fileSchemeSetup"
                      accept=".csv"
                      label="Scheme Setup CSV *"
                      variant="outlined"
                      density="compact"
                      prepend-icon="mdi-file-delimited"
                      @update:model-value="onFileChange('scheme_setup', $event)"
                    />
                    <div v-if="previews.scheme_setup.length > 0" class="mt-3">
                      <div class="text-caption text-grey mb-1"
                        >Scheme Setup preview (first 10 rows of
                        {{ previews.scheme_setup.length }} total)</div
                      >
                      <v-table
                        density="compact"
                        class="elevation-1"
                        style="max-height: 300px; overflow: auto"
                      >
                        <thead>
                          <tr>
                            <th
                              v-for="col in Object.keys(
                                previews.scheme_setup[0]
                              )"
                              :key="col"
                              class="text-caption"
                              >{{ col }}</th
                            >
                          </tr>
                        </thead>
                        <tbody>
                          <tr
                            v-for="(row, i) in previews.scheme_setup.slice(
                              0,
                              10
                            )"
                            :key="i"
                          >
                            <td
                              v-for="col in Object.keys(
                                previews.scheme_setup[0]
                              )"
                              :key="col"
                              class="text-caption"
                              >{{ row[col] }}</td
                            >
                          </tr>
                        </tbody>
                      </v-table>
                    </div>
                  </v-card>
                </v-stepper-window-item>

                <!-- Step 2: Categories -->
                <v-stepper-window-item :value="2">
                  <v-card flat class="pa-4">
                    <div class="d-flex align-center mb-3">
                      <div class="text-h6">Step 2: Scheme Categories</div>
                      <v-spacer />
                      <v-btn
                        size="small"
                        variant="outlined"
                        prepend-icon="mdi-download"
                        @click="downloadTemplate('scheme_categories')"
                      >
                        Download Template
                      </v-btn>
                    </div>
                    <v-file-input
                      v-model="fileCategories"
                      accept=".csv"
                      label="Scheme Categories CSV *"
                      variant="outlined"
                      density="compact"
                      prepend-icon="mdi-file-delimited"
                      @update:model-value="onFileChange('categories', $event)"
                    />
                    <div v-if="previews.categories.length > 0" class="mt-3">
                      <div class="text-caption text-grey mb-1"
                        >Categories preview (first 10 rows of
                        {{ previews.categories.length }} total)</div
                      >
                      <v-table
                        density="compact"
                        class="elevation-1"
                        style="max-height: 300px; overflow: auto"
                      >
                        <thead>
                          <tr>
                            <th
                              v-for="col in Object.keys(previews.categories[0])"
                              :key="col"
                              class="text-caption"
                              >{{ col }}</th
                            >
                          </tr>
                        </thead>
                        <tbody>
                          <tr
                            v-for="(row, i) in previews.categories.slice(0, 10)"
                            :key="i"
                          >
                            <td
                              v-for="col in Object.keys(previews.categories[0])"
                              :key="col"
                              class="text-caption"
                              >{{ row[col] }}</td
                            >
                          </tr>
                        </tbody>
                      </v-table>
                    </div>
                  </v-card>
                </v-stepper-window-item>

                <!-- Step 3: Member Data -->
                <v-stepper-window-item :value="3">
                  <v-card flat class="pa-4">
                    <div class="d-flex align-center mb-3">
                      <div class="text-h6">Step 3: Member Data</div>
                      <v-spacer />
                      <v-btn
                        size="small"
                        variant="outlined"
                        prepend-icon="mdi-download"
                        @click="downloadTemplate('member_data')"
                      >
                        Download Template
                      </v-btn>
                    </div>
                    <v-file-input
                      v-model="fileMemberData"
                      accept=".csv"
                      label="Member Data CSV *"
                      variant="outlined"
                      density="compact"
                      prepend-icon="mdi-file-delimited"
                      @update:model-value="onFileChange('member_data', $event)"
                    />
                    <div v-if="previews.member_data.length > 0" class="mt-3">
                      <div class="text-caption text-grey mb-1"
                        >Members preview (first 10 rows of
                        {{ previews.member_data.length }} total)</div
                      >
                      <v-table
                        density="compact"
                        class="elevation-1"
                        style="max-height: 300px; overflow: auto"
                      >
                        <thead>
                          <tr>
                            <th
                              v-for="col in Object.keys(
                                previews.member_data[0]
                              )"
                              :key="col"
                              class="text-caption"
                              >{{ col }}</th
                            >
                          </tr>
                        </thead>
                        <tbody>
                          <tr
                            v-for="(row, i) in previews.member_data.slice(
                              0,
                              10
                            )"
                            :key="i"
                          >
                            <td
                              v-for="col in Object.keys(
                                previews.member_data[0]
                              )"
                              :key="col"
                              class="text-caption"
                              >{{ row[col] }}</td
                            >
                          </tr>
                        </tbody>
                      </v-table>
                    </div>
                  </v-card>
                </v-stepper-window-item>

                <!-- Step 4: Beneficiaries (optional) -->
                <v-stepper-window-item :value="4">
                  <v-card flat class="pa-4">
                    <div class="d-flex align-center mb-3">
                      <div class="text-h6"
                        >Step 4: Beneficiaries (Optional)</div
                      >
                      <v-spacer />
                      <v-btn
                        size="small"
                        variant="outlined"
                        prepend-icon="mdi-download"
                        @click="downloadTemplate('beneficiaries')"
                      >
                        Download Template
                      </v-btn>
                    </div>
                    <v-file-input
                      v-model="fileBeneficiaries"
                      accept=".csv"
                      label="Beneficiaries CSV"
                      variant="outlined"
                      density="compact"
                      prepend-icon="mdi-file-delimited"
                      @update:model-value="
                        onFileChange('beneficiaries', $event)
                      "
                    />
                    <div v-if="previews.beneficiaries.length > 0" class="mt-3">
                      <div class="text-caption text-grey mb-1"
                        >Beneficiaries preview (first 10 rows of
                        {{ previews.beneficiaries.length }} total)</div
                      >
                      <v-table
                        density="compact"
                        class="elevation-1"
                        style="max-height: 300px; overflow: auto"
                      >
                        <thead>
                          <tr>
                            <th
                              v-for="col in Object.keys(
                                previews.beneficiaries[0]
                              )"
                              :key="col"
                              class="text-caption"
                              >{{ col }}</th
                            >
                          </tr>
                        </thead>
                        <tbody>
                          <tr
                            v-for="(row, i) in previews.beneficiaries.slice(
                              0,
                              10
                            )"
                            :key="i"
                          >
                            <td
                              v-for="col in Object.keys(
                                previews.beneficiaries[0]
                              )"
                              :key="col"
                              class="text-caption"
                              >{{ row[col] }}</td
                            >
                          </tr>
                        </tbody>
                      </v-table>
                    </div>
                  </v-card>
                </v-stepper-window-item>

                <!-- Step 5: Claims Experience (optional) -->
                <v-stepper-window-item :value="5">
                  <v-card flat class="pa-4">
                    <div class="d-flex align-center mb-3">
                      <div class="text-h6"
                        >Step 5: Claims Experience (Optional)</div
                      >
                      <v-spacer />
                      <v-btn
                        size="small"
                        variant="outlined"
                        prepend-icon="mdi-download"
                        @click="downloadTemplate('claims_experience')"
                      >
                        Download Template
                      </v-btn>
                    </div>
                    <v-file-input
                      v-model="fileClaimsExperience"
                      accept=".csv"
                      label="Claims Experience CSV"
                      variant="outlined"
                      density="compact"
                      prepend-icon="mdi-file-delimited"
                      @update:model-value="
                        onFileChange('claims_experience', $event)
                      "
                    />
                    <div
                      v-if="previews.claims_experience.length > 0"
                      class="mt-3"
                    >
                      <div class="text-caption text-grey mb-1"
                        >Claims Experience preview (first 10 rows of
                        {{ previews.claims_experience.length }} total)</div
                      >
                      <v-table
                        density="compact"
                        class="elevation-1"
                        style="max-height: 300px; overflow: auto"
                      >
                        <thead>
                          <tr>
                            <th
                              v-for="col in Object.keys(
                                previews.claims_experience[0]
                              )"
                              :key="col"
                              class="text-caption"
                              >{{ col }}</th
                            >
                          </tr>
                        </thead>
                        <tbody>
                          <tr
                            v-for="(row, i) in previews.claims_experience.slice(
                              0,
                              10
                            )"
                            :key="i"
                          >
                            <td
                              v-for="col in Object.keys(
                                previews.claims_experience[0]
                              )"
                              :key="col"
                              class="text-caption"
                              >{{ row[col] }}</td
                            >
                          </tr>
                        </tbody>
                      </v-table>
                    </div>
                  </v-card>
                </v-stepper-window-item>

                <!-- Step 6: Validate & Execute -->
                <v-stepper-window-item :value="6">
                  <v-card flat class="pa-4">
                    <div class="text-h6 mb-3">Step 6: Validate & Execute</div>

                    <!-- Summary -->
                    <v-row class="mb-4">
                      <v-col
                        v-for="item in summaryItems"
                        :key="item.label"
                        cols="6"
                        md="2"
                      >
                        <v-card variant="tonal" class="pa-3 text-center">
                          <div class="text-h5">{{ item.count }}</div>
                          <div class="text-caption">{{ item.label }}</div>
                        </v-card>
                      </v-col>
                    </v-row>

                    <!-- Validation -->
                    <div class="d-flex align-center mb-3">
                      <v-btn
                        color="primary"
                        variant="elevated"
                        prepend-icon="mdi-check-circle"
                        :loading="validating"
                        :disabled="!canValidate"
                        @click="validate"
                      >
                        Validate All
                      </v-btn>
                      <v-chip
                        v-if="validationResult && validationResult.valid"
                        color="success"
                        class="ml-3"
                        variant="elevated"
                      >
                        <v-icon start>mdi-check</v-icon>
                        Validation Passed
                      </v-chip>
                      <v-chip
                        v-else-if="validationResult && !validationResult.valid"
                        color="error"
                        class="ml-3"
                        variant="elevated"
                      >
                        <v-icon start>mdi-close</v-icon>
                        {{ validationResult.errors.length }} Error(s)
                      </v-chip>
                    </div>

                    <!-- Validation Errors -->
                    <v-expansion-panels
                      v-if="
                        validationResult &&
                        validationResult.errors &&
                        validationResult.errors.length > 0
                      "
                      class="mb-4"
                    >
                      <v-expansion-panel>
                        <v-expansion-panel-title>
                          <v-icon start color="error">mdi-alert-circle</v-icon>
                          Validation Errors ({{
                            validationResult.errors.length
                          }})
                        </v-expansion-panel-title>
                        <v-expansion-panel-text>
                          <v-table density="compact">
                            <thead>
                              <tr>
                                <th>Template</th>
                                <th>Row</th>
                                <th>Column</th>
                                <th>Message</th>
                              </tr>
                            </thead>
                            <tbody>
                              <tr
                                v-for="(error, idx) in validationResult.errors"
                                :key="idx"
                              >
                                <td>{{ error.template }}</td>
                                <td>{{ error.row || '-' }}</td>
                                <td>{{ error.column || '-' }}</td>
                                <td>{{ error.message }}</td>
                              </tr>
                            </tbody>
                          </v-table>
                        </v-expansion-panel-text>
                      </v-expansion-panel>
                    </v-expansion-panels>

                    <!-- Execute -->
                    <v-btn
                      color="success"
                      variant="elevated"
                      size="large"
                      prepend-icon="mdi-rocket-launch"
                      :loading="executing"
                      :disabled="!canExecute"
                      @click="execute"
                    >
                      Execute Migration
                    </v-btn>
                  </v-card>
                </v-stepper-window-item>
              </v-stepper-window>

              <v-stepper-actions
                :disabled="
                  currentStep === 1
                    ? 'prev'
                    : currentStep === 6
                      ? 'next'
                      : undefined
                "
                @click:prev="currentStep--"
                @click:next="currentStep++"
              />
            </v-stepper>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Success Dialog -->
    <v-dialog v-model="showSuccessDialog" max-width="600" persistent>
      <v-card>
        <v-card-title
          class="text-h6"
          style="background-color: #223f54; color: white"
        >
          <v-icon class="mr-2" color="white">mdi-check-circle</v-icon>
          Migration Complete
        </v-card-title>
        <v-card-text class="pa-4">
          <v-alert type="success" variant="tonal" class="mb-4">
            All schemes have been successfully migrated and are now in-force.
          </v-alert>

          <div v-if="executionResult">
            <v-row>
              <v-col cols="6" md="4">
                <div class="text-h5 text-center">{{
                  executionResult.scheme_count
                }}</div>
                <div class="text-caption text-center">Schemes Created</div>
              </v-col>
              <v-col cols="6" md="4">
                <div class="text-h5 text-center">{{
                  executionResult.member_count
                }}</div>
                <div class="text-caption text-center">Members Enrolled</div>
              </v-col>
              <v-col cols="6" md="4">
                <div class="text-h5 text-center">{{
                  executionResult.beneficiary_count
                }}</div>
                <div class="text-caption text-center">Beneficiaries Added</div>
              </v-col>
            </v-row>

            <div
              v-if="
                executionResult.created_scheme_ids &&
                executionResult.created_scheme_ids.length
              "
              class="mt-4"
            >
              <div class="text-subtitle-2 mb-2">Created Scheme IDs:</div>
              <v-chip
                v-for="id in executionResult.created_scheme_ids"
                :key="id"
                class="mr-2 mb-1"
                size="small"
                color="primary"
              >
                {{ id }}
              </v-chip>
            </div>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" variant="elevated" @click="goToSchemes">
            View Schemes
          </v-btn>
          <v-btn variant="text" @click="resetMigration"> New Migration </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script lang="ts">
// Module-level store — persists across stepper mount/unmount cycles
import { reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import Papa from 'papaparse'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const migrationStore = reactive<{
  files: Record<string, File | null>
  data: Record<string, any[]>
}>({
  files: {
    scheme_setup: null,
    categories: null,
    member_data: null,
    beneficiaries: null,
    claims_experience: null
  },
  data: {
    scheme_setup: [],
    categories: [],
    member_data: [],
    beneficiaries: [],
    claims_experience: []
  }
})

function handleFileSelected(key: string, fileInput: File | File[] | null) {
  // Vuetify v-file-input emits File (single) or File[] (multiple) or null
  const file =
    fileInput instanceof File
      ? fileInput
      : Array.isArray(fileInput) &&
          fileInput.length > 0 &&
          fileInput[0] instanceof File
        ? fileInput[0]
        : null
  if (!file) {
    return
  }
  migrationStore.files[key] = file
  const reader = new FileReader()
  reader.onload = (e) => {
    const text = e.target?.result as string
    if (!text) {
      migrationStore.data[key] = []
      return
    }
    Papa.parse(text, {
      header: true,
      skipEmptyLines: true,
      complete: (results: any) => {
        migrationStore.data[key] = results.data || []
      },
      error: () => {
        migrationStore.data[key] = []
      }
    })
  }
  reader.onerror = () => {
    migrationStore.data[key] = []
  }
  reader.readAsText(file)
}

function resetStore() {
  for (const key of Object.keys(migrationStore.files)) {
    migrationStore.files[key] = null
    migrationStore.data[key] = []
  }
}
</script>

<script setup lang="ts">
const router = useRouter()
const currentStep = ref(1)

const fileSchemeSetup = ref<File[] | null>(null)
const fileCategories = ref<File[] | null>(null)
const fileMemberData = ref<File[] | null>(null)
const fileBeneficiaries = ref<File[] | null>(null)
const fileClaimsExperience = ref<File[] | null>(null)

const previews = computed(() => migrationStore.data)

const validating = ref(false)
const executing = ref(false)
const validationResult = ref<any>(null)
const executionResult = ref<any>(null)
const showSuccessDialog = ref(false)

const summaryItems = computed(() => [
  { label: 'Schemes', count: migrationStore.data.scheme_setup.length },
  { label: 'Categories', count: migrationStore.data.categories.length },
  { label: 'Members', count: migrationStore.data.member_data.length },
  { label: 'Beneficiaries', count: migrationStore.data.beneficiaries.length },
  {
    label: 'Claims Rows',
    count: migrationStore.data.claims_experience.length
  }
])

const canValidate = computed(() => {
  return (
    migrationStore.data.scheme_setup.length > 0 &&
    migrationStore.data.categories.length > 0 &&
    migrationStore.data.member_data.length > 0
  )
})

const canExecute = computed(() => {
  return (
    validationResult.value && validationResult.value.valid && !executing.value
  )
})

function onFileChange(key: string, fileInput: File | File[] | null) {
  handleFileSelected(key, fileInput)
  validationResult.value = null
}

function buildFormData(): FormData {
  const fd = new FormData()
  const keys = [
    'scheme_setup',
    'categories',
    'member_data',
    'beneficiaries',
    'claims_experience'
  ]
  for (const key of keys) {
    const file = migrationStore.files[key]
    if (file) {
      fd.append(key, file)
    }
  }
  return fd
}

async function validate() {
  validating.value = true
  validationResult.value = null
  try {
    const { data } =
      await GroupPricingService.validateMigration(buildFormData())
    validationResult.value = data
  } catch (err: any) {
    validationResult.value = {
      valid: false,
      errors: [
        {
          template: 'server',
          row: 0,
          column: '',
          message:
            err?.data?.error || err?.message || 'Validation request failed'
        }
      ]
    }
  } finally {
    validating.value = false
  }
}

async function execute() {
  executing.value = true
  try {
    const { data } = await GroupPricingService.executeMigration(buildFormData())
    executionResult.value = data
    showSuccessDialog.value = true
  } catch (err: any) {
    validationResult.value = {
      valid: false,
      errors: [
        {
          template: 'server',
          row: 0,
          column: '',
          message:
            err?.data?.error || err?.message || 'Migration execution failed'
        }
      ]
    }
  } finally {
    executing.value = false
  }
}

async function downloadTemplate(name: string) {
  try {
    const response = await GroupPricingService.downloadMigrationTemplate(name)
    const blob = new Blob([response.data], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = name + '.csv'
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    window.URL.revokeObjectURL(url)
  } catch (err) {
    console.error('Failed to download template:', err)
  }
}

function goToSchemes() {
  showSuccessDialog.value = false
  router.push({ name: 'group-pricing-schemes' })
}

function resetMigration() {
  showSuccessDialog.value = false
  currentStep.value = 1
  fileSchemeSetup.value = null
  fileCategories.value = null
  fileMemberData.value = null
  fileBeneficiaries.value = null
  fileClaimsExperience.value = null
  resetStore()
  validationResult.value = null
  executionResult.value = null
}
</script>
