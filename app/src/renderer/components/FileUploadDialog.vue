<template>
  <v-dialog v-model="localIsDialogOpen" max-width="600px">
    <base-card>
      <template #header>
        <span class="headline">{{ uploadTitle }}</span>
      </template>
      <template #default>
        <v-form ref="form">
          <v-file-input
            v-model="file"
            variant="outlined"
            density="compact"
            label="File"
            prepend-icon="mdi-paperclip"
            accept=".csv"
            :rules="[rules.required, rules.csvOnly]"
            required
          ></v-file-input>
          <v-text-field
            v-if="showModelPoint"
            v-model="fileName"
            density="compact"
            variant="outlined"
            :label="mpLabel"
            :rules="[rules.required]"
            required
          ></v-text-field>

          <v-select
            v-if="showForTables()"
            v-model="year"
            variant="outlined"
            density="compact"
            :label="yearLabel"
            :items="years"
            :rules="[rules.required]"
            required
          ></v-select>
        </v-form>
      </template>
      <template #actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" @click="upload">Upload</v-btn>
        <v-btn @click="close">Cancel</v-btn>
      </template>
    </base-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseCard from './BaseCard.vue'

interface Props {
  isDialogOpen: boolean
  yearLabel: string
  table: string | undefined
  uploadTitle: string
  mpLabel: string
  showModelPoint: boolean
}

const props = withDefaults(defineProps<Props>(), {
  yearLabel: 'Select a year',
  uploadTitle: 'Upload File Data',
  table: undefined
})

const showForTables = () => {
  return (
    props.table === 'Accidental_Mortality' ||
    props.table === 'Mortality' ||
    props.table === 'Lapse' ||
    props.table === 'Lapse_Margins' ||
    props.table === 'Model Points' ||
    props.table === 'Disability_Incidence' ||
    props.table === 'Critical_Illness' ||
    props.table === 'Retrenchment'
  )
}

// const props = defineProps<{
//   showProductCode: boolean;
//   isDialogOpen: boolean;
//   yearLabel: {
//     type: string,
//     default: 'Select a year',
//   }
// }>();

const years = ref<number[]>(
  Array.from({ length: 10 }, (v, k) => new Date().getFullYear() - k)
)

const emit = defineEmits<{
  (
    e: 'upload',
    payload: {
      file: File | null
      fileName: string
      productCode: string
      year: number | null
    }
  ): void
  (e: 'update:isDialogOpen', value: boolean): void
}>()

const localIsDialogOpen = ref(props.isDialogOpen)

const file = ref<File | null>(null)
const fileName = ref('')
const productCode = ref('')
const year = ref<number | null>(null)

const rules = {
  required: (value: any) => !!value || 'Required.',
  csvOnly: (value: File | null) => {
    if (!value) return true
    return (
      value.name.toLowerCase().endsWith('.csv') ||
      'Only CSV files are accepted.'
    )
  }
}

const upload = async () => {
  const { valid } = await form.value.validate()
  if (!valid) return
  emit('upload', {
    file: file.value,
    fileName: fileName.value,
    productCode: productCode.value,
    year: year.value
  })
  close()
}

const close = () => {
  localIsDialogOpen.value = false
}

const form = ref()

watch(
  () => props.isDialogOpen,
  (newVal) => {
    localIsDialogOpen.value = newVal
  }
)

watch(localIsDialogOpen, (newVal) => {
  emit('update:isDialogOpen', newVal)
})

watch(
  () => props.isDialogOpen,
  (newVal) => {
    if (!newVal) {
      file.value = null
      fileName.value = ''
      productCode.value = ''
      year.value = null
    }
  }
)
</script>

<style scoped>
/* Add any scoped styles here */
</style>
