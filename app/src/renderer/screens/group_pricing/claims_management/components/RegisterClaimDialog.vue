<template>
  <v-dialog
    :model-value="modelValue"
    persistent
    max-width="900px"
    @update:model-value="(v) => emit('update:modelValue', v)"
  >
    <base-card>
      <template #header>
        <span class="headline">Register New Claim</span>
      </template>
      <template #default>
        <claim-registration-form
          :schemes="schemes"
          :prefilled-member="prefilledMember"
          @save="onSave"
          @cancel="onCancel"
        />
      </template>
    </base-card>
  </v-dialog>
</template>

<script setup lang="ts">
import BaseCard from '@/renderer/components/BaseCard.vue'
import ClaimRegistrationForm from './ClaimRegistrationForm.vue'

interface Scheme {
  id: number
  name: string
  [key: string]: any
}

interface PrefilledMember {
  member_id_number: string
  member_name?: string
  [key: string]: any
}

interface Props {
  modelValue: boolean
  schemes: Scheme[]
  prefilledMember?: PrefilledMember | null
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'save', payload: FormData): void
  (e: 'cancel'): void
}

withDefaults(defineProps<Props>(), {
  prefilledMember: null
})

const emit = defineEmits<Emits>()

const onSave = (payload: FormData) => {
  emit('save', payload)
}

const onCancel = () => {
  emit('cancel')
  emit('update:modelValue', false)
}
</script>
