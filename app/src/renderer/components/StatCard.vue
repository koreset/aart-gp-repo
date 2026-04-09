<template>
  <v-card variant="tonal" :color="color" class="stat-card">
    <v-card-text class="pa-3">
      <div class="d-flex justify-space-between align-start">
        <div class="flex-grow-1">
          <p class="text-caption text-medium-emphasis mb-1">{{ title }}</p>
          <v-skeleton-loader v-if="loading" type="text" class="mt-1" />
          <template v-else>
            <p class="text-h6 font-weight-bold">{{ value }}</p>
            <p v-if="subtitle" class="text-caption mt-1">{{ subtitle }}</p>
          </template>
        </div>
        <v-icon v-if="icon" size="32" :color="color" class="ml-2 mt-1">{{
          icon
        }}</v-icon>
      </div>
    </v-card-text>
  </v-card>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  value?: string | number
  icon?: string
  color?: string
  loading?: boolean
  subtitle?: string
}>()
</script>

<style scoped>
.stat-card {
  min-height: 80px;
}
</style>
