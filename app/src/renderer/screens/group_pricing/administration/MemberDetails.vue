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
                  Member Details
                  <template v-if="selectedMember?.member_name">
                    - {{ selectedMember.member_name }}
                  </template>
                </span>
              </div>
              <div v-if="selectedMember">
                <v-btn
                  size="small"
                  rounded
                  color="white"
                  variant="outlined"
                  class="mr-2"
                  @click="editMember"
                >
                  Edit Member
                </v-btn>
                <v-btn
                  rounded
                  size="small"
                  color="error"
                  variant="outlined"
                  @click="deactivateMember"
                >
                  Deactivate
                </v-btn>
              </div>
            </div>
          </template>
          <template #default>
            <v-row v-if="loading && !selectedMember">
              <v-col cols="12" class="text-center py-8">
                <v-progress-circular indeterminate color="primary" />
                <div class="mt-2 text-body-2 text-medium-emphasis"
                  >Loading member...</div
                >
              </v-col>
            </v-row>

            <v-row v-else-if="!selectedMember">
              <v-col cols="12" class="text-center py-8">
                <v-icon size="64" color="grey-lighten-1" class="mb-4"
                  >mdi-account-off</v-icon
                >
                <h3 class="text-h6 text-grey-darken-1 mb-2">Member not found</h3>
                <v-btn color="primary" variant="outlined" @click="goBack">
                  Back to Member Management
                </v-btn>
              </v-col>
            </v-row>

            <member-detail-view
              v-else
              :member="selectedMember"
              :beneficiaries="memberBeneficiaries"
              :schemes="schemes"
              @manage-beneficiaries="openBeneficiaryManagement"
              @view-claims="viewMemberClaims"
              @member-updated="handleMemberUpdated"
              @claim-registered="handleClaimRegistered"
              @notify="handleNotify"
            />
          </template>
        </base-card>
      </v-col>
    </v-row>

    <!-- Edit Member Dialog -->
    <v-dialog v-model="editMemberDialog" persistent max-width="800px">
      <base-card>
        <template #header>
          <span class="headline">Edit Member</span>
        </template>
        <template #default>
          <member-enrollment-form
            :member="selectedMember"
            :schemes="schemes"
            :is-edit-mode="true"
            :preselected-scheme-id="selectedMember?.scheme_id ?? null"
            @save="handleMemberSave"
            @cancel="editMemberDialog = false"
          />
        </template>
      </base-card>
    </v-dialog>

    <!-- Exit Date Confirmation Dialog -->
    <v-dialog v-model="exitDateDialog" persistent max-width="500px">
      <base-card>
        <template #header>
          <span class="headline">Confirm Member Deactivation</span>
        </template>
        <template #default>
          <div class="pa-4">
            <div class="mb-4">
              <p class="text-body-1 mb-2">
                You are about to deactivate
                <strong>{{ selectedMember?.member_name }}</strong> from the
                <strong>{{ selectedMember?.scheme_name }}</strong> scheme.
              </p>
              <p class="text-body-2 text-medium-emphasis">
                Please specify the effective exit date for this member.
              </p>
            </div>

            <v-form
              v-model="exitDateValid"
              @submit.prevent="confirmDeactivation"
            >
              <v-text-field
                v-model="exitDate"
                label="Effective Exit Date"
                type="date"
                variant="outlined"
                :rules="exitDateRules"
                :min="
                  selectedMember?.entry_date
                    ? new Date(selectedMember.entry_date)
                        .toISOString()
                        .substr(0, 10)
                    : undefined
                "
                required
                class="mb-4"
              />

              <div class="text-caption text-medium-emphasis mb-4">
                <v-icon size="small" class="mr-1">mdi-information</v-icon>
                Entry Date:
                {{
                  selectedMember?.entry_date
                    ? new Date(selectedMember.entry_date).toLocaleDateString()
                    : 'Not specified'
                }}
              </div>
            </v-form>
          </div>
        </template>
        <template #actions>
          <v-btn color="grey" @click="cancelDeactivation">Cancel</v-btn>
          <v-btn
            color="error"
            :disabled="!exitDateValid"
            @click="confirmDeactivation"
          >
            Deactivate Member
          </v-btn>
        </template>
      </base-card>
    </v-dialog>

    <!-- Snackbar for notifications -->
    <v-snackbar v-model="snackbar" :timeout="4000" :color="snackbarColor">
      {{ snackbarMessage }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '@/renderer/components/BaseCard.vue'
import MemberEnrollmentForm from './components/MemberEnrollmentForm.vue'
import MemberDetailView from './components/MemberDetailView.vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

interface Member {
  id?: number
  member_name: string
  member_id_number: string
  member_id_type: string
  scheme_name: string
  gender: string
  date_of_birth: Date | null
  email?: string
  phone_number?: string
  employee_number?: string
  scheme_id: number | null
  scheme_category: string
  entry_date: Date | null
  annual_salary: number
  status?: string
  effective_exit_date?: string | Date | null
  occupation?: string
  occupational_class?: string
  address_line_1?: string
  address_line_2?: string
  city?: string
  province?: string
  postal_code?: string
  benefits: {
    gla_enabled: boolean
    gla_multiple?: number
    sgla_enabled: boolean
    sgla_multiple?: number
    ptd_enabled: boolean
    ptd_multiple?: number
    ci_enabled: boolean
    ci_multiple?: number
    ttd_enabled: boolean
    ttd_multiple?: number
    phi_enabled: boolean
    phi_multiple?: number
    gff_enabled: boolean
  }
}

interface Scheme {
  id: number
  name: string
}

const props = defineProps<{
  id: string | number
}>()

const router = useRouter()

const loading = ref(false)
const selectedMember = ref<Member | null>(null)
const memberBeneficiaries = ref<any[]>([])
const schemes = ref<Scheme[]>([])

const editMemberDialog = ref(false)
const exitDateDialog = ref(false)

const exitDate = ref('')
const exitDateValid = ref(false)
const exitDateRules = [
  (v: string) => !!v || 'Exit date is required',
  (v: string) => {
    if (!v) return true
    const selectedDate = new Date(v)
    const today = new Date()
    const thirtyDaysAgo = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000)
    const entryDate = selectedMember.value?.entry_date
      ? new Date(selectedMember.value.entry_date)
      : null

    if (selectedDate < thirtyDaysAgo) {
      return 'Exit date cannot be more than 30 days in the past'
    }
    if (entryDate && selectedDate <= entryDate) {
      return 'Exit date must be after entry date'
    }
    return true
  }
]

const snackbar = ref(false)
const snackbarMessage = ref('')
const snackbarColor = ref('success')

const showSnackbar = (message: string, color: string = 'success') => {
  snackbarMessage.value = message
  snackbarColor.value = color
  snackbar.value = true
}

const memberId = (): number => {
  return typeof props.id === 'string' ? parseInt(props.id, 10) : props.id
}

const loadSchemes = async () => {
  try {
    const response = await GroupPricingService.getSchemesInforce()
    schemes.value = response.data
  } catch (error) {
    console.error('Error loading schemes:', error)
  }
}

const loadMember = async () => {
  const id = memberId()
  if (!id || Number.isNaN(id)) {
    showSnackbar('Invalid member id', 'error')
    return
  }

  loading.value = true
  try {
    const memberResponse = await GroupPricingService.getMemberInfo(id)
    const member: Member = memberResponse.data
    if (member && !member.scheme_name) {
      const scheme = schemes.value.find((s) => s.id === member.scheme_id)
      member.scheme_name = scheme?.name || 'Unknown Scheme'
    }
    selectedMember.value = member
    await loadMemberBeneficiaries(id)
  } catch (error) {
    console.error('Error loading member:', error)
    showSnackbar('Failed to load member', 'error')
  } finally {
    loading.value = false
  }
}

const loadMemberBeneficiaries = async (id: number) => {
  try {
    const response = await GroupPricingService.getMemberBeneficiaries(id)
    memberBeneficiaries.value = response.data
  } catch (error) {
    console.error('Error loading beneficiaries:', error)
  }
}

const goBack = () => {
  router.push({ name: 'group-pricing-member-management' })
}

const editMember = () => {
  editMemberDialog.value = true
}

const handleMemberSave = async (memberData: any) => {
  try {
    await GroupPricingService.updateMember(selectedMember.value?.id, memberData)
    showSnackbar('Member updated successfully', 'success')
    editMemberDialog.value = false
    await loadMember()
  } catch (error: any) {
    console.error('Error saving member:', error)
    showSnackbar(
      error?.response?.data || error?.message || 'Failed to save member',
      'error'
    )
  }
}

const deactivateMember = () => {
  if (!selectedMember.value) return
  exitDate.value = new Date().toISOString().substr(0, 10)
  exitDateValid.value = false
  exitDateDialog.value = true
}

const confirmDeactivation = async () => {
  if (!selectedMember.value || !exitDateValid.value) return

  try {
    const today = new Date().toISOString().slice(0, 10)
    const exitInFuture = exitDate.value > today

    const memberToUpdate: any = {
      ...selectedMember.value,
      effective_exit_date: exitDate.value,
      status: exitInFuture ? 'ACTIVE' : 'INACTIVE'
    }
    await GroupPricingService.removeMemberFromScheme(
      selectedMember.value.scheme_id,
      memberToUpdate
    )
    showSnackbar(
      exitInFuture
        ? `Deactivation scheduled for ${exitDate.value}`
        : 'Member deactivated successfully',
      'success'
    )
    exitDateDialog.value = false
    await loadMember()
  } catch (error) {
    console.error('Error deactivating member:', error)
    showSnackbar('Error deactivating member', 'error')
  }
}

const cancelDeactivation = () => {
  exitDateDialog.value = false
  exitDate.value = ''
  exitDateValid.value = false
}

const openBeneficiaryManagement = () => {
  if (!selectedMember.value?.id) return
  router.push({
    name: 'group-pricing-beneficiary-management',
    params: { memberId: selectedMember.value.id }
  })
}

const viewMemberClaims = () => {
  if (!selectedMember.value?.id) return
  router.push({
    name: 'group-pricing-claims-management',
    query: { memberId: selectedMember.value.id }
  })
}

const handleMemberUpdated = async (updatedMember: any) => {
  selectedMember.value = updatedMember
  await loadMember()
  showSnackbar('Member updated successfully', 'success')
}

const handleClaimRegistered = async () => {
  // Member-level data is unchanged by claim registration; nothing to refresh
  // here today. Hook left in place so future per-member claim views can react.
}

const handleNotify = (message: string, color: string = 'success') => {
  showSnackbar(message, color)
}

onMounted(async () => {
  await loadSchemes()
  await loadMember()
})
</script>

<style scoped></style>
