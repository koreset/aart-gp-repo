<template>
  <base-card>
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-percent-outline" class="mr-3" />
        <span>Discount Calculation Method</span>
        <v-spacer />
        <v-chip
          :color="discountMethod ? 'success' : 'warning'"
          variant="outlined"
          size="small"
        >
          {{ chipLabel }}
        </v-chip>
      </div>
    </template>
    <template #default>
      <p class="text-body-2 text-medium-emphasis mb-4">
        This setting controls how a discount is applied to every quote
        system-wide. It only takes effect on the next time a discount is applied
        to a quote — existing quotes are not retroactively recomputed.
      </p>
      <p v-if="lastUpdatedLabel" class="text-caption text-medium-emphasis mb-4">
        {{ lastUpdatedLabel }}
      </p>
      <v-radio-group
        v-model="discountMethod"
        :disabled="loading"
        density="compact"
      >
        <v-radio
          label="Loading Adjustment (default)"
          value="loading_adjustment"
        ></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          The discount is added to the scheme loading denominator, reducing only
          the loading-driven gross-up. Final office premium = risk / (1 −
          (scheme loading − discount)).
        </p>
        <v-radio label="Prorata" value="prorata"></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          The discount is applied as a flat (1 − d) factor to the entire
          pre-commission office premium, proportionally reducing every component
          (risk + each loading). Final premium per benefit = office premium × (1
          − d) + commission.
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

type DiscountMethod = 'loading_adjustment' | 'prorata'

const flash = useFlashStore()

const discountMethod = ref<DiscountMethod>('loading_adjustment')
const updatedAt = ref<string | null>(null)
const updatedBy = ref<string>('')
const loading = ref(false)
const saving = ref(false)

const chipLabel = computed(() =>
  discountMethod.value === 'prorata' ? 'Prorata' : 'Loading Adjustment'
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
    discountMethod.value =
      data?.discount_method === 'prorata' ? 'prorata' : 'loading_adjustment'
    updatedAt.value = data?.discount_method_updated_at ?? null
    updatedBy.value = data?.discount_method_updated_by ?? ''
  } catch (err: any) {
    flash.show(
      'Failed to load discount method: ' + (err?.message ?? 'unknown error'),
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
      discount_method: discountMethod.value
    })
    updatedAt.value = data?.discount_method_updated_at ?? updatedAt.value
    updatedBy.value = data?.discount_method_updated_by ?? updatedBy.value
    flash.show('Discount method saved', 'success')
  } catch (err: any) {
    flash.show(
      'Failed to save discount method: ' + (err?.message ?? 'unknown error'),
      'error'
    )
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
