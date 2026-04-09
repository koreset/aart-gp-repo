<template>
  <v-container>
    <v-row>
      <v-col>
        <v-card v-if="!isServiceAvailable" class="mb-5">
          <v-card-text>
            <h3>Server Unavailable</h3>
            <p
              >Looks like the connection to the server is broken. Please reload
              this view</p
            >
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, computed } from 'vue'
import { useNetworkStatusStore } from '../store/network_status'

const networkStatusStore = useNetworkStatusStore()
const isOnline = computed(() => networkStatusStore.isOnline)
const isServiceAvailable = computed(() => networkStatusStore.isServiceAvailable)

const updateOnlineStatus = () => {
  networkStatusStore.updateOnlineStatus()
}

onMounted(() => {
  window.addEventListener('online', updateOnlineStatus)
  window.addEventListener('offline', updateOnlineStatus)
  if (isOnline.value) {
    networkStatusStore.startCheckingService()
  }
})

onUnmounted(() => {
  window.removeEventListener('online', updateOnlineStatus)
  window.removeEventListener('offline', updateOnlineStatus)
  networkStatusStore.stopCheckingService()
})
</script>

<style scoped></style>
