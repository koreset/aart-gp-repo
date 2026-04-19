<script setup lang="tsx">
import NavigationDrawer from '@/renderer/screens/NavigationDrawer.vue'
import { useRoute, useRouter } from 'vue-router'
import { ref, computed } from 'vue'
import ServerUnavailable from '../ServerUnavailable.vue'
import log from 'electron-log'
// import { ipcRenderer } from 'electron'
import ConfirmDialog from '../ConfirmDialog.vue'
import StatusBar from './StatusBar.vue'
import NotificationBell from '@/renderer/components/NotificationBell.vue'
import { useAppStore } from '@/renderer/store/app'
import { useGroupUserPermissionsStore } from '@/renderer/store/group_user'

const drawer = ref(true)
const confirmationDialog = ref()
const appStore = useAppStore()
const permissionsStore = useGroupUserPermissionsStore()
const router = useRouter()
const expiryDismissed = ref(false)

const licenseData = computed(() => appStore.getLicenseData)

const licenseName = computed(() => {
  if (!licenseData.value) return ''
  const meta =
    licenseData.value.data?.attributes?.metadata ||
    licenseData.value.attributes?.metadata
  return meta?.organization || meta?.userName || ''
})

const licenseExpiry = computed(() => {
  if (!licenseData.value) return null
  const expiry =
    licenseData.value.data?.attributes?.expiry ||
    licenseData.value.attributes?.expiry
  return expiry ? new Date(expiry) : null
})

const daysUntilExpiry = computed(() => {
  if (!licenseExpiry.value) return null
  const now = new Date()
  const diff = licenseExpiry.value.getTime() - now.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

const showExpiryWarning = computed(() => {
  if (expiryDismissed.value) return false
  return (
    daysUntilExpiry.value !== null &&
    daysUntilExpiry.value <= 30 &&
    daysUntilExpiry.value > 0
  )
})

const expiryAlertType = computed(() =>
  daysUntilExpiry.value !== null && daysUntilExpiry.value <= 7
    ? 'error'
    : 'warning'
)

const expiryFormatted = computed(() => {
  if (!licenseExpiry.value) return ''
  return licenseExpiry.value.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const route: any = useRoute()
const titleKey: string = (route?.meta?.titleKey || 'title.main') as string

const toggleDrawer = (): void => {
  drawer.value = !drawer.value
}

window.mainApi?.on('update_available', async () => {
  log.info('Update Available')
  const res = await confirmationDialog.value?.open(
    'Update Available',
    'A new version of the application is available. Do you want to update now? You will be notified when the download completes.'
  )
  if (res) {
    window.mainApi?.send('msgRestartApplication', false)
  }
})

window.mainApi?.on('update_downloaded', async () => {
  log.info('Update Downloaded')
  const res = await confirmationDialog.value?.open(
    'Update Downloaded',
    'The update has been downloaded. Do you want to restart the application now?'
  )
  if (res) {
    window.mainApi?.send('msgRestartApplication', true)
  }
})

window.mainApi?.on('download_progress', (event: any, progress: any) => {
  log.info('Download Progress:', progress)
})

window.mainApi?.on('update_error', async (event: any, error: any) => {
  log.error('Update Error Event:', event)
  log.error('Update Error:', error)
})

window.mainApi?.on('logout', async () => {
  log.info('Logout')
  const res = await confirmationDialog.value?.open(
    'Logout',
    'Are you sure you want to logout?'
  )
  if (res) {
    try {
      // Clear frontend state
      appStore.clearAll()
      permissionsStore.clearPermissions()

      // Send logout message to main process
      window.mainApi?.send('msgLogout')
    } catch (error) {
      console.error('Error during logout:', error)
    }
  }
})
</script>

<template>
  <v-app>
    <v-app-bar
      class="custom-app"
      color="primary"
      density="compact"
      elevation="2"
    >
      <v-app-bar-nav-icon
        variant="text"
        @click="toggleDrawer()"
      ></v-app-bar-nav-icon>
      <v-app-bar-title class="app-title">{{ $t(titleKey) }}</v-app-bar-title>
      <template #append>
        <v-chip
          v-if="licenseName"
          size="small"
          variant="tonal"
          color="white"
          class="mr-2 license-chip"
          @click="router.push({ name: 'app-settings' })"
        >
          <v-icon start size="14">mdi-license</v-icon>
          {{ licenseName }}
          <v-tooltip activator="parent" location="bottom">
            <template v-if="licenseExpiry">
              Expires: {{ expiryFormatted }}
            </template>
            <template v-else> License info </template>
          </v-tooltip>
        </v-chip>
        <NotificationBell />
      </template>
    </v-app-bar>
    <NavigationDrawer :drawer="drawer" />
    <v-main>
      <v-alert
        v-if="showExpiryWarning"
        :type="expiryAlertType"
        variant="flat"
        density="compact"
        closable
        prominent
        class="mx-4 mt-2 mb-0"
        @click:close="expiryDismissed = true"
      >
        Your license expires on <strong>{{ expiryFormatted }}</strong>
        (<strong>{{ daysUntilExpiry }} day{{ daysUntilExpiry === 1 ? '' : 's' }}</strong>
        remaining). Please contact your administrator to renew.
      </v-alert>
      <server-unavailable />

      <slot />
    </v-main>
    <StatusBar />
    <confirm-dialog ref="confirmationDialog" />
  </v-app>
</template>

<style scoped>
.custom-app {
  height: 40px;
}

.app-title {
  font-size: 14px;
  font-weight: 500;
}

.license-chip {
  cursor: pointer;
  font-size: 11px;
}
</style>
