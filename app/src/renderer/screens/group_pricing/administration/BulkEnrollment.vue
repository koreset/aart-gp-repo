<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center w-100">
              <div class="d-flex align-center">
                <v-btn
                  icon="mdi-arrow-left"
                  size="small"
                  variant="text"
                  class="mr-2"
                  @click="goBack"
                />
                <div>
                  <span class="headline">Bulk Member Enrollment</span>
                  <div class="text-caption breadcrumb-on-dark">
                    Member Management → Bulk Enrollment
                  </div>
                </div>
              </div>
              <v-btn
                size="small"
                variant="outlined"
                color="white"
                rounded
                prepend-icon="mdi-download"
                @click="downloadTemplate"
              >
                Download CSV Template
              </v-btn>
            </div>
          </template>
          <template #default>
            <v-form ref="form" @submit.prevent="handleUpload">
              <v-row>
                <v-col cols="12" md="6">
                  <v-autocomplete
                    v-model="selectedScheme"
                    :items="schemes"
                    item-title="name"
                    item-value="id"
                    label="Target Scheme *"
                    variant="outlined"
                    density="compact"
                    :rules="[rules.required]"
                    clearable
                    auto-select-first
                    required
                    :disabled="!!createdBatchId"
                  />
                </v-col>
                <v-col cols="12" md="6">
                  <v-file-input
                    v-model="uploadFile"
                    accept=".csv"
                    label="Member Data CSV File *"
                    variant="outlined"
                    density="compact"
                    prepend-icon="mdi-upload"
                    :rules="[rules.required, rules.fileType]"
                    show-size
                    required
                    :disabled="!!createdBatchId"
                  />
                </v-col>

                <v-col cols="12" md="6">
                  <v-checkbox
                    v-model="skipDuplicates"
                    label="Skip duplicate ID numbers"
                    color="primary"
                    hide-details
                    :disabled="!!createdBatchId"
                  />
                </v-col>
                <v-col cols="12" md="6">
                  <v-checkbox
                    v-model="runExternalCheckIdAfterUpload"
                    label="Run external CheckID validation after upload"
                    color="primary"
                    hide-details
                    :disabled="!!createdBatchId"
                  />
                </v-col>

                <!-- Upload Progress -->
                <v-col v-if="uploadProgress.show" cols="12">
                  <v-card variant="outlined">
                    <v-card-text>
                      <div class="text-subtitle-2 mb-2">Upload Progress</div>
                      <v-progress-linear
                        :model-value="uploadProgress.percentage"
                        color="primary"
                        height="8"
                        class="mb-2"
                      />
                      <div class="text-body-2">
                        {{ uploadProgress.message }}
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>

                <!-- Results Summary -->
                <v-col v-if="createdBatch" cols="12">
                  <v-card variant="outlined">
                    <v-card-title class="text-subtitle-2"
                      >Batch Submitted for Approval</v-card-title
                    >
                    <v-card-text>
                      <v-row>
                        <v-col cols="6" md="3">
                          <div class="text-center">
                            <div class="text-h4 text-success">{{
                              createdBatch.valid_count
                            }}</div>
                            <div class="text-caption">Valid</div>
                          </div>
                        </v-col>
                        <v-col cols="6" md="3">
                          <div class="text-center">
                            <div class="text-h4 text-error">{{
                              createdBatch.blocking_count
                            }}</div>
                            <div class="text-caption">Blocking</div>
                          </div>
                        </v-col>
                        <v-col cols="6" md="3">
                          <div class="text-center">
                            <div class="text-h4 text-warning">{{
                              createdBatch.soft_error_count
                            }}</div>
                            <div class="text-caption">Soft Warnings</div>
                          </div>
                        </v-col>
                        <v-col cols="6" md="3">
                          <div class="text-center">
                            <div class="text-h4 text-info">{{
                              createdBatch.member_count
                            }}</div>
                            <div class="text-caption">Total Rows</div>
                          </div>
                        </v-col>
                      </v-row>

                      <v-alert
                        class="mt-4"
                        type="info"
                        variant="tonal"
                        density="compact"
                      >
                        Members are saved as drafts and will not appear in the
                        live members list until an approver reviews and accepts
                        the batch.
                      </v-alert>

                      <MemberUploadErrorReport
                        v-if="validationReport.length > 0"
                        class="mt-4"
                        :blocking-errors="blockingValidationRows"
                        :soft-errors="softValidationRows"
                        :context-label="selectedSchemeName"
                        filename-prefix="bulk-enrollment-errors"
                      />
                    </v-card-text>
                  </v-card>
                </v-col>

                <!-- Action Buttons -->
                <v-col cols="12">
                  <v-card-actions class="justify-end">
                    <v-btn
                      v-if="!createdBatchId"
                      size="small"
                      rounded
                      color="grey"
                      variant="outlined"
                      :disabled="uploading"
                      @click="goBack"
                      >Cancel</v-btn
                    >
                    <v-btn
                      v-if="!createdBatchId"
                      size="small"
                      rounded
                      color="primary"
                      type="submit"
                      :loading="uploading"
                      :disabled="!isFormValid"
                    >
                      Upload as Draft Batch
                    </v-btn>
                    <v-btn
                      v-if="createdBatchId"
                      size="small"
                      rounded
                      color="primary"
                      @click="viewBatch"
                    >
                      Open Batch Review
                    </v-btn>
                    <v-btn
                      v-if="createdBatchId"
                      size="small"
                      rounded
                      color="info"
                      variant="outlined"
                      @click="goBack"
                    >
                      Back to Members
                    </v-btn>
                  </v-card-actions>
                </v-col>
              </v-row>
            </v-form>

            <v-expansion-panels class="mt-4">
              <v-expansion-panel>
                <v-expansion-panel-title>
                  <v-icon left>mdi-information-outline</v-icon>
                  CSV Format Reference
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <div class="text-body-2">
                    <p
                      >Your CSV file must contain the following columns in this
                      exact order. Benefit multiples become required when the
                      scheme's in-force quote does not have "Use Global Salary
                      Multiple" enabled.</p
                    >
                    <v-table density="compact" class="mt-3">
                      <thead>
                        <tr>
                          <th>Column Name</th>
                          <th>Required</th>
                          <th>Format</th>
                          <th>Example</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr>
                          <td>member_name</td>
                          <td
                            ><v-chip size="x-small" color="error"
                              >Required</v-chip
                            ></td
                          >
                          <td>Text</td>
                          <td>John Doe</td>
                        </tr>
                        <tr>
                          <td>member_id_number</td>
                          <td
                            ><v-chip size="x-small" color="error"
                              >Required</v-chip
                            ></td
                          >
                          <td>13 digits for RSA ID</td>
                          <td>8001015009087</td>
                        </tr>
                        <tr>
                          <td>member_id_type</td>
                          <td
                            ><v-chip size="x-small" color="error"
                              >Required</v-chip
                            ></td
                          >
                          <td>RSA_ID, PASSPORT, OTHER</td>
                          <td>RSA_ID</td>
                        </tr>
                        <tr>
                          <td>gender</td>
                          <td
                            ><v-chip size="x-small" color="error"
                              >Required</v-chip
                            ></td
                          >
                          <td>Male / Female</td>
                          <td>Male</td>
                        </tr>
                        <tr>
                          <td>date_of_birth</td>
                          <td
                            ><v-chip size="x-small" color="error"
                              >Required</v-chip
                            ></td
                          >
                          <td>YYYY-MM-DD</td>
                          <td>1980-12-25</td>
                        </tr>
                        <tr>
                          <td>annual_salary</td>
                          <td
                            ><v-chip size="x-small" color="error"
                              >Required</v-chip
                            ></td
                          >
                          <td>Positive number</td>
                          <td>500000</td>
                        </tr>
                        <tr>
                          <td>entry_date</td>
                          <td
                            ><v-chip size="x-small" color="error"
                              >Required</v-chip
                            ></td
                          >
                          <td>YYYY-MM-DD</td>
                          <td>2024-01-15</td>
                        </tr>
                        <tr>
                          <td>scheme_category</td>
                          <td
                            ><v-chip size="x-small" color="error"
                              >Required</v-chip
                            ></td
                          >
                          <td>Must match a category on the in-force quote</td>
                          <td>General</td>
                        </tr>
                        <tr>
                          <td>occupation</td>
                          <td
                            ><v-chip size="x-small" color="error"
                              >Required</v-chip
                            ></td
                          >
                          <td>Must match an Occupation Class entry</td>
                          <td>Software Developer</td>
                        </tr>
                        <tr>
                          <td>email</td>
                          <td
                            ><v-chip size="x-small" color="warning"
                              >Optional</v-chip
                            ></td
                          >
                          <td>Email address</td>
                          <td>john@company.com</td>
                        </tr>
                        <tr>
                          <td>phone_number</td>
                          <td
                            ><v-chip size="x-small" color="warning"
                              >Optional</v-chip
                            ></td
                          >
                          <td>Digits</td>
                          <td>0821234567</td>
                        </tr>
                        <tr>
                          <td>employee_number</td>
                          <td
                            ><v-chip size="x-small" color="warning"
                              >Optional</v-chip
                            ></td
                          >
                          <td>Text</td>
                          <td>EMP001</td>
                        </tr>
                        <tr>
                          <td>benefits_gla_multiple</td>
                          <td
                            ><v-chip size="x-small" color="warning"
                              >Conditional</v-chip
                            ></td
                          >
                          <td>Number (decimal)</td>
                          <td>5.0</td>
                        </tr>
                        <tr>
                          <td>benefits_sgla_multiple</td>
                          <td
                            ><v-chip size="x-small" color="warning"
                              >Conditional</v-chip
                            ></td
                          >
                          <td>Number (decimal)</td>
                          <td>4.5</td>
                        </tr>
                        <tr>
                          <td>benefits_ptd_multiple</td>
                          <td
                            ><v-chip size="x-small" color="warning"
                              >Conditional</v-chip
                            ></td
                          >
                          <td>Number (decimal)</td>
                          <td>4.0</td>
                        </tr>
                        <tr>
                          <td>benefits_ci_multiple</td>
                          <td
                            ><v-chip size="x-small" color="warning"
                              >Conditional</v-chip
                            ></td
                          >
                          <td>Number (decimal)</td>
                          <td>3.0</td>
                        </tr>
                        <tr>
                          <td>benefits_ttd_multiple</td>
                          <td
                            ><v-chip size="x-small" color="warning"
                              >Conditional</v-chip
                            ></td
                          >
                          <td>Number (decimal)</td>
                          <td>0.56</td>
                        </tr>
                        <tr>
                          <td>benefits_phi_multiple</td>
                          <td
                            ><v-chip size="x-small" color="warning"
                              >Conditional</v-chip
                            ></td
                          >
                          <td>Number (decimal)</td>
                          <td>0.7</td>
                        </tr>
                      </tbody>
                    </v-table>
                  </div>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :timeout="4000" :color="snackbarColor">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Papa from 'papaparse'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import MemberUploadErrorReport from '@/renderer/screens/group_pricing/shared/MemberUploadErrorReport.vue'

interface ValidationRow {
  row: number
  field?: string
  severity: string
  message: string
  member_id_number?: string
  member_name?: string
}

const router = useRouter()
const form = ref<any>(null)
const uploading = ref(false)
const selectedScheme = ref<number | null>(null)
const uploadFile = ref<File | null>(null)
const skipDuplicates = ref(true)
const runExternalCheckIdAfterUpload = ref(false)

const schemes = ref<Array<{ id: number; name: string }>>([])
const createdBatch = ref<any>(null)
const createdBatchId = computed(() => createdBatch.value?.id ?? null)
const validationReport = ref<ValidationRow[]>([])

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref<'success' | 'error' | 'info'>('success')

const selectedSchemeName = computed(() => {
  const s = schemes.value.find((x) => x.id === selectedScheme.value)
  return s?.name ?? ''
})

const blockingValidationRows = computed(() =>
  validationReport.value.filter((r) => r.severity === 'blocking')
)
const softValidationRows = computed(() =>
  validationReport.value
    .filter((r) => r.severity === 'soft')
    .map((r) => ({ row: r.row, message: r.message }))
)

const uploadProgress = ref({ show: false, percentage: 0, message: '' })

const rules = {
  required: (value: any) => !!value || 'Field is required',
  fileType: (value: any) => {
    if (!value) return true
    return value?.type === 'text/csv' || 'File must be CSV format'
  }
}

const isFormValid = computed(
  () => !!selectedScheme.value && uploadFile.value instanceof File
)

onMounted(async () => {
  try {
    const res = await GroupPricingService.getSchemesInforcev2()
    schemes.value = res?.data ?? []
  } catch (err) {
    console.error('Failed to load schemes', err)
  }
})

const goBack = () => {
  router.push({ name: 'group-pricing-member-management' })
}

const viewBatch = () => {
  if (!createdBatchId.value) return
  router.push({
    name: 'group-pricing-bulk-enrollment-batch',
    params: { batchId: String(createdBatchId.value) }
  })
}

const showSnack = (
  msg: string,
  color: 'success' | 'error' | 'info' = 'success'
) => {
  snackbarMessage.value = msg
  snackbarColor.value = color
  snackbar.value = true
}

// Reads the file once with FileReader. Reading the File handle a second time
// (e.g. via file.arrayBuffer() after Papa.parse(file)) trips Electron's "file
// could not be read" permission error, so callers must reuse this text.
const readFileAsText = (file: File): Promise<string> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result ?? ''))
    reader.onerror = () =>
      reject(reader.error ?? new Error('Failed to read file'))
    reader.readAsText(file)
  })

const parseCSVText = (text: string): any[] => {
  const result = Papa.parse<any>(text, {
    header: true,
    skipEmptyLines: true
  })
  if (result.errors && result.errors.length > 0) {
    throw new Error('CSV parsing errors: ' + result.errors[0].message)
  }
  return result.data
}

const toMemberPayload = (row: any) => {
  const num = (v: any) =>
    v === undefined || v === null || String(v).trim() === '' ? 0 : Number(v)
  return {
    member_name: row.member_name,
    member_id_number: row.member_id_number,
    member_id_type: row.member_id_type,
    scheme_category: row.scheme_category || '',
    gender: row.gender,
    date_of_birth: row.date_of_birth,
    annual_salary: num(row.annual_salary),
    email: row.email || '',
    phone_number: row.phone_number || '',
    employee_number: row.employee_number || '',
    occupation: row.occupation || '',
    occupational_class: row.occupational_class || '',
    entry_date: row.entry_date,
    benefits: {
      gla_multiple: num(row.benefits_gla_multiple),
      sgla_multiple: num(row.benefits_sgla_multiple),
      ptd_multiple: num(row.benefits_ptd_multiple),
      ci_multiple: num(row.benefits_ci_multiple),
      ttd_multiple: num(row.benefits_ttd_multiple),
      phi_multiple: num(row.benefits_phi_multiple)
    }
  }
}

const sha256HexOfText = async (text: string) => {
  try {
    const buf = new TextEncoder().encode(text)
    const hash = await crypto.subtle.digest('SHA-256', buf)
    return Array.from(new Uint8Array(hash))
      .map((b) => b.toString(16).padStart(2, '0'))
      .join('')
  } catch {
    return ''
  }
}

const handleUpload = async () => {
  if (!form.value) return
  const { valid } = await (form.value as any).validate()
  if (!valid) return

  uploading.value = true
  uploadProgress.value.show = true
  uploadProgress.value.percentage = 0
  uploadProgress.value.message = 'Reading file...'

  try {
    if (!uploadFile.value) throw new Error('No file selected')
    const file = uploadFile.value

    const text = await readFileAsText(file)
    const csvRows = parseCSVText(text)
    uploadProgress.value.percentage = 30
    uploadProgress.value.message = 'Uploading members as draft batch...'

    const checksum = await sha256HexOfText(text)
    const members = csvRows.map(toMemberPayload)

    const res = await GroupPricingService.createBulkEnrollmentBatch(
      selectedScheme.value,
      members,
      {
        skipDuplicates: skipDuplicates.value,
        fileName: file.name,
        fileSizeBytes: file.size,
        fileChecksum: checksum
      }
    )

    createdBatch.value = res?.data?.batch ?? null
    validationReport.value = res?.data?.validation_report ?? []
    uploadProgress.value.percentage = 70
    uploadProgress.value.message = 'Batch saved'

    if (runExternalCheckIdAfterUpload.value && createdBatchId.value) {
      uploadProgress.value.message = 'Running external CheckID validation...'
      const ext = await GroupPricingService.runExternalIdCheckOnBatch(
        createdBatchId.value
      )
      createdBatch.value = ext?.data?.batch ?? createdBatch.value
      validationReport.value =
        ext?.data?.validation_report ?? validationReport.value
    }

    uploadProgress.value.percentage = 100
    uploadProgress.value.message = 'Batch submitted for approval'
    showSnack('Batch submitted for approval', 'success')
  } catch (error: any) {
    console.error('Upload error:', error)
    const msg =
      error?.response?.data?.error || error?.message || 'Upload failed'
    uploadProgress.value.percentage = 100
    uploadProgress.value.message = msg
    showSnack(msg, 'error')
  } finally {
    uploading.value = false
  }
}

const downloadTemplate = () => {
  const templateData = [
    {
      member_name: 'John Doe',
      member_id_number: '8001015009087',
      member_id_type: 'RSA_ID',
      gender: 'Male',
      date_of_birth: '1980-12-25',
      annual_salary: 500000,
      entry_date: '2024-01-15',
      scheme_category: 'General',
      email: 'john@company.com',
      phone_number: '0821234567',
      employee_number: 'EMP001',
      occupation: 'Software Developer',
      benefits_gla_multiple: 5.0,
      benefits_sgla_multiple: 4.5,
      benefits_ptd_multiple: 4.0,
      benefits_ci_multiple: 3.0,
      benefits_ttd_multiple: 0.56,
      benefits_phi_multiple: 0.7
    }
  ]
  const csv = Papa.unparse(templateData, {
    quotes: true,
    quoteChar: '"',
    escapeChar: '"',
    delimiter: ',',
    header: true
  })
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', 'bulk_enrollment_template.csv')
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
</script>

<style scoped>
.headline {
  font-size: 1.25rem;
  font-weight: 500;
}

.breadcrumb-on-dark {
  color: rgba(255, 255, 255, 0.75);
}
</style>
