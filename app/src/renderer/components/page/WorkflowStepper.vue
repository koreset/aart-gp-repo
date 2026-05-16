<template>
  <div class="workflow-stepper">
    <div
      v-for="(step, idx) in steps"
      :key="step.label"
      class="workflow-stepper__step"
      :class="`workflow-stepper__step--${step.tone}`"
    >
      <div class="workflow-stepper__num">{{ idx + 1 }}</div>
      <div class="workflow-stepper__text">
        <div class="workflow-stepper__label">{{ step.label }}</div>
        <div class="workflow-stepper__sub">{{ step.sub }}</div>
      </div>
      <v-icon
        v-if="idx < steps.length - 1"
        class="workflow-stepper__arrow"
        size="18"
        >mdi-chevron-right</v-icon
      >
    </div>
  </div>
</template>

<script setup lang="ts">
export interface WorkflowStep {
  label: string
  sub: string
  tone:
    | 'muted'
    | 'primary'
    | 'accent'
    | 'success'
    | 'warning'
    | 'error'
    | 'info'
}

defineProps<{
  steps: WorkflowStep[]
}>()
</script>

<style scoped>
.workflow-stepper {
  display: flex;
  align-items: stretch;
  flex-wrap: wrap;
  gap: 8px;
  padding: 14px 16px;
  background: rgba(var(--v-theme-primary), 0.04);
  border: 1px solid rgba(var(--v-theme-primary), 0.1);
  border-radius: 10px;
}

.workflow-stepper__step {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1 1 0;
  min-width: 150px;
}

.workflow-stepper__num {
  flex: 0 0 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.78rem;
  font-weight: 700;
  color: #fff;
  background: rgba(var(--v-theme-on-surface), 0.35);
}

.workflow-stepper__step--muted .workflow-stepper__num {
  background: rgba(var(--v-theme-on-surface), 0.35);
}
.workflow-stepper__step--primary .workflow-stepper__num {
  background: rgb(var(--v-theme-primary));
}
.workflow-stepper__step--warning .workflow-stepper__num {
  background: rgb(var(--v-theme-warning));
}
.workflow-stepper__step--accent .workflow-stepper__num {
  background: rgb(var(--v-theme-accent));
}
.workflow-stepper__step--info .workflow-stepper__num {
  background: rgb(var(--v-theme-info));
}
.workflow-stepper__step--success .workflow-stepper__num {
  background: rgb(var(--v-theme-success));
}
.workflow-stepper__step--error .workflow-stepper__num {
  background: rgb(var(--v-theme-error));
}

.workflow-stepper__text {
  flex: 1;
  min-width: 0;
}

.workflow-stepper__label {
  font-size: 0.82rem;
  font-weight: 600;
  color: rgba(var(--v-theme-on-surface), 0.88);
  line-height: 1.2;
}

.workflow-stepper__sub {
  font-size: 0.7rem;
  color: rgba(var(--v-theme-on-surface), 0.55);
  line-height: 1.2;
}

.workflow-stepper__arrow {
  color: rgba(var(--v-theme-primary), 0.35);
  flex: 0 0 auto;
}
</style>
