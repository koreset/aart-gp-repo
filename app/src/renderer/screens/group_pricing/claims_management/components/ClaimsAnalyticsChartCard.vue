<template>
  <div class="chart-card">
    <div class="chart-card__head">
      <div class="chart-card__titles">
        <div class="chart-card__title">{{ title }}</div>
        <div v-if="subtitle" class="chart-card__subtitle">{{ subtitle }}</div>
      </div>
      <div class="chart-card__head-actions">
        <slot name="actions" />
        <v-menu v-if="actions && actions.length" location="bottom end">
          <template #activator="{ props: menuProps }">
            <v-btn
              v-bind="menuProps"
              icon="mdi-menu"
              size="x-small"
              variant="text"
              density="comfortable"
              :aria-label="`${title} actions`"
              class="chart-card__menu-btn"
            />
          </template>
          <v-list density="compact">
            <v-list-item
              v-for="(item, idx) in actions"
              :key="idx"
              :prepend-icon="item.icon"
              :title="item.label"
              @click="item.onClick"
            />
          </v-list>
        </v-menu>
      </div>
    </div>
    <div class="chart-card__body" :style="{ height: bodyHeight }">
      <div v-if="loading" class="chart-card__loading">
        <v-progress-circular indeterminate size="36" color="primary" />
      </div>
      <div v-else-if="empty" class="chart-card__empty">
        <v-icon size="32" color="grey-lighten-1">mdi-chart-line</v-icon>
        <div class="chart-card__empty-text">{{ emptyText }}</div>
      </div>
      <slot v-else />
    </div>
  </div>
</template>

<script setup lang="ts">
interface ChartAction {
  label: string
  icon?: string
  onClick: () => void
}

defineProps<{
  title: string
  subtitle?: string
  loading?: boolean
  empty?: boolean
  emptyText?: string
  bodyHeight?: string
  actions?: ChartAction[]
}>()
</script>

<style scoped>
.chart-card {
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.chart-card__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 16px 10px;
  border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.06);
}

.chart-card__head-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.chart-card__menu-btn {
  opacity: 0.7;
  transition: opacity 0.15s ease;
}

.chart-card__menu-btn:hover {
  opacity: 1;
}

.chart-card__titles {
  min-width: 0;
}

.chart-card__title {
  font-size: 0.95rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.9);
  line-height: 1.25;
}

.chart-card__subtitle {
  font-size: 0.72rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
  margin-top: 2px;
}

.chart-card__body {
  flex: 1;
  padding: 8px 12px 12px;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.chart-card__loading,
.chart-card__empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 8px;
  color: rgba(var(--v-theme-on-surface), 0.45);
  font-size: 0.85rem;
}

.chart-card__empty-text {
  text-align: center;
}
</style>
