<template>
  <v-card variant="outlined">
    <v-card-text>
      <v-alert
        v-if="blockingErrors.length > 0"
        type="error"
        variant="tonal"
        density="compact"
        class="mb-3"
      >
        <div class="font-weight-bold">Upload blocked</div>
        <div class="text-body-2 mt-1">
          {{ blockingErrors.length }} member{{ blockingErrors.length === 1 ? '' : 's' }}
          {{ blockingErrors.length === 1 ? 'is' : 'are' }} missing required fields
          (gender, date of birth, annual salary). Please correct the source file
          and re-upload.
        </div>
      </v-alert>

      <div v-if="blockingErrors.length > 0" class="d-flex justify-end mb-2">
        <v-btn
          size="small"
          rounded
          color="primary"
          variant="outlined"
          @click="downloadErrorReport"
        >
          <v-icon left>mdi-download</v-icon>
          Download error report
        </v-btn>
      </div>

      <v-expansion-panels v-if="blockingErrors.length > 0">
        <v-expansion-panel>
          <v-expansion-panel-title>
            <v-icon left color="error">mdi-alert</v-icon>
            Blocking errors ({{ blockingErrors.length }})
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-list density="compact">
              <v-list-item
                v-for="(err, idx) in blockingErrors"
                :key="`block-${idx}`"
              >
                <v-list-item-title class="text-caption text-error">
                  Row {{ err.row }}<span v-if="err.member_id_number">
                    ({{ err.member_id_number }})</span>: {{ err.message }}
                </v-list-item-title>
              </v-list-item>
            </v-list>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>

      <v-expansion-panels v-if="softErrors.length > 0" class="mt-2">
        <v-expansion-panel>
          <v-expansion-panel-title>
            <v-icon left color="warning">mdi-alert-outline</v-icon>
            Other errors ({{ softErrors.length }})
          </v-expansion-panel-title>
          <v-expansion-panel-text>
            <v-list density="compact">
              <v-list-item
                v-for="(err, idx) in softErrors"
                :key="`soft-${idx}`"
              >
                <v-list-item-title class="text-caption text-warning">
                  Row {{ err.row }}: {{ err.message }}
                </v-list-item-title>
              </v-list-item>
            </v-list>
          </v-expansion-panel-text>
        </v-expansion-panel>
      </v-expansion-panels>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import Papa from 'papaparse'

export interface BlockingError {
  row: number
  field?: string
  message: string
  member_id_number?: string
  member_name?: string
}

export interface SoftError {
  row: number
  message: string
}

interface Props {
  blockingErrors: BlockingError[]
  softErrors?: SoftError[]
  filenamePrefix?: string
  contextLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  softErrors: () => [],
  filenamePrefix: 'member-upload-errors',
  contextLabel: ''
})

const downloadErrorReport = () => {
  const rows = props.blockingErrors.map((e) => ({
    row_number: e.row,
    member_id_number: e.member_id_number ?? '',
    member_name: e.member_name ?? '',
    missing_field: e.field ?? '',
    error_message: e.message
  }))
  const csv = Papa.unparse(rows)
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  const stamp = new Date().toISOString().replace(/[:.]/g, '-')
  const ctx = props.contextLabel
    ? `-${props.contextLabel.replace(/[^A-Za-z0-9_-]+/g, '_')}`
    : ''
  link.href = url
  link.download = `${props.filenamePrefix}${ctx}-${stamp}.csv`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
}
</script>
