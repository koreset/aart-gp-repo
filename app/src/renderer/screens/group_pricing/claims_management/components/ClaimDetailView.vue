<template>
  <v-container>
    <v-row>
      <!-- Claim Summary Card -->
      <v-col cols="12" md="9">
        <v-card class="mb-4">
          <v-card-title class="bg-primary text-white">
            <div class="d-flex justify-space-between align-center w-100">
              <span>{{ claim?.claim_number }}</span>
              <div class="d-flex align-center gap-2">
                <v-chip :color="getStatusColor(claim?.status)" size="small">
                  {{ formatStatus(claim?.status) }}
                </v-chip>
                <v-btn
                  icon="mdi-close"
                  variant="text"
                  color="white"
                  size="small"
                  @click="$emit('close')"
                />
              </div>
            </div>
          </v-card-title>
          <v-card-text>
            <v-row class="mt-2">
              <v-col cols="12" md="6">
                <div class="mb-3">
                  <strong>Member:</strong> {{ claim?.member_name }}<br />
                  <strong>ID Number:</strong> {{ claim?.member_id_number
                  }}<br />
                  <strong>Scheme:</strong> {{ claim?.scheme_name }}
                </div>
              </v-col>
              <v-col cols="12" md="6">
                <div class="mb-3">
                  <strong>Benefit Type:</strong> {{ claim?.benefit_alias
                  }}<br />
                  <strong>Member Type: </strong> {{ claim?.member_type }} <br />
                  <strong>Claim Amount:</strong>
                  {{ formatCurrency(claim?.claim_amount) }}<br />
                  <strong>Priority:</strong>
                  <v-chip
                    :color="getPriorityColor(claim?.priority)"
                    size="x-small"
                    class="ml-1"
                  >
                    {{ claim?.priority?.toUpperCase() }}
                  </v-chip>
                </div>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <strong>Date of Event:</strong>
                {{ formatDate(claim?.date_of_event) }}
              </v-col>
              <v-col cols="12" md="6">
                <strong>Date Notified:</strong>
                {{ formatDate(claim?.date_notified) }} <br />
                <strong>Date of Transaction:</strong>
                {{ formatDate(claim?.date_registered) }}
              </v-col>
            </v-row>
            <!-- Show claimant information if available -->
            <v-row
              v-if="
                claim?.claimant_name ||
                claim?.claimant_id_number ||
                claim?.claimant_contact_number
              "
            >
              <v-col cols="12">
                <v-divider class="my-2" />
                <div class="text-subtitle-2 text-primary mb-2">
                  <v-icon start color="primary">mdi-account</v-icon>
                  Claimant Information
                </div>
              </v-col>
              <v-col v-if="claim?.claimant_name" cols="12" md="6">
                <strong>Claimant Name:</strong> {{ claim.claimant_name }}
              </v-col>
              <v-col v-if="claim?.claimant_id_number" cols="12" md="6">
                <strong>Claimant ID Number:</strong>
                {{ claim.claimant_id_number }}
              </v-col>
              <v-col v-if="claim?.claimant_contact_number" cols="12" md="6">
                <strong>Contact Number:</strong>
                {{ claim.claimant_contact_number }}
              </v-col>
              <v-col v-if="claim?.relationship_to_member" cols="12" md="6">
                <strong>Relationship:</strong>
                {{ claim.relationship_to_member }}
              </v-col>
            </v-row>
            <!-- Show decline information if claim is declined -->
            <v-row v-if="claim?.status === 'declined'">
              <v-col cols="12">
                <v-divider class="my-2" />
                <div class="bg-error-lighten-4 pa-3 rounded">
                  <div class="text-subtitle-2 text-error mb-2">
                    <v-icon start color="error">mdi-information</v-icon>
                    Decline Information
                  </div>
                  <div v-if="claim?.decline_reason" class="mb-2">
                    <strong>Reason:</strong>
                    {{ getDeclineReasonText(claim.decline_reason) }}
                  </div>
                  <div v-if="claim?.decline_details" class="mb-2">
                    <strong>Details:</strong> {{ claim.decline_details }}
                  </div>
                  <div v-if="claim?.declined_date" class="text-caption">
                    <strong>Declined on:</strong>
                    {{ formatDate(claim.declined_date) }}
                  </div>
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Action Panel -->
      <v-col cols="12" md="3">
        <v-card class="mb-4">
          <v-card-title class="bg-secondary text-white">Actions</v-card-title>
          <v-card-text>
            <div class="d-flex mt-4 flex-column gap-2">
              <v-btn
                v-if="hasPermission('claims:assess')"
                rounded
                color="primary"
                variant="outlined"
                prepend-icon="mdi-pencil"
                :disabled="!canAssess"
                @click="assessmentDialog = true"
              >
                Assess Claim
              </v-btn>
              <v-btn
                v-if="hasPermission('claims:approve')"
                rounded
                color="success"
                variant="outlined"
                prepend-icon="mdi-check"
                :disabled="!canApprove"
                @click="showApproveDialog"
              >
                Approve
              </v-btn>
              <v-btn
                v-if="hasPermission('claims:reject')"
                rounded
                color="error"
                variant="outlined"
                prepend-icon="mdi-close"
                :disabled="!canDecline"
                @click="showDeclineDialog"
              >
                Decline
              </v-btn>
              <v-btn
                v-if="hasPermission('claims:assess')"
                rounded
                color="warning"
                variant="outlined"
                prepend-icon="mdi-information"
                :disabled="!canRequestInfo"
                @click="requestInfoDialog = true"
              >
                Request Info
              </v-btn>
              <v-btn
                rounded
                color="info"
                variant="outlined"
                prepend-icon="mdi-file-document"
                @click="generateReport"
              >
                Generate Report
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Assessment History -->
      <v-col v-if="assessmentHistory.length > 0" cols="12">
        <v-card class="mb-4">
          <v-card-title class="bg-info text-white"
            >Assessment History</v-card-title
          >
          <v-card-text>
            <v-timeline density="compact">
              <v-timeline-item
                v-for="(entry, index) in assessmentHistory"
                :key="index"
                :dot-color="getTimelineColor(entry.action)"
                size="small"
              >
                <div class="d-flex justify-space-between">
                  <div>
                    <div class="font-weight-medium">{{ entry.action }}</div>
                    <div class="text-body-2">{{ entry.description }}</div>
                    <div class="text-caption text-grey"
                      >By {{ entry.assessor }}</div
                    >
                  </div>
                  <div class="text-caption text-grey">
                    {{ formatDateTime(entry.timestamp) }}
                  </div>
                </div>
              </v-timeline-item>
            </v-timeline>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Supporting Documents -->
      <v-col cols="12" md="6">
        <v-card class="mb-4">
          <v-card-title class="bg-warning text-white"
            >Supporting Documents</v-card-title
          >
          <v-card-text>
            <v-list density="compact">
              <v-list-item
                v-for="(doc, index) in supportingDocuments"
                :key="index"
                class="cursor-pointer"
                @click="viewDocument(doc)"
              >
                <template #prepend>
                  <v-icon>{{ getDocumentIcon(doc.content_type) }}</v-icon>
                </template>
                <v-list-item-title>{{ doc.document_name }}</v-list-item-title>
                <v-list-item-subtitle>{{
                  doc.uploaded_at
                }}</v-list-item-subtitle>
                <template #append>
                  <v-btn
                    icon="mdi-download"
                    variant="text"
                    size="small"
                    @click.stop="downloadDocument(doc)"
                  />
                </template>
              </v-list-item>
            </v-list>
            <v-divider class="my-2" />
            <v-btn
              color="primary"
              variant="text"
              prepend-icon="mdi-plus"
              @click="uploadDialog = true"
            >
              Add Document
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Missing Documents -->
      <v-col v-if="missingDocs.length !== 0" cols="12" md="6">
        <v-card class="mb-4">
          <v-card-title class="bg-info text-white"
            >Missing Documents</v-card-title
          >
          <v-card-text>
            <v-list density="compact">
              <v-list-item
                v-for="(doc, index) in missingDocs"
                :key="index"
                class="cursor-pointer"
              >
                <template #prepend>
                  <v-icon>{{ getDocumentIcon(doc.content_type) }}</v-icon>
                </template>
                <v-list-item-title>{{ doc.name }}</v-list-item-title>
              </v-list-item>
            </v-list>
            <v-divider class="my-2" />
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Communication Log -->
      <v-col cols="12" md="6">
        <v-card class="mb-4">
          <v-card-title class="bg-teal text-white">
            <div class="d-flex justify-space-between align-center w-100">
              <span>Communication Log</span>
              <v-btn
                icon="mdi-refresh"
                variant="text"
                color="white"
                size="small"
                :loading="communicationsLoading"
                @click="loadCommunications"
              />
            </div>
          </v-card-title>
          <v-card-text>
            <!-- Loading State -->
            <v-container v-if="communicationsLoading" class="text-center py-4">
              <v-progress-circular indeterminate color="teal" size="32" />
              <div class="mt-2 text-caption">Loading communications...</div>
            </v-container>

            <!-- Error State -->
            <v-alert
              v-else-if="communicationsError"
              type="error"
              variant="tonal"
              class="mb-3"
            >
              {{ communicationsError }}
            </v-alert>

            <!-- Communications List -->
            <v-list
              v-else
              density="compact"
              max-height="300"
              style="overflow-y: auto"
            >
              <v-list-item
                v-if="communications.length === 0"
                class="text-center text-grey"
              >
                <v-list-item-title>No communications yet</v-list-item-title>
              </v-list-item>
              <v-list-item
                v-for="(comm, index) in communications"
                :key="comm.id || index"
              >
                <v-list-item-title class="text-body-2">{{
                  comm.message
                }}</v-list-item-title>
                <v-list-item-subtitle>
                  {{ comm.method }} -
                  {{ formatDateTime(comm.timestamp || comm.created_at) }}
                  <span v-if="comm.created_by" class="text-caption ml-2"
                    >by {{ comm.created_by }}</span
                  >
                </v-list-item-subtitle>
              </v-list-item>
            </v-list>
            <v-divider class="my-2" />
            <v-btn
              color="primary"
              variant="text"
              prepend-icon="mdi-message"
              :disabled="communicationsLoading"
              @click="communicationDialog = true"
            >
              Add Communication
            </v-btn>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Detailed Assessments -->
      <v-col v-if="claimAssessments.length > 0" cols="12">
        <v-card class="mb-4">
          <v-card-title class="bg-purple text-white">
            <div class="d-flex justify-space-between align-center w-100">
              <span>Detailed Assessments ({{ claimAssessments.length }})</span>
              <v-btn
                icon="mdi-refresh"
                variant="text"
                color="white"
                size="small"
                @click="loadAssessments"
              />
            </div>
          </v-card-title>
          <v-card-text>
            <v-expansion-panels variant="accordion">
              <v-expansion-panel
                v-for="assessment in claimAssessments"
                :key="assessment.id"
                :title="`Assessment by ${assessment.assessor_name} - ${formatDate(assessment.assessment_date)}`"
              >
                <template #text>
                  <div class="pa-2">
                    <v-row>
                      <v-col cols="12" md="6">
                        <div class="mb-3">
                          <strong>Outcome:</strong>
                          <v-chip
                            :color="
                              getAssessmentOutcomeColor(
                                assessment.assessment_outcome
                              )
                            "
                            size="small"
                            class="ml-2"
                          >
                            {{ assessment.assessment_outcome }}
                          </v-chip>
                        </div>
                        <div class="mb-3">
                          <strong>Recommended Amount:</strong>
                          {{ formatCurrency(assessment.recommended_amount)
                          }}<br />
                          <strong>Risk Level:</strong>
                          <v-chip
                            :color="
                              getRiskLevelColor(assessment.fraud_risk_level)
                            "
                            size="x-small"
                            class="ml-1"
                          >
                            {{ assessment.fraud_risk_level?.toUpperCase() }}
                          </v-chip>
                        </div>
                      </v-col>
                      <v-col cols="12" md="6">
                        <div v-if="assessment.medical_officer" class="mb-3">
                          <strong>Medical Officer:</strong>
                          {{ assessment.medical_officer }}<br />
                          <strong>Medical Condition:</strong>
                          {{ assessment.medical_condition }}<br />
                          <strong>Disability %:</strong>
                          {{ assessment.disability_percentage }}
                        </div>
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col cols="12">
                        <v-textarea
                          :model-value="assessment.assessment_notes"
                          label="Assessment Notes"
                          variant="outlined"
                          readonly
                          rows="3"
                        />
                        <v-textarea
                          v-if="assessment.next_actions"
                          :model-value="assessment.next_actions"
                          label="Next Actions"
                          variant="outlined"
                          readonly
                          rows="2"
                        />
                      </v-col>
                    </v-row>
                    <v-row v-if="assessment.documents_verified">
                      <v-col cols="12">
                        <div class="text-subtitle-2 mb-2"
                          >Document Verification Status:</div
                        >
                        <div class="d-flex flex-wrap gap-2">
                          <v-chip
                            v-for="(
                              verified, docType
                            ) in assessment.documents_verified"
                            :key="String(docType)"
                            :color="verified ? 'success' : 'error'"
                            size="small"
                          >
                            <v-icon
                              start
                              :icon="verified ? 'mdi-check' : 'mdi-close'"
                            ></v-icon>
                            {{ formatDocumentType(String(docType)) }}
                          </v-chip>
                        </div>
                      </v-col>
                    </v-row>
                  </div>
                </template>
              </v-expansion-panel>
            </v-expansion-panels>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Assessment Dialog -->
    <v-dialog v-model="assessmentDialog" persistent max-width="800px">
      <v-card>
        <v-card-title class="bg-primary text-white"
          >Claim Assessment</v-card-title
        >
        <v-card-text class="pt-4">
          <claim-assessment-form
            :claim="claim"
            @save="handleAssessmentSave"
            @cancel="assessmentDialog = false"
          />
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- Request Additional Info Dialog -->
    <v-dialog v-model="requestInfoDialog" persistent max-width="600px">
      <v-card>
        <v-card-title>Request Additional Information</v-card-title>
        <v-card-text>
          <v-textarea
            v-model="requestInfoMessage"
            label="Information Required"
            rows="4"
            variant="outlined"
            placeholder="Please specify what additional information is needed..."
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="grey" variant="text" @click="requestInfoDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="primary" @click="handleRequestInfo">Send Request</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Upload Document Dialog -->
    <v-dialog v-model="uploadDialog" persistent max-width="500px">
      <v-card>
        <v-card-title>Upload Document</v-card-title>
        <v-card-text>
          <v-file-input
            v-model="uploadFile"
            label="Select File"
            variant="outlined"
            accept=".pdf,.jpg,.jpeg,.png,.doc,.docx"
            show-size
          />
          <!-- <v-text-field
            v-model="uploadDescription"
            label="Document Description"
            variant="outlined"
            class="mt-2"
          /> -->
          <v-select
            v-model="selectedDocType"
            variant="outlined"
            density="compact"
            :items="missingDocs"
            item-title="name"
            label="Document Type"
            placeholder="Select Document Type"
            return-object
          >
          </v-select>
        </v-card-text>

        <v-card-actions>
          <v-spacer />
          <v-btn color="grey" variant="text" @click="uploadDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="primary" :disabled="!uploadFile" @click="handleUpload"
            >Upload</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Communication Dialog -->
    <v-dialog v-model="communicationDialog" persistent max-width="600px">
      <v-card>
        <v-card-title>Add Communication</v-card-title>
        <v-card-text>
          <v-select
            v-model="newCommunication.method"
            :items="['Email', 'Phone', 'SMS', 'Letter', 'In Person']"
            label="Communication Method"
            variant="outlined"
          />
          <v-textarea
            v-model="newCommunication.message"
            label="Message/Notes"
            rows="4"
            variant="outlined"
            class="mt-2"
          />
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            color="grey"
            variant="text"
            @click="communicationDialog = false"
            >Cancel</v-btn
          >
          <v-btn color="primary" @click="handleAddCommunication">Add</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Decline Claim Dialog -->
    <v-dialog v-model="declineDialog" persistent max-width="700px">
      <v-card>
        <v-card-title class="bg-error text-white">
          <v-icon start>mdi-close-circle</v-icon>
          Decline Claim - {{ claim?.claim_number }}
        </v-card-title>
        <v-card-text class="pt-4">
          <v-alert type="warning" variant="tonal" class="mb-4">
            <template #title>Important Notice</template>
            Declining this claim will permanently change its status and may
            trigger member notifications. Please ensure all decline reasons are
            properly documented.
          </v-alert>

          <v-row>
            <v-col cols="12" md="6">
              <v-select
                v-model="declineData.primary_reason"
                :items="declineReasons"
                item-title="text"
                item-value="value"
                label="Primary Decline Reason *"
                variant="outlined"
                :rules="[(v) => !!v || 'Primary reason is required']"
              />
            </v-col>
            <v-col cols="12" md="6">
              <v-text-field
                v-model="declineData.assessment_reference"
                label="Assessment Reference"
                variant="outlined"
                placeholder="Reference specific assessment findings"
                hint="Optional reference to assessment that supports this decline"
              />
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <v-textarea
                v-model="declineData.detailed_reason"
                label="Detailed Decline Reason *"
                variant="outlined"
                rows="4"
                :rules="[(v) => !!v || 'Detailed reason is required']"
                placeholder="Provide specific details about why this claim is being declined. This information may be shared with the member."
                hint="Be specific and professional as this may be shared with the member"
              />
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <v-textarea
                v-model="declineData.internal_notes"
                label="Internal Notes"
                variant="outlined"
                rows="3"
                placeholder="Internal notes for future reference (not shared with member)"
                hint="Optional internal notes for audit trail and future reference"
              />
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <v-checkbox
                v-model="declineData.requires_member_notification"
                label="Send decline notification to member"
                hint="Automatically create communication record and notify member of decline"
              />
            </v-col>
          </v-row>

          <v-divider class="my-3" />

          <v-row>
            <v-col cols="12">
              <div class="text-subtitle-2 mb-2">Claim Summary:</div>
              <div class="bg-grey-lighten-4 pa-3 rounded">
                <div class="text-body-2">
                  <strong>Member:</strong> {{ claim?.member_name }}<br />
                  <strong>Benefit:</strong> {{ claim?.benefit_type }}<br />
                  <strong>Amount:</strong>
                  {{ formatCurrency(claim?.claim_amount) }}<br />
                  <strong>Current Status:</strong>
                  {{ formatStatus(claim?.status) }}
                </div>
              </div>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="px-6 pb-4">
          <v-spacer />
          <v-btn color="grey" variant="text" @click="handleDeclineCancel">
            Cancel
          </v-btn>
          <v-btn
            color="error"
            variant="elevated"
            :disabled="
              !declineData.primary_reason || !declineData.detailed_reason
            "
            @click="handleDecline"
          >
            <v-icon start>mdi-close-circle</v-icon>
            Decline Claim
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Document Viewer Dialog -->
    <v-dialog v-model="documentViewerDialog" max-width="900px">
      <v-card>
        <v-card-title class="bg-primary text-white">
          <div class="d-flex justify-space-between align-center w-100">
            <span>{{
              currentDocument?.filename || currentDocument?.name
            }}</span>
            <v-btn
              icon="mdi-close"
              variant="text"
              color="white"
              size="small"
              @click="documentViewerDialog = false"
            />
          </div>
        </v-card-title>
        <v-card-text class="pa-0">
          <div v-if="documentUrl" style="height: 600px">
            <!-- PDF Viewer -->
            <iframe
              v-if="isPdfDocument(currentDocument)"
              :src="documentUrl"
              style="width: 100%; height: 100%; border: none"
              title="Document Viewer"
            />
            <!-- Image Viewer -->
            <div
              v-else-if="isImageDocument(currentDocument)"
              class="d-flex justify-center align-center"
              style="height: 100%; background-color: #f5f5f5"
            >
              <v-img
                :src="documentUrl"
                :alt="currentDocument?.filename || currentDocument?.name"
                contain
                max-height="580px"
                max-width="100%"
              />
            </div>
            <!-- Unsupported file type -->
            <div
              v-else
              class="d-flex flex-column justify-center align-center"
              style="height: 100%"
            >
              <v-icon size="64" color="grey">mdi-file-document</v-icon>
              <div class="text-h6 mt-4 text-grey">Preview not available</div>
              <div class="text-body-2 text-grey">{{
                getFileTypeDescription(currentDocument)
              }}</div>
              <v-btn
                color="primary"
                class="mt-4"
                prepend-icon="mdi-download"
                @click="downloadDocument(currentDocument)"
              >
                Download to View
              </v-btn>
            </div>
          </div>
          <div
            v-else
            class="d-flex flex-column justify-center align-center"
            style="height: 400px"
          >
            <v-icon size="64" color="error">mdi-alert-circle</v-icon>
            <div class="text-h6 mt-4 text-error">Document not found</div>
            <div class="text-body-2 text-grey"
              >The document URL is not available</div
            >
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            color="primary"
            prepend-icon="mdi-download"
            @click="downloadDocument(currentDocument)"
          >
            Download
          </v-btn>
          <v-btn
            color="grey"
            variant="text"
            @click="documentViewerDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Confirmation Dialogs -->
    <v-dialog v-model="confirmDialog" persistent max-width="400px">
      <v-card>
        <v-card-title>{{ confirmTitle }}</v-card-title>
        <v-card-text>{{ confirmMessage }}</v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="grey" variant="text" @click="confirmDialog = false"
            >Cancel</v-btn
          >
          <v-btn :color="confirmColor" @click="confirmAction">{{
            confirmButtonText
          }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import ClaimAssessmentForm from './ClaimAssessmentForm.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { usePermissionCheck } from '@/renderer/composables/usePermissionCheck'

interface Props {
  claim: any
}

interface Emits {
  (e: 'update', claim: any): void
  (e: 'close'): void
  (e: 'assessment-created', assessment: any): void
}

interface DocumentMetadata {
  code: string
  name: string
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const { hasPermission } = usePermissionCheck()

// Dialog states
const assessmentDialog = ref(false)
const requestInfoDialog = ref(false)
const uploadDialog = ref(false)
const communicationDialog = ref(false)
const confirmDialog = ref(false)
const documentViewerDialog = ref(false)
const declineDialog = ref(false)

// Form data
const requestInfoMessage = ref('')
const uploadFile = ref<File | null>(null)
// const uploadDescription = ref('')
const newCommunication = ref({
  method: '',
  message: ''
})

const declineData = ref({
  primary_reason: '',
  detailed_reason: '',
  assessment_reference: '',
  requires_member_notification: true,
  internal_notes: ''
})

// Confirmation dialog data
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmColor = ref('primary')
const confirmButtonText = ref('Confirm')
const confirmCallback = ref<(() => void) | null>(null)

// Document viewer data
const currentDocument = ref<any>(null)
const documentUrl = ref('')
const documentBlobUrls = ref<string[]>([])

// Assessments data
const claimAssessments = ref<any[]>([])

// Communication data and loading states
const communications = ref<any[]>([])
const communicationsLoading = ref(false)
const communicationsError = ref('')

// Assessment history data
const assessmentHistory = ref<any[]>([])

const supportingDocuments = ref<any[]>([])
const selectedDocType = ref<DocumentMetadata | null>(null)

// Mock fallback data
const mockDocuments = [
  {
    id: 1,
    filename: 'Death Certificate.pdf',
    type: 'pdf',
    uploadDate: '2024-01-20',
    size: '2.3 MB'
  },
  {
    id: 2,
    name: 'ID Document.jpg',
    type: 'image',
    uploadDate: '2024-01-20',
    size: '1.1 MB'
  },
  {
    id: 3,
    name: 'Banking Details.pdf',
    type: 'pdf',
    uploadDate: '2024-01-21',
    size: '0.8 MB'
  }
]

// Computed properties
const canAssess = computed(() => {
  return ['pending', 'under_assessment', 'additional_info_required'].includes(
    props.claim?.status
  )
})

const declineReasons = computed(() => [
  { value: 'insufficient_documentation', text: 'Insufficient Documentation' },
  { value: 'policy_exclusion', text: 'Policy Exclusion' },
  { value: 'pre_existing_condition', text: 'Pre-existing Condition' },
  { value: 'fraud_suspected', text: 'Fraud Suspected' },
  { value: 'coverage_expired', text: 'Coverage Expired' },
  { value: 'waiting_period_not_met', text: 'Waiting Period Not Met' },
  { value: 'benefit_limit_exceeded', text: 'Benefit Limit Exceeded' },
  { value: 'incorrect_provider', text: 'Incorrect Provider' },
  { value: 'duplicate_claim', text: 'Duplicate Claim' },
  { value: 'late_submission', text: 'Late Submission' },
  { value: 'medical_necessity', text: 'Medical Necessity Not Established' },
  { value: 'non_covered_service', text: 'Non-covered Service' },
  { value: 'investigation_findings', text: 'Investigation Findings' },
  { value: 'other', text: 'Other (Specify in notes)' }
])

const documentTypesMapping = {
  GLA: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_deceased',
      name: 'Certified ID - Deceased',
      required: true
    },
    {
      code: 'certified_id_claimant',
      name: 'Certified ID - Claimant/Beneficiaries',
      required: true
    },
    {
      code: 'death_certificate',
      name: 'Death Certificate (BI-5)',
      required: true
    },
    {
      code: 'dha_notification',
      name: 'DHA-1663 Notification of Death',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: true
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'post_mortem',
      name: 'Post-mortem / Final BI-1680/1683',
      required: false
    }
  ],
  SGLA: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_deceased',
      name: 'Certified ID - Deceased',
      required: true
    },
    {
      code: 'certified_id_claimant',
      name: 'Certified ID - Claimant/Beneficiaries',
      required: true
    },
    {
      code: 'death_certificate',
      name: 'Death Certificate (BI-5)',
      required: true
    },
    {
      code: 'dha_notification',
      name: 'DHA-1663 Notification of Death',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: true
    },
    {
      code: 'relationship_proof',
      name: 'Proof of Relationship (Spouse/Child)',
      required: true
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'post_mortem',
      name: 'Post-mortem / Final BI-1680/1683',
      required: false
    }
  ],
  GFF: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_deceased',
      name: 'Certified ID - Deceased',
      required: true
    },
    {
      code: 'certified_id_claimant',
      name: 'Certified ID - Claimant/Beneficiaries',
      required: true
    },
    {
      code: 'death_certificate',
      name: 'Death Certificate (BI-5)',
      required: true
    },
    {
      code: 'dha_notification',
      name: 'DHA-1663 Notification of Death',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: true
    },
    {
      code: 'relationship_proof',
      name: 'Proof of Relationship (Spouse/Child)',
      required: true
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'post_mortem',
      name: 'Post-mortem / Final BI-1680/1683',
      required: false
    }
  ],
  PTD: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_member',
      name: 'Certified ID - Member',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: false
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'medical_reports',
      name: 'Medical Reports - treating doctor report',
      required: true
    },
    {
      code: 'attending_doctor_statement',
      name: "Attending Doctor's Statement (Disability/CI Report)",
      required: true
    },
    {
      code: 'specialist_report',
      name: 'Specialist Medical Report (e.g., Oncologist, Cardiologist, Neurologist)',
      required: false
    },
    {
      code: 'employer_duties_statement',
      name: 'Employer Statement of Duties / Job Description',
      required: true
    },
    {
      code: 'functional_capacity_assessment',
      name: 'Functional Capacity Assessment (FCE)',
      required: false
    },
    {
      code: 'occupational_therapist_report',
      name: 'Occupational Therapist Report',
      required: false
    }
  ],
  CI: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_member',
      name: 'Certified ID - Member',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: false
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'medical_reports',
      name: 'Medical Reports - treating doctor report',
      required: true
    },
    {
      code: 'attending_doctor_statement',
      name: "Attending Doctor's Statement (Disability/CI Report)",
      required: true
    },
    {
      code: 'specialist_report',
      name: 'Specialist Medical Report (e.g., Oncologist, Cardiologist, Neurologist)',
      required: true
    },
    {
      code: 'employer_duties_statement',
      name: 'Employer Statement of Duties / Job Description',
      required: false
    },
    {
      code: 'functional_capacity_assessment',
      name: 'Functional Capacity Assessment (FCE)',
      required: false
    },
    {
      code: 'occupational_therapist_report',
      name: 'Occupational Therapist Report',
      required: false
    }
  ],
  PHI: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_member',
      name: 'Certified ID - Member',
      required: true
    },
    {
      code: 'beneficiary_form',
      name: 'Beneficiary Nomination Form / Employer Beneficiary Statement',
      required: false
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'accident_report',
      name: 'Accident Report / Police Report (if accidental cause)',
      required: false
    },
    {
      code: 'medical_reports',
      name: 'Medical Reports - treating doctor report',
      required: true
    },
    {
      code: 'specialist_report',
      name: 'Specialist Medical Report (e.g., Oncologist, Cardiologist, Neurologist)',
      required: true
    },
    {
      code: 'employer_duties_statement',
      name: 'Employer Statement of Duties / Job Description',
      required: true
    },
    {
      code: 'functional_capacity_assessment',
      name: 'Functional Capacity Assessment (FCE)',
      required: false
    },
    {
      code: 'occupational_therapist_report',
      name: 'Occupational Therapist Report',
      required: false
    },
    {
      code: 'psychiatric_report',
      name: 'Psychiatric Report (if mental illness claim)',
      required: false
    },
    {
      code: 'income_loss_proof',
      name: 'Proof of Income Loss / Sick Leave Records',
      required: true
    }
  ],
  TTD: [
    {
      code: 'claim_form',
      name: 'Claim Form (official insurer form)',
      required: true
    },
    {
      code: 'certified_id_member',
      name: 'Certified ID - Member',
      required: true
    },
    {
      code: 'medical_reports',
      name: 'Medical Reports - treating doctor report',
      required: true
    },
    {
      code: 'banking_details',
      name: 'Banking Details - beneficiary or member',
      required: true
    },
    {
      code: 'employment_proof',
      name: 'Proof of Employment / HR Letter',
      required: true
    },
    {
      code: 'salary_confirmation',
      name: 'Salary Confirmation / CTC / Pensionable Salary',
      required: true
    }
  ]
}

const missingDocs = computed(() => {
  const allDocTypes = documentTypesMapping[props.claim.benefit_code] || []
  const uploadedDocs = new Set(
    supportingDocuments.value.map((d) => d.document_type)
  )

  return allDocTypes.filter((d) => !uploadedDocs.has(d.code))
})

const canApprove = computed(() => {
  return [
    'under_assessment',
    'additional_info_required',
    'referred_to_committee'
  ].includes(props.claim?.status)
})

const canDecline = computed(() => {
  return [
    'under_assessment',
    'additional_info_required',
    'referred_to_committee'
  ].includes(props.claim?.status)
})

const canRequestInfo = computed(() => {
  return [
    'pending',
    'under_assessment',
    'additional_info_required',
    'referred_to_committee'
  ].includes(props.claim?.status)
})

// Methods
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    pending: 'info',
    under_assessment: 'warning',
    additional_info_required: 'orange',
    referred_to_committee: 'deep-purple',
    approved: 'success',
    declined: 'error',
    paid: 'teal',
    payment_failed: 'error',
    cancelled: 'grey'
  }
  return colors[status] || 'default'
}

const getPriorityColor = (priority: string) => {
  const colors: Record<string, string> = {
    low: 'success',
    medium: 'warning',
    high: 'error',
    critical: 'purple'
  }
  return colors[priority] || 'default'
}

const getTimelineColor = (action: string) => {
  if (action.includes('Register')) return 'primary'
  if (action.includes('Review') || action.includes('Assess')) return 'info'
  if (action.includes('Approve')) return 'success'
  if (action.includes('Decline')) return 'error'
  return 'grey'
}

const getDocumentIcon = (type: string) => {
  const icons: Record<string, string> = {
    pdf: 'mdi-file-pdf-box',
    image: 'mdi-image',
    doc: 'mdi-file-word-box',
    docx: 'mdi-file-word-box',
    'image/png': 'mdi-image'
  }
  return icons[type] || 'mdi-file-document'
}

const formatStatus = (status: string) => {
  return (
    status
      ?.split('_')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ') || ''
  )
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('en-ZA', {
    style: 'currency',
    currency: 'ZAR'
  }).format(amount || 0)
}

const formatDate = (dateString: string) => {
  return dateString ? new Date(dateString).toLocaleDateString() : ''
}

const formatDateTime = (dateString: string) => {
  return dateString ? new Date(dateString).toLocaleString() : ''
}

const showApproveDialog = () => {
  confirmTitle.value = 'Approve Claim'
  confirmMessage.value = `Are you sure you want to approve this claim for ${formatCurrency(props.claim?.claim_amount)}?`
  confirmColor.value = 'success'
  confirmButtonText.value = 'Approve'
  confirmCallback.value = handleApprove
  confirmDialog.value = true
}

const showDeclineDialog = () => {
  // Reset decline form
  declineData.value = {
    primary_reason: '',
    detailed_reason: '',
    assessment_reference: '',
    requires_member_notification: true,
    internal_notes: ''
  }
  declineDialog.value = true
}

const handleApprove = () => {
  const updatedClaim = {
    ...props.claim,
    status: 'approved',
    approved_date: new Date().toISOString().split('T')[0]
  }
  emit('update', updatedClaim)
  confirmDialog.value = false
}

const handleDeclineCancel = () => {
  declineDialog.value = false
  // Reset form data
  declineData.value = {
    primary_reason: '',
    detailed_reason: '',
    assessment_reference: '',
    requires_member_notification: true,
    internal_notes: ''
  }
}

const handleDecline = async () => {
  if (!declineData.value.primary_reason || !declineData.value.detailed_reason) {
    // TODO: Show validation error
    return
  }

  try {
    const declineTimestamp = new Date().toISOString()

    // Create decline record
    const declineRecord = {
      claim_id: props.claim?.id,
      primary_reason: declineData.value.primary_reason,
      detailed_reason: declineData.value.detailed_reason,
      assessment_reference: declineData.value.assessment_reference,
      requires_member_notification:
        declineData.value.requires_member_notification,
      internal_notes: declineData.value.internal_notes,
      declined_by: 'current_user', // This should come from auth context
      declined_at: declineTimestamp
    }

    // Save decline reasons to backend
    await GroupPricingService.createClaimDeclineRecord(
      props.claim.id,
      declineRecord
    )

    // Update claim status
    const updatedClaim = {
      ...props.claim,
      status: 'declined',
      declined_date: declineTimestamp.split('T')[0],
      decline_reason: declineData.value.primary_reason,
      decline_details: declineData.value.detailed_reason,
      last_updated: declineTimestamp
    }

    // Save updated claim to backend
    await GroupPricingService.updateClaim(props.claim?.id, updatedClaim)

    // Add to assessment history
    assessmentHistory.value.unshift({
      action: 'Claim Declined',
      description: `Declined: ${declineReasons.value.find((r) => r.value === declineData.value.primary_reason)?.text || declineData.value.primary_reason}. ${declineData.value.detailed_reason}`,
      assessor: 'current_user',
      timestamp: declineTimestamp
    })

    // Create communication record if member notification is required
    if (declineData.value.requires_member_notification) {
      const communicationData = {
        claim_id: props.claim.id,
        method: 'Email',
        message: `Claim declined: ${declineReasons.value.find((r) => r.value === declineData.value.primary_reason)?.text || declineData.value.primary_reason}. ${declineData.value.detailed_reason}`,
        timestamp: declineTimestamp,
        created_by: 'current_user'
      }

      try {
        const response = await GroupPricingService.createClaimCommunication(
          props.claim.id,
          communicationData
        )
        communications.value.unshift(response.data)
      } catch (error) {
        console.error('Failed to create decline communication:', error)
      }
    }

    emit('update', updatedClaim)
    declineDialog.value = false
  } catch (error) {
    console.error('Failed to decline claim:', error)
    // TODO: Show user-friendly error message
  }
}

const handleRequestInfo = async () => {
  if (!props.claim?.id || !requestInfoMessage.value) return

  try {
    const updatedClaim = {
      ...props.claim,
      status: 'additional_info_required'
    }

    // Save communication to backend
    const communicationData = {
      claim_id: props.claim.id,
      method: 'Email',
      message: `Additional information requested: ${requestInfoMessage.value}`,
      timestamp: new Date().toISOString(),
      created_by: 'current_user'
    }

    const response = await GroupPricingService.createClaimCommunication(
      props.claim.id,
      communicationData
    )
    const savedCommunication = response.data

    // Add to local communication log
    communications.value.unshift(savedCommunication)

    // Update claim status
    await GroupPricingService.updateClaim(props.claim.id, updatedClaim)

    emit('update', updatedClaim)
    requestInfoDialog.value = false
    requestInfoMessage.value = ''
  } catch (error: any) {
    console.error('Failed to request additional info:', error)
    communicationsError.value =
      error.response?.data?.message || 'Failed to save information request'
    // TODO: Show user-friendly error message
  }
}

const loadCommunications = async () => {
  if (!props.claim?.id) return

  communicationsLoading.value = true
  communicationsError.value = ''

  try {
    const response = await GroupPricingService.getClaimCommunications(
      props.claim.id
    )
    communications.value = response.data || []
  } catch (error: any) {
    console.error('Failed to load communications:', error)
    communicationsError.value =
      error.response?.data?.message || 'Failed to load communications'
    // Keep existing communications on error
  } finally {
    communicationsLoading.value = false
  }
}

const handleAssessmentSave = async (assessmentData: any) => {
  try {
    // Create assessment payload for backend
    const assessmentPayload = {
      claim_id: props.claim?.id,
      ...assessmentData,
      assessment_timestamp:
        assessmentData.assessment_timestamp || new Date().toISOString()
    }

    // Save assessment to backend
    const assessmentResponse = await GroupPricingService.createClaimAssessment(
      assessmentPayload.claim_id,
      assessmentPayload
    )
    const savedAssessment = assessmentResponse.data

    // Determine claim status based on assessment outcome
    let newStatus = 'under_assessment'
    const outcome = assessmentData.assessment_outcome
    if (outcome === 'refer_to_committee') {
      newStatus = 'referred_to_committee'
    } else if (outcome === 'committee_approved') {
      newStatus = 'approved'
    } else if (outcome === 'committee_declined') {
      newStatus = 'declined'
    } else if (outcome === 'requires_additional_info') {
      newStatus = 'additional_info_required'
    }

    // Update claim with latest assessment info
    const updatedClaim = {
      ...props.claim,
      status: newStatus,
      last_updated: new Date().toISOString()
    }

    // Save updated claim to backend
    await GroupPricingService.updateClaim(props.claim?.id, updatedClaim)

    // Update local state only after successful backend saves
    claimAssessments.value.unshift(savedAssessment)

    // Add to assessment history timeline
    assessmentHistory.value.unshift({
      action: 'Assessment Update',
      description: assessmentData.assessment_notes || 'Assessment completed',
      assessor: assessmentData.assessor_name,
      timestamp: assessmentData.assessment_timestamp || new Date().toISOString()
    })

    // Emit both the updated claim and the new assessment
    emit('update', updatedClaim)
    emit('assessment-created', savedAssessment)

    assessmentDialog.value = false
  } catch (error) {
    console.error('Failed to save assessment:', error)
    // TODO: Show user-friendly error message
    // You might want to show a toast/snackbar here
  }
}

const handleUpload = async () => {
  if (!uploadFile.value || !props.claim?.id || !selectedDocType.value) return

  try {
    const file = uploadFile.value
    const formData = new FormData()
    formData.append('file', file)
    // formData.append('description', uploadDescription.value)
    formData.append('claim_id', props.claim.id.toString())
    formData.append('document_type', selectedDocType.value.code)
    formData.append('document_name', selectedDocType.value.name)

    // Upload file to backend
    const response = await GroupPricingService.uploadClaimDocument(
      props.claim.id,
      formData
    )
    const uploadedDocument = response.data

    // Add the uploaded document to the local list
    supportingDocuments.value.unshift(uploadedDocument[0])

    uploadDialog.value = false
    uploadFile.value = null
    // uploadDescription.value = ''
    selectedDocType.value = null
  } catch (error: any) {
    console.error('Failed to upload document:', error)
    // TODO: Show user-friendly error message (toast/snackbar)
    // You might want to show an error dialog or toast here
  }
}

const handleAddCommunication = async () => {
  if (!newCommunication.value.method || !newCommunication.value.message) return
  if (!props.claim?.id) return

  try {
    const communicationData = {
      claim_id: props.claim.id,
      method: newCommunication.value.method,
      message: newCommunication.value.message,
      timestamp: new Date().toISOString(),
      created_by: 'current_user' // This should come from auth context
    }

    // Save to backend
    const response = await GroupPricingService.createClaimCommunication(
      props.claim.id,
      communicationData
    )
    const savedCommunication = response.data

    // Add to local state
    communications.value.unshift(savedCommunication)

    communicationDialog.value = false
    newCommunication.value = { method: '', message: '' }
  } catch (error: any) {
    console.error('Failed to save communication:', error)
    communicationsError.value =
      error.response?.data?.message || 'Failed to save communication'
    // TODO: Show user-friendly error message (toast/snackbar)
  }
}

const viewDocument = async (doc: any) => {
  currentDocument.value = doc
  const baseUrl = window.mainApi.sendSync('msgGetBaseUrl')
  const token = window.mainApi.sendSync('msgGetAccessToken')
  // strip trailing slash if present
  const normalizedBaseUrl = baseUrl.endsWith('/')
    ? baseUrl.slice(0, -1)
    : baseUrl

  // Check if document has a URL or path
  let url = ''
  if (doc.viewer_url) {
    url = normalizedBaseUrl + doc.viewer_url
  } else if (doc.path) {
    url = doc.path
  } else if (doc.file_path) {
    url = doc.file_path
  } else {
    // Generate a placeholder URL or handle missing URL
    console.warn('Document has no URL or path:', doc)
    documentUrl.value = ''
    documentViewerDialog.value = true
    return
  }

  // Fetch document with authorization headers and create blob URL
  try {
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(
        `Failed to fetch document: ${response.status} ${response.statusText}`
      )
    }

    const blob = await response.blob()
    const blobUrl = URL.createObjectURL(blob)
    documentUrl.value = blobUrl

    // Keep track of blob URLs for cleanup
    documentBlobUrls.value.push(blobUrl)
  } catch (error) {
    console.error('Error fetching document:', error)
    documentUrl.value = ''
  }

  documentViewerDialog.value = true
}

const downloadDocument = async (doc: any) => {
  if (!doc) return

  const baseUrl = window.mainApi.sendSync('msgGetBaseUrl')
  const token = window.mainApi.sendSync('msgGetAccessToken')
  const normalizedBaseUrl = baseUrl.endsWith('/')
    ? baseUrl.slice(0, -1)
    : baseUrl

  // Determine the document URL
  let url = ''
  if (doc.viewer_url) {
    url = normalizedBaseUrl + doc.viewer_url
  } else if (doc.path) {
    url = doc.path
  } else if (doc.file_path) {
    url = doc.file_path
  } else {
    console.warn('Document has no URL or path for download:', doc)
    return
  }

  try {
    // Fetch document with authorization headers
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })

    if (!response.ok) {
      throw new Error(
        `Failed to download document: ${response.status} ${response.statusText}`
      )
    }

    // Get the blob
    const blob = await response.blob()

    // Determine filename
    const filename = doc.filename || doc.name || `document_${Date.now()}`

    // Create download link
    const downloadUrl = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename

    // Trigger download
    document.body.appendChild(link)
    link.click()

    // Cleanup
    document.body.removeChild(link)
    URL.revokeObjectURL(downloadUrl)
  } catch (error) {
    console.error('Error downloading document:', error)
    // TODO: Show user-friendly error message (toast/snackbar)
  }
}

const generateReport = async () => {
  if (!props.claim) {
    console.error('No claim data available for report generation')
    return
  }

  try {
    // Generate report content
    const reportData = {
      claim: props.claim,
      assessments: claimAssessments.value,
      assessmentHistory: assessmentHistory.value,
      supportingDocuments: supportingDocuments.value,
      communications: communications.value,
      generatedAt: new Date().toISOString(),
      generatedBy: 'Claims Management System'
    }

    // Create HTML content for the report
    const htmlContent = generateReportHTML(reportData)

    // Create and download the report
    const blob = new Blob([htmlContent], { type: 'text/html' })
    const url = URL.createObjectURL(blob)

    const link = document.createElement('a')
    link.href = url
    link.download = `Claim_Report_${props.claim.claim_number}_${new Date().toISOString().split('T')[0]}.html`

    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    // Cleanup
    URL.revokeObjectURL(url)
  } catch (error) {
    console.error('Error generating claim report:', error)
  }
}

const generateReportHTML = (data: any) => {
  const {
    claim,
    assessments,
    assessmentHistory,
    supportingDocuments,
    communications,
    generatedAt
  } = data

  return `
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <title>Claim Report - ${claim.claim_number}</title>
      <style>
        body {
          font-family: Arial, sans-serif;
          line-height: 1.6;
          color: #333;
          max-width: 1000px;
          margin: 0 auto;
          padding: 20px;
        }
        .header {
          text-align: center;
          border-bottom: 2px solid #2196F3;
          padding-bottom: 20px;
          margin-bottom: 30px;
        }
        .section {
          margin-bottom: 30px;
        }
        .section h2 {
          color: #2196F3;
          border-left: 4px solid #2196F3;
          padding-left: 10px;
          margin-bottom: 15px;
        }
        .info-grid {
          display: grid;
          grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
          gap: 20px;
          margin-bottom: 20px;
        }
        .info-item {
          background: #f5f5f5;
          padding: 15px;
          border-radius: 5px;
        }
        .info-item strong {
          color: #1976D2;
        }
        .status-badge {
          display: inline-block;
          padding: 4px 8px;
          border-radius: 12px;
          font-size: 12px;
          font-weight: bold;
          color: white;
        }
        .status-pending { background-color: #2196F3; }
        .status-approved { background-color: #4CAF50; }
        .status-declined { background-color: #f44336; }
        .status-under-assessment { background-color: #FF9800; }
        .priority-high { background-color: #f44336; }
        .priority-medium { background-color: #FF9800; }
        .priority-low { background-color: #4CAF50; }
        table {
          width: 100%;
          border-collapse: collapse;
          margin-top: 10px;
        }
        th, td {
          padding: 12px;
          text-align: left;
          border-bottom: 1px solid #ddd;
        }
        th {
          background-color: #2196F3;
          color: white;
        }
        .assessment-item {
          background: #f9f9f9;
          border-left: 4px solid #2196F3;
          padding: 15px;
          margin-bottom: 15px;
          border-radius: 0 5px 5px 0;
        }
        .footer {
          margin-top: 40px;
          text-align: center;
          color: #666;
          border-top: 1px solid #ddd;
          padding-top: 20px;
        }
        @media print {
          body { print-color-adjust: exact; }
        }
      </style>
    </head>
    <body>
      <div class="header">
        <h1>Claims Management Report</h1>
        <p><strong>Claim Number:</strong> ${claim.claim_number}</p>
        <p><strong>Generated:</strong> ${new Date(generatedAt).toLocaleDateString()} at ${new Date(generatedAt).toLocaleTimeString()}</p>
      </div>

      <div class="section">
        <h2>Claim Summary</h2>
        <div class="info-grid">
          <div class="info-item">
            <strong>Member Name:</strong><br>${claim.member_name || 'N/A'}<br>
            <strong>ID Number:</strong><br>${claim.member_id_number || 'N/A'}<br>
            <strong>Scheme:</strong><br>${claim.scheme_name || 'N/A'}
          </div>
          <div class="info-item">
            <strong>Benefit Type:</strong><br>${claim.benefit_type || 'N/A'}<br>
            <strong>Claim Amount:</strong><br>${formatCurrency(claim.claim_amount || 0)}<br>
            <strong>Priority:</strong><br>
            <span class="status-badge priority-${claim.priority || 'low'}">${(claim.priority || 'low').toUpperCase()}</span>
          </div>
          <div class="info-item">
            <strong>Date of Event:</strong><br>${formatDate(claim.date_of_event) || 'N/A'}<br>
            <strong>Date Notified:</strong><br>${formatDate(claim.date_notified) || 'N/A'}<br>
            <strong>Status:</strong><br>
            <span class="status-badge status-${claim.status?.replace('_', '-') || 'pending'}">${formatStatus(claim.status || 'pending')}</span>
            ${
              claim.status === 'declined' && claim.decline_reason
                ? `
              <div style="margin-top: 10px; padding: 10px; background-color: #ffebee; border-left: 4px solid #f44336; border-radius: 4px;">
                <div style="font-weight: bold; color: #c62828; margin-bottom: 5px;">Decline Information:</div>
                <div><strong>Reason:</strong> ${getDeclineReasonText(claim.decline_reason)}</div>
                ${claim.decline_details ? `<div><strong>Details:</strong> ${claim.decline_details}</div>` : ''}
                ${claim.declined_date ? `<div class="text-caption"><strong>Declined on:</strong> ${formatDate(claim.declined_date)}</div>` : ''}
              </div>
              `
                : ''
            }
          </div>
        </div>
      </div>

      ${
        assessments.length > 0
          ? `
      <div class="section">
        <h2>Assessment Details</h2>
        ${assessments
          .map(
            (assessment) => `
          <div class="assessment-item">
            <h3>Assessment by ${assessment.assessor_name || 'Unknown'} - ${formatDate(assessment.assessment_date)}</h3>
            <div class="info-grid">
              <div>
                <strong>Outcome:</strong> 
                <span class="status-badge status-${assessment.assessment_outcome?.toLowerCase().replace(' ', '-') || 'pending'}">
                  ${assessment.assessment_outcome || 'Pending'}
                </span><br>
                <strong>Recommended Amount:</strong> ${formatCurrency(assessment.recommended_amount || 0)}<br>
                <strong>Risk Level:</strong> ${assessment.fraud_risk_level || 'N/A'}
              </div>
              ${
                assessment.medical_officer
                  ? `
              <div>
                <strong>Medical Officer:</strong> ${assessment.medical_officer}<br>
                <strong>Medical Condition:</strong> ${assessment.medical_condition || 'N/A'}<br>
                <strong>Disability %:</strong> ${assessment.disability_percentage || 'N/A'}%
              </div>
              `
                  : ''
              }
            </div>
            ${
              assessment.assessment_notes
                ? `
            <div style="margin-top: 15px;">
              <strong>Assessment Notes:</strong><br>
              <div style="background: white; padding: 10px; border-radius: 3px; margin-top: 5px;">
                ${assessment.assessment_notes}
              </div>
            </div>
            `
                : ''
            }
            ${
              assessment.next_actions
                ? `
            <div style="margin-top: 10px;">
              <strong>Next Actions:</strong><br>
              <div style="background: white; padding: 10px; border-radius: 3px; margin-top: 5px;">
                ${assessment.next_actions}
              </div>
            </div>
            `
                : ''
            }
          </div>
        `
          )
          .join('')}
      </div>
      `
          : ''
      }

      ${
        supportingDocuments.length > 0
          ? `
      <div class="section">
        <h2>Supporting Documents</h2>
        <table>
          <thead>
            <tr>
              <th>Document Name</th>
              <th>Type</th>
              <th>Upload Date</th>
              <th>Size</th>
            </tr>
          </thead>
          <tbody>
            ${supportingDocuments
              .map(
                (doc) => `
              <tr>
                <td>${doc.filename || doc.name || 'Unknown'}</td>
                <td>${doc.content_type || doc.type || 'Unknown'}</td>
                <td>${formatDate(doc.uploaded_at || doc.uploadDate) || 'N/A'}</td>
                <td>${doc.size || 'N/A'}</td>
              </tr>
            `
              )
              .join('')}
          </tbody>
        </table>
      </div>
      `
          : ''
      }

      ${
        communications.length > 0
          ? `
      <div class="section">
        <h2>Communication Log</h2>
        ${communications
          .map(
            (comm) => `
          <div class="assessment-item">
            <strong>${comm.method || 'Unknown Method'}:</strong> ${formatDateTime(comm.timestamp || comm.created_at)}<br>
            ${comm.created_by ? `<div style="margin-top: 2px; font-size: 0.9em; color: #666;">By: ${comm.created_by}</div>` : ''}
            <div style="margin-top: 5px; background: white; padding: 10px; border-radius: 3px;">
              ${comm.message || 'No message content'}
            </div>
          </div>
        `
          )
          .join('')}
      </div>
      `
          : ''
      }

      ${
        assessmentHistory.length > 0
          ? `
      <div class="section">
        <h2>Timeline History</h2>
        ${assessmentHistory
          .map(
            (entry) => `
          <div class="assessment-item">
            <strong>${entry.action}:</strong> ${formatDateTime(entry.timestamp)}<br>
            <div style="margin-top: 5px;">By: ${entry.assessor || 'System'}</div>
            ${
              entry.description
                ? `
            <div style="margin-top: 5px; background: white; padding: 10px; border-radius: 3px;">
              ${entry.description}
            </div>
            `
                : ''
            }
          </div>
        `
          )
          .join('')}
      </div>
      `
          : ''
      }

      <div class="footer">
        <p><em>This report was generated automatically by the Claims Management System.</em></p>
        <p><small>Report generated on ${new Date(generatedAt).toLocaleDateString()} at ${new Date(generatedAt).toLocaleTimeString()}</small></p>
      </div>
    </body>
    </html>
  `
}

const confirmAction = () => {
  if (confirmCallback.value) {
    confirmCallback.value()
    confirmCallback.value = null
  }
}

const loadAssessments = async () => {
  if (!props.claim?.id) return

  try {
    const response = await GroupPricingService.getClaimAssessments(
      props.claim.id
    )
    claimAssessments.value = response.data || []

    if (claimAssessments.value.length > 0) {
      assessmentHistory.value = claimAssessments.value.map(
        (assessment: any) => ({
          action: 'Assessment Update',
          description: assessment.assessment_notes || 'Assessment completed',
          assessor: assessment.assessor_name,
          timestamp: assessment.assessment_date
        })
      )
    }
  } catch (error) {
    console.error('Failed to load assessments:', error)
    // Fallback to mock data for development/demo purposes
    if (claimAssessments.value.length === 0) {
      claimAssessments.value = [
        {
          id: 1,
          claim_id: props.claim?.id,
          assessor_name: 'Sarah Johnson',
          assessment_date: '2024-01-21',
          assessment_outcome: 'Requires Medical Review',
          recommended_amount: props.claim?.claim_amount,
          fraud_risk_level: 'low',
          requires_investigation: false,
          assessment_notes:
            'Initial assessment completed. All documents verified. Medical review required due to nature of claim.',
          next_actions: 'Forward to medical team for specialist review.',
          documents_verified: {
            death_certificate: true,
            id_documents: true,
            banking_details: true,
            claim_form: true
          },
          created_at: '2024-01-21T14:30:00Z'
        }
      ]
    }
  }
}

const getAssessmentOutcomeColor = (outcome: string) => {
  const colors: Record<string, string> = {
    'Recommended for Approval': 'success',
    'Recommended for Partial Approval': 'warning',
    'Recommended for Decline': 'error',
    'Requires Additional Information': 'orange',
    'Requires Medical Review': 'info',
    'Requires Investigation': 'purple'
  }
  return colors[outcome] || 'default'
}

const getRiskLevelColor = (riskLevel: string) => {
  const colors: Record<string, string> = {
    low: 'success',
    medium: 'warning',
    high: 'error',
    critical: 'purple'
  }
  return colors[riskLevel] || 'default'
}

const formatDocumentType = (docType: string) => {
  return docType
    .split('_')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

const getDeclineReasonText = (reasonValue: string) => {
  const reason = declineReasons.value.find((r) => r.value === reasonValue)
  return reason ? reason.text : reasonValue
}

const isPdfDocument = (doc: any) => {
  if (!doc) return false
  const contentType = doc.content_type || doc.type || ''
  const filename = doc.filename || doc.name || ''
  return contentType.includes('pdf') || filename.toLowerCase().endsWith('.pdf')
}

const isImageDocument = (doc: any) => {
  if (!doc) return false
  const contentType = doc.content_type || doc.type || ''
  const filename = doc.filename || doc.name || ''
  return (
    contentType.includes('image') ||
    /\.(jpg|jpeg|png|gif|bmp|webp)$/i.test(filename)
  )
}

const getFileTypeDescription = (doc: any) => {
  if (!doc) return 'Unknown file type'
  const contentType = doc.content_type || doc.type || ''
  const filename = doc.filename || doc.name || ''

  if (contentType.includes('pdf') || filename.toLowerCase().endsWith('.pdf')) {
    return 'PDF Document'
  }
  if (contentType.includes('word') || /\.(doc|docx)$/i.test(filename)) {
    return 'Word Document'
  }
  if (contentType.includes('excel') || /\.(xls|xlsx)$/i.test(filename)) {
    return 'Excel Spreadsheet'
  }
  if (contentType.includes('image')) {
    return 'Image File'
  }

  return 'Document File'
}

const loadSupportingDocuments = () => {
  if (props.claim?.attachments && Array.isArray(props.claim.attachments)) {
    supportingDocuments.value = props.claim.attachments
  } else {
    // Fallback to mock data if no attachments or not an array
    supportingDocuments.value = mockDocuments
  }
}

// Watch for claim changes to reload data
watch(
  () => props.claim,
  () => {
    loadSupportingDocuments()
    loadCommunications()
  },
  { deep: true }
)

// Cleanup blob URLs to prevent memory leaks
const cleanupBlobUrls = () => {
  documentBlobUrls.value.forEach((url) => {
    URL.revokeObjectURL(url)
  })
  documentBlobUrls.value = []
}

// Load data when component mounts
onMounted(() => {
  loadAssessments()
  loadSupportingDocuments()
  loadCommunications()
})

// Cleanup when component unmounts
onBeforeUnmount(() => {
  cleanupBlobUrls()
})
</script>

<style scoped>
.cursor-pointer {
  cursor: pointer;
}

.gap-2 > * + * {
  margin-top: 8px;
}
</style>
