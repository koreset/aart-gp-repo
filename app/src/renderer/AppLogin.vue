<template>
  <v-app>
    <!-- Background with gradient -->
    <div class="login-background">
      <v-container fluid class="fill-height">
        <v-row justify="center" align="center" class="fill-height">
          <!-- Logo Section -->
          <v-col
            cols="12"
            sm="6"
            md="5"
            lg="4"
            class="text-center mb-8 mb-md-0"
          >
            <div class="logo-container">
              <img
                class="logo-image"
                src="/images/aart-logo-02.png"
                alt="AART Logo"
                @load="logoLoaded = true"
              />
              <h1 class="app-title mt-6">AART Enterprise</h1>
              <p class="app-subtitle">Actuarial, Analytics & Reporting Tool</p>
            </div>
          </v-col>

          <!-- Login Form Section -->
          <v-col cols="12" sm="6" md="5" lg="4">
            <v-card
              class="login-card elevation-12"
              :class="{ 'card-loading': isLoading }"
            >
              <!-- Login Form -->
              <v-card-text v-if="!forgotPasswordFlow" class="pa-8">
                <div class="text-center mb-6">
                  <h2 class="login-title">Welcome Back</h2>
                  <p class="login-subtitle">Sign in to your account</p>
                </div>

                <!-- Identity Provider Selection -->
                <div
                  v-if="showProviderSelection && availableProviders.length > 1"
                  class="mb-6"
                >
                  <v-select
                    v-model="selectedProvider"
                    label="Select Identity Provider"
                    variant="outlined"
                    :items="availableProviders"
                    item-title="displayName"
                    item-value="id"
                    prepend-inner-icon="mdi-domain"
                    class="mb-4"
                    :disabled="isLoading"
                    @update:model-value="onProviderChange"
                  >
                    <template #item="{ props, item }">
                      <v-list-item
                        v-bind="props"
                        :prepend-icon="getProviderIcon(item.raw.type)"
                      >
                        <v-list-item-title>{{
                          item.raw.displayName
                        }}</v-list-item-title>
                        <v-list-item-subtitle>{{
                          item.raw.description
                        }}</v-list-item-subtitle>
                      </v-list-item>
                    </template>
                  </v-select>
                </div>

                <v-form
                  ref="loginForm"
                  v-model="isFormValid"
                  :disabled="isLoading"
                  @submit.prevent="login"
                >
                  <!-- Credential fields - only show for providers that require them -->
                  <template v-if="requiresCredentials">
                    <v-text-field
                      v-model="username"
                      label="Username or Email"
                      variant="outlined"
                      prepend-inner-icon="mdi-account-outline"
                      :rules="usernameRules"
                      :error-messages="fieldErrors.username"
                      class="mb-4"
                      autocomplete="username"
                      :disabled="isLoading"
                      @keyup.enter="login"
                      @input="clearFieldError('username')"
                    ></v-text-field>

                    <v-text-field
                      v-model="password"
                      label="Password"
                      variant="outlined"
                      prepend-inner-icon="mdi-lock-outline"
                      :append-inner-icon="
                        showPassword ? 'mdi-eye' : 'mdi-eye-off'
                      "
                      :type="showPassword ? 'text' : 'password'"
                      :rules="passwordRules"
                      :error-messages="fieldErrors.password"
                      class="mb-6"
                      autocomplete="current-password"
                      :disabled="isLoading"
                      @click:append-inner="showPassword = !showPassword"
                      @keyup.enter="login"
                      @input="clearFieldError('password')"
                    ></v-text-field>
                  </template>

                  <!-- SSO Information - show for external providers -->
                  <template v-else>
                    <v-card variant="outlined" class="mb-6 pa-4">
                      <div class="d-flex align-center">
                        <v-icon
                          :icon="
                            getProviderIcon(selectedProviderInfo?.type || '')
                          "
                          size="24"
                          class="mr-3"
                        ></v-icon>
                        <div>
                          <div class="text-body-1 font-weight-medium">{{
                            selectedProviderInfo?.displayName
                          }}</div>
                          <div class="text-body-2 text-medium-emphasis">{{
                            selectedProviderInfo?.description
                          }}</div>
                        </div>
                      </div>
                    </v-card>
                  </template>

                  <!-- Standard Login Button for credential-based providers -->
                  <v-btn
                    v-if="requiresCredentials"
                    type="submit"
                    block
                    size="large"
                    color="primary"
                    variant="elevated"
                    class="mb-4 login-btn"
                    :loading="isLoading"
                    :disabled="!isFormValid || isLoading"
                  >
                    <v-icon start>mdi-login</v-icon>
                    Sign In
                  </v-btn>

                  <!-- SSO Login Button for external providers -->
                  <v-btn
                    v-else
                    block
                    size="large"
                    color="primary"
                    variant="elevated"
                    class="mb-4 login-btn"
                    :loading="isLoading"
                    :disabled="isLoading"
                    @click="loginWithSSO"
                  >
                    <v-icon start>{{
                      getProviderIcon(selectedProviderInfo?.type || '')
                    }}</v-icon>
                    Sign in with {{ selectedProviderInfo?.displayName }}
                  </v-btn>

                  <div class="text-center">
                    <v-btn
                      variant="text"
                      color="primary"
                      size="small"
                      :disabled="isLoading"
                      @click="showForgotPassword"
                    >
                      Forgot your password?
                    </v-btn>
                  </div>
                </v-form>
              </v-card-text>

              <!-- Forgot Password Form -->
              <v-card-text v-else class="pa-8">
                <div class="text-center mb-6">
                  <v-icon size="48" color="primary" class="mb-4"
                    >mdi-email-outline</v-icon
                  >
                  <h2 class="login-title">Reset Password</h2>
                  <p class="login-subtitle"
                    >Enter your email to receive a reset link</p
                  >
                </div>

                <v-form
                  ref="forgotForm"
                  v-model="isForgotFormValid"
                  :disabled="isLoading"
                  @submit.prevent="doForgotPassword"
                >
                  <v-text-field
                    v-model="userEmail"
                    label="Email Address"
                    variant="outlined"
                    prepend-inner-icon="mdi-email-outline"
                    :rules="emailRules"
                    :error-messages="fieldErrors.email"
                    class="mb-6"
                    autocomplete="email"
                    :disabled="isLoading"
                    @keyup.enter="doForgotPassword"
                    @input="clearFieldError('email')"
                  ></v-text-field>

                  <v-btn
                    type="submit"
                    block
                    size="large"
                    color="primary"
                    variant="elevated"
                    class="mb-4"
                    :loading="isLoading"
                    :disabled="!isForgotFormValid || isLoading"
                  >
                    <v-icon start>mdi-email-send</v-icon>
                    Send Reset Link
                  </v-btn>

                  <v-btn
                    variant="text"
                    color="primary"
                    size="small"
                    block
                    :disabled="isLoading"
                    @click="showLogin"
                  >
                    <v-icon start>mdi-arrow-left</v-icon>
                    Back to Login
                  </v-btn>
                </v-form>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </div>

    <!-- Snackbar for notifications -->
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :timeout="snackbar.timeout"
      location="top"
    >
      {{ snackbar.message }}
      <template #actions>
        <v-btn color="white" variant="text" @click="snackbar.show = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script setup lang="ts">
import {
  onBeforeMount,
  onMounted,
  ref,
  reactive,
  nextTick,
  computed
} from 'vue'
import { AuthService } from '@/renderer/auth/AuthService'
import type { IAuthProvider } from '@/renderer/types/auth'

// Component state
const authUrl = import.meta.env.VITE_APP_AUTH_URL
const authService = AuthService.getInstance()
const forgotPasswordFlow = ref(false)
const availableProviders = ref<IAuthProvider[]>([])
const selectedProvider = ref<string | null>(null)
const showProviderSelection = ref(false)
const isLoading = ref(false)
const isFormValid = ref(false)
const isForgotFormValid = ref(false)
const showPassword = ref(false)
const logoLoaded = ref(false)

// Form refs
const loginForm = ref()
const forgotForm = ref()

// Form data
const username = ref('')
const password = ref('')
const userEmail = ref('')

// Error handling
const fieldErrors = reactive({
  username: [] as string[],
  password: [] as string[],
  email: [] as string[]
})

// Snackbar for notifications
const snackbar = reactive({
  show: false,
  message: '',
  color: 'info',
  timeout: 5000
})

// Validation rules
const usernameRules = [
  (v: string) => !!v || 'Username is required',
  (v: string) => v.length >= 3 || 'Username must be at least 3 characters'
]

const passwordRules = [
  (v: string) => !!v || 'Password is required',
  (v: string) => v.length >= 6 || 'Password must be at least 6 characters'
]

const emailRules = [
  (v: string) => !!v || 'Email is required',
  (v: string) => /.+@.+\..+/.test(v) || 'Please enter a valid email address'
]

// Computed properties
const requiresCredentials = computed(() => {
  if (!selectedProvider.value) return true
  const provider = availableProviders.value.find(
    (p) => p.id === selectedProvider.value
  )
  return provider?.type === 'internal' || provider?.type === 'ldap'
})

const selectedProviderInfo = computed(() => {
  return availableProviders.value.find((p) => p.id === selectedProvider.value)
})

// Lifecycle hooks
onBeforeMount(() => {
  window.mainApi?.send('msgResizeWindow', 1200, 700, false)
})

onMounted(async () => {
  // Initialize authentication service
  await initializeAuth()

  // Focus on username field after component mounts
  nextTick(() => {
    const usernameField = document.querySelector(
      'input[autocomplete="username"]'
    ) as HTMLInputElement
    if (usernameField) {
      usernameField.focus()
    }
  })
})

// Utility functions
const showSnackbar = (message: string, color: string = 'info') => {
  snackbar.message = message
  snackbar.color = color
  snackbar.show = true
}

const clearFieldError = (field: string) => {
  if (fieldErrors[field]) {
    fieldErrors[field] = []
  }
}

const clearAllErrors = () => {
  Object.keys(fieldErrors).forEach((field) => {
    fieldErrors[field] = []
  })
}

// Navigation functions
const showForgotPassword = () => {
  clearAllErrors()
  forgotPasswordFlow.value = true
  nextTick(() => {
    const emailField = document.querySelector(
      'input[autocomplete="email"]'
    ) as HTMLInputElement
    if (emailField) {
      emailField.focus()
    }
  })
}

const showLogin = () => {
  clearAllErrors()
  forgotPasswordFlow.value = false
  userEmail.value = ''
  nextTick(() => {
    const usernameField = document.querySelector(
      'input[autocomplete="username"]'
    ) as HTMLInputElement
    if (usernameField) {
      usernameField.focus()
    }
  })
}

// Authentication initialization
const initializeAuth = async () => {
  try {
    // Initialize auth service with client detection
    await authService.initialize()

    // Load available providers
    availableProviders.value = authService.getAvailableProviders()
    showProviderSelection.value =
      authService.allowsProviderSelection() &&
      availableProviders.value.length > 1

    // Set default provider
    if (availableProviders.value.length > 0) {
      const config = authService.getConfiguration()
      selectedProvider.value =
        config?.defaultProvider || availableProviders.value[0].id
    }
  } catch (error) {
    console.error('Failed to initialize authentication:', error)
    showSnackbar('Failed to load authentication configuration', 'error')
  }
}

// Provider management
const onProviderChange = () => {
  clearAllErrors()
  username.value = ''
  password.value = ''
}

const getProviderIcon = (type: string) => {
  const icons = {
    saml: 'mdi-shield-key',
    oauth2: 'mdi-account-key',
    oidc: 'mdi-openid',
    ldap: 'mdi-server-network',
    'azure-ad': 'mdi-microsoft',
    google: 'mdi-google',
    okta: 'mdi-shield-account',
    internal: 'mdi-account'
  }
  return icons[type] || 'mdi-domain'
}

// SSO Authentication
const loginWithSSO = async () => {
  if (!selectedProvider.value) return

  isLoading.value = true
  clearAllErrors()

  try {
    const result = await authService.login(selectedProvider.value)

    if (result.success) {
      showSnackbar('Login successful! Redirecting...', 'success')

      // Store authenticated user
      const storeResult = await window.mainApi?.sendSync(
        'msgStoreAuthenticatedUser',
        result.token
      )
      if (storeResult) {
        setTimeout(() => {
          window.mainApi?.send('msgRestartApplication')
        }, 1000)
      }
    } else if (result.error === 'SSO_REDIRECT') {
      showSnackbar('Redirecting to external authentication provider...', 'info')
      // SSO redirect initiated, user will authenticate externally
    } else {
      showSnackbar(result.error || 'SSO authentication failed', 'error')
    }
  } catch (error) {
    console.error('SSO login error:', error)
    showSnackbar(
      'Network error. Please check your connection and try again.',
      'error'
    )
  } finally {
    isLoading.value = false
  }
}

// Authentication functions
const login = async () => {
  if (!loginForm.value?.validate()) {
    return
  }

  isLoading.value = true
  clearAllErrors()

  try {
    if (!selectedProvider.value) {
      showSnackbar('Please select an authentication provider', 'error')
      return
    }

    const credentials = {
      username: username.value,
      password: password.value,
      provider: selectedProvider.value
    }

    const result = await authService.login(selectedProvider.value, credentials)

    if (result.success) {
      showSnackbar('Login successful! Redirecting...', 'success')

      // Store authenticated user
      const storeResult = await window.mainApi?.sendSync(
        'msgStoreAuthenticatedUser',
        result.token
      )
      if (storeResult) {
        // Small delay to show success message
        setTimeout(() => {
          window.mainApi?.send('msgRestartApplication')
        }, 1000)
      }
    } else {
      fieldErrors.password = [result.error || 'Login failed']
      showSnackbar(
        result.error || 'Login failed. Please check your credentials.',
        'error'
      )
    }
  } catch (error) {
    console.error('Login error:', error)
    showSnackbar(
      'Network error. Please check your connection and try again.',
      'error'
    )
  } finally {
    isLoading.value = false
  }
}

const doForgotPassword = async () => {
  if (!forgotForm.value?.validate()) {
    return
  }

  isLoading.value = true
  clearAllErrors()

  try {
    const response = await fetch(authUrl + '/forgotPassword', {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ recipient: userEmail.value })
    })

    if (response.ok) {
      showSnackbar('Password reset link sent to your email!', 'success')
      // Auto-return to login after success
      setTimeout(() => {
        showLogin()
      }, 2000)
    } else {
      const data = await response.json()
      fieldErrors.email = [data.message || 'Failed to send reset email']
      showSnackbar('Failed to send reset email. Please try again.', 'error')
    }
  } catch (error) {
    console.error('Forgot password error:', error)
    showSnackbar(
      'Network error. Please check your connection and try again.',
      'error'
    )
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
/* Main background with gradient */
.login-background {
  min-height: 100vh;
  background: linear-gradient(135deg, #1a3a4f 0%, #2e566e 50%, #003f58 100%);
  position: relative;
  overflow: hidden;
}

.login-background::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background:
    radial-gradient(
      circle at 20% 20%,
      rgba(120, 119, 198, 0.3) 0%,
      transparent 50%
    ),
    radial-gradient(
      circle at 80% 80%,
      rgba(255, 255, 255, 0.15) 0%,
      transparent 50%
    ),
    radial-gradient(
      circle at 40% 40%,
      rgba(120, 119, 198, 0.2) 0%,
      transparent 50%
    );
  animation: backgroundShift 20s ease-in-out infinite;
}

@keyframes backgroundShift {
  0%,
  100% {
    transform: scale(1) rotate(0deg);
  }
  50% {
    transform: scale(1.1) rotate(5deg);
  }
}

/* Logo section */
.logo-container {
  padding: 2rem;
  text-align: center;
}

.logo-image {
  max-width: 280px;
  width: 100%;
  height: auto;
  filter: drop-shadow(0 8px 32px rgba(0, 0, 0, 0.1));
  animation: logoFloat 6s ease-in-out infinite;
}

@keyframes logoFloat {
  0%,
  100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-10px);
  }
}

.app-title {
  font-size: 2.5rem;
  font-weight: 300;
  color: white;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  margin-bottom: 0.5rem;
}

.app-subtitle {
  font-size: 1.1rem;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 300;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

/* Login card */
.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow:
    0 32px 64px rgba(0, 0, 0, 0.15),
    0 8px 32px rgba(0, 0, 0, 0.1);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  max-width: 420px;
  margin: 0 auto;
}

.login-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #003f58, #2e566e, #1a3a4f);
  background-size: 200% 100%;
  animation: gradientShift 3s ease infinite;
}

@keyframes gradientShift {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.login-card:hover {
  transform: translateY(-4px);
  box-shadow:
    0 40px 80px rgba(0, 0, 0, 0.2),
    0 16px 48px rgba(0, 0, 0, 0.15);
}

.card-loading {
  pointer-events: none;
  opacity: 0.9;
}

/* Typography */
.login-title {
  font-size: 1.75rem;
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.login-subtitle {
  font-size: 0.95rem;
  color: #64748b;
  font-weight: 400;
  line-height: 1.5;
}

/* Form elements */
.v-text-field {
  margin-bottom: 1rem;
}

:deep(.v-field--outlined) {
  border-radius: 12px;
}

:deep(.v-field--outlined .v-field__outline) {
  --v-field-border-opacity: 0.2;
}

:deep(.v-field--focused .v-field__outline) {
  --v-field-border-opacity: 1;
  --v-field-border-width: 2px;
}

:deep(.v-text-field .v-input__prepend-inner) {
  padding-top: 8px;
}

/* Login button */
.login-btn {
  border-radius: 12px !important;
  text-transform: none !important;
  font-weight: 600 !important;
  font-size: 1rem !important;
  height: 48px !important;
  background: linear-gradient(to right, #003f58, #2e566e) !important;
  box-shadow: 0 4px 12px rgba(0, 63, 88, 0.4) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

.login-btn:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 8px 20px rgba(0, 63, 88, 0.5) !important;
}

.login-btn:active {
  transform: translateY(0) !important;
}

/* Responsive design */
@media (max-width: 960px) {
  .logo-container {
    padding: 1rem;
    margin-bottom: 2rem;
  }

  .app-title {
    font-size: 2rem;
  }

  .app-subtitle {
    font-size: 1rem;
  }

  .login-card {
    margin: 1rem;
    border-radius: 16px;
  }
}

@media (max-width: 600px) {
  .login-background {
    padding: 1rem;
  }

  .logo-image {
    max-width: 200px;
  }

  .app-title {
    font-size: 1.75rem;
  }

  .login-card {
    margin: 0.5rem;
  }

  .login-title {
    font-size: 1.5rem;
  }
}

/* Loading states */
:deep(.v-btn--loading) {
  pointer-events: none;
}

/* Accessibility improvements */
:deep(.v-btn:focus-visible) {
  outline: 2px solid #2e566e;
  outline-offset: 2px;
}

:deep(.v-text-field:focus-within .v-field__outline) {
  --v-field-border-color: #2e566e;
}

/* Animation for form transitions */
.v-card-text {
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Error state styling */
:deep(.v-text-field--error .v-field__outline) {
  --v-field-border-color: #ef4444;
}

/* Success state for snackbar */
:deep(.v-snackbar--variant-elevated) {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}
</style>
