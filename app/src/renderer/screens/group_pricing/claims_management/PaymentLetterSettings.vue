<template>
  <v-container>
    <v-row>
      <v-col cols="12" md="10">
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center">
              <v-btn
                icon="mdi-arrow-left"
                variant="text"
                class="mr-2"
                @click="goBack"
              />
              <span class="headline">Payment Letter Settings</span>
            </div>
          </template>

          <template #default>
            <v-alert
              color="primary"
              variant="tonal"
              density="compact"
              icon="mdi-information-outline"
              class="mb-4"
            >
              <div class="text-body-2">
                These details appear on every claim payment confirmation
                letter — letterhead at the top, signatory block at the bottom,
                and intro / closing paragraphs in the body. Tokens you can use
                in the intro / closing text:
                <code v-pre>{{claimant_name}}</code>,
                <code v-pre>{{claim_number}}</code>,
                <code v-pre>{{amount}}</code>,
                <code v-pre>{{paid_at}}</code>,
                <code v-pre>{{member_name}}</code>.
              </div>
            </v-alert>

            <v-card variant="outlined" class="mb-4">
              <v-card-title class="text-subtitle-1 bg-grey-lighten-4"
                >Company</v-card-title
              >
              <v-card-text>
                <v-row dense>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.company_name"
                      label="Company name"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.phone"
                      label="Phone"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.email"
                      label="Email"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.website"
                      label="Website"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12">
                    <v-text-field
                      v-model="form.address_line1"
                      label="Address line 1"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.address_line2"
                      label="Address line 2"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.address_line3"
                      label="Address line 3"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-text-field
                      v-model="form.city"
                      label="City"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-text-field
                      v-model="form.postal_code"
                      label="Postal code"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="4">
                    <v-text-field
                      v-model="form.country"
                      label="Country"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <v-card variant="outlined" class="mb-4">
              <v-card-title class="text-subtitle-1 bg-grey-lighten-4"
                >Branding</v-card-title
              >
              <v-card-text>
                <div class="d-flex align-center mb-3">
                  <div v-if="state.has_logo" class="mr-4">
                    <v-chip color="success" size="small" variant="tonal">
                      <v-icon size="small" class="mr-1">mdi-check</v-icon>
                      Logo on file ({{ state.logo_mime_type || 'image' }})
                    </v-chip>
                  </div>
                  <v-file-input
                    v-model="logoFile"
                    label="Replace logo (PNG, JPEG, SVG, max 2 MB)"
                    accept="image/png,image/jpeg,image/svg+xml"
                    variant="outlined"
                    density="compact"
                    show-size
                    hide-details
                    prepend-icon=""
                    prepend-inner-icon="mdi-image-outline"
                  />
                </div>
              </v-card-text>
            </v-card>

            <v-card variant="outlined" class="mb-4">
              <v-card-title class="text-subtitle-1 bg-grey-lighten-4"
                >Signatory</v-card-title
              >
              <v-card-text>
                <v-row dense>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.signatory_name"
                      label="Signatory name"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="form.signatory_title"
                      label="Signatory title"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12">
                    <div class="d-flex align-center">
                      <div v-if="state.has_signature" class="mr-4">
                        <v-chip color="success" size="small" variant="tonal">
                          <v-icon size="small" class="mr-1">mdi-check</v-icon>
                          Signature image on file
                        </v-chip>
                      </div>
                      <v-file-input
                        v-model="signatureFile"
                        label="Replace signature image (PNG, JPEG, max 2 MB)"
                        accept="image/png,image/jpeg"
                        variant="outlined"
                        density="compact"
                        show-size
                        hide-details
                        prepend-icon=""
                        prepend-inner-icon="mdi-draw-pen"
                      />
                    </div>
                  </v-col>
                </v-row>
              </v-card-text>
            </v-card>

            <v-card variant="outlined" class="mb-4">
              <v-card-title class="text-subtitle-1 bg-grey-lighten-4"
                >Letter body</v-card-title
              >
              <v-card-text>
                <v-textarea
                  v-model="form.letter_intro_template"
                  label="Intro paragraph"
                  variant="outlined"
                  density="compact"
                  rows="3"
                  hint="Appears after the salutation. Leave blank to use the default 'We confirm…' wording."
                  persistent-hint
                  class="mb-3"
                />
                <v-textarea
                  v-model="form.letter_closing_template"
                  label="Closing paragraph"
                  variant="outlined"
                  density="compact"
                  rows="3"
                  hint="Appears just above the signature. Leave blank for the default 'Should you have any queries…' wording."
                  persistent-hint
                />
              </v-card-text>
            </v-card>

            <v-alert
              v-if="saveError"
              type="error"
              variant="tonal"
              density="compact"
              class="mb-3"
              >{{ saveError }}</v-alert
            >

            <div class="d-flex">
              <v-btn variant="text" @click="reset">Reset</v-btn>
              <v-spacer />
              <v-btn
                color="primary"
                :loading="saving"
                prepend-icon="mdi-content-save-outline"
                @click="save"
                >Save settings</v-btn
              >
            </div>
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface SettingsForm {
  company_name: string
  address_line1: string
  address_line2: string
  address_line3: string
  city: string
  postal_code: string
  country: string
  phone: string
  email: string
  website: string
  signatory_name: string
  signatory_title: string
  letter_intro_template: string
  letter_closing_template: string
}

interface SettingsState extends SettingsForm {
  has_logo: boolean
  logo_mime_type: string
  has_signature: boolean
  signature_mime_type: string
  updated_at?: string
  updated_by?: string
}

const router = useRouter()

const blankForm = (): SettingsForm => ({
  company_name: '',
  address_line1: '',
  address_line2: '',
  address_line3: '',
  city: '',
  postal_code: '',
  country: '',
  phone: '',
  email: '',
  website: '',
  signatory_name: '',
  signatory_title: '',
  letter_intro_template: '',
  letter_closing_template: ''
})

const form = ref<SettingsForm>(blankForm())
const state = ref<SettingsState>({
  ...blankForm(),
  has_logo: false,
  logo_mime_type: '',
  has_signature: false,
  signature_mime_type: ''
})

const logoFile = ref<File | File[] | null>(null)
const signatureFile = ref<File | File[] | null>(null)
const saving = ref(false)
const saveError = ref('')
const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref<'success' | 'error'>('success')

function goBack() {
  router.back()
}

function reset() {
  form.value = { ...state.value }
  logoFile.value = null
  signatureFile.value = null
}

onMounted(async () => {
  await load()
})

async function load() {
  try {
    const res = await GroupPricingService.getPaymentLetterSettings()
    const data = unwrap<SettingsState>(res.data)
    state.value = { ...state.value, ...data }
    form.value = {
      company_name: data.company_name || '',
      address_line1: data.address_line1 || '',
      address_line2: data.address_line2 || '',
      address_line3: data.address_line3 || '',
      city: data.city || '',
      postal_code: data.postal_code || '',
      country: data.country || '',
      phone: data.phone || '',
      email: data.email || '',
      website: data.website || '',
      signatory_name: data.signatory_name || '',
      signatory_title: data.signatory_title || '',
      letter_intro_template: data.letter_intro_template || '',
      letter_closing_template: data.letter_closing_template || ''
    }
  } catch (err: any) {
    saveError.value = extractError(err)
  }
}

async function save() {
  saving.value = true
  saveError.value = ''
  try {
    const fd = new FormData()
    for (const [k, v] of Object.entries(form.value)) {
      fd.append(k, v ?? '')
    }
    const logo = unwrapFile(logoFile.value)
    if (logo) fd.append('logo', logo)
    const sig = unwrapFile(signatureFile.value)
    if (sig) fd.append('signature', sig)

    await GroupPricingService.updatePaymentLetterSettings(fd)
    await load()
    logoFile.value = null
    signatureFile.value = null
    snackbarColor.value = 'success'
    snackbarMessage.value = 'Settings saved'
    snackbar.value = true
  } catch (err: any) {
    saveError.value = extractError(err)
    snackbarColor.value = 'error'
    snackbarMessage.value = saveError.value
    snackbar.value = true
  } finally {
    saving.value = false
  }
}

function unwrapFile(v: File | File[] | null): File | null {
  if (!v) return null
  return Array.isArray(v) ? v[0] || null : v
}

function unwrap<T>(payload: any): T {
  if (payload && typeof payload === 'object' && 'data' in payload) {
    return payload.data as T
  }
  return payload as T
}

function extractError(err: any): string {
  return (
    err?.response?.data?.error ||
    err?.response?.data?.message ||
    err?.message ||
    'Unexpected error'
  )
}
</script>
