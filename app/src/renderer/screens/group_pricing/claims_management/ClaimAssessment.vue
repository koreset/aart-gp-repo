<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <base-card :show-actions="false">
          <template #header>
            <div class="d-flex align-center">
              <v-btn
                icon="mdi-arrow-left"
                variant="text"
                class="mr-2"
                @click="goBack"
              />
              <span class="headline">
                Claim Assessment
                <template v-if="claim?.claim_number">
                  - {{ claim.claim_number }}
                </template>
              </span>
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
                  Back to Claim Details
                </v-btn>
              </v-col>
            </v-row>

            <v-row v-else>
              <!-- Assessment Form (left) -->
              <v-col cols="12" md="7">
                <claim-assessment-form
                  :claim="claim"
                  @save="handleSave"
                  @cancel="goBack"
                />
              </v-col>

              <!-- Documents Panel (right) -->
              <v-col cols="12" md="5">
                <v-card variant="outlined" class="mb-4">
                  <v-card-title class="bg-primary text-white"
                    >Supporting Documents</v-card-title
                  >
                  <v-card-text>
                    <v-list v-if="supportingDocuments.length" density="compact">
                      <v-list-item
                        v-for="doc in supportingDocuments"
                        :key="doc.id ?? doc.document_name"
                        :active="previewDoc?.id === doc.id"
                        class="cursor-pointer"
                        @click="previewDocument(doc)"
                      >
                        <template #prepend>
                          <v-icon>{{ getDocumentIcon(doc) }}</v-icon>
                        </template>
                        <v-list-item-title>{{
                          doc.document_name || doc.filename || doc.name
                        }}</v-list-item-title>
                        <v-list-item-subtitle v-if="doc.uploaded_at">{{
                          formatDate(doc.uploaded_at)
                        }}</v-list-item-subtitle>
                      </v-list-item>
                    </v-list>
                    <div v-else class="text-body-2 text-medium-emphasis py-2">
                      No supporting documents uploaded yet.
                    </div>
                  </v-card-text>
                </v-card>

                <v-card v-if="previewDoc" variant="outlined" class="mb-4">
                  <v-card-title class="bg-primary text-white">
                    <div
                      class="d-flex justify-space-between align-center w-100"
                    >
                      <span>{{
                        previewDoc.document_name ||
                        previewDoc.filename ||
                        previewDoc.name
                      }}</span>
                      <v-btn
                        icon="mdi-close"
                        variant="text"
                        color="white"
                        size="small"
                        @click="clearPreview"
                      />
                    </div>
                  </v-card-title>
                  <v-card-text class="pa-0" style="height: 480px">
                    <div
                      v-if="previewLoading"
                      class="d-flex flex-column justify-center align-center"
                      style="height: 100%"
                    >
                      <v-progress-circular indeterminate color="primary" />
                      <div class="mt-2 text-body-2">Loading document...</div>
                    </div>
                    <div
                      v-else-if="previewError"
                      class="d-flex flex-column justify-center align-center text-error pa-4"
                      style="height: 100%"
                    >
                      <v-icon size="48" color="error">mdi-alert-circle</v-icon>
                      <div class="mt-2 text-body-2 text-center">{{
                        previewError
                      }}</div>
                    </div>
                    <iframe
                      v-else-if="isPdfDocument(previewDoc) && previewUrl"
                      :src="previewUrl"
                      style="width: 100%; height: 100%; border: none"
                      title="Document Preview"
                    />
                    <div
                      v-else-if="isImageDocument(previewDoc) && previewUrl"
                      class="d-flex justify-center align-center"
                      style="height: 100%; background-color: #f5f5f5"
                    >
                      <v-img
                        :src="previewUrl"
                        :alt="
                          previewDoc?.document_name ||
                          previewDoc?.filename ||
                          previewDoc?.name
                        "
                        contain
                        max-height="460px"
                        max-width="100%"
                      />
                    </div>
                    <div
                      v-else
                      class="d-flex flex-column justify-center align-center text-medium-emphasis pa-4"
                      style="height: 100%"
                    >
                      <v-icon size="48" color="grey">mdi-file-document</v-icon>
                      <div class="mt-2 text-body-2"
                        >Preview not available for this file type</div
                      >
                    </div>
                  </v-card-text>
                </v-card>

                <v-card variant="outlined">
                  <v-card-title class="bg-primary text-white"
                    >Missing Documents</v-card-title
                  >
                  <v-card-text>
                    <v-list v-if="missingDocs.length" density="compact">
                      <v-list-item
                        v-for="missing in missingDocs"
                        :key="missing.code"
                      >
                        <template #prepend>
                          <v-icon>mdi-file-document-outline</v-icon>
                        </template>
                        <v-list-item-title>{{
                          missing.name
                        }}</v-list-item-title>
                        <template v-if="missing.required" #append>
                          <v-chip color="error" size="x-small" variant="tonal"
                            >Required</v-chip
                          >
                        </template>
                      </v-list-item>
                    </v-list>
                    <div v-else class="text-body-2 text-medium-emphasis py-2">
                      All required documents are present.
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
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
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import ClaimAssessmentForm from './components/ClaimAssessmentForm.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { claimDocumentTypes } from './claimDocumentTypes'

const props = defineProps<{
  id: string | number
}>()

const router = useRouter()

const loading = ref(false)
const saving = ref(false)
const claim = ref<any>(null)

const previewDoc = ref<any>(null)
const previewUrl = ref('')
const previewLoading = ref(false)
const previewError = ref('')
const previewBlobUrls = ref<string[]>([])

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref<'success' | 'error' | 'info'>('success')

const showSnackbar = (
  message: string,
  color: 'success' | 'error' | 'info' = 'success'
) => {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

const claimId = (): number => {
  return typeof props.id === 'string' ? parseInt(props.id, 10) : props.id
}

const supportingDocuments = computed<any[]>(() => {
  const attachments = claim.value?.attachments
  return Array.isArray(attachments) ? attachments : []
})

const missingDocs = computed(() => {
  if (!claim.value) return []
  const required = claimDocumentTypes[claim.value.benefit_code] || []
  const uploaded = new Set(
    supportingDocuments.value.map((d) => d.document_type)
  )
  return required.filter((d) => !uploaded.has(d.code))
})

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

const goBack = () => {
  router.push({
    name: 'group-pricing-claim-details',
    params: { id: claimId() }
  })
}

const handleSave = async (assessmentData: any) => {
  if (!claim.value || saving.value) return

  saving.value = true
  try {
    const payload = {
      claim_id: claim.value.id,
      ...assessmentData,
      assessment_timestamp:
        assessmentData.assessment_timestamp || new Date().toISOString()
    }

    await GroupPricingService.createClaimAssessment(claim.value.id, payload)

    const outcome = assessmentData.assessment_outcome
    let newStatus = 'under_assessment'
    if (outcome === 'refer_to_committee') newStatus = 'referred_to_committee'
    else if (outcome === 'committee_approved') newStatus = 'approved'
    else if (outcome === 'committee_declined') newStatus = 'declined'
    else if (outcome === 'requires_additional_info')
      newStatus = 'additional_info_required'

    const updatedClaim = {
      ...claim.value,
      status: newStatus,
      last_updated: new Date().toISOString()
    }
    await GroupPricingService.updateClaim(claim.value.id, updatedClaim)

    showSnackbar('Assessment saved', 'success')
    router.push({
      name: 'group-pricing-claim-details',
      params: { id: claim.value.id }
    })
  } catch (error) {
    console.error('Failed to save assessment:', error)
    showSnackbar('Failed to save assessment. Please try again.', 'error')
  } finally {
    saving.value = false
  }
}

const isPdfDocument = (doc: any): boolean => {
  if (!doc) return false
  const contentType = doc.content_type || doc.type || ''
  const filename = doc.document_name || doc.filename || doc.name || ''
  return contentType.includes('pdf') || filename.toLowerCase().endsWith('.pdf')
}

const isImageDocument = (doc: any): boolean => {
  if (!doc) return false
  const contentType = doc.content_type || doc.type || ''
  const filename = doc.document_name || doc.filename || doc.name || ''
  return (
    contentType.includes('image') ||
    /\.(jpg|jpeg|png|gif|bmp|webp)$/i.test(filename)
  )
}

const getDocumentIcon = (doc: any): string => {
  if (isPdfDocument(doc)) return 'mdi-file-pdf-box'
  if (isImageDocument(doc)) return 'mdi-image'
  return 'mdi-file-document'
}

const formatDate = (value: string | Date | undefined): string => {
  if (!value) return ''
  const date = typeof value === 'string' ? new Date(value) : value
  if (Number.isNaN(date.getTime())) return String(value)
  return date.toLocaleDateString()
}

const clearPreview = () => {
  previewDoc.value = null
  previewUrl.value = ''
  previewError.value = ''
}

const resolveDocumentUrl = (doc: any): string | null => {
  const baseUrl: string = window.mainApi.sendSync('msgGetBaseUrl')
  const normalizedBaseUrl = baseUrl.endsWith('/')
    ? baseUrl.slice(0, -1)
    : baseUrl
  if (doc.viewer_url) return normalizedBaseUrl + doc.viewer_url
  if (doc.path) return doc.path
  if (doc.file_path) return doc.file_path
  return null
}

const previewDocument = async (doc: any) => {
  previewDoc.value = doc
  previewUrl.value = ''
  previewError.value = ''
  previewLoading.value = true

  const url = resolveDocumentUrl(doc)
  if (!url) {
    previewLoading.value = false
    previewError.value = 'Document URL is not available'
    return
  }

  try {
    const token: string = window.mainApi.sendSync('msgGetAccessToken')
    const response = await fetch(url, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!response.ok) {
      throw new Error(`${response.status} ${response.statusText}`)
    }
    const blob = await response.blob()
    const blobUrl = URL.createObjectURL(blob)
    previewBlobUrls.value.push(blobUrl)
    previewUrl.value = blobUrl
  } catch (error: any) {
    console.error('Error fetching document:', error)
    previewError.value = `Failed to load document: ${error?.message ?? 'unknown error'}`
  } finally {
    previewLoading.value = false
  }
}

onMounted(() => {
  loadClaim()
})

onBeforeUnmount(() => {
  previewBlobUrls.value.forEach((url) => URL.revokeObjectURL(url))
  previewBlobUrls.value = []
})
</script>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}
</style>
