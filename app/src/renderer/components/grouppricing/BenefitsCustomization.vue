<template>
  <base-card :show-actions="false">
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-clipboard-list" class="mr-3" color="white" />
        <span>Benefits Customization</span>
        <v-spacer />
        <v-chip
          color="info"
          variant="outlined"
          size="small"
          class="white--text"
        >
          Customize Plan Benefits
        </v-chip>
      </div>
    </template>

    <template #default>
      <h5 class="section-title">
        <v-icon icon="mdi-cog-outline" size="small" class="mr-2" />
        Benefit Configuration
      </h5>

      <p class="text-body-2 text-medium-emphasis mb-4">
        Configure custom benefit names and mapping for your insurance products.
        These settings will be applied across all quotations and policy
        documents.
      </p>

      <benefit-mapper :loading="isLoading" @submit="handleSubmit" />
    </template>
  </base-card>

  <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="3000">
    {{ snackbarMessage }}
    <template #actions>
      <v-btn color="white" variant="text" @click="hideNotification"
        >Close</v-btn
      >
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import BenefitMapper from '@/renderer/components/grouppricing/BenefitMapper.vue'
import { useNotifications } from '@/renderer/composables/useNotifications'
import { useErrorHandler } from '@/renderer/composables/useErrorHandler'
import { VALIDATION_MESSAGES } from '@/renderer/constants/metadata'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '../BaseCard.vue'

const emit = defineEmits<{
  benefitsSaved: []
}>()

const {
  snackbar,
  snackbarMessage,
  snackbarColor,
  showSuccess,
  showError,
  hideNotification
} = useNotifications()
const { handleApiError } = useErrorHandler()
const isLoading = ref(false)

const handleSubmit = async (benefits: any) => {
  try {
    isLoading.value = true
    const res = await GroupPricingService.saveBenefitMap(benefits)

    if (res.status === 201) {
      showSuccess(VALIDATION_MESSAGES.BENEFITS_SAVED)
      emit('benefitsSaved')
    } else {
      showError(VALIDATION_MESSAGES.BENEFITS_FAILED)
    }
  } catch (error: any) {
    const errorMessage = handleApiError(
      error,
      VALIDATION_MESSAGES.BENEFITS_FAILED
    )
    showError(errorMessage)
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.benefits-card {
  border-radius: 16px;
  overflow: hidden;
}

.section-header {
  background: linear-gradient(135deg, #ff9800 0%, #f57c00 100%);
  color: white;
  padding: 1.5rem;
  font-size: 1.25rem;
  font-weight: 600;
}

.form-section {
  padding: 2rem;
  background: #fefefe;
  border-radius: 12px;
  border-left: 4px solid #ff9800;
}

.section-title {
  color: #ff9800;
  font-weight: 600;
  font-size: 1rem;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
}
</style>
