<template>
  <div class="kpi-card" :class="`kpi-card--${tone || 'muted'}`">
    <div class="kpi-card__stripe" />
    <div class="kpi-card__body">
      <div class="kpi-card__head">
        <span class="kpi-card__label">{{ label }}</span>
        <v-icon v-if="icon" class="kpi-card__icon" size="20">{{
          icon
        }}</v-icon>
      </div>
      <div class="kpi-card__value">{{ value }}</div>
      <div v-if="hint" class="kpi-card__hint">{{ hint }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  value: string | number
  hint?: string
  icon?: string
  tone?:
    | 'primary'
    | 'accent'
    | 'success'
    | 'warning'
    | 'error'
    | 'info'
    | 'muted'
}>()
</script>

<style scoped>
.kpi-card {
  position: relative;
  display: flex;
  background: rgb(var(--v-theme-surface));
  border: 1px solid rgba(var(--v-theme-on-surface), 0.08);
  border-radius: 10px;
  overflow: hidden;
  min-height: 108px;
  height: 100%;
  transition:
    transform 0.15s ease,
    box-shadow 0.15s ease,
    border-color 0.15s ease;
}

.kpi-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 14px rgba(0, 63, 88, 0.08);
  border-color: rgba(var(--v-theme-on-surface), 0.14);
}

.kpi-card__stripe {
  flex: 0 0 4px;
  background: rgba(var(--v-theme-on-surface), 0.2);
}

.kpi-card__body {
  flex: 1;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.kpi-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.kpi-card__label {
  font-size: 0.72rem;
  font-weight: 600;
  letter-spacing: 0.4px;
  text-transform: uppercase;
  color: rgba(var(--v-theme-on-surface), 0.6);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.kpi-card__icon {
  color: rgba(var(--v-theme-on-surface), 0.45);
  flex: 0 0 auto;
}

.kpi-card__value {
  font-size: 1.6rem;
  font-weight: 700;
  line-height: 1.15;
  color: rgba(var(--v-theme-on-surface), 0.95);
  word-break: break-word;
}

.kpi-card__hint {
  margin-top: auto;
  padding-top: 4px;
  font-size: 0.72rem;
  color: rgba(var(--v-theme-on-surface), 0.5);
  line-height: 1.3;
}

/* Tone variants — accent stripe + icon tint within theme family */
.kpi-card--primary .kpi-card__stripe {
  background: rgb(var(--v-theme-primary));
}
.kpi-card--primary .kpi-card__icon {
  color: rgb(var(--v-theme-primary));
}

.kpi-card--accent .kpi-card__stripe {
  background: rgb(var(--v-theme-accent));
}
.kpi-card--accent .kpi-card__icon {
  color: rgb(var(--v-theme-accent));
}

.kpi-card--success .kpi-card__stripe {
  background: rgb(var(--v-theme-success));
}
.kpi-card--success .kpi-card__icon {
  color: rgb(var(--v-theme-success));
}

.kpi-card--warning .kpi-card__stripe {
  background: rgb(var(--v-theme-warning));
}
.kpi-card--warning .kpi-card__icon {
  color: rgb(var(--v-theme-warning));
}
.kpi-card--warning .kpi-card__value {
  color: rgb(var(--v-theme-warning));
}

.kpi-card--error .kpi-card__stripe {
  background: rgb(var(--v-theme-error));
}
.kpi-card--error .kpi-card__icon {
  color: rgb(var(--v-theme-error));
}
.kpi-card--error .kpi-card__value {
  color: rgb(var(--v-theme-error));
}

.kpi-card--info .kpi-card__stripe {
  background: rgb(var(--v-theme-info));
}
.kpi-card--info .kpi-card__icon {
  color: rgb(var(--v-theme-info));
}

.kpi-card--muted .kpi-card__stripe {
  background: rgba(var(--v-theme-on-surface), 0.18);
}
</style>
