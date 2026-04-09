<template>
  <v-chip
    :size="small ? 'x-small' : 'small'"
    variant="tonal"
    :style="chipStyle"
  >
    {{ label }}
  </v-chip>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { STATUS_COLORS } from '@/renderer/constants/designTokens'

const props = defineProps<{
  status: string
  small?: boolean
}>()

const token = computed(() => {
  const key = (props.status ?? '').toLowerCase().replace(/\s+/g, '_')
  return STATUS_COLORS[key] ?? { hex: '#9E9E9E', label: props.status || '—' }
})

const label = computed(() => token.value.label)

const chipStyle = computed(() => ({
  background: `${token.value.hex}22`,
  color: token.value.hex,
  border: `1px solid ${token.value.hex}55`,
  fontWeight: '600'
}))
</script>
