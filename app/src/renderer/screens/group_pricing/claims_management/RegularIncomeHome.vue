<template>
  <div class="ri-home">
    <v-tabs
      v-model="activeTab"
      color="primary"
      align-tabs="start"
      density="comfortable"
      class="ri-home-tabs"
      @update:model-value="syncTabToUrl"
    >
      <v-tab value="claims" prepend-icon="mdi-cash-multiple">
        Regular Income Claims
      </v-tab>
      <v-tab value="cpi" prepend-icon="mdi-chart-bar">CPI Index</v-tab>
    </v-tabs>
    <v-window v-model="activeTab" class="ri-home-window">
      <v-window-item value="claims">
        <regular-income-claims />
      </v-window-item>
      <v-window-item value="cpi">
        <cpi-index-management />
      </v-window-item>
    </v-window>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import RegularIncomeClaims from './RegularIncomeClaims.vue'
import CpiIndexManagement from './CpiIndexManagement.vue'

const route = useRoute()
const router = useRouter()

const validTabs = ['claims', 'cpi'] as const
type TabKey = (typeof validTabs)[number]

function initialTab(): TabKey {
  const q = route.query.tab
  const value = Array.isArray(q) ? q[0] : q
  if (value && validTabs.includes(value as TabKey)) {
    return value as TabKey
  }
  return 'claims'
}

const activeTab = ref<TabKey>(initialTab())

function syncTabToUrl(val: string | number | null) {
  if (typeof val !== 'string') return
  if (route.query.tab === val) return
  router.replace({ query: { ...route.query, tab: val } })
}

watch(
  () => route.query.tab,
  (next) => {
    const value = Array.isArray(next) ? next[0] : next
    if (value && validTabs.includes(value as TabKey)) {
      activeTab.value = value as TabKey
    } else if (!value) {
      activeTab.value = 'claims'
    }
  }
)
</script>

<style scoped>
.ri-home {
  width: 100%;
}

.ri-home-tabs {
  border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  padding-left: 8px;
}

.ri-home-window {
  background: transparent;
}
</style>
