<template>
  <v-app>
    <v-container fluid class="fill-height no-internet-bg">
      <v-row align="center" justify="center">
        <v-col cols="12" sm="8" md="5">
          <v-card class="pa-6 text-center" elevation="8" rounded="lg">
            <v-icon
              size="64"
              :color="serverUnreachable ? 'error' : 'warning'"
              class="mb-4"
            >
              {{ serverUnreachable ? 'mdi-server-off' : 'mdi-wifi-off' }}
            </v-icon>
            <h2 class="text-h5 mb-2">
              {{
                serverUnreachable
                  ? 'License Server Unavailable'
                  : 'No Internet Connection'
              }}
            </h2>
            <p class="text-body-2 text-medium-emphasis mb-6">
              <template v-if="serverUnreachable">
                Your internet connection is working but the license server could
                not be reached. This may be a temporary issue. Please try again
                or contact your administrator if the problem persists.
              </template>
              <template v-else>
                An internet connection is required to validate your license.
                Please check your network connection and try again.
              </template>
            </p>
            <v-btn color="primary" rounded :loading="retrying" @click="retry">
              Retry
            </v-btn>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-app>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue'

defineProps<{
  serverUnreachable?: boolean
}>()

const retrying = ref(false)

onBeforeMount(() => {
  window.mainApi?.send('msgResizeWindow', 1024, 600, false)
})

const retry = () => {
  retrying.value = true
  window.mainApi?.send('msgRestartApplication')
}
</script>

<style scoped>
.no-internet-bg {
  background: linear-gradient(135deg, #003f58 0%, #2e566e 100%);
}
</style>
