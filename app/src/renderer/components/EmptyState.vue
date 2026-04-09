<template>
  <div class="empty-state d-flex flex-column align-center justify-center pa-8">
    <v-icon :icon="icon" size="56" color="grey-lighten-1" class="mb-4" />
    <p class="text-subtitle-1 font-weight-medium text-medium-emphasis mb-1">{{
      title
    }}</p>
    <p v-if="message" class="text-body-2 text-disabled text-center mb-4">{{
      message
    }}</p>
    <v-btn
      v-if="actionLabel && (actionRoute || actionFn)"
      color="primary"
      variant="tonal"
      size="small"
      :to="actionRoute"
      @click="actionFn && actionFn()"
    >
      {{ actionLabel }}
    </v-btn>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    icon?: string
    title: string
    message?: string
    actionLabel?: string
    actionRoute?: string | object
    actionFn?: () => void
  }>(),
  {
    icon: 'mdi-database-off-outline'
  }
)
</script>

<style scoped>
.empty-state {
  min-height: 200px;
}
</style>
