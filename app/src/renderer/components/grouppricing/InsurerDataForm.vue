<template>
  <base-card>
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-domain" class="mr-3" />
        <span>Insurer Information</span>
        <v-spacer />
        <v-chip
          :color="formData.name ? 'success' : 'warning'"
          variant="outlined"
          size="small"
        >
          {{ formData.name ? 'Configured' : 'Pending Setup' }}
        </v-chip>
      </div>
    </template>
    <template #default>
      <h5 class="section-title">
        <v-icon icon="mdi-information-outline" size="small" class="mr-2" />
        Basic Information
      </h5>
      <v-row>
        <v-col cols="12" md="6">
          <v-text-field
            v-model="formData.name"
            variant="outlined"
            density="compact"
            label="Insurer Name"
            placeholder="Enter insurer name"
            prepend-inner-icon="mdi-domain"
          />
        </v-col>
        <v-col cols="12" md="6">
          <v-select
            v-model="formData.year_end_month"
            placeholder="Select year end month"
            label="Year End Month"
            :items="YEAR_END_MONTHS"
            variant="outlined"
            density="compact"
            item-title="name"
            item-value="month_number"
            prepend-inner-icon="mdi-calendar"
          />
        </v-col>
      </v-row>
      <h5 class="section-title">
        <v-icon icon="mdi-map-marker-outline" size="small" class="mr-2" />
        Address Information
      </h5>
      <v-row>
        <v-col cols="12" md="6">
          <v-text-field
            v-model="formData.address_line_1"
            variant="outlined"
            density="compact"
            label="Address Line 1"
            placeholder="Enter primary address"
            prepend-inner-icon="mdi-home-outline"
          />
        </v-col>
        <v-col cols="12" md="6">
          <v-text-field
            v-model="formData.address_line_2"
            variant="outlined"
            density="compact"
            label="Address Line 2"
            placeholder="Enter secondary address (optional)"
            prepend-inner-icon="mdi-home-plus-outline"
          />
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="formData.city"
            variant="outlined"
            density="compact"
            label="City"
            placeholder="Enter city"
            prepend-inner-icon="mdi-city"
          />
        </v-col>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="formData.province"
            variant="outlined"
            density="compact"
            label="Province"
            placeholder="Enter province"
            prepend-inner-icon="mdi-map"
          />
        </v-col>
        <v-col cols="12" md="4">
          <v-text-field
            v-model="formData.post_code"
            variant="outlined"
            density="compact"
            label="Postal Code"
            placeholder="Enter postal code"
            prepend-inner-icon="mdi-mailbox"
          />
        </v-col>
      </v-row>
      <h5 class="section-title">
        <v-icon icon="mdi-phone-outline" size="small" class="mr-2" />
        Contact Information
      </h5>
      <v-row>
        <v-col cols="12" md="6">
          <v-text-field
            v-model="formData.email"
            variant="outlined"
            density="compact"
            label="Email Address"
            placeholder="Enter email address"
            prepend-inner-icon="mdi-email-outline"
            type="email"
          />
        </v-col>
        <v-col cols="12" md="6">
          <v-text-field
            v-model="formData.telephone"
            variant="outlined"
            density="compact"
            label="Telephone"
            placeholder="Enter phone number"
            prepend-inner-icon="mdi-phone"
          />
        </v-col>
      </v-row>
      <h5 class="section-title">
        <v-icon icon="mdi-image-outline" size="small" class="mr-2" />
        Company Logo
      </h5>
      <v-row>
        <v-col cols="12" md="6">
          <v-file-input
            v-model="logoFile"
            label="Upload Company Logo"
            accept="image/*"
            variant="outlined"
            prepend-icon="mdi-camera"
            show-size
            @update:model-value="onFileChange"
          />
        </v-col>
        <v-col cols="12" md="6">
          <v-card v-if="imagePreview" class="logo-preview" elevation="2">
            <v-card-text class="d-flex justify-center align-center pa-4">
              <img
                :src="imagePreview"
                alt="Logo Preview"
                class="preview-image"
              />
            </v-card-text>
          </v-card>
          <div v-else class="logo-placeholder">
            <v-icon icon="mdi-image-plus" size="48" color="grey-lighten-2" />
            <p class="text-caption text-medium-emphasis mt-2"
              >Logo preview will appear here</p
            >
          </div>
        </v-col>
      </v-row>
      <v-divider class="my-6"></v-divider>
      <h5 class="section-title">
        <v-icon icon="mdi-text-box-outline" size="small" class="mr-2" />
        Document Content
      </h5>
      <v-row>
        <v-col cols="12">
          <v-textarea
            v-model="formData.introductory_text"
            label="Quotation Introductory Text"
            placeholder="Enter introductory text for quotations"
            variant="outlined"
            rows="3"
            auto-grow
            prepend-inner-icon="mdi-text"
          />
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12">
          <v-textarea
            v-model="formData.general_provisions_text"
            label="Underwriting and General Provisions Text"
            placeholder="Enter underwriting and general provisions text"
            variant="outlined"
            rows="3"
            auto-grow
            prepend-inner-icon="mdi-file-document-outline"
          />
        </v-col>
      </v-row>
    </template>
    <template #actions>
      <v-spacer />
      <v-btn
        color="primary"
        size="large"
        prepend-icon="mdi-content-save"
        :loading="isSaving"
        @click="handleSave"
      >
        Save Insurer Details
      </v-btn>
    </template>
  </base-card>
  <confirm-dialog ref="confirmActionDialog" />
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import { useNotifications } from '@/renderer/composables/useNotifications'
import { useImageUpload } from '@/renderer/composables/useImageUpload'
import { YEAR_END_MONTHS } from '@/renderer/constants/metadata'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import ConfirmDialog from '../ConfirmDialog.vue'
import BaseCard from '../BaseCard.vue'
const emit = defineEmits<{
  saved: []
}>()

const { showSuccess, showError } = useNotifications()
const { imagePreview, logoFile, handleFileChange } = useImageUpload()

const isSaving = ref(false)
const confirmActionDialog: any = ref()

// Simple form data without validation
const formData = reactive({
  name: '',
  address_line_1: '',
  address_line_2: '',
  city: '',
  province: '',
  post_code: '',
  telephone: '',
  email: '',
  year_end_month: null,
  introductory_text: '',
  general_provisions_text: ''
})

const onFileChange = () => {
  console.log('File change detected, logoFile.value:', logoFile.value)
  try {
    if (logoFile.value) {
      console.log('File available:', logoFile.value)
      handleFileChange(logoFile.value)
      console.log(
        'After handleFileChange, imagePreview.value:',
        imagePreview.value
      )
    } else {
      console.log('No file selected')
    }
  } catch (error: any) {
    console.error('Error in onFileChange:', error)
    showError(error.message)
  }
}

const handleSave = async () => {
  try {
    const confirmed = await confirmActionDialog.value.open(
      'Confirm Save',
      'Are you sure you want to save the insurer details?'
    )
    if (!confirmed) {
      return
    }
  } catch (error) {
    console.error('Confirmation dialog failed to open:', error)
    return
  }

  isSaving.value = true
  try {
    const payload = new FormData()

    // Add form data
    Object.keys(formData).forEach((key) => {
      if (formData[key] !== null && formData[key] !== undefined) {
        payload.append(key, formData[key])
      }
    })

    // Add logo if selected
    if (logoFile.value) {
      payload.append('logo', logoFile.value)
    }

    const response = await GroupPricingService.createInsurer(payload)

    if (response.status === 201) {
      showSuccess('Insurer information saved successfully')
      emit('saved')
    }
  } catch (error: any) {
    showError('Failed to save insurer information')
  } finally {
    isSaving.value = false
  }
}

onMounted(async () => {
  try {
    const response = await GroupPricingService.getInsurer()
    if (response?.data) {
      Object.assign(formData, response.data)
      if (response.data.logo) {
        imagePreview.value = `data:image/*;base64,${response.data.logo}`
      }
    }
  } catch (error) {
    // Error loading data - continue with empty form
  }
})
</script>

<style scoped>
.insurer-card {
  border-radius: 16px;
  overflow: hidden;
}

.section-header {
  background: linear-gradient(135deg, #1976d2 0%, #1565c0 100%);
  color: white;
  padding: 1.5rem;
  font-size: 1.25rem;
  font-weight: 600;
}

.form-section {
  margin-bottom: 2rem;
  padding: 1rem;
  background: #fafafa;
  border-radius: 12px;
  border-left: 4px solid #1976d2;
}

.section-title {
  color: #1976d2;
  font-weight: 600;
  font-size: 1rem;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
}

.logo-preview {
  border-radius: 12px;
  max-height: 200px;
  overflow: hidden;
}

.preview-image {
  max-width: 100%;
  max-height: 160px;
  object-fit: contain;
  border-radius: 8px;
}

.logo-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  background: #f5f5f5;
  border-radius: 12px;
  border: 2px dashed #e0e0e0;
}

/* Form field enhancements */
:deep(.v-text-field .v-field__outline) {
  border-radius: 8px;
}

:deep(.v-textarea .v-field__outline) {
  border-radius: 8px;
}

:deep(.v-select .v-field__outline) {
  border-radius: 8px;
}

:deep(.v-file-input .v-field__outline) {
  border-radius: 8px;
}

/* Hover effects */
:deep(.v-text-field:hover .v-field__outline) {
  border-color: #1976d2;
}

/* Focus effects */
:deep(.v-text-field.v-field--focused .v-field__outline) {
  border-color: #1976d2;
  border-width: 2px;
}
</style>
