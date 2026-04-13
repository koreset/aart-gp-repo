<template>
  <v-container fluid>
    <v-row class="top-banner primary">
      <v-col cols="3">
        <img
          class="ml-3 mb-3 mt-0"
          width="100%"
          :src="'./images/aart-logo-02.png'"
        />
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-alert
          v-if="statusMessage"
          :type="statusType"
          variant="tonal"
          class="mx-4 mt-4"
          prominent
        >
          <div class="text-body-1 font-weight-medium">{{ statusTitle }}</div>
          <div class="text-body-2">{{ statusMessage }}</div>
        </v-alert>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-form>
          <v-container>
            <v-row>
              <v-col><h1>License Registration</h1></v-col>
            </v-row>
            <v-row>
              <v-col cols="8" class="d-flex flex-column">
                <v-text-field
                  v-model:model-value="license"
                  variant="outlined"
                  density="compact"
                  label="Enter license"
                  placeholder="XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX-XXXX"
                  :maxlength="39"
                  :error-messages="keyErrors"
                  :messages="successMessages"
                  :disabled="activating"
                  @update:modelValue="addDashes"
                ></v-text-field>
              </v-col>
              <v-col cols="4">
                <v-btn
                  class="mr-3"
                  color="primary"
                  rounded
                  size="small"
                  :loading="activating"
                  :disabled="!license || license.length < 39"
                  @click="activateLicense"
                >
                  Activate
                </v-btn>
                <v-btn
                  v-if="validFlow"
                  color="primary"
                  rounded
                  size="small"
                  @click="restartApp"
                >
                  Restart App
                </v-btn>
              </v-col>
            </v-row>
          </v-container>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
import { computed, onBeforeMount, ref } from 'vue'

const props = defineProps<{
  licenseStatus?: string
}>()

const keyErrors = ref<string[]>([])
const successMessages = ref<string[]>([])
const validFlow = ref(false)
const license = ref('')
const activating = ref(false)

const statusTitle = computed(() => {
  switch (props.licenseStatus) {
    case 'EXPIRED':
      return 'License Expired'
    case 'SUSPENDED':
      return 'License Suspended'
    case 'OVERDUE':
      return 'License Check-In Overdue'
    default:
      return 'License Issue'
  }
})

const statusMessage = computed(() => {
  switch (props.licenseStatus) {
    case 'EXPIRED':
      return 'Your license has expired. Please enter a new license key or contact your administrator to renew.'
    case 'SUSPENDED':
      return 'Your license has been suspended. Please contact your administrator for assistance.'
    case 'OVERDUE':
      return 'Your license is overdue for a check-in. Please re-register to continue using the application.'
    case 'INVALID':
      return 'Your license could not be validated. Please enter a valid license key.'
    default:
      return ''
  }
})

const statusType = computed(() => {
  switch (props.licenseStatus) {
    case 'EXPIRED':
      return 'warning'
    case 'SUSPENDED':
      return 'error'
    case 'OVERDUE':
      return 'warning'
    default:
      return 'info'
  }
})

const restartApp = () => {
  window.mainApi?.send('msgRestartApplication')
}

const addDashes = () => {
  license.value = license.value.toUpperCase()
  const positions = [4, 9, 14, 19, 24, 29, 34]
  if (positions.includes(license.value.length)) {
    license.value += '-'
  }
}

onBeforeMount(() => {
  window.mainApi?.send('msgResizeWindow', 1024, 600, false)
})

const activateLicense = async () => {
  keyErrors.value = []
  successMessages.value = []
  activating.value = true

  try {
    const result = window.mainApi?.sendSync('msgActivateLicense', license.value)

    if (!result) {
      keyErrors.value.push(
        'Unable to reach the license server. Please check your connection and try again.'
      )
      return
    }

    if (result.valid) {
      window.mainApi?.send('msgSetUserLicense', result.license, result.machine)
      successMessages.value = ['License activated successfully']
      validFlow.value = true
    } else {
      switch (result.status) {
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
</script>

<style scoped>
.input-info {
  justify-content: flex-end !important;
}
.v-input input {
  text-align: center !important;
}

.top-banner {
  height: 100px;
  background: rgba(40, 78, 103, 0.7);
}
</style>
