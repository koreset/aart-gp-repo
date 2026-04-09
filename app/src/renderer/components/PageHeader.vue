<template>
  <div class="page-header mb-4">
    <v-breadcrumbs
      v-if="breadcrumbs?.length"
      :items="breadcrumbItems"
      density="compact"
      class="pa-0 mb-2"
    >
      <template #divider>
        <v-icon size="small">mdi-chevron-right</v-icon>
      </template>
    </v-breadcrumbs>
    <div class="d-flex align-center">
      <v-btn
        icon
        variant="text"
        size="small"
        class="mr-2"
        @click="$router.back()"
      >
        <v-icon>mdi-chevron-left</v-icon>
      </v-btn>
      <v-icon v-if="icon" :icon="icon" size="28" color="primary" class="mr-3" />
      <div>
        <h3 class="text-h6 font-weight-bold mb-0">{{ title }}</h3>
        <p v-if="subtitle" class="text-caption text-medium-emphasis mb-0">{{
          subtitle
        }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Breadcrumb {
  text: string
  to?: string | object
}

const props = defineProps<{
  title: string
  subtitle?: string
  icon?: string
  breadcrumbs?: Breadcrumb[]
}>()

const breadcrumbItems = computed(() =>
  (props.breadcrumbs ?? []).map((b) => ({
    title: b.text,
    to: b.to,
    disabled: !b.to
  }))
)
</script>

<style scoped>
.page-header {
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
  padding-bottom: 12px;
}
</style>
