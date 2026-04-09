<template>
  <v-snackbar v-model="show" color="error" location="bottom" :timeout="-1">
    <v-icon start>mdi-server-network-off</v-icon>
    Service unavailable. Please try again later.
    <template #actions>
      <v-btn variant="text" color="white" @click="show = false">Dismiss</v-btn>
    </template>
  </v-snackbar>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { useNetworkStatusStore } from '../store/network_status'

const networkStatusStore = useNetworkStatusStore()
const isServiceAvailable = computed(() => networkStatusStore.isServiceAvailable)
const show = ref(false)

watch(isServiceAvailable, (v) => {
  show.value = !v
})

const updateOnlineStatus = () => {
  networkStatusStore.updateOnlineStatus()
}

onMounted(() => {
  window.addEventListener('online', updateOnlineStatus)
  window.addEventListener('offline', updateOnlineStatus)
  networkStatusStore.startCheckingService()
})

onUnmounted(() => {
  window.removeEventListener('online', updateOnlineStatus)
  window.removeEventListener('offline', updateOnlineStatus)
  networkStatusStore.stopCheckingService()
})
</script>
