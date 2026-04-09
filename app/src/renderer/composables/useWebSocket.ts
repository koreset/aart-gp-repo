import { ref } from 'vue'

type WSHandler = (payload: any) => void

let wsInstance: WebSocket | null = null
const isConnected = ref(false)
let reconnectAttempts = 0
const maxReconnectAttempts = 10
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
const handlers = new Map<string, Set<WSHandler>>()

function getWsUrl(): string {
  const baseUrl: string = window.mainApi.sendSync('msgGetBaseUrl') || ''
  // Convert http(s) to ws(s) and strip trailing slash to avoid double-slash
  const wsBase = baseUrl.replace(/^http/, 'ws').replace(/\/+$/, '')
  const token: string = window.mainApi.sendSync('msgGetAccessToken') || ''
  return `${wsBase}/ws?token=${encodeURIComponent(token)}`
}

function connect() {
  if (
    wsInstance &&
    (wsInstance.readyState === WebSocket.OPEN ||
      wsInstance.readyState === WebSocket.CONNECTING)
  ) {
    return
  }

  try {
    const url = getWsUrl()
    wsInstance = new WebSocket(url)

    wsInstance.onopen = () => {
      isConnected.value = true
      reconnectAttempts = 0
      console.log('[WS] Connected')
    }

    wsInstance.onclose = () => {
      isConnected.value = false
      console.log('[WS] Disconnected')
      scheduleReconnect()
    }

    wsInstance.onerror = (err) => {
      console.error('[WS] Error:', err)
    }

    wsInstance.onmessage = (event) => {
      try {
        const envelope = JSON.parse(event.data)
        const type = envelope.type as string
        const payload = envelope.payload
        const typeHandlers = handlers.get(type)
        if (typeHandlers) {
          typeHandlers.forEach((handler) => handler(payload))
        }
      } catch (e) {
        console.error('[WS] Failed to parse message:', e)
      }
    }
  } catch (e) {
    console.error('[WS] Failed to connect:', e)
    scheduleReconnect()
  }
}

function disconnect() {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (wsInstance) {
    wsInstance.onclose = null // prevent reconnect
    wsInstance.close()
    wsInstance = null
  }
  isConnected.value = false
}

function scheduleReconnect() {
  if (reconnectAttempts >= maxReconnectAttempts) {
    console.log('[WS] Max reconnect attempts reached')
    return
  }
  const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000)
  reconnectAttempts++
  console.log(`[WS] Reconnecting in ${delay}ms (attempt ${reconnectAttempts})`)
  reconnectTimer = setTimeout(connect, delay)
}

function on(type: string, handler: WSHandler) {
  if (!handlers.has(type)) {
    handlers.set(type, new Set())
  }
  handlers.get(type)!.add(handler)
}

function off(type: string, handler: WSHandler) {
  handlers.get(type)?.delete(handler)
}

function send(type: string, payload: any) {
  if (wsInstance && wsInstance.readyState === WebSocket.OPEN) {
    wsInstance.send(JSON.stringify({ type, payload }))
  }
}

// Reconnect on network recovery (Electron fires the 'online' event)
if (typeof window !== 'undefined') {
  window.addEventListener('online', () => {
    if (!isConnected.value) {
      reconnectAttempts = 0
      connect()
    }
  })
}

export function useWebSocket() {
  return { connect, disconnect, on, off, send, isConnected }
}
