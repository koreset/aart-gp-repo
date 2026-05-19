<template>
  <v-form ref="form" @submit.prevent="handleUpload">
    <v-row>
      <v-col cols="12">
        <v-alert type="info" variant="tonal" class="mb-4">
          <div class="text-subtitle-2">Bulk Member Enrollment</div>
          <div class="text-body-2 mt-2">
            Upload a CSV file containing member data to enroll multiple members
            at once. All members will be added to the selected scheme.
          </div>
        </v-alert>
      </v-col>

      <v-col cols="12">
        <v-autocomplete
          v-model="selectedScheme"
          :items="inForceSchemes"
          item-title="name"
          item-value="id"
          label="Target Scheme *"
          variant="outlined"
          density="compact"
          :rules="[rules.required]"
          clearable
          auto-select-first
          required
        />
      </v-col>

      <v-col cols="12">
        <v-file-input
          v-model="uploadFile"
          accept=".csv"
          label="Member Data CSV File *"
          variant="outlined"
          density="compact"
          prepend-icon="mdi-upload"
          :rules="[rules.required, rules.fileType]"
          required
        />
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
                      <td>Male, Female</td>
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
                      <td>1980-01-01</td>
                    </tr>
                    <tr>
                      <td>annual_salary</td>
                      <td
                        ><v-chip size="x-small" color="error"
                          >Required</v-chip
                        ></td
                      >
                      <td>Number</td>
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
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Management, Administration, General</td>
                      <td>General</td>
                    </tr>
                    <tr>
                      <td>email</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Valid email</td>
                      <td>john@company.com</td>
                    </tr>
                    <tr>
                      <td>phone_number</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Phone number</td>
                      <td>0821234567</td>
                    </tr>
                    <tr>
                      <td>employee_number</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Text/Number</td>
                      <td>EMP001</td>
                    </tr>
                    <tr>
                      <td>occupation</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Text</td>
                      <td>Software Developer</td>
                    </tr>
                    <tr>
                      <td>premium</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Number</td>
                      <td>2500.50</td>
                    </tr>
                    <tr>
                      <td>benefits_gla_multiple</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Number (decimal)</td>
                      <td>5.0</td>
                    </tr>
                    <tr>
                      <td>benefits_sgla_multiple</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Number (decimal)</td>
                      <td>4.5</td>
                    </tr>
                    <tr>
                      <td>benefits_ptd_multiple</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Number (decimal)</td>
                      <td>4.0</td>
                    </tr>
                    <tr>
                      <td>benefits_ci_multiple</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Number (decimal)</td>
                      <td>3.0</td>
                    </tr>
                    <tr>
                      <td>benefits_ttd_multiple</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
                        ></td
                      >
                      <td>Number (decimal)</td>
                      <td>0.56</td>
                    </tr>
                    <tr>
                      <td>benefits_phi_multiple</td>
                      <td
                        ><v-chip size="x-small" color="warning"
                          >Optional</v-chip
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
      </v-col>

      <v-col cols="12">
        <v-btn
          size="small"
          rounded
          color="info"
          variant="outlined"
          class="mb-3"
          @click="downloadTemplate"
        >
          <v-icon left>mdi-download</v-icon>
          Download CSV Template
        </v-btn>
      </v-col>

      <!-- Validation Options -->
      <v-col cols="12">
        <v-card variant="outlined">
          <v-card-title class="text-subtitle-2">Upload Options</v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="validateOnly"
                  label="Validate only (don't import)"
                  color="primary"
                  hide-details
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="skipDuplicates"
                  label="Skip duplicate ID numbers"
                  color="primary"
                  hide-details
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
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
      <v-col v-if="uploadResults" cols="12">
        <v-card variant="outlined">
          <v-card-title class="text-subtitle-2">Upload Results</v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="6" md="3">
                <div class="text-center">
                  <div class="text-h4 text-success">{{
                    uploadResults.success
                  }}</div>
                  <div class="text-caption">Successful</div>
                </div>
              </v-col>
              <v-col cols="6" md="3">
                <div class="text-center">
                  <div class="text-h4 text-error">{{
                    uploadResults.failed
                  }}</div>
                  <div class="text-caption">Failed</div>
                </div>
              </v-col>
              <v-col cols="6" md="3">
                <div class="text-center">
                  <div class="text-h4 text-warning">{{
                    uploadResults.duplicates
                  }}</div>
                  <div class="text-caption">Duplicates</div>
                </div>
              </v-col>
              <v-col cols="6" md="3">
                <div class="text-center">
                  <div class="text-h4 text-info">{{ uploadResults.total }}</div>
                  <div class="text-caption">Total Rows</div>
                </div>
              </v-col>
            </v-row>

            <MemberUploadErrorReport
              v-if="
                uploadResults.blockingErrors.length > 0 ||
                uploadResults.errors.length > 0
              "
              class="mt-4"
              :blocking-errors="uploadResults.blockingErrors"
              :soft-errors="uploadResults.errors"
              :context-label="selectedSchemeName"
              filename-prefix="member-upload-errors"
            />
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Action Buttons -->
      <v-col cols="12">
        <v-card-actions class="justify-end">
          <v-btn
            size="small"
            rounded
            color="grey"
            variant="outlined"
            @click="$emit('cancel')"
            >Cancel</v-btn
          >
          <v-btn
            size="small"
            rounded
            color="primary"
            type="submit"
            :loading="uploading"
            :disabled="!isFormValid"
          >
            {{ validateOnly ? 'Validate File' : 'Upload Members' }}
          </v-btn>
        </v-card-actions>
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Papa from 'papaparse'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import MemberUploadErrorReport from '@/renderer/screens/group_pricing/shared/MemberUploadErrorReport.vue'
// import { format } from 'date-fns'

interface Props {
  schemes: Array<any>
}

interface Emits {
  (e: 'upload-complete', result: any): void
  (e: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Interfaces
interface BlockingError {
  row: number
  field?: string
  message: string
  member_id_number?: string
  member_name?: string
}

interface UploadResult {
  success: number
  failed: number
  duplicates: number
  total: number
  errors: Array<{ row: number; message: string }>
  blockingErrors: BlockingError[]
}

interface UploadProgress {
  show: boolean
  percentage: number
  message: string
}

const form = ref(null)
const uploading = ref(false)
const selectedScheme = ref<number | null>(null)
const uploadFile = ref<File | null>(null)
const validateOnly = ref(false)
const skipDuplicates = ref(true)

const inForceSchemes = computed(() =>
  props.schemes.filter((s) => s?.status === 'in_force')
)

const selectedSchemeName = computed(() => {
  const s = props.schemes.find((x) => x.id === selectedScheme.value)
  return s?.name ?? ''
})

const uploadProgress = ref<UploadProgress>({
  show: false,
  percentage: 0,
  message: ''
})

const uploadResults = ref<UploadResult | null>(null)

// Validation rules
const rules = {
  required: (value: any) => !!value || 'Field is required',
  fileType: (value: any) => {
    if (!value) return true
    return value?.type === 'text/csv' || 'File must be CSV format'
  }
}

// Computed properties
const isFormValid = computed(() => {
  const hasScheme = !!selectedScheme.value
  const hasFile = uploadFile.value instanceof File

  return hasScheme && hasFile
})

// Methods
const handleUpload = async () => {
  if (!form.value) return

  const { valid } = await (form.value as any).validate()
  if (!valid) return

  uploading.value = true
  uploadProgress.value.show = true
  uploadProgress.value.percentage = 0
  uploadProgress.value.message = 'Reading file...'

  try {
    if (!uploadFile.value) {
      throw new Error('No file selected')
    }

    const file = uploadFile.value
    const csvData = await parseCSVFile(file)

    uploadProgress.value.percentage = 25
    uploadProgress.value.message = 'Validating data...'

    const validationResult = validateCSVData(csvData)

    console.log('Validation Result:', validationResult)

    if (validateOnly.value) {
      // Dry run: report blocking + soft errors but never upload.
      uploadResults.value = {
        success: validationResult.valid.length,
        failed:
          validationResult.errors.length +
          validationResult.blockingErrors.length,
        duplicates: validationResult.duplicates.length,
        total: csvData.length,
        errors: validationResult.errors,
        blockingErrors: validationResult.blockingErrors
      }

      uploadProgress.value.percentage = 100
      uploadProgress.value.message =
        validationResult.blockingErrors.length > 0
          ? 'Validation found blocking errors — fix and re-validate'
          : 'Validation complete'
    } else if (validationResult.blockingErrors.length > 0) {
      // GATE: any row missing gender / date_of_birth / annual_salary aborts
      // the entire upload. No HTTP call is made — the user must correct the
      // CSV and re-upload. This avoids wasted cloud round-trips and partial
      // scheme state.
      uploadResults.value = {
        success: 0,
        failed: validationResult.blockingErrors.length,
        duplicates: validationResult.duplicates.length,
        total: csvData.length,
        errors: validationResult.errors,
        blockingErrors: validationResult.blockingErrors
      }
      uploadProgress.value.percentage = 100
      uploadProgress.value.message =
        'Upload blocked — correct source file and re-upload'
    } else {
      // Process the upload via the bulk endpoint — one HTTP request, one
      // transaction, one audit-log entry.
      uploadProgress.value.percentage = 50
      uploadProgress.value.message = 'Processing members...'

      const processResult = await processMemberUploads(validationResult.valid)

      console.log('Process Result:', processResult)

      uploadResults.value = {
        success: processResult.success,
        failed: processResult.failed + validationResult.errors.length,
        duplicates: validationResult.duplicates.length,
        total: csvData.length,
        errors: [...validationResult.errors, ...processResult.errors],
        blockingErrors: processResult.blockingErrors
      }

      uploadProgress.value.percentage = 100
      uploadProgress.value.message =
        processResult.blockingErrors.length > 0
          ? 'Upload blocked by server validation — correct source file and re-upload'
          : 'Upload complete'

      // Emit completion event
      console.log(
        'Emitting upload-complete event with results:',
        uploadResults.value
      )
      emit('upload-complete', uploadResults.value)
    }
  } catch (error: any) {
    console.error('Upload error:', error)
    uploadResults.value = {
      success: 0,
      failed: 1,
      duplicates: 0,
      total: 1,
      errors: [
        {
          row: 1,
          message:
            'Failed to process file: ' + (error?.message || 'Unknown error')
        }
      ],
      blockingErrors: []
    }
  } finally {
    uploading.value = false
  }
}

const parseCSVFile = (file: File): Promise<any[]> => {
  return new Promise((resolve, reject) => {
    Papa.parse(file, {
      header: true,
      skipEmptyLines: true,
      complete: (results: any) => {
        if (results.errors.length > 0) {
          reject(new Error('CSV parsing errors: ' + results.errors[0].message))
        } else {
          resolve(results.data)
        }
      },
      error: (error: any) => {
        reject(error)
      }
    })
  })
}

const validateCSVData = (data: any[]) => {
  const valid: any[] = []
  const errors: Array<{ row: number; message: string }> = []
  const blockingErrors: BlockingError[] = []
  const duplicates: Array<{ row: number; message: string }> = []
  const seenIds = new Set()

  const blockingFields: Array<keyof any> = [
    'gender',
    'date_of_birth',
    'annual_salary'
  ]
  const otherRequiredFields = [
    'member_name',
    'member_id_number',
    'member_id_type',
    'entry_date'
  ]

  const scheme = props.schemes.find((s) => s.id === selectedScheme.value)
  const schemeCommencement = scheme?.commencement_date
    ? new Date(scheme.commencement_date)
    : null

  data.forEach((row, index) => {
    const rowNumber = index + 2 // +2 because CSV header is row 1, data starts at row 2

    // BLOCKING: gender / date_of_birth / annual_salary missing → reject the
    // whole upload. These drive pricing/valuation maths and silent partial
    // imports would corrupt downstream calculations.
    let rowBlocked = false
    for (const field of blockingFields as string[]) {
      const v = row[field]
      const isEmpty =
        v === undefined ||
        v === null ||
        (typeof v === 'string' && v.trim() === '')
      if (isEmpty) {
        blockingErrors.push({
          row: rowNumber,
          field,
          message: `${field} is required`,
          member_id_number: row.member_id_number,
          member_name: row.member_name
        })
        rowBlocked = true
      }
    }
    // Non-positive salary also blocks (a 0 or negative salary is meaningless
    // for pricing and a common source of dirty data).
    if (
      !rowBlocked &&
      row.annual_salary !== undefined &&
      row.annual_salary !== null &&
      String(row.annual_salary).trim() !== ''
    ) {
      const salary = Number(row.annual_salary)
      if (isNaN(salary) || salary <= 0) {
        blockingErrors.push({
          row: rowNumber,
          field: 'annual_salary',
          message: 'annual_salary must be a positive number',
          member_id_number: row.member_id_number,
          member_name: row.member_name
        })
        rowBlocked = true
      }
    }
    if (rowBlocked) {
      return
    }

    // Other required fields → existing per-row error (row skipped, upload
    // still proceeds for the remaining valid rows).
    const missingFields = otherRequiredFields.filter(
      (field) => !row[field]?.trim()
    )
    if (missingFields.length > 0) {
      errors.push({
        row: rowNumber,
        message: `Missing required fields: ${missingFields.join(', ')}`
      })
      return
    }

    // Check for duplicate ID numbers
    if (seenIds.has(row.member_id_number)) {
      duplicates.push({
        row: rowNumber,
        message: `Duplicate ID number: ${row.member_id_number}`
      })
      if (!skipDuplicates.value) {
        errors.push({
          row: rowNumber,
          message: `Duplicate ID number: ${row.member_id_number}`
        })
        return
      }
    }
    seenIds.add(row.member_id_number)

    // Validate ID number format
    if (row.member_id_type === 'ID' && row.member_id_number.length !== 13) {
      errors.push({
        row: rowNumber,
        message: 'RSA ID number must be 13 digits'
      })
      return
    }

    // Validate gender
    if (!['Male', 'Female'].includes(row.gender)) {
      errors.push({
        row: rowNumber,
        message: 'Gender must be Male or Female'
      })
      return
    }

    // Validate date format
    if (!isValidDate(row.date_of_birth)) {
      errors.push({
        row: rowNumber,
        message: 'Invalid date format for date_of_birth (use YYYY-MM-DD)'
      })
      return
    }

    if (!isValidDate(row.entry_date)) {
      errors.push({
        row: rowNumber,
        message: 'Invalid date format for entry_date (use YYYY-MM-DD)'
      })
      return
    }

    if (schemeCommencement && new Date(row.entry_date) < schemeCommencement) {
      errors.push({
        row: rowNumber,
        message: `Entry date (${row.entry_date}) cannot be before scheme commencement date (${schemeCommencement.toISOString().slice(0, 10)})`
      })
      return
    }

    // Validate optional premium field
    if (
      row.premium &&
      (isNaN(Number(row.premium)) || Number(row.premium) < 0)
    ) {
      errors.push({
        row: rowNumber,
        message: 'Premium must be a valid positive number or zero'
      })
      return
    }

    // Validate optional benefit multiple fields
    const benefitMultipleFields = [
      'benefits_gla_multiple',
      'benefits_sgla_multiple',
      'benefits_ptd_multiple',
      'benefits_ci_multiple',
      'benefits_ttd_multiple',
      'benefits_phi_multiple'
    ]

    for (const field of benefitMultipleFields) {
      if (row[field] && (isNaN(Number(row[field])) || Number(row[field]) < 0)) {
        errors.push({
          row: rowNumber,
          message: `${field} must be a valid positive number or zero`
        })
        return
      }
    }

    // If all validations pass, add to valid array
    const memberData: any = {
      ...row,
      rowNumber,
      annual_salary: Number(row.annual_salary),
      scheme_id: selectedScheme.value,
      scheme_category: row.scheme_category || 'General'
    }

    // Add optional premium field
    if (row.premium) {
      memberData.premium = Number(row.premium)
    }

    // Add optional benefit multiple fields
    benefitMultipleFields.forEach((field) => {
      if (row[field]) {
        memberData[field] = Number(row[field])
      }
    })

    valid.push(memberData)
  })

  return { valid, errors, duplicates, blockingErrors }
}

const toMemberPayload = (member: any) => {
  const payload: any = {
    ...member,
    benefits: {
      gla_multiple: member.benefits_gla_multiple,
      sgla_multiple: member.benefits_sgla_multiple,
      ptd_multiple: member.benefits_ptd_multiple,
      ci_multiple: member.benefits_ci_multiple,
      ttd_multiple: member.benefits_ttd_multiple,
      phi_multiple: member.benefits_phi_multiple
    }
  }
  delete payload.benefits_gla_multiple
  delete payload.benefits_sgla_multiple
  delete payload.benefits_ptd_multiple
  delete payload.benefits_ci_multiple
  delete payload.benefits_ttd_multiple
  delete payload.benefits_phi_multiple
  delete payload.rowNumber
  return payload
}

const processMemberUploads = async (validMembers: any[]) => {
  const results = {
    success: 0,
    failed: 0,
    errors: [] as Array<{ row: number; message: string }>,
    blockingErrors: [] as BlockingError[]
  }

  if (validMembers.length === 0) {
    return results
  }

  uploadProgress.value.percentage = 70
  uploadProgress.value.message = `Uploading ${validMembers.length} members in one request...`

  const members = validMembers.map(toMemberPayload)

  try {
    const response = await GroupPricingService.addMembersBulk(
      selectedScheme.value,
      members,
      skipDuplicates.value
    )
    results.success = response?.data?.inserted ?? validMembers.length
  } catch (error: any) {
    const data = error?.response?.data
    // Server returned structured blocking errors → surface them to the user
    // so they can fix the CSV and re-upload. Nothing was inserted.
    if (data?.blocking_errors && Array.isArray(data.blocking_errors)) {
      results.blockingErrors = data.blocking_errors
      results.failed = data.blocking_errors.length
    } else {
      results.failed = validMembers.length
      results.errors.push({
        row: 0,
        message: `Bulk upload failed: ${data?.error || error?.message || 'Unknown error'}`
      })
    }
  }

  return results
}

const isValidDate = (dateString: string) => {
  const regex = /^\d{4}[-/]\d{2}[-/]\d{2}$/
  if (!regex.test(dateString)) return false

  const date = new Date(dateString)
  return date instanceof Date && !isNaN(date.getTime())
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
      premium: 2500.5,
      benefits_gla_multiple: 5.0,
      benefits_sgla_multiple: 4.5,
      benefits_ptd_multiple: 4.0,
      benefits_ci_multiple: 3.0,
      benefits_ttd_multiple: 0.56,
      benefits_phi_multiple: 0.7
    }
  ]

  const formattedTemplateData = templateData.map((row) => ({
    ...row
  }))

  const csv = Papa.unparse(formattedTemplateData, {
    quotes: true, // Force quotes around all fields
    quoteChar: '"',
    escapeChar: '"',
    delimiter: ',',
    header: true
  })
  console.log('Generated CSV Template:\n', csv)
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')

  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', 'member_enrollment_template.csv')
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
</script>

<style scoped>
.v-card-title {
  padding: 12px 16px;
}

.v-card-text {
  padding: 16px;
}
</style>
