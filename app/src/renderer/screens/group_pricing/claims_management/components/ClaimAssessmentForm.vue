<template>
  <v-row>
    <!-- Claim Details  -->
    <v-col col="12">
      <v-card variant="outlined" class="mb-4">
        <v-card-title class="bg-primary text-white">
          Claim Details</v-card-title
        >
        <v-card-text>
          <v-row>
            <v-col cols="12" md="6">
              <strong>Member Type: </strong> {{ claim?.member_type }} <br />
            </v-col>
            <v-col>
              <strong>Date of Transaction: </strong>
              {{ claim?.date_registered }} <br />
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
  <v-form ref="form" @submit.prevent="handleSubmit">
    <v-row>
      <!-- Assessment Summary -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="bg-primary text-white">
            Assessment Summary
          </v-card-title>
          <v-card-text class="pt-4">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.assessor_name"
                  label="Assessor Name *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.assessment_date"
                  label="Assessment Date *"
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
                <v-select
                  v-model="formData.assessment_outcome"
                  :items="assessmentOutcomes"
                  label="Assessment Outcome *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model.number="formData.recommended_amount"
                  label="Recommended Amount"
                  type="number"
                  prefix="R"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
            </v-row>
            <!-- Committee referral fields -->
            <v-expand-transition>
              <div v-if="isCommitteeOutcome">
                <v-divider class="my-3" />
                <v-row>
                  <v-col cols="12">
                    <v-alert
                      v-if="
                        formData.assessment_outcome === 'refer_to_committee'
                      "
                      type="info"
                      variant="tonal"
                      density="compact"
                      class="mb-3"
                    >
                      This claim will be referred to the Claims Committee for
                      adjudication. The claim status will change to "Referred to
                      Committee".
                    </v-alert>
                    <v-alert
                      v-else-if="
                        formData.assessment_outcome === 'committee_approved'
                      "
                      type="success"
                      variant="tonal"
                      density="compact"
                      class="mb-3"
                    >
                      Record the Claims Committee's decision to approve this
                      claim.
                    </v-alert>
                    <v-alert
                      v-else-if="
                        formData.assessment_outcome === 'committee_declined'
                      "
                      type="error"
                      variant="tonal"
                      density="compact"
                      class="mb-3"
                    >
                      Record the Claims Committee's decision to decline this
                      claim.
                    </v-alert>
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="formData.committee_date"
                      :label="
                        formData.assessment_outcome === 'refer_to_committee'
                          ? 'Scheduled Committee Date'
                          : 'Committee Meeting Date'
                      "
                      type="date"
                      variant="outlined"
                      density="compact"
                    />
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-text-field
                      v-model="formData.committee_reason"
                      :label="
                        formData.assessment_outcome === 'refer_to_committee'
                          ? 'Reason for Referral *'
                          : 'Committee Decision Reference'
                      "
                      variant="outlined"
                      density="compact"
                      :rules="
                        formData.assessment_outcome === 'refer_to_committee'
                          ? [rules.required]
                          : []
                      "
                    />
                  </v-col>
                  <v-col
                    v-if="formData.assessment_outcome !== 'refer_to_committee'"
                    cols="12"
                  >
                    <v-textarea
                      v-model="formData.committee_decision_notes"
                      label="Committee Decision Notes *"
                      variant="outlined"
                      rows="3"
                      :rules="[rules.required]"
                      placeholder="Record the committee's rationale, conditions, or stipulations..."
                    />
                  </v-col>
                </v-row>
              </div>
            </v-expand-transition>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Medical Assessment (for relevant claim types) -->
      <v-col v-if="requiresMedicalAssessment" cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="bg-info text-white">
            Medical Assessment
          </v-card-title>
          <v-card-text class="pt-4">
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.medical_officer"
                  label="Medical Officer"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="formData.medical_assessment_date"
                  label="Medical Assessment Date"
                  type="date"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.disability_percentage"
                  :items="disabilityPercentages"
                  label="Disability Percentage"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.medical_condition"
                  :items="medicalConditions"
                  label="Primary Condition"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-textarea
                  v-model="formData.medical_notes"
                  label="Medical Assessment Notes"
                  variant="outlined"
                  rows="3"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Document Verification -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="bg-warning text-white">
            Document Verification
          </v-card-title>
          <v-card-text class="pt-4">
            <v-row>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-if="claim?.benefit_type?.includes('Life')"
                  v-model="formData.documents_verified.death_certificate"
                  label="Death Certificate Verified"
                  color="success"
                />
                <v-checkbox
                  v-if="requiresMedicalAssessment"
                  v-model="formData.documents_verified.medical_reports"
                  label="Medical Reports Verified"
                  color="success"
                />
                <v-checkbox
                  v-model="formData.documents_verified.id_documents"
                  label="ID Documents Verified"
                  color="success"
                />
                <v-checkbox
                  v-model="formData.documents_verified.proof_of_relationship"
                  label="Proof of Relationship Verified"
                  color="success"
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.documents_verified.banking_details"
                  label="Banking Details Verified"
                  color="success"
                />
                <v-checkbox
                  v-model="formData.documents_verified.employment_records"
                  label="Employment Records Verified"
                  color="success"
                />
                <v-checkbox
                  v-model="formData.documents_verified.claim_form"
                  label="Claim Form Complete"
                  color="success"
                />
                <v-checkbox
                  v-model="formData.documents_verified.beneficiary_nomination"
                  label="Beneficiary Nomination Verified"
                  color="success"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Risk Assessment -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="bg-error text-white">
            Risk Assessment
          </v-card-title>
          <v-card-text class="pt-4">
            <v-row>
              <v-col cols="12" md="6">
                <v-select
                  v-model="formData.fraud_risk_level"
                  :items="riskLevels"
                  label="Fraud Risk Level *"
                  variant="outlined"
                  density="compact"
                  :rules="[rules.required]"
                  required
                />
              </v-col>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.requires_investigation"
                  label="Requires Special Investigation"
                  color="error"
                />
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12">
                <v-textarea
                  v-model="formData.risk_notes"
                  label="Risk Assessment Notes"
                  variant="outlined"
                  rows="2"
                  placeholder="Any concerns or red flags identified during assessment..."
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Assessment Notes -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title class="bg-grey text-white">
            Assessment Notes & Recommendations
          </v-card-title>
          <v-card-text class="pt-4">
            <v-textarea
              v-model="formData.assessment_notes"
              label="Detailed Assessment Notes *"
              variant="outlined"
              rows="4"
              :rules="[rules.required]"
              placeholder="Provide detailed assessment findings, rationale for decision, and any recommendations..."
              required
            />
            <v-textarea
              v-model="formData.next_actions"
              label="Next Actions Required"
              variant="outlined"
              rows="2"
              placeholder="Specify any follow-up actions, additional information needed, or next steps..."
            />
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Assessment Checklist -->
      <v-col cols="12">
        <v-card variant="outlined" class="mb-4">
          <v-card-title
            class="bg-success text-white d-flex justify-space-between align-center"
          >
            Assessment Checklist
            <v-chip
              v-if="autoCheckLoading"
              size="small"
              color="white"
              variant="outlined"
            >
              <v-progress-circular
                indeterminate
                size="14"
                width="2"
                class="mr-2"
              />
              Verifying...
            </v-chip>
            <v-chip
              v-else-if="autoCheckCompleted"
              size="small"
              color="white"
              variant="outlined"
            >
              <v-icon size="small" class="mr-1">mdi-check-circle</v-icon>
              Auto-verified
            </v-chip>
          </v-card-title>
          <v-card-text class="pt-4">
            <!-- Auto-check results banner -->
            <v-alert
              v-if="autoCheckCompleted && autoCheckNotes.length > 0"
              type="info"
              variant="tonal"
              density="compact"
              class="mb-3"
            >
              <div class="text-subtitle-2 mb-1"
                >Pre-assessment verification notes:</div
              >
              <ul class="text-caption pl-4">
                <li v-for="(note, i) in autoCheckNotes" :key="i">{{ note }}</li>
              </ul>
            </v-alert>
            <v-row>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.checklist.policy_in_force"
                  :color="
                    formData.checklist.policy_in_force ? 'success' : 'error'
                  "
                  density="compact"
                >
                  <template #label>
                    <span>Policy was in force at time of event</span>
                    <v-chip
                      v-if="autoCheckCompleted"
                      :color="
                        formData.checklist.policy_in_force ? 'success' : 'error'
                      "
                      size="x-small"
                      variant="tonal"
                      class="ml-2"
                    >
                      {{
                        formData.checklist.policy_in_force
                          ? 'Verified'
                          : 'Failed'
                      }}
                    </v-chip>
                  </template>
                </v-checkbox>
                <v-checkbox
                  v-model="formData.checklist.premiums_paid"
                  :color="
                    formData.checklist.premiums_paid ? 'success' : 'error'
                  "
                  density="compact"
                >
                  <template #label>
                    <span>Premiums up to date</span>
                    <v-chip
                      v-if="autoCheckCompleted"
                      :color="
                        formData.checklist.premiums_paid ? 'success' : 'error'
                      "
                      size="x-small"
                      variant="tonal"
                      class="ml-2"
                    >
                      {{
                        formData.checklist.premiums_paid ? 'Verified' : 'Failed'
                      }}
                    </v-chip>
                  </template>
                </v-checkbox>
                <v-checkbox
                  v-model="formData.checklist.waiting_period_satisfied"
                  :color="
                    formData.checklist.waiting_period_satisfied
                      ? 'success'
                      : 'error'
                  "
                  density="compact"
                >
                  <template #label>
                    <span>Waiting period satisfied</span>
                    <v-chip
                      v-if="autoCheckCompleted"
                      :color="
                        formData.checklist.waiting_period_satisfied
                          ? 'success'
                          : 'error'
                      "
                      size="x-small"
                      variant="tonal"
                      class="ml-2"
                    >
                      {{
                        formData.checklist.waiting_period_satisfied
                          ? 'Verified'
                          : 'Failed'
                      }}
                    </v-chip>
                  </template>
                </v-checkbox>
                <v-checkbox
                  v-model="formData.checklist.benefit_applicable"
                  :color="
                    formData.checklist.benefit_applicable ? 'success' : 'error'
                  "
                  density="compact"
                >
                  <template #label>
                    <span>Benefit applicable for condition/event</span>
                    <v-chip
                      v-if="autoCheckCompleted"
                      :color="
                        formData.checklist.benefit_applicable
                          ? 'success'
                          : 'error'
                      "
                      size="x-small"
                      variant="tonal"
                      class="ml-2"
                    >
                      {{
                        formData.checklist.benefit_applicable
                          ? 'Verified'
                          : 'Failed'
                      }}
                    </v-chip>
                  </template>
                </v-checkbox>
              </v-col>
              <v-col cols="12" md="6">
                <v-checkbox
                  v-model="formData.checklist.exclusions_checked"
                  label="Policy exclusions reviewed"
                  color="success"
                  density="compact"
                />
                <v-checkbox
                  v-model="formData.checklist.fraud_checks_completed"
                  label="Fraud checks completed"
                  color="success"
                  density="compact"
                />
                <v-checkbox
                  v-model="formData.checklist.legal_requirements_met"
                  label="Legal requirements satisfied"
                  color="success"
                  density="compact"
                />
                <v-checkbox
                  v-model="formData.checklist.supervisor_review"
                  label="Supervisor review completed (if required)"
                  color="success"
                  density="compact"
                />
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
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
            Save Assessment
          </v-btn>
        </v-card-actions>
      </v-col>
    </v-row>
  </v-form>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface Props {
  claim: any
}

interface Emits {
  (e: 'save', assessment: any): void
  (e: 'cancel'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const form = ref(null)
const loading = ref(false)
const currentUser = ref<any>(null)

// Form data
const formData = ref({
  assessor_name: (() => {
    try {
      currentUser.value = window.mainApi?.sendSync('msgGetAuthenticatedUser')
      return currentUser.value?.full_name || currentUser.value?.username || ''
    } catch (error) {
      console.warn('Could not get current user for assessor name:', error)
      return ''
    }
  })(),
  assessment_date: new Date().toISOString().split('T')[0],
  assessment_outcome: '',
  recommended_amount: props.claim?.claim_amount || 0,

  // Medical assessment
  medical_officer: '',
  medical_assessment_date: null,
  disability_percentage: null,
  medical_condition: '',
  medical_notes: '',

  // Document verification
  documents_verified: {
    death_certificate: false,
    medical_reports: false,
    id_documents: false,
    proof_of_relationship: false,
    banking_details: false,
    employment_records: false,
    claim_form: false,
    beneficiary_nomination: false
  },

  // Risk assessment
  fraud_risk_level: 'low',
  requires_investigation: false,
  risk_notes: '',

  // Committee referral
  committee_reason: '',
  committee_date: '',
  committee_decision_notes: '',

  // Assessment notes
  assessment_notes: '',
  next_actions: '',

  // Checklist
  checklist: {
    policy_in_force: false,
    premiums_paid: false,
    waiting_period_satisfied: false,
    benefit_applicable: false,
    exclusions_checked: false,
    fraud_checks_completed: false,
    legal_requirements_met: false,
    supervisor_review: false
  }
})

// Form options
const assessmentOutcomes = [
  { title: 'Recommended for Approval', value: 'recommended_approval' },
  {
    title: 'Recommended for Partial Approval',
    value: 'recommended_partial_approval'
  },
  { title: 'Recommended for Decline', value: 'recommended_decline' },
  { title: 'Refer to Claims Committee', value: 'refer_to_committee' },
  { title: 'Committee Approved', value: 'committee_approved' },
  { title: 'Committee Declined', value: 'committee_declined' },
  {
    title: 'Requires Additional Information',
    value: 'requires_additional_info'
  },
  { title: 'Requires Medical Review', value: 'requires_medical_review' },
  { title: 'Requires Investigation', value: 'requires_investigation' }
]

// Show committee reason field when referring to or returning from committee
const isCommitteeOutcome = computed(() => {
  const v = formData.value.assessment_outcome
  return (
    v === 'refer_to_committee' ||
    v === 'committee_approved' ||
    v === 'committee_declined'
  )
})

const riskLevels = [
  { title: 'Low Risk', value: 'low' },
  { title: 'Medium Risk', value: 'medium' },
  { title: 'High Risk', value: 'high' },
  { title: 'Critical Risk', value: 'critical' }
]

const disabilityPercentages = ['0%', '25%', '50%', '75%', '100%']

const medicalConditions = [
  'Cardiovascular Disease',
  'Cancer',
  'Neurological Disorder',
  'Musculoskeletal Disorder',
  'Mental Health Condition',
  'Respiratory Disease',
  'Diabetes',
  'Kidney Disease',
  'Liver Disease',
  'Other'
]

// Validation rules
const rules = {
  required: (value: any) => !!value || 'Field is required'
}

// Computed properties
const requiresMedicalAssessment = computed(() => {
  const benefitType = props.claim?.benefit_type || ''
  return (
    benefitType.includes('Disability') ||
    benefitType.includes('Critical Illness')
  )
})

const isFormValid = computed(() => {
  return (
    formData.value.assessor_name !== null &&
    formData.value.assessment_date &&
    formData.value.assessment_outcome &&
    formData.value.fraud_risk_level &&
    formData.value.assessment_notes
  )
})

// Auto-check state
const autoCheckLoading = ref(false)
const autoCheckCompleted = ref(false)
const autoCheckNotes = ref<string[]>([])

// Waiting period defaults (in days) per benefit code
const waitingPeriodDays: Record<string, number> = {
  GLA: 180,
  SGLA: 180,
  GFF: 180,
  PTD: 365,
  CI: 365,
  TTD: 90,
  PHI: 90
}

// Run pre-assessment verification checks
async function runAutoChecks() {
  if (!props.claim) return

  autoCheckLoading.value = true
  autoCheckCompleted.value = false
  autoCheckNotes.value = []

  try {
    // Fetch scheme data
    let scheme: any = null
    if (props.claim.scheme_id) {
      try {
        const res = await GroupPricingService.getScheme(props.claim.scheme_id)
        scheme = res.data?.data ?? res.data
      } catch {
        autoCheckNotes.value.push(
          'Could not fetch scheme data — manual verification required'
        )
      }
    }

    // Fetch member data
    let member: any = null
    if (props.claim.member_id_number) {
      try {
        const res = await GroupPricingService.getMemberByIdNumber(
          props.claim.member_id_number
        )
        member = res.data?.data ?? res.data
      } catch {
        autoCheckNotes.value.push(
          'Could not fetch member data — manual verification required'
        )
      }
    }

    // 1. Policy in force at time of event
    if (scheme && props.claim.date_of_event) {
      const eventDate = new Date(props.claim.date_of_event)
      const coverStart = scheme.cover_start_date
        ? new Date(scheme.cover_start_date)
        : null
      const coverEnd = scheme.cover_end_date
        ? new Date(scheme.cover_end_date)
        : null
      const schemeInForce =
        scheme.in_force === true || scheme.status?.toLowerCase() === 'in_force'

      if (
        schemeInForce &&
        coverStart &&
        eventDate >= coverStart &&
        (!coverEnd || eventDate <= coverEnd)
      ) {
        formData.value.checklist.policy_in_force = true
        autoCheckNotes.value.push(
          `Policy confirmed in force: cover ${coverStart.toLocaleDateString('en-ZA')} to ${coverEnd ? coverEnd.toLocaleDateString('en-ZA') : 'ongoing'}, event on ${eventDate.toLocaleDateString('en-ZA')}`
        )
      } else {
        formData.value.checklist.policy_in_force = false
        if (!schemeInForce) {
          autoCheckNotes.value.push(
            'Scheme status is not in force — policy check FAILED'
          )
        } else if (coverStart && eventDate < coverStart) {
          autoCheckNotes.value.push(
            `Event date (${eventDate.toLocaleDateString('en-ZA')}) is before cover start (${coverStart.toLocaleDateString('en-ZA')}) — policy check FAILED`
          )
        } else if (coverEnd && eventDate > coverEnd) {
          autoCheckNotes.value.push(
            `Event date (${eventDate.toLocaleDateString('en-ZA')}) is after cover end (${coverEnd.toLocaleDateString('en-ZA')}) — policy check FAILED`
          )
        }
      }
    }

    // 2. Premiums up to date
    if (member) {
      // Check member's in_force status as proxy for premium status
      const memberInForce =
        member.in_force === true ||
        member.status?.toLowerCase() === 'active' ||
        member.status?.toLowerCase() === 'in_force'
      if (memberInForce) {
        formData.value.checklist.premiums_paid = true
        autoCheckNotes.value.push(
          'Member is active/in-force — premiums considered up to date'
        )
      } else {
        formData.value.checklist.premiums_paid = false
        autoCheckNotes.value.push(
          `Member status is "${member.status || 'unknown'}" — premiums check FAILED`
        )
      }
    } else if (scheme) {
      // Fallback: if scheme is in force, assume premiums OK
      const schemeInForce =
        scheme.in_force === true || scheme.status?.toLowerCase() === 'in_force'
      formData.value.checklist.premiums_paid = schemeInForce
      autoCheckNotes.value.push(
        schemeInForce
          ? 'Scheme is in force — premiums assumed up to date (verify member-level)'
          : 'Scheme not in force — premiums check FAILED'
      )
    }

    // 3. Waiting period satisfied
    if (props.claim.date_of_event && props.claim.benefit_code) {
      const eventDate = new Date(props.claim.date_of_event)
      const memberJoinDate = member?.date_joined
        ? new Date(member.date_joined)
        : member?.commencement_date
          ? new Date(member.commencement_date)
          : scheme?.cover_start_date
            ? new Date(scheme.cover_start_date)
            : null

      const requiredDays = waitingPeriodDays[props.claim.benefit_code] || 180

      if (memberJoinDate) {
        const daysSinceJoin = Math.floor(
          (eventDate.getTime() - memberJoinDate.getTime()) /
            (1000 * 60 * 60 * 24)
        )
        if (daysSinceJoin >= requiredDays) {
          formData.value.checklist.waiting_period_satisfied = true
          autoCheckNotes.value.push(
            `Waiting period satisfied: ${daysSinceJoin} days since join (${requiredDays} required for ${props.claim.benefit_code})`
          )
        } else {
          formData.value.checklist.waiting_period_satisfied = false
          autoCheckNotes.value.push(
            `Waiting period NOT met: ${daysSinceJoin} days since join, ${requiredDays} required for ${props.claim.benefit_code} — FAILED`
          )
        }
      } else {
        autoCheckNotes.value.push(
          'Could not determine member join date — waiting period requires manual check'
        )
      }
    }

    // 4. Benefit applicable for condition/event
    if (props.claim.benefit_code && props.claim.member_type) {
      const benefitCode = props.claim.benefit_code
      const memberType = props.claim.member_type?.toLowerCase()

      // Benefit-to-member-type applicability rules
      const applicabilityMap: Record<string, string[]> = {
        GLA: ['member'],
        SGLA: ['spouse'],
        GFF: ['member', 'spouse', 'child', 'parent', 'dependant'],
        PTD: ['member'],
        CI: ['member'],
        TTD: ['member'],
        PHI: ['member']
      }

      const allowedTypes = applicabilityMap[benefitCode] || []
      if (allowedTypes.includes(memberType)) {
        formData.value.checklist.benefit_applicable = true
        autoCheckNotes.value.push(
          `Benefit ${benefitCode} is applicable for ${memberType} claims`
        )
      } else {
        formData.value.checklist.benefit_applicable = false
        autoCheckNotes.value.push(
          `Benefit ${benefitCode} is NOT applicable for ${memberType} claims (allowed: ${allowedTypes.join(', ')}) — FAILED`
        )
      }
    }
  } catch (error) {
    console.error('Auto-check error:', error)
    autoCheckNotes.value.push(
      'An error occurred during auto-verification — manual checks required'
    )
  } finally {
    autoCheckLoading.value = false
    autoCheckCompleted.value = true
  }
}

onMounted(() => {
  runAutoChecks()
})

// Methods
const handleSubmit = async () => {
  if (!form.value) return

  const { valid } = await (form.value as any).validate()
  if (!valid) return

  loading.value = true
  try {
    const assessmentData = {
      ...formData.value,
      claim_id: props.claim?.id,
      assessment_timestamp: new Date().toISOString()
    }

    emit('save', assessmentData)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.v-card-title {
  padding: 12px 16px;
}

.v-card-text {
  padding: 16px;
}
</style>
