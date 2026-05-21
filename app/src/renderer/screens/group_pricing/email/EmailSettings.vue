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
                <span class="headline">Email Settings</span>
              </div>
              <div class="d-flex align-center gap-2">
                <v-btn
                  class="mr-2"
                  variant="outlined"
                  size="small"
                  rounded
                  prepend-icon="mdi-send-check-outline"
                  :loading="testing"
                  :disabled="!canTest"
                  @click="onTest"
                >
                  Send Test Email
                </v-btn>
                <v-btn
                  size="small"
                  variant="outlined"
                  rounded
                  prepend-icon="mdi-content-save-outline"
                  :loading="saving"
                  @click="onSave"
                >
                  Save
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <v-alert
              v-if="loadError"
              type="error"
              variant="tonal"
              class="mb-4"
              closable
              @click:close="loadError = ''"
            >
              {{ loadError }}
            </v-alert>

            <v-form ref="formRef">
              <v-row>
                <v-col cols="12" md="6">
                  <v-select
                    v-model="form.provider"
                    :items="providers"
                    label="Provider"
                    variant="outlined"
                    density="compact"
                  />
                </v-col>
              </v-row>

              <template v-if="form.provider === 'smtp'">
                <v-row>
                  <v-col cols="12" md="8">
                    <v-text-field
                      v-model="form.host"
                      label="SMTP Host"
                      placeholder="smtp-relay.brevo.com"
                      variant="outlined"
                      density="compact"
                      :rules="[required]"
                    />
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-text-field
                      v-model.number="form.port"
                      label="Port"
                      type="number"
                      variant="outlined"
                      density="compact"
                      :rules="[required]"
                    />
                  </v-col>
                </v-row>

                <v-row>
                  <v-col cols="12" md="4">
                    <v-select
                      v-model="form.tls_mode"
                      :items="tlsModes"
                      label="TLS Mode"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-text-field
                      v-model="form.auth_user"
                      label="Auth Username"
                      variant="outlined"
                      density="compact"
                      autocomplete="off"
                    />
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-text-field
                      v-model="form.auth_password"
                      :label="
                        hasPassword
                          ? 'Password (leave blank to keep)'
                          : 'Password'
                      "
                      :type="showPassword ? 'text' : 'password'"
                      :append-inner-icon="
                        showPassword ? 'mdi-eye-off' : 'mdi-eye'
                      "
                      variant="outlined"
                      density="compact"
                      autocomplete="new-password"
                      @click:append-inner="showPassword = !showPassword"
                    />
                  </v-col>
                </v-row>
              </template>

              <template v-else>
                <v-row>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.graph_tenant_id"
                      label="Tenant ID"
                      placeholder="Directory (tenant) ID"
                      variant="outlined"
                      density="compact"
                      :rules="[required]"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.graph_client_id"
                      label="Client ID"
                      placeholder="Application (client) ID"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                </v-row>

                <v-row>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.graph_client_secret"
                      :label="
                        hasGraphSecret
                          ? 'Client Secret (leave blank to keep)'
                          : 'Client Secret'
                      "
                      :type="showGraphSecret ? 'text' : 'password'"
                      :append-inner-icon="
                        showGraphSecret ? 'mdi-eye-off' : 'mdi-eye'
                      "
                      variant="outlined"
                      density="compact"
                      autocomplete="new-password"
                      @click:append-inner="showGraphSecret = !showGraphSecret"
                    />
                  </v-col>
                </v-row>
              </template>

              <v-row>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model="form.from_address"
                    label="From Address"
                    placeholder="bordereaux@your-insurer.com"
                    variant="outlined"
                    density="compact"
                    :rules="[required, emailRule]"
                  />
                </v-col>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model="form.from_name"
                    label="From Display Name"
                    placeholder="Your Insurer — Group Pricing"
                    variant="outlined"
                    density="compact"
                  />
                </v-col>
              </v-row>

              <v-row>
                <v-col cols="12" md="6">
                  <v-text-field
                    v-model="form.reply_to"
                    label="Reply-To (optional)"
                    placeholder="no-reply@your-insurer.com"
                    variant="outlined"
                    density="compact"
                    :rules="form.reply_to ? [emailRule] : []"
                  />
                </v-col>
              </v-row>

              <v-row v-if="lastUpdatedBy">
                <v-col>
                  <p class="text-caption text-medium-emphasis">
                    Last updated by {{ lastUpdatedBy }}.
                  </p>
                </v-col>
              </v-row>
            </v-form>
          </template>
        </base-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useFlashStore } from '@/renderer/store/flash'

const flash = useFlashStore()

interface SettingsForm {
  provider: 'smtp' | 'microsoft_graph'
  host: string
  port: number
  tls_mode: 'starttls' | 'tls' | 'none'
  auth_user: string
  auth_password: string
  graph_tenant_id: string
  graph_client_id: string
  graph_client_secret: string
  from_address: string
  from_name: string
  reply_to: string
}

const form = ref<SettingsForm>({
  provider: 'smtp',
  host: '',
  port: 587,
  tls_mode: 'starttls',
  auth_user: '',
  auth_password: '',
  graph_tenant_id: '',
  graph_client_id: '',
  graph_client_secret: '',
  from_address: '',
  from_name: '',
  reply_to: ''
})

const providers = [
  { title: 'SMTP relay (SendGrid, Postmark, Brevo…)', value: 'smtp' },
  { title: 'Microsoft 365 (Graph API)', value: 'microsoft_graph' }
]

const tlsModes = [
  { title: 'STARTTLS (recommended for Office 365)', value: 'starttls' },
  { title: 'Implicit TLS', value: 'tls' },
  { title: 'None (insecure)', value: 'none' }
]

const formRef = ref<any>(null)
const saving = ref(false)
const testing = ref(false)
const loadError = ref('')
const hasPassword = ref(false)
const hasGraphSecret = ref(false)
const lastUpdatedBy = ref('')
const showPassword = ref(false)
const showGraphSecret = ref(false)

const required = (v: any) => (!!v && String(v).trim() !== '') || 'Required'
const emailRule = (v: string) =>
  !v || /.+@.+\..+/.test(v) || 'Invalid email address'

const canTest = computed(() => {
  if (form.value.provider === 'microsoft_graph') {
    return !!(
      form.value.from_address &&
      form.value.graph_tenant_id &&
      form.value.graph_client_id &&
      (hasGraphSecret.value || form.value.graph_client_secret)
    )
  }
  return !!(
    form.value.host &&
    form.value.from_address &&
    (hasPassword.value || form.value.auth_password)
  )
})

const load = async () => {
  loadError.value = ''
  try {
    const { data } = await GroupPricingService.getEmailSettings()
    form.value = {
      provider: data.provider || 'smtp',
      host: data.host || '',
      port: data.port || 587,
      tls_mode: data.tls_mode || 'starttls',
      auth_user: data.auth_user || '',
      auth_password: '',
      graph_tenant_id: data.graph_tenant_id || '',
      graph_client_id: data.graph_client_id || '',
      graph_client_secret: '',
      from_address: data.from_address || '',
      from_name: data.from_name || '',
      reply_to: data.reply_to || ''
    }
    hasPassword.value = !!data.has_password
    hasGraphSecret.value = !!data.has_graph_secret
    lastUpdatedBy.value = data.updated_by || ''
  } catch (err: any) {
    if (err?.response?.status === 404) {
      // No settings yet — leave defaults.
      return
    }
    loadError.value =
      err?.response?.data?.error ||
      err.message ||
      'Failed to load email settings'
  }
}

const onSave = async () => {
  const valid = await formRef.value?.validate?.()
  if (valid && valid.valid === false) return
  saving.value = true
  try {
    const payload = { ...form.value }
    if (!payload.auth_password) delete (payload as any).auth_password
    if (!payload.graph_client_secret)
      delete (payload as any).graph_client_secret
    const { data } = await GroupPricingService.saveEmailSettings(payload)
    hasPassword.value = !!data.has_password
    hasGraphSecret.value = !!data.has_graph_secret
    lastUpdatedBy.value = data.updated_by || ''
    form.value.auth_password = ''
    form.value.graph_client_secret = ''
    flash.show('Email settings saved', 'success')
  } catch (err: any) {
    flash.show(
      err?.response?.data?.error || err.message || 'Failed to save',
      'error'
    )
  } finally {
    saving.value = false
  }
}

const onTest = async () => {
  testing.value = true
  try {
    const { data } = await GroupPricingService.sendTestEmail()
    flash.show(`Test email queued (outbox id ${data.outbox_id})`, 'success')
  } catch (err: any) {
    flash.show(
      err?.response?.data?.error ||
        err.message ||
        'Failed to queue test email — is a template with code "system_test" configured?',
      'error'
    )
  } finally {
    testing.value = false
  }
}

onMounted(load)
</script>
