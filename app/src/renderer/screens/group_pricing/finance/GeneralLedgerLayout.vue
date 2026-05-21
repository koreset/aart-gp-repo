<template>
  <v-tabs
    :model-value="activeTab"
    color="primary"
    density="compact"
    show-arrows
    class="finance-tabs"
  >
    <v-tab
      value="group-pricing-gl-journals"
      :to="{ name: 'group-pricing-gl-journals' }"
      prepend-icon="mdi-book-open-page-variant-outline"
    >
      Journals
    </v-tab>
    <v-tab
      value="group-pricing-gl-trial-balance"
      :to="{ name: 'group-pricing-gl-trial-balance' }"
      prepend-icon="mdi-scale-balance"
    >
      Trial Balance
    </v-tab>
    <v-tab
      value="group-pricing-gl-chart-of-accounts"
      :to="{ name: 'group-pricing-gl-chart-of-accounts' }"
      prepend-icon="mdi-format-list-bulleted-type"
    >
      Chart of Accounts
    </v-tab>
    <v-tab
      value="group-pricing-gl-periods"
      :to="{ name: 'group-pricing-gl-periods' }"
      prepend-icon="mdi-calendar-month-outline"
    >
      Periods
    </v-tab>
    <v-tab
      value="group-pricing-gl-posting-rules"
      :to="{ name: 'group-pricing-gl-posting-rules' }"
      prepend-icon="mdi-routes"
    >
      Posting Rules
    </v-tab>
    <v-tab
      value="group-pricing-gl-audit-log"
      :to="{ name: 'group-pricing-gl-audit-log' }"
      prepend-icon="mdi-history"
    >
      Audit Log
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
.finance-tabs {
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}
</style>
