<template>
  <base-card>
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-alert-octagon-outline" class="mr-3" />
        <span>Risk Watchlist Thresholds</span>
        <v-spacer />
        <v-chip color="primary" variant="outlined" size="small">
          ALR &gt; {{ alrCeiling }}% · Δ &gt; {{ alrDelta }}pp
        </v-chip>
      </div>
    </template>
    <template #default>
      <p class="text-body-2 text-medium-emphasis mb-4">
        These two thresholds drive the <strong>Deteriorating Schemes</strong>
        watchlist on the In-Force Performance &amp; Risk dashboard. They
        represent the company's stated risk appetite — analysts can still dial
        them down on the dashboard for ad-hoc what-if exploration, but the
        values saved here are the official defaults the dashboard always loads
        with.
      </p>
      <p v-if="lastUpdatedLabel" class="text-caption text-medium-emphasis mb-4">
        {{ lastUpdatedLabel }}
      </p>

      <v-row dense>
        <v-col cols="12" md="6">
          <v-text-field
            v-model.number="alrCeiling"
            label="ALR ceiling (%)"
            type="number"
            min="0"
            max="1000"
            step="1"
            variant="outlined"
            density="compact"
            hint="Flag schemes whose ITD Actual Loss Ratio exceeds this percentage."
            persistent-hint
            :disabled="loading"
          />
        </v-col>
        <v-col cols="12" md="6">
          <v-text-field
            v-model.number="alrDelta"
            label="ALR − ELR delta (pp)"
            type="number"
            min="0"
            max="1000"
            step="1"
            variant="outlined"
            density="compact"
            hint="Flag schemes whose ITD ALR exceeds the Expected Loss Ratio by more than this many percentage points."
            persistent-hint
            :disabled="loading"
          />
        </v-col>
        <v-col cols="12" md="6">
          <v-text-field
            v-model.number="profileVariationTolerance"
            label="Quote Risk Profile Variation tolerance (%)"
            type="number"
            min="0"
            max="100"
            step="0.5"
            variant="outlined"
            density="compact"
            hint="Printed on the Acceptance Form: if the member data profile at implementation differs by more than this %, the insurer reserves the right to revise rates."
            persistent-hint
            :disabled="loading"
          />
        </v-col>
      </v-row>

      <v-alert
        v-if="rangeWarning"
        type="warning"
        variant="tonal"
        density="compact"
        class="mt-3"
      >
        {{ rangeWarning }}
      </v-alert>

      <div class="d-flex justify-end mt-4">
        <v-btn
          color="primary"
          :loading="saving"
          :disabled="loading || !canSave"
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

const flash = useFlashStore()

const alrCeiling = ref<number>(100)
const alrDelta = ref<number>(20)
const profileVariationTolerance = ref<number>(7)
const updatedAt = ref<string | null>(null)
const updatedBy = ref<string>('')
const loading = ref(false)
const saving = ref(false)

const canSave = computed(
  () =>
    Number.isFinite(alrCeiling.value) &&
    Number.isFinite(alrDelta.value) &&
    Number.isFinite(profileVariationTolerance.value) &&
    alrCeiling.value >= 0 &&
    alrDelta.value >= 0 &&
    profileVariationTolerance.value >= 0 &&
    profileVariationTolerance.value <= 100
)

const rangeWarning = computed(() => {
  const issues: string[] = []
  if (alrCeiling.value < 50 || alrCeiling.value > 300) {
    issues.push('ALR ceiling outside the typical 50–300% range')
  }
  if (alrDelta.value < 5 || alrDelta.value > 100) {
    issues.push('Delta outside the typical 5–100pp range')
  }
  if (
    profileVariationTolerance.value < 1 ||
    profileVariationTolerance.value > 25
  ) {
    issues.push(
      'Profile variation tolerance outside the typical 1–25% range'
    )
  }
  return issues.length > 0
    ? `${issues.join('; ')}. Save anyway if intentional.`
    : ''
})

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
    alrCeiling.value = Number(data?.risk_alr_ceiling_pct ?? 100)
    alrDelta.value = Number(data?.risk_alr_delta_pp ?? 20)
    profileVariationTolerance.value = Number(
      data?.risk_profile_variation_tolerance_pct ?? 7
    )
    updatedAt.value = data?.risk_thresholds_updated_at ?? null
    updatedBy.value = data?.risk_thresholds_updated_by ?? ''
  } catch (err: any) {
    flash.show(
      'Failed to load risk thresholds: ' + (err?.message ?? 'unknown error'),
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
      risk_alr_ceiling_pct: alrCeiling.value,
      risk_alr_delta_pp: alrDelta.value,
      risk_profile_variation_tolerance_pct: profileVariationTolerance.value
    })
    updatedAt.value = data?.risk_thresholds_updated_at ?? updatedAt.value
    updatedBy.value = data?.risk_thresholds_updated_by ?? updatedBy.value
    flash.show('Risk thresholds saved', 'success')
  } catch (err: any) {
    flash.show(
      'Failed to save risk thresholds: ' + (err?.message ?? 'unknown error'),
      'error'
    )
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
