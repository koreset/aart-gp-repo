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
import { useUpdaterStore } from '@/renderer/store/updater'

const drawer = ref(true)
const confirmationDialog = ref()
const appStore = useAppStore()
const permissionsStore = useGroupUserPermissionsStore()
const updaterStore = useUpdaterStore()
const router = useRouter()
const expiryDismissed = ref(false)

const licenseData = computed(() => appStore.getLicenseData)

const licenseName = computed<string>(() => appStore.getOrganisationName)

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
    // Kick off the actual download. When it finishes, electron-updater
    // emits `update-downloaded` and the listener below prompts the user
    // to restart (which is what triggers `quitAndInstall`).
    window.mainApi?.send('msgStartDownload')
  }
})

window.mainApi?.on('update_downloaded', async (_event: any, info: any) => {
  log.info('Update Downloaded')
  // Drive the status-bar indicator into "ready" state so it stays visible
  // even if the user dismisses the restart dialog below — they can come
  // back to it later and the indicator will remind them.
  updaterStore.setDownloaded(info?.version || '')
  const res = await confirmationDialog.value?.open(
    'Update Downloaded',
    'The update has been downloaded. Do you want to restart the application now?'
  )
  if (res) {
    window.mainApi?.send('msgRestartApplication', true)
  }
})

window.mainApi?.on('download_progress', (_event: any, progress: any) => {
  // Push every tick into the updater store; the store rounds to integer
  // percent so re-renders are bounded to ~100 over a full download.
  updaterStore.setProgress(
    progress?.percent ?? 0,
    progress?.bytesPerSecond ?? 0
  )
})

window.mainApi?.on('update_error', async (_event: any, error: any) => {
  // electron-updater serializes the error before sending across IPC,
  // so `error` here is usually a plain object with a `message` field.
  const message =
    (error && (error.message || error.stack)) || String(error) || 'Unknown'

  // A 404 on `latest-<platform>.yml` (or "Cannot find channel ... update
  // info") just means the update server hasn't published a manifest for
  // this build yet. From the user's perspective that's identical to
  // "no update available", so swallow it silently — no log noise, no
  // store error, no UI surface. Real errors still flow through below.
  const isNoManifestError =
    /Cannot find channel/i.test(message) ||
    /latest-[a-z]+\.yml/i.test(message) ||
    /\b404\b/.test(message)
  if (isNoManifestError) {
    return
  }

  // Genuine errors (network failures, signature mismatch, install
  // problems) are recorded for debugging but no longer pop a dialog —
  // interrupting the user's flow with a wall of stack-trace text was
  // unhelpful. The status bar handles user-visible state if we ever
  // need to surface it; until then this stays quiet.
  log.error('Update Error:', message)
  updaterStore.setError(message)
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
        Your license expires on <strong>{{ expiryFormatted }}</strong> (<strong
          >{{ daysUntilExpiry }} day{{
            daysUntilExpiry === 1 ? '' : 's'
          }}</strong
        >
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
