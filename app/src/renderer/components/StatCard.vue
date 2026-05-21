<template>
  <div class="stat-card" :class="`stat-card--${resolvedTone}`">
    <div class="stat-card__stripe" />
    <div class="stat-card__body">
      <div class="stat-card__head">
        <span class="stat-card__label">{{ title }}</span>
        <v-icon v-if="icon" class="stat-card__icon" size="20">{{
          icon
        }}</v-icon>
      </div>
      <v-skeleton-loader
        v-if="loading"
        type="text"
        class="stat-card__skeleton"
      />
      <template v-else>
        <div class="stat-card__value">{{ value ?? '—' }}</div>
        <div v-if="subtitle" class="stat-card__hint">{{ subtitle }}</div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  title: string
  value?: string | number
  icon?: string
  color?: string
  loading?: boolean
  subtitle?: string
}>()

const toneMap: Record<string, string> = {
  primary: 'primary',
  accent: 'accent',
  info: 'info',
  success: 'success',
  warning: 'warning',
  error: 'error',
  grey: 'muted',
  '': 'muted'
}

const resolvedTone = computed(() => toneMap[props.color ?? ''] ?? 'muted')
</script>

<style scoped>
.stat-card {
  position: relative;
  display: flex;
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 10px;
  overflow: hidden;
  min-height: 108px;
  transition:
    transform 0.15s ease,
    box-shadow 0.15s ease,
    border-color 0.15s ease;
}

.stat-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 14px rgba(0, 63, 88, 0.08);
  border-color: rgba(var(--v-theme-on-surface), 0.14);
}

.stat-card__stripe {
  flex: 0 0 4px;
  background: rgba(var(--v-theme-on-surface), 0.2);
}

.stat-card__body {
  flex: 1;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.stat-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.stat-card__label {
  font-size: 0.72rem;
  font-weight: 600;
  letter-spacing: 0.4px;
  text-transform: uppercase;
  color: rgba(var(--v-theme-on-surface), 0.6);
  line-height: 1.2;
}

.stat-card__icon {
  color: rgba(var(--v-theme-on-surface), 0.45);
  flex-shrink: 0;
}

.stat-card__value {
  font-size: 1.6rem;
  font-weight: 700;
  line-height: 1.15;
  color: rgba(var(--v-theme-on-surface), 0.95);
  word-break: break-word;
}

.stat-card__hint {
  margin-top: auto;
  padding-top: 4px;
  font-size: 0.72rem;
  color: rgba(var(--v-theme-on-surface), 0.5);
  line-height: 1.3;
}

.stat-card__skeleton {
  background: transparent !important;
}

/* Tone variants */
.stat-card--primary .stat-card__stripe {
  background: rgb(var(--v-theme-primary));
}
.stat-card--primary .stat-card__icon {
  color: rgb(var(--v-theme-primary));
}

.stat-card--accent .stat-card__stripe {
  background: rgb(var(--v-theme-accent));
}
.stat-card--accent .stat-card__icon {
  color: rgb(var(--v-theme-accent));
}

.stat-card--info .stat-card__stripe {
  background: rgb(var(--v-theme-info));
}
.stat-card--info .stat-card__icon {
  color: rgb(var(--v-theme-info));
}

.stat-card--success .stat-card__stripe {
  background: rgb(var(--v-theme-success));
}
.stat-card--success .stat-card__icon {
  color: rgb(var(--v-theme-success));
}
.stat-card--success .stat-card__value {
  color: rgb(var(--v-theme-success));
}

.stat-card--warning .stat-card__stripe {
  background: rgb(var(--v-theme-warning));
}
.stat-card--warning .stat-card__icon {
  color: rgb(var(--v-theme-warning));
}
.stat-card--warning .stat-card__value {
  color: rgb(var(--v-theme-warning));
}

.stat-card--error .stat-card__stripe {
  background: rgb(var(--v-theme-error));
}
.stat-card--error .stat-card__icon {
  color: rgb(var(--v-theme-error));
}
.stat-card--error .stat-card__value {
  color: rgb(var(--v-theme-error));
}

.stat-card--muted .stat-card__stripe {
  background: rgba(var(--v-theme-on-surface), 0.18);
}
</style>
