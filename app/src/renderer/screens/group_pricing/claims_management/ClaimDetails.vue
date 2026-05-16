<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex justify-space-between align-center">
              <div class="d-flex align-center">
                <v-btn
                  icon="mdi-arrow-left"
                  variant="text"
                  class="mr-2"
                  @click="goBack"
                />
                <span class="headline">
                  Claim Details
                  <template v-if="selectedClaim?.claim_number">
                    - {{ selectedClaim.claim_number }}
                  </template>
                </span>
              </div>
            </div>
          </template>
          <template #default>
            <v-row v-if="loading && !selectedClaim">
              <v-col cols="12" class="text-center py-8">
                <v-progress-circular indeterminate color="primary" />
                <div class="mt-2 text-body-2 text-medium-emphasis"
                  >Loading claim...</div
                >
              </v-col>
            </v-row>

            <v-row v-else-if="!selectedClaim">
              <v-col cols="12" class="text-center py-8">
                <v-icon size="64" color="grey-lighten-1" class="mb-4"
                  >mdi-file-document-remove-outline</v-icon
                >
                <h3 class="text-h6 text-grey-darken-1 mb-2">Claim not found</h3>
                <v-btn color="primary" variant="outlined" @click="goBack">
                  Back to Claims Management
                </v-btn>
              </v-col>
            </v-row>

            <claim-detail-view
              v-else
              :claim="selectedClaim"
              @update="handleClaimUpdate"
              @close="goBack"
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :timeout="3000" :color="snackbarColor">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import ClaimDetailView from './components/ClaimDetailView.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const props = defineProps<{
  id: string | number
}>()

const router = useRouter()

const loading = ref(false)
const selectedClaim = ref<any>(null)

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const showSnackbar = (message: string, color: string = 'success') => {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

const claimId = (): number => {
  return typeof props.id === 'string' ? parseInt(props.id, 10) : props.id
}

const loadClaim = async () => {
  const id = claimId()
  if (!id || Number.isNaN(id)) {
    showSnackbar('Invalid claim id', 'error')
    return
  }

  loading.value = true
  try {
    const response = await GroupPricingService.getClaim(id)
    selectedClaim.value = response.data
  } catch (error) {
    console.error('Error loading claim:', error)
    showSnackbar('Failed to load claim', 'error')
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push({ name: 'group-pricing-claims-management' })
}

const handleClaimUpdate = async (updatedClaim: any) => {
  try {
    await GroupPricingService.updateClaim(updatedClaim.id, updatedClaim)
    showSnackbar('Claim updated successfully', 'success')
    await loadClaim()
  } catch (error) {
    console.error('Error updating claim:', error)
    showSnackbar('Error updating claim. Please try again.', 'error')
  }
}

onMounted(() => {
  loadClaim()
})
</script>
