<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useStatusBarStore } from '@/renderer/store/statusBar'
import { useNetworkStatusStore } from '@/renderer/store/network_status'

const displayName = ref<string>('')
const statusBarStore = useStatusBarStore()
const networkStore = useNetworkStatusStore()

const apiDotColor = computed(() =>
  networkStore.isServiceAvailable ? '#4CAF50' : '#F44336'
)
const apiDotTitle = computed(() =>
  networkStore.isServiceAvailable ? 'API Connected' : 'API Offline'
)

onMounted(() => {
  const user = window.mainApi?.sendSync('msgGetAuthenticatedUser')
  displayName.value = user?.full_name || user?.username || ''
  networkStore.startCheckingService()
})
</script>

<template>
  <v-footer app height="24" class="status-bar px-3">
    <!-- API connection dot -->
    <span
      class="api-dot mr-3"
      :style="{ background: apiDotColor }"
      :title="apiDotTitle"
    />

    <!-- Context items injected by active screen -->
    <template v-for="(item, i) in statusBarStore.items" :key="i">
      <span v-if="i > 0" class="separator mx-2">·</span>
      <v-icon
        v-if="item.icon"
        size="12"
        class="mr-1"
        :color="
          item.severity === 'error'
            ? '#ff6b6b'
            : item.severity === 'warn'
              ? '#ffc107'
              : 'rgba(255,255,255,0.8)'
        "
        >{{ item.icon }}</v-icon
      >
      <span
        class="status-text"
        :style="
          item.severity === 'error'
            ? 'color:#ff6b6b'
            : item.severity === 'warn'
              ? 'color:#ffc107'
              : ''
        "
        >{{ item.text }}</span
      >
    </template>

    <span class="status-spacer" />

    <!-- Logged-in user -->
    <v-icon size="13" class="mr-1">mdi-account</v-icon>
    <span class="status-text">{{ displayName }}</span>
  </v-footer>
</template>

<style scoped>
.status-bar {
  background: var(--color-card-header);
  color: white;
  display: flex;
  align-items: center;
  font-size: 11px;
}
.api-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  display: inline-block;
  flex-shrink: 0;
}
.separator {
  opacity: 0.4;
}
.status-spacer {
  flex: 1;
}
.status-text {
  font-size: 11px;
  white-space: nowrap;
}
</style>
