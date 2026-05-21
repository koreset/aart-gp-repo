<template>
  <div class="claims-home">
    <v-tabs
      v-model="activeTab"
      color="primary"
      align-tabs="start"
      density="comfortable"
      class="claims-home-tabs"
      @update:model-value="syncTabToUrl"
    >
      <v-tab value="dashboard" prepend-icon="mdi-view-dashboard-outline">
        Dashboard
      </v-tab>
      <v-tab value="claims" prepend-icon="mdi-file-document-multiple-outline">
        Claims
      </v-tab>
    </v-tabs>
    <v-window v-model="activeTab" class="claims-home-window">
      <v-window-item value="dashboard">
        <claims-analytics />
      </v-window-item>
      <v-window-item value="claims">
        <claims-management />
      </v-window-item>
    </v-window>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ClaimsAnalytics from './ClaimsAnalytics.vue'
import ClaimsManagement from './ClaimsManagement.vue'

const route = useRoute()
const router = useRouter()

const validTabs = ['dashboard', 'claims'] as const
type TabKey = (typeof validTabs)[number]

function initialTab(): TabKey {
  const q = route.query.tab
  const value = Array.isArray(q) ? q[0] : q
  if (value && validTabs.includes(value as TabKey)) {
    return value as TabKey
  }
  return 'dashboard'
}

const activeTab = ref<TabKey>(initialTab())

function syncTabToUrl(val: string | number | null) {
  if (typeof val !== 'string') return
  if (route.query.tab === val) return
  router.replace({ query: { ...route.query, tab: val } })
}

// Keep tab in sync if the user navigates with a different ?tab param
// (e.g. via a deep-link click while already on this page).
watch(
  () => route.query.tab,
  (next) => {
    const value = Array.isArray(next) ? next[0] : next
    if (value && validTabs.includes(value as TabKey)) {
      activeTab.value = value as TabKey
    } else if (!value) {
      activeTab.value = 'dashboard'
    }
  }
)
</script>

<style scoped>
.claims-home {
  width: 100%;
}

.claims-home-tabs {
  border-bottom: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  padding-left: 8px;
}

.claims-home-window {
  background: transparent;
}
</style>
