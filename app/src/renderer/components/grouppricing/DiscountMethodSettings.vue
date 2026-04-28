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
          :disabled="!isDirty || loading"
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
import { useNotifications } from '@/renderer/composables/useNotifications'

type DiscountMethod = 'loading_adjustment' | 'prorata'

const { showSuccess, showError } = useNotifications()

const discountMethod = ref<DiscountMethod>('loading_adjustment')
const initialMethod = ref<DiscountMethod>('loading_adjustment')
const loading = ref(false)
const saving = ref(false)

const isDirty = computed(() => discountMethod.value !== initialMethod.value)

const chipLabel = computed(() =>
  discountMethod.value === 'prorata' ? 'Prorata' : 'Loading Adjustment'
)

async function load() {
  loading.value = true
  try {
    const { data } = await GroupPricingService.getGroupPricingSettings()
    const method: DiscountMethod =
      data?.discount_method === 'prorata' ? 'prorata' : 'loading_adjustment'
    discountMethod.value = method
    initialMethod.value = method
  } catch (err: any) {
    showError(
      'Failed to load discount method: ' + (err?.message ?? 'unknown error')
    )
  } finally {
    loading.value = false
  }
}

async function save() {
  saving.value = true
  try {
    await GroupPricingService.updateGroupPricingSettings({
      discount_method: discountMethod.value
    })
    initialMethod.value = discountMethod.value
    showSuccess('Discount method saved')
  } catch (err: any) {
    showError(
      'Failed to save discount method: ' + (err?.message ?? 'unknown error')
    )
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
