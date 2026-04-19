import axios from 'axios'

const baseUrl = window.mainApi.sendSync('msgGetBaseUrl')
const accessToken = window.mainApi.sendSync('msgGetAccessToken')

const instance = axios.create({
  baseURL: baseUrl,
  withCredentials: false,
  timeout: 300000,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
    'ngrok-skip-browser-warning': 'true'
  }
})

instance.interceptors.request.use(
  (config) => {
    if (accessToken) {
      config.headers.Authorization = `Bearer ${accessToken}`
    }

    // Cache-busting for GETs:
    //
    // Without this, Chromium (the Electron renderer) can serve previously
    // cached responses for plain JSON GETs after a recalculation — the user
    // sees stale numbers in the result / output / premium / reinsurance
    // summary screens until they force-reload (Ctrl+Shift+R). The API also
    // sets Cache-Control: no-store on these responses, but appending a
    // unique timestamp to the URL guarantees a fresh request even when the
    // header path is broken (intermediary proxy, stale binary, service
    // worker, etc.). Mutating verbs are skipped — they aren't cached and
    // some endpoints look at exact query strings (e.g. signed URLs).
    if ((config.method || 'get').toLowerCase() === 'get') {
      config.params = { ...(config.params || {}), _t: Date.now() }
      // Belt-and-braces: also tell any in-process cache (e.g. service
      // workers) we want a fresh copy.
      config.headers = {
        ...(config.headers || {}),
        'Cache-Control': 'no-cache, no-store',
        Pragma: 'no-cache'
      }
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

instance.interceptors.response.use(
  (response) => {
    return response
  },
  function (error) {
    return Promise.reject(error.response)
  }
)

export default instance
