<template>
  <v-row v-if="props.loadingData">
    <v-col>
      <p class="mt-3">{{ props.label || 'Loading data...' }}</p>
      <v-progress-linear
        buffer-value="20"
        color="primary"
        stream
        value="10"
      ></v-progress-linear>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
const props = defineProps({
  loadingData: {
    type: Boolean,
    default: false
  },
  label: {
    type: String,
    default: ''
  }
})
</script>

<style></style>
