<template>
  <div class="ag-status-name-value">
    <span class="component-label">Visible Rows (incl. headers): </span>
    <span class="component-value">{{ count }}</span>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps(['params'])
const count = ref(0)

const updateStatusBar = () => {
  if (props.params.isInfinite) {
    count.value = props.params.api.getDisplayedRowCount() - 1
  } else {
    count.value = props.params.api.getDisplayedRowCount()
  }
}

onMounted(() => {
  // Listen to grid events to update the status dynamically
  props.params.api.addEventListener('modelUpdated', updateStatusBar)
  updateStatusBar()
})

onUnmounted(() => {
  // Clean up the event listener when the component is destroyed
  props.params.api.removeEventListener('modelUpdated', updateStatusBar)
})
</script>

<style scoped>
.component-label {
  font-weight: bold;
}
.ag-status-name-value {
  margin: 0 10px;
}
</style>
