<template>
  <v-form ref="form" @submit.prevent="handleUpload">
    <v-row>
      <v-col cols="12">
        <v-alert type="info" variant="tonal" class="mb-4">
          <div class="text-subtitle-2">Bulk Claims Upload</div>
          <div class="text-body-2 mt-2">
            Upload a CSV file containing claims data to register multiple claims
            at once. All claims will be processed and validated before
            submission.
          </div>
        </v-alert>
      </v-col>

      <v-col cols="12">
        <v-file-input
          v-model="uploadFile"
          accept=".csv"
          label="Claims Data CSV File *"
          variant="outlined"
          density="compact"
          prepend-icon="mdi-upload"
          :rules="[rules.required, rules.fileType]"
          :loading="processing"
          required
          @change="handleFileSelect"
        />
      </v-col>

      <!-- File Processing Status -->
      <v-col v-if="processing" cols="12">
        <v-card variant="outlined" class="pa-4">
          <div class="d-flex align-center">
            <v-progress-circular indeterminate size="24" class="mr-3" />
            <span>Processing CSV file...</span>
          </div>
        </v-card>
      </v-col>

      <!-- Preview Data -->
      <v-col v-if="previewData.length > 0 && !processing" cols="12">
        <v-card variant="outlined">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Preview Data ({{ previewData.length }} claims found)
          </v-card-title>
          <v-card-text class="pa-0">
            <div class="table-container">
              <v-data-table
                :headers="previewHeaders"
                :items="previewData.slice(0, 10)"
                density="compact"
                :items-per-page="10"
                hide-default-footer
              />
            </div>
            <div
              v-if="previewData.length > 10"
              class="pa-3 text-caption text-center"
            >
              Showing first 10 rows of {{ previewData.length }} total claims
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Validation Errors -->
      <v-col v-if="validationErrors.length > 0" cols="12">
        <v-card color="error" variant="tonal">
          <v-card-title class="text-subtitle-1">
            <v-icon left>mdi-alert</v-icon>
            Validation Errors ({{ validationErrors.length }})
          </v-card-title>
          <v-card-text>
            <div class="max-height-200 overflow-y-auto">
              <div
                v-for="(error, index) in validationErrors.slice(0, 20)"
                :key="index"
                class="text-body-2 mb-1"
              >
                Row {{ error.row }}: {{ error.message }}
              </div>
              <div
                v-if="validationErrors.length > 20"
                class="text-caption mt-2"
              >
                ... and {{ validationErrors.length - 20 }} more errors
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12">
        <v-expansion-panels>
          <v-expansion-panel>
            <v-expansion-panel-title>
              <v-icon left>mdi-information-outline</v-icon>
              CSV Format Requirements
            </v-expansion-panel-title>
            <v-expansion-panel-text>
              <div class="text-body-2">
                <p
                  >Your CSV file must contain the following columns in this
                  exact order:</p
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
                      <td>claim_number</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Text (auto-generated if empty)</td>
                      <td>CLM-2024-001234</td>
                    </tr>
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
                      <td>Text</td>
                      <td>8001015009087</td>
                    </tr>
                    <tr>
                      <td>member_type</td>
                      <td
                        ><v-chip size="x-small" color="error"
                          >Required</v-chip
                        ></td
                      >
                      <td>Text</td>
                      <td>Member</td>
                    </tr>

                    <tr>
                      <td>scheme_name</td>
                      <td
                        ><v-chip size="x-small" color="error"
                          >Required</v-chip
                        ></td
                      >
                      <td>Text (must match existing scheme)</td>
                      <td>ABC Company Group Scheme</td>
                    </tr>
                    <tr>
                      <td>benefit_type</td>
                      <td
                        ><v-chip size="x-small" color="error"
                          >Required</v-chip
                        ></td
                      >
                      <td>Predefined values</td>
                      <td>Group Life Assurance (GLA)</td>
                    </tr>
                    <tr>
                      <td>claim_amount</td>
                      <td
                        ><v-chip size="x-small" color="error"
                          >Required</v-chip
                        ></td
                      >
                      <td>Number (no currency symbols)</td>
                      <td>150000</td>
                    </tr>
                    <tr>
                      <td>date_of_event</td>
                      <td
                        ><v-chip size="x-small" color="error"
                          >Required</v-chip
                        ></td
                      >
                      <td>YYYY-MM-DD</td>
                      <td>2024-12-01</td>
                    </tr>
                    <tr>
                      <td>date_notified</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>YYYY-MM-DD (defaults to today)</td>
                      <td>2024-12-06</td>
                    </tr>
                    <tr>
                      <td>priority</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>low, medium, high (defaults to medium)</td>
                      <td>high</td>
                    </tr>
                    <tr>
                      <td>description</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Text</td>
                      <td>Death claim for natural causes</td>
                    </tr>
                  </tbody>
                </v-table>

                <div class="mt-4">
                  <v-alert type="warning" variant="tonal" density="compact">
                    <div class="text-caption">
                      <strong>Important Notes:</strong>
                      <ul class="mt-2 ml-4">
                        <li>CSV file must have headers in the first row</li>
                        <li
                          >Benefit types must match exactly: Group Life
                          Assurance (GLA), Spouse Group Life Assurance (SGLA),
                          Permanent Total Disability (PTD), Critical Illness
                          (CI), Temporary Total Disability (TTD), Personal
                          Health Insurance (PHI), Group Family Funeral (GFF)</li
                        >
                        <li
                          >Member types must match exactly: Member, Spouse,
                          Child, Parent, Dependant</li
                        >

                        <li
                          >Scheme names must match existing schemes exactly</li
                        >
                        <li
                          >Member ID numbers will be validated against existing
                          members</li
                        >
                        <li
                          >All claims will be created with "pending" status
                          initially</li
                        >
                        <li>Maximum file size: 10MB</li>
                        <li>Maximum rows: 1000 claims per upload</li>
                      </ul>
                    </div>
                  </v-alert>
                </div>

                <div class="mt-3">
                  <v-btn
                    color="info"
                    variant="outlined"
                    size="small"
                    prepend-icon="mdi-download"
                    @click="downloadTemplate"
                  >
                    Download CSV Template
                  </v-btn>
                </div>
              </div>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-col>
    </v-row>

    <v-row class="mt-4">
      <v-col cols="12" class="d-flex justify-end gap-2">
        <v-btn color="grey" variant="outlined" @click="handleCancel">
          Cancel
        </v-btn>
        <v-btn
          type="submit"
          color="primary"
          :disabled="!canUpload"
          :loading="uploading"
        >
          Upload {{ previewData.length }} Claims
        </v-btn>
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

// Interfaces
interface ClaimRow {
  claim_number?: string
  member_name: string
  member_id_number: string
  scheme_name: string
  benefit_type: string
  claim_amount: number | string
  date_of_event: string
  date_notified?: string
  priority?: string
  description?: string
}

interface ValidationError {
  row: number
  field: string
  message: string
}

interface Props {
  schemes: Array<{ id: number; name: string }>
}

interface Emits {
  (
    e: 'upload-complete',
    result: { successful: number; failed: number; errors: ValidationError[] }
  ): void
  (e: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form ref
const form = ref()

// State
const uploadFile = ref<File[] | null>(null)
const processing = ref(false)
const uploading = ref(false)
const previewData = ref<ClaimRow[]>([])
const validationErrors = ref<ValidationError[]>([])

// Validation rules
const rules = {
  required: (v: any) => !!v || 'This field is required',
  fileType: (v: File[] | null) => {
    if (!v || (Array.isArray(v) && v.length === 0)) return true
    const file = Array.isArray(v) ? v[0] : v
    if (!file || !file.name) return 'Invalid file'
    return file.name.toLowerCase().endsWith('.csv') || 'File must be a CSV'
  }
}

// Valid benefit types
const validBenefitTypes = [
  'Group Life Assurance (GLA)',
  'Spouse Group Life Assurance (SGLA)',
  'Permanent Total Disability (PTD)',
  'Critical Illness (CI)',
  'Temporary Total Disability (TTD)',
  'Personal Health Insurance (PHI)',
  'Group Family Funeral (GFF)'
]

// Valid priorities
const validPriorities = ['low', 'medium', 'high']

// Preview table headers
const previewHeaders = [
  { title: 'Claim #', value: 'claim_number', width: 150 },
  { title: 'Member', value: 'member_name', width: 180 },
  { title: 'ID Number', value: 'member_id_number', width: 150 },
  { title: 'Scheme', value: 'scheme_name', width: 200 },
  { title: 'Benefit Type', value: 'benefit_type', width: 200 },
  { title: 'Amount', value: 'claim_amount', width: 120 },
  { title: 'Event Date', value: 'date_of_event', width: 120 },
  { title: 'Priority', value: 'priority', width: 100 }
]

// Computed
const canUpload = computed(() => {
  return (
    previewData.value.length > 0 &&
    validationErrors.value.length === 0 &&
    !processing.value &&
    !uploading.value
  )
})

// Methods
const handleFileSelect = async () => {
  if (
    !uploadFile.value ||
    (Array.isArray(uploadFile.value) && uploadFile.value.length === 0)
  ) {
    previewData.value = []
    validationErrors.value = []
    return
  }

  const file = Array.isArray(uploadFile.value)
    ? uploadFile.value[0]
    : uploadFile.value

  // Additional safety check for file object
  if (!file || !file.name) {
    validationErrors.value = [
      { row: 0, field: 'file', message: 'Invalid file selected' }
    ]
    return
  }

  // Check file size (10MB limit)
  if (file.size > 10 * 1024 * 1024) {
    validationErrors.value = [
      { row: 0, field: 'file', message: 'File size exceeds 10MB limit' }
    ]
    return
  }

  processing.value = true
  validationErrors.value = []
  previewData.value = []

  try {
    const csvText = await file.text()
    const { data, errors } = parseCSV(csvText)

    previewData.value = data
    validationErrors.value = errors
  } catch (error) {
    console.error('Error parsing CSV:', error)
    validationErrors.value = [
      { row: 0, field: 'file', message: 'Error reading CSV file' }
    ]
  } finally {
    processing.value = false
  }
}

const parseCSV = (csvText: string) => {
  const lines = csvText.trim().split('\n')
  const data: ClaimRow[] = []
  const errors: ValidationError[] = []

  if (lines.length === 0) {
    errors.push({ row: 0, field: 'file', message: 'CSV file is empty' })
    return { data, errors }
  }

  // Parse headers
  const headers = lines[0].split(',').map((h) => h.trim().replace(/"/g, ''))
  const requiredHeaders = [
    'member_name',
    'member_id_number',
    'member_type',
    'scheme_name',
    'benefit_type',
    'claim_amount',
    'date_of_event'
  ]

  // Check required headers
  const missingHeaders = requiredHeaders.filter((h) => !headers.includes(h))
  if (missingHeaders.length > 0) {
    errors.push({
      row: 0,
      field: 'headers',
      message: `Missing required columns: ${missingHeaders.join(', ')}`
    })
    return { data, errors }
  }

  // Check row count limit
  if (lines.length > 1001) {
    // 1000 data rows + 1 header row
    errors.push({
      row: 0,
      field: 'file',
      message: 'Maximum 1000 claims allowed per upload'
    })
    return { data, errors }
  }

  // Parse data rows
  for (let i = 1; i < lines.length; i++) {
    const row = lines[i].split(',').map((cell) => cell.trim().replace(/"/g, ''))

    if (row.length !== headers.length) {
      errors.push({
        row: i + 1,
        field: 'row',
        message: `Row has ${row.length} columns, expected ${headers.length}`
      })
      continue
    }

    const rowData: any = {}
    headers.forEach((header, index) => {
      rowData[header] = row[index] || ''
    })

    // Validate required fields
    const rowErrors = validateClaimRow(rowData, i + 1)
    errors.push(...rowErrors)

    // Add to preview data even if there are errors
    data.push(rowData as ClaimRow)
  }

  return { data, errors }
}

const validateClaimRow = (row: any, rowNumber: number): ValidationError[] => {
  const errors: ValidationError[] = []

  // Required field validation
  if (!row.member_name?.trim()) {
    errors.push({
      row: rowNumber,
      field: 'member_name',
      message: 'Member name is required'
    })
  }

  if (!row.member_id_number?.trim()) {
    errors.push({
      row: rowNumber,
      field: 'member_id_number',
      message: 'Member ID number is required'
    })
  }

  if (!row.scheme_name?.trim()) {
    errors.push({
      row: rowNumber,
      field: 'scheme_name',
      message: 'Scheme name is required'
    })
  } else {
    // Validate scheme exists
    const schemeExists = props.schemes.some(
      (s) => s.name === row.scheme_name.trim()
    )
    if (!schemeExists) {
      errors.push({
        row: rowNumber,
        field: 'scheme_name',
        message: `Scheme "${row.scheme_name}" not found`
      })
    }
  }

  if (!row.benefit_type?.trim()) {
    errors.push({
      row: rowNumber,
      field: 'benefit_type',
      message: 'Benefit type is required'
    })
  } else {
    // Validate benefit type
    if (!validBenefitTypes.includes(row.benefit_type.trim())) {
      errors.push({
        row: rowNumber,
        field: 'benefit_type',
        message: `Invalid benefit type "${row.benefit_type}"`
      })
    }
  }

  if (!row.claim_amount) {
    errors.push({
      row: rowNumber,
      field: 'claim_amount',
      message: 'Claim amount is required'
    })
  } else {
    const amount = parseFloat(row.claim_amount.toString())
    if (isNaN(amount) || amount <= 0) {
      errors.push({
        row: rowNumber,
        field: 'claim_amount',
        message: 'Claim amount must be a positive number'
      })
    }
  }

  if (!row.date_of_event?.trim()) {
    errors.push({
      row: rowNumber,
      field: 'date_of_event',
      message: 'Date of event is required'
    })
  } else {
    // Validate date format
    const dateRegex = /^\d{4}[-/]\d{2}[-/]\d{2}$/
    if (!dateRegex.test(row.date_of_event.trim())) {
      errors.push({
        row: rowNumber,
        field: 'date_of_event',
        message:
          'Date must be in Date must be in YYYY-MM-DD or YYYY/MM/DD format'
      })
    } else {
      const eventDate = new Date(row.date_of_event.trim())
      if (isNaN(eventDate.getTime())) {
        errors.push({
          row: rowNumber,
          field: 'date_of_event',
          message: 'Invalid date'
        })
      } else if (eventDate > new Date()) {
        errors.push({
          row: rowNumber,
          field: 'date_of_event',
          message: 'Event date cannot be in the future'
        })
      }
    }
  }

  // Optional field validation
  if (row.date_notified?.trim()) {
    const dateRegex = /^\d{4}[-/]\d{2}[-/]\d{2}$/
    if (!dateRegex.test(row.date_notified.trim())) {
      errors.push({
        row: rowNumber,
        field: 'date_notified',
        message: 'Notification Date must be in YYYY-MM-DD or YYYY/MM/DD format'
      })
    }
  }

  if (
    row.priority?.trim() &&
    !validPriorities.includes(row.priority.trim().toLowerCase())
  ) {
    errors.push({
      row: rowNumber,
      field: 'priority',
      message: `Priority must be one of: ${validPriorities.join(', ')}`
    })
  }

  return errors
}

const handleUpload = async () => {
  if (!canUpload.value) return

  uploading.value = true

  try {
    // Prepare claims data for upload
    const claimsData = previewData.value.map((row) => ({
      ...row,
      claim_amount: parseFloat(row.claim_amount.toString()),
      date_notified:
        row.date_notified || new Date().toISOString().split('T')[0],
      priority: row.priority?.toLowerCase() || 'medium',
      status: 'pending'
    }))

    const response = await GroupPricingService.bulkUploadClaims(claimsData)

    emit('upload-complete', {
      successful: response.data.successful || claimsData.length,
      failed: response.data.failed || 0,
      errors: response.data.errors || []
    })
  } catch (error) {
    console.error('Bulk upload error:', error)
    emit('upload-complete', {
      successful: 0,
      failed: previewData.value.length,
      errors: [
        { row: 0, field: 'upload', message: 'Upload failed. Please try again.' }
      ]
    })
  } finally {
    uploading.value = false
  }
}

const handleCancel = () => {
  emit('cancel')
}

const downloadTemplate = () => {
  const headers = [
    'claim_number',
    'member_name',
    'member_id_number',
    'member_type',
    'scheme_name',
    'benefit_type',
    'claim_amount',
    'date_of_event',
    'date_notified',
    'priority',
    'description'
  ]

  const sampleData = [
    '',
    'John Doe',
    '8001015009087',
    'Member',
    'ABC Company Group Scheme',
    'Group Life Assurance (GLA)',
    '150000',
    '2024-12-01',
    '2024-12-06',
    'medium',
    'Death claim for natural causes'
  ]

  const csvContent = [headers.join(','), sampleData.join(',')].join('\n')

  const blob = new Blob([csvContent], { type: 'text/csv' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'claims_upload_template.csv'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  window.URL.revokeObjectURL(url)
}
</script>

<style scoped>
.table-container {
  max-height: 400px;
  overflow: auto;
}

.max-height-200 {
  max-height: 200px;
}

.overflow-y-auto {
  overflow-y: auto;
}

:deep(.v-data-table) {
  font-size: 0.875rem;
}

:deep(.v-data-table__td) {
  padding: 8px 16px !important;
}

:deep(.v-expansion-panel-text__wrapper) {
  padding: 16px;
}
</style>
