<template>
  <v-card class="mb-5">
    <v-card-title class="grey mb-4 lighten-2">
      <slot name="header">
        <div class="d-flex align-center">
          <v-icon v-if="icon" size="small" class="mr-2">{{ icon }}</v-icon>
          <span>{{ defaultHeader }}</span>
        </div>
        <div v-if="subtitle" class="text-caption font-weight-regular mt-1">{{
          subtitle
        }}</div>
      </slot>
    </v-card-title>
    <v-card-text>
      <slot />
    </v-card-text>
    <v-card-actions v-if="showActions">
      <slot name="actions"></slot>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
defineProps({
  defaultHeader: {
    type: String,
    default: 'Default Header'
  },
  showActions: {
    type: Boolean,
    default: true
  },
  icon: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  }
})
</script>

<style scoped>
.v-card-title {
  font-size: small;
  border-bottom: 1px solid rgba(0, 0, 0, 0.12);
  background-color: var(--color-card-header, #223f54);
  color: white;
}
</style>
