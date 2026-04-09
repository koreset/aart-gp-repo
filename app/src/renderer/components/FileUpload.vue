<template>
  <v-dialog
    :model-value="isDialogOpen"
    max-width="600px"
    @update:model-value="emit('update:isDialogOpen', $event)"
  >
    <base-card>
      <template #header>
        <span class="headline">Upload File</span>
      </template>
      <template #default>
        <v-form ref="form">
          <v-text-field v-model="fileName" label="File Name"></v-text-field>
          <v-btn size="small" @click="upload">Submit</v-btn>
        </v-form>
      </template>
    </base-card>
  </v-dialog>
</template>

<script setup lang="ts">
import BaseCard from './BaseCard.vue'
import { defineProps, watch, ref, defineEmits } from 'vue'

const props = defineProps<{
  isDialogOpen: boolean
}>()

const fileName = ref('')
const productCode = ref('')

const upload = () => {
  emit('upload', { fileName: fileName.value })
}

const emit = defineEmits<{
  (e: 'upload', payload: { fileName: string }): void
  (e: 'update:isDialogOpen', value: boolean): void
}>()

watch(
  () => props.isDialogOpen,
  (newVal) => {
    if (!newVal) {
      productCode.value = ''
    }
  }
)
</script>

<style scoped></style>
