import { createVuetify } from 'vuetify'
import {
  ko,
  en,
  zhHans,
  zhHant,
  de,
  es,
  ja,
  fr,
  ru,
  pt,
  nl
} from 'vuetify/locale'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.min.css'
import colors from 'vuetify/util/colors'

export default createVuetify({
  locale: {
    messages: { ko, en, zhHans, zhHant, de, es, ja, fr, ru, pt, nl },
    locale: 'en',
    fallback: 'en'
  },
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: {
      mdi
    }
  },
  theme: {
    themes: {
      light: {
        dark: false,
        colors: {
          // Slate-neutral with sharper accents. The bright Material
          // defaults (info / warning / error / success) shipped by
          // Vuetify were drowning out content on every screen that used
          // `bg-info` / `bg-warning` etc. Anchored on the existing
          // corporate dark teal — the rest are tuned to sit beside it
          // without competing.
          primary: '#003F58', // existing corporate dark teal
          secondary: '#64748B', // slate — neutral chrome
          accent: '#006C8C', // lighter sibling of primary
          info: '#4338CA', // indigo — informational headers
          success: '#059669', // emerald — positive states
          warning: '#D97706', // amber — caution, not alarm
          error: '#DC2626' // red — destructive
        }
      },
      dark: {
        dark: true,
        colors: {
          primary: colors.green.darken4
        }
      }
    }
  }
})
