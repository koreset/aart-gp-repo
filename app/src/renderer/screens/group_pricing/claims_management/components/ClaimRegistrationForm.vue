<template>
  <!-- Missing Documents Confirmation Dialog -->
  <v-dialog v-model="showMissingDocsDialog" max-width="500px" persistent>
    <v-card>
      <v-card-title class="text-h6 bg-warning text-white">
        <v-icon start>mdi-alert-triangle</v-icon>
        Missing Required Documents
      </v-card-title>
      <v-card-text class="pt-4">
        <p class="mb-3">
          You are about to submit this claim without the following required
          documents:
        </p>
        <v-list density="compact">
          <v-list-item
            v-for="docType in missingRequiredDocuments"
            :key="docType.code"
            prepend-icon="mdi-file-alert"
          >
            <v-list-item-title>{{ docType.name }}</v-list-item-title>
          </v-list-item>
        </v-list>
        <v-alert type="info" variant="tonal" class="mt-4">
          <strong>Note:</strong> Claims with missing documents may experience
          delays in processing. You can upload these documents later during the
          claims assessment process.
        </v-alert>
      </v-card-text>
      <v-card-actions class="justify-end pa-4">
        <v-btn
          color="grey"
          variant="outlined"
          @click="showMissingDocsDialog = false"
        >
          Cancel
        </v-btn>
        <v-btn color="warning" @click="confirmSubmitWithMissingDocs">
          Submit Anyway
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-form ref="form" @submit.prevent="handleSubmit">
    <v-row>
      <!-- Member Information Section -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Member Information
          </v-card-title>
          <v-card-text class="pt-4">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.member_id_number"
                  label="Member ID/Passport Number *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required, rules.idOrPassport]"
                  required
                  @blur="lookupMember"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="memberInfo.member_name"
                  label="Member Name"
                  variant="outlined"
                  density="compact"
                  readonly
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="memberInfo.scheme_name"
                  label="Scheme"
                  variant="outlined"
                  density="compact"
                  readonly
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="memberInfo.annual_salary"
                  label="Annual Salary"
                  prefix="R"
                  variant="outlined"
                  density="compact"
                  readonly
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Claim Details Section -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Claim Information
          </v-card-title>
          <v-card-text class="pt-4">
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.member_type"
                  :items="memberTypes"
                  label="Member Type *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.benefit_type"
                  :items="benefitTypes"
                  label="Benefit Type *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                  item-title="title"
                  item-value="value"
                  return-object
                  @update:model-value="calculateClaimAmount"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.date_of_event"
                  label="Date of Event *"
                  type="date"
                  variant="outlined"
                  density="compact"
                  :rules="[
                    rules.required,
                    rules.dateNotFuture,
                    rules.eventDateMin
                  ]"
                  :hint="eventDateWarning"
                  persistent-hint
                  hint-class="warning-hint"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.date_notified"
                  label="Date Notified *"
                  type="date"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="formData.claim_amount"
                  label="Claim Amount *"
                  type="number"
                  prefix="R"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required, rules.claimAmount]"
                  :loading="claimAmountLoading"
                  readonly
                  required
                  :append-inner-icon="
                    claimAmountLoading ? 'mdi-loading mdi-spin' : undefined
                  "
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.priority"
                  :items="priorities"
                  label="Priority *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.cause_type"
                  :items="accidentTypes"
                  label="Claim Type *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Claimant Information Section (for dependants) -->
      <v-col v-if="formData.member_type !== 'member'" cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Claimant Information
          </v-card-title>
          <v-card-text class="pt-4">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.claimant_name"
                  label="Claimant Name *"
                  variant="outlined"
                  density="compact"
                  :rules="
                    formData.member_type !== 'member' ? [rules.required] : []
                  "
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.claimant_id_number"
                  label="Claimant ID Number *"
                  variant="outlined"
                  density="compact"
                  :rules="
                    formData.member_type !== 'member'
                      ? [rules.required, rules.idOrPassport]
                      : []
                  "
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.relationship_to_member"
                  :items="relationships"
                  label="Relationship to Member *"
                  variant="outlined"
                  density="compact"
                  :rules="
                    formData.member_type !== 'member' ? [rules.required] : []
                  "
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.claimant_contact_number"
                  label="Contact Number"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Banking Details Section -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title
            class="text-subtitle-1 bg-grey-lighten-4 d-flex justify-space-between align-center"
          >
            Banking Details
            <v-chip
              v-if="formData.bank_verification_status === 'verified'"
              color="success"
              size="small"
              variant="flat"
            >
              <v-icon size="small" class="mr-1">mdi-check-circle</v-icon>
              Verified
            </v-chip>
            <v-chip
              v-else-if="formData.bank_verification_status === 'failed'"
              color="error"
              size="small"
              variant="flat"
            >
              <v-icon size="small" class="mr-1">mdi-alert-circle</v-icon>
              Verification Failed
            </v-chip>
          </v-card-title>
          <v-card-text class="pt-4">
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.bank_name"
                  :items="bankOptions"
                  label="Bank Name *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  @update:model-value="onBankSelected"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.bank_branch_code"
                  label="Branch Code (Universal)"
                  variant="outlined"
                  density="compact"
                  :readonly="formData.bank_name !== 'Other'"
                  hint="Auto-filled from bank selection"
                  persistent-hint
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.bank_account_number"
                  label="Account Number *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required, rules.numericOnly]"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.bank_account_type"
                  :items="accountTypeOptions"
                  label="Account Type *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.account_holder_name"
                  label="Account Holder Name *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  hint="Pre-filled with claimant/member name"
                  persistent-hint
                />
              </v-col>
              <v-col cols="12" md="6" class="d-flex align-center">
                <v-btn
                  color="primary"
                  variant="outlined"
                  size="small"
                  prepend-icon="mdi-bank-check"
                  :loading="bankVerifying"
                  :disabled="!canVerifyBanking"
                  @click="verifyBankingDetails"
                >
                  Verify Account
                </v-btn>
              </v-col>
            </v-row>

            <!-- Verification result -->
            <v-expand-transition>
              <div v-if="formData.bank_verification_status">
                <v-divider class="my-3" />
                <v-alert
                  v-if="formData.bank_verification_status === 'verified'"
                  type="success"
                  variant="tonal"
                  density="compact"
                >
                  <div class="d-flex justify-space-between align-center">
                    <div>
                      <strong>Account verified successfully.</strong>
                      Account holder name matches. Account is open and active.
                    </div>
                    <div class="text-caption text-medium-emphasis">
                      {{ formData.bank_verification_date }}
                      <span v-if="formData.bank_verification_reference">
                        &middot; Ref: {{ formData.bank_verification_reference }}
                      </span>
                    </div>
                  </div>
                </v-alert>
                <v-alert
                  v-else-if="formData.bank_verification_status === 'failed'"
                  type="error"
                  variant="tonal"
                  density="compact"
                >
                  <div class="d-flex justify-space-between align-center">
                    <div>
                      <strong>Verification failed.</strong>
                      {{ bankVerificationError }}
                    </div>
                    <v-btn
                      size="x-small"
                      variant="text"
                      color="error"
                      @click="verifyBankingDetails"
                    >
                      Retry
                    </v-btn>
                  </div>
                </v-alert>
                <v-alert
                  v-else-if="formData.bank_verification_status === 'pending'"
                  type="warning"
                  variant="tonal"
                  density="compact"
                >
                  Verification in progress...
                </v-alert>

                <!--
                  Per-check TriState breakdown. 'unknown' renders with a
                  warning question mark so users can distinguish "bank
                  could not confirm" from a hard "no".
                -->
                <div
                  v-if="verificationChecks.length > 0"
                  class="mt-3 d-flex flex-wrap ga-2"
                >
                  <v-chip
                    v-for="check in verificationChecks"
                    :key="check.label"
                    :color="check.color"
                    size="small"
                    variant="tonal"
                  >
                    <v-icon start size="small">{{ check.icon }}</v-icon>
                    {{ check.label }}: {{ check.stateLabel }}
                  </v-chip>
                </div>
              </div>
            </v-expand-transition>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Supporting Documentation -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Supporting Documentation
            <v-chip
              v-if="formData.benefit_type"
              size="small"
              color="primary"
              variant="tonal"
              class="ml-2"
            >
              {{ formData.benefit_type.value?.benefit_alias }}
            </v-chip>
          </v-card-title>
          <v-card-text class="pt-4">
            <v-alert
              v-if="!formData.benefit_type"
              type="info"
              variant="tonal"
              class="mb-4"
            >
              Please select a benefit type to see required documents.
            </v-alert>

            <div v-else>
              <!-- Document upload sections for each required document type -->
              <v-row
                v-for="docType in requiredDocumentTypes"
                :key="docType.code"
              >
                <v-col cols="12">
                  <v-card variant="outlined" class="mb-3">
                    <v-card-title
                      class="text-body-1 py-3 px-3 d-flex align-center"
                      :class="
                        docType.required
                          ? 'bg-blue-lighten-5'
                          : 'bg-grey-lighten-2'
                      "
                    >
                      <v-icon
                        :color="docType.required ? 'primary' : 'grey-darken-1'"
                        size="small"
                        class="mr-2"
                      >
                        {{
                          docType.required
                            ? 'mdi-asterisk'
                            : 'mdi-circle-outline'
                        }}
                      </v-icon>
                      {{ docType.name }}
                      <v-chip
                        v-if="docType.required"
                        size="x-small"
                        color="primary"
                        variant="flat"
                        text-color="white"
                        class="ml-2"
                      >
                        Required
                      </v-chip>
                      <v-chip
                        v-else
                        size="x-small"
                        color="grey-darken-1"
                        variant="flat"
                        text-color="white"
                        class="ml-2"
                      >
                        Optional
                      </v-chip>
                      <v-spacer />
                      <v-icon
                        v-if="
                          formData.supporting_documents[docType.code]?.length >
                          0
                        "
                        color="success"
                        size="small"
                      >
                        mdi-check-circle
                      </v-icon>
                    </v-card-title>
                    <v-card-text class="pt-3 pb-3">
                      <v-file-input
                        :model-value="null"
                        :label="`Upload ${docType.name}`"
                        multiple
                        accept=".pdf,.jpg,.jpeg,.png,.doc,.docx"
                        variant="outlined"
                        density="comfortable"
                        prepend-icon="mdi-paperclip"
                        show-size
                        color="primary"
                        @update:model-value="
                          (files) => handleDocumentUpload(docType.code, files)
                        "
                      />

                      <!-- Display uploaded files for this document type -->
                      <div
                        v-if="
                          formData.supporting_documents[docType.code]?.length >
                          0
                        "
                        class="mt-2"
                      >
                        <v-chip-group column>
                          <v-chip
                            v-for="(file, fileIndex) in formData
                              .supporting_documents[docType.code]"
                            :key="`${docType.code}-${file.name}-${fileIndex}`"
                            closable
                            color="primary"
                            variant="tonal"
                            size="small"
                            @click:close="
                              removeDocumentFile(docType.code, fileIndex)
                            "
                          >
                            <v-icon start size="small"
                              >mdi-file-document</v-icon
                            >
                            {{ file.name }} ({{ formatFileSize(file.size) }})
                          </v-chip>
                        </v-chip-group>
                      </div>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>

              <!-- Document completion summary -->
              <v-alert
                v-if="allRequiredDocumentsUploaded"
                type="success"
                variant="tonal"
                class="mt-4"
              >
                <v-icon start>mdi-check-circle</v-icon>
                All required documents have been uploaded.
              </v-alert>
              <v-alert
                v-else-if="requiredDocumentTypes.some((d) => d.required)"
                type="warning"
                variant="tonal"
                class="mt-4"
              >
                <v-icon start>mdi-alert-circle</v-icon>
                {{ missingRequiredDocuments.length }} required document(s) still
                need to be uploaded:
                <ul class="mt-1">
                  <li
                    v-for="docType in missingRequiredDocuments"
                    :key="docType.code"
                    class="text-caption"
                  >
                    {{ docType.name }}
                  </li>
                </ul>
              </v-alert>
            </div>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Claim Description -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="text-subtitle-1 bg-grey-lighten-4">
            Claim Description
          </v-card-title>
          <v-card-text class="pt-4">
            <v-textarea
              v-model="formData.description"
              label="Detailed Description of Claim"
              variant="outlined"
              rows="4"
            />
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Registration blockers -->
      <v-col v-if="!isFormValid && formData.member_id_number" cols="12">
        <v-alert type="warning" variant="tonal" density="compact" class="mb-0">
          <div class="text-subtitle-2 mb-1"
            >Registration blocked — resolve the following:</div
          >
          <ul class="text-caption pl-4">
            <li v-if="!bankingDetailsComplete">
              <strong>Banking details:</strong>
              {{
                !formData.bank_name ||
                !formData.bank_account_number ||
                !formData.bank_account_type ||
                !formData.account_holder_name
                  ? 'Complete all banking fields'
                  : 'Verify the account (click "Verify Account")'
              }}
            </li>
            <li v-for="doc in missingHardRequiredDocs" :key="doc.code">
              <strong>Upload required:</strong> {{ doc.name }}
            </li>
          </ul>
        </v-alert>
      </v-col>

      <!-- Action Buttons -->
      <v-col cols="12">
        <v-card-actions class="justify-end">
          <v-btn color="grey" variant="outlined" @click="$emit('cancel')"
            >Cancel</v-btn
          >
          <v-btn
            color="primary"
            type="submit"
            :loading="loading"
            :disabled="!isFormValid"
          >
            Register Claim
          </v-btn>
        </v-card-actions>
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import {
  triStateIcon,
  type VerifyResult,
  type TriState
} from '@/renderer/types/bav'

interface Props {
  schemes: Array<any>
}

interface Emits {
  (e: 'save', claim: any): void
  (e: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const form = ref(null)
const loading = ref(false)
const claimAmountLoading = ref(false)
const serverBenefitsResponse = ref<string[]>([])
const showMissingDocsDialog = ref(false)

// Document types mapping based on benefit codes
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

// Form data
const formData = ref({
  member_id_number: '',
  scheme_id: null,
  benefit_type: null as any,
  member_type: 'member',
  date_of_event: null,
  date_notified: new Date().toISOString().split('T')[0],
  claim_amount: 0,
  priority: 'medium',
  cause_type: 'non-accidental',
  claimant_name: '',
  claimant_id_number: '',
  relationship_to_member: '',
  claimant_contact_number: '',
  bank_name: '',
  bank_branch_code: '',
  bank_account_number: '',
  bank_account_type: '',
  account_holder_name: '',
  bank_verification_status: '' as '' | 'pending' | 'verified' | 'failed',
  bank_verification_date: '',
  bank_verification_reference: '',
  supporting_documents: {} as Record<string, File[]>,
  description: ''
})

// Member information
const memberInfo: any = ref({
  member_name: '',
  annual_salary: 0,
  scheme_name: '',
  scheme_id: null,
  benefits: {},
  scheme_category_details: {}
})

// Form options
const benefitMaps: any = ref([])

const benefitTypes = computed(() => {
  if (!benefitMaps.value || benefitMaps.value.length === 0) {
    return []
  }

  let allowedCodes: string[] = []
  switch (formData.value.member_type) {
    case 'member':
      allowedCodes = ['GLA', 'PTD', 'CI', 'TTD', 'PHI', 'GFF']
      break
    case 'spouse':
      allowedCodes = ['SGLA', 'GFF']
      break
    case 'child':
    case 'parent':
    case 'dependant':
      allowedCodes = ['GFF']
      break
    default:
      allowedCodes = ['GLA', 'SGLA', 'PTD', 'CI', 'TTD', 'PHI', 'GFF']
  }

  return benefitMaps.value
    .filter((benefit: any) => allowedCodes.includes(benefit.benefit_code))
    .map((benefit: any) => {
      const benefitName = benefit.benefit_alias
        ? benefit.benefit_alias
        : benefit.benefit_name
      return {
        title: benefitName,
        value: benefit
      }
    })
})

const memberTypes = [
  { title: 'Member Claim', value: 'member' },
  { title: 'Spouse Claim', value: 'spouse' },
  { title: 'Child Claim', value: 'child' },
  { title: 'Parent Claim', value: 'parent' },
  { title: 'Dependant Claim', value: 'dependant' }
]

const priorities = [
  { title: 'Low', value: 'low' },
  { title: 'Medium', value: 'medium' },
  { title: 'High', value: 'high' },
  { title: 'Critical', value: 'critical' }
]

const accidentTypes = [
  { title: 'Non-accidental', value: 'non-accidental' },
  { title: 'Accidental', value: 'accidental' }
]

const relationships = ['Spouse', 'Child', 'Parent', 'Sibling', 'Other']

const bankOptions = [
  'FNB',
  'Standard Bank',
  'ABSA',
  'Nedbank',
  'Capitec',
  'Investec',
  'African Bank',
  'TymeBank',
  'Discovery Bank',
  'Bank Zero',
  'Other'
]

const universalBranchCodes: Record<string, string> = {
  FNB: '250655',
  'Standard Bank': '051001',
  ABSA: '632005',
  Nedbank: '198765',
  Capitec: '470010',
  Investec: '580105',
  'African Bank': '430000',
  TymeBank: '678910',
  'Discovery Bank': '679000',
  'Bank Zero': '888000'
}

const accountTypeOptions = [
  { title: 'Current/Cheque', value: '1' },
  { title: 'Savings', value: '2' },
  { title: 'Transmission', value: '3' }
]

const onBankSelected = (bankName: string) => {
  if (bankName && bankName !== 'Other') {
    formData.value.bank_branch_code = universalBranchCodes[bankName] || ''
  } else {
    formData.value.bank_branch_code = ''
  }
  // Reset verification when banking details change
  formData.value.bank_verification_status = ''
  formData.value.bank_verification_date = ''
  formData.value.bank_verification_reference = ''
}

// Banking verification
const bankVerifying = ref(false)
const bankVerificationError = ref('')

const canVerifyBanking = computed(() => {
  return (
    formData.value.bank_name &&
    formData.value.bank_account_number &&
    formData.value.bank_account_type &&
    formData.value.account_holder_name &&
    !bankVerifying.value
  )
})

const accountTypeMap: Record<string, string> = {
  '1': 'Cheque',
  '2': 'Savings',
  '3': 'Transmission'
}

const lastVerification = ref<VerifyResult | null>(null)

// describeIssue maps a TriState for a specific check to a short problem
// message. 'no' is a hard negative; 'unknown' is flagged so the user knows
// the bank couldn't confirm the check either way.
const describeIssue = (
  label: string,
  notFoundMsg: string,
  t: TriState
): string | null => {
  if (t === 'no') return notFoundMsg
  if (t === 'unknown') return `${label}: bank could not confirm`
  return null
}

// Async poll configuration for providers that return status=pending.
// 3s interval / 60s ceiling matches the Phase 6 plan.
const BAV_POLL_INTERVAL_MS = 3000
const BAV_POLL_TIMEOUT_MS = 60_000

/**
 * pollUntilResolved repeatedly hits the status endpoint every
 * BAV_POLL_INTERVAL_MS ms until the result is no longer pending or the
 * BAV_POLL_TIMEOUT_MS deadline is reached. Returns the last observed result.
 */
const pollUntilResolved = async (
  jobId: string,
  initial: VerifyResult
): Promise<VerifyResult> => {
  let current = initial
  const deadline = Date.now() + BAV_POLL_TIMEOUT_MS
  while (current.status === 'pending' && Date.now() < deadline) {
    await new Promise((resolve) => setTimeout(resolve, BAV_POLL_INTERVAL_MS))
    const res = await GroupPricingService.getBankVerificationStatus(jobId)
    current = res.data?.data as VerifyResult
    lastVerification.value = current
  }
  return current
}

/**
 * Verify banking details via the provider-agnostic BAV v2 endpoint.
 * Handles sync providers (immediate complete/failed) and async providers
 * (status=pending → poll until resolved or timeout).
 */
const verifyBankingDetails = async () => {
  if (!canVerifyBanking.value) return

  bankVerifying.value = true
  formData.value.bank_verification_status = 'pending'
  bankVerificationError.value = ''
  lastVerification.value = null

  try {
    const nameParts = (formData.value.account_holder_name || '')
      .trim()
      .split(/\s+/)
    const firstName = nameParts.slice(0, -1).join(' ') || nameParts[0] || ''
    const surname = nameParts.length > 1 ? nameParts[nameParts.length - 1] : ''

    const idNumber =
      formData.value.claimant_id_number || formData.value.member_id_number
    const identityType = /^\d{13}$/.test(idNumber) ? 'IDNumber' : 'Passport'

    const res = await GroupPricingService.verifyBankAccount({
      first_name: firstName,
      surname,
      identity_number: idNumber,
      identity_type: identityType,
      bank_account_number: formData.value.bank_account_number,
      bank_branch_code: formData.value.bank_branch_code,
      bank_account_type:
        accountTypeMap[formData.value.bank_account_type] ||
        formData.value.bank_account_type
    })

    let result = res.data?.data as VerifyResult
    lastVerification.value = result

    if (result.status === 'pending' && result.providerJobId) {
      result = await pollUntilResolved(result.providerJobId, result)
    }

    if (result.status === 'pending') {
      formData.value.bank_verification_status = 'failed'
      bankVerificationError.value =
        'Verification is still in progress. Please try again in a few minutes.'
      return
    }

    if (result.verified) {
      const now = new Date()
      formData.value.bank_verification_status = 'verified'
      formData.value.bank_verification_date =
        now.toISOString().split('T')[0] +
        ' ' +
        now.toLocaleTimeString('en-ZA', { hour: '2-digit', minute: '2-digit' })
      formData.value.bank_verification_reference = result.providerRequestId
    } else {
      formData.value.bank_verification_status = 'failed'
      const issues = [
        describeIssue(
          'Account found',
          'Account not found',
          result.accountFound
        ),
        describeIssue(
          'Account open',
          'Account is not active',
          result.accountOpen
        ),
        describeIssue(
          'Identity match',
          'Identity does not match account holder',
          result.identityMatch
        ),
        describeIssue(
          'Account type',
          'Account type mismatch',
          result.accountTypeMatch
        )
      ].filter((x): x is string => x !== null)
      bankVerificationError.value =
        issues.length > 0
          ? issues.join('. ') + '.'
          : result.summary || 'Verification failed.'
    }
  } catch (error: any) {
    formData.value.bank_verification_status = 'failed'
    bankVerificationError.value =
      error?.response?.data?.message ||
      error?.data?.message ||
      error?.message ||
      'Verification service unavailable. Please try again or verify manually.'
  } finally {
    bankVerifying.value = false
  }
}

// verificationChecks surfaces the per-field TriState results for display.
// Each row uses triStateIcon() to pick a tick / cross / question icon so
// 'unknown' is visually distinct from a hard 'no'.
const verificationChecks = computed(() => {
  const r = lastVerification.value
  if (!r) return []
  const rows: Array<{ label: string; state: TriState }> = [
    { label: 'Account found', state: r.accountFound },
    { label: 'Account open', state: r.accountOpen },
    { label: 'Identity match', state: r.identityMatch },
    { label: 'Account type match', state: r.accountTypeMatch },
    { label: 'Accepts credits', state: r.acceptsCredits }
  ]
  return rows.map((row) => {
    const presentation = triStateIcon(row.state)
    return {
      label: row.label,
      icon: presentation.icon,
      color: presentation.color,
      stateLabel: presentation.label
    }
  })
})

// Required document types for current benefit
const requiredDocumentTypes = computed(() => {
  if (!formData.value.benefit_type?.value?.benefit_code) {
    return []
  }

  const benefitCode = formData.value.benefit_type.value.benefit_code
  return documentTypesMapping[benefitCode] || []
})

// Validation rules
const rules = {
  required: (value: any) => !!value || 'Field is required',
  idOrPassport: (value: string) => {
    if (!value) return true
    // RSA ID number: 13 digits
    if (value.length === 13 && /^\d{13}$/.test(value)) {
      return true
    }
    // Passport: 6-12 characters, letters and numbers
    if (
      value.length >= 6 &&
      value.length <= 12 &&
      /^[A-Za-z0-9]+$/.test(value)
    ) {
      return true
    }
    return 'Invalid ID number or passport format'
  },
  dateNotFuture: (value: string) => {
    if (!value) return true
    const eventDate = new Date(value)
    const today = new Date()
    return eventDate <= today || 'Date cannot be in the future'
  },
  claimAmount: (value: number) => {
    if (!value) return 'Claim amount is required'
    return value > 0 || 'Claim amount must be greater than 0'
  },
  numericOnly: (value: string) => {
    if (!value) return true
    return /^\d+$/.test(value) || 'Only numeric characters allowed'
  },
  eventDateMin: (value: string) => {
    if (!value) return true
    const eventDate = new Date(value)
    const minDate = new Date(memberScheme.value.commencement_date)
    return eventDate >= minDate || 'Date cannot be before scheme commencement'
  }
}

// Computed properties
// Documents that absolutely must be uploaded before registration (no bypass)
const hardRequiredDocCodes = [
  'claim_form',
  'banking_details',
  'certified_id_deceased',
  'certified_id_claimant',
  'certified_id_member'
]

const missingHardRequiredDocs = computed(() => {
  return requiredDocumentTypes.value.filter(
    (doc) =>
      hardRequiredDocCodes.includes(doc.code) &&
      (!formData.value.supporting_documents[doc.code] ||
        formData.value.supporting_documents[doc.code].length === 0)
  )
})

const bankingDetailsComplete = computed(() => {
  return (
    formData.value.bank_name &&
    formData.value.bank_account_number &&
    formData.value.bank_account_type &&
    formData.value.account_holder_name &&
    formData.value.bank_verification_status === 'verified'
  )
})

const isFormValid = computed(() => {
  const benefitTypeValid =
    formData.value.benefit_type &&
    (formData.value.benefit_type.value
      ? formData.value.benefit_type.value.benefit_code &&
        formData.value.benefit_type.value.benefit_code !== ''
      : formData.value.benefit_type.benefit_code &&
        formData.value.benefit_type.benefit_code !== '')

  return (
    formData.value.member_id_number &&
    formData.value.scheme_id &&
    benefitTypeValid &&
    formData.value.member_type &&
    formData.value.date_of_event &&
    formData.value.date_notified &&
    formData.value.claim_amount > 0 &&
    formData.value.priority &&
    bankingDetailsComplete.value &&
    missingHardRequiredDocs.value.length === 0 &&
    (formData.value.member_type === 'member' ||
      (formData.value.claimant_name &&
        formData.value.claimant_id_number &&
        formData.value.relationship_to_member))
  )
})

const memberScheme = ref()

// Document validation computed properties
const allRequiredDocumentsUploaded = computed(() => {
  const requiredDocs = requiredDocumentTypes.value.filter((doc) => doc.required)
  return requiredDocs.every(
    (doc) => formData.value.supporting_documents[doc.code]?.length > 0
  )
})

const missingRequiredDocuments = computed(() => {
  return requiredDocumentTypes.value.filter(
    (doc) =>
      doc.required &&
      (!formData.value.supporting_documents[doc.code] ||
        formData.value.supporting_documents[doc.code].length === 0)
  )
})

const eventDateWarning = computed(() => {
  if (formData.value.date_of_event !== null) {
    const eventDate = new Date(formData.value.date_of_event)
    const coverStart = new Date(memberScheme.value.cover_start_date)

    if (eventDate < coverStart) {
      return "Scheme's cover start date has not passed"
    }
    return ''
  }
  return ''
})
// Methods
const lookupMember = async () => {
  if (
    !formData.value.member_id_number ||
    formData.value.member_id_number.length < 6
  ) {
    return
  }

  try {
    // This would be replaced with actual API call
    const response = await GroupPricingService.getMemberByIdNumber(
      formData.value.member_id_number
    )

    serverBenefitsResponse.value = response.data.benefits
    memberInfo.value = response.data
    formData.value.scheme_id = memberInfo.value.scheme_id
    // Pre-fill account holder name with member name
    if (!formData.value.account_holder_name) {
      formData.value.account_holder_name = memberInfo.value.member_name || ''
    }
    // Set scheme based on member data
    memberScheme.value = props.schemes.find(
      (s) => s.name === memberInfo.value.scheme_name
    )

    if (memberScheme.value) {
      formData.value.scheme_id = memberInfo.value.scheme_id
    }

    const res = await GroupPricingService.getBenefitMapsBySchemeAndCategory(
      memberScheme.value.id,
      memberInfo.value.scheme_category_details.id
    )
    benefitMaps.value = res.data

    await calculateClaimAmount()
  } catch (error) {
    console.error('Error looking up member:', error)
    memberInfo.value = {
      member_name: '',
      annual_salary: 0,
      scheme_name: '',
      scheme_id: null,
      benefits: {},
      scheme_category_details: {}
    }
  }
}

const fetchUpdatedClaimAmount = async () => {
  const benefitData =
    formData.value.benefit_type?.value || formData.value.benefit_type

  if (
    !formData.value.member_id_number ||
    !benefitData ||
    !benefitData.benefit_name ||
    !benefitData.benefit_code ||
    !formData.value.member_type
  ) {
    return
  }

  claimAmountLoading.value = true
  try {
    const response = await GroupPricingService.getUpdatedClaimAmount({
      member_id_number: formData.value.member_id_number,
      scheme_id: formData.value.scheme_id,
      benefit_type: benefitData.benefit_name,
      benefit_code: benefitData.benefit_code,
      member_type: formData.value.member_type
    })

    if (response.data && response.data.updated_claim_amount !== null) {
      formData.value.claim_amount = response.data.updated_claim_amount

      // Show warning if amount was reduced due to previous claims
      if (response.data.previous_claims_amount > 0) {
        console.warn(
          `Claim amount adjusted from R${response.data.original_amount} to R${response.data.updated_claim_amount} due to previous claims totaling R${response.data.previous_claims_amount}`
        )
      }
    } else {
      // Fallback to original calculation if API fails
      calculateClaimAmountLocal()
    }
  } catch (error) {
    console.error('Error fetching updated claim amount:', error)
    // Fallback to original calculation
    calculateClaimAmountLocal()
  } finally {
    claimAmountLoading.value = false
  }
}

const calculateClaimAmountLocal = () => {
  const benefitData =
    formData.value.benefit_type?.value || formData.value.benefit_type

  if (
    !benefitData ||
    !benefitData.benefit_code ||
    !memberInfo.value.annual_salary
  ) {
    formData.value.claim_amount = 0
    return
  }

  const benefits = memberInfo.value.benefits as any
  let multiple = 0
  const benefitCode = benefitData.benefit_code

  switch (benefitCode) {
    case 'GLA':
      multiple = benefits.gla_multiple || 0
      break
    case 'SGLA':
      multiple = benefits.sgla_multiple || 0
      break
    case 'PTD':
      multiple = benefits.ptd_multiple || 0
      break
    case 'CI':
      multiple = benefits.ci_multiple || 0
      break
    case 'TTD':
      multiple = benefits.ttd_multiple || 0
      break
    case 'PHI':
      multiple = benefits.phi_multiple || 0
      break
    case 'GFF':
      // Fixed amount for funeral benefit
      switch (formData.value.member_type) {
        case 'member':
          formData.value.claim_amount =
            memberInfo.value.scheme_category_details.family_funeral_main_member_funeral_sum_assured
          break
        case 'spouse':
          formData.value.claim_amount =
            memberInfo.value.scheme_category_details.family_funeral_spouse_funeral_sum_assured
          break
        case 'child':
          formData.value.claim_amount =
            memberInfo.value.scheme_category_details.family_funeral_children_funeral_sum_assured
          break
        case 'parent':
          formData.value.claim_amount =
            memberInfo.value.scheme_category_details.family_funeral_parent_funeral_sum_assured
          break
        case 'dependant':
          formData.value.claim_amount =
            memberInfo.value.scheme_category_details.family_funeral_adult_dependant_sum_assured
          break
        default:
          formData.value.claim_amount = 0
      }
      return
  }

  formData.value.claim_amount = Math.round(
    memberInfo.value.annual_salary * multiple
  )
}

const calculateClaimAmount = async () => {
  // Use API to get updated amount that accounts for previous claims
  if (formData.value.member_id_number && formData.value.benefit_type) {
    await fetchUpdatedClaimAmount()
  } else {
    // Fallback to local calculation
    calculateClaimAmountLocal()
  }
}

const handleDocumentUpload = (
  documentType: string,
  newFiles: File | File[] | null
) => {
  if (newFiles) {
    // Normalize to array format
    const filesArray = Array.isArray(newFiles) ? newFiles : [newFiles]

    if (filesArray.length > 0) {
      // Initialize array if it doesn't exist
      if (!formData.value.supporting_documents[documentType]) {
        formData.value.supporting_documents[documentType] = []
      }

      // Add new files to existing ones for this document type, avoiding duplicates
      const existingFileNames = formData.value.supporting_documents[
        documentType
      ].map((f: File) => f.name)
      const uniqueNewFiles = filesArray.filter(
        (file: File) => !existingFileNames.includes(file.name)
      )

      formData.value.supporting_documents[documentType].push(...uniqueNewFiles)
    }
  }
}

const removeDocumentFile = (documentType: string, fileIndex: number) => {
  if (formData.value.supporting_documents[documentType]) {
    formData.value.supporting_documents[documentType].splice(fileIndex, 1)

    // Clean up empty arrays
    if (formData.value.supporting_documents[documentType].length === 0) {
      delete formData.value.supporting_documents[documentType]
    }
  }
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const handleSubmit = async () => {
  if (!form.value) return

  const { valid } = await (form.value as any).validate()
  if (!valid) return

  // Check for soft-required documents (excludes hard-required which are already blocked by isFormValid)
  const softMissing = missingRequiredDocuments.value.filter(
    (doc) => !hardRequiredDocCodes.includes(doc.code)
  )
  if (softMissing.length > 0) {
    showMissingDocsDialog.value = true
    return
  }

  // Proceed with submission if all required documents are present
  await submitClaim()
}

const confirmSubmitWithMissingDocs = async () => {
  showMissingDocsDialog.value = false
  await submitClaim()
}

const submitClaim = async () => {
  loading.value = true
  try {
    // Generate claim number
    const claimNumber = `CLM-${new Date().getFullYear()}-${String(Date.now()).slice(-6)}`

    // Create FormData for multipart submission
    const formDataForSubmission = new FormData()

    // Add basic claim data
    const claimBasicData = {
      ...formData.value,
      claim_number: claimNumber,
      member_name: memberInfo.value.member_name,
      scheme_name: memberInfo.value.scheme_name,
      benefit_name:
        formData.value.benefit_type?.value?.benefit_name ||
        formData.value.benefit_type?.benefit_name ||
        '',
      benefit_code:
        formData.value.benefit_type?.value?.benefit_code ||
        formData.value.benefit_type?.benefit_code ||
        '',
      benefit_alias:
        formData.value.benefit_type?.value?.benefit_alias ||
        formData.value.benefit_type?.benefit_alias ||
        '',
      status: 'pending',
      date_registered: new Date().toISOString().split('T')[0],
      missing_required_documents: missingRequiredDocuments.value.map(
        (doc) => doc.name
      )
    }

    // Remove supporting_documents from basic data since we'll handle files separately
    const { supporting_documents: _, ...claimDataWithoutFiles } = claimBasicData

    // Add non-file data as JSON string
    formDataForSubmission.append(
      'claim_data',
      JSON.stringify(claimDataWithoutFiles)
    )

    // Add files and their metadata separately
    const documentMetadata: Array<{
      document_type: string
      document_name: string
      file_index: number
    }> = []

    let fileIndex = 0
    for (const [documentType, files] of Object.entries(
      formData.value.supporting_documents
    )) {
      for (const file of files) {
        // Add the actual file
        formDataForSubmission.append(`files`, file)

        // Track metadata for this file
        const benefitCode =
          formData.value.benefit_type?.value?.benefit_code ||
          formData.value.benefit_type?.benefit_code ||
          ''
        documentMetadata.push({
          document_type: documentType,
          document_name:
            documentTypesMapping[benefitCode]?.find(
              (d) => d.code === documentType
            )?.name || documentType,
          file_index: fileIndex
        })
        fileIndex++
      }
    }

    // // Add document metadata as JSON
    formDataForSubmission.append(
      'document_metadata',
      JSON.stringify(documentMetadata)
    )

    emit('save', formDataForSubmission)
  } finally {
    loading.value = false
  }
}

// Reset bank verification when key banking fields change
watch(
  () => [
    formData.value.bank_account_number,
    formData.value.account_holder_name,
    formData.value.bank_account_type
  ],
  () => {
    if (formData.value.bank_verification_status) {
      formData.value.bank_verification_status = ''
      formData.value.bank_verification_date = ''
      formData.value.bank_verification_reference = ''
    }
  }
)

// Watchers
watch(
  () => formData.value.benefit_type,
  async () => {
    await calculateClaimAmount()
  }
)

watch(
  () => formData.value.member_type,
  async (newType) => {
    if (newType === 'member') {
      // Clear claimant fields when switching to member claim
      formData.value.claimant_name = ''
      formData.value.claimant_id_number = ''
      formData.value.relationship_to_member = ''
      formData.value.claimant_contact_number = ''
    }

    // Reset benefit type if it's not valid for the new claim type
    const availableBenefits = benefitTypes.value
    const benefitData =
      formData.value.benefit_type?.value || formData.value.benefit_type

    if (
      formData.value.benefit_type &&
      !availableBenefits.some((b) => {
        const currentBenefitCode = benefitData?.benefit_code || ''
        return b.value.benefit_code === currentBenefitCode
      })
    ) {
      formData.value.benefit_type = null
      formData.value.claim_amount = 0
      // Clear documents when benefit type changes
      formData.value.supporting_documents = {}
    } else if (formData.value.benefit_type) {
      // Recalculate claim amount for new claim type
      await calculateClaimAmount()
    }
  }
)

onMounted(async () => {
  const res = await GroupPricingService.getBenefitMaps()
  benefitMaps.value = res.data
  // Initial member lookup if ID number is pre-filled
  if (formData.value.member_id_number) {
    lookupMember()
  }
})
</script>

<style scoped>
.v-card-title {
  padding: 12px 16px;
}

.v-card-text {
  padding: 16px;
}

ul {
  padding-left: 20px;
}

li {
  margin-bottom: 4px;
}
</style>
