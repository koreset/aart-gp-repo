<!-- eslint-disable vue/first-attribute-linebreak -->
<template>
  <v-navigation-drawer
    v-model="internalDrawer"
    class="custom-navbar nav-bg drawer-text"
    dark
    color="navbar"
    location="left"
  >
    <v-container>
      <v-row>
        <v-col>
          <img class="ml-3" width="100%" :src="'./images/aart-logo-02.png'" />
          <p class="nav-text">App Version: {{ appVersion }}</p>
        </v-col>
      </v-row>
    </v-container>
    <v-list class="nav-text smaller-font">
      <!-- Dashboard -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:view_gp_dashboard') }"
        prepend-icon="mdi-view-dashboard-outline"
        @click="
          navigate('group-pricing-dashboard', 'navigation:view_gp_dashboard')
        "
      >
        <v-list-item-title>Dashboard</v-list-item-title>
      </v-list-item>

      <!-- Quotes -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:view_quotes') }"
        prepend-icon="mdi-file-document-outline"
        @click="navigate('group-pricing-quotes', 'navigation:view_quotes')"
      >
        <v-list-item-title>Quotes</v-list-item-title>
      </v-list-item>

      <!-- Tables -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:group_tables') }"
        prepend-icon="mdi-table"
        @click="navigate('group-pricing-tables', 'navigation:group_tables')"
      >
        <v-list-item-title>Tables</v-list-item-title>
      </v-list-item>

      <!-- Metadata -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:view_metadata') }"
        prepend-icon="mdi-tag-outline"
        @click="navigate('group-pricing-metadata', 'navigation:view_metadata')"
      >
        <v-list-item-title>Metadata</v-list-item-title>
      </v-list-item>

      <!-- Scheme Management -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:view_schemes') }"
        prepend-icon="mdi-domain"
        @click="navigate('group-pricing-schemes', 'navigation:view_schemes')"
      >
        <v-list-item-title>Scheme Management</v-list-item-title>
      </v-list-item>

      <!-- Member Management -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:manage_members') }"
        prepend-icon="mdi-account-group-outline"
        @click="
          navigate(
            'group-pricing-member-management',
            'navigation:manage_members'
          )
        "
      >
        <v-list-item-title>Member Management</v-list-item-title>
      </v-list-item>

      <!-- Scheme Migration -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:manage_scheme_migration') }"
        prepend-icon="mdi-database-import-outline"
        @click="
          navigate(
            'group-pricing-scheme-migration',
            'navigation:manage_scheme_migration'
          )
        "
      >
        <v-list-item-title>Scheme Migration</v-list-item-title>
      </v-list-item>

      <!-- Claims Management -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:manage_claims') }"
        prepend-icon="mdi-file-clock-outline"
        @click="
          navigate(
            'group-pricing-claims-management',
            'navigation:manage_claims'
          )
        "
      >
        <v-list-item-title>Claims Management</v-list-item-title>
      </v-list-item>

      <!-- Claims Analytics -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:view_claims_analytics') }"
        prepend-icon="mdi-chart-line"
        @click="
          navigate('group-pricing-claims-analytics', 'navigation:view_claims_analytics')
        "
      >
        <v-list-item-title>Claims Analytics</v-list-item-title>
      </v-list-item>

      <!-- Bordereaux Management -->
      <v-list-item
        :class="{ 'disabled-item': !canAccess('navigation:manage_bordereaux') }"
        prepend-icon="mdi-folder-table-outline"
        @click="
          navigate(
            'group-pricing-bordereaux-management',
            'navigation:manage_bordereaux'
          )
        "
      >
        <v-list-item-title>Bordereaux Management</v-list-item-title>
      </v-list-item>

      <v-divider class="my-2"></v-divider>

      <!-- PHI -->
      <v-list-group v-model="expandedGroups" value="PHI">
        <template #activator="{ props }">
          <v-list-item
            v-bind="props"
            :class="{ 'disabled-item': !canAccess('navigation:view_phi') }"
            prepend-icon="mdi-heart-pulse"
            title="PHI"
          ></v-list-item>
        </template>
        <v-list-item
          :class="{
            'disabled-item': !canAccess('navigation:view_phi_tables')
          }"
          prepend-icon="mdi-table-heart"
          @click="
            navigateGroup(
              'group-pricing-phi-tables',
              'PHI',
              'navigation:view_phi_tables'
            )
          "
        >
          <v-list-item-title>Tables</v-list-item-title>
        </v-list-item>
        <v-list-item
          :class="{
            'disabled-item': !canAccess('navigation:view_phi_run_settings')
          }"
          prepend-icon="mdi-cog-play-outline"
          @click="
            navigateGroup(
              'group-pricing-phi-run-settings',
              'PHI',
              'navigation:view_phi_run_settings'
            )
          "
        >
          <v-list-item-title>Run Settings</v-list-item-title>
        </v-list-item>
        <v-list-item
          :class="{
            'disabled-item': !canAccess('navigation:view_phi_shock_settings')
          }"
          prepend-icon="mdi-weather-lightning"
          @click="
            navigateGroup(
              'group-pricing-phi-shock-settings',
              'PHI',
              'navigation:view_phi_shock_settings'
            )
          "
        >
          <v-list-item-title>Shock Settings</v-list-item-title>
        </v-list-item>
        <v-list-item
          :class="{
            'disabled-item': !canAccess('navigation:view_phi_run_results')
          }"
          prepend-icon="mdi-chart-bar"
          @click="
            navigateGroup(
              'group-pricing-phi-run-results',
              'PHI',
              'navigation:view_phi_run_results'
            )
          "
        >
          <v-list-item-title>Run Results</v-list-item-title>
        </v-list-item>
      </v-list-group>

      <!-- Premiums -->
      <v-list-group v-model="expandedGroups" value="Premiums">
        <template #activator="{ props }">
          <v-list-item
            v-bind="props"
            :class="{
              'disabled-item': !canAccess('navigation:view_premium_dashboard')
            }"
            prepend-icon="mdi-cash-multiple"
            title="Premiums"
          ></v-list-item>
        </template>
        <v-list-item
          :class="{
            'disabled-item': !canAccess('navigation:view_premium_dashboard')
          }"
          prepend-icon="mdi-view-dashboard-outline"
          @click="
            navigateGroup(
              'group-pricing-premium-dashboard',
              'Premiums',
              'navigation:view_premium_dashboard'
            )
          "
        >
          <v-list-item-title>Dashboard</v-list-item-title>
        </v-list-item>
        <v-list-item
          :class="{
            'disabled-item': !canAccess('navigation:manage_premium_schedules')
          }"
          prepend-icon="mdi-calendar-month-outline"
          @click="
            navigateGroup(
              'group-pricing-premium-schedules',
              'Premiums',
              'navigation:manage_premium_schedules'
            )
          "
        >
          <v-list-item-title>Premium Schedules</v-list-item-title>
        </v-list-item>
        <v-list-item
          :class="{ 'disabled-item': !canAccess('navigation:manage_invoices') }"
          prepend-icon="mdi-receipt-text-outline"
          @click="
            navigateGroup(
              'group-pricing-invoices',
              'Premiums',
              'navigation:manage_invoices'
            )
          "
        >
          <v-list-item-title>Invoices</v-list-item-title>
        </v-list-item>
        <v-list-item
          :class="{ 'disabled-item': !canAccess('navigation:manage_payments') }"
          prepend-icon="mdi-bank-transfer"
          @click="
            navigateGroup(
              'group-pricing-payments',
              'Premiums',
              'navigation:manage_payments'
            )
          "
        >
          <v-list-item-title>Payments</v-list-item-title>
        </v-list-item>
        <v-list-item
          :class="{
            'disabled-item': !canAccess(
              'navigation:manage_premium_reconciliation'
            )
          }"
          prepend-icon="mdi-scale-balance"
          @click="
            navigateGroup(
              'group-pricing-premium-reconciliation',
              'Premiums',
              'navigation:manage_premium_reconciliation'
            )
          "
        >
          <v-list-item-title>Reconciliation</v-list-item-title>
        </v-list-item>
        <v-list-item
          :class="{ 'disabled-item': !canAccess('navigation:manage_arrears') }"
          prepend-icon="mdi-clock-alert-outline"
          @click="
            navigateGroup(
              'group-pricing-arrears',
              'Premiums',
              'navigation:manage_arrears'
            )
          "
        >
          <v-list-item-title>Arrears</v-list-item-title>
        </v-list-item>
        <v-list-item
          :class="{ 'disabled-item': !canAccess('navigation:view_statements') }"
          prepend-icon="mdi-file-chart-outline"
          @click="
            navigateGroup(
              'group-pricing-statements',
              'Premiums',
              'navigation:view_statements'
            )
          "
        >
          <v-list-item-title>Statements</v-list-item-title>
        </v-list-item>
      </v-list-group>

      <v-divider class="my-2"></v-divider>

      <!-- User Management -->
      <v-list-group v-model="expandedGroups" value="User Management">
        <template #activator="{ props }">
          <v-list-item
            v-bind="props"
            :class="{ 'disabled-item': !canAccess('navigation:manage_users') }"
            prepend-icon="mdi-account-cog-outline"
            title="User Management"
          ></v-list-item>
        </template>
        <v-list-item
          class="second-level-item"
          :class="{ 'disabled-item': !canAccess('navigation:manage_users') }"
          prepend-icon="mdi-account-multiple-outline"
          @click="
            navigateGroup(
              'user-management-list',
              'User Management',
              'navigation:manage_users'
            )
          "
        >
          <v-list-item-title>Users List</v-list-item-title>
        </v-list-item>
        <v-list-item
          class="second-level-item"
          :class="{ 'disabled-item': !canAccess('navigation:manage_users') }"
          prepend-icon="mdi-shield-account-outline"
          @click="
            navigateGroup(
              'user-management-roles',
              'User Management',
              'navigation:manage_users'
            )
          "
        >
          <v-list-item-title>Roles</v-list-item-title>
        </v-list-item>
      </v-list-group>

      <!-- Messages -->
      <v-list-item
        :to="{ name: 'messages-inbox' }"
        prepend-icon="mdi-message-text-outline"
      >
        <v-list-item-title>Messages</v-list-item-title>
      </v-list-item>

      <!-- Documentation -->
      <v-list-item
        :to="{ name: 'documentation' }"
        prepend-icon="mdi-book-open-page-variant-outline"
      >
        <v-list-item-title>Documentation</v-list-item-title>
      </v-list-item>

      <v-list-item
        :to="{ name: 'app-settings' }"
        :prepend-icon="'mdi-cog-outline'"
      >
        <v-list-item-title>Application Settings</v-list-item-title>
      </v-list-item>

      <v-divider class="my-2"></v-divider>

      <v-list-item
        :prepend-icon="'mdi-logout'"
        class="logout-item"
        @click="handleLogout"
      >
        <v-list-item-title>Logout</v-list-item-title>
      </v-list-item>
    </v-list>

    <v-snackbar v-model="snackbar" :timeout="timeout">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>

    <confirm-dialog ref="confirmationDialog" />
  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { watchEffect, ref, onMounted, watch } from 'vue'
import { useGroupUserPermissionsStore } from '@/renderer/store/group_user'
import GroupPricingService from '@/renderer/api/GroupPricingService'
import { useRouter } from 'vue-router'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import { useAppStore } from '@/renderer/store/app'

const router = useRouter()

const navProps = defineProps({
  drawer: {
    type: Boolean,
    required: true
  }
})

const appStore = useAppStore()
const permissionsStore = useGroupUserPermissionsStore()
const snackbar = ref(false)
const snackbarMessage = ref('')
const timeout = ref(3000)

// Track expanded groups to prevent them from collapsing
const expandedGroups = ref<string[]>([])

// Flag to prevent saving during initial load
const isInitialLoad = ref(true)

// Load expanded groups from localStorage on component mount
const loadExpandedGroups = () => {
  try {
    const saved = localStorage.getItem('navigation-expanded-groups')
    if (saved) {
      expandedGroups.value = JSON.parse(saved)
    }
  } catch (error) {
    console.error('Error loading expanded groups from localStorage:', error)
  }
  isInitialLoad.value = false
}

// Save expanded groups to localStorage whenever they change
const saveExpandedGroups = () => {
  if (isInitialLoad.value) return
  try {
    localStorage.setItem(
      'navigation-expanded-groups',
      JSON.stringify(expandedGroups.value)
    )
  } catch (error) {
    console.error('Error saving expanded groups to localStorage:', error)
  }
}

watchEffect(() => {
  saveExpandedGroups()
})

// Accordion: only one group expanded at a time
watch(
  expandedGroups,
  (newValue) => {
    if (newValue.length > 1) {
      expandedGroups.value = [newValue[newValue.length - 1]]
    }
  },
  { deep: true }
)

loadExpandedGroups()

/**
 * Check if user has the given permission slug.
 * Returns true if:
 *  - permissions haven't loaded yet (graceful)
 *  - no role assigned (empty permissions)
 *  - user has system:admin
 *  - user has the specific permission
 */
const canAccess = (permission: string): boolean => {
  if (!permissionsStore.loaded) return true
  const perms = permissionsStore.permissions
  if (Object.keys(perms).length === 0) return true
  if (perms['system:admin']) return true
  return !!perms[permission]
}

const navigate = (routeName: string, permission?: string) => {
  if (permission && !canAccess(permission)) {
    snackbarMessage.value = `You don't have permission to access this feature`
    snackbar.value = true
    return
  }
  router.push({ name: routeName })
}

const expandGroup = (groupName: string) => {
  expandedGroups.value = [groupName]
}

const navigateGroup = (
  routeName: string,
  groupName: string,
  permission?: string
) => {
  navigate(routeName, permission)
  expandGroup(groupName)
}

const appVersion = ref('')

onMounted(async () => {
  appVersion.value = await window.mainApi?.sendSync('msgGetAppVersion')

  const result = await window.mainApi?.sendSync('msgGetUserLicense')
  console.log('[RBAC] Fetching role for license:', result?.data?.id)
  GroupPricingService.getRoleForUser(result?.data?.id)
    .then((response: any) => {
      console.log(
        '[RBAC] getRoleForUser response:',
        JSON.stringify(response.data)
      )
      const role = response.data?.data ?? response.data
      if (role && role.permissions && Array.isArray(role.permissions)) {
        const permMap: Record<string, boolean> = {}
        for (const perm of role.permissions) {
          permMap[perm.slug] = true
        }
        console.log('[RBAC] Permission map:', JSON.stringify(permMap))
        permissionsStore.setPermissions(permMap)
      } else {
        console.log('[RBAC] No permissions found in role response')
        permissionsStore.markLoaded()
      }
    })
    .catch((error: any) => {
      console.error('[RBAC] Error fetching user role:', error)
      permissionsStore.markLoaded()
    })
})

const internalDrawer = ref(navProps.drawer)
const confirmationDialog = ref()

const handleLogout = async () => {
  const confirmed = await confirmationDialog.value?.open(
    'Logout',
    'Are you sure you want to logout?'
  )
  if (confirmed) {
    try {
      appStore.clearAll()
      permissionsStore.clearPermissions()
      await window.mainApi?.sendSync('msgLogout')
    } catch (error) {
      console.error('Error during logout:', error)
    }
  }
}

watchEffect(() => {
  internalDrawer.value = navProps.drawer
})
</script>

<style scoped>
.custom-navbar {
  font-size: 10px !important;
}

.nav-text {
  color: white !important;
  font-size: 12px !important;
}

.nav-bg {
  background-color: #2e566e !important;
}

.smaller-font :deep(.v-list-item-title),
.smaller-font :deep(.v-list-item__append) {
  font-size: 14px;
  padding-left: 0 !important;
}
.smaller-font.v-list-item {
  min-height: unset;
}
.first-level-group :deep(.v-list-group__items) {
  padding-left: 0 !important;
  --indent-padding: calc(var(--parent-padding) - 16px) !important;
}

.v-list-item.v-list-item--active {
  border-top-right-radius: 32px !important;
  border-bottom-right-radius: 32px !important;
}

.v-list-group__items .v-list-item {
  padding-inline-start: calc(var(--indent-padding)) !important;
}

.disabled-item {
  pointer-events: auto;
  opacity: 0.5;
  cursor: not-allowed;
}

.v-list-group .disabled-item:hover {
  background-color: transparent !important;
}

.logout-item {
  margin-top: 8px;
}

.logout-item:hover {
  background-color: rgba(255, 255, 255, 0.1) !important;
}
</style>
