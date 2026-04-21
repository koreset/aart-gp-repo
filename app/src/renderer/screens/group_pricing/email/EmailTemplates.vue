<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center justify-space-between">
              <div class="d-flex align-center">
                <v-btn
                  class="mr-3"
                  size="small"
                  variant="text"
                  prepend-icon="mdi-arrow-left"
                  @click="$router.back()"
                >
                  Back
                </v-btn>
                <span class="headline">Email Templates</span>
              </div>
              <div class="d-flex align-center gap-2">
                <v-btn
                  v-if="templates.length === 0"
                  class="mr-2"
                  variant="outlined"
                  size="small"
                  rounded
                  prepend-icon="mdi-sparkles"
                  :loading="seeding"
                  @click="seedStarterTemplates"
                >
                  Seed Starter Templates
                </v-btn>
                <v-btn
                  size="small"
                  variant="outlined"
                  rounded
                  prepend-icon="mdi-plus"
                  @click="createNew"
                >
                  New Template
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <v-row>
              <!-- List -->
              <v-col cols="12" md="4">
                <v-card variant="outlined">
                  <v-list density="compact" nav>
                    <v-list-item
                      v-for="t in templates"
                      :key="t.code"
                      :active="selected?.code === t.code"
                      :title="t.name || t.code"
                      :subtitle="t.code"
                      @click="selectTemplate(t.code)"
                    >
                      <template #append>
                        <v-chip
                          size="x-small"
                          :color="t.status === 'active' ? 'success' : 'warning'"
                          variant="tonal"
                        >
                          {{ t.status }}
                        </v-chip>
                      </template>
                    </v-list-item>
                    <v-list-item v-if="templates.length === 0">
                      <p class="text-caption text-medium-emphasis">
                        No templates yet — create one or seed the starter set.
                      </p>
                    </v-list-item>
                  </v-list>
                </v-card>
              </v-col>

              <!-- Editor + preview -->
              <v-col cols="12" md="8">
                <v-card v-if="selected" variant="outlined">
                  <v-card-text>
                    <v-row>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="selected.code"
                          label="Code"
                          variant="outlined"
                          density="compact"
                          :disabled="!isNew"
                          :hint="
                            isNew
                              ? 'Unique, e.g. bordereaux_submission'
                              : 'Code is immutable'
                          "
                          persistent-hint
                        />
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="selected.name"
                          label="Name"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12">
                        <v-textarea
                          v-model="selected.description"
                          label="Description"
                          variant="outlined"
                          density="compact"
                          rows="2"
                          auto-grow
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12">
                        <v-text-field
                          v-model="selected.subject_template"
                          label="Subject template"
                          variant="outlined"
                          density="compact"
                          hint="text/template — {{ .scheme_name }}"
                          persistent-hint
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12">
                        <v-textarea
                          v-model="selected.body_template"
                          label="Body template (HTML)"
                          variant="outlined"
                          density="compact"
                          rows="10"
                          auto-grow
                          hint="html/template — auto-escapes values; use {{ .var }}"
                          persistent-hint
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12" md="6">
                        <v-select
                          v-model="selected.status"
                          :items="statuses"
                          label="Status"
                          variant="outlined"
                          density="compact"
                        />
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="selected.attachments_spec"
                          label="Attachments spec (JSON)"
                          variant="outlined"
                          density="compact"
                          hint='[{"kind":"file","path":"..."}]'
                          persistent-hint
                        />
                      </v-col>
                    </v-row>

                    <v-divider class="my-4" />

                    <v-row>
                      <v-col cols="12">
                        <v-textarea
                          v-model="previewVarsJson"
                          label="Preview variables (JSON)"
                          variant="outlined"
                          density="compact"
                          rows="4"
                          auto-grow
                        />
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col>
                        <v-btn
                          variant="outlined"
                          size="small"
                          rounded
                          prepend-icon="mdi-magnify"
                          :loading="previewing"
                          @click="runPreview"
                        >
                          Preview
                        </v-btn>
                      </v-col>
                    </v-row>
                    <v-row v-if="preview.subject || preview.body">
                      <v-col cols="12">
                        <v-card variant="tonal" class="mb-2">
                          <v-card-title class="text-subtitle-2">
                            Subject
                          </v-card-title>
                          <v-card-text class="text-body-2">
                            {{ preview.subject }}
                          </v-card-text>
                        </v-card>
                        <v-card variant="tonal">
                          <v-card-title class="text-subtitle-2">
                            Body
                          </v-card-title>
                          <v-card-text>
                            <div v-html="preview.body" />
                          </v-card-text>
                        </v-card>
                      </v-col>
                    </v-row>
                  </v-card-text>

                  <v-card-actions>
                    <v-btn
                      color="error"
                      variant="text"
                      size="small"
                      :disabled="isNew"
                      @click="onDelete"
                    >
                      Delete
                    </v-btn>
                    <v-spacer />
                    <v-btn
                      variant="outlined"
                      size="small"
                      rounded
                      prepend-icon="mdi-content-save-outline"
                      :loading="saving"
                      @click="onSave"
                    >
                      {{ isNew ? 'Create' : 'Save' }}
                    </v-btn>
                  </v-card-actions>
                </v-card>

                <v-card
                  v-else
                  variant="outlined"
                  class="pa-6 text-center text-medium-emphasis"
                >
                  Select a template on the left or create a new one.
                </v-card>
              </v-col>
            </v-row>
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useFlashStore } from '@/renderer/store/flash'

interface Template {
  code: string
  name: string
  description: string
  subject_template: string
  body_template: string
  attachments_spec: string
  status: 'draft' | 'active'
}

const flash = useFlashStore()
const templates = ref<Template[]>([])
const selected = ref<Template | null>(null)
const isNew = ref(false)
const saving = ref(false)
const seeding = ref(false)
const previewing = ref(false)
const previewVarsJson = ref('{\n  "scheme_name": "Example Scheme"\n}')
const preview = ref<{ subject: string; body: string }>({
  subject: '',
  body: ''
})

const statuses = [
  { title: 'Draft', value: 'draft' },
  { title: 'Active', value: 'active' }
]

const load = async () => {
  try {
    const { data } = await GroupPricingService.listEmailTemplates()
    templates.value = data || []
  } catch (err: any) {
    flash.show(
      err?.response?.data?.error || 'Failed to load templates',
      'error'
    )
  }
}

const selectTemplate = async (code: string) => {
  try {
    const { data } = await GroupPricingService.getEmailTemplate(code)
    selected.value = { ...blankTemplate(), ...data }
    isNew.value = false
    preview.value = { subject: '', body: '' }
  } catch (err: any) {
    flash.show(err?.response?.data?.error || 'Failed to load template', 'error')
  }
}

const createNew = () => {
  selected.value = blankTemplate()
  isNew.value = true
  preview.value = { subject: '', body: '' }
}

const blankTemplate = (): Template => ({
  code: '',
  name: '',
  description: '',
  subject_template: '',
  body_template: '',
  attachments_spec: '',
  status: 'draft'
})

const onSave = async () => {
  if (!selected.value) return
  if (!selected.value.code) {
    flash.show('Code is required', 'warning')
    return
  }
  saving.value = true
  try {
    if (isNew.value) {
      const { data } = await GroupPricingService.createEmailTemplate(
        selected.value
      )
      templates.value.push(data)
      isNew.value = false
    } else {
      const { data } = await GroupPricingService.updateEmailTemplate(
        selected.value.code,
        selected.value
      )
      const idx = templates.value.findIndex((t) => t.code === data.code)
      if (idx !== -1) templates.value[idx] = data
    }
    flash.show('Template saved', 'success')
  } catch (err: any) {
    flash.show(err?.response?.data?.error || 'Failed to save', 'error')
  } finally {
    saving.value = false
  }
}

const onDelete = async () => {
  if (!selected.value || isNew.value) return
  if (!confirm(`Delete template "${selected.value.code}"?`)) return
  try {
    await GroupPricingService.deleteEmailTemplate(selected.value.code)
    templates.value = templates.value.filter(
      (t) => t.code !== selected.value!.code
    )
    selected.value = null
    flash.show('Template deleted', 'success')
  } catch (err: any) {
    flash.show(err?.response?.data?.error || 'Failed to delete', 'error')
  }
}

const runPreview = async () => {
  if (!selected.value) return
  previewing.value = true
  try {
    let vars: Record<string, any> = {}
    if (previewVarsJson.value.trim()) {
      try {
        vars = JSON.parse(previewVarsJson.value)
      } catch {
        flash.show('Preview variables must be valid JSON', 'warning')
        previewing.value = false
        return
      }
    }
    const payload = {
      vars,
      subject_template: selected.value.subject_template,
      body_template: selected.value.body_template
    }
    const { data } = await GroupPricingService.previewEmailTemplate(
      selected.value.code || 'new',
      payload
    )
    preview.value = data
  } catch (err: any) {
    flash.show(err?.response?.data?.error || 'Preview failed', 'error')
  } finally {
    previewing.value = false
  }
}

// The starter set covers the immediate use cases we know about today. New
// templates can be added any time via the New Template button.
const starterTemplates: Template[] = [
  {
    code: 'system_test',
    name: 'System Test',
    description: 'Sent by the "Send Test Email" button in Email Settings.',
    subject_template: 'AART — Test email',
    body_template:
      '<p>Hi {{ .user_name }},</p>' +
      '<p>This is a test email from AART confirming your SMTP settings are correct.</p>' +
      '<p>If you received this, email delivery is working.</p>',
    attachments_spec: '',
    status: 'active'
  },
  {
    code: 'bordereaux_submission',
    name: 'Bordereaux Submission',
    description: 'Sent to the scheme contact when a bordereaux is submitted.',
    subject_template: 'Bordereaux for {{ .scheme_name }} — {{ .period }}',
    body_template:
      '<p>Hi {{ .contact_name }},</p>' +
      '<p>Please find attached the bordereaux for <strong>{{ .scheme_name }}</strong> covering the period <strong>{{ .period }}</strong>.</p>' +
      '<p>Kind regards,<br/>{{ .sender_name }}</p>',
    attachments_spec: '',
    status: 'active'
  }
]

const seedStarterTemplates = async () => {
  seeding.value = true
  try {
    for (const t of starterTemplates) {
      try {
        await GroupPricingService.createEmailTemplate(t)
      } catch (err: any) {
        // Ignore "already exists" so seeding is idempotent.
        if (!/already exists/i.test(err?.response?.data?.error || '')) throw err
      }
    }
    await load()
    flash.show('Starter templates seeded', 'success')
  } catch (err: any) {
    flash.show(err?.response?.data?.error || 'Seeding failed', 'error')
  } finally {
    seeding.value = false
  }
}

onMounted(load)
</script>
