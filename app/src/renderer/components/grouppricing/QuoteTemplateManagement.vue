<template>
  <base-card>
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-file-word-outline" class="mr-3" />
        <span>Quote Template</span>
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
        Upload a Word (.docx) template to customise quote documents for this
        insurer. Tokens like <code v-pre>{{ scheme_name }}</code> and
        <code v-pre>{{#categories}}...{{/categories}}</code> are filled with
        quote data at generation time. Download the sample template to see all
        supported tokens.
      </p>

      <!-- Active template card -->
      <v-card variant="tonal" class="mb-4">
        <v-card-text>
          <div v-if="loading" class="d-flex align-center">
            <v-progress-circular indeterminate size="20" class="mr-3" />
            Loading active template…
          </div>
          <div v-else-if="active">
            <div class="d-flex align-center">
              <v-icon
                icon="mdi-file-word"
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
            No template uploaded — quote documents use the system's default
            built-in layout.
          </div>
        </v-card-text>
      </v-card>

      <!-- Action buttons -->
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
        <input
          ref="fileInput"
          type="file"
          accept=".docx"
          style="display: none"
          @change="handleFileSelected"
        />
      </div>

      <!-- Previous versions -->
      <div v-if="previousVersions.length > 0">
        <h5 class="section-title">
          <v-icon icon="mdi-history" size="small" class="mr-2" />
          Previous versions
        </h5>
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
              <td class="text-right">
                <v-btn
                  size="small"
                  variant="text"
                  prepend-icon="mdi-download"
                  @click="downloadVersion(t)"
                >
                  Download
                </v-btn>
                <v-btn
                  size="small"
                  variant="text"
                  color="primary"
                  prepend-icon="mdi-check"
                  @click="activateVersion(t)"
                >
                  Activate
                </v-btn>
              </td>
            </tr>
          </tbody>
        </v-table>
      </div>
    </template>
  </base-card>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
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
const fileInput = ref<HTMLInputElement | null>(null)

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
    // Active (404 is normal — means no template yet)
    try {
      const r = await GroupPricingService.getActiveInsurerQuoteTemplate(
        insurerId.value
      )
      active.value = r.data
    } catch (e: any) {
      if (e?.response?.status === 404) active.value = null
      else throw e
    }
    // Versions
    const v = await GroupPricingService.listInsurerQuoteTemplateVersions(
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
    await GroupPricingService.uploadInsurerQuoteTemplate(insurerId.value, file)
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
    const r = await GroupPricingService.downloadInsurerQuoteTemplate(t.id)
    saveAs(r.data, t.filename || `template_v${t.version}.docx`)
  } catch (e) {
    console.error(e)
    showError('Failed to download template')
  }
}

async function downloadSample() {
  try {
    const r = await GroupPricingService.downloadSampleQuoteTemplate()
    saveAs(r.data, 'sample_quote_template.docx')
  } catch (e) {
    console.error(e)
    showError('Failed to download sample template')
  }
}

async function activateVersion(t: any) {
  if (!insurerId.value) return
  try {
    await GroupPricingService.activateInsurerQuoteTemplate(
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

onMounted(async () => {
  // Discover the insurer ID from the existing single-insurer endpoint.
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
</style>
