<script setup lang="tsx">
import { DefaultLayout } from '@/renderer/components/layout'
import { useAppStore } from '@/renderer/store/app'
import { onBeforeMount, onMounted, onUnmounted, ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Toast from './components/Toast.vue'
import Api from '@/renderer/api/Api'
import { useWebSocket } from '@/renderer/composables/useWebSocket'
import { useNotificationStore } from '@/renderer/store/notifications'
import { useConversationStore } from '@/renderer/store/conversations'

const router = useRouter()
const route = useRoute()

// Routes that should always render, even while entitlements are loading
const bypassLoadingRoutes = ['app-settings', 'no-entitlements']
const shouldShowContent = computed(
  () =>
    !loadingEntitlements.value ||
    bypassLoadingRoutes.includes(route.name as string)
)
const appStore = useAppStore()
const notificationStore = useNotificationStore()
const conversationStore = useConversationStore()
const {
  connect: wsConnect,
  disconnect: wsDisconnect,
  on: wsOn
} = useWebSocket()

const licenseUrl = import.meta.env.VITE_APP_LICENSE_SERVER
const loadingEntitlements = ref(true)
const usingCachedEntitlements = ref(false)
const licenseRevoked = ref(false)
const licenseRevokedStatus = ref('')
let licenseCheckInterval: ReturnType<typeof setInterval> | null = null

const LICENSE_CHECK_INTERVAL = 5 * 60 * 1000 // 5 minutes

const checkLicensePeriodically = () => {
  licenseCheckInterval = setInterval(async () => {
    const status = await window.mainApi?.sendSync('msgCheckLicenseValidity')
    if (
      status === 'SUSPENDED' ||
      status === 'EXPIRED' ||
      status === 'INVALID' ||
      status === 'OVERDUE'
    ) {
      licenseRevoked.value = true
      licenseRevokedStatus.value = status
      // Stop further checks
      if (licenseCheckInterval) {
        clearInterval(licenseCheckInterval)
        licenseCheckInterval = null
      }
    }
  }, LICENSE_CHECK_INTERVAL)
}

// eslint-disable-next-line no-unused-vars
const handleLicenseRestart = () => {
  window.mainApi?.send('msgRestartApplication', false)
}

const getEntitlements = async (licenseId, licenseKey?) => {
  if (!licenseId) {
    return
  }
  try {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      Accept: 'application/json'
    }
    if (licenseKey) {
      headers.Authorization = 'Bearer ' + licenseKey
    }

    const response = await fetch(
      licenseUrl + '/licenses/' + licenseId + '/get-entitlements',
      {
        method: 'GET',
        headers
      }
    )

    const rs = await response.json()
    const entitlementList: any = []
    if (rs && rs.data && rs.data.length > 0) {
      rs.data.forEach((entitlement: any) => {
        entitlementList.push(entitlement.name)
      })
    }

    appStore.setEntitlements(entitlementList)

    // Cache entitlements for offline use
    if (entitlementList.length > 0) {
      window.mainApi?.sendSync('msgCacheEntitlements', entitlementList)
    }
  } catch (error) {
    console.error('Failed to fetch entitlements:', error)

    // Fall back to cached entitlements
    const cached = window.mainApi?.sendSync('msgGetCachedEntitlements')
    if (cached && cached.length > 0) {
      appStore.setEntitlements(cached)
      usingCachedEntitlements.value = true
    }
  }
}

onBeforeMount(async () => {
  try {
    window.mainApi?.send('msgResizeWindow', 1024, 600, true)
  } catch (error) {
    console.error('Error in onBeforeMount:', error)
  }
})

// Register a dynamic axios interceptor that attaches the user's entitlements
// as a header so the backend can enforce entitlement gates per route group.
Api.interceptors.request.use((config) => {
  const ents: string[] = appStore.getEntitlements() || []
  if (ents.length > 0) {
    config.headers['X-Entitlements'] = ents.join(',')
  }
  return config
})

onMounted(async () => {
  if (window.electronAPI?.onNavigate) {
    window.electronAPI.onNavigate((routeName: string) => {
      router.push({ name: routeName })
    })
  }

  const result = await window.mainApi?.sendSync('msgGetUserLicense')
  if (result) {
    appStore.setLicense(result)
  }

  // check if there are any entitlements
  let entitlements: any = appStore.entitlements

  if (entitlements.length === 0) {
    const licenseId = result?.data?.id || result?.id
    const licenseKey =
      result?.data?.attributes?.key || result?.attributes?.key || result?.key
    await getEntitlements(licenseId, licenseKey)
    entitlements = appStore.entitlements
  }

  loadingEntitlements.value = false

  if (entitlements.length === 0) {
    await router.push({ name: 'no-entitlements' })
  }

  // Start periodic license re-validation
  checkLicensePeriodically()

  // Connect WebSocket and register global message handlers
  wsConnect()

  // Handle incoming notifications globally
  wsOn('notification', (payload: any) => {
    notificationStore.addNotification(payload)
  })

  // Handle incoming conversation messages globally
  wsOn('conversation_message', (payload: any) => {
    conversationStore.addIncomingMessage(payload)
    // Also refresh inbox list so unread indicators update
    conversationStore.fetchInbox()
  })
})

onUnmounted(() => {
  if (licenseCheckInterval) {
    clearInterval(licenseCheckInterval)
    licenseCheckInterval = null
  }
  wsDisconnect()
})
</script>

<template>
  <DefaultLayout>
    <v-alert
      v-if="usingCachedEntitlements"
      type="info"
      variant="tonal"
      density="compact"
      closable
      class="mx-4 mt-2 mb-0"
    >
      License server unreachable — using cached entitlements. Some features may
      be out of date.
    </v-alert>
    <template v-if="loadingEntitlements && !shouldShowContent">
      <v-container>
        <v-row>
          <v-col>
            <v-skeleton-loader type="card" class="mb-4" />
            <v-skeleton-loader type="table-heading, table-row@3" />
          </v-col>
        </v-row>
      </v-container>
    </template>
    <template v-else>
      <router-view />
    </template>
    <Toast />

    <v-dialog v-model="licenseRevoked" persistent max-width="480">
      <v-card>
        <v-card-title class="text-h6">
          {{
            licenseRevokedStatus === 'SUSPENDED'
              ? 'License Suspended'
              : licenseRevokedStatus === 'EXPIRED'
                ? 'License Expired'
                : licenseRevokedStatus === 'OVERDUE'
                  ? 'License Check-In Overdue'
                  : 'License Invalid'
          }}
        </v-card-title>
        <v-card-text>
          {{
            licenseRevokedStatus === 'SUSPENDED'
              ? 'Your license has been suspended by your administrator. The application will restart so you can re-activate.'
              : licenseRevokedStatus === 'EXPIRED'
                ? 'Your license has expired. The application will restart so you can enter a new license key.'
                : licenseRevokedStatus === 'OVERDUE'
                  ? 'Your license is overdue for a check-in. The application will restart so you can re-register.'
                  : 'Your license is no longer valid. The application will restart so you can re-activate.'
          }}
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" @click="handleLicenseRestart">
            Restart Now
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </DefaultLayout>
</template>

<style>
:root {
  --color-primary: #003f58;
  --color-nav-bg: #2e566e;
  --color-card-header: #223f54;
  --color-main-bg: rgba(97, 136, 162, 0.7);
}

html {
  overflow-y: auto !important;
  user-select: none;
}
html,
body {
  width: 100%;
  height: 100%;
}
/* Do not force capitalization of button text */
.v-btn {
  text-transform: unset !important;
}

.v-main {
  background: rgba(97, 136, 162, 0.7);
}

.table-col {
  min-width: 120px;
  font-size: 12px;
  white-space: nowrap;
}

.table-row {
  background-color: #223f54;
  color: white;
  white-space: nowrap;
}

.v-table {
  border: 1px solid rgb(208, 207, 207);
}

.v-toolbar {
  height: 48px !important;
}
.v-expansion-panel--active > .v-expansion-panel-title {
  background-color: rgba(97, 136, 162, 0.7) !important;
  color: white !important;
  border-bottom-left-radius: 0 !important;
  border-bottom-right-radius: 0 !important;
  min-height: 50px !important;
}

/* AG Grid — Excel-like density & font */
.ag-theme-balham .ag-cell {
  color: #000000 !important;
  font-size: 13px !important;
  padding-top: 0 !important;
  padding-bottom: 0 !important;
  padding-left: 8px !important;
  padding-right: 8px !important;
  line-height: 22px !important;
}

.ag-theme-balham .ag-header-cell,
.ag-theme-balham .ag-header-group-cell {
  font-size: 13px !important;
  padding-top: 0 !important;
  padding-bottom: 0 !important;
}

/* AG Grid — dark navy header */
.ag-theme-balham .ag-header {
  background-color: var(--color-card-header) !important;
}
.ag-theme-balham .ag-header-cell-text {
  color: white !important;
}
.ag-theme-balham .ag-header-icon {
  color: rgba(255, 255, 255, 0.7) !important;
}
.ag-theme-balham .ag-row {
  height: 32px !important;
}
.ag-theme-balham .ag-cell {
  line-height: 32px !important;
}
.ag-theme-balham .ag-row-hover {
  background-color: rgba(0, 63, 88, 0.06) !important;
}

/* Row-level highlight helpers */
.ag-row-overdue {
  background-color: rgba(244, 67, 54, 0.08) !important;
}
.ag-row-warning {
  background-color: rgba(255, 152, 0, 0.08) !important;
}
</style>
