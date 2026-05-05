<template>
  <base-card>
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-calculator-variant" class="mr-3" />
        <span>PTD Base Rate Calculation Method</span>
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
        This setting controls how <code>BasePtdRate</code> is calculated for
        every quote system-wide. It only takes effect on the next quote
        calculation — existing quote outputs are not retroactively recomputed.
      </p>
      <p v-if="lastUpdatedLabel" class="text-caption text-medium-emphasis mb-4">
        {{ lastUpdatedLabel }}
      </p>
      <v-radio-group v-model="method" :disabled="loading" density="compact">
        <v-radio label="PTD only (default)" value="ptd_only"></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          <code>BasePtdRate = ptd_rate × (1 + PtdIndustryLoading + PtdRegionLoading)</code>.
          GLA AIDS rate is excluded — historical behaviour.
        </p>
        <v-radio label="PTD + GLA AIDS" value="ptd_plus_gla_aids"></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          Mirrors the GLA pattern by adding a separately-loaded GLA AIDS rate
          component:
          <code>BasePtdRate = ptd_rate × (1 + PtdIndustryLoading + PtdRegionLoading) + gla_aids_rate × (1 + GlaAidsRegionLoading)</code>.
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

type PtdBaseRateMethod = 'ptd_only' | 'ptd_plus_gla_aids'

const flash = useFlashStore()

const method = ref<PtdBaseRateMethod>('ptd_only')
const updatedAt = ref<string | null>(null)
const updatedBy = ref<string>('')
const loading = ref(false)
const saving = ref(false)

const chipLabel = computed(() =>
  method.value === 'ptd_plus_gla_aids' ? 'PTD + GLA AIDS' : 'PTD only'
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
      data?.ptd_base_rate_method === 'ptd_plus_gla_aids'
        ? 'ptd_plus_gla_aids'
        : 'ptd_only'
    updatedAt.value = data?.ptd_base_rate_method_updated_at ?? null
    updatedBy.value = data?.ptd_base_rate_method_updated_by ?? ''
  } catch (err: any) {
    flash.show(
      'Failed to load PTD base rate method: ' +
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
      ptd_base_rate_method: method.value
    })
    updatedAt.value =
      data?.ptd_base_rate_method_updated_at ?? updatedAt.value
    updatedBy.value =
      data?.ptd_base_rate_method_updated_by ?? updatedBy.value
    flash.show('PTD base rate method saved', 'success')
  } catch (err: any) {
    flash.show(
      'Failed to save PTD base rate method: ' +
        (err?.message ?? 'unknown error'),
      'error'
    )
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
