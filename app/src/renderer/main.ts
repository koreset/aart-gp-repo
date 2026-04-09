/* eslint-disable vue/one-component-per-file */
import { createApp, h, ref, onMounted, defineComponent } from 'vue'
import { createPinia } from 'pinia'

import App from '@/renderer/App.vue'
import AppLogin from '@/renderer/AppLogin.vue'
import AppSetup from '@/renderer/AppSetup.vue'
import AppNoInternet from '@/renderer/AppNoInternet.vue'
import LicenseWindow from '@/renderer/LicenseWindow.vue'
import router from '@/renderer/router'
import vuetify from '@/renderer/plugins/vuetify'
import i18n from '@/renderer/plugins/i18n'
import Vuelidate from 'vuelidate'
import { LicenseManager } from 'ag-grid-enterprise'
import upperFirst from 'lodash/upperFirst'
import camelCase from 'lodash/camelCase'

LicenseManager.setLicenseKey(
  'DownloadDevTools_COM_NDEwMjM0NTgwMDAwMA==59158b5225400879a12a96634544f5b6'
)

// Add API key defined in contextBridge to window object type
declare global {
  // eslint-disable-next-line no-unused-vars
  interface Window {
    mainApi?: any
    node?: any
    electronAPI?: {
      onNavigate: (callback: (routeName: string) => void) => void
    }
  }
}

// Splash screen component — shown while license checks run
const AppSplash = defineComponent({
  props: {
    message: { type: String, default: 'Starting application...' }
  },
  setup(props) {
    return () =>
      h(
        'div',
        {
          style: {
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            height: '100vh',
            background: 'linear-gradient(135deg, #003f58 0%, #2e566e 100%)',
            color: 'white',
            fontFamily: 'Roboto, sans-serif'
          }
        },
        [
          h('img', {
            src: '/images/aart-logo-02.png',
            style: { width: '200px', marginBottom: '32px' }
          }),
          h('div', {
            style: {
              width: '40px',
              height: '40px',
              border: '3px solid rgba(255,255,255,0.3)',
              borderTop: '3px solid white',
              borderRadius: '50%',
              animation: 'spin 1s linear infinite',
              marginBottom: '16px'
            }
          }),
          h(
            'p',
            { style: { fontSize: '14px', opacity: '0.8' } },
            props.message
          ),
          h(
            'style',
            {},
            '@keyframes spin { to { transform: rotate(360deg); } }'
          )
        ]
      )
  }
})

// save the license server url to the store
window.mainApi?.sendSync(
  'msgSetLicenseServerUrl',
  import.meta.env.VITE_APP_LICENSE_SERVER
)

// Bootstrap wrapper that shows splash, resolves license state, then swaps to real app
const Bootstrap = defineComponent({
  setup() {
    const resolvedApp = ref<any>(null)
    const resolvedProps = ref<Record<string, any>>({})
    const splashMessage = ref('Checking license...')

    onMounted(async () => {
      // Wait for the browser to paint the splash before blocking with sendSync
      await new Promise((resolve) =>
        requestAnimationFrame(() => requestAnimationFrame(resolve))
      )

      const activated = window.mainApi?.sendSync('msgGetAppStatus')

      if (!activated) {
        resolvedApp.value = AppSetup
        return
      }

      splashMessage.value = 'Validating license...'
      const validLicense = window.mainApi?.sendSync('msgCheckLicenseValidity')

      switch (validLicense) {
        case 'NO_INTERNET':
        case 'SERVER_UNREACHABLE':
          resolvedApp.value = AppNoInternet
          resolvedProps.value = {
            serverUnreachable: validLicense === 'SERVER_UNREACHABLE'
          }
          return
        case 'VALID':
          break
        case 'SUSPENDED':
        case 'EXPIRED':
        case 'OVERDUE':
        case 'INVALID':
          resolvedApp.value = LicenseWindow
          resolvedProps.value = { licenseStatus: validLicense }
          return
        default:
          resolvedApp.value = LicenseWindow
          resolvedProps.value = { licenseStatus: 'INVALID' }
          return
      }

      // License valid — check authentication
      splashMessage.value = 'Loading...'
      const authenticatedUser = window.mainApi?.sendSync(
        'msgGetAuthenticatedUser'
      )
      resolvedApp.value = authenticatedUser ? App : AppLogin
    })

    return () => {
      if (!resolvedApp.value) {
        return h(AppSplash, { message: splashMessage.value })
      }
      return h(resolvedApp.value, resolvedProps.value)
    }
  }
})

const app = createApp(Bootstrap)

const comps = import.meta.glob('./components/**/*.vue')
for (const path in comps) {
  let name = upperFirst(
    camelCase(path.replace(/^\.\//, '').replace(/\.\w+$/, ''))
  )
  name = name
    .replace(/Components/g, '')
    .replace(/^\.\//, '')
    .replace(/\.\w+$/, '')
    .replace(/([a-z0-9]|(?=[A-Z]))([A-Z])/g, '$1-$2')
    .replace(/^-/, '')
    .toLowerCase()
  app.component(name, comps[path])
}

app.use(createPinia()).use(vuetify).use(Vuelidate).use(i18n).use(router)

app.mount('#app')
