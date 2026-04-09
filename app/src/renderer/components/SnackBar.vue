<template>
  <v-snackbar v-model="snackbar" :timeout="timeout">
    {{ text }}

    <template #actions>
      <v-btn color="white" variant="text" @click="snackbar = false">
        Close
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps({
  text: {
    type: String,
    required: true
  },
  timeout: {
    type: Number,
    default: 6000
  },
  snackbar: {
    type: Boolean,
    default: false
  }
})

const snackbar = ref(props.snackbar)
</script>

<style scoped></style>
