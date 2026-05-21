<template>
  <v-tabs
    :model-value="activeTab"
    color="primary"
    density="compact"
    show-arrows
    class="finance-tabs"
  >
    <v-tab
      value="group-pricing-gl-bank-accounts"
      :to="{ name: 'group-pricing-gl-bank-accounts' }"
      prepend-icon="mdi-bank-outline"
    >
      Bank Accounts
    </v-tab>
    <v-tab
      value="group-pricing-gl-bank-reconciliation"
      :to="{ name: 'group-pricing-gl-bank-reconciliation' }"
      prepend-icon="mdi-bank-transfer"
    >
      Reconciliation
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
