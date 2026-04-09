<template>
  <div>
    <v-card>
      <v-card-text>
        <v-stepper v-model="step" flat alt-labels>
          <v-stepper-header>
            <v-stepper-item
              :value="1"
              title="License Key"
              :complete="step > 1"
            />
            <v-divider />
            <v-stepper-item
              :value="2"
              title="API Server"
              :complete="step > 2"
            />
            <v-divider />
            <v-stepper-item :value="3" title="Complete" :complete="step > 3" />
          </v-stepper-header>

          <v-stepper-window>
            <!-- Step 1: License Activation -->
            <v-stepper-window-item :value="1">
              <v-container>
                <v-row>
                  <v-col>
                    <p class="text-body-2 mb-4 text-medium-emphasis">
                      Enter the license key that was provided to you via email
                      or from your company administrator.
                    </p>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="8">
                    <v-text-field
                      v-model="license"
                      variant="outlined"
                      density="compact"
                      label="License key"
                      placeholder="XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"
                      :maxlength="39"
                      :error-messages="keyErrors"
                      :disabled="activating"
                      @update:modelValue="addDashes"
                    />
                  </v-col>
                  <v-col cols="4">
                    <v-btn
                      color="primary"
                      rounded
                      size="small"
                      :loading="activating"
                      :disabled="!license || license.length < 39"
                      @click="activateLicense"
                    >
                      Activate license
                    </v-btn>
                  </v-col>
                </v-row>
              </v-container>
            </v-stepper-window-item>

            <!-- Step 2: API Server URL -->
            <v-stepper-window-item :value="2">
              <v-container>
                <v-row>
                  <v-col>
                    <v-alert
                      type="success"
                      variant="tonal"
                      density="compact"
                      class="mb-4"
                    >
                      License activated successfully.
                    </v-alert>
                    <p class="text-body-2 mb-4 text-medium-emphasis">
                      Enter the URL of the AART API server your organization
                      uses.
                    </p>
                  </v-col>
                </v-row>
                <v-row>
                  <v-col cols="8">
                    <v-text-field
                      v-model="baseApi"
                      label="API Server URL"
                      variant="outlined"
                      density="compact"
                      placeholder="https://api.example.com/api/v1/"
                    />
                  </v-col>
                  <v-col cols="4">
                    <v-btn
                      color="primary"
                      rounded
                      size="small"
                      :disabled="!baseApi"
                      @click="saveApiUrl"
                    >
                      Save & Continue
                    </v-btn>
                  </v-col>
                </v-row>
              </v-container>
            </v-stepper-window-item>

            <!-- Step 3: Complete -->
            <v-stepper-window-item :value="3">
              <v-container class="text-center py-8">
                <v-icon size="64" color="success" class="mb-4"
                  >mdi-check-circle-outline</v-icon
                >
                <h3 class="text-h6 mb-2">Setup Complete</h3>
                <p class="text-body-2 text-medium-emphasis mb-6">
                  Your license has been activated and the API server configured.
                  The application will restart to apply your settings.
                </p>
                <v-btn color="primary" rounded @click="completeSetup">
                  Launch Application
                </v-btn>
              </v-container>
            </v-stepper-window-item>
          </v-stepper-window>
        </v-stepper>
      </v-card-text>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const licenseUrl = import.meta.env.VITE_APP_LICENSE_SERVER
const fingerPrint = window.mainApi?.sendSync('msgGetMachineFingerprint')

const step = ref(1)
const license = ref('')
const baseApi = ref('')
const keyErrors = ref<string[]>([])
const activating = ref(false)

const addDashes = () => {
  license.value = license.value.toUpperCase()
  const positions = [4, 9, 14, 19, 24, 29, 34]
  if (positions.includes(license.value.length)) {
    license.value += '-'
  }
}

const activateLicense = async () => {
  activating.value = true
  keyErrors.value = []

  try {
    const validation = await fetch(licenseUrl + '/activate-key', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json'
      },
      body: JSON.stringify({
        key: license.value,
        fingerprint: fingerPrint
      })
    })

    const rs = await validation.json()

    if (rs.valid) {
      window.mainApi?.send('msgSetUserLicense', rs.license, rs.machine)
      step.value = 2
    } else {
      switch (rs.status) {
        case 'SUSPENDED':
          keyErrors.value.push(
            'This license has been suspended. Please provide a valid license.'
          )
          break
        case 'EXPIRED':
          keyErrors.value.push(
            'This license has expired. Please provide a valid license.'
          )
          break
        case 'NOT_FOUND':
          keyErrors.value.push(
            'This license does not exist. Please contact your admin for a valid license.'
          )
          break
        case 'FINGERPRINT_SCOPE_MISMATCH':
          keyErrors.value.push(
            'This license cannot be registered on this machine. Please provide a valid license.'
          )
          break
        case 'TOO_MANY_MACHINES':
          keyErrors.value.push(
            'This license has reached its maximum number of activated machines.'
          )
          break
        default:
          keyErrors.value.push(
            'This license cannot be registered successfully. Please provide a valid license.'
          )
      }
    }
  } catch (error) {
    keyErrors.value.push(
      'Unable to reach the license server. Please check your connection and try again.'
    )
  } finally {
    activating.value = false
  }
}

const saveApiUrl = () => {
  window.mainApi?.send('msgSaveBaseUrl', baseApi.value)
  step.value = 3
}

const completeSetup = () => {
  window.mainApi?.send('msgRestartApplication')
}
</script>

<style scoped></style>
