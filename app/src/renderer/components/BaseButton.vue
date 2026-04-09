<template>
  <v-btn
    :size="size"
    :color="color"
    :icon="icon"
    :disabled="disabled"
    @click="handleClick"
  >
    <template v-if="icon">
      <v-icon>{{ icon }}</v-icon>
    </template>
    <template v-else>
      {{ label }}
    </template>
  </v-btn>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

// eslint-disable-next-line no-unused-vars
const props = defineProps<{
  label: string
  color?: string
  icon?: string
  disabled?: boolean
  size?: string
}>()

const emit = defineEmits<{
  (e: 'click'): void
}>()

const handleClick = () => {
  emit('click')
}
</script>

<style scoped>
/* Add any scoped styles here if needed */
</style>
