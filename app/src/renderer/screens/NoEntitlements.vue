<template>
  <v-container>
    <v-row justify="center" class="mt-8">
      <v-col cols="12" sm="8" md="6">
        <v-card class="pa-6 text-center" elevation="4" rounded="lg">
          <v-icon size="64" color="warning" class="mb-4"
            >mdi-shield-alert-outline</v-icon
          >
          <h2 class="text-h5 mb-2">No Entitlements Found</h2>
          <p class="text-body-2 text-medium-emphasis mb-2">
            Your license does not have any feature entitlements assigned to it.
          </p>
          <p class="text-body-2 text-medium-emphasis mb-6">
            Please contact your administrator to assign the required
            entitlements to your license, then retry.
          </p>
          <v-btn
            color="primary"
            rounded
            :loading="retrying"
            class="mr-2"
            @click="retryEntitlements"
          >
            Retry
          </v-btn>
          <v-btn variant="outlined" rounded @click="goToSettings">
            Settings
          </v-btn>
          <v-alert
            v-if="errorMessage"
            type="error"
            variant="tonal"
            density="compact"
            class="mt-4 text-left"
          >
            {{ errorMessage }}
          </v-alert>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/renderer/store/app'

const router = useRouter()
const appStore = useAppStore()
const retrying = ref(false)
const errorMessage = ref('')
const licenseUrl = import.meta.env.VITE_APP_LICENSE_SERVER

const retryEntitlements = async () => {
  retrying.value = true
  errorMessage.value = ''

  try {
    const license = window.mainApi?.sendSync('msgGetUserLicense')
    if (!license?.data?.id) {
      errorMessage.value =
        'No license data found. Please re-activate your license.'
      return
    }

    const licenseKey = license.data.attributes?.key
    const response = await fetch(
      licenseUrl + '/licenses/' + license.data.id + '/get-entitlements',
      {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          ...(licenseKey ? { Authorization: 'Bearer ' + licenseKey } : {})
        }
      }
    )

    const rs = await response.json()
    const entitlementList: string[] = []
    if (rs?.data?.length > 0) {
      rs.data.forEach((entitlement: any) => {
        entitlementList.push(entitlement.name)
      })
    }

    if (entitlementList.length > 0) {
      appStore.setEntitlements(entitlementList)
      await router.push({ name: 'dashboard' })
    } else {
      errorMessage.value =
        'Still no entitlements found. Please contact your administrator.'
    }
  } catch (error) {
    errorMessage.value =
      'Unable to reach the license server. Please check your connection and try again.'
  } finally {
    retrying.value = false
  }
}

const goToSettings = () => {
  router.push({ name: 'app-settings' })
}
</script>

<style scoped></style>
