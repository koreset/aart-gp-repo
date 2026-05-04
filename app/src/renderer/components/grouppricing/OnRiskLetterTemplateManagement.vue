<template>
  <base-card>
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-file-document-check-outline" class="mr-3" />
        <span>On Risk Letter Template</span>
        <v-spacer />
        <v-chip
          :color="active ? 'success' : 'warning'"
          variant="outlined"
          size="small"
        >
          {{
            active ? `Active v${active.version}` : 'Using default formatting'
          }}
        </v-chip>
      </div>
    </template>

    <template #default>
      <p class="text-subtitle-2 text-medium-emphasis mb-4">
        Upload a Word (.docx) template to customise the On Risk letter sent when
        a quote is accepted. Tokens like
        <code v-pre>{{ scheme_name }}</code
        >, <code v-pre>{{ commencement_date }}</code> and
        <code v-pre>{{#benefit_summary}}...{{/benefit_summary}}</code> are
        filled with the accepted quote's data at generation time. Download the
        sample template to see all supported tokens.
      </p>

      <v-card variant="tonal" class="mb-4">
        <v-card-text>
          <div v-if="loading" class="d-flex align-center">
            <v-progress-circular indeterminate size="20" class="mr-3" />
            Loading active template…
          </div>
          <div v-else-if="active">
            <div class="d-flex align-center">
              <v-icon
                icon="mdi-file-document-check"
                color="primary"
                size="large"
                class="mr-3"
              />
              <div class="flex-grow-1">
                <div class="text-subtitle-1 font-weight-medium">
                  {{ active.filename }}
                </div>
                <div class="text-caption text-medium-emphasis">
                  Version {{ active.version }} ·
                  {{ formatBytes(active.size_bytes) }} · uploaded
                  {{ formatDate(active.uploaded_at) }}
                  <span v-if="active.uploaded_by">
                    by {{ active.uploaded_by }}</span
                  >
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-medium-emphasis">
            <v-icon icon="mdi-information-outline" class="mr-2" />
            No template uploaded — On Risk letters use the system's default
            built-in layout.
          </div>
        </v-card-text>
      </v-card>

      <div class="d-flex flex-wrap ga-2 mb-4">
        <v-btn
          color="primary"
          prepend-icon="mdi-cloud-upload"
          :loading="uploading"
          :disabled="!insurerId"
          @click="triggerFileInput"
        >
          Upload new template
        </v-btn>
        <v-btn
          color="info"
          variant="tonal"
          prepend-icon="mdi-download"
          :disabled="!active"
          @click="downloadActive"
        >
          Download active template
        </v-btn>
        <v-btn
          color="secondary"
          variant="outlined"
          prepend-icon="mdi-file-document-outline"
          @click="downloadSample"
        >
          Download sample template
        </v-btn>
        <v-btn
          color="error"
          variant="outlined"
          prepend-icon="mdi-trash-can-outline"
          :disabled="!active"
          :loading="deletingActive"
          @click="confirmDeleteActive"
        >
          Delete active template
        </v-btn>
        <input
          ref="fileInput"
          type="file"
          accept=".docx"
          style="display: none"
          @change="handleFileSelected"
        />
      </div>

      <div v-if="previousVersions.length > 0">
        <div class="d-flex align-center mb-2">
          <h5 class="section-title mb-0">
            <v-icon icon="mdi-history" size="small" class="mr-2" />
            Previous versions
          </h5>
          <v-spacer />
          <v-btn
            size="small"
            variant="outlined"
            color="error"
            prepend-icon="mdi-trash-can-outline"
            :loading="deletingInactive"
            @click="confirmDeleteAllInactive"
          >
            Delete all previous
          </v-btn>
        </div>
        <v-table density="compact">
          <thead>
            <tr>
              <th>Version</th>
              <th>Filename</th>
              <th>Uploaded</th>
              <th>Uploaded by</th>
              <th class="text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="t in previousVersions" :key="t.id">
              <td>v{{ t.version }}</td>
              <td>{{ t.filename }}</td>
              <td>{{ formatDate(t.uploaded_at) }}</td>
              <td>{{ t.uploaded_by || '—' }}</td>
              <td class="actions-cell text-right">
                <v-btn
                  size="small"
                  variant="text"
                  icon="mdi-download"
                  density="comfortable"
                  @click="downloadVersion(t)"
                >
                  <v-icon>mdi-download</v-icon>
                  <v-tooltip activator="parent" location="top">
                    Download
                  </v-tooltip>
                </v-btn>
                <v-btn
                  size="small"
                  variant="text"
                  color="primary"
                  icon="mdi-check"
                  density="comfortable"
                  @click="activateVersion(t)"
                >
                  <v-icon>mdi-check</v-icon>
                  <v-tooltip activator="parent" location="top">
                    Activate
                  </v-tooltip>
                </v-btn>
                <v-btn
                  size="small"
                  variant="text"
                  color="error"
                  icon="mdi-trash-can-outline"
                  density="comfortable"
                  :loading="deletingId === t.id"
                  @click="confirmDelete(t)"
                >
                  <v-icon>mdi-trash-can-outline</v-icon>
                  <v-tooltip activator="parent" location="top">
                    Delete
                  </v-tooltip>
                </v-btn>
              </td>
            </tr>
          </tbody>
        </v-table>
      </div>

      <!-- Delete confirmation dialog -->
      <v-dialog v-model="deleteDialog.open" max-width="480">
        <v-card>
          <v-card-title class="text-h6">
            <v-icon
              icon="mdi-alert-circle-outline"
              color="error"
              class="mr-2"
            />
            {{ deleteDialog.title }}
          </v-card-title>
          <v-card-text>{{ deleteDialog.message }}</v-card-text>
          <v-card-actions>
            <v-spacer />
            <v-btn variant="text" @click="deleteDialog.open = false">
              Cancel
            </v-btn>
            <v-btn color="error" variant="flat" @click="deleteDialog.confirm">
              Delete
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
    </template>
  </base-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import { saveAs } from 'file-saver'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useNotifications } from '@/renderer/composables/useNotifications'
import BaseCard from '../BaseCard.vue'

const { showSuccess, showError } = useNotifications()

const insurerId = ref<number | null>(null)
const active = ref<any>(null)
const versions = ref<any[]>([])
const loading = ref(false)
const uploading = ref(false)
const deletingId = ref<number | null>(null)
const deletingInactive = ref(false)
const deletingActive = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

// Single dialog used for both single-version and bulk deletes.
const deleteDialog = reactive<{
  open: boolean
  title: string
  message: string
  confirm: () => void
}>({
  open: false,
  title: '',
  message: '',
  confirm: () => {}
})

const previousVersions = computed(() =>
  versions.value.filter((v) => !v.is_active)
)

function formatBytes(n: number): string {
  if (!n) return '0 B'
  if (n < 1024) return `${n} B`
  if (n < 1024 * 1024) return `${(n / 1024).toFixed(1)} KB`
  return `${(n / (1024 * 1024)).toFixed(1)} MB`
}

function formatDate(d: string): string {
  if (!d) return ''
  const date = new Date(d)
  return date.toLocaleString()
}

async function loadAll() {
  if (!insurerId.value) return
  loading.value = true
  try {
    try {
      const r = await GroupPricingService.getActiveInsurerOnRiskLetterTemplate(
        insurerId.value
      )
      active.value = r.data
    } catch (e: any) {
      if (e?.response?.status === 404) active.value = null
      else throw e
    }
    const v = await GroupPricingService.listInsurerOnRiskLetterTemplateVersions(
      insurerId.value
    )
    versions.value = v.data || []
  } catch (e: any) {
    console.error('Failed to load templates', e)
    showError('Failed to load templates')
  } finally {
    loading.value = false
  }
}

function triggerFileInput() {
  fileInput.value?.click()
}

async function handleFileSelected(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file || !insurerId.value) return
  if (!file.name.toLowerCase().endsWith('.docx')) {
    showError('Template must be a .docx file')
    return
  }
  uploading.value = true
  try {
    await GroupPricingService.uploadInsurerOnRiskLetterTemplate(
      insurerId.value,
      file
    )
    showSuccess(`Template "${file.name}" uploaded successfully`)
    await loadAll()
  } catch (e: any) {
    console.error(e)
    showError('Failed to upload template')
  } finally {
    uploading.value = false
    if (target) target.value = ''
  }
}

async function downloadActive() {
  if (!active.value) return
  await downloadVersion(active.value)
}

async function downloadVersion(t: any) {
  try {
    const r = await GroupPricingService.downloadInsurerOnRiskLetterTemplate(
      t.id
    )
    saveAs(r.data, t.filename || `on_risk_letter_template_v${t.version}.docx`)
  } catch (e) {
    console.error(e)
    showError('Failed to download template')
  }
}

async function downloadSample() {
  try {
    const r = await GroupPricingService.downloadSampleOnRiskLetterTemplate()
    saveAs(r.data, 'sample_on_risk_letter_template.docx')
  } catch (e) {
    console.error(e)
    showError('Failed to download sample template')
  }
}

async function activateVersion(t: any) {
  if (!insurerId.value) return
  try {
    await GroupPricingService.activateInsurerOnRiskLetterTemplate(
      insurerId.value,
      t.id
    )
    showSuccess(`Activated version ${t.version}`)
    await loadAll()
  } catch (e) {
    console.error(e)
    showError('Failed to activate version')
  }
}

function confirmDelete(t: any) {
  deleteDialog.title = `Delete version ${t.version}?`
  deleteDialog.message = `"${t.filename}" will be permanently removed. This cannot be undone.`
  deleteDialog.confirm = () => {
    deleteDialog.open = false
    deleteVersion(t)
  }
  deleteDialog.open = true
}

async function deleteVersion(t: any) {
  if (!insurerId.value) return
  deletingId.value = t.id
  try {
    await GroupPricingService.deleteInsurerOnRiskLetterTemplate(
      insurerId.value,
      t.id
    )
    showSuccess(`Deleted version ${t.version}`)
    await loadAll()
  } catch (e: any) {
    console.error(e)
    const msg =
      e?.response?.status === 409
        ? e.response.data?.error ||
          'Cannot delete the active template. Activate another version first.'
        : 'Failed to delete version'
    showError(msg)
  } finally {
    deletingId.value = null
  }
}

function confirmDeleteActive() {
  if (!active.value) return
  const t = active.value
  deleteDialog.title = `Delete active template (v${t.version})?`
  deleteDialog.message = `"${t.filename}" is the active template. Deleting it will revert On Risk letters to the system's default built-in layout. This cannot be undone.`
  deleteDialog.confirm = () => {
    deleteDialog.open = false
    deleteActive()
  }
  deleteDialog.open = true
}

async function deleteActive() {
  if (!insurerId.value || !active.value) return
  deletingActive.value = true
  try {
    await GroupPricingService.deleteInsurerOnRiskLetterTemplate(
      insurerId.value,
      active.value.id,
      { force: true }
    )
    showSuccess('Active template deleted')
    await loadAll()
  } catch (e) {
    console.error(e)
    showError('Failed to delete active template')
  } finally {
    deletingActive.value = false
  }
}

function confirmDeleteAllInactive() {
  const n = previousVersions.value.length
  deleteDialog.title = `Delete all ${n} previous version${n === 1 ? '' : 's'}?`
  deleteDialog.message = `The active template is kept. All ${n} previous version${n === 1 ? '' : 's'} will be permanently removed. This cannot be undone.`
  deleteDialog.confirm = () => {
    deleteDialog.open = false
    deleteAllInactive()
  }
  deleteDialog.open = true
}

async function deleteAllInactive() {
  if (!insurerId.value) return
  deletingInactive.value = true
  try {
    const r =
      await GroupPricingService.deleteInactiveInsurerOnRiskLetterTemplates(
        insurerId.value
      )
    const n = r?.data?.deleted ?? 0
    showSuccess(`Deleted ${n} previous version${n === 1 ? '' : 's'}`)
    await loadAll()
  } catch (e) {
    console.error(e)
    showError('Failed to delete previous versions')
  } finally {
    deletingInactive.value = false
  }
}

onMounted(async () => {
  try {
    const r = await GroupPricingService.getInsurer()
    if (r?.data?.id) {
      insurerId.value = r.data.id
      await loadAll()
    }
  } catch (e) {
    console.error('Could not load insurer details', e)
  }
})
</script>

<style scoped>
.section-title {
  color: #1976d2;
  font-weight: 600;
  font-size: 1rem;
  margin: 1rem 0 0.5rem;
  display: flex;
  align-items: center;
}
code {
  background: rgba(0, 0, 0, 0.05);
  padding: 1px 4px;
  border-radius: 3px;
  font-size: 0.85em;
}
.actions-cell {
  white-space: nowrap;
  width: 1%;
}
.actions-cell .v-btn {
  margin-left: 4px;
}
</style>
