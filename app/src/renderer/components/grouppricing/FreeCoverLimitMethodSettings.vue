<template>
  <base-card>
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-shield-check-outline" class="mr-3" />
        <span>Free Cover Limit Calculation Method</span>
        <v-spacer />
        <v-chip
          :color="fclMethod ? 'success' : 'warning'"
          variant="outlined"
          size="small"
        >
          {{ chipLabel }}
        </v-chip>
      </div>
    </template>
    <template #default>
      <p class="text-body-2 text-medium-emphasis mb-4">
        This setting controls how the Free Cover Limit is calculated for every
        quote system-wide. It only takes effect on the next time a quote is
        recomputed — existing quotes are not retroactively recomputed.
      </p>
      <p v-if="lastUpdatedLabel" class="text-caption text-medium-emphasis mb-4">
        {{ lastUpdatedLabel }}
      </p>
      <v-radio-group v-model="fclMethod" :disabled="loading" density="compact">
        <v-radio label="Percentile (default)" value="percentile"></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          Free cover limit is the lesser of the scaled mean salary (scaling × √N
          × mean salary) and the configured percentile of member sum assured.
        </p>
        <v-radio label="Statistical Outlier" value="outlier"></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          Free cover limit is the lesser of the scaled mean salary, the
          log-normal +3σ upper bound of sum assured (exp(μ + 3σ) of ln(SA)), and
          the largest member sum assured times the maximum-cover scaling factor.
        </p>
      </v-radio-group>
      <v-alert type="info" variant="tonal" density="compact" class="mb-4">
        <p class="text-body-2 mb-0">
          When a free cover limit is set on the quote, that value is used
          directly — except where it exceeds the maximum allowed free cover
          limit (configured per risk rate on the Restrictions table) by more
          than the override tolerance below, in which case the limit is clamped.
          This applies under both calculation methods.
        </p>
      </v-alert>
      <v-text-field
        v-model.number="tolerancePct"
        :disabled="loading"
        type="number"
        min="0"
        max="100"
        step="1"
        suffix="%"
        label="Override Tolerance"
        density="compact"
        variant="outlined"
        hide-details="auto"
        class="mb-4"
        :hint="`Allow quote-level overrides up to ${tolerancePct}% above the maximum allowed free cover limit before clamping.`"
        persistent-hint
      />
      <div class="d-flex justify-end">
        <v-btn
          color="primary"
          :loading="saving"
          :disabled="loading || !isValid"
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

type FCLMethod = 'percentile' | 'outlier'

const DEFAULT_TOLERANCE_PCT = 20

const flash = useFlashStore()

const fclMethod = ref<FCLMethod>('percentile')
const tolerancePct = ref<number>(DEFAULT_TOLERANCE_PCT)
const updatedAt = ref<string | null>(null)
const updatedBy = ref<string>('')
const loading = ref(false)
const saving = ref(false)

const isValid = computed(
  () =>
    Number.isFinite(tolerancePct.value) &&
    tolerancePct.value >= 0 &&
    tolerancePct.value <= 100
)

const chipLabel = computed(() =>
  fclMethod.value === 'outlier' ? 'Statistical Outlier' : 'Percentile'
)

const lastUpdatedLabel = computed(() => {
  if (!updatedAt.value) return ''
  const ts = new Date(updatedAt.value)
  if (Number.isNaN(ts.getTime())) return ''
  const formatted = ts.toLocaleString()
  return updatedBy.value
    ? `Method last updated ${formatted} by ${updatedBy.value}`
    : `Method last updated ${formatted}`
})

async function load() {
  loading.value = true
  try {
    const { data } = await GroupPricingService.getGroupPricingSettings()
    fclMethod.value = data?.fcl_method === 'outlier' ? 'outlier' : 'percentile'

    const fraction = Number(data?.fcl_override_tolerance)
    tolerancePct.value =
      Number.isFinite(fraction) && fraction > 0
        ? Math.round(fraction * 100)
        : DEFAULT_TOLERANCE_PCT

    updatedAt.value = data?.fcl_method_updated_at ?? null
    updatedBy.value = data?.fcl_method_updated_by ?? ''
  } catch (err: any) {
    flash.show(
      'Failed to load FCL method: ' + (err?.message ?? 'unknown error'),
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
      fcl_method: fclMethod.value,
      fcl_override_tolerance: tolerancePct.value / 100
    })
    updatedAt.value = data?.fcl_method_updated_at ?? updatedAt.value
    updatedBy.value = data?.fcl_method_updated_by ?? updatedBy.value
    flash.show('Free cover limit settings saved', 'success')
  } catch (err: any) {
    flash.show(
      'Failed to save FCL settings: ' + (err?.message ?? 'unknown error'),
      'error'
    )
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
