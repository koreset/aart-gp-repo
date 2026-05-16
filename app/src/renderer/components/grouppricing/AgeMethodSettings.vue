<template>
  <base-card>
    <template #header>
      <div class="d-flex align-center">
        <v-icon icon="mdi-calendar-account-outline" class="mr-3" />
        <span>Age Calculation Method</span>
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
        This setting controls how each member's age is derived from their date
        of birth and the scheme commencement date for every quote system-wide.
        It only takes effect on the next quote calculation — existing quote
        outputs are not retroactively recomputed.
      </p>
      <p v-if="lastUpdatedLabel" class="text-caption text-medium-emphasis mb-4">
        {{ lastUpdatedLabel }}
      </p>
      <v-radio-group v-model="method" :disabled="loading" density="compact">
        <v-radio
          label="Age Next Birthday (default)"
          value="age_next_birthday"
        ></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          Member is one year older on/after their birthday this year; before it,
          they are at the current year-difference age. Historical behaviour.
        </p>
        <v-radio label="Age Last Birthday" value="age_last_birthday"></v-radio>
        <p class="text-body-2 text-medium-emphasis ml-8 mb-4">
          Floored months-between-dates:
          <code
            >ROUNDDOWN((12*(YEAR(CommenDate)-YEAR(DoB)) +
            (MONTH(CommenDate)-MONTH(DoB)))/12, 0)</code
          >.
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

type AgeMethod = 'age_next_birthday' | 'age_last_birthday'

const flash = useFlashStore()

const method = ref<AgeMethod>('age_next_birthday')
const updatedAt = ref<string | null>(null)
const updatedBy = ref<string>('')
const loading = ref(false)
const saving = ref(false)

const chipLabel = computed(() =>
  method.value === 'age_last_birthday'
    ? 'Age Last Birthday'
    : 'Age Next Birthday'
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
      data?.age_method === 'age_last_birthday'
        ? 'age_last_birthday'
        : 'age_next_birthday'
    updatedAt.value = data?.age_method_updated_at ?? null
    updatedBy.value = data?.age_method_updated_by ?? ''
  } catch (err: any) {
    flash.show(
      'Failed to load age calculation method: ' +
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
      age_method: method.value
    })
    updatedAt.value = data?.age_method_updated_at ?? updatedAt.value
    updatedBy.value = data?.age_method_updated_by ?? updatedBy.value
    flash.show('Age calculation method saved', 'success')
  } catch (err: any) {
    flash.show(
      'Failed to save age calculation method: ' +
        (err?.message ?? 'unknown error'),
      'error'
    )
  } finally {
    saving.value = false
  }
}

onMounted(load)
</script>
