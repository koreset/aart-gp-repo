<template>
  <section class="page-section" :class="{ 'page-section--last': last }">
    <div class="section-header">
      <span class="section-label">{{ label }}</span>
      <v-chip
        v-if="count !== undefined"
        size="x-small"
        variant="tonal"
        color="primary"
        class="ml-2"
      >
        {{ count }}
      </v-chip>
      <span class="section-divider" />
      <slot name="actions" />
    </div>
    <slot />
  </section>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  count?: number | string
  last?: boolean
}>()
</script>

<style scoped>
.page-section {
  margin-bottom: 28px;
}
.page-section--last {
  margin-bottom: 0;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.section-label {
  font-size: 0.7rem;
  font-weight: 700;
  letter-spacing: 1.2px;
  text-transform: uppercase;
  color: rgb(var(--v-theme-primary));
}

.section-divider {
  flex: 1;
  height: 1px;
  background: linear-gradient(
    to right,
    rgba(var(--v-theme-primary), 0.18),
    rgba(var(--v-theme-primary), 0.02)
  );
}
</style>
