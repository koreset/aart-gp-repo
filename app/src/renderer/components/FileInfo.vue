<template>
  <v-dialog v-model="isDialogOpen" width="70%" height="80vh">
    <base-card :showActions="false">
      <template #header>
        <span class="headline">{{ tableTitle }}</span>
      </template>
      <template #default>
        <v-row>
          <v-col class="text-right">
            <v-btn size="small" variant="plain" rounded @click="upload"
              >Close</v-btn
            ></v-col
          >
        </v-row>
        <data-grid
          :showExport="showExport"
          :columnDefs="columnDefs"
          :rowData="rowData"
        />
      </template>
    </base-card>
  </v-dialog>
</template>

<script setup lang="ts">
import BaseCard from './BaseCard.vue'
import { ref, watch } from 'vue'
import DataGrid from './tables/DataGrid.vue'

interface Props {
  isDialogOpen: boolean
  rowData: any
  columnDefs: any
  tableTitle: string | null
  showExport?: boolean
}

const emit = defineEmits<{
  (e: 'update:isInfoDialogOpen', value: boolean): void
}>()

const props = defineProps<Props>()

const isDialogOpen = ref(props.isDialogOpen)

watch(
  () => props.isDialogOpen,
  (value) => {
    isDialogOpen.value = value
  }
)

watch(isDialogOpen, (newVal) => {
  emit('update:isInfoDialogOpen', newVal)
})

const upload = () => {
  isDialogOpen.value = false
}
</script>

<style scoped></style>
