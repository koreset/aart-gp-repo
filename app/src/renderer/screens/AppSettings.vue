<template>
  <v-container>
    <v-row>
      <v-col>
        <base-card>
          <template #header>
            <span class="headline"> Application Configuration </span>
          </template>
          <template #default>
            <v-expansion-panels v-model="panels" class="mt-3" multiple>
              <v-expansion-panel>
                <v-expansion-panel-title> Licenses </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <v-row v-if="license !== null">
                    <v-col>
                      <p>
                        Licensed to
                        {{
                          license?.data?.attributes?.metadata?.userName ||
                          license?.attributes?.metadata?.userName ||
                          'N/A'
                        }}
                        ({{
                          license?.data?.attributes?.metadata?.userEmail ||
                          license?.attributes?.metadata?.userEmail ||
                          'N/A'
                        }}) [{{ maskedKey }}]
                      </p>
                    </v-col>
                  </v-row>
                  <v-row>
                    <v-col cols="4">
                      <v-btn
                        class="primary"
                        variant="outlined"
                        size="small"
                        rounded
                        @click="showLicenseDialog = true"
                        >Change License</v-btn
                      >
                    </v-col>
                  </v-row>
                  <v-dialog v-model="showLicenseDialog" max-width="500">
                    <v-card>
                      <v-card-title>Change License</v-card-title>
                      <v-card-text>
                        <v-textarea
                          v-model="newLicenseKey"
                          variant="outlined"
                          density="compact"
                          label="Paste new license key here"
                          rows="1"
                          auto-grow
                        ></v-textarea>
                      </v-card-text>
                      <v-card-actions>
                        <v-spacer></v-spacer>
                        <v-btn color="primary" @click="submitNewLicense"
                          >Submit</v-btn
                        >
                        <v-btn color="grey" @click="showLicenseDialog = false"
                          >Cancel</v-btn
                        >
                      </v-card-actions>
                    </v-card>
                  </v-dialog>
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  Activated Machines
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <v-row v-if="loadingMachines" justify="center">
                    <v-col cols="auto">
                      <v-progress-circular
                        indeterminate
                        size="24"
                        class="mr-2"
                      />
                      Loading machines...
                    </v-col>
                  </v-row>
                  <v-row v-else-if="machines.length === 0">
                    <v-col>
                      <p class="text-medium-emphasis"
                        >No machines activated for this license.</p
                      >
                    </v-col>
                  </v-row>
                  <template v-else>
                    <v-table density="compact">
                      <thead>
                        <tr>
                          <th>Fingerprint</th>
                          <th>Name</th>
                          <th>Platform</th>
                          <th>Created</th>
                          <th>Actions</th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr v-for="m in machines" :key="m.id">
                          <td class="text-caption"
                            >{{
                              m.attributes.fingerprint?.substring(0, 16)
                            }}...</td
                          >
                          <td>{{ m.attributes.name || '-' }}</td>
                          <td>{{ m.attributes.platform || '-' }}</td>
                          <td>{{
                            new Date(m.attributes.created).toLocaleDateString()
                          }}</td>
                          <td>
                            <v-btn
                              size="x-small"
                              variant="text"
                              color="error"
                              :loading="deactivatingMachine === m.id"
                              @click="deactivateMachine(m)"
                            >
                              Deactivate
                            </v-btn>
                          </td>
                        </tr>
                      </tbody>
                    </v-table>
                  </template>
                  <v-row class="mt-2">
                    <v-col>
                      <v-btn
                        variant="outlined"
                        size="small"
                        rounded
                        :loading="loadingMachines"
                        @click="fetchMachines"
                      >
                        Refresh
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-expansion-panel-text>
              </v-expansion-panel>

              <v-expansion-panel>
                <v-expansion-panel-title>
                  API Server Configuration
                </v-expansion-panel-title>
                <v-expansion-panel-text>
                  <v-row>
                    <v-col cols="5">
                      <v-text-field
                        v-model="apiServerUrl"
                        variant="outlined"
                        density="compact"
                        :error-messages="apiServerUrlErrors"
                        placeholder="Enter a valid url for the API Service"
                        label="API Server Url"
                        :messages="baseApiHint"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="4">
                      <v-btn
                        class="primary"
                        variant="outlined"
                        size="small"
                        rounded
                        @click="updateApiUrl"
                        >Save</v-btn
                      ></v-col
                    >
                  </v-row>
                </v-expansion-panel-text>
              </v-expansion-panel>
            </v-expansion-panels>
          </template>
        </base-card>
      </v-col>
    </v-row>
    <v-snackbar v-model="snackbar" :timeout="timeout">
      {{ snackbarMessage }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>

    <v-dialog v-model="showRestartDialog" max-width="400" persistent>
      <v-card>
        <v-card-title>Restart Required</v-card-title>
        <v-card-text>
          {{ restartReason }}
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" @click="confirmRestart">Restart Now</v-btn>
          <v-btn color="grey" @click="showRestartDialog = false">Later</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script setup lang="ts">
import BaseCard from '../components/BaseCard.vue'
import { computed, onMounted, ref } from 'vue'

// data
const panels = ref(0)
const license: any = ref(null)
const apiServerUrl: any = ref('')
const apiServerUrlErrors = ref([])
const baseApiHint = ref('')
const snackbar = ref(false)
const snackbarMessage = ref('')
const timeout = ref(6000)
const showLicenseDialog = ref(false)
const newLicenseKey = ref('')
const keyErrors = ref<Array<string>>([])
const successMessages = ref<Array<string>>([])
const validFlow = ref(false)
const showRestartDialog = ref(false)
const restartReason = ref('')
const machines = ref<any[]>([])
const loadingMachines = ref(false)
const deactivatingMachine = ref<string | null>(null)

const licenseUrl = import.meta.env.VITE_APP_LICENSE_SERVER

const maskedKey = computed(() => {
  const key =
    license.value?.data?.attributes?.key ||
    license.value?.attributes?.key ||
    license.value?.key
  if (!key) return 'N/A'
  const parts = key.split('-')
  if (parts.length <= 2) return key
  return (
    parts[0] +
    '-' +
    parts
      .slice(1, -1)
      .map(() => '****')
      .join('-') +
    '-' +
    parts[parts.length - 1]
  )
})

const confirmRestart = () => {
  showRestartDialog.value = false
  window.mainApi?.send('msgRestartApplication')
}

const fetchMachines = async () => {
  const licenseId = license.value?.data?.id || license.value?.id
  if (!licenseId) return
  loadingMachines.value = true
  try {
    const licenseKey =
      license.value?.data?.attributes?.key ||
      license.value?.attributes?.key ||
      license.value?.key
    const response = await fetch(
      licenseUrl + '/licenses/' + licenseId + '/machines',
      {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          Accept: 'application/json',
          Authorization: 'Bearer ' + licenseKey
        }
      }
    )
    const rs = await response.json()
    machines.value = rs.data || []
  } catch (error) {
    console.error('Failed to fetch machines:', error)
    snackbarMessage.value = 'Failed to load machines.'
    snackbar.value = true
  } finally {
    loadingMachines.value = false
  }
}

const deactivateMachine = async (machine: any) => {
  deactivatingMachine.value = machine.id
  try {
    const licenseKey =
      license.value?.data?.attributes?.key ||
      license.value?.attributes?.key ||
      license.value?.key
    const response = await fetch(licenseUrl + '/machines/' + machine.id, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
        Authorization: 'Bearer ' + licenseKey
      }
    })
    if (response.ok) {
      machines.value = machines.value.filter((m: any) => m.id !== machine.id)
      snackbarMessage.value = 'Machine deactivated successfully.'
      snackbar.value = true
    } else {
      const rs = await response.json()
      snackbarMessage.value = rs.message || 'Failed to deactivate machine.'
      snackbar.value = true
    }
  } catch (error) {
    snackbarMessage.value = 'Failed to deactivate machine.'
    snackbar.value = true
  } finally {
    deactivatingMachine.value = null
  }
}

onMounted(async () => {
  apiServerUrl.value = await window.mainApi?.sendSync('msgGetBaseUrl')
  license.value = await window.mainApi?.sendSync('msgGetUserLicense')
  await fetchMachines()
})

// const getLicenseDetails = () => {
//   showLicenseDialog.value = true
// }

const submitNewLicense = () => {
  if (!newLicenseKey.value.trim()) {
    snackbarMessage.value = 'Please enter a license key.'
    snackbar.value = true
    return
  }

  const result = window.mainApi?.sendSync(
    'msgActivateLicense',
    newLicenseKey.value.trim()
  )

  if (!result) {
    keyErrors.value.push(
      'Unable to reach the license server. Please check your connection and try again.'
    )
  } else if (result.valid) {
    window.mainApi?.send('msgSetUserLicense', result.license, result.machine)
    successMessages.value = ['License activated successfully']
    validFlow.value = true
    showLicenseDialog.value = false
    restartReason.value =
      'Your license has been updated. The application needs to restart to apply the changes.'
    showRestartDialog.value = true
  } else {
    switch (result.status) {
      case 'SUSPENDED':
        keyErrors.value.push(
          'This license has been suspended. Please provide a valid license'
        )
        break
      case 'EXPIRED':
        keyErrors.value.push(
          'This license has expired. Please provide a valid license'
        )
        break
      case 'NOT_FOUND':
        keyErrors.value.push(
          'This license does not exist. Please contact your admin for a valid license'
        )
        break
      case 'FINGERPRINT_SCOPE_MISMATCH':
        keyErrors.value.push(
          'This license is already associated with another machine. Please provide a valid license'
        )
        break
      case 'TOO_MANY_MACHINES':
        keyErrors.value.push(
          'This license has reached its maximum number of activated machines'
        )
        break
      default:
        keyErrors.value.push(
          'This license cannot be registered successfully. Please provide a valid license'
        )
    }
  }
  if (keyErrors.value.length > 0) {
    snackbarMessage.value = keyErrors.value.join(' ')
    snackbar.value = true
  } else {
    snackbarMessage.value = successMessages.value.join(' ')
    snackbar.value = true
  }
}

const updateApiUrl = () => {
  window.mainApi?.send('msgSaveBaseUrl', apiServerUrl.value)
  restartReason.value =
    'The API server URL has been updated. The application needs to restart to apply the changes.'
  showRestartDialog.value = true
}
</script>

<style></style>
