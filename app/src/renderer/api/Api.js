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
