<template>
  <v-tooltip :text="tooltipText" location="top">
    <template #activator="{ props: tooltipProps }">
      <v-chip
        v-bind="tooltipProps"
        :color="chipColor"
        :size="size"
        variant="flat"
        label
        class="probability-chip"
      >
        <span v-if="loading">–</span>
        <span
          v-else-if="score === null || score === undefined"
          class="text-disabled"
          >–</span
        >
        <span v-else>{{ score.toFixed(1) }}%</span>
      </v-chip>
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  score: number | null | undefined
  size?: 'x-small' | 'small' | 'default' | 'large'
  loading?: boolean
}>()

const chipColor = computed(() => {
  if (props.loading || props.score === null || props.score === undefined)
    return 'grey-lighten-1'
  if (props.score >= 75) return 'success'
  if (props.score >= 50) return 'info'
  if (props.score >= 25) return 'warning'
  return 'error'
})

const bandLabel = computed(() => {
  if (props.score === null || props.score === undefined)
    return 'Insufficient data to score'
  if (props.score >= 75) return 'Very high likelihood of conversion'
  if (props.score >= 50) return 'High likelihood of conversion'
  if (props.score >= 25) return 'Medium likelihood of conversion'
  return 'Low likelihood of conversion'
})

const tooltipText = computed(() => {
  if (props.loading) return 'Calculating win probability…'
  if (props.score === null || props.score === undefined)
    return 'Score will appear once enough historical quotes are available'
  return `Win probability: ${props.score.toFixed(1)}% — ${bandLabel.value}`
})
</script>

<style scoped>
.probability-chip {
  font-weight: 600;
  min-width: 52px;
  justify-content: center;
}
</style>
