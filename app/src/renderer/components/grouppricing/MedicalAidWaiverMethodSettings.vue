<template>
  <base-card>
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-medical-bag" class="mr-3" />
        <span>Medical Aid Waiver Calculation Method</span>
        <v-spacer />
        <v-chip
          :color="method ? 'success' : 'warning'"
          variant="outlined"
          size="small"
        >
          {{ chipLabel }}
        </v-chip>
      </div>
    </template>
    <template #default>
      <p class="text-body-2 text-medium-emphasis mb-4">
        This setting controls how the PHI medical aid waiver
        (<code>phi_medical_aid_waiver</code>) is calculated for every quote
        system-wide. It only takes effect on the next quote calculation —
        existing quote outputs are not retroactively recomputed.
      </p>
      <p v-if="lastUpdatedLabel" class="text-caption text-medium-emphasis mb-4">
        {{ lastUpdatedLabel }}
      </p>
      <v-radio-group v-model="method" :disabled="loading" density="compact">
        <v-radio label="Formula (default)" value="formula"></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          Annual salary × medical aid waiver proportion + medical aid waiver
          amount, capped at the maximum medical aid waiver restriction.
        </p>
        <v-radio label="Table Lookup" value="table_lookup"></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          Reads <code>medicalwaiver_sum_at_risk</code> from the
          <code>medical_waivers</code> table by risk rate code, gender, age next
          birthday, and income level. Members with no matching row receive a
          zero waiver.
        </p>
      </v-radio-group>
      <div class="d-flex justify-end">
        <v-btn
          color="primary"
          :loading="saving"
          :disabled="loading"
          @click="save"
        >
          Save
        </v-btn>
      </div>
    </template>
  </base-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import BaseCard from '@/renderer/components/BaseCard.vue'
import { useFlashStore } from '@/renderer/store/flash'

type MedicalAidWaiverMethod = 'formula' | 'table_lookup'

const flash = useFlashStore()

const method = ref<MedicalAidWaiverMethod>('formula')
const updatedAt = ref<string | null>(null)
const updatedBy = ref<string>('')
const loading = ref(false)
const saving = ref(false)

const chipLabel = computed(() =>
  method.value === 'table_lookup' ? 'Table Lookup' : 'Formula'
)

const lastUpdatedLabel = computed(() => {
  if (!updatedAt.value) return ''
  const ts = new Date(updatedAt.value)
  if (Number.isNaN(ts.getTime())) return ''
  const formatted = ts.toLocaleString()
  return updatedBy.value
    ? `Last updated ${formatted} by ${updatedBy.value}`
    : `Last updated ${formatted}`
})

async function load() {
  loading.value = true
  try {
    const { data } = await GroupPricingService.getGroupPricingSettings()
    method.value =
      data?.medical_aid_waiver_method === 'table_lookup'
        ? 'table_lookup'
        : 'formula'
    updatedAt.value = data?.medical_aid_waiver_method_updated_at ?? null
    updatedBy.value = data?.medical_aid_waiver_method_updated_by ?? ''
  } catch (err: any) {
    flash.show(
      'Failed to load medical aid waiver method: ' +
        (err?.message ?? 'unknown error'),
      'error'
    )
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    const { data } = await GroupPricingService.updateGroupPricingSettings({
      medical_aid_waiver_method: method.value
    })
    updatedAt.value =
      data?.medical_aid_waiver_method_updated_at ?? updatedAt.value
    updatedBy.value =
      data?.medical_aid_waiver_method_updated_by ?? updatedBy.value
    flash.show('Medical aid waiver method saved', 'success')
  } catch (err: any) {
    flash.show(
      'Failed to save medical aid waiver method: ' +
        (err?.message ?? 'unknown error'),
      'error'
    )
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
