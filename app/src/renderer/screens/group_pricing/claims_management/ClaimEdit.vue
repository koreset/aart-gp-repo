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
                  Edit Claim
                  <template v-if="claim?.claim_number">
                    - {{ claim.claim_number }}
                  </template>
                </span>
              </div>
            </div>
          </template>
          <template #default>
            <v-row v-if="loading && !claim">
              <v-col cols="12" class="text-center py-8">
                <v-progress-circular indeterminate color="primary" />
                <div class="mt-2 text-body-2 text-medium-emphasis"
                  >Loading claim...</div
                >
              </v-col>
            </v-row>

            <v-row v-else-if="!claim">
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

            <v-row v-else-if="!isEditable">
              <v-col cols="12" class="text-center py-8">
                <v-icon size="64" color="warning" class="mb-4"
                  >mdi-lock-outline</v-icon
                >
                <h3 class="text-h6 text-grey-darken-1 mb-2"
                  >This claim can no longer be edited</h3
                >
                <p class="text-body-2 text-medium-emphasis mb-4">
                  Only claims that are still in draft, pending, under
                  assessment, or awaiting additional information can be
                  edited. This claim is currently
                  <strong>{{ formattedStatus }}</strong
                  >.
                </p>
                <v-btn color="primary" variant="outlined" @click="goToDetails">
                  View claim details
                </v-btn>
              </v-col>
            </v-row>

            <claim-registration-form
              v-else
              :schemes="schemes"
              :claim="claim"
              mode="edit"
              @save="handleSave"
              @cancel="goBack"
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import ClaimRegistrationForm from './components/ClaimRegistrationForm.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { isEditableClaimStatus } from '@/renderer/utils/claimStatus'

const props = defineProps<{
  id: string | number
}>()

const router = useRouter()

const loading = ref(false)
const claim = ref<any>(null)
const schemes = ref<any[]>([])

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const isEditable = computed(() => isEditableClaimStatus(claim.value?.status))

const formattedStatus = computed(() => {
  if (!claim.value?.status) return ''
  return claim.value.status
    .split('_')
    .map((part: string) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(' ')
})

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
    claim.value = response.data
  } catch (error) {
    console.error('Error loading claim:', error)
    showSnackbar('Failed to load claim', 'error')
  } finally {
    loading.value = false
  }
}

const loadSchemes = async () => {
  try {
    const response = await GroupPricingService.getSchemesInforce()
    schemes.value = response.data || []
  } catch (error) {
    console.error('Error loading schemes:', error)
    schemes.value = []
  }
}

const handleSave = async (payload: FormData) => {
  const id = claimId()
  loading.value = true
  try {
    await GroupPricingService.updateClaim(id, payload)
    showSnackbar('Claim updated successfully', 'success')
    router.push({
      name: 'group-pricing-claim-details',
      params: { id }
    })
  } catch (error: any) {
    console.error('Error updating claim:', error)
    const status = error?.response?.status
    if (status === 409) {
      showSnackbar(
        'This claim can no longer be edited in its current status.',
        'error'
      )
    } else {
      showSnackbar('Error updating claim. Please try again.', 'error')
    }
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  const id = claimId()
  if (id) {
    router.push({ name: 'group-pricing-claim-details', params: { id } })
  } else {
    router.push({ name: 'group-pricing-claims-management' })
  }
}

const goToDetails = () => {
  router.push({
    name: 'group-pricing-claim-details',
    params: { id: claimId() }
  })
}

onMounted(() => {
  loadClaim()
  loadSchemes()
})
</script>
