<template>
  <v-card variant="outlined" rounded="lg">
    <v-card-title class="d-flex align-center font-weight-bold">
      <v-icon class="mr-2">mdi-briefcase-check-outline</v-icon>
      Actively at work attestation
    </v-card-title>
    <v-card-text>
      <p class="text-caption mb-3"
        >This attestation confirms the member is actively performing the duties
        of their occupation on the cover date. Required before cover commences
        for any tier-1+ member.</p
      >

      <v-row dense>
        <v-col cols="12" md="6">
          <v-text-field
            v-model="form.attested_by_name"
            label="Signed by (typed name)"
            density="compact"
            variant="outlined"
          />
        </v-col>
        <v-col cols="12" md="6">
          <v-select
            v-model="form.attested_by_role"
            :items="roleOptions"
            label="Role"
            density="compact"
            variant="outlined"
          />
        </v-col>
      </v-row>

      <v-checkbox
        v-model="confirmed"
        label="I confirm the member is actively at work on the cover date and not absent due to illness, injury or leave."
        density="compact"
        hide-details
      />

      <v-alert
        v-if="error"
        type="error"
        variant="tonal"
        density="compact"
        class="mt-3"
        >{{ error }}</v-alert
      >
      <v-btn
        class="mt-3"
        color="primary"
        :loading="busy"
        :disabled="!canSubmit"
        @click="submit"
        >Sign &amp; submit</v-btn
      >

      <v-divider v-if="attestations.length" class="my-4" />
      <p
        v-if="attestations.length"
        class="text-subtitle-2 font-weight-bold mb-2"
        >Signed attestations</p
      >
      <v-list density="compact">
        <v-list-item
          v-for="a in attestations"
          :key="a.id"
          :title="a.attested_by_name"
          :subtitle="`${a.attested_by_role || 'self'} · ${formatDate(a.attested_at)}`"
        >
          <template #append>
            <v-tooltip location="top">
              <template #activator="{ props: tipProps }">
                <v-icon v-bind="tipProps" size="small" color="success"
                  >mdi-shield-check-outline</v-icon
                >
              </template>
              <span
                >Signature {{ a.signature_hash?.slice(0, 12) }}… ·
                {{ a.ip_address }}</span
              >
            </v-tooltip>
          </template>
        </v-list-item>
      </v-list>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'

const props = defineProps<{
  caseId: number
  quoteId: number
  memberName?: string
  memberIdNumber?: string
}>()
const emit = defineEmits<{
  (e: 'submitted'): void
}>()

interface Attestation {
  id: number
  attested_by_name: string
  attested_by_role: string
  attested_at: string
  signature_hash: string
  ip_address: string
}

const attestations = ref<Attestation[]>([])
const confirmed = ref(false)
const busy = ref(false)
const error = ref('')

const form = ref({
  attested_by_name: '',
  attested_by_role: 'self'
})

const roleOptions = [
  { title: 'Self (the member)', value: 'self' },
  { title: 'Employer / HR', value: 'employer_hr' },
  { title: 'Broker', value: 'broker' }
]

const canSubmit = computed(
  () =>
    confirmed.value &&
    form.value.attested_by_name.trim().length > 0 &&
    form.value.attested_by_role.length > 0
)

const formatDate = (s: string) => (s ? new Date(s).toLocaleString() : '—')

const load = async () => {
  try {
    const res = await GroupPricingService.listActivelyAtWork({
      case_id: props.caseId
    })
    attestations.value = res.data || []
  } catch (err) {
    console.warn('Failed to load attestations', err)
  }
}

const submit = async () => {
  error.value = ''
  busy.value = true
  try {
    await GroupPricingService.submitActivelyAtWork({
      case_id: props.caseId,
      quote_id: props.quoteId,
      member_name: props.memberName,
      member_id_number: props.memberIdNumber,
      attested_by_name: form.value.attested_by_name.trim(),
      attested_by_role: form.value.attested_by_role
    })
    form.value.attested_by_name = ''
    confirmed.value = false
    emit('submitted')
    await load()
  } catch (err: any) {
    error.value =
      err?.response?.data || err?.message || 'Failed to submit attestation'
  } finally {
    busy.value = false
  }
}

watch(
  () => props.caseId,
  () => load(),
  { immediate: false }
)

onMounted(load)
</script>
