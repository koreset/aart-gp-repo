<template>
  <v-tabs
    :model-value="activeTab"
    color="primary"
    density="compact"
    show-arrows
    class="premium-receipts-tabs"
  >
    <v-tab
      value="group-pricing-premium-dashboard"
      :to="{ name: 'group-pricing-premium-dashboard' }"
      prepend-icon="mdi-view-dashboard-outline"
    >
      Dashboard
    </v-tab>
    <v-tab
      value="group-pricing-premium-schedules"
      :to="{ name: 'group-pricing-premium-schedules' }"
      prepend-icon="mdi-calendar-month-outline"
    >
      Schedules
    </v-tab>
    <v-tab
      value="group-pricing-invoices"
      :to="{ name: 'group-pricing-invoices' }"
      prepend-icon="mdi-receipt-text-outline"
    >
      Invoices
    </v-tab>
    <v-tab
      value="group-pricing-payments"
      :to="{ name: 'group-pricing-payments' }"
      prepend-icon="mdi-bank-transfer"
    >
      Payments
    </v-tab>
    <v-tab
      value="group-pricing-premium-reconciliation"
      :to="{ name: 'group-pricing-premium-reconciliation' }"
      prepend-icon="mdi-scale-balance"
    >
      Reconciliation
    </v-tab>
    <v-tab
      value="group-pricing-arrears"
      :to="{ name: 'group-pricing-arrears' }"
      prepend-icon="mdi-clock-alert-outline"
    >
      Arrears
    </v-tab>
    <v-tab
      value="group-pricing-statements"
      :to="{ name: 'group-pricing-statements' }"
      prepend-icon="mdi-file-chart-outline"
    >
      Statements
    </v-tab>
  </v-tabs>

  <router-view />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const activeTab = computed(() => (route.name as string) ?? '')
</script>

<style scoped>
.premium-receipts-tabs {
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}
</style>
